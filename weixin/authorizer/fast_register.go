package authorizer

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

// 快速注册小程序

const (
	apiFastRegister         = "/cgi-bin/account/fastregister"
	apiComponentreBindAdmin = "/cgi-bin/account/componentrebindadmin"
)

// 复用公众号主体快速注册小程序
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Register_Mini_Programs/fast_registration_of_mini_program.html
func (api *Authorizer) GetFastRegisterAuthUrl(
	copyWxVerify, redirectUri string,
) string {
	params := url.Values{}
	params.Add("component_appid", api.ComponentAppid)
	params.Add("appid", api.Appid)
	params.Add("copy_wx_verify", copyWxVerify) // 是否复用公众号的资质进行微信认证(1:申请复用资质进行微信 认证 0:不申请)
	// 用户扫码授权后，MP 扫码页面将跳转到该地址(注:1.链接需 urlencode 2.Host 需和第三方平台在微信开放平台上面填写的登 录授权的发起页域名一致)
	params.Add("redirect_uri", redirectUri)

	return "https://mp.weixin.qq.com/cgi-bin/fastregisterauth?" + params.Encode()
}

// 快速注册 API 完成注册
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Register_Mini_Programs/fast_registration_of_mini_program.html
type FastRegisterResult struct {
	utils.WeixinError
	AppID             string `json:"appid"`              // 新创建小程序的 appid
	AuthorizationCode string `json:"authorization_code"` // 新创建小程序的授权码
	IsWxVerifySucc    bool   `json:"is_wx_verify_succ"`  // 复用公众号微信认证小程序是否成功
	IsLinkSucc        bool   `json:"is_link_succ"`       // 小程序是否和公众号关联成功
}

func (api *Authorizer) FastRegister(
	ctx context.Context, ticket string,
) (*FastRegisterResult, error) {
	params := map[string]string{
		"ticket": ticket,
	}
	var result FastRegisterResult
	err := api.Client.HTTPPostJson(ctx, apiFastRegister, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 换绑小程序管理员接口
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Register_Mini_Programs/fast_registration_of_mini_program.html
func (api *Authorizer) GetComponentreBindAdminUrl(
	copyWxVerify, redirectUri string,
) string {
	params := url.Values{}
	params.Add("component_appid", api.ComponentAppid)
	params.Add("appid", api.Appid)
	// 新管理员信息填写完成点击提交后，将跳转到该地址(注：1.链接需 urlencode 2.Host 需和第三方平台在微信开放平台上面填写的登录授权的发起页域名一致)
	params.Add("redirect_uri", redirectUri)

	return "https://mp.weixin.qq.com/wxopen/componentrebindadmin?" + params.Encode()
}

// 快速注册 API 完成管理员换绑
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Register_Mini_Programs/fast_registration_of_mini_program.html
func (api *Authorizer) ComponentreBindAdmin(
	ctx context.Context, taskID string,
) error {
	params := map[string]string{
		"taskid": taskID,
	}
	return api.Client.HTTPPostJson(ctx, apiComponentreBindAdmin, params, nil)
}
