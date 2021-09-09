package authorizer

import (
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

var (
	ErrTokenUpdateForbidden = errors.New("can NOT refresh&update token in wxopen lite mode")
)

type Authorizer struct {
	ComponentAppid   string
	Appid            string
	Client           *utils.Client
	accessTokenCache *utils.AccessTokenCache // 用于支持手动刷新Token， Client不对外暴露该对象
}

func New(
	cache utils.Cache,
	locker utils.Lock,
	componentAppid, appid string,
	accessTokenGetter RefreshAccessToken,
) *Authorizer {
	accessTokenCache := utils.NewAccessTokenCache(
		newAdapter(componentAppid, appid, accessTokenGetter), cache, locker,
	)
	return &Authorizer{
		ComponentAppid:   componentAppid,
		Appid:            appid,
		Client:           utils.NewClient(WXServerUrl, accessTokenCache),
		accessTokenCache: accessTokenCache,
	}
}

// 不支持刷新Token的简化模式
func NewLite(
	cache utils.Cache,
	locker utils.Lock,
	componentAppid, appid string,
) *Authorizer {
	accessTokenCache := utils.NewAccessTokenCache(
		newAdapter(componentAppid, appid, func() (string, int, error) {
			return "", 0, fmt.Errorf(
				"can NOT refresh token in lite mod, appid(%s , %s), %w",
				componentAppid, appid, ErrTokenUpdateForbidden,
			)
		}), cache, locker,
	)
	return &Authorizer{
		ComponentAppid: componentAppid,
		Appid:          appid,
		Client:         utils.NewClient(WXServerUrl, accessTokenCache),
	}
}

func (authorizer *Authorizer) RefreshAccessToken(expireBefore int) (string, error) {
	if authorizer.accessTokenCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrTokenUpdateForbidden,
		)
	}
	return authorizer.accessTokenCache.RefreshAccessToken(expireBefore)
}
