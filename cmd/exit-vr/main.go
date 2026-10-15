package main

import (
	"fmt"
	"os"
	"sync"

	"git.friedmicro-lab.org/Homelab_Infra_Public/os-remote-control/internal/osctl"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go baseStateToggle(&wg)
	go markOutOfVR(&wg)
	wg.Wait()

	if err := osctl.ShutdownThickClient(); err != nil {
		fmt.Fprintln(os.Stderr, "exit-vr:", err)
		os.Exit(1)
	}
}

func baseStateToggle(wg *sync.WaitGroup) {
	if err := osctl.Webhook(osctl.WebhookVR); err != nil {
		fmt.Fprintln(os.Stderr, "exit-vr: webhook:", err)
	}
	defer wg.Done()
}

func markOutOfVR(wg *sync.WaitGroup) {
	if err := os.Remove("/tmp/windows-bare-metal"); err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "exit-vr: remove marker:", err)
	}
	defer wg.Done()
}
