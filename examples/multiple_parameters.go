package main

import (
	"context"
	"fmt"
	"syscall"

	"github.com/raksly/runner"
)

func runSomething(ctx context.Context, a, b, c int) {
	fmt.Println("Entering runSomething", a, b, c)
	<-ctx.Done()
	fmt.Println("Exiting runSomething")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	runner := runner.New(ctx)

	select {
	case <-runner.Run(func() { runSomething(ctx, 1, 2, 3) }):
		fmt.Println("Exited runSomething")
	case sig := <-runner.RunSigs(syscall.SIGINT, syscall.SIGTERM):
		fmt.Println("Received signal", sig)
		cancel()
	}

	runner.Wait()
}
