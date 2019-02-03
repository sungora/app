package main

import (
	"fmt"
	"os"

	_ "github.com/sungora/app/cmd/workflow/worker"
	"github.com/sungora/app/startup"
	"github.com/sungora/app/workflow"
)

func main() {
	// custom
	pool := workflow.NewPool(200, 9)
	pool.TaskAdd(SampleTask("foo"))
	pool.TaskAdd(SampleTask("bar"))
	for i := 0; i < 20; i++ {
		pool.TaskAdd(SampleTask(fmt.Sprintf("additional_%d", i+1)))
	}
	pool.Wait()
	// inside program
	// workflow.TaskAdd(SampleTask("foo"))
	// workflow.TaskAdd(SampleTask("bar"))
	os.Exit(startup.Start())
}

// Пример задачи
type SampleTask string

func (e SampleTask) Execute() {
	fmt.Println("execute: ", string(e))
}
