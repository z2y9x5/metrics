package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/z2y9x5/metrics/internal/agent"
)

func main() {
	cnf := agent.NewConfig()
	cnf.ApplyCLIArgs()
	cnf.ApplyEnvArgs()
	cnfApp := cnf.GetAppConfig()

	db := agent.NewMemoryStorage()
	collector := agent.NewCollector(cnfApp, db)
	sender := agent.NewSender(cnfApp, db)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go collector.Run(ctx, &wg)
	wg.Add(1)
	go sender.Run(ctx, &wg)

	<-sigChan
	cancel()
	wg.Wait()
}
