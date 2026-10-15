package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.SSHExec(osctl.ThickClientHost, "qm stop 101"); err != nil {
		fmt.Fprintln(os.Stderr, "windows-vm-shutdown:", err)
		os.Exit(1)
	}
	if err := osctl.WindowsHostShutdown(); err != nil {
		fmt.Fprintln(os.Stderr, "windows-vm-shutdown: windows-host-shutdown:", err)
		os.Exit(1)
	}
}
