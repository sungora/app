package workfour

import (
	"fmt"
	"time"

	"github.com/sungora/app/workflow"
)

// Пример задачи работающей по расписанию
type SampleTaskFour struct {
}

// Manager режим работы задачи
func (task *SampleTaskFour) Manager() workflow.Manager {
	return workflow.Manager{
		Name:      "SampleTaskFour",
		IsExecute: true,
		Minute:    "*/4",
		Hour:      "*",
		Day:       "*",
		Month:     "*",
		Week:      "*",
	}
}

// Execute выполняемая задача
func (task *SampleTaskFour) Execute() {
	fmt.Println(time.Now().Format(time.RFC3339) + " execute: " + task.Manager().Name)
}
