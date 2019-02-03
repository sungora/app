package workone

import (
	"fmt"
	"time"
)

// Пример задачи работающей по расписанию
type SampleTaskOne struct {
}

func (self *SampleTaskOne) Execute() {
	fmt.Println(time.Now().Format(time.RFC3339) + " execute: SampleTaskOne")
}
