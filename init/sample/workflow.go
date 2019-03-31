package sample

import (
	"github.com/sungora/app/internal/sample/worker/workfour"
	"github.com/sungora/app/internal/sample/worker/workone"
	"github.com/sungora/app/internal/sample/worker/worktwo"
	"github.com/sungora/app/workflow"
)

func workers() {
	workflow.TaskAddCron(&workone.SampleTaskOne{})
	workflow.TaskAddCron(&worktwo.SampleTaskTwo{})
	workflow.TaskAddCron(&workfour.SampleTaskFour{})
}
