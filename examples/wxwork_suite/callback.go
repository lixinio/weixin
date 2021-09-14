package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/wxwork_suite"
)

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func serveData(serverApi *wxwork_suite.WxWorkSuite) http.HandlerFunc {
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
		case *wxwork_suite.EventSuiteTicket:
			fmt.Printf("Suite Ticket : %s\n", v.SuiteTicket)
			// todo save to db
			// 刷新Ticket到Cache
			serverApi.UpdateTicket(v.SuiteTicket)
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

		io.WriteString(w, "success")
	}
}

func weixinCallback(serverApi *wxwork_suite.WxWorkSuite) http.HandlerFunc {
	f := serveData(serverApi)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "get" {
			serverApi.ServeEcho(w, r)
		} else if strings.ToLower(r.Method) == "post" {
			f(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, http.StatusText(http.StatusBadRequest))
		}
	}
}
