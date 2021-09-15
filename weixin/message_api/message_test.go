package message_api

import (
	"context"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/weixin/official_account"
)

type messageItem struct {
	OpenID string
	Client *utils.Client
}

func initOfficialAccount() *messageItem {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(redis, redis, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})
	return &messageItem{
		OpenID: test.OfficialAccountOpenid,
		Client: officialAccount.Client,
	}
}

func initAuthorizer() *messageItem {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxopenOA := authorizer.NewLite(
		redis, redis,
		test.WxOpenAppid,
		test.WxOpenOAAppid,
	)
	return &messageItem{
		OpenID: test.WxOpenOAOpenID,
		Client: wxopenOA.Client,
	}
}

func TestCustomerMessage(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*messageItem{
		initOfficialAccount(),
		initAuthorizer(),
	} {
		messageApi := NewApi(client.Client)
		messageApi.SendCustomTextMessage(ctx, client.OpenID, "发多了开发")
	}
}
