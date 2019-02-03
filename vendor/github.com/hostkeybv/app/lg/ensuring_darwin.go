package lg

import (
	"os"
	"syscall"
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
	return syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_EX)
}

func (lock fLock) TryLock() (bool, error) {
	err := syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	switch err {
	case nil:
		return true, nil
	case syscall.EWOULDBLOCK:
		return false, nil
	}
	return false, err
}

func (lock fLock) Unlock() error {
	lock.fh.Close()
	return syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_UN)
}
