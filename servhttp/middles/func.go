package middles

import (
	"crypto/rand"
	"io"
)

const (
	num     = "0123456789"
	strdown = "abcdefghijklmnopqrstuvwxyz"
	strup   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbol  = "~!@#$%^&*_+-="
)

// newRandomString generates password key of a specified length (a-z0-9.)
func newRandomString(length int) string {
	return randChar(length, []byte(strdown+strup+num))
}

func randChar(length int, chars []byte) string {
	pword := make([]byte, length)
	data := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, data); err != nil {
			panic(err)
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
	panic("unreachable")
}
