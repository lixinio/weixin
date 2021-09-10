package main

import (
	"context"

	"github.com/lixinio/weixin/wxwork/authorizer"
	"github.com/lixinio/weixin/wxwork_suite"
)

func GetAuthorizerAccessToken(
	suite *wxwork_suite.WxWorkSuite,
	tokenCache TokenCache,
	corpID string,
) authorizer.RefreshAccessToken {
	return func() (string, int, error) {
		permanentCode, err := tokenCache.GetPermanentCode()
		if err != nil {
			return "", 0, err
		}
		resp, err := suite.GetCorpToken(
			context.TODO(),
			corpID,
			permanentCode,
		)
		if err != nil {
			return "", 0, err
		}
		return resp.AccessToken, resp.ExpiresIn, nil
	}
}
