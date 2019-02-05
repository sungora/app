package workflow

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	core.ComponentReg(component)
}

// компонент
type componentTyp struct {
	p               *pool
	cronTaskManager map[string]*manager
	cronTaskRun     map[string]Task
	cronControlCH   chan struct{}
}

var (
	config    *configMain   // конфигурация
	component *componentTyp // компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *core.ConfigRoot) (err error) {
	sep := string(os.PathSeparator)
	config = new(configMain)

	// читаем конфигурацию
	path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, config); err != nil {
		return
	}

	// читаем задачи из конфигурации
	path = cfg.DirConfig + sep + cfg.ServiceName + "_workflow.toml"
	if _, err := toml.DecodeFile(path, &comp.cronTaskManager); err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
	}

	comp.cronTaskRun = make(map[string]Task)

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	var (
		t           time.Time
		index       string
		taskManager *manager
		task        Task
		ok          bool
	)
	comp.cronControlCH = make(chan struct{})
	comp.p = NewPool(config.Workflow.LimitCh, config.Workflow.LimitPool)
	comp.p.wg.Add(1)
	go func() {
		defer comp.p.wg.Done()
		for {
			// таймаут
			select {
			case <-comp.cronControlCH:
				return
			case <-time.After(time.Minute):
				t = time.Now()
			}
			//
			for index, taskManager = range comp.cronTaskManager {
				task, ok = comp.cronTaskRun[index]
				if ok == false {
					lg.Error("not found cron task [%s]", index)
					continue
				}
				if taskManager.IsExecute == false {
					continue
				}
				if checkRuntime(t.Minute(), taskManager.Minute) == false {
					continue
				}
				if checkRuntime(t.Hour(), taskManager.Hour) == false {
					continue
				}
				if checkRuntime(t.Day(), taskManager.Day) == false {
					continue
				}
				if checkRuntime(int(t.Month()), taskManager.Month) == false {
					continue
				}
				if checkRuntime(int(t.Weekday()), taskManager.Week) == false {
					continue
				}
				TaskAdd(task)
			}
		}
	}()
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	comp.cronControlCH <- struct{}{}
	time.Sleep(time.Second)
	comp.p.Wait()
	return
}
