package tpl

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

const nameDinamicTpl = "dinamic"

// ParseFile компиляция html шаблона из указанного файла
func ParseFile(viewPath string, functions, variables map[string]interface{}) (ret bytes.Buffer, err error) {
	if _, err = os.Stat(viewPath); err != nil {
		return
	}
	var tpl *template.Template
	if tpl, err = template.New(filepath.Base(viewPath)).Funcs(functions).ParseFiles(viewPath); err != nil {
		return
	}
	if err = tpl.Execute(&ret, variables); err != nil {
		return
	}
	return
}

// ParseText компиляция html шаблона переданного в строке
func ParseText(view string, functions, variables map[string]interface{}) (ret bytes.Buffer, err error) {
	var tpl *template.Template
	if tpl, err = template.New(nameDinamicTpl).Funcs(functions).Parse(view); err != nil {
		return
	}
	if err = tpl.Execute(&ret, variables); err != nil {
		return
	}
	return
}
