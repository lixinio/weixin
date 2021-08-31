package main

import (
	"context"
	"fmt"

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
		fmt.Printf("refresh appid (%s) token (%s)\n", appid, resp.AccessToken)
		return resp.AccessToken, resp.ExpiresIn, nil
	}
}
