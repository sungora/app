package worker

import (
	"PKGAPPNAME/workers/sample"

	"gopkg.in/sungora/app.v1/workflow"
)

func init() {
	workflow.TaskAddCron("TaskSample", &sample.TaskSample{})
}
