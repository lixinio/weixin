package redis

// https://github.com/silenceper/wechat/blob/master/cache/redis.go

import (
	"encoding/json"
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
	MaxIdle     int   `yml:"max_idle" json:"max_idle"`
	MaxActive   int   `yml:"max_active" json:"max_active"`
	IdleTimeout int32 `yml:"idle_timeout" json:"idle_timeout"` //second

	// Host        string `yml:"host" json:"host"`
	// Password    string `yml:"password" json:"password"`
	// Database    int    `yml:"database" json:"database"`
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
func (r *Redis) Get(key string) interface{} {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	var err error
	if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		return nil
	}
	var reply interface{}
	if err = json.Unmarshal(data, &reply); err != nil {
		return nil
	}

	return reply
}

//Set 设置一个值
func (r *Redis) Set(key string, val interface{}, timeout time.Duration) (err error) {
	conn := r.conn.Get()
	defer conn.Close()

	var data []byte
	if data, err = json.Marshal(val); err != nil {
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
	if i > 0 {
		return true
	}
	return false
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
