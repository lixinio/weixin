package tag_api

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/stretchr/testify/require"
)

func TestTag(t *testing.T) {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, redis, redis, &agent.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})
	ctx := context.Background()

	tagApi := NewApi(agent.Client)
	taglist, err := tagApi.List(ctx)
	require.Equal(t, nil, err)
	for _, item := range taglist.TagList {
		fmt.Println(item.TagID, item.TagName)
	}
}
