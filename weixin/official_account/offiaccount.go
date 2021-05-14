package official_account

import (
	"fmt"

	"github.com/lixinio/weixin"
)

var (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

/*
公众号配置
*/
type Config struct {
	Appid          string
	Secret         string
	Token          string
	EncodingAESKey string
}

type OfficialAccount struct {
	Config *Config
	Client *weixin.Client
}

func New(cache weixin.Cache, config *Config) *OfficialAccount {
	instance := &OfficialAccount{
		Config: config,
	}
	instance.Client = weixin.NewClient(WXServerUrl, weixin.NewAccessTokenCache(instance, cache, 0))
	return instance
}

// GetAccessToken 接口 weixin.AccessTokenGetter 实现
func (officialAccount *OfficialAccount) GetAccessToken() (accessToken string, expiresIn int, err error) {
	accessToken, expiresIn, err = officialAccount.refreshAccessTokenFromWXServer()
	return
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (officialAccount *OfficialAccount) GetAccessTokenKey() string {
	return fmt.Sprintf(
		"access-token:officialaccount:%s",
		officialAccount.Config.Appid,
	)
}
