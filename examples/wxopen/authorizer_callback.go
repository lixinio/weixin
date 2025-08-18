package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/server_api"
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
		case *server_api.MessageVoice:
			fmt.Printf("MessageVoice : %s %s\n", v.Format, v.MediaId)
			return serverApi.ResponseVoice(w, r, &server_api.ReplyMessageVoice{
				ReplyMessage: *v.Reply(),
				Voice: struct {
					MediaId server_api.CDATA
				}{
					MediaId: server_api.CDATA(v.MediaId),
				},
			})
		case *server_api.MessageVideo:
			fmt.Printf("MessageVideo : %s %s\n", v.MediaId, v.ThumbMediaId)
			// return  serverApi.ResponseVideo(w, r, &server_api.ReplyMessageVideo{
			// 	ReplyMessage: *v.Reply(),
			// 	Video: struct {
			// 		MediaId     server_api.CDATA
			// 		Title       server_api.CDATA
			// 		Description server_api.CDATA
			// 	}{
			// 		// 通过素材管理中的接口上传多媒体文件，得到的id
			// 		// 直接用似乎不行
			// 		MediaId:     server_api.CDATA(v.MediaId),
			// 		Title:       "Title",
			// 		Description: "Description",
			// 	},
			// })
		case *server_api.MessageLocation:
			fmt.Printf("MessageLocation : %s %sX%s\n", v.Label, v.Location_X, v.Location_Y)
			news := &server_api.ReplyMessageNewsItem{
				Title:       "欢迎关注",
				Description: "hello",
				PicUrl:      "https://mat1.gtimg.com/pingjs/ext2020/qqindex2018/dist/img/qq_logo_2x.png",
				URL:         "https://www.baidu.com",
			}
			msg := &server_api.ReplyMessageNews{
				ReplyMessage: *v.Reply(),
				ArticleCount: "1",
				Articles: struct {
					Item []server_api.ReplyMessageNewsItem `xml:"item"`
				}{
					Item: []server_api.ReplyMessageNewsItem{*news},
				},
			}
			msg.ReplyMessage.MsgType = server_api.ReplyMsgTypeNews
			return serverApi.ResponseNews(w, r, msg)
		case *server_api.MessageLink:
			fmt.Printf("MessageLink : %s %s %s\n", v.Title, v.Url, v.Description)
		case *server_api.MessageFile:
			fmt.Printf("MessageFile : %s %s\n", v.Title, v.Description)
		case *server_api.EventSubscribe:
			fmt.Printf("EventSubscribe : %s %s\n", v.FromUserName, v.EventKey)
		case *server_api.EventUnsubscribe:
			fmt.Printf("EventUnsubscribe : %s\n", v.FromUserName)
		case *server_api.EventTemplateSendJobFinish:
			fmt.Printf("EventTemplateSendJobFinish : %s %s\n", v.MsgID, v.Status)
		case *server_api.EventMenuClick:
			fmt.Printf("EventMenuClick : %s\n", v.EventKey)
		case *server_api.EventAuthorizeInvoice:
			fmt.Printf("EventMenuClick : %s %s %s %s\n", v.SuccOrderId, v.FailOrderId, v.AuthorizeAppId, v.Source)
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
