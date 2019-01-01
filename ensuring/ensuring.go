package ensuring

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

var fl fLock

// PidFileCreate Создание PID файла с ID текущего процесса
func PidFileCreate(fileName string) (err error) {
	var unlocked bool
	fl, err = newFLock(fileName)
	if err == nil {
		unlocked, err = fl.TryLock()
		if unlocked {
			// Пишем PID текущего процесса
			fl.fh.Seek(0, 0)
			fl.fh.Truncate(0)
			fmt.Fprintf(fl.fh, "%d", syscall.Getpid())
			err = fl.Lock()

		} else {
			err = errors.New("Запущен другой процесс либо PID файл заблокирован")
		}
	}
	return
}

// PidFileUnlock Снятие блокировки с PID файла
func PidFileUnlock() error {
	return fl.Unlock()
}

// CheckMemory Проверка наличия свободной памяти
func CheckMemory(minMem int) (err error) {
	defer func() {
		if errPanic := recover(); errPanic != nil {
			err = errors.New(fmt.Sprintf("%v", errPanic))
		}
	}()
	var oneMb []byte = make([]byte, 1024*1024, 1024*1024)
	var mymem []string = make([]string, minMem)

	oneMb = []byte(strings.Repeat("A", 1024*1024))
	for i := 0; i < minMem; i++ {
		mymem = append(mymem, string(oneMb))
	}
	mymem = nil
	oneMb = nil
	return
}
