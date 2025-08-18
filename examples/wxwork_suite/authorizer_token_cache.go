package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

type authorizerAccessTokenAdaptor struct {
	suiteID string
	corpID  string
	agentID int
}

func (ta *authorizerAccessTokenAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	return "", 0, errors.New("can NOT update authorizer access token")
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerAccessTokenAdaptor) GetAccessTokenKey() string {
	return fmt.Sprintf("access-token.suite-authorizer.%s.%s.%d", ta.suiteID, ta.corpID, ta.agentID)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerAccessTokenAdaptor) GetAccessTokenLockKey() string {
	return fmt.Sprintf(
		"access-token.suite-authorizer.%s.%s.%d.lock",
		ta.suiteID,
		ta.corpID,
		ta.agentID,
	)
}

type authorizerPermanentCodeAdaptor struct {
	suiteID string
	corpID  string
	agentID int
}

func (ta *authorizerPermanentCodeAdaptor) GetAccessToken(
	ctx context.Context,
) (accessToken string, expiresIn int, err error) {
	return "", 0, errors.New("can NOT refresh authorizer permanent code")
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerPermanentCodeAdaptor) GetAccessTokenKey() string {
	return fmt.Sprintf(
		"permanent-code.suite-authorizer.%s.%s.%d",
		ta.suiteID,
		ta.corpID,
		ta.agentID,
	)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (ta *authorizerPermanentCodeAdaptor) GetAccessTokenLockKey() string {
	return fmt.Sprintf(
		"permanent-code.suite-authorizer.%s.%s.%d.lock",
		ta.suiteID,
		ta.corpID,
		ta.agentID,
	)
}

type AuthorizerTokenCache struct {
	accessToken   *utils.AccessTokenCache
	permanentCode *utils.AccessTokenCache
}

func (tc *AuthorizerTokenCache) SetAccessToken(
	ctx context.Context, token string, expiresIn int,
) error {
	_, err := tc.accessToken.UpdateAccessToken(ctx, token, expiresIn)
	return err
}

func (tc *AuthorizerTokenCache) GetAccessToken(ctx context.Context) (string, error) {
	return tc.accessToken.GetAccessToken(ctx)
}

func (tc *AuthorizerTokenCache) SetPermanentCode(
	ctx context.Context, token string,
) error {
	_, err := tc.permanentCode.UpdateAccessToken(ctx, token, 3600*24*365) // 一年， 长期有效
	return err
}

func (tc *AuthorizerTokenCache) GetPermanentCode(ctx context.Context) (string, error) {
	return tc.permanentCode.GetAccessToken(ctx)
}

func NewAuthorizerTokenCache(
	suiteID, corpID string, agentID int,
	cache utils.Cache, locker utils.Lock,
) *AuthorizerTokenCache {
	return &AuthorizerTokenCache{
		accessToken: utils.NewAccessTokenCache(&authorizerAccessTokenAdaptor{
			suiteID: suiteID, corpID: corpID, agentID: agentID,
		}, cache, locker),
		permanentCode: utils.NewAccessTokenCache(&authorizerPermanentCodeAdaptor{
			suiteID: suiteID, corpID: corpID, agentID: agentID,
		}, cache, locker),
	}
}
