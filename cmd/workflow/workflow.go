package main

import (
	"fmt"

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/workflow"
)

type ExampleTask string

func (e ExampleTask) Execute() {
	fmt.Println("executing:", string(e))
}


type SampleCron struct {
}

func (self *SampleCron) Execute() {
	fmt.Println("RUN SampleCron")
}


func main() {

	// custom
	fmt.Println("custom")
	pool := workflow.NewPool(200, 9)
	pool.TaskAdd(ExampleTask("foo"))
	pool.TaskAdd(ExampleTask("bar"))
	for i := 0; i < 20; i++ {
		pool.TaskAdd(ExampleTask(fmt.Sprintf("additional_%d", i+1)))
	}
	pool.Wait()

	// system
	fmt.Println("\nsystem")

	// cron task
	workflow.TaskAddCron("SampleCron", new(SampleCron))

	c := conf.Workflow{
		LimitCh:   200,
		LimitPool: 9,
	}
	workflow.Start(c)
	workflow.TaskAdd(ExampleTask("foo"))
	workflow.TaskAdd(ExampleTask("bar"))
	for i := 0; i < 20; i++ {
		workflow.TaskAdd(ExampleTask(fmt.Sprintf("additional_%d", i+1)))
	}
	workflow.Wait()
}
