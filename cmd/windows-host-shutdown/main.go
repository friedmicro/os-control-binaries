package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.WindowsHostShutdown(); err != nil {
		fmt.Fprintln(os.Stderr, "windows-host-shutdown:", err)
		os.Exit(1)
	}
}
