package utils

import (
	"time"
)

// https://github.com/silenceper/wechat/blob/master/cache/cache.go

//Cache interface
type Cache interface {
	Get(string, interface{}) (bool, error) // 不存在的情况(false,nil)
	Set(string, interface{}, time.Duration) error
	IsExist(string) bool
	Delete(string) error
	TTL(string) (int, error)
}
