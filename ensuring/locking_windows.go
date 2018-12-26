package ensuring

import (
)

func (lock FLock) Lock() error {
	return nil
}

func (lock FLock) TryLock() (bool, error) {
	return true, nil
}

func (lock FLock) Unlock() error {
	return nil
}
