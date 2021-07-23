package utils

import "time"

type Lock interface {
	// key , 超时时间
	Lock(string, time.Duration) (bool, error)
	// key
	UnLock(string) error
	// key ， 超时时间， 等待总时间， 失败后休眠时长
	LockTimeout(string, time.Duration, time.Duration, time.Duration) (bool, error)
}
