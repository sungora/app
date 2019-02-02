package setup

import (
	"accounter/workers/tariffhourly"

	"github.com/hostkeybv/app/workflow"
)

func setupWorker() {

	workflow.TaskAddCron("TariffHourly", &tariffhourly.TaskTariffHourly{})

}
