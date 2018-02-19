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
	r := runner.Runner{Ctx: ctx}

	select {
	case <-r.Run(func() { runSomething(ctx, 1, 2, 3) }):
		fmt.Println("Exited runSomething")
	case sig := <-r.RunSigs(syscall.SIGINT, syscall.SIGTERM):
		fmt.Println("Received signal", sig)
	}

	cancel()
	r.Wait()
}
