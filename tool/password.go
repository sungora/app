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

// NewKey generates key of a specified length (a-z0-9)
func NewKey(length int, char ...string) string {
	if 0 < len(char) {
		return randChar(length, []byte(char[0]))
	}
	return randChar(length, []byte(STRDOWN+NUM))
}

// NewPass generates password key of a specified length (a-z0-9.)
func NewPass(length int, char ...string) string {
	if 0 < len(char) {
		return randChar(length, []byte(char[0]))
	}
	return randChar(length, []byte(STRDOWN+STRUP+SYMBOL+NUM))
}

// NewKeyAPI generates keys such kind: uuu-xxxx-zzzzz
func NewKeyAPI(length int, char ...string) string {
	if 0 < len(char) {
		return NewKey(length, char[0]) + "-" + NewKey(length+1, char[0]) + "-" + NewKey(length+2, char[0])
	}
	return NewKey(length) + "-" + NewKey(length+1) + "-" + NewKey(length+2)
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
//func CreatePassword() string {
//	c := 10
//	b := make([]byte, c)
//	n, err := io.ReadFull(rand.Reader, b)
//	if n != len(b) || err != nil {
//		fmt.Println("error:", err)
//	}
//	return fmt.Sprintf("%x", b)
//}

// CreatePasswordHash make password hash
//func CreatePasswordHash(password string) string {
//	shaCoo := sha256.New()
//	shaCoo.Write([]byte(password))
//	return fmt.Sprintf("%x", shaCoo.Sum(nil))
//}
