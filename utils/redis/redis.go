package redis

// https://github.com/silenceper/wechat/blob/master/cache/redis.go
import (
	"fmt"
	"reflect"
	"time"

	"github.com/gomodule/redigo/redis"
)

//Redis redis cache
type Redis struct {
	conn *redis.Pool
}

// Config redis 连接属性
type Config struct {
	RedisUrl    string
	MaxIdle     int   `yml:"max_idle"     json:"max_idle"`
	MaxActive   int   `yml:"max_active"   json:"max_active"`
	IdleTimeout int32 `yml:"idle_timeout" json:"idle_timeout"` //second
}

//NewRedis 实例化
func NewRedis(opts *Config) *Redis {
	redisDB, redisHost, redisPwd, err := parseRedisURL(opts.RedisUrl)
	if err != nil {
		panic(err)
	}

	pool := &redis.Pool{
		MaxActive:   opts.MaxActive,
		MaxIdle:     opts.MaxIdle,
		IdleTimeout: time.Second * time.Duration(opts.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisHost,
				redis.DialDatabase(redisDB),
				redis.DialPassword(redisPwd),
			)
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
	return &Redis{pool}
}

//SetConn 设置conn
func (r *Redis) SetConn(conn *redis.Pool) {
	r.conn = conn
}

//Get 获取一个值
func (r *Redis) Get(key string, value interface{}) (exist bool, err error) {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		if err == redis.ErrNil {
			// 不存在特殊处理
			return false, nil
		}
	}

	v, ok := value.(*string)
	if !ok {
		err = fmt.Errorf("value must be pointer to string")
		return
	}
	*v = string(data)

	return true, nil
}

//Set 设置一个值
func (r *Redis) Set(key string, val interface{}, timeout time.Duration) (err error) {
	conn := r.conn.Get()
	defer conn.Close()

	data, ok := val.(string)
	if !ok {
		err = fmt.Errorf("val must be string")
		return
	}

	_, err = conn.Do("SETEX", key, int64(timeout/time.Second), data)

	return
}

//IsExist 判断key是否存在
func (r *Redis) IsExist(key string) bool {
	conn := r.conn.Get()
	defer conn.Close()

	a, _ := conn.Do("EXISTS", key)
	i := a.(int64)
	return i > 0
}

//Delete 删除
func (r *Redis) Delete(key string) error {
	conn := r.conn.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}

	return nil
}

// 获得剩余时间(秒)
func (r *Redis) TTL(key string) (int, error) {
	conn := r.conn.Get()
	defer conn.Close()

	if reply, err := conn.Do("TTL", key); err != nil {
		return -1, err
	} else {
		if ttl, ok := reply.(int64); ok {
			// 如果不存在, 返回-2
			return int(ttl), nil
		} else {
			return -1, fmt.Errorf(
				"invalid ttl reply type '%s'", reflect.TypeOf(reply).String(),
			)
		}
	}
}

// https://www.programmersought.com/article/85921351841/
// http://xiaorui.cc/archives/3028
func (r *Redis) Lock(key string, expire time.Duration) (bool, error) {
	conn := r.conn.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("set", key, 1, "ex", int(expire/time.Second), "nx"))
	if err != nil {
		if err == redis.ErrNil {
			// The lock was not successful, it already exists.
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Redis) LockTimeout(key string, expire, timeout, sleep time.Duration) (bool, error) {
	var total time.Duration = 0
	for total < timeout {
		result, err := r.Lock(key, expire)
		if err == nil {
			if result {
				// lock success
				return result, err
			} else {
				// lock fail
				time.Sleep(sleep)
				total += sleep
			}
		} else {
			// error
			return result, err
		}
	}
	// lock fail
	return false, nil
}

func (r *Redis) UnLock(key string) error {
	conn := r.conn.Get()
	defer conn.Close()

	_, err := conn.Do("del", key)
	if err != nil {
		return err
	}
	return nil
}
