package miniprogram

import (
	"context"
	"fmt"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/stretchr/testify/require"
	"testing"
)

func initSubscribeClient() *utils.Client {
	r := redis.NewRedis(&redis.Config{RedisUrl: "redis://@localhost:6379/7"})
	mpClient := New(r, r, &Config{
		Appid:  "appid",
		Secret: "app secret",
	})

	return mpClient.Client
}

type messageItem struct {
	TmplID string
	OpenID string
	Client *utils.Client
}

func TestSubscribeMessage(t *testing.T) {
	ctx := context.Background()

	for _, client := range []*messageItem{
		{
			TmplID: "subscribe template id",
			OpenID: "odUML5P7gbNiUIJz-eZA5IM7sEpg",
			Client: initSubscribeClient(),
		},
	} {
		messageApi := NewApi(client.Client)
		respMessage, err := messageApi.SendSubscribeMessage(
			ctx,
			client.TmplID,
			client.OpenID,
			map[string]DataValue{},
		)

		require.Equal(t, nil, err)
		fmt.Println(respMessage)
	}
}
