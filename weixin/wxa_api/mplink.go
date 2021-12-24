package wxa_api

import (
	"context"
	"net/url"
	"strconv"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGetShowWxaItem      = "/wxa/getshowwxaitem"
	apiGetWxaMplinkForShow = "/wxa/getwxamplinkforshow"
	apiUpdateShowWxaItem   = "/wxa/updateshowwxaitem"
)

// 公众号关注组件。当用户扫小程序码打开小程序时，开发者可在小程序内配置公众号关注组件，方便用户快捷关注公众号，可嵌套在原生组件内。
// https://developers.weixin.qq.com/miniprogram/dev/component/official-account.html

/**
获取展示的公众号信息
使用本接口可以获取扫码关注组件所展示的公众号信息
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/subscribe_component/getshowwxaitem.html
**/
type WxaMpInfo struct {
	AppID    string `json:"appid"`
	NickName string `json:"nickname"`
	HeadImg  string `json:"headimg"`
}

type ShowWxaMpInfo struct {
	utils.WeixinError
	CanOpen int `json:"can_open"`
	IsOpen  int `json:"is_open"`
	WxaMpInfo
}

func (api *WxaApi) GetShowWxaItem(ctx context.Context) (*ShowWxaMpInfo, error) {
	result := &ShowWxaMpInfo{}
	if err := api.Client.HTTPGet(ctx, apiGetShowWxaItem, result); err != nil {
		return nil, err
	}
	return result, nil
}

/**
获取可以用来设置的公众号列表
通过本接口可以获取扫码关注组件允许展示的公众号列表
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/subscribe_component/getwxamplinkforshow.html
**/
type ShowWxaMpList struct {
	utils.WeixinError
	TotalNum    int         `json:"total_num"`
	BizInfoList []WxaMpInfo `json:"biz_info_list"`
}

func (api *WxaApi) GetWxaMplinkForShow(
	ctx context.Context, page, count int,
) (*ShowWxaMpList, error) {
	result := &ShowWxaMpList{}
	if err := api.Client.HTTPGetWithParams(
		ctx, apiGetShowWxaItem, func(query url.Values) {
			query.Add("page", strconv.Itoa(page))
			if count > 0 {
				count = 20
			}
			query.Add("num", strconv.Itoa(count))
		}, result,
	); err != nil {
		return nil, err
	}
	return result, nil
}

/**
使用本接口可以设置扫码关注组件所展示的公众号信息
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/subscribe_component/updateshowwxaitem.html
**/
func (api *WxaApi) UpdateShowWxaItem(
	ctx context.Context, subscribe int, appid string,
) error {
	param := &struct {
		Subscribe int    `json:"wxa_subscribe_biz_flag"`
		AppID     string `json:"appid"`
	}{subscribe, appid}
	if err := api.Client.HTTPPostJson(ctx, apiUpdateShowWxaItem, param, nil); err != nil {
		return err
	}
	return nil
}
