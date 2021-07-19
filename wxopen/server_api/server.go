package server_api

// https://github.com/fastwego/wxopen/blob/master/server.go

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
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/official_account"
)

type ServerApi struct {
	*utils.Client
	OAConfig       *official_account.Config
	Token          string
	EncodingAESKey string
}

func NewOfficialAccountApi(token, encodingAESKey string, officialAccount *official_account.OfficialAccount) *ServerApi {
	return &ServerApi{
		Client:         officialAccount.Client,
		OAConfig:       officialAccount.Config,
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}
}

func calcSignature(timestamp, nonce, token string) string {
	strs := []string{timestamp, nonce, token}
	sort.Strings(strs)

	h := sha1.New()
	_, _ = io.WriteString(h, strings.Join(strs, ""))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func calcSignatureFromHttp(r *http.Request, token string) string {
	return calcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		token,
	)
}

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func (s *ServerApi) ServeEcho(w http.ResponseWriter, r *http.Request) {
	signature := calcSignatureFromHttp(r, s.Token)
	echoStr := r.URL.Query().Get("echostr")
	if echoStr != "" && signature == r.URL.Query().Get("signature") {
		io.WriteString(w, echoStr)
	} else {
		httpAbort(w, http.StatusBadRequest)
	}
}

func (s *ServerApi) ServeData(w http.ResponseWriter, r *http.Request, processor http.HandlerFunc) {
	signature := calcSignatureFromHttp(r, s.Token)
	if signature == r.URL.Query().Get("signature") {
		processor(w, r)
	} else {
		httpAbort(w, http.StatusBadRequest)
	}
}

// ParseXML 解析微信推送过来的消息/事件
func (s *ServerApi) ParseXML(body []byte) (m interface{}, err error) {

	// 是否加密消息
	encryptMsg := EncryptMessage{}
	err = xml.Unmarshal(body, &encryptMsg)
	if err != nil {
		return
	}

	// 需要解密
	if encryptMsg.Encrypt != "" {
		var xmlMsg []byte
		_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, s.EncodingAESKey)
		if err != nil {
			return
		}
		body = xmlMsg

	}

	message := Message{}
	err = xml.Unmarshal(body, &message)
	//fmt.Println(message)
	if err != nil {
		return
	}

	switch message.MsgType {
	case MsgTypeText:
		msg := MessageText{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeImage:
		msg := MessageImage{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeVoice:
		msg := MessageVoice{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeVideo:
		msg := MessageVideo{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeShortVideo:
		msg := MessageShortVideo{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeLocation:
		msg := MessageLocation{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeLink:
		msg := MessageLink{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case MsgTypeFile:
		msg := MessageFile{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
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
	event := Event{}
	err = xml.Unmarshal(body, &event)
	if err != nil {
		return
	}
	switch event.Event {

	// 关注事件
	case EventTypeSubscribe:
		msg := EventSubscribe{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUnsubscribe:
		msg := EventUnsubscribe{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeScan:
		msg := EventScan{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeLocation:
		msg := EventLocation{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil

	// 菜单事件
	case EventTypeMenuClick:
		msg := EventMenuClick{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuView:
		msg := EventMenuView{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuScanCodePush:
		msg := EventMenuScanCodePush{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuScanCodeWaitMsg:
		msg := EventMenuScanCodeWaitMsg{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuPicSysPhoto:
		msg := EventMenuPicSysPhoto{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuPicSysPhotoOrAlbum:
		msg := EventMenuPicSysPhotoOrAlbum{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuPicWeixin:
		msg := EventMenuPicWeixin{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuLocationSelect:
		msg := EventMenuLocationSelect{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeMenuViewMiniprogram:
		msg := EventMenuViewMiniprogram{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil

	// 资质事件
	case EventTypeQualificationVerifySuccess:
		msg := EventQualificationVerifySuccess{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeQualificationVerifyFail:
		msg := EventQualificationVerifyFail{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeNamingVerifySuccess:
		msg := EventNamingVerifySuccess{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeNamingVerifyFail:
		msg := EventNamingVerifyFail{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeAnnualRenew:
		msg := EventAnnualRenew{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeVerifyExpired:
		msg := EventVerifyExpired{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil

		// 卡券事件
	case EventTypeCardPassChecke:
		msg := EventCardPassChecke{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeCardNotPassChecke:
		msg := EventCardNotPassChecke{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserGetCard:
		msg := EventUserGetCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserGiftingCard:
		msg := EventUserGiftingCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserDelCard:
		msg := EventUserDelCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserConsumeCard:
		msg := EventUserConsumeCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserPayFromPayCell:
		msg := EventUserPayFromPayCell{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserViewCard:
		msg := EventUserViewCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUserEnterSessionFromCard:
		msg := EventUserEnterSessionFromCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateMemberCard:
		msg := EventUpdateMemberCard{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeCardSkuRemind:
		msg := EventCardSkuRemind{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeCardPayOrder:
		msg := EventCardPayOrder{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeSubmitMembercardUserInfo:
		msg := EventSubmitMembercardUserInfo{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil

		// 导购事件
	case EventTypeGuideQrcodeScan:
		msg := EventGuideQrcodeScan{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil

		// 模版消息发送任务完成
	case EventTypeTemplateSendJobFinish:
		msg := EventTemplateSendJobFinish{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	}

	return
}

// Response 响应微信消息 (自动判断是否要加密)
func (s *ServerApi) response(w http.ResponseWriter, r *http.Request, reply interface{}) (err error) {

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
		fmt.Println(string(output))

		// 加密
		if r.URL.Query().Get("encrypt_type") == "aes" {
			var message *ReplyEncryptMessage
			message, err = s.encryptReplyMessage(output)
			if err != nil {
				return
			}
			output, err = xml.Marshal(message)
			if err != nil {
				return
			}
		}
	}

	_, err = w.Write(output)

	return
}

// encryptReplyMessage 加密回复消息
func (s *ServerApi) encryptReplyMessage(rawXmlMsg []byte) (*ReplyEncryptMessage, error) {
	cipherText, err := utils.AESEncryptMsg([]byte(utils.GetRandString(16)), rawXmlMsg, s.OAConfig.Appid, s.EncodingAESKey)
	if err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := utils.GetRandString(6)

	strs := []string{
		timestamp,
		nonce,
		s.Token,
		cipherText,
	}
	sort.Strings(strs)
	h := sha1.New()
	_, _ = io.WriteString(h, strings.Join(strs, ""))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	return &ReplyEncryptMessage{
		Encrypt:      cipherText,
		MsgSignature: signature,
		TimeStamp:    timestamp,
		Nonce:        nonce,
	}, nil
}

func (s *ServerApi) ResponseText(w http.ResponseWriter, r *http.Request, message *ReplyMessageText) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseImage(w http.ResponseWriter, r *http.Request, message *ReplyMessageImage) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseVoice(w http.ResponseWriter, r *http.Request, message *ReplyMessageVoice) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseVideo(w http.ResponseWriter, r *http.Request, message *ReplyMessageVideo) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseMusic(w http.ResponseWriter, r *http.Request, message *ReplyMessageMusic) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseNews(w http.ResponseWriter, r *http.Request, message *ReplyMessageNews) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseTransferCustomerService(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageTransferCustomerService,
) (err error) {
	return s.response(w, r, message)
}
