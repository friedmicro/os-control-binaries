package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.ShutdownThickClient(); err != nil {
		fmt.Fprintln(os.Stderr, "shutdown-thick-client:", err)
		os.Exit(1)
	}
}
