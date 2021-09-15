package wxwork_suite

import (
	"encoding/xml"
	"errors"
	"fmt"
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

func (suite *WxWorkSuite) ServeData(
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

	if encryptMsg.ToUserName != suite.Config.SuiteID {
		return fmt.Errorf("invalid tousername %s", encryptMsg.ToUserName)
	}

	// 验证签名
	signature := utils.CalcSignature(
		r.URL.Query().Get("timestamp"),
		r.URL.Query().Get("nonce"),
		encryptMsg.Encrypt,
		suite.Config.Token,
	)
	if msgSignature := r.URL.Query().Get("msg_signature"); signature != msgSignature {
		err = fmt.Errorf("signature dismatch %s != %s", signature, msgSignature)
		return err
	}

	// 解密
	var xmlMsg []byte
	_, xmlMsg, _, err = utils.AESDecryptMsg(encryptMsg.Encrypt, suite.Config.EncodingAESKey)
	if err != nil {
		return err
	}
	return processor(w, r, xmlMsg)
}

// ParseXML 解析微信推送过来的消息/事件
func (suite *WxWorkSuite) ParseXML(body []byte) (m interface{}, err error) {
	event := &Event{}
	if err = xml.Unmarshal(body, event); err != nil {
		return
	}

	// 继续校验数据
	if event.SuiteId != suite.Config.SuiteID {
		return nil, fmt.Errorf("invalid suite id %s,%s", event.SuiteId, suite.Config.SuiteID)
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
	return
}

func parseChangeContactEvent(body []byte, event *Event) (m interface{}, err error) {
	switch event.ChangeType {
	case EventTypeCreateUser:
		msg := &EventCreateUser{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateUser:
		msg := &EventUpdateUser{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeDeleteUser:
		msg := &EventDeleteUser{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeCreateParty:
		msg := &EventCreateParty{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateParty:
		msg := &EventUpdateParty{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeDeleteParty:
		msg := &EventDeleteParty{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	case EventTypeUpdateTag:
		msg := &EventUpdateTag{}
		if err = xml.Unmarshal(body, msg); err != nil {
			return
		}
		return msg, nil
	}
	return
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

func (suite *WxWorkSuite) ServeEcho(w http.ResponseWriter, r *http.Request) error {
	signature, echoStr := calcSignatureFromHttp(r, suite.Config.Token)
	if echoStr != "" && signature == r.URL.Query().Get("msg_signature") {
		// 解密 echoStr
		_, msg, _, err := utils.AESDecryptMsg(echoStr, suite.Config.EncodingAESKey)
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
