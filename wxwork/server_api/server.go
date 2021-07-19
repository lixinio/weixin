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
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/agent"
)

type ServerApi struct {
	*utils.Client
	AgentConfig    *agent.Config
	Token          string // 接收消息服务器配置（Token）
	EncodingAESKey string // 接收消息服务器配置（EncodingAESKey）
}

func NewAgentApi(token, encodingAESKey string, agent *agent.Agent) *ServerApi {
	return &ServerApi{
		Client:         agent.Client,
		AgentConfig:    agent.Config,
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}
}

func calcSignature(timestamp, nonce, echostr, token string) string {
	strs := []string{timestamp, nonce, token, echostr}
	sort.Strings(strs)

	h := sha1.New()
	_, _ = io.WriteString(h, strings.Join(strs, ""))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func calcSignatureFromHttp(r *http.Request, token string) (string, string) {
	echostr := r.URL.Query().Get("echostr")
	return calcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		echostr,
		token,
	), echostr
}

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func (s *ServerApi) ServeEcho(w http.ResponseWriter, r *http.Request) {
	signature, echoStr := calcSignatureFromHttp(r, s.Token)
	if echoStr != "" && signature == r.URL.Query().Get("msg_signature") {
		// 解密 echoStr
		_, msg, _, err := utils.AESDecryptMsg(echoStr, s.EncodingAESKey)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}
		w.Write(msg)
	} else {
		httpAbort(w, http.StatusBadRequest)
	}
}

func (s *ServerApi) ServeData(w http.ResponseWriter, r *http.Request, processor utils.XmlHandlerFunc) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	// 加密格式 的消息
	encryptMsg := EncryptMessage{}
	err = xml.Unmarshal(body, &encryptMsg)
	if err != nil {
		return
	}

	// 验证签名
	signature := calcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		encryptMsg.Encrypt,
		s.Token,
	)

	if msgSignature := r.URL.Query().Get("msg_signature"); signature != msgSignature {
		err = fmt.Errorf("%s != %s", signature, msgSignature)
		return
	}

	// 解密
	var xmlMsg []byte
	_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, s.EncodingAESKey)
	if err != nil {
		return
	}
	processor(w, r, xmlMsg)
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
func (s *ServerApi) ParseXML(body []byte) (m interface{}, err error) {
	message := Message{}
	err = xml.Unmarshal(body, &message)
	// fmt.Println(message)
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
	case MsgTypeEvent:
		return parseEvent(body)
	}

	return
}

// ParseEvent 解析微信推送过来的事件
func parseEvent(body []byte) (m interface{}, err error) {
	event := Event{}
	err = xml.Unmarshal(body, &event)
	if err != nil {
		return
	}
	switch event.Event {
	// 事件
	case EventTypeChangeContact:
		msg := EventChangeContact{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		switch msg.ChangeType {
		case EventTypeChangeContactCreateUser:
			msg := EventChangeContactCreateUser{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactUpdateUser:
			msg := EventChangeContactUpdateUser{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactDeleteUser:
			msg := EventChangeContactDeleteUser{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactCreateParty:
			msg := EventChangeContactCreateParty{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactUpdateParty:
			msg := EventChangeContactUpdateParty{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactDeleteParty:
			msg := EventChangeContactDeleteParty{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeContactUpdateTag:
			msg := EventChangeContactUpdateTag{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		}
	case EventTypeBatchJobResult:
		msg := EventBatchJobResult{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeApproval:
		msg := EventApproval{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeChangeExternalContact:
		msg := EventChangeExternalContact{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		switch msg.ChangeType {
		case EventTypeChangeExternalContactAddExternalContact:
			msg := EventChangeExternalContactAddExternalContact{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeExternalContactAddHalfExternalContact:
			msg := EventChangeExternalContactAddHalfExternalContact{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeExternalContactChangeExternalChat:
			msg := EventChangeExternalContactChangeExternalChat{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeExternalContactDelExternalContact:
			msg := EventChangeExternalContactDelExternalContact{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeExternalContactEditExternalContact:
			msg := EventChangeExternalContactEditExternalContact{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		case EventTypeChangeExternalContactDelFollowUser:
			msg := EventChangeExternalContactDelFollowUser{}
			err = xml.Unmarshal(body, &msg)
			if err != nil {
				return
			}
			return msg, nil
		}
	case EventTypeTaskCardClick:
		msg := EventTaskCardClick{}
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
	case EventTypeMenuClick:
		msg := EventMenuClick{}
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
	}

	return
}

// Response 响应微信消息
func (s *ServerApi) response(w http.ResponseWriter, r *http.Request, reply interface{}) (err error) {

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
func (s *ServerApi) encryptReplyMessage(rawXmlMsg []byte) (replyEncryptMessage *ReplyEncryptMessage, err error) {
	cipherText, err := utils.AESEncryptMsg(
		[]byte(utils.GetRandString(16)),
		rawXmlMsg,
		s.AgentConfig.AgentId,
		s.EncodingAESKey,
	)
	if err != nil {
		return
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

func (s *ServerApi) ResponseTaskCard(w http.ResponseWriter, r *http.Request, message *ReplyMessageTaskCard) (err error) {
	return s.response(w, r, message)
}

func (s *ServerApi) ResponseNews(w http.ResponseWriter, r *http.Request, message *ReplyMessageNews) (err error) {
	return s.response(w, r, message)
}
