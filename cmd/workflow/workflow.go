package main

import (
	"fmt"

	"github.com/sungora/app/workflow"
)

// Пример задачи
type SampleTask string

func (e SampleTask) Execute() {
	fmt.Println("execute: ", string(e))
}

// Пример задачи работающей по расписанию
type SampleTaskCron struct {
}

func (self *SampleTaskCron) Execute() {
	fmt.Println("execute: SampleTaskCron")
}

func main() {

	var err error

	// custom
	fmt.Println("Sample run task")
	pool := workflow.NewPool(200, 9)
	pool.TaskAdd(SampleTask("foo"))
	pool.TaskAdd(SampleTask("bar"))
	for i := 0; i < 20; i++ {
		pool.TaskAdd(SampleTask(fmt.Sprintf("additional_%d", i+1)))
	}
	pool.Wait()

	// cron task
	fmt.Println("\nSample run crontab task")

	workflow.TaskAddCron("SampleTaskCron", new(SampleTaskCron))

	c := workflow.Config{
		LimitCh:   200,
		LimitPool: 9,
	}
	if err = workflow.Init(c); err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}
	workflow.TaskAdd(SampleTask("foo"))
	workflow.TaskAdd(SampleTask("bar"))
	for i := 0; i < 20; i++ {
		workflow.TaskAdd(SampleTask(fmt.Sprintf("additional_%d", i+1)))
	}
	// time.Sleep(time.Minute * 3)
	workflow.Wait()
}
