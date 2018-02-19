package main

import (
	"context"
	"fmt"
	"time"

	"github.com/raksly/runner"
)

func runTCP(ctx context.Context) {
	fmt.Println("Entering runTCP")
	time.Sleep(time.Second)
	fmt.Println("Exiting runTCP")
}

func runHTTP(ctx context.Context) {
	fmt.Println("Entering runHTTP")
	<-ctx.Done()
	fmt.Println("Exiting runHTTP")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	r := runner.Runner{Ctx: ctx}

	select {
	case <-r.RunContext(runHTTP):
		fmt.Println("Exited runHTTP")
	case <-r.RunContext(runTCP):
		fmt.Println("Exited runTCP")
	}

	cancel()
	r.Wait()
}
