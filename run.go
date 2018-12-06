package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"gopkg.in/sungora/app.v1/tool"
)

type run struct {
	path    string
	nameApp string
	cmd     *exec.Cmd
}

func newRun(nameApp string) *run {
	sep := string(os.PathSeparator)
	self := new(run)
	self.path = os.Getenv("GOPATH") + sep + "src" + sep + nameApp
	self.nameApp = nameApp
	os.Chdir(self.path)
	return self
}

func (self *run) Control() {
	chanelAppControl := make(chan os.Signal, 1)
	signal.Notify(chanelAppControl, os.Interrupt)

	fs := tool.NewControlFS(self.path, ".go")
	fs.CheckSumMd5()

	self.reBuild()
	buffError, buffOk := self.start()

	for {
		time.Sleep(time.Second * 1)
		select {
		case <-chanelAppControl:
			self.stop()
			goto end
		default:
			if isChange, _ := fs.CheckSumMd5(); isChange == true {
				self.stop()
				self.reBuild()
				buffError, buffOk = self.start()
			}
		}
		fmt.Printf("%s", string(buffOk.Next(buffOk.Len())))
		fmt.Printf("%s", string(buffError.Next(buffError.Len())))
	}
end:
}

func (self *run) reBuild() {
	fmt.Print("Build: ")
	self.cmd = exec.Command("go", "build", "-i")
	var buffError bytes.Buffer
	var buffOk bytes.Buffer
	self.cmd.Stderr = &buffError
	self.cmd.Stdout = &buffOk
	if err := self.cmd.Start(); err != nil {
		fmt.Println("error command build: " + err.Error())
		return
	}
	if err := self.cmd.Wait(); err != nil {
		fmt.Print("error build: " + buffError.String())
		return
	}
	if buffOk.String() != "" {
		fmt.Print(buffOk.String())
	} else {
		fmt.Println("OK")
	}
	return
}

func (self *run) start() (buffError, buffOk *bytes.Buffer) {
	fmt.Print("Start: ")
	self.cmd = exec.Command("./" + self.nameApp)
	buffError = &bytes.Buffer{}
	buffOk = &bytes.Buffer{}
	self.cmd.Stderr = buffError
	self.cmd.Stdout = buffOk
	if err := self.cmd.Start(); err != nil {
		fmt.Println("error command start: " + err.Error())
		return
	}
	fmt.Println("OK")
	return
}

func (self *run) stop() {
	fmt.Print("Stop: ")
	self.cmd.Process.Kill()
	// self.cmd.Process.Signal(os.Kill)
	if err := self.cmd.Wait(); err != nil {
		fmt.Println("error command stop: " + err.Error())
		return
	}
	fmt.Println("OK")
	return
}
