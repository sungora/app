package worktwo

import (
	"fmt"
	"time"
)

// Пример задачи работающей по расписанию
type SampleTaskTwo struct {
}

func (self *SampleTaskTwo) Execute() {
	fmt.Println(time.Now().Format(time.RFC3339) + " execute: SampleTaskTwo")
}
