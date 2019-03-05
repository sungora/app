package core

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

// TplCompilation компиляция html шаблона
func TplCompilation(pathTpl string, functions, variables map[string]interface{}) (ret bytes.Buffer, err error) {
	if _, err = os.Stat(pathTpl); err != nil {
		return
	}
	var tpl *template.Template
	if tpl, err = template.New(filepath.Base(pathTpl)).Funcs(functions).ParseFiles(pathTpl); err != nil {
		return
	}
	if err = tpl.Execute(&ret, variables); err != nil {
		return
	}
	return
}
