package authorizer

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

// 开放平台
const (
	apiWxOpenCreate = "/cgi-bin/open/create"
	apiWxOpenGet    = "/cgi-bin/open/get"
	apiWxOpenBind   = "/cgi-bin/open/bind"
	apiWxOpenUnbind = "/cgi-bin/open/unbind"
	apiWxOpenHave   = "/cgi-bin/open/have"
	apiRidGet       = "/cgi-bin/openapi/rid/get"
)

// 创建开放平台帐号并绑定公众号/小程序
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/account/create.html
func (api *Authorizer) WxOpenCreate(
	ctx context.Context, appid string,
) (string, error) {
	params := map[string]string{
		"appid": appid,
	}
	result := struct {
		utils.WeixinError
		OpenAppID string `json:"open_appid"`
	}{}
	err := api.Client.HTTPPostJson(ctx, apiWxOpenCreate, params, &result)
	if err != nil {
		return "", err
	}
	return result.OpenAppID, nil
}

// 将公众号/小程序绑定到开放平台帐号下
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/account/bind.html
func (api *Authorizer) WxOpenBind(
	ctx context.Context, appid, openAppid string,
) error {
	params := map[string]string{
		"appid":      appid,
		"open_appid": openAppid,
	}
	return api.Client.HTTPPostJson(ctx, apiWxOpenBind, params, nil)
}

// 将公众号/小程序从开放平台帐号下解绑
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/account/unbind.html
func (api *Authorizer) WxOpenUnBind(
	ctx context.Context, appid, openAppid string,
) error {
	params := map[string]string{
		"appid":      appid,
		"open_appid": openAppid,
	}
	return api.Client.HTTPPostJson(ctx, apiWxOpenUnbind, params, nil)
}

// 获取公众号/小程序所绑定的开放平台帐号
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/account/get.html
func (api *Authorizer) WxOpenGet(
	ctx context.Context, appid string,
) (string, error) {
	params := map[string]string{
		"appid": appid,
	}
	result := struct {
		utils.WeixinError
		OpenAppID string `json:"open_appid"`
	}{}
	err := api.Client.HTTPPostJson(ctx, apiWxOpenGet, params, &result)
	if err != nil {
		return "", err
	}
	return result.OpenAppID, nil
}

// 查询公众号/小程序是否绑定open帐号
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/getbindopeninfo.html
func (api *Authorizer) WxOpenHave(ctx context.Context) (bool, error) {
	result := struct {
		utils.WeixinError
		HaveOpen bool `json:"have_open"`
	}{}
	err := api.Client.HTTPGet(ctx, apiWxOpenHave, &result)
	if err != nil {
		return false, err
	}
	return result.HaveOpen, nil
}

// 查询rid信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/openApi/get_rid_info.html
type RidResult struct {
	InvokeTime   int64  `json:"invoke_time"`
	CostInMs     int    `json:"cost_in_ms"`
	RequestUrl   string `json:"request_url"`
	RequestBody  string `json:"request_body"`
	ResponseBody string `json:"response_body"`
}

func (api *Authorizer) RidGet(
	ctx context.Context, rid string,
) (*RidResult, error) {
	params := map[string]string{
		"rid": rid,
	}
	result := struct {
		utils.WeixinError
		Request *RidResult `json:"request"`
	}{}
	err := api.Client.HTTPPostJson(ctx, apiRidGet, params, &result)
	if err != nil {
		return nil, err
	}
	return result.Request, nil
}
