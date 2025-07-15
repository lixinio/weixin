package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/server_api"
)

func serveAuthorizerData(serverApi *server_api.ServerApi) utils.XmlHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, body []byte) error {
		_, content, err := serverApi.ParseXML(body)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return err
		}

		switch v := content.(type) {
		case *server_api.MessageText:
			fmt.Printf("MsgTypeText : %s\n", v.Content)
			return serverApi.ResponseText(w, r, &server_api.ReplyMessageText{
				ReplyMessage: *v.Reply(),
				Content:      server_api.CDATA(v.Content),
			})
		case *server_api.MessageImage:
			fmt.Printf("MessageImage : %s %s\n", v.MediaId, v.PicUrl)
			return serverApi.ResponseImage(w, r, &server_api.ReplyMessageImage{
				ReplyMessage: *v.Reply(),
				Image: struct {
					MediaId server_api.CDATA
				}{
					MediaId: server_api.CDATA(v.MediaId),
				},
			})
		case *server_api.EventApproval:
			fmt.Printf(
				"审批变更 : %s %s %s\n",
				v.ApprovalInfo.ThirdNo,
				v.ApprovalInfo.OpenSpName,
				v.ApprovalInfo.OpenSpStatus,
			)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		_, err = io.WriteString(w, "success")
		return err
	}
}

func authorizerCallback(serverApi *server_api.ServerApi) http.HandlerFunc {
	f := serveAuthorizerData(serverApi)
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
