package server_api

// https://github.com/fastwego/offiaccount/blob/master/server.go

// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/lixinio/weixin/utils"
)

type ServerApi struct {
	*utils.Client
	AppID          string
	Token          string
	EncodingAESKey string
}

func NewApi(
	appid, token, encodingAESKey string,
	client *utils.Client,
) *ServerApi {
	return &ServerApi{
		Client:         client,
		AppID:          appid,
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}
}

func calcSignatureFromHttp(r *http.Request, token string) string {
	return utils.CalcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		token,
	)
}

func (s *ServerApi) ServeEcho(w http.ResponseWriter, r *http.Request) error {
	signature := calcSignatureFromHttp(r, s.Token)
	echoStr := r.URL.Query().Get("echostr")
	if echoStr != "" && signature == r.URL.Query().Get("signature") {
		_, err := io.WriteString(w, echoStr)
		return err
	} else {
		utils.HttpAbortBadRequest(w)
		return errors.New("signature dismatch")
	}
}

func (s *ServerApi) ServeData(
	w http.ResponseWriter,
	r *http.Request,
	processor utils.XmlHandlerFunc,
) error {
	signature := calcSignatureFromHttp(r, s.Token)
	if signature != r.URL.Query().Get("signature") {
		utils.HttpAbortBadRequest(w)
		return fmt.Errorf(
			"signature dismatch %s != %s",
			signature, r.URL.Query().Get("signature"),
		)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.HttpAbortBadRequest(w)
		return err
	}

	// 是否加密消息
	encryptMsg := &EncryptMessage{}
	if err = xml.Unmarshal(body, encryptMsg); err != nil {
		return err
	}

	// 需要解密
	if encryptMsg.Encrypt != "" {
		msgSignature := r.URL.Query().Get("msg_signature")
		timestamp := r.URL.Query().Get("timestamp")
		nonce := r.URL.Query().Get("nonce")

		if utils.CalcSignature(
			s.Token, timestamp, nonce, encryptMsg.Encrypt,
		) != msgSignature {
			return errors.New("invalid msg signature")
		}

		var xmlMsg []byte
		_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, s.EncodingAESKey)
		if err != nil {
			return err
		}
		body = xmlMsg
	}

	return processor(w, r, body)
}

