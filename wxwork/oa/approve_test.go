package oa_api

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/stretchr/testify/require"
)

func getOaApi() *OaApi {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, redis, redis, &agent.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})
	return NewApi(agent.Client)
}

func TestGetApprovalInfo(t *testing.T) {
	// https://blog.csdn.net/qq_40600379/article/details/116104313
	// 需要特定的agent
	oaApi := getOaApi()
	ctx := context.Background()

	result, err := oaApi.GetApprovalInfo(
		ctx,
		strconv.FormatInt(time.Now().Add(-20*24*time.Hour).Unix(), 10),
		strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10),
		0, 100, nil,
	)
	require.Equal(t, nil, err)
	fmt.Println(result.SpNolist)
}

func TestGetOpenApprovalData(t *testing.T) {
	oaApi := getOaApi()
	ctx := context.Background()

	result, err := oaApi.GetOpenApprovalData(
		ctx, test.AgentOAThirdNo,
	)
	require.Equal(t, nil, err)
	fmt.Println(result)
}
