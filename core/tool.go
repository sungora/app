package core

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
)

const (
	NUM     = "0123456789"
	STRDOWN = "abcdefghijklmnopqrstuvwxyz"
	STRUP   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SYMBOL  = "~!@#$%^&*_+-="
)

// NewRandomString generates password key of a specified length (a-z0-9.)
func NewRandomString(length int) string {
	return randChar(length, []byte(STRDOWN+STRUP+NUM))
}

func randChar(length int, chars []byte) string {
	pword := make([]byte, length)
	data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, data); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return ""
		}
		for _, c := range data {
			if c >= maxrb {
				continue
			}
			pword[i] = chars[c%clen]
			i++
			if i == length {
				return string(pword)
			}
		}
	}
	fmt.Fprintln(os.Stderr, "unreachable")
	return ""
}

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
