// Управление запуском и остановкой приложения
package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"gopkg.in/sungora/app.v1/conf"
	"gopkg.in/sungora/app.v1/lg"
	"gopkg.in/sungora/app.v1/server/core"
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
		err      error
		store    net.Listener
		confMain *conf.Config
	)

	// configuration
	if confMain, err = conf.GetConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// logs
	if err = lg.Start(confMain.Log, confMain.NameApp, confMain.TimeZone); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer lg.Wait()

	// base controller
	if err = core.Start(confMain); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer core.Wait()

	// workflow
	if confMain.Isworkflow == true {
		if err = workflow.Start(confMain.Workflow); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return 1
		}
		defer workflow.Wait()
	}

	// web server - application
	if store, err = newWeb(confMain.Server); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	defer store.Close()
	fmt.Fprintln(os.Stdout, "web app start success")

	// The correctness of the application is closed by a signal
	signal.Notify(chanelAppControl, os.Interrupt)
	<-chanelAppControl

	return
}

// stop Stop an application
func Stop() {
	chanelAppControl <- os.Interrupt
	<-chanelAppStop
}
