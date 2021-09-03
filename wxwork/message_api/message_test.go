package message_api

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

func TestSendMessage(t *testing.T) {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, redis, redis, &agent.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})
	ctx := context.Background()

	messageApi := NewApi(agent.Client, agent.Config.AgentID)
	result, err := messageApi.SendTextMessage(
		ctx,
		&MessageHeader{ToUser: test.AgentUserID},
		"你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。",
	)
	require.Equal(t, nil, err)
	fmt.Println(result.MsgID)
}
