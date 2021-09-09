package wxopen

import (
	"context"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/stretchr/testify/require"
)

func initWxOpen() *WxOpen {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxopen := New(redis, redis, &Config{
		Appid:          test.WxOpenAppid,
		Secret:         test.WxOpenSecret,
		Token:          test.WxOpenToken,
		EncodingAESKey: test.WxOpenEncodingAESKey,
	})
	return wxopen
}

func TestStartPushTicket(t *testing.T) {
	open := initWxOpen()
	err := open.StartPushTicket(context.Background())
	require.Empty(t, err)
}

// func TestRefreshToken(t *testing.T) {
// 	open := initWxOpen()
// 	token, _, err := open.refreshAccessTokenFromWXServer()
// 	fmt.Println(token)
// 	require.Empty(t, err)
// }
