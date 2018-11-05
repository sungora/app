package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"gopkg.in/sungora/app.v1/tool"

	"github.com/logrusorgru/aurora"
)

type run struct {
	path    string
	pathApp string
	cmd     *exec.Cmd
}

func NewRun(nameApp string) *run {
	sep := string(os.PathSeparator)
	self := new(run)
	self.path = os.Getenv("GOPATH") + sep + "src" + sep + nameApp
	self.pathApp = self.path + sep + nameApp
	return self
}

func (self *run) Control() {

	chanelAppControl := make(chan os.Signal, 1)
	signal.Notify(chanelAppControl, os.Interrupt)

	fs := tool.NewControlFS()
	fs.CheckSumMd5(self.path, ".go")

	self.reBuild()
	self.start()

	for {
		time.Sleep(time.Second * 1)
		select {
		case <-chanelAppControl:
			self.stop()
			goto end
		default:
			if isChange, _ := fs.CheckSumMd5(self.path, ".go"); isChange == true {
				self.stop()
				self.reBuild()
				self.start()
			}
		}
	}
end:
}

func (self *run) reBuild() (err error) {
	fmt.Print("Build: ")
	self.cmd = exec.Command("go", "build", "-i")
	var buffError bytes.Buffer
	var buffOk bytes.Buffer
	self.cmd.Stderr = &buffError
	self.cmd.Stdout = &buffOk
	if err = self.cmd.Start(); err != nil {
		fmt.Println(aurora.Magenta("error command build: " + err.Error()))
		return
	}
	if err = self.cmd.Wait(); err != nil {
		fmt.Print(aurora.Red("error build: " + buffError.String()))
		return
	}
	if buffOk.String() != "" {
		fmt.Print(aurora.Green(buffOk.String()))
	} else {
		fmt.Println(aurora.Bold(aurora.Green("OK")))
	}
	return
}

func (self *run) start() (buffError, buffOk bytes.Buffer, err error) {
	fmt.Print("Start: ")
	self.cmd = exec.Command(self.pathApp)
	// var buffError bytes.Buffer
	// var buffOk bytes.Buffer
	self.cmd.Stderr = &buffError
	self.cmd.Stdout = &buffOk
	if err = self.cmd.Start(); err != nil {
		fmt.Println(aurora.Magenta("error command start: " + err.Error()))
		return
	}
	// if err = self.cmd.Wait(); err != nil {
	// 	fmt.Print(aurora.Red("error start: " + buffError.String()))
	// 	return
	// }
	fmt.Println(aurora.Bold(aurora.Green("OK")))
	// if buffOk.String() != "" {
	// 	fmt.Print(aurora.Green(buffOk.String()))
	// } else {
	// 	fmt.Println(aurora.Bold(aurora.Green("OK")))
	// }
	return
}

func (self *run) stop() (err error) {
	fmt.Print("Stop: ")
	self.cmd.Process.Signal(os.Interrupt)
	if err = self.cmd.Wait(); err != nil {
		fmt.Println(aurora.Magenta("error command stop: " + err.Error()))
		return
	}
	fmt.Println(aurora.Bold(aurora.Green("OK")))
	return
}
