package workfour

import (
	"fmt"
	"time"
)

// Пример задачи работающей по расписанию
type SampleTaskFour struct {
}

func (self *SampleTaskFour) Execute() {
	fmt.Println(time.Now().Format(time.RFC3339) + " execute: SampleTaskFour")
}
