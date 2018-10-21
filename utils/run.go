package utils

import (
	"bytes"
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type run struct {
	path    string
	pathApp string
	cmd     *exec.Cmd
}

func NewRun() *run {
	self := new(run)
	sep := string(os.PathSeparator)
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	sl := strings.Split(dir, sep)
	program := sl[len(sl)-1]
	if sep != "/" {
		program += ".exe"
	}
	self.path = dir
	self.pathApp = dir + sep + program
	return self
}

func (self *run) Control() {

	chanelAppControl := make(chan os.Signal, 1)
	signal.Notify(chanelAppControl, os.Interrupt)

	fs := NewControlFS()
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

func (self *run) start() (err error) {
	fmt.Print("Start: ")
	self.cmd = exec.Command(self.pathApp)
	var buffError bytes.Buffer
	var buffOk bytes.Buffer
	self.cmd.Stderr = &buffError
	self.cmd.Stdout = &buffOk
	if err = self.cmd.Start(); err != nil {
		fmt.Println(aurora.Magenta("error command start: " + err.Error()))
		return
	}
	if err = self.cmd.Wait(); err != nil {
		fmt.Print(aurora.Red("error start: " + buffError.String()))
		return
	}
	if buffOk.String() != "" {
		fmt.Print(aurora.Green(buffOk.String()))
	} else {
		fmt.Println(aurora.Bold(aurora.Green("OK")))
	}
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
