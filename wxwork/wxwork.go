package wxwork

import "github.com/lixinio/weixin/utils"

var (
	QyWXServerUrl = "https://qyapi.weixin.qq.com"
	UserAgent     = "lixinio/wxwork"
)

type Config struct {
	Corpid string // 企业ID
}

type WxWork struct {
	Config *Config
}

func New(config *Config) (corp *WxWork) {
	instance := WxWork{
		Config: config,
	}
	return &instance
}

func (corp *WxWork) NewClient(accessTokenCache *utils.AccessTokenCache) *utils.Client {
	return utils.NewClient(QyWXServerUrl, accessTokenCache)
}
