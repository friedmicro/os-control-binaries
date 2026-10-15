package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.SSHExec(osctl.ThickClientHost, "qm stop 102"); err != nil {
		fmt.Fprintln(os.Stderr, "stop-dev-ai:", err)
		os.Exit(1)
	}
	if err := osctl.RemoteReboot(osctl.ThickClientHost); err != nil {
		fmt.Fprintln(os.Stderr, "stop-dev-ai: reboot:", err)
		os.Exit(1)
	}
	if err := osctl.PollOSRelease(osctl.ThickClientHost, []string{"Arch"}, osctl.PollOSReleaseTimeout); err != nil {
		fmt.Fprintln(os.Stderr, "stop-dev-ai:", err)
		os.Exit(1)
	}
	if err := osctl.ShutdownThickClient(); err != nil {
		fmt.Fprintln(os.Stderr, "stop-dev-ai: shutdown-thick-client:", err)
		os.Exit(1)
	}
}
