package agent_api

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

func getClient() []*item {
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
	return []*item{
		{
			client:  agent.Client,
			agentID: agent.Config.AgentID,
		},
		{
			client:  authorizer.Client,
			agentID: authorizer.AgentID,
		},
	}
}

func TestSetWorkbenchImageTemplate(t *testing.T) {
	items := getClient()
	ctx := context.Background()

	for _, item := range items {
		api := NewApi(item.client)
		{
			err := api.SetWorkbenchTemplate(ctx, &WorkbenchTemplateParam{
				AgentID: item.agentID,
				WorkbenchTemplate: WorkbenchTemplate{
					Type:            KeyTypeImage,
					ReplaceUserData: true,
					WorkBenchImage: &WorkBenchImageItem{
						URL:     "https://www.baidu.com",
						JumpURL: "https://www.baidu.com",
					},
				},
			})
			require.Equal(t, nil, err)
		}
	}
}

func TestSetWorkbenchKeyDataTemplate(t *testing.T) {
	items := getClient()
	ctx := context.Background()

	for _, item := range items {
		api := NewApi(item.client)
		{
			err := api.SetWorkbenchTemplate(ctx, &WorkbenchTemplateParam{
				AgentID: item.agentID,
				WorkbenchTemplate: WorkbenchTemplate{
					Type:            KeyTypeKeyData,
					ReplaceUserData: true,
					WorkBenchKeyData: &WorkBenchKeyDataItem{
						Items: []*WorkBenchKeyDataItemSection{
							{
								Key:     "待审批",
								Data:    "2",
								JumpURL: "http://www.qq.com",
							},
							{
								Key:     "带批阅作业",
								Data:    "4",
								JumpURL: "http://www.qq.com",
							},
							{
								Key:     "成绩录入",
								Data:    "45",
								JumpURL: "http://www.qq.com",
							},
							{
								Key:     "综合评价",
								Data:    "98",
								JumpURL: "http://www.qq.com",
							},
						},
					},
				},
			})
			require.Equal(t, nil, err)
		}
	}
}

func TestSetWorkbenchListTemplate(t *testing.T) {
	items := getClient()
	ctx := context.Background()

	for _, item := range items {
		api := NewApi(item.client)
		{
			err := api.SetWorkbenchTemplate(ctx, &WorkbenchTemplateParam{
				AgentID: item.agentID,
				WorkbenchTemplate: WorkbenchTemplate{
					Type:            KeyTypeList,
					ReplaceUserData: true,
					WorkBenchList: &WorkBenchListItem{
						Items: []*WorkBenchListItemSection{
							{
								Title:   "智慧校园新版设计的小程序要来啦",
								JumpURL: "http://www.qq.com",
							},
							{
								Title:   "植物百科，这是什么花",
								JumpURL: "http://www.qq.com",
							},
							{
								Title:   "周一升旗通知，全体学生必须穿校服",
								JumpURL: "http://www.qq.com",
							},
						},
					},
				},
			})
			require.Equal(t, nil, err)
		}
	}
}

func TestGetWorkbenchTemplate(t *testing.T) {
	items := getClient()
	ctx := context.Background()

	for _, item := range items {
		api := NewApi(item.client)
		{
			resp, err := api.GetWorkbenchTemplate(ctx, item.agentID)
			require.Equal(t, nil, err)
			fmt.Println("resp", resp)
		}
	}
}
