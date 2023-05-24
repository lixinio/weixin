package message_api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
)

func initWxa() *messageItem {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxa := official_account.New(redis, redis, &official_account.Config{
		Appid:  test.WxaAppid,
		Secret: test.WxaSecret,
	})

	return &messageItem{
		OpenID: test.OfficialAccountOpenid,
		Client: wxa.Client,
	}
}

/*
125 停车服务
884 维修保养
892 汽车经销商/4S店
894 汽车厂商
1129
131 车辆出场通知 3 125
0
1 车牌号 粤A12345
2 停车场 广州白云万达
4 停车费用 20元
5 提示说明 您的车已驶出停车场
6 停车时间 2015/1/1 12:00:00~2015/1/1 15:00:00
*/
func TestSubscribe(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*messageItem{
		initWxa(),
	} {
		messageApi := NewApi(client.Client)
		categories, err := messageApi.SubscribeGetCategory(ctx)
		require.Equal(t, nil, err)

		categoryIDs := []string{}
		for _, category := range categories {
			fmt.Println(category.ID, category.Name)
			categoryIDs = append(categoryIDs, strconv.Itoa(category.ID))
		}

		subscribePubTemplate, total, err := messageApi.SubscribeGetPubTemplateTitles(
			ctx, strings.Join(categoryIDs, ","), 1, 30)
		require.Equal(t, nil, err)
		fmt.Println(total)
		for _, tmpl := range subscribePubTemplate {
			fmt.Println(tmpl.Tid, tmpl.Title, tmpl.Type, tmpl.CategoryID)
		}

		keywords, total, err := messageApi.SubscribeGetPubTemplateKeywords(
			ctx, subscribePubTemplate[0].Tid,
		)
		require.Equal(t, nil, err)
		fmt.Println(total)
		for _, keyword := range keywords {
			fmt.Println(keyword.Kid, keyword.Name, keyword.Example)
		}
	}
}

/*
690 洗车提醒 2 884

1 车牌号 京TIN68
2 洗车指数 较适宜洗车
3 服务项目 蜡洗
4 客服热线 4008822339
5 洗车地址 软件产业基地5D2栋
6 温馨提示 设备即将为您进行洗车，请注意安全
7 洗车时间 2020-07-21 15:20
*/
func TestSubscribeTemplate(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*messageItem{
		initWxa(),
	} {
		messageApi := NewApi(client.Client)

		tmplates, err := messageApi.SubscribeGetTemplate(ctx)
		require.Equal(t, nil, err)

		for _, tmpl := range tmplates {
			fmt.Println(tmpl.PriTmplID, tmpl.Type, tmpl.Title, tmpl.Content, tmpl.Example)
		}

		priTmplID, err := messageApi.SubscribeAddTemplate(ctx, 690, []int{1, 3, 6}, "提醒车主洗车")
		require.Equal(t, nil, err)
		fmt.Println(priTmplID)

		err = messageApi.SubscribeDelTemplate(ctx, priTmplID)
		require.Equal(t, nil, err)
	}
}
