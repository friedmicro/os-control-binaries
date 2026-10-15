package main

import (
	"fmt"
	"os"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	if err := osctl.OpenNicPort(); err != nil {
		fmt.Fprintln(os.Stderr, "open-nic-port:", err)
		os.Exit(1)
	}
}
