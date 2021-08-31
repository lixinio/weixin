package main

import (
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/stretchr/testify/require"
)

func TestTokenCache(t *testing.T) {
	componentAppid, appid := "abcdefg", "hijklmn"
	token := "1234567890"
	token2 := token + "refresh"
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	manager := NewAuthorizerRefreshTokenManager(redis, redis)
	authorizer, err := manager.GetTokenCache(componentAppid, appid)
	require.Equal(t, err, nil)

	err = authorizer.SetAccessToken(token, 3600)
	require.Equal(t, err, nil)
	err = authorizer.SetRefreshToken(token2)
	require.Equal(t, err, nil)

	newToken, err := authorizer.GetAceessToken()
	require.Equal(t, err, nil)
	require.Equal(t, newToken, token)

	newToken, err = authorizer.GetRefreshToken()
	require.Equal(t, err, nil)
	require.Equal(t, newToken, token2)
}
