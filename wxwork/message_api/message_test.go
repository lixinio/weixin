package message_api

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

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkAgent(),
		initWxWorkSuiteAuthorizer(),
	} {
		messageApi := NewApi(cli.client, cli.agentID)
		result, err := messageApi.SendTextMessage(
			ctx,
			&MessageHeader{ToUser: test.AgentUserID},
			"你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。",
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}

}
