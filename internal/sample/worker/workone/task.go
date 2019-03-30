package workone

import (
	"fmt"
	"time"

	"github.com/sungora/app/workflow"
)

// Пример задачи работающей по расписанию
type SampleTaskOne struct {
}

// Manager режим работы задачи
func (task *SampleTaskOne) Manager() workflow.Manager {
	return workflow.Manager{
		Name:      "SampleTaskOne",
		IsExecute: true,
		Minute:    "*",
		Hour:      "*",
		Day:       "*",
		Month:     "*",
		Week:      "*",
	}
}

// Execute выполняемая задача
func (task *SampleTaskOne) Execute() {
	fmt.Println(time.Now().Format(time.RFC3339) + " execute: " + task.Manager().Name)
}
