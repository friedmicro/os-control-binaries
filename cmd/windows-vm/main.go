package main

import (
	"fmt"
	"os"
	"time"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	_ = osctl.AthenaRemoveTracking("lair.friedmicro-lab.org")
	if err := osctl.WindowsHost(); err != nil {
		fmt.Fprintln(os.Stderr, "windows-vm: windows-host:", err)
		os.Exit(1)
	}
	if err := osctl.SSHExec(osctl.ThickClientHost, "qm start 101"); err != nil {
		fmt.Fprintln(os.Stderr, "windows-vm:", err)
		os.Exit(1)
	}
	time.Sleep(osctl.WindowsVMStartWait)
}
