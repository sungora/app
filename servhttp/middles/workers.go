package middles

import (
	"net/http"

	"github.com/sungora/app"
	"github.com/sungora/app/workflow"
)

// Логирование выполение запроса
type TaskLogRequest struct {
	Request *http.Request
	Status  int
}

// Manager режим работы задачи
func (task *TaskLogRequest) Manager() workflow.Manager {
	return workflow.Manager{
		Name:      "TaskLogRequest",
		IsExecute: true,
	}
}

// Execute выполняемая задача
func (task *TaskLogRequest) Execute() {
	app.Dumper(task)
}
