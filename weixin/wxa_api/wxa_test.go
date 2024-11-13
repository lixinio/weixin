package wxa_api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
)

func initWxa() *utils.Client {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxa := official_account.New(redis, redis, &official_account.Config{
		Appid:  test.WxaAppid,
		Secret: test.WxaSecret,
	})
	return wxa.Client
}

// func initAuthorizer() *utils.Client {
// 	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
// 	wxopenOA := authorizer.NewLite(
// 		redis, redis,
// 		test.WxOpenAppid,
// 		test.WxOpenOAAppid,
// 	)
// 	return wxopenOA.Client
// }

func TestUrlLink(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initWxa(),
		// initAuthorizer(),
	} {
		wxaApi := NewApi(client)
		url, err := wxaApi.GenerateUrlLink(ctx, &GenerateUrlLinkRequest{
			Path:       "/modules/usedcar/Showroom/index",
			IsExpire:   true,
			ExpireTime: time.Now().Add(time.Hour).Unix(),
		})
		require.Equal(t, nil, err)
		fmt.Println(url)
		// url := "https://wxaurl.cn/74306i0l9ug"

		body, err := wxaApi.GetUrlLink(ctx, url)
		require.Equal(t, nil, err)
		fmt.Println(body)
	}
}

func TestSchema(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initWxa(),
		// initAuthorizer(),
	} {
		wxaApi := NewApi(client)
		url, err := wxaApi.GenerateScheme(ctx, &GenerateSchemeRequest{
			JumpWxa: &struct {
				Path       string `json:"path"`
				Query      string `json:"query,omitempty"`
				EnvVersion string `json:"env_version,omitempty"`
			}{
				Path: "/modules/usedcar/Showroom/index",
			},
			ExpireType: 0,
			ExpireTime: time.Now().Add(time.Hour).Unix(),
		})
		require.Equal(t, nil, err)
		fmt.Println(url)
		// url := "weixin://dl/business/?t=hTMDg0hg3hu"

		body, err := wxaApi.GetSchema(ctx, url)
		require.Equal(t, nil, err)
		fmt.Println(body)
	}
}
