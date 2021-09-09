package official_account

import (
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

/*
公众号配置
*/
type Config struct {
	Appid  string
	Secret string
}

type OfficialAccount struct {
	Config *Config
	Client *utils.Client
}

func New(cache utils.Cache, locker utils.Lock, config *Config) *OfficialAccount {
	instance := &OfficialAccount{
		Config: config,
	}
	instance.Client = utils.NewClient(
		WXServerUrl,
		utils.NewAccessTokenCache(instance, cache, locker),
	)
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
		"weixin.access_token.%s",
		officialAccount.Config.Appid,
	)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (officialAccount *OfficialAccount) GetAccessTokenLockKey() string {
	return fmt.Sprintf(
		"weixin.access_token.%s.lock",
		officialAccount.Config.Appid,
	)
}
