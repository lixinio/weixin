package message_api

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/lixinio/weixin/wxwork/authorizer"
	"github.com/lixinio/weixin/wxwork/material_api"
	"github.com/lixinio/weixin/wxwork/user_api"
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
		userApi := user_api.NewApi(cli.client)
		userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
		require.Equal(t, nil, err)

		messageApi := NewApi(cli.client, cli.agentID)

		result, err := messageApi.SendTextMessage(
			ctx,
			NewMessageHeaderByUsers([]string{userid}),
			"你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。",
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)

		result, err = messageApi.SendTextCardMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			"领奖通知",
			"<div class=\"gray\">2016年9月26日</div> <div class=\"normal\">恭喜你抽中iPhone 7一台，领奖码：xxxx</div><div class=\"highlight\">请于2016年10月10日前联系行政同事领取</div>",
			"https://www.baidu.com",
			"更多",
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)

		result, err = messageApi.SendNewsMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			[]*NewsMessageParam{
				{
					Title:       "中秋节礼品领取",
					Description: "今年中秋节公司有豪礼相送",
					URL:         "https://www.baidu.com",
					PicURL:      "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
				},
				{
					Title:       "端午节礼品领取",
					Description: "今年端午节公司有豪礼相送",
					URL:         "https://www.baidu.com",
					PicURL:      "http://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png",
				},
			},
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)

		result, err = messageApi.SendMarkdownMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			`
您的会议室已经预定，稍后会同步到**邮箱**
>**事项详情** 
>事　项：<font color=\"info\">开会</font> 
>组织者：@miglioguan 
>参与者：@miglioguan、@kunliu、@jamdeezhou、@kanexiong、@kisonwang 
> 
>会议室：<font color=\"info\">广州TIT 1楼 301</font> 
>日　期：<font color=\"warning\">2018年5月18日</font> 
>时　间：<font color=\"comment\">上午9:00-11:00</font> 
> 
>请准时参加会议。 
> 
>如需修改会议信息，请点击：[修改会议信息](https://work.weixin.qq.com)
			`,
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}
}

func TestSendImageMessage(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkAgent(),
		initWxWorkSuiteAuthorizer(),
	} {
		materialApi := material_api.NewApi(cli.client)
		messageApi := NewApi(cli.client, cli.agentID)

		userApi := user_api.NewApi(cli.client)
		userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
		require.Equal(t, nil, err)
		file, err := os.Open(test.ImagePath)
		require.Empty(t, err)
		defer file.Close()

		materialResult, err := materialApi.Upload(ctx, test.ImagePath, file, material_api.MediaTypeImage)
		require.Empty(t, err)
		require.Equal(t, materialResult.Type, material_api.MediaTypeImage)

		result, err := messageApi.SendImageMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			materialResult.MediaID,
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)

		result, err = messageApi.SendFileMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			materialResult.MediaID,
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}
}

func TestSendVideoMessage(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkAgent(),
		initWxWorkSuiteAuthorizer(),
	} {
		materialApi := material_api.NewApi(cli.client)
		messageApi := NewApi(cli.client, cli.agentID)

		userApi := user_api.NewApi(cli.client)
		userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
		require.Equal(t, nil, err)
		file, err := os.Open(test.VideoPath)
		require.Empty(t, err)
		defer file.Close()

		materialResult, err := materialApi.Upload(ctx, test.VideoPath, file, material_api.MediaTypeVideo)
		require.Empty(t, err)
		require.Equal(t, materialResult.Type, material_api.MediaTypeVideo)

		result, err := messageApi.SendVideoMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			materialResult.MediaID,
			"",
			"",
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)

		result, err = messageApi.SendFileMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			materialResult.MediaID,
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}
}

func TestSendAudioMessage(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkAgent(),
		initWxWorkSuiteAuthorizer(),
	} {
		materialApi := material_api.NewApi(cli.client)
		messageApi := NewApi(cli.client, cli.agentID)

		userApi := user_api.NewApi(cli.client)
		userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
		require.Equal(t, nil, err)
		file, err := os.Open(test.AudioPath)
		require.Empty(t, err)
		defer file.Close()

		materialResult, err := materialApi.Upload(ctx, test.AudioPath, file, material_api.MediaTypeVoice)
		require.Empty(t, err)
		require.Equal(t, materialResult.Type, material_api.MediaTypeVoice)

		result, err := messageApi.SendVoiceMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			materialResult.MediaID,
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)

		result, err = messageApi.SendFileMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			materialResult.MediaID,
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}
}

func TestSendMpNewsMessage(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkAgent(),
		initWxWorkSuiteAuthorizer(),
	} {
		materialApi := material_api.NewApi(cli.client)
		messageApi := NewApi(cli.client, cli.agentID)

		userApi := user_api.NewApi(cli.client)
		userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
		require.Equal(t, nil, err)
		file, err := os.Open(test.ImagePath)
		require.Empty(t, err)
		defer file.Close()

		materialResult, err := materialApi.Upload(ctx, test.ImagePath, file, material_api.MediaTypeImage)
		require.Empty(t, err)
		require.Equal(t, materialResult.Type, material_api.MediaTypeImage)

		file.Seek(0, 0)
		url, err := materialApi.UploadImg(ctx, test.ImagePath, file)
		require.Empty(t, err)
		fmt.Print(url)

		result, err := messageApi.SendMpNewsMessage(
			ctx,
			NewMessageHeaderByUser(userid),
			[]*MpNewsMessageParam{
				{
					Title:            "Title",
					ThumbMediaID:     materialResult.MediaID,
					Author:           "Author",
					ContentSourceURL: "https://www.baidu.com",
					Content:          fmt.Sprintf("Content <img src=\"%s\">", url),
					Digest:           "Digest description",
				},
			},
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}
}

func TestSendMpNoticeMessage(t *testing.T) {
	ctx := context.Background()
	for _, cli := range []*item{
		initWxWorkSuiteAuthorizer(),
	} {
		userApi := user_api.NewApi(cli.client)
		userid, err := userApi.MobileGetUserId(ctx, test.AgentUserMobile)
		require.Equal(t, nil, err)

		messageApi := NewApi(cli.client, cli.agentID)

		result, err := messageApi.SendMpNoticeMessage(
			ctx,
			NewMessageHeaderByUsers([]string{userid}),
			&MpNoticeMessageParam{
				AppID:             "wx123123123123123",
				Page:              "",
				Title:             "会议室预订成功通知",
				Description:       "4月27日 16:16",
				EmphasisFirstItem: true,
				ContentItem: []*MpNoticeItem{
					{
						Key:   "会议室",
						Value: "402",
					},
					{
						Key:   "会议地点",
						Value: "广州TIT-402会议室",
					},
					{
						Key:   "会议时间",
						Value: "2018年8月1日 09:00-09:30",
					},
					{
						Key:   "参与人员",
						Value: "周剑轩",
					},
				},
			},
		)
		require.Equal(t, nil, err)
		fmt.Println(result.MsgID)
	}
}
