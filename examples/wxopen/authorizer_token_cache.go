package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

type authorizerAccessTokenAdaptor struct {
	componentAppid string
	appid          string
}

func (ta *authorizerAccessTokenAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	return "", 0, errors.New("can NOT update authorizer access token")
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerAccessTokenAdaptor) GetAccessTokenKey() string {
	return fmt.Sprintf("access-token.authorizer.%s.%s", ta.componentAppid, ta.appid)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerAccessTokenAdaptor) GetAccessTokenLockKey() string {
	return fmt.Sprintf("access-token.authorizer.%s.%s.lock", ta.componentAppid, ta.appid)
}

type authorizerRefreshTokenAdaptor struct {
	componentAppid string
	appid          string
}

func (ta *authorizerRefreshTokenAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	return "", 0, errors.New("can NOT refresh authorizer refresh code")
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerRefreshTokenAdaptor) GetAccessTokenKey() string {
	return fmt.Sprintf("refresh-token.authorizer.%s.%s", ta.componentAppid, ta.appid)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerRefreshTokenAdaptor) GetAccessTokenLockKey() string {
	return fmt.Sprintf("refresh-token.authorizer.%s.%s.lock", ta.componentAppid, ta.appid)
}

type AuthorizerTokenCache struct {
	accessToken  *utils.AccessTokenCache
	refreshToken *utils.AccessTokenCache
}

func (tc *AuthorizerTokenCache) SetAccessToken(
	ctx context.Context, token string, expiresIn int,
) error {
	_, err := tc.accessToken.UpdateAccessToken(ctx, token, expiresIn)
	return err
}

func (tc *AuthorizerTokenCache) GetAccessToken(
	ctx context.Context,
) (string, error) {
	return tc.accessToken.GetAccessToken(ctx)
}

func (tc *AuthorizerTokenCache) SetRefreshToken(
	ctx context.Context, token string,
) error {
	_, err := tc.refreshToken.UpdateAccessToken(ctx, token, 3600*24*365) // 一年， 长期有效
	return err
}

func (tc *AuthorizerTokenCache) GetRefreshToken(
	ctx context.Context,
) (string, error) {
	return tc.refreshToken.GetAccessToken(ctx)
}

func NewAuthorizerTokenCache(
	componentAppid, appid string,
	cache utils.Cache,
	locker utils.Lock,
) *AuthorizerTokenCache {
	return &AuthorizerTokenCache{
		accessToken: utils.NewAccessTokenCache(&authorizerAccessTokenAdaptor{
			componentAppid: componentAppid, appid: appid,
		}, cache, locker),
		refreshToken: utils.NewAccessTokenCache(&authorizerRefreshTokenAdaptor{
			componentAppid: componentAppid, appid: appid,
		}, cache, locker),
	}
}
