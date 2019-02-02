package worker

import (
	"github.com/sungora/app/startup"

	"github.com/sungora/app/cmd/workflow/worker/workfour"
	"github.com/sungora/app/cmd/workflow/worker/workone"
	"github.com/sungora/app/cmd/workflow/worker/worktwo"
	"github.com/sungora/app/workflow"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	startup.SetComponent(component)
}

// компонент
type componentTyp struct {
}

var (
	component *componentTyp // компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init() (err error) {

	workflow.TaskAddCron("SampleTaskOne", &workone.SampleTaskOne{})
	workflow.TaskAddCron("SampleTaskTwo", &worktwo.SampleTaskTwo{})
	workflow.TaskAddCron("SampleTaskFour", &workfour.SampleTaskFour{})

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	return
}
