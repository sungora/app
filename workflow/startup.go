package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/lg"
	"github.com/sungora/app/startup"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	startup.SetComponent(component)
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
func (comp *componentTyp) Init() (err error) {
	sep := string(os.PathSeparator)
	config = new(configMain)

	// техническое имя приложения
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), ext)
		config.ServiceName = sl[0]
	} else {
		config.ServiceName = filepath.Base(os.Args[0])
	}

	// читаем конфигурацию
	dirWork, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := dirWork + sep + "config" + sep + config.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, &config); err != nil {
		return
	}

	// читаем задачи из конфигурации
	path = dirWork + sep + "config" + sep + config.ServiceName + "_workflow.toml"
	if _, err := toml.DecodeFile(path, &comp.cronTaskManager); err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
	}

	comp.cronTaskRun = make(map[string]Task)
	comp.cronControlCH = make(chan struct{})
	comp.p = NewPool(config.Workflow.LimitCh, config.Workflow.LimitPool)

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
