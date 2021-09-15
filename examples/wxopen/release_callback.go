package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/server_api"
	"github.com/lixinio/weixin/wxopen"
)

func releaseCallback(wxopenApi *wxopen.WxOpen, serverApi *server_api.ServerApi) http.HandlerFunc {
	f := wxopenApi.ServeRelease(wxopenApi, serverApi)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "get" {
			if err := serverApi.ServeEcho(w, r); err != nil {
				fmt.Printf("serve echo fail %v", err)
			}
		} else if strings.ToLower(r.Method) == "post" {
			if err := serverApi.ServeData(w, r, f); err != nil {
				fmt.Printf("serve release fail %v", err)
			}
		} else {
			utils.HttpAbortBadRequest(w)
		}
	}
}
