package main

import (
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
	var r runner.Runner

	r.Run(runHTTP)
	r.Run(runTCP)

	r.Wait()
}
