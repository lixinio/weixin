package wxopen

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/message_api"
	"github.com/lixinio/weixin/weixin/server_api"
)

type ReleaseApp struct {
	UserName string // 例如 gh_8dad206e9538
	IsMp     bool   // true: 小程序，false，公众号
}

type ReleaseApps map[string]*ReleaseApp

var ReleaseAppIDS = ReleaseApps{
	"wx570bc396a51b8ff8": {
		UserName: "gh_3c884a361561",
		IsMp:     false,
	},
	"wx9252c5e0bb1836fc": {
		UserName: "gh_c0f28a78b318",
		IsMp:     false,
	},
	"wx8e1097c5bc82cde9": {
		UserName: "gh_3f222ed8d140",
		IsMp:     false,
	},
	"wx14550af28c71a144": {
		UserName: "gh_26128078e9ab",
		IsMp:     false,
	},
	"wxa35b9c23cfe664eb": {
		UserName: "gh_2b3713f184a6",
		IsMp:     false,
	},
	"wxd101a85aa106f53e": {
		UserName: "gh_8dad206e9538",
		IsMp:     true,
	},
	"wxc39235c15087f6f3": {
		UserName: "gh_905ae9d01059",
		IsMp:     true,
	},
	"wx7720d01d4b2a4500": {
		UserName: "gh_393666f1fdf4",
		IsMp:     true,
	},
	"wx05d483572dcd5d8b": {
		UserName: "gh_39abb5d4e1b7",
		IsMp:     true,
	},
	"wx5910277cae6fd970": {
		UserName: "gh_7818dcb60240",
		IsMp:     true,
	},
}

// 全网发布
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/operation/thirdparty/releases_instructions.html
func (api *WxOpen) ServeRelease(
	serverApi *server_api.ServerApi,
) utils.XmlHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, body []byte) (err error) {
		message := &server_api.Message{}
		if err = xml.Unmarshal(body, message); err != nil {
			utils.HttpAbortBadRequest(w)
			return
		}

		switch message.MsgType {
		case server_api.MsgTypeText:
			msg := &server_api.MessageText{}
			if err = xml.Unmarshal(body, msg); err != nil {
				utils.HttpAbortBadRequest(w)
				return
			}
			if msg.Content == "TESTCOMPONENT_MSG_TYPE_TEXT" {
				serverApi.ResponseText(w, r, &server_api.ReplyMessageText{
					ReplyMessage: *msg.Reply(),
					Content:      server_api.CDATA("TESTCOMPONENT_MSG_TYPE_TEXT_callback"),
				})
				return nil
			}

			if err = processOtherText(r.Context(), api, msg); err != nil {
				return
			}
			w.WriteHeader(http.StatusOK)
		case server_api.MsgTypeEvent:
			event := &server_api.Event{}
			if err = xml.Unmarshal(body, event); err != nil {
				utils.HttpAbortBadRequest(w)
				return
			}
			return serverApi.ResponseText(w, r, &server_api.ReplyMessageText{
				ReplyMessage: *event.Reply(),
				Content:      server_api.CDATA(fmt.Sprintf("%sfrom_callback", event.Event)),
			})
		}
		return nil
	}
}

func newStaticClient(accessToken string) *utils.Client {
	client := utils.NewClient(
		WXServerUrl, utils.StaticClientAccessTokenGetter(accessToken),
	)
	client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	return client
}

func processOtherText(
	ctx context.Context, api *WxOpen, msg *server_api.MessageText,
) (err error) {
	items := strings.Split(msg.Content, ":")
	if len(items) != 2 || items[0] != "QUERY_AUTH_CODE" {
		return fmt.Errorf("invalid query auth code %s", items)
	}

	authCode := items[1]
	authInfo, err := api.QueryAuth(ctx, authCode)
	if err != nil {
		return err
	}

	client := newStaticClient(authInfo.AuthorizerAccessToken)
	messageApi := message_api.NewApi(client)
	err = messageApi.SendCustomTextMessage(
		ctx,
		msg.FromUserName,
		fmt.Sprintf("%s_from_api", authCode),
	)
	if err != nil {
		return fmt.Errorf("error send custom text message by token '%s', error %w",
			authInfo.AuthorizerAccessToken, err,
		)
	}
	return
}
