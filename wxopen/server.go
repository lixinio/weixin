package wxopen

import (
	"encoding/xml"
	"errors"

	"github.com/lixinio/weixin/utils"
)

// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/authorize_event.html
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/token/component_verify_ticket.html
/*
启用 加密模式 后 收到的 消息格式
<xml>
    <ToUserName><![CDATA[]]></ToUserName>
    <Encrypt><![CDATA[]]></Encrypt>
</xml>
*/
type EncryptMessage struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
}

// ParseXML 解析微信推送过来的消息/事件
func (wxopen *WxOpen) ParseXML(
	body []byte,
	msgSignature, timestamp, nonce string,
) (m interface{}, err error) {

	// 是否加密消息
	encryptMsg := EncryptMessage{}
	err = xml.Unmarshal(body, &encryptMsg)
	if err != nil {
		return
	}

	// 需要解密
	if encryptMsg.Encrypt != "" {
		if utils.CalcSignature(
			wxopen.Config.Token,
			timestamp,
			nonce,
			encryptMsg.Encrypt,
		) != msgSignature {
			return nil, errors.New("invalid msg signature")
		}

		var xmlMsg []byte
		_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, wxopen.Config.EncodingAESKey)
		if err != nil {
			return
		}
		body = xmlMsg
	} else {
		return nil, errors.New("invalid encrypt msg")
	}

	event := Event{}
	err = xml.Unmarshal(body, &event)
	if err != nil {
		return
	}

	switch event.InfoType {

	case EventTypeComponentVerifyTicket:
		msg := EventComponentVerifyTicket{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeAuthorized:
		msg := EventAuthorized{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUnauthorized:
		msg := EventUnauthorized{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateAuthorized:
		msg := EventUpdateAuthorized{}
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			return
		}
		return msg, nil
	}
	return
}
