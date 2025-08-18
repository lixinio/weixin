package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxopen"
)

func serveData(serverApi *wxopen.WxOpen, apps wxopen.ReleaseApps) utils.XmlHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, body []byte) error {
		_, content, err := serverApi.ParseXML(body)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return err
		}

		switch v := content.(type) {
		case *wxopen.EventComponentVerifyTicket:
			fmt.Printf("ComponentVerifyTicket : %s\n", v.ComponentVerifyTicket)
			// todo save to db
			// 刷新Ticket到Cache
			serverApi.UpdateTicket(r.Context(), v.ComponentVerifyTicket)
		case *wxopen.EventAuthorized:
			if app, ok := apps[v.AuthorizerAppid]; ok {
				fmt.Printf("release app %s %v\n", app.UserName, app.IsMp)
				// ignore
			}
			fmt.Printf("EventAuthorized : %s\n", v.AuthorizerAppid)
		case *wxopen.EventUnauthorized:
			if app, ok := apps[v.AuthorizerAppid]; ok {
				fmt.Printf("release app %s %v\n", app.UserName, app.IsMp)
				// ignore
			}
			fmt.Printf("EventUnauthorized : %s\n", v.AuthorizerAppid)
		case *wxopen.EventUpdateAuthorized:
			if app, ok := apps[v.AuthorizerAppid]; ok {
				fmt.Printf("release app %s %v\n", app.UserName, app.IsMp)
				// ignore
			}
			fmt.Printf("EventUpdateAuthorized : %s\n", v.AuthorizerAppid)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		_, err = io.WriteString(w, "success")
		return err
	}
}

func weixinCallback(serverApi *wxopen.WxOpen, apps wxopen.ReleaseApps) http.HandlerFunc {
	f := serveData(serverApi, apps)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "post" {
			if err := serverApi.ServeData(w, r, f); err != nil {
				fmt.Printf("serve data fail %v", err)
			}
		} else {
			utils.HttpAbortBadRequest(w)
		}
	}
}
