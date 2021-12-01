package user_api

import (
	"context"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/lixinio/weixin/wxwork/authorizer"
	"github.com/stretchr/testify/require"
)

func getClient() []*utils.Client {
	rds := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, rds, rds, &agent.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})

	authorizer := authorizer.NewLite(
		rds, rds, test.WxWorkSuiteID,
		test.WxWorkSuiteCorpID, test.WxWorkSuiteAgentID,
	)
	return []*utils.Client{
		agent.Client,
		authorizer.Client,
	}
}

func TestUser(t *testing.T) {
	clis := getClient()
	ctx := context.Background()

	for _, cli := range clis {
		userApi := NewApi(cli)
		{
			userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
			require.Equal(t, nil, err)

			resp, err := userApi.Get(ctx, userid)
			require.Equal(t, nil, err)

			openID, err := userApi.ConvertToOpenId(ctx, resp.UserID)
			require.Equal(t, nil, err)

			newUserId, err := userApi.ConvertToUserId(ctx, openID)
			require.Equal(t, nil, err)
			require.Equal(t, userid, newUserId)

			_, err = userApi.SimpleList(ctx, test.AgentRootDep, 1)
			require.Equal(t, nil, err)

			_, err = userApi.List(ctx, test.AgentRootDep, 1)
			require.Equal(t, nil, err)
		}
	}
}
