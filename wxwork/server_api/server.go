package server_api

// https://github.com/fastwego/wxwork/blob/master/corporation/server.go

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
	"net/http"
	"strconv"
	"time"

	"github.com/lixinio/weixin/utils"
)

type ServerApi struct {
	AgentID        string
	Token          string // 接收消息服务器配置（Token）
	EncodingAESKey string // 接收消息服务器配置（EncodingAESKey）
}

func NewApi(
	agentID int,
	token, encodingAESKey string,
) *ServerApi {
	return &ServerApi{
		AgentID:        strconv.Itoa(agentID),
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}
}

func calcSignatureFromHttp(r *http.Request, token string) (string, string) {
	echostr := r.URL.Query().Get("echostr")
	return utils.CalcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		echostr,
		token,
	), echostr
}

func (s *ServerApi) ServeEcho(w http.ResponseWriter, r *http.Request) error {
	signature, echoStr := calcSignatureFromHttp(r, s.Token)
	if echoStr != "" && signature == r.URL.Query().Get("msg_signature") {
		// 解密 echoStr
		_, msg, _, err := utils.AESDecryptMsg(echoStr, s.EncodingAESKey)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return err
		}
		_, err = w.Write(msg)
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// 加密格式 的消息
	encryptMsg := &EncryptMessage{}
	if err = xml.Unmarshal(body, encryptMsg); err != nil {
		return err
	}

	// 验证签名
	signature := utils.CalcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		encryptMsg.Encrypt,
		s.Token,
	)

	if msgSignature := r.URL.Query().Get("msg_signature"); signature != msgSignature {
		err = fmt.Errorf("signature dismatch %s != %s", signature, msgSignature)
		return err
	}

	// 解密
	var xmlMsg []byte
	_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, s.EncodingAESKey)
	if err != nil {
		return err
	}
	return processor(w, r, xmlMsg)
}

/*
ParseXML 解析微信推送过来的消息/事件

POST /api/callback?msg_signature=ASDFQWEXZCVAQFASDFASDFSS
&timestamp=13500001234
&nonce=123412323

<xml>

	<ToUserName><![CDATA[toUser]]></ToUserName>
	<AgentID><![CDATA[toAgentID]]></AgentID>
	<Encrypt><![CDATA[msg_encrypt]]></Encrypt>

</xml>
*/
func (s *ServerApi) ParseXML(body []byte) (message *Message, m interface{}, err error) {
	message = &Message{}
	if err = xml.Unmarshal(body, message); err != nil {
		return
	}

	switch message.MsgType {
	case MsgTypeText:
		msg := &MessageText{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return message, msg, nil
	case MsgTypeImage:
		msg := &MessageImage{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return message, msg, nil
	case MsgTypeVoice:
		msg := &MessageVoice{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return message, msg, nil
	case MsgTypeVideo:
		msg := &MessageVideo{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return message, msg, nil
	case MsgTypeLocation:
		msg := &MessageLocation{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return message, msg, nil
	case MsgTypeLink:
		msg := &MessageLink{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return message, msg, nil
	case MsgTypeEvent:
		m, err = parseEvent(body)
	}

	return
}

// ParseEvent 解析微信推送过来的事件
func parseEvent(body []byte) (m interface{}, err error) {
	event := &Event{}
	if err = xml.Unmarshal(body, event); err != nil {
		return
	}
	switch event.Event {
	// 事件
	case EventTypeChangeContact:
		msg := &EventChangeContact{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		switch msg.ChangeType {
		case EventTypeChangeContactCreateUser:
			msg := &EventChangeContactCreateUser{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactUpdateUser:
			msg := &EventChangeContactUpdateUser{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactDeleteUser:
			msg := &EventChangeContactDeleteUser{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactCreateParty:
			msg := &EventChangeContactCreateParty{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactUpdateParty:
			msg := &EventChangeContactUpdateParty{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactDeleteParty:
			msg := &EventChangeContactDeleteParty{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactUpdateTag:
			msg := &EventChangeContactUpdateTag{}
			if err = xml.Unmarshal(body, msg); err != nil {
				return
			}
			return msg, nil
		}
	case EventTypeBatchJobResult:
		msg := &EventBatchJobResult{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeApproval:
		msg := &EventApproval{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeTaskCardClick:
		msg := &EventTaskCardClick{}
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
	case EventTypeMenuClick:
		msg := &EventMenuClick{}
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
	case EventTypeEnterAgent:
		msg := &EventEnterAgent{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	}

	return
}

// Response 响应微信消息
func (s *ServerApi) response(
	w http.ResponseWriter,
	r *http.Request,
	reply interface{},
) (err error) {
	output := []byte("") // 默认回复
	if reply != nil {
		output, err = xml.Marshal(reply)
		if err != nil {
			return
		}

		// 加密
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

	_, err = w.Write(output)

	return
}

// encryptReplyMessage 加密回复消息
func (s *ServerApi) encryptReplyMessage(
	rawXmlMsg []byte,
) (replyEncryptMessage *ReplyEncryptMessage, err error) {
	cipherText, err := utils.AESEncryptMsg(
		[]byte(utils.GetRandString(16)),
		rawXmlMsg,
		s.AgentID,
		s.EncodingAESKey,
	)
	if err != nil {
		return
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

func (s *ServerApi) ResponseTaskCard(
	w http.ResponseWriter,
	r *http.Request,
	message *ReplyMessageTaskCard,
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
