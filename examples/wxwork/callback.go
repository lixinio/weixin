package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/server_api"
)

func serveData(serverApi *server_api.ServerApi) utils.XmlHandlerFunc {
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
			return serverApi.ResponseVideo(w, r, &server_api.ReplyMessageVideo{
				ReplyMessage: *v.Reply(),
				Video: struct {
					MediaId     server_api.CDATA
					Title       server_api.CDATA
					Description server_api.CDATA
				}{
					MediaId:     server_api.CDATA(v.MediaId),
					Title:       "Title",
					Description: "Description",
				},
			})
		case *server_api.MessageLocation:
			fmt.Printf("MessageLocation : %s %s %s %sX%s\n", v.Label, v.Scale, v.AppType, v.Location_X, v.Location_Y)
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
			fmt.Printf("MessageLink : %s %s %s %s\n", v.Title, v.Url, v.PicUrl, v.Description)
			msg := &server_api.ReplyMessageTaskCard{
				ReplyMessage: *v.Reply(),

				TaskCard: struct {
					ReplaceName server_api.CDATA
				}{
					ReplaceName: "task_card",
				},
			}
			msg.ReplyMessage.MsgType = server_api.ReplyMsgTypeTaskCard
			return serverApi.ResponseTaskCard(w, r, msg)
		case *server_api.EventChangeContactCreateUser:
			fmt.Print("create user", " ", v.ChangeType, " ",
				v.UserID, " ", v.Name, " ", v.Mobile, " ", v.Alias, " ", v.Email, " ",
				v.Position, " ", v.Telephone, " ", v.Address, " ",
				v.Avatar, " ", v.Gender, " ", v.Status, " ", v.Department, " ",
				v.MainDepartment, " ", v.IsLeaderInDept, "\n",
			)
		case *server_api.EventChangeContactUpdateUser:
			fmt.Print("update user", " ", v.ChangeType, " ",
				v.UserID, " ", v.Name, " ", v.Mobile, " ", v.Alias, " ", v.Email, " ",
				v.Position, " ", v.Telephone, " ", v.Address, " ",
				v.Avatar, " ", v.Gender, " ", v.Status, " ", v.Department, " ",
				v.MainDepartment, " ", v.IsLeaderInDept, " ",
				v.NewUserID, "\n",
			)
		case *server_api.EventChangeContactDeleteUser:
			fmt.Print("delete user", " ", v.ChangeType, " ", v.UserID, "\n")
		case *server_api.EventChangeContactCreateParty:
			fmt.Print("create party", " ", v.ChangeType, " ", v.ID, " ", v.Name, " ", v.Order, " ", v.ParentId, "\n")
		case *server_api.EventChangeContactUpdateParty:
			fmt.Print("create party", " ", v.ChangeType, " ", v.ID, " ", v.Name, " ", v.ParentId, "\n")
		case *server_api.EventChangeContactDeleteParty:
			fmt.Print("delete user", " ", v.ChangeType, " ", v.ID, "\n")
		case *server_api.EventMenuClick:
			fmt.Print("EventMenuClick", " ", v.AgentID, " ", v.EventKey, "\n")
			msg := &server_api.ReplyMessageText{
				ReplyMessage: *v.Reply(),
				Content:      server_api.CDATA(v.EventKey),
			}
			msg.MsgType = server_api.ReplyMsgTypeText
			return serverApi.ResponseText(w, r, msg)
		case *server_api.EventMenuView:
			fmt.Print("EventMenuView", " ", v.AgentID, " ", v.EventKey, "\n")
		case *server_api.EventMenuScanCodePush:
			fmt.Print("EventMenuScanCodePush", " ",
				v.AgentID, " ", v.EventKey,
				v.ScanCodeInfo.ScanType, " ", v.ScanCodeInfo.ScanResult, "\n",
			)
		case *server_api.EventMenuScanCodeWaitMsg:
			fmt.Print("EventMenuScanCodeWaitMsg", " ",
				v.AgentID, " ", v.EventKey,
				v.ScanCodeInfo.ScanType, " ", v.ScanCodeInfo.ScanResult, "\n",
			)
			msg := &server_api.ReplyMessageText{
				ReplyMessage: *v.Reply(),
				Content:      server_api.CDATA(v.ScanCodeInfo.ScanResult),
			}
			msg.MsgType = server_api.ReplyMsgTypeText
			return serverApi.ResponseText(w, r, msg)
		case *server_api.EventMenuPicSysPhoto:
			fmt.Print("EventMenuPicSysPhoto", " ", v.AgentID, " ", v.EventKey, "\n")
		case *server_api.EventMenuPicSysPhotoOrAlbum:
			fmt.Print("EventMenuPicSysPhotoOrAlbum", " ", v.AgentID, " ", v.EventKey, "\n")
		case *server_api.EventMenuPicWeixin:
			fmt.Print("EventMenuPicWeixin", " ", v.AgentID, " ", v.EventKey, "\n")
		case *server_api.EventMenuLocationSelect:
			fmt.Print("EventMenuLocationSelect", " ", v.AgentID, " ", v.EventKey,
				v.SendLocationInfo.Label, " ", v.SendLocationInfo.Text, " ", "\n")
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

func msgCallback(serverApi *server_api.ServerApi) http.HandlerFunc {
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
