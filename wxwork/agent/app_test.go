package agent

import (
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := New(corp, cache, &Config{
		AgentId: test.AgentID,
		Secret:  test.AgentSecret,
	})

	menu := []MenuEntryObj{
		{
			Name: "常规",
			SubButton: []*MenuEntryObj{
				{
					Type: MenuTypeView,
					Name: "百度",
					Url:  "https://www.baidu.com",
				},
				{
					Type: MenuTypeClick,
					Name: "点击",
					Key:  "V1001_TODAY_MUSIC",
				},
				{
					Type: MenuTypeLocationSelect,
					Name: "位置",
					Key:  "MenuTypeLocationSelect",
				},
				// {
				// 	Type:  MenuTypeViewMiniPrograme,
				// 	Name:  "小程序",
				// 	AppID: "wx0dfb2786bfd2b513",
				// },
			},
		},
		{
			Name: "扫码",
			SubButton: []*MenuEntryObj{
				{
					Type: MenuTypeScanCodeWaitmsg,
					Name: "扫码带Msg",
					Key:  "MenuTypeScanCodeWaitmsg",
				},
				{
					Type: MenuTypeScanCodePush,
					Name: "扫码带Push",
					Key:  "MenuTypeScanCodePush",
				},
				{
					Type: MenuTypePicSysPhoto,
					Name: "系统选图",
					Key:  "MenuTypePicSysPhoto",
				},
				{
					Type: MenuTypePicPhotoOrAlbum,
					Name: "系统选图2",
					Key:  "MenuTypePicPhotoOrAlbum",
				},
				{
					Type: MenuTypePicWeixin,
					Name: "系统选图3",
					Key:  "MenuTypePicWeixin",
				},
			},
		},
	}

	require.Equal(t, nil, agent.MenuDelete(agent.Config.AgentId))
	require.Equal(t, nil, agent.MenuCreate(agent.Config.AgentId, menu))
}
