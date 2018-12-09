package tool

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

func HtmlCompilation(pathTpl string, functions, variables map[string]interface{}) (data string, err error) {
	if _, err = os.Lstat(pathTpl); err != nil {
		return
	}
	var tpl *template.Template
	tpl, err = template.New(filepath.Base(pathTpl)).Funcs(functions).ParseFiles(pathTpl)
	if err != nil {
		return "", err
	}
	var ret bytes.Buffer
	if err = tpl.Execute(&ret, variables); err != nil {
		return "", err
	}
	return ret.String(), nil
}
