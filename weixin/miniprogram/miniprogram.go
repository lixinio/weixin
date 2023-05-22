package miniprogram

import (
	"context"
	"github.com/lixinio/weixin/utils"
	"net/url"
)

const (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

type Config struct {
	Appid  string
	Secret string
}

type MiniProgram struct {
	Config *Config
	Client *utils.Client
}

func New(cache utils.Cache, locker utils.Lock, config *Config) *MiniProgram {
	instance := &MiniProgram{
		Config: config,
	}

	instance.Client = utils.NewClient(
		WXServerUrl,
		utils.NewAccessTokenCache(
			newAdapter(config.Appid, instance.refreshAccessTokenFromWXServer),
			cache, locker,
		),
	)

	return instance
}

func (m *MiniProgram) refreshAccessTokenFromWXServer() (accessToken string, expiresIn int, err error) {
	var result utils.TokenResponse
	if err := m.Client.HTTPGetToken(context.TODO(), "/cgi-bin/token", func(params url.Values) {
		params.Add("appid", m.Config.Appid)
		params.Add("secret", m.Config.Secret)
		params.Add("grant_type", "client_credential")
	}, &result); err != nil {
		return "", 0, err
	}

	return result.AccessToken, result.ExpiresIn, nil
}
