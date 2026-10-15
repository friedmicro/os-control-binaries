package main

import (
	"fmt"
	"os"
	"time"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.OpenNicPortNonPrivilege(); err != nil {
		fmt.Fprintln(os.Stderr, "start-linux-application-stream: open-nic-port:", err)
		os.Exit(1)
	}
	_ = osctl.AthenaStop(osctl.WindowsVM)
	time.Sleep(osctl.StreamShutdownWait)
	if err := osctl.Webhook(osctl.WebhookWake); err != nil {
		fmt.Fprintln(os.Stderr, "start-linux-application-stream: webhook:", err)
		os.Exit(1)
	}
	if err := osctl.WaitForPing(osctl.ThickClientHost, osctl.WaitForPingTimeout); err != nil {
		fmt.Fprintln(os.Stderr, "start-linux-application-stream:", err)
		os.Exit(1)
	}
}
