package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	key   = "key123"
	value = "value456"
	ttl   = 60
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	redis := NewRedis(&Config{RedisUrl: "redis://127.0.0.1:6379/1"})
	err := redis.Set(ctx, key, value, time.Second*ttl)
	require.Equal(t, err, nil)

	var val string
	exist, err := redis.Get(ctx, key, &val)
	require.Equal(t, err, nil)
	require.Equal(t, exist, true)
	require.Equal(t, val, value)

	ttl, err := redis.TTL(ctx, key)
	require.Equal(t, err, nil)
	fmt.Println("ttl", ttl)

	err = redis.Delete(ctx, key)
	require.Equal(t, err, nil)

	exist, err = redis.Get(ctx, key, &val)
	require.Equal(t, err, nil)
	require.Equal(t, exist, false)

	ttl, err = redis.TTL(ctx, key)
	require.Equal(t, err, nil)
	require.Less(t, ttl, 0)
	fmt.Println("ttl", ttl)
}
