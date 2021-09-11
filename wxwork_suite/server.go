package wxwork_suite

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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

// ParseXML 解析微信推送过来的消息/事件
func (suite *WxWorkSuite) ParseXML(
	body []byte,
	msgSignature, timestamp, nonce string,
) (m interface{}, err error) {
	// 是否加密消息
	encryptMsg := &EncryptMessage{}
	if err = xml.Unmarshal(body, encryptMsg); err != nil {
		return
	}

	if encryptMsg.ToUserName != suite.Config.SuiteID {
		return nil, fmt.Errorf("invalid tousername %s", encryptMsg.ToUserName)
	}

	if utils.CalcSignature(
		suite.Config.Token,
		timestamp,
		nonce,
		encryptMsg.Encrypt,
	) != msgSignature {
		return nil, errors.New("invalid msg signature")
	}

	var xmlMsg []byte
	_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, suite.Config.EncodingAESKey)
	if err != nil {
		return
	}

	body = xmlMsg
	event := &Event{}
	if err = xml.Unmarshal(body, event); err != nil {
		return
	}

	// 继续校验数据
	if event.SuiteId != suite.Config.SuiteID {
		return nil, fmt.Errorf("invalid tousername %s", encryptMsg.ToUserName)
	}

	switch event.InfoType {
	case EventTypeSuiteTicket:
		msg := &EventSuiteTicket{}
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
	case EventTypeChangeAuthorized:
		msg := &EventUpdateAuthorized{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeChangeContact:
		return parseChangeContactEvent(body, event)
	}
	return nil, nil
}

func parseChangeContactEvent(body []byte, event *Event) (interface{}, error) {
	switch event.ChangeType {
	case EventTypeCreateUser:
		msg := &EventCreateUser{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	case EventTypeUpdateUser:
		msg := &EventUpdateUser{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	case EventTypeDeleteUser:
		msg := &EventDeleteUser{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	case EventTypeCreateParty:
		msg := &EventCreateParty{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	case EventTypeUpdateParty:
		msg := &EventUpdateParty{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	case EventTypeDeleteParty:
		msg := &EventDeleteParty{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	case EventTypeUpdateTag:
		msg := &EventUpdateTag{}
		if err := xml.Unmarshal(body, msg); err != nil {
			return nil, err
		}
		return msg, nil
	}
	return nil, nil
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

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func (suite *WxWorkSuite) ServeEcho(w http.ResponseWriter, r *http.Request) {
	signature, echoStr := calcSignatureFromHttp(r, suite.Config.Token)
	if echoStr != "" && signature == r.URL.Query().Get("msg_signature") {
		// 解密 echoStr
		_, msg, _, err := utils.AESDecryptMsg(echoStr, suite.Config.EncodingAESKey)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}
		w.Write(msg)
	} else {
		httpAbort(w, http.StatusBadRequest)
	}
}
