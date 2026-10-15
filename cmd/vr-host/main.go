package main

import (
	"fmt"
	"os"
	"sync"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	go baseStationToggle(&wg)
	go markInVR(&wg)
	go wakeThickClient(&wg)
	wg.Wait()

	if err := osctl.SSHExec(osctl.ThickClientHost, "vr-host"); err != nil {
		fmt.Fprintln(os.Stderr, "vr-host:", err)
		os.Exit(1)
	}
}

func baseStationToggle(wg *sync.WaitGroup) {
	if err := osctl.Webhook(osctl.WebhookVR); err != nil {
		fmt.Fprintln(os.Stderr, "vr-host: webhook:", err)
	}
	defer wg.Done()
}

func markInVR(wg *sync.WaitGroup) {
	if err := os.WriteFile("/tmp/windows-bare-metal", nil, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "vr-host: touch marker:", err)
	}
	defer wg.Done()
}

func wakeThickClient(wg *sync.WaitGroup) {
	if err := osctl.WakeThickClient(); err != nil {
		fmt.Fprintln(os.Stderr, "vr-host: wake-thick-client:", err)
	}
	defer wg.Done()
}
