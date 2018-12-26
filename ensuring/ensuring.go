package ensuring

import (
	"errors"
	"fmt"
	"strings"
	"syscall"
)

var fl FLock

// PidFileCreate Создание PID файла с ID текущего процесса
func PidFileCreate(fileName string) (err error) {
	var unlocked bool
	fl, err = NewFLock(fileName)
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
