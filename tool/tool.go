package tool

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"time"
)

func init() {
	DirWork, _ = os.Getwd()
	sep := string(os.PathSeparator)
	DirConfig = DirWork + sep + "config"
	DirWww = DirWork + sep + "www"
	DirTpl = DirWork + sep + "tpl"
}

var (
	DirWork      string
	DirConfig    string
	DirWww       string
	DirTpl       string
	TimeLocation *time.Location
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
