package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type app struct {
	nameApp string
	pathIn  string
	pathOut string
}

func newApp(nameApp, version string) *app {
	sep := string(os.PathSeparator)
	path := os.Getenv("GOPATH")
	pathList := strings.Split(path, ";")

	self := new(app)
	self.nameApp = nameApp
	self.pathOut = pathList[0] + sep + nameApp

	for _, p := range pathList {
		p1 := p + sep + "src" + sep + "gopkg.in" + sep + "sungora" + sep + "app." + version + sep + "cmd" + sep + "blankapp"
		if fi, err := os.Stat(p1); err == nil {
			if fi.IsDir() == true {
				self.pathIn = p1
				return self
			}
		}
		p1 = p + sep + "src" + sep + "vendor" + sep + "gopkg.in" + sep + "sungora" + sep + "app." + version + sep + "cmd" + sep + "blankapp"
		if fi, err := os.Stat(p1); err == nil {
			if fi.IsDir() == true {
				self.pathIn = p1
				return self
			}
		}
	}
	return nil
}

// New создание нового приложения
func (self *app) New() (err error) {
	if err = os.RemoveAll(self.pathOut); err != nil {
		return err
	}
	err = filepath.Walk(self.pathIn, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		// вычисляем и создаем целевой путь
		pathTarget := strings.Replace(path, self.pathIn, self.pathOut, 1)
		pathTarget = strings.Replace(pathTarget, ".gogogo", ".go", 1)
		if err := os.MkdirAll(filepath.Dir(pathTarget), 0777); err != nil {
			return err
		}
		// читаем файл
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// делаем необходимые замены
		str := string(data)
		str = strings.Replace(str, "zzzzzzzzz", self.nameApp, -1)
		data = []byte(str)
		// сохраняем файл
		err = ioutil.WriteFile(pathTarget, data, 0777)
		return err
	})
	return err
}

// Blank создание бланка приложения
func (self *app) Blank() (err error) {
	if err = os.RemoveAll(self.pathIn); err != nil {
		return err
	}
	err = filepath.Walk(self.pathOut, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		// вычисляем и создаем целевой путь
		pathTarget := strings.Replace(path, self.pathOut, self.pathIn, 1)
		pathTarget = strings.Replace(pathTarget, ".go", ".gogogo", 1)
		if err := os.MkdirAll(filepath.Dir(pathTarget), 0777); err != nil {
			return err
		}
		// читаем файл
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// делаем необходимые замены
		// str := string(data)
		// str = strings.Replace(str, "zzzzzzzzz", self.nameApp, -1)
		// data = []byte(str)
		// сохраняем файл
		err = ioutil.WriteFile(pathTarget, data, 0777)
		return err
	})
	return err
}
