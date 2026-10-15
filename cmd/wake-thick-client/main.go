package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.WakeThickClient(); err != nil {
		fmt.Fprintln(os.Stderr, "wake-thick-client:", err)
		os.Exit(1)
	}
}
