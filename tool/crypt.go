package tool

import (
	"crypto/rand"
	"io"
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

// CreatePassword make random password
// func CreatePassword() string {
// 	c := 10
// 	b := make([]byte, c)
// 	n, err := io.ReadFull(rand.Reader, b)
// 	if n != len(b) || err != nil {
// 		fmt.Println("error:", err)
// 	}
// 	return fmt.Sprintf("%x", b)
// }

// CreatePasswordHash make password hash
// func CreatePasswordHash(password string) string {
// 	shaCoo := sha256.New()
// 	shaCoo.Write([]byte(password))
// 	return fmt.Sprintf("%x", shaCoo.Sum(nil))
// }
