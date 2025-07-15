package user_api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lixinio/weixin/utils"
	"github.com/stretchr/testify/require"
)

func TestUserTag(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
		initAuthorizer(),
	} {
		userApi := NewApi(client)
		// 获得第一个用户
		resp, err := userApi.Get(ctx, "")
		require.Equal(t, nil, err)
		require.NotEmpty(t, resp.Data.OpenIDs)
		require.Equal(t, len(resp.Data.OpenIDs), resp.Count)

		openid := resp.Data.OpenIDs[0]
		tagName := time.Now().Format("20060201150405")
		tagID := 0

		{
			// 创建标签
			resp, err := userApi.CreateTag(ctx, tagName)
			require.Equal(t, nil, err)
			tagID = resp.Tag.ID
			require.Equal(t, tagName, resp.Tag.Name)
		}

		{
			// 获取标签
			resp, err := userApi.GetTag(ctx)
			require.Equal(t, nil, err)

			var tag *TagItem = nil
			for _, item := range resp.Tags {
				if item.ID == tagID {
					tag = &item
				}
			}
			require.NotEmpty(t, tag)
			require.Equal(t, tagName, tag.Name)
		}

		tagName = fmt.Sprintf("%snew", tagName)
		{
			// 编辑标签
			require.Equal(t, nil, userApi.UpdateTag(ctx, tagID, tagName))
		}

		{
			// 重新 获取标签
			resp, err := userApi.GetTag(ctx)
			require.Equal(t, nil, err)

			var tag *TagItem = nil
			for _, item := range resp.Tags {
				if item.ID == tagID {
					tag = &item
				}
			}
			require.NotEmpty(t, tag)
			require.Equal(t, tagName, tag.Name)
		}

		{
			// 给粉丝打标签
			require.Equal(t, nil, userApi.BatchTagging(ctx, tagID, []string{openid}))
		}

		{
			// 获取标签下粉丝列表
			resp, err := userApi.GetUsersByTag(ctx, tagID, "")
			require.Equal(t, nil, err)
			require.Greater(t, resp.Count, 0)
			require.Contains(t, resp.Data.OpenIDs, openid)
		}

		{
			// 获取用户身上的标签列表
			resp, err := userApi.GetTagIdList(ctx, openid)
			require.Equal(t, nil, err)
			require.Contains(t, resp.TagIDList, tagID)
		}

		{
			// 删除粉丝标签
			require.Equal(t, nil, userApi.BatchUnTagging(ctx, tagID, []string{openid}))
		}

		{
			// 验证删除粉丝标签
			resp, err := userApi.GetTagIdList(ctx, openid)
			require.Equal(t, nil, err)
			require.NotContains(t, resp.TagIDList, tagID)
		}

		{
			// 删除标签
			require.Equal(t, nil, userApi.DeleteTag(ctx, tagID))
		}

		{
			// 验证删除标签
			resp, err := userApi.GetTag(ctx)
			require.Equal(t, nil, err)

			var tag *TagItem = nil
			for _, item := range resp.Tags {
				if item.ID == tagID {
					tag = &item
				}
			}
			require.Empty(t, tag)
		}
	}
}
