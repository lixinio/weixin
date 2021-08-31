package authorizer

import (
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

// 需要通过wxopen对象刷新authorizer Access token
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/api_authorizer_token.html
type RefreshAccessToken func() (string, int, error) // 直接获取token， 不做任何缓存

// utils.AccessTokenGetter 接口实现
type authorizerAccessTokenGetterAdapter struct {
	accessTokenKey     string
	accessTokenLockKey string
	accessTokenGetter  RefreshAccessToken
}

// GetAccessToken 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessToken() (string, int, error) {
	return adapter.accessTokenGetter()
}

// GetAccessTokenKey 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessTokenKey() string {
	return adapter.accessTokenKey
}

// GetAccessTokenLockKey 接口 utils.AccessTokenGetter 实现
func (adapter *authorizerAccessTokenGetterAdapter) GetAccessTokenLockKey() string {
	return adapter.accessTokenLockKey
}

func newAdapter(
	componentAppid, appid string,
	accessTokenGetter RefreshAccessToken,
) utils.AccessTokenGetter {
	return &authorizerAccessTokenGetterAdapter{
		accessTokenGetter: accessTokenGetter,
		accessTokenKey: fmt.Sprintf(
			"access-token:authorizer:%s:%s",
			componentAppid,
			appid,
		),
		accessTokenLockKey: fmt.Sprintf(
			"access-token:authorizer:%s:%s.lock",
			componentAppid,
			appid,
		),
	}
}

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
		newAdapter(componentAppid, appid, accessTokenGetter),
		cache,
		locker,
		0,
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
				"can NOT refresh token in bare mod, appid(%s , %s)",
				componentAppid,
				appid,
			)
		}),
		cache,
		locker,
		0,
	)
	return &Authorizer{
		ComponentAppid: componentAppid,
		Appid:          appid,
		Client:         utils.NewClient(WXServerUrl, accessTokenCache),
	}
}

func (authorizer *Authorizer) RefreshAccessToken() error {
	if authorizer.accessTokenCache != nil {
		_, err := authorizer.accessTokenCache.GetAccessToken()
		return err
	} else {
		// bare模式，只能从cache读取token， 无法刷新
		return fmt.Errorf(
			"can NOT refresh token in bare mod, appid(%s , %s)",
			authorizer.ComponentAppid,
			authorizer.Appid,
		)
	}
}
