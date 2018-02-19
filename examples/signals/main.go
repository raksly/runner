package main

import (
	"context"
	"fmt"
	"syscall"

	"github.com/raksly/runner"
)

func runTCP(ctx context.Context) {
	fmt.Println("Entering runTCP")
	<-ctx.Done()
	fmt.Println("Exiting runTCP")
}

func runHTTP(ctx context.Context) {
	fmt.Println("Entering runHTTP")
	<-ctx.Done()
	fmt.Println("Exiting runHTTP")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	runner := runner.New(ctx)

	select {
	case <-runner.RunContext(runHTTP):
		fmt.Println("Exited runHTTP")
	case <-runner.RunContext(runTCP):
		fmt.Println("Exited runTCP")
	case sig := <-runner.RunSigs(syscall.SIGINT, syscall.SIGTERM):
		fmt.Println("Received signal", sig)
	}

	cancel()
	runner.Wait()
}
