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
	for _, p := range strings.Split(path, ";") {
		p1 := p + sep + "src" + sep + "gopkg.in" + sep + "sungora" + sep + "app." + version + sep + "cmd" + sep + "blankapp"
		if fi, err := os.Stat(p1); err == nil {
			if fi.IsDir() == true {
				self := new(app)
				self.nameApp = nameApp
				self.pathIn = p1
				self.pathOut, _ = os.Getwd()
				return self
			}
		}
		p1 = p + sep + "src" + sep + "vendor" + sep + "gopkg.in" + sep + "sungora" + sep + "app." + version + sep + "cmd" + sep + "blankapp"
		if fi, err := os.Stat(p1); err == nil {
			if fi.IsDir() == true {
				self := new(app)
				self.nameApp = nameApp
				self.pathIn = p1
				self.pathOut, _ = os.Getwd()
				return self
			}
		}
	}
	return nil
}

func (self *app) Copy() (err error) {
	sep := string(os.PathSeparator)
	err = filepath.Walk(self.pathIn, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		pathTarget := strings.Replace(path, self.pathIn, self.pathOut+sep+self.nameApp, 1)
		if err := os.MkdirAll(filepath.Dir(pathTarget), 0777); err != nil {
			return err
		}
		//
		str := string(data)
		str = strings.Replace(str, "PKGAPPNAME", self.nameApp, -1)
		data = []byte(str)
		//
		err = ioutil.WriteFile(pathTarget, data, 0777)
		return err
	})
	return err
}
