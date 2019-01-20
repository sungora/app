package tool

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"time"
)

var (
	DirConfig    string         //
	DirWork      string         //
	DirLog       string         //
	DirWww       string         //
	DirTpl       string         //
	DirPid       string         //
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
