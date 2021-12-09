package official_account

import (
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

var (
	ErrTokenUpdateForbidden = errors.New("can NOT refresh&update token in offiaccount lite mode")
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
		WXServerUrl, utils.NewAccessTokenCache(
			newAdapter(config.Appid, instance.refreshAccessTokenFromWXServer),
			cache, locker,
		),
	)
	return instance
}

func NewLite(cache utils.Cache, locker utils.Lock, appid string) *OfficialAccount {
	client := utils.NewClient(
		WXServerUrl, utils.NewAccessTokenCache(
			newAdapter(appid, func() (string, int, error) {
				return "", 0, fmt.Errorf(
					"can NOT refresh token in lite mod, appid(%s), %w",
					appid, ErrTokenUpdateForbidden,
				)
			}),
			cache, locker,
		),
	)
	return &OfficialAccount{
		Client: client,
		Config: &Config{
			Appid: appid,
		},
	}
}
