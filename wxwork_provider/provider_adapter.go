package wxwork_provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

var (
	ErrTokenUpdateForbidden = errors.New("can NOT refresh&update token in wxwork provider lite mode")
)

// access token的cache
type accessTokenAdaptor struct {
	config    *Config
	client    *utils.Client
	tokenKey  string
	lockerKey string
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *accessTokenAdaptor) GetAccessTokenKey() string {
	return ta.tokenKey
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *accessTokenAdaptor) GetAccessTokenLockKey() string {
	return ta.lockerKey
}

func (ta *accessTokenAdaptor) GetAccessToken() (accessToken string, expiresIn int, err error) {
	if ta.config.ProviderSecret == "" {
		return "", 0, fmt.Errorf(
			"wxopen appid : %s, error: %w", ta.config.CorpID, ErrTokenUpdateForbidden,
		)
	}

	// AccessToken 和其他地方 字段不一致
	result := struct {
		utils.WeixinError
		AccessToken string `json:"provider_access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}{}

	payload := map[string]string{
		"corpid":          ta.config.CorpID,
		"provider_secret": ta.config.ProviderSecret,
	}
	if err := ta.client.HTTPPostToken(
		context.TODO(), apiGetProviderToken, payload, &result,
	); err != nil {
		return "", 0, err
	}
	return result.AccessToken, result.ExpiresIn, nil
}

func newAccessTokenAdaptor(config *Config) *accessTokenAdaptor {
	adaptor := &accessTokenAdaptor{
		config:    config,
		client:    utils.NewClient(WXServerUrl, utils.EmptyClientAccessTokenGetter(0)),
		tokenKey:  fmt.Sprintf("qywx.provider_access_token.%s", config.CorpID),
		lockerKey: fmt.Sprintf("qywx.provider_access_token.%s.lock", config.CorpID),
	}
	return adaptor
}
