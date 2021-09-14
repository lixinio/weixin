package main

import (
	"fmt"
	"time"

	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/wxopen"
)

func RefreshWxOpenToken(wxOpen *wxopen.WxOpen) {
	for {
		endTime := time.After(10 * time.Second)
		<-endTime
		token, err := wxOpen.RefreshAccessToken(0)
		if err != nil {
			fmt.Println("refresh wxopen token fail", err)
		} else {
			fmt.Printf("refresh token success '%s'\n", token)
		}
	}
}

func RefreshAuthorizerToken(authorizers []*authorizer.Authorizer) {
	for {
		endTime := time.After(10 * time.Second)
		<-endTime
		for _, auth := range authorizers {
			token, err := auth.RefreshAccessToken(0)
			if err != nil {
				fmt.Printf(
					"refresh authorizer(%s %s) fail, error %s\n",
					auth.ComponentAppid, auth.Appid, err.Error(),
				)
			} else {
				fmt.Printf(
					"refresh authorizer(%s %s) token success '%s'\n",
					auth.ComponentAppid, auth.Appid, token,
				)
			}
		}
	}
}
