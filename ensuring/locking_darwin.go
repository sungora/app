package ensuring

import (
	"syscall"
)

func (lock FLock) Lock() error {
	return syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_EX)
}

func (lock FLock) TryLock() (bool, error) {
	err := syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	switch err {
	case nil:
		return true, nil
	case syscall.EWOULDBLOCK:
		return false, nil
	}
	return false, err
}

func (lock FLock) Unlock() error {
	lock.fh.Close()
	return syscall.Flock(int(lock.fh.Fd()), syscall.LOCK_UN)
}
