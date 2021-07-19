package open

import (
	"fmt"

	"github.com/lixinio/weixin/utils"
)

var (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

type ComponentVerifyTicketGetter func(component_appid string) string

/*
公众号配置
*/
type Config struct {
	ComponentAppid  string
	ComponentSecret string
}

type Open struct {
	Config                         *Config
	Client                         *utils.Client
	component_verify_ticket_getter ComponentVerifyTicketGetter
}

func New(cache utils.Cache, config *Config, component_verify_ticket_getter ComponentVerifyTicketGetter) *Open {
	instance := &Open{
		Config:                         config,
		component_verify_ticket_getter: component_verify_ticket_getter,
	}
	instance.Client = utils.NewClient(WXServerUrl, utils.NewAccessTokenCache(instance, cache, 0))
	return instance
}

// GetAccessToken 接口 weixin.AccessTokenGetter 实现
func (open *Open) GetAccessToken() (accessToken string, expiresIn int, err error) {
	accessToken, expiresIn, err = open.refreshAccessTokenFromWXServer()
	return
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (open *Open) GetAccessTokenKey() string {
	return fmt.Sprintf(
		"access-token:wxopen:%s",
		open.Config.ComponentAppid,
	)
}
