package osctl

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	ThickClientHost = "192.168.1.8"
	WindowsVM       = "192.168.1.27"
)

const (
	WebhookVR       = "http://serpent.friedmicro-lab.org:8123/api/webhook/-S8d5lmJYlgGHhEzQTe_B6j_2"
	WebhookShutdown = "http://serpent.friedmicro-lab.org:8123/api/webhook/-6MtudTLIw34yJxCMYknhw9Iv"
	WebhookWake     = "http://serpent.friedmicro-lab.org:8123/api/webhook/-D-oaFaHosZ8BDzkbnwNAYDls"
)

const (
	WaitForPingTimeout   = 5 * time.Minute
	PollOSReleaseTimeout = 10 * time.Minute
	WaitForNoVMSTimeout  = 6 * time.Minute
	RebootSleep          = 20 * time.Second
	ShutdownSnapshotWait = 60 * time.Second
	PostShutdownWait     = 20 * time.Second
	SnapshotPingWait     = 5 * time.Second
	StreamShutdownWait   = 20 * time.Second
	WindowsVMStartWait   = 60 * time.Second
)

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	if err := c.Run(); err != nil {
		return fmt.Errorf("%s %s: %w: %s", cmd, strings.Join(args, " "), err, strings.TrimSpace(out.String()))
	}
	return nil
}

func runOutput(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	if err := c.Run(); err != nil {
		return "", fmt.Errorf("%s %s: %w: %s", cmd, strings.Join(args, " "), err, strings.TrimSpace(out.String()))
	}
	return out.String(), nil
}

func privilegedRun(bin string, args ...string) error {
	if os.Geteuid() == 0 {
		return run(bin, args...)
	}
	return run("sudo", append([]string{bin}, args...)...)
}

func sshArgs(host string) []string {
	return []string{
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"root@" + host,
	}
}

func SSHExec(host, command string) error {
	return run("ssh", append(sshArgs(host), command)...)
}

func SSHOutput(host, command string) (string, error) {
	return runOutput("ssh", append(sshArgs(host), command)...)
}

func PingHost(host string) bool {
	return run("ping", "-c", "1", "-W", "1", host) == nil
}

func WaitForPing(host string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if PingHost(host) {
			return nil
		}
		time.Sleep(time.Second)
	}
	return fmt.Errorf("host %s did not answer ping within %s", host, timeout)
}

func Webhook(url string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("webhook %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook %s: unexpected status %s", url, resp.Status)
	}
	return nil
}

func CloseNicPort() error {
	if err := privilegedRun("ip", "link", "set", "ens18", "up"); err != nil {
		return err
	}
	if err := privilegedRun("ip", "link", "set", "ens27", "down"); err != nil {
		return err
	}
	return privilegedRun("systemctl", "start", "ci")
}

func OpenNicPortNonPrivilege() error {
	if err := run("sudo", "open-nic-port"); err != nil {
		return err
	}
	return nil
}

func OpenNicPort() error {
	if err := privilegedRun("ip", "link", "set", "ens18", "down"); err != nil {
		return err
	}
	if err := privilegedRun("ip", "link", "set", "ens27", "up"); err != nil {
		return err
	}
	return privilegedRun("systemctl", "start", "ci")
}

func AthenaStop(host string) error {
	return run("athena", "remote", "stop", host)
}

func AthenaRemoveTracking(host string) error {
	return run("athena", "remote", "remove_tracking", host)
}

func PollOSRelease(host string, matches []string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if out, err := SSHOutput(host, "cat /etc/os-release"); err == nil {
			for _, m := range matches {
				if strings.Contains(out, m) {
					return nil
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("host %s os-release did not contain %v within %s", host, matches, timeout)
}

func WaitForNoVMs(host string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		out, err := SSHOutput(host, "qm list")
		if err == nil && !hasRunningVM(out) {
			return nil
		}
		if time.Now().After(deadline) {
			fmt.Fprintf(os.Stderr, "osctl: VMs still running on %s after %s, proceeding anyway\n", host, timeout)
			return nil
		}
		time.Sleep(5 * time.Second)
	}
}

func hasRunningVM(qmList string) bool {
	sc := bufio.NewScanner(strings.NewReader(qmList))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 3 {
			continue
		}
		if fields[0] == "VMID" {
			continue
		}
		if fields[2] != "stopped" {
			return true
		}
	}
	return false
}

func RebootIfNotArch(host string) error {
	out, err := SSHOutput(host, "cat /etc/os-release")
	if err != nil {
		return err
	}
	if strings.Contains(out, "Arch") {
		return nil
	}
	return RemoteReboot(host)
}

func TakeSnapshot(host string) error {
	return SSHExec(host, "take-snapshot")
}

func RemoteShutdown(host string) error {
	return SSHExec(host, "/usr/bin/shutdown 0")
}

func RemoteReboot(host string) error {
	return SSHExec(host, "/usr/sbin/reboot")
}

func WakeThickClient() error {
	_ = AthenaStop(WindowsVM)
	if PingHost(ThickClientHost) {
		_ = RebootIfNotArch(ThickClientHost)
		time.Sleep(RebootSleep)
	}
	if err := Webhook(WebhookWake); err != nil {
		return err
	}
	return WaitForPing(ThickClientHost, WaitForPingTimeout)
}

func ShutdownThickClient() error {
	if err := run("sudo", "close-nic-port"); err != nil {
		return err
	}
	time.Sleep(SnapshotPingWait)
	if err := WaitForPing(ThickClientHost, WaitForPingTimeout); err != nil {
		return err
	}
	if err := TakeSnapshot(ThickClientHost); err != nil {
		return err
	}
	if err := RemoteShutdown(ThickClientHost); err != nil {
		return err
	}
	time.Sleep(ShutdownSnapshotWait)
	if err := Webhook(WebhookShutdown); err != nil {
		return err
	}
	time.Sleep(PostShutdownWait)
	return nil
}

func WindowsHost() error {
	if err := WakeThickClient(); err != nil {
		return err
	}
	if err := SSHExec(ThickClientHost, "windows-host"); err != nil {
		return err
	}
	return PollOSRelease(ThickClientHost, []string{"Debian", "Proxmox"}, PollOSReleaseTimeout)
}

func WindowsHostShutdown() error {
	if err := WaitForNoVMs(ThickClientHost, WaitForNoVMSTimeout); err != nil {
		return err
	}
	if err := RemoteReboot(ThickClientHost); err != nil {
		return err
	}
	if err := PollOSRelease(ThickClientHost, []string{"Arch"}, PollOSReleaseTimeout); err != nil {
		return err
	}
	return ShutdownThickClient()
}
