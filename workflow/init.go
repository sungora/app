package workflow

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/core"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	core.ComponentReg(component)
}

var (
	config    *configFile   // конфигурация
	component *componentTyp // компонент
)

// компонент
type componentTyp struct {
	p             *pool
	cronTaskRun   []Task
	cronControlCH chan struct{}
}

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *core.ConfigRoot) (err error) {
	sep := string(os.PathSeparator)
	config = new(configFile)

	// читаем конфигурацию
	path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, config); err != nil {
		return
	}

	// читаем задачи из конфигурации
	// cronTaskManager map[string]*Manager
	// path = cfg.DirConfig + sep + cfg.ServiceName + "_workflow.toml"
	// if _, err := toml.DecodeFile(path, &cronTaskManager); err != nil {
	// 	fmt.Fprintln(os.Stdout, err.Error())
	// }
	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	var (
		t           time.Time
		taskManager Manager
		task        Task
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
			for _, task = range comp.cronTaskRun {
				taskManager = task.Manager()
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
	comp.p.Wait()
	return
}
