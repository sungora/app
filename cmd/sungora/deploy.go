package main

import (
	"bytes"
	"fmt"
	"gopkg.in/sungora/app.v1/tool"
	"os"
	"os/exec"
	"time"
)

func Save(s string) {
	path, _ := os.Executable()
	fp, _ := os.Create(path + "_deploy.txt")
	fp.WriteString(s)
	fp.Close()
}

func Deploy(nameApp string) {
	path, _ := os.Executable()

	fs := tool.NewControlFS()
	fs.CheckSumMd5(path, "")
	for {
		fmt.Println("iteration deploy ", path)
		time.Sleep(time.Second * 3)
		if isChange, _ := fs.CheckSumMd5(path, ""); isChange == true {

			l := time.Now().Format("2006-01-02 15:04:05")
			fmt.Println(l, " rebuild: ", nameApp)

			cmd := exec.Command(path, "restart")
			var buffError bytes.Buffer
			cmd.Stderr = &buffError
			if err := cmd.Start(); err != nil {
				fp, _ := os.Create(path + "_deploy1.txt")
				fp.WriteString(err.Error())
				fp.Close()
				return
			}
			if err := cmd.Wait(); err != nil {
				fp, _ := os.Create(path + "_deploy2.txt")
				fp.WriteString(err.Error())
				fp.Close()
				return
			}
			fp, _ := os.Create(path + "_deploy3.txt")
			fp.WriteString("ok")
			fp.Close()
			return
		}
	}
}
