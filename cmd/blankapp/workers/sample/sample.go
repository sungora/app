package sample

import (
	"fmt"
	"time"

	"gopkg.in/sungora/app.v1/conf"
)

type TaskSample struct {
}

// Execute Задача
func (self *TaskSample) Execute() {

	timeLabel := time.Now().In(conf.TimeLocation).Format("2006-01-02 15:04:05")
	fmt.Println(timeLabel + " Sample Task")

}