// ParseXML 解析微信推送过来的消息/事件
func (s *ServerApi) ParseXML(body []byte) (m interface{}, err error) {
	message := &Message{}
	if err = xml.Unmarshal(body, message); err != nil {
		return
	}

	switch message.MsgType {
	case MsgTypeText:
		msg := &MessageText{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeImage:
		msg := &MessageImage{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeVoice:
		msg := &MessageVoice{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeVideo:
		msg := &MessageVideo{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeShortVideo:
		msg := &MessageShortVideo{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeLocation:
		msg := &MessageLocation{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeLink:
		msg := &MessageLink{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeFile:
		msg := &MessageFile{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case MsgTypeEvent:
		return parseEvent(body)
	}
	return
}

// parseEvent 解析微信推送过来的事件
func parseEvent(body []byte) (m interface{}, err error) {
	event := &Event{}
	if err = xml.Unmarshal(body, event); err != nil {
		return
	}
	switch event.Event {

	// 关注事件
	case EventTypeSubscribe:
		msg := &EventSubscribe{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUnsubscribe:
		msg := &EventUnsubscribe{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeScan:
		msg := &EventScan{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeLocation:
		msg := &EventLocation{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil

	// 菜单事件
	case EventTypeMenuClick:
		msg := &EventMenuClick{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuView:
		msg := &EventMenuView{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuScanCodePush:
		msg := &EventMenuScanCodePush{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuScanCodeWaitMsg:
		msg := &EventMenuScanCodeWaitMsg{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuPicSysPhoto:
		msg := &EventMenuPicSysPhoto{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuPicSysPhotoOrAlbum:
		msg := &EventMenuPicSysPhotoOrAlbum{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuPicWeixin:
		msg := &EventMenuPicWeixin{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuLocationSelect:
		msg := &EventMenuLocationSelect{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuViewMiniprogram:
		msg := &EventMenuViewMiniprogram{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil

	// 资质事件
	case EventTypeQualificationVerifySuccess:
		msg := &EventQualificationVerifySuccess{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeQualificationVerifyFail:
		msg := &EventQualificationVerifyFail{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeNamingVerifySuccess:
		msg := &EventNamingVerifySuccess{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeNamingVerifyFail:
		msg := &EventNamingVerifyFail{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeAnnualRenew:
		msg := &EventAnnualRenew{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeVerifyExpired:
		msg := &EventVerifyExpired{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil

		// 卡券事件
	case EventTypeCardPassChecke:
		msg := &EventCardPassChecke{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeCardNotPassChecke:
		msg := &EventCardNotPassChecke{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserGetCard:
		msg := &EventUserGetCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserGiftingCard:
		msg := &EventUserGiftingCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserDelCard:
		msg := &EventUserDelCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserConsumeCard:
		msg := &EventUserConsumeCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserPayFromPayCell:
		msg := &EventUserPayFromPayCell{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserViewCard:
		msg := &EventUserViewCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUserEnterSessionFromCard:
		msg := &EventUserEnterSessionFromCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateMemberCard:
		msg := &EventUpdateMemberCard{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeCardSkuRemind:
		msg := &EventCardSkuRemind{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeCardPayOrder:
		msg := &EventCardPayOrder{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeSubmitMembercardUserInfo:
		msg := &EventSubmitMembercardUserInfo{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil

		// 导购事件
	case EventTypeGuideQrcodeScan:
		msg := &EventGuideQrcodeScan{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil

		// 模版消息发送任务完成
	case EventTypeTemplateSendJobFinish:
		msg := &EventTemplateSendJobFinish{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil

	case EventTypeAuthorizeInvoice:
		msg := &EventAuthorizeInvoice{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	}

	return
}

// Response 响应微信消息 (自动判断是否要加密)
func (s *ServerApi) response(
	w http.ResponseWriter,
	r *http.Request,
	reply interface{},
) (err error) {

	// 如果 开启加密，微信服务器 发过来的请求 带有 如下参数
	//signature=c44d29564aa1d57bd0e274c37baa92bd5b3da5bd
	//&timestamp=1596184957
	//&nonce=1250398014
	//&openid=oEnxesxpxWw-PKkz-vW5IMdfcQaE
	//&encrypt_type=aes
	//&msg_signature=cc24cc38467417603fc3689170e8b0fd3c9bf4a2

	output := []byte("success") // 默认回复
	if reply != nil {
		output, err = xml.Marshal(reply)
		if err != nil {
			return
		}

		// 加密
		if r.URL.Query().Get("encrypt_type") == "aes" {
			var message *ReplyEncryptMessage
			message, err = s.encryptReplyMessage(output)
			if err != nil {
				fmt.Println("encryptReplyMessage", err)
				return
			}
			output, err = xml.Marshal(message)
			if err != nil {
				fmt.Println("marshal encryptReplyMessage", err)
				return
			}
		}
	}

	_, err = w.Write(output)

	return
}

// encryptReplyMessage 加密回复消息
func (s *ServerApi) encryptReplyMessage(rawXmlMsg []byte) (*ReplyEncryptMessage, error) {
	cipherText, err := utils.AESEncryptMsg(
		[]byte(utils.GetRandString(16)),
		rawXmlMsg,
		s.AppID,
		s.EncodingAESKey,
	)
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := utils.GetRandString(6)
	signature := utils.CalcSignature(timestamp, nonce, s.Token, cipherText)

	return &ReplyEncryptMessage{
		Encrypt:      CDATA(cipherText),
		MsgSignature: CDATA(signature),
		TimeStamp:    timestamp,
		Nonce:        CDATA(nonce),
	}, nil
}

func (s *ServerApi) ResponseText(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageText,
) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseImage(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageImage,
) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseVoice(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageVoice,
) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseVideo(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageVideo,
) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseMusic(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageMusic,
) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseNews(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageNews,
) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseTransferCustomerService(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageTransferCustomerService,
) (err error) {
	return s.response(w, r, message)
}
