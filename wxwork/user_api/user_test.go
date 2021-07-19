package user_api

import (
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, cache, &agent.Config{
		AgentId: test.AgentID,
		Secret:  test.AgentSecret,
	})

	userApi := NewAgentApi(agent)
	{
		resp, err := userApi.Get(test.AgentUserID)
		require.Equal(t, nil, err)
		require.Equal(t, test.AgentUserID, resp.UserID)
	}
}
