package department_api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	agentApi "github.com/lixinio/weixin/wxwork/agent"
	"github.com/stretchr/testify/require"
)

func TestDepartment(t *testing.T) {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agentApi.New(corp, redis, redis, &agentApi.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})
	agentContact := agentApi.New(corp, redis, redis, &agentApi.Config{
		AgentID: 0,
		Secret:  test.AgentContactSecret,
	})
	ctx := context.Background()

	departmentName := time.Now().Format("20060201150405")
	departmentApi := NewApi(agent.Client)
	departmentContactApi := NewApi(agentContact.Client)

	rootDepartmentID := 0
	{
		// 获取根目录
		resp, err := departmentApi.List(ctx, 0)
		require.Equal(t, nil, err)
		require.Greater(t, len(resp.Department), 0)
		rootDepartmentID = resp.Department[0].ID
	}

	newDepartmentID := 0
	{
		// 创建部门
		resp, err := departmentContactApi.Create(ctx, &CreateParam{
			Parentid: rootDepartmentID,
			Name:     departmentName,
		})
		require.Equal(t, nil, err)
		require.Equal(t, "created", resp.ErrMsg)
		require.Greater(t, resp.ID, 0)
		newDepartmentID = resp.ID
	}

	{
		// 更新部门
		require.Equal(t, nil, departmentContactApi.Update(ctx, &UpdateParam{
			ID:   newDepartmentID,
			Name: fmt.Sprintf("%snew", departmentName),
		}))
	}

	{
		// 确认新部门
		resp, err := departmentApi.List(ctx, newDepartmentID)
		require.Equal(t, nil, err)
		require.Greater(t, len(resp.Department), 0)
		require.Equal(t, newDepartmentID, resp.Department[0].ID)
		require.Equal(t, fmt.Sprintf("%snew", departmentName), resp.Department[0].Name)
	}

	{
		// 删除部门
		require.Equal(t, nil, departmentContactApi.Delete(ctx, newDepartmentID))
	}

	{
		// 确认删除结果
		resp, err := departmentApi.List(ctx, 0)
		require.Equal(t, nil, err)
		require.Greater(t, len(resp.Department), 0)

		var department *DepartmentItem = nil
		for _, item := range resp.Department {
			if item.ID == newDepartmentID {
				department = &item
			}
		}

		require.Empty(t, department)
	}
}
