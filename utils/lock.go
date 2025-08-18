package utils

import (
	"context"
	"time"
)

type Lock interface {
	// key , 超时时间
	Lock(context.Context, string, time.Duration) (bool, error)
	// key
	UnLock(context.Context, string) error
	// key ， 超时时间， 等待总时间， 失败后休眠时长
	LockTimeout(context.Context, string, time.Duration, time.Duration, time.Duration) (bool, error)
}
