package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/lixinio/weixin/weixin/server_api"
)

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func serveData(serverApi *server_api.ServerApi) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			httpAbort(w, http.StatusBadRequest)
			return
		}
		content, err := serverApi.ParseXML(body)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}

		switch v := content.(type) {
		case server_api.MessageText:
			fmt.Printf("MsgTypeText : %s\n", v.Content)
			serverApi.ResponseText(w, r, &server_api.ReplyMessageText{
				ReplyMessage: *v.Reply(),
				Content:      server_api.CDATA(v.Content),
			})
			return
		case server_api.MessageImage:
			fmt.Printf("MessageImage : %s %s\n", v.MediaId, v.PicUrl)
			serverApi.ResponseImage(w, r, &server_api.ReplyMessageImage{
				ReplyMessage: *v.Reply(),
				Image: struct {
					MediaId server_api.CDATA
				}{
					MediaId: server_api.CDATA(v.MediaId),
				},
			})
		case server_api.MessageVoice:
			fmt.Printf("MessageVoice : %s %s\n", v.Format, v.MediaId)
			serverApi.ResponseVoice(w, r, &server_api.ReplyMessageVoice{
				ReplyMessage: *v.Reply(),
				Voice: struct {
					MediaId server_api.CDATA
				}{
					MediaId: server_api.CDATA(v.MediaId),
				},
			})
			return
		case server_api.MessageVideo:
			fmt.Printf("MessageVideo : %s %s\n", v.MediaId, v.ThumbMediaId)
			// serverApi.ResponseVideo(w, r, &server_api.ReplyMessageVideo{
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
			// return
		case server_api.MessageLocation:
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
			serverApi.ResponseNews(w, r, msg)
			return
		case server_api.MessageLink:
			fmt.Printf("MessageLink : %s %s %s\n", v.Title, v.Url, v.Description)
		case server_api.MessageFile:
			fmt.Printf("MessageFile : %s %s\n", v.Title, v.Description)
		case server_api.EventSubscribe:
			fmt.Printf("EventSubscribe : %s %s\n", v.FromUserName, v.EventKey)
			return
		case server_api.EventUnsubscribe:
			fmt.Printf("EventUnsubscribe : %s\n", v.FromUserName)
		case server_api.EventTemplateSendJobFinish:
			fmt.Printf("EventTemplateSendJobFinish : %s %s\n", v.MsgID, v.Status)
		case server_api.EventMenuClick:
			fmt.Printf("EventMenuClick : %s\n", v.EventKey)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

		io.WriteString(w, "success")
		return
	}
}

func weixinCallback(serverApi *server_api.ServerApi) http.HandlerFunc {
	f := serveData(serverApi)
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "get" {
			serverApi.ServeEcho(w, r)
		} else if strings.ToLower(r.Method) == "post" {
			serverApi.ServeData(w, r, f)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, http.StatusText(http.StatusBadRequest))
		}
	}
}

func main() {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})
	serverApi := server_api.NewOfficialAccountApi(
		test.OfficialAccountToken,
		test.OfficialAccountAESKey,
		officialAccount,
	)

	http.HandleFunc(fmt.Sprintf("/weixin/%s", test.OfficialAccountAppid), weixinCallback(serverApi))

	http.ListenAndServe(":5000", nil)
}
