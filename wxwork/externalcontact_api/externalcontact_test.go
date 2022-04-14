package externalcontact_api

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/lixinio/weixin/wxwork/authorizer"
	"github.com/stretchr/testify/require"
)

type item struct {
	client  *utils.Client
	agentID int
}

func initWxWorkAgent() *item {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, redis, redis, &agent.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})

	return &item{
		client:  agent.Client,
		agentID: agent.Config.AgentID,
	}
}

func initWxWorkSuiteAuthorizer() *item {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := authorizer.NewLite(
		redis, redis,
		test.WxWorkSuiteID,
		test.WxWorkSuiteCorpID,
		test.WxWorkSuiteAgentID,
	)

	return &item{
		client:  corp.Client,
		agentID: corp.AgentID,
	}
}

func TestGetFollowUserList(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkAgent(),
		initWxWorkSuiteAuthorizer(),
	} {
		api := NewApi(cli.client)

		result, err := api.GetFollowUserList(ctx)
		require.Equal(t, nil, err)
		fmt.Println(result.FollowUsers)
	}
}
