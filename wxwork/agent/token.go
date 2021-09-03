package agent

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

/*
从微信服务器获取新的 AccessToken

See: https://developers.weixin.qq.com/doc/corporation/Basic_Information/Get_access_token.html
*/
func (agent *Agent) refreshAccessTokenFromWXServer() (accessToken string, expiresIn int, err error) {
	var result utils.TokenResponse
	if err := agent.Client.HTTPGetToken(context.TODO(), "/cgi-bin/gettoken", func(params url.Values) {
		params.Add("corpid", agent.wxwork.Config.Corpid)
		params.Add("corpsecret", agent.Config.Secret)
	}, &result); err != nil {
		return "", 0, err
	}
	return result.AccessToken, result.ExpiresIn, nil
}
