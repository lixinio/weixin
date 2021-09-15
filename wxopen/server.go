package wxopen

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"

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

func (wxopen *WxOpen) ServeData(
	w http.ResponseWriter,
	r *http.Request,
	processor utils.XmlHandlerFunc,
) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// 加密格式 的消息
	encryptMsg := &EncryptMessage{}
	if err = xml.Unmarshal(body, encryptMsg); err != nil {
		return err
	}

	msgSignature := r.URL.Query().Get("msg_signature")
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	if utils.CalcSignature(
		wxopen.Config.Token,
		timestamp,
		nonce,
		encryptMsg.Encrypt,
	) != msgSignature {
		return errors.New("invalid msg signature")
	}

	// 解密
	var xmlMsg []byte
	_, xmlMsg, _, err = utils.AESDecryptMsg(
		encryptMsg.Encrypt, wxopen.Config.EncodingAESKey,
	)
	if err != nil {
		return err
	}
	return processor(w, r, xmlMsg)
}

// ParseXML 解析微信推送过来的消息/事件
func (wxopen *WxOpen) ParseXML(body []byte) (m interface{}, err error) {
	event := &Event{}
	if err = xml.Unmarshal(body, event); err != nil {
		return
	}

	switch event.InfoType {
	case EventTypeComponentVerifyTicket:
		msg := &EventComponentVerifyTicket{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeAuthorized:
		msg := &EventAuthorized{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUnauthorized:
		msg := &EventUnauthorized{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateAuthorized:
		msg := &EventUpdateAuthorized{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	}
	return
}
