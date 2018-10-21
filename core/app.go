// Управление запуском и остановкой приложения
package core

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"gopkg.in/sungora/app.v1/core/web"
	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/workflow"
)

// Каналы управления запуском и остановкой приложения
var (
	chanelAppStop    = make(chan os.Signal, 1)
	chanelAppControl = make(chan os.Signal, 1)
)

// Start Launch an application
func Start() (code int) {
	defer func() { // контроль завершение работы приложения
		chanelAppStop <- os.Interrupt
	}()
	var (
		err    error
		mode   string
		store  net.Listener
		config *conf.Config
	)

	mode, err = conf.GetCmdArgs() // входные данные командной строки
	if err != nil {
		return 0
	}

	// configuration
	if config, err = conf.GetConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// logs
	if err = lg.Start(config.Log, config.NameApp, config.TimeZone); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer lg.Wait()

	// web - base controller
	if err = web.Start(config); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer web.Wait()

	// workflow
	if config.Isworkflow == true {
		if err = workflow.Start(config.Workflow); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer workflow.Wait()
	}

	// web server - application
	if store, err = newWeb(config.Server); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer store.Close()
	fmt.Fprintln(os.Stdout, "web app start success")

	// auto deploy
	// go utils.Deploy(config.NameApp)

	// The correctness of the application is closed by a signal
	if mode == "" || mode == "run" {
		signal.Notify(chanelAppControl, os.Interrupt)
		<-chanelAppControl
	}
	return
}

// stop Stop an application
func Stop() {
	chanelAppControl <- os.Interrupt
	<-chanelAppStop
}
