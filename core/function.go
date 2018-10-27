package core

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

// CreatePassword make random password
func CreatePassword() string {
	c := 10
	b := make([]byte, c)
	n, err := io.ReadFull(rand.Reader, b)
	if n != len(b) || err != nil {
		fmt.Println("error:", err)
	}
	return fmt.Sprintf("%x", b)
}

// CreatePasswordHash make password hash
func CreatePasswordHash(password string) string {
	shaCoo := sha256.New()
	shaCoo.Write([]byte(password))
	return fmt.Sprintf("%x", shaCoo.Sum(nil))
}
