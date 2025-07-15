package official_account

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

/*
从微信服务器获取新的 AccessToken
See: https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
See: https://github.com/fastwego/offiaccount/blob/master/client.go#L268
*/
func (officialAccount *OfficialAccount) refreshAccessTokenFromWXServer(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	var result utils.TokenResponse
	if err := officialAccount.Client.HTTPGetToken(
		utils.NewStripContext(ctx, "secret"),
		"/cgi-bin/token",
		func(params url.Values) {
			params.Add("appid", officialAccount.Config.Appid)
			params.Add("secret", officialAccount.Config.Secret)
			params.Add("grant_type", "client_credential")
		},
		&result,
	); err != nil {
		return "", 0, err
	}
	return result.AccessToken, result.ExpiresIn, nil
}
