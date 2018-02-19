package main

import (
	"context"
	"fmt"
	"time"

	"github.com/raksly/runner"
)

func runTCP() {
	fmt.Println("Entering runTCP")
	time.Sleep(time.Second)
	fmt.Println("Exiting runTCP")
}

func runHTTP() {
	fmt.Println("Entering runTCP")
	time.Sleep(time.Second)
	fmt.Println("Exiting runTCP")
}

func main() {
	runner := runner.New(context.Background())

	runner.Run(runHTTP)
	runner.Run(runTCP)

	runner.Wait()
}
