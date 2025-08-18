package wxopen

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiJscode2Session = "/sns/component/jscode2session"
)

type MpSession struct {
	utils.WeixinError
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	SessionKey string `json:"session_key"`
}

// 小程序登录
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/others/WeChat_login.html
func (api *WxOpen) Jscode2Session(
	ctx context.Context,
	authorizerAppID, jsCode string,
) (*MpSession, error) {
	// 无需 access token
	result := &MpSession{}
	if err := api.Client.HTTPGetWithParams(
		ctx,
		apiJscode2Session,
		func(params url.Values) {
			params.Add("appid", authorizerAppID)
			params.Add("component_appid", api.Config.Appid)
			params.Add("js_code", jsCode)
			params.Add("grant_type", "authorization_code")
		}, result); err != nil {
		return nil, err
	}
	return result, nil
}
