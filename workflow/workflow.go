// Пул обработчиков для паралельных задач
package workflow

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/tool"
)

type manager struct {
	Name      string
	IsExecute bool
	Minute    string
	Hour      string
	Day       string
	Month     string
	Week      string
}

type Config struct {
	IsWorkflow bool
	LimitCh    int // Лимит канала задач
	LimitPool  int // Лимит пула (количество воркеров)
}

var p *pool
var cronTaskManager = make(map[string]*manager)
var cronTaskRun = make(map[string]Task)
var cronControlCH chan struct{}

// Start Создаем пул воркеров указанного размера на уровне пакета
func Start(c Config) (err error) {
	var cronTaskPath string
	if cronTaskPath, err = os.Getwd(); err != nil {
		return
	}
	sep := string(os.PathSeparator)
	cronTaskPath = cronTaskPath + sep + "config" + sep + "cron.toml"

	p = NewPool(c.LimitCh, c.LimitPool)

	if err = reloadTasks(cronTaskPath); err != nil {
		return
	}

	cronControlCH = make(chan struct{})
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			// таймаут минуту
			for i := 0; i < 60; i += 5 {
				select {
				case <-cronControlCH:
					return
				case <-time.After(5 * time.Second):
				}
			}
			// обновляем задачи
			reloadTasks(cronTaskPath)
			//
			minute := time.Now().Minute()
			hour := time.Now().Hour()
			day := time.Now().Day()
			month := int(time.Now().Month())
			week := int(time.Now().Weekday())
			for index, t := range cronTaskManager {
				if t.IsExecute == false {
					continue
				}
				if checkRuntime(minute, t.Minute) == false {
					continue
				}
				if checkRuntime(hour, t.Hour) == false {
					continue
				}
				if checkRuntime(day, t.Day) == false {
					continue
				}
				if checkRuntime(month, t.Month) == false {
					continue
				}
				if checkRuntime(week, t.Week) == false {
					continue
				}
				if task, ok := cronTaskRun[index]; ok {
					TaskAdd(task)
				} else {
					lg.Error("not found cron task [%s]", index)
				}
			}
		}
	}()
	return
}

// TaskAdd Добавляем задачу в пул
func TaskAdd(task Task) {
	p.tasks <- task
}

func TaskAddCron(name string, task Task) {
	cronTaskRun[name] = task
}

// Wait Завершаем работу пула
func Wait() {
	cronControlCH <- struct{}{}
	p.Wait()
}

var controlTask = tool.NewControlFS()

// ReadToml Функция читает конфигурационный файл в формате toml. Отдельный конфиг не связанный с beego.
func reloadTasks(cronTaskPath string) (err error) {
	var isChange bool
	if isChange, err = controlTask.CheckSumMd5(cronTaskPath, ""); err == nil && isChange {
		if _, err = toml.DecodeFile(cronTaskPath, &cronTaskManager); err != nil {
			return
		}
	}
	return
}

func checkRuntime(val int, mask string) bool {
	var number int
	var sl []string
	//  any valid value or exact match
	number, _ = strconv.Atoi(mask)
	if "*" == mask || val == number {
		return true
	}
	//  range
	sl = strings.Split(mask, "-")
	if 1 < len(sl) {
		number1, _ := strconv.Atoi(sl[0])
		number2, _ := strconv.Atoi(sl[1])
		if number1 <= val && val <= number2 {
			return true
		}
		return false
	}
	//  fold
	sl = strings.Split(mask, "/")
	if 1 < len(sl) {
		number, _ = strconv.Atoi(sl[1])
		if 0 < val%number {
			return false
		}
		return true
	}
	//  list
	sl = strings.Split(mask, ",")
	if 1 < len(sl) {
		for _, v := range sl {
			number, _ = strconv.Atoi(v)
			if number == val {
				return true
			}
		}
		return false
	}
	//
	return false
}
