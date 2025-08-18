package user_api

import (
	"context"
	"testing"
	"time"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
)

func initOfficialAccount() *utils.Client {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(redis, redis, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})
	return officialAccount.Client
}

func initAuthorizer() *utils.Client {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxopenOA := authorizer.NewLite(
		redis, redis,
		test.WxOpenAppid,
		test.WxOpenOAAppid,
	)
	return wxopenOA.Client
}

func TestUser(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
		initAuthorizer(),
	} {
		userApi := NewApi(client)
		// 用户列表
		resp, err := userApi.Get(ctx, "")
		require.Equal(t, nil, err)
		require.NotEmpty(t, resp.Data.OpenIDs)
		require.Equal(t, len(resp.Data.OpenIDs), resp.Count)

		openid := resp.Data.OpenIDs[0]

		remark := time.Now().Format("2006-01-02 15:04:05")
		{
			// 设置备注
			err := userApi.UpdateRemark(ctx, openid, remark)
			require.Equal(t, nil, err)
		}

		{
			// 用户详细信息
			resp, err := userApi.GetUserInfo(ctx, openid, "")
			require.Equal(t, nil, err)
			require.Equal(t, openid, resp.OpenID)
			// 备注一样
			require.Equal(t, remark, resp.Remark)
		}

		{
			// 批量获取
			resp, err := userApi.BatchGetUserInfo(ctx, &BatchGetUserParams{
				UserList: []struct {
					OpenID string `json:"openid"`
					Lang   string `json:"lang"`
				}{
					{
						OpenID: openid,
						Lang:   "",
					},
				},
			})
			require.Equal(t, nil, err)
			exist := false
			for _, user := range resp.UserInfoList {
				if user.OpenID == openid {
					exist = true
					break
				}
			}
			require.True(t, exist)
		}

		{
			// 拉黑
			err := userApi.BatchBlackList(ctx, []string{openid})
			require.Equal(t, nil, err)
			err = userApi.BatchBlackList(ctx, []string{openid})
			require.Equal(t, nil, err)
		}

		{
			// 获取拉黑列表
			resp, err := userApi.GetBlackList(ctx, "")
			require.Equal(t, nil, err)
			// 在黑名单
			require.Contains(t, resp.Data.OpenIDs, openid)
		}

		{
			// 取消拉黑
			err := userApi.BatchUnBlackList(ctx, []string{openid})
			require.Equal(t, nil, err)
			err = userApi.BatchUnBlackList(ctx, []string{openid})
			require.Equal(t, nil, err)
		}

		{
			// 获取拉黑列表
			resp, err := userApi.GetBlackList(ctx, "")
			require.Equal(t, nil, err)
			// 不在黑名单
			require.NotContains(t, resp.Data.OpenIDs, openid)
		}
	}
}
