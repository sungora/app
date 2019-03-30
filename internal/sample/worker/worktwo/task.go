package worktwo

import (
	"fmt"
	"time"

	"github.com/sungora/app/workflow"
)

// Пример задачи работающей по расписанию
type SampleTaskTwo struct {
}

// Manager режим работы задачи
func (task *SampleTaskTwo) Manager() workflow.Manager {
	return workflow.Manager{
		Name:      "SampleTaskTwo",
		IsExecute: true,
		Minute:    "*/2",
		Hour:      "*",
		Day:       "*",
		Month:     "*",
		Week:      "*",
	}
}

// Execute выполняемая задача
func (task *SampleTaskTwo) Execute() {
	fmt.Println(time.Now().Format(time.RFC3339) + " execute: " + task.Manager().Name)
}
