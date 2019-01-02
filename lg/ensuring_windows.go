package lg

import (
	"os"
)

type fLock struct {
	fh *os.File
}

func newFLock(path string) (fLock, error) {
	var ret fLock
	var err error

	ret.fh, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	return ret, err
}

func (lock fLock) Lock() error {
	return nil
}

func (lock fLock) TryLock() (bool, error) {
	return true, nil
}

func (lock fLock) Unlock() error {
	return nil
}
