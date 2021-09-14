package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/wxopen"
)

func serveData(serverApi *wxopen.WxOpen) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v\n", err)
			httpAbort(w, http.StatusBadRequest)
			return
		}
		content, err := serverApi.ParseXML(
			body,
			r.URL.Query().Get("msg_signature"),
			r.URL.Query().Get("timestamp"),
			r.URL.Query().Get("nonce"),
		)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}

		switch v := content.(type) {
		case *wxopen.EventComponentVerifyTicket:
			fmt.Printf("ComponentVerifyTicket : %s\n", v.ComponentVerifyTicket)
			// todo save to db
			// 刷新Ticket到Cache
			serverApi.UpdateTicket(v.ComponentVerifyTicket)
		case *wxopen.EventAuthorized:
			fmt.Printf("EventAuthorized : %s\n", v.AuthorizerAppid)
		case *wxopen.EventUnauthorized:
			fmt.Printf("EventUnauthorized : %s\n", v.AuthorizerAppid)
		case *wxopen.EventUpdateAuthorized:
			fmt.Printf("EventUpdateAuthorized : %s\n", v.AuthorizerAppid)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		io.WriteString(w, "success")
	}
}

func weixinCallback(serverApi *wxopen.WxOpen) http.HandlerFunc {
	f := serveData(serverApi)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "post" {
			f(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, http.StatusText(http.StatusBadRequest))
		}
	}
}
