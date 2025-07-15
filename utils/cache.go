package utils

import (
	"context"
	"time"
)

// https://github.com/silenceper/wechat/blob/master/cache/cache.go

// Cache interface
type Cache interface {
	Get(context.Context, string, interface{}) (bool, error) // 不存在的情况(false,nil)
	Set(context.Context, string, interface{}, time.Duration) error
	IsExist(context.Context, string) bool
	Delete(context.Context, string) error
	TTL(context.Context, string) (int, error)
}
