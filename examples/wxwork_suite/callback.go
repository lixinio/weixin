package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork_suite"
)

func serveData(serverApi *wxwork_suite.WxWorkSuite) utils.XmlHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, body []byte) error {
		_, content, err := serverApi.ParseXML(body)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return err
		}

		switch v := content.(type) {
		case *wxwork_suite.EventSuiteTicket:
			fmt.Printf("Suite Ticket : %s\n", v.SuiteTicket)
			// todo save to db
			// 刷新Ticket到Cache
			serverApi.UpdateTicket(r.Context(), v.SuiteTicket)
		case *wxwork_suite.EventAuthorized:
			fmt.Printf("EventAuthorized : %s\n", v.SuiteId)
		case *wxwork_suite.EventUnauthorized:
			fmt.Printf("EventUnauthorized : %s\n", v.SuiteId)
		case *wxwork_suite.EventUpdateAuthorized:
			fmt.Printf("EventUpdateAuthorized : %s\n", v.SuiteId)
		case *wxwork_suite.EventCreateUser:
			fmt.Printf("EventCreateUser : %s %s %s\n", v.SuiteId, v.UserID, v.Alias)
		case *wxwork_suite.EventUpdateUser:
			fmt.Printf("EventUpdateUser : %s %s %s\n", v.SuiteId, v.UserID, v.Alias)
		case *wxwork_suite.EventDeleteUser:
			fmt.Printf("EventDeleteUser : %s %s\n", v.SuiteId, v.UserID)
		case *wxwork_suite.EventCreateParty:
			fmt.Printf("EventCreateParty : %s %d %d\n", v.SuiteId, v.Id, v.ParentId)
		case *wxwork_suite.EventUpdateParty:
			fmt.Printf("EventUpdateParty : %s %d %d\n", v.SuiteId, v.Id, v.ParentId)
		case *wxwork_suite.EventDeleteParty:
			fmt.Printf("EventDeleteParty : %s %d\n", v.SuiteId, v.Id)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		_, err = io.WriteString(w, "success")
		return err
	}
}

func weixinCallback(serverApi *wxwork_suite.WxWorkSuite) http.HandlerFunc {
	f := serveData(serverApi)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "get" {
			if err := serverApi.ServeEcho(w, r); err != nil {
				fmt.Printf("serve echo fail %v", err)
			}
		} else if strings.ToLower(r.Method) == "post" {
			if err := serverApi.ServeData(w, r, f); err != nil {
				fmt.Printf("serve data fail %v", err)
			}
		} else {
			utils.HttpAbortBadRequest(w)
		}
	}
}
