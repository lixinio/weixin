package message_api

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
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

/*
测试号模板
{{title.DATA}}
{{msg1.DATA}}
{{msg2.DATA}}
{{msg3.DATA}}
{{msg4.DATA}}
{{msg5.DATA}}
*/
func TestTemplateMessage(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*messageItem{
		initOfficialAccount(),
		initAuthorizer(),
	} {
		messageApi := NewApi(client.Client)
		id, err := messageApi.SendTemplateMessage(ctx, &TemplateMessage{
			ToUser:     client.OpenID,
			TemplateID: "pBP6oKHSz4UOkadAYsLw-ug5G495JRZsEVu626eJEmo",
			Datas: map[string]*TemplateMessageData{
				"title": {
					Value: "标题",
				},
				"msg1": {
					Value: "标题msg1",
				},
				"msg2": {
					Value: "标题msg2",
				},
				"msg3": {
					Value: "标题msg3",
				},
				"msg4": {
					Value: "标题msg4",
				},
				"msg5": {
					Value: "标题msg5",
				},
			},
		})
		require.Equal(t, nil, err)
		fmt.Println(id)
	}
}
