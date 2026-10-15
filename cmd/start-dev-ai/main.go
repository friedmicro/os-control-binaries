package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.WindowsHost(); err != nil {
		fmt.Fprintln(os.Stderr, "start-dev-ai: windows-host:", err)
		os.Exit(1)
	}
	if err := osctl.SSHExec(osctl.ThickClientHost, "qm start 102"); err != nil {
		fmt.Fprintln(os.Stderr, "start-dev-ai:", err)
		os.Exit(1)
	}
	if err := osctl.CloseNicPort(); err != nil {
		fmt.Fprintln(os.Stderr, "start-dev-ai: close-nic-port:", err)
		os.Exit(1)
	}
}
