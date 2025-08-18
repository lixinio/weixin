package main

import (
	"context"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/stretchr/testify/require"
)

func TestTokenCache(t *testing.T) {
	ctx := context.Background()
	componentAppid, appid := "abcdefg", "hijklmn"
	token := "1234567890"
	token2 := token + "refresh"
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	manager := NewAuthorizerRefreshTokenManager(redis, redis)
	authorizer, err := manager.GetTokenCache(componentAppid, appid)
	require.Equal(t, err, nil)

	err = authorizer.SetAccessToken(ctx, token, 3600)
	require.Equal(t, err, nil)
	err = authorizer.SetRefreshToken(ctx, token2)
	require.Equal(t, err, nil)

	newToken, err := authorizer.GetAccessToken(ctx)
	require.Equal(t, err, nil)
	require.Equal(t, newToken, token)

	newToken, err = authorizer.GetRefreshToken(ctx)
	require.Equal(t, err, nil)
	require.Equal(t, newToken, token2)
}
