package ensuring

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var AlreadyLocked = errors.New("Уже заблокирован")

type FLock struct {
	fh *os.File
}

func NewFLock(path string) (FLock, error) {
	var ret FLock
	var err error

	ret.fh, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	return ret, err
}

type FLocks []FLock

func FLockDirs(dirs ...string) (FLocks, error) {
	locks := make([]FLock, 0, len(dirs))
	allright := false
	defer func() {
		if !allright {
			for _, lock := range locks {
				lock.Unlock()
			}
		}
	}()
	var (
		err  error
		ok   bool
		lock FLock
	)
	for _, path := range dirs {
		if lock, err = NewFLock(path); err != nil {
			return nil, err
		}
		if ok, err = lock.TryLock(); err != nil {
			return nil, err
		} else if !ok {
			return nil, AlreadyLocked
		}
		locks = append(locks, lock)
	}
	allright = true
	return FLocks(locks), nil
}

func (locks FLocks) Unlock() {
	for _, lock := range locks {
		lock.Unlock()
	}
}

type DirLock string

func NewDirLock(path string) (DirLock, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return DirLock(""), err
	}
	if fi.IsDir() {
		path = filepath.Join(path, ".lock")
	} else {
		path = path + ".lock"
	}
	return DirLock(path), nil
}

func (lock DirLock) Lock() error {
	var (
		ok  bool
		err error
	)
	for {
		if ok, err = lock.TryLock(); ok && err == nil {
			return nil
		}
		if err != nil {
			return err
		}
		time.Sleep(1)
	}
}

func (lock DirLock) TryLock() (bool, error) {
	err := os.Mkdir(string(lock), 0600)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (lock DirLock) Unlock() error {
	return os.Remove(string(lock))
}

type PortLock struct {
	hostport string
	ln       net.Listener
}

func NewPortLock(port int) *PortLock {
	return &PortLock{hostport: net.JoinHostPort("127.0.0.1", strconv.Itoa(port))}
}

func (p *PortLock) Lock() {
	var err error
	t := 1 * time.Second
	for {
		if p.ln, err = net.Listen("tcp", p.hostport); err == nil {
			return
		}
		time.Sleep(t)
		t = time.Duration(float32(t) * 1.2)
	}
}

func (p *PortLock) Unlock() {
	if p.ln != nil {
		p.ln.Close()
	}
}
