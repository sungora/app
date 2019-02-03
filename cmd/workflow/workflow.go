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
	pool.TaskAdd(&SampleTask{Sample: "foo"})
	pool.TaskAdd(&SampleTask{Sample: "bar"})
	for i := 0; i < 20; i++ {
		pool.TaskAdd(&SampleTask{Sample: fmt.Sprintf("additional_%d", i+1)})
	}
	pool.Wait()
	// inside program
	// workflow.TaskAdd(SampleTask("foo"))
	// workflow.TaskAdd(SampleTask("bar"))
	os.Exit(startup.Start())
}

// Пример задачи
type SampleTask struct {
	Sample string
}

func (e *SampleTask) Execute() {
	fmt.Println("execute: ", e.Sample)
}
