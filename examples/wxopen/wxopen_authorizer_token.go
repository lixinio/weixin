package main

import (
	"context"

	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/wxopen"
)

func GetAuthorizerAccessToken(
	wxOpen *wxopen.WxOpen,
	tokenCache TokenCache,
	appid string,
) authorizer.RefreshAccessToken {
	return func() (string, int, error) {
		refreshToken, err := tokenCache.GetRefreshToken()
		if err != nil {
			return "", 0, err
		}
		resp, err := wxOpen.GetAuthorizerToken(
			context.TODO(),
			appid,
			refreshToken,
		)
		if err != nil {
			return "", 0, err
		}
		tokenCache.SetRefreshToken(resp.RefreshToken) // noqa
		return resp.AccessToken, resp.ExpiresIn, nil
	}
}
