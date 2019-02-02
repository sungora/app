package setup

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/hostkeybv/app"
	"github.com/hostkeybv/app/lg"
	"github.com/hostkeybv/app/tool"
	"github.com/hostkeybv/app/workflow"

	"accounter/config"
)

type Config struct {
	App        app.Config
	Mysql      config.Mysql
	Postgresql config.Postgresql
	Log        lg.Config
	Workflow   workflow.Config
	BillingNL  config.Billing
	BillingRU  config.Billing
}

var (
	chanelAppControl = make(chan os.Signal, 1) // Канал управления остановкой приложения
)

// Run Launch an application
func Run() (code int) {
	defer func() { // контроль завершение работы приложения
		chanelAppControl <- os.Interrupt
	}()

	// config
	conf, err := configuration()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// logs
	setupLg()
	if err = lg.Init(conf.Log); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer lg.Wait()

	// model
	if err = initModel(conf); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer waitModel()

	// controller
	setupController()

	// workflow
	setupWorker()
	if err = workflow.Init(conf.Workflow); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer workflow.Wait()

	// service - application
	if err = app.Init(conf.App); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer app.Wait()

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
	<-chanelAppControl
	return
}

// Wait an application
func Wait() {
	chanelAppControl <- os.Interrupt
	<-chanelAppControl
}

func configuration() (c *Config, err error) {
	sep := string(os.PathSeparator)

	// техническое имя приложения
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), ext)
		tool.ServiceName = sl[0]
	} else {
		tool.ServiceName = filepath.Base(os.Args[0])
	}

	// рабочие директории
	tool.DirWork, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	// tool.DirWork = filepath.Dir(tool.DirWork)
	tool.DirLog = tool.DirWork + sep + "log"
	tool.DirWww = tool.DirWork + sep + "www"
	tool.DirTpl = tool.DirWork + sep + "tpl"
	tool.DirPid = tool.DirWork + sep + "pid"
	tool.DirConfig = tool.DirWork + sep + "conf"

	// конфигурация
	// var searchList []string
	// searchList = append(searchList, tool.DirWork+sep+"conf")
	// searchList = append(searchList, "/etc/"+tool.ServiceName)
	// for _, path := range searchList {
	// 	info, err := os.Stat(path)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	if info.IsDir() == false {
	// 		continue
	// 	}
	// 	tool.DirConfig = path
	// 	break
	// }
	path := tool.DirConfig + sep + tool.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, &c); err != nil {
		return
	}

	// Инициализация временной зоны
	if loc, err := time.LoadLocation(c.App.TimeZone); err == nil {
		tool.TimeLocation = loc
	} else {
		tool.TimeLocation = time.UTC
	}

	c.App.SessionTimeout *= time.Second

	// проектные конфигурации
	config.BillingNL = &c.BillingNL
	config.BillingRU = &c.BillingRU

	return
}
