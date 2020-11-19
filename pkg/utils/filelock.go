// +build !windows

package utils

import (
	"github.com/juju/fslock"
	"time"
)

var (
	ErrTimeout = fslock.ErrTimeout
	ErrLocked  = fslock.ErrLocked
)

type FileLock struct {
	inner *fslock.Lock
}

func NewFileLock(fp string) *FileLock {
	fs := FileLock{
		fslock.New(fp),
	}
	return &fs
}
func (f *FileLock) LockWithTimeout(duration time.Duration) error {
	return f.inner.LockWithTimeout(duration)
}
func (f *FileLock) Unlock() error {
	return f.inner.Unlock()
}
