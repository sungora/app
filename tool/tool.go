package tool

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	sep := string(os.PathSeparator)
	DirWork, _ = os.Getwd()
	DirLog = DirWork + sep + "log"
	DirWww = DirWork + sep + "www"
	DirTpl = DirWork + sep + "tpl"
	DirConfig = DirWork + sep + "config"
	if ext := filepath.Ext(os.Args[0]); ext != "" {
		sl := strings.Split(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
		ServiceName = sl[0]
	} else {
		ServiceName = filepath.Base(os.Args[0])
	}
}

var (
	DirWork      string         //
	DirConfig    string         //
	DirLog       string         //
	DirWww       string         //
	DirTpl       string         //
	ServiceName  string         // Техническое название приложения
	TimeLocation *time.Location // Временная зона
)

// Dump all variables to STDOUT
func Dumper(idl ...interface{}) string {
	ret := dump(idl...)
	fmt.Print(ret.String())
	return ret.String()
}

// Dump all variables to bytes.Buffer
func dump(idl ...interface{}) bytes.Buffer {
	var buf bytes.Buffer
	var wr io.Writer

	wr = io.MultiWriter(&buf)
	for _, field := range idl {
		fset := token.NewFileSet()
		ast.Fprint(wr, fset, field, ast.NotNilFilter)
	}
	return buf
}
