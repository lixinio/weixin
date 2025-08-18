package wxwork_suite

import (
	"encoding/xml"
)

const (
	EventTypeSuiteTicket           = "suite_ticket"
	EventTypeAuthorized            = "create_auth"
	EventTypeUnauthorized          = "cancel_auth"
	EventTypeChangeAuthorized      = "change_auth"
	EventTypeChangeContact         = "change_contact"
	EventTypeChangeExternalContact = "change_external_contact" // 企业客户变更
	EventTypeChangeExternalChat    = "change_external_chat"    // 客户群事件
	EventTypeChangeExternalTag     = "change_external_tag"     // 企业客户标签事件
)

type Event struct {
	XMLName    xml.Name `xml:"xml"`
	SuiteId    string
	TimeStamp  string
	InfoType   string
	ChangeType string // 成员通知事件 & 部门通知事件 还会细分， 这里一次性读出避免重复序列化
}

// https://open.work.weixin.qq.com/api/doc/90001/90143/90628
/*
<xml>
	<SuiteId><![CDATA[ww4asffe99e54c0fxxxx]]></SuiteId>
	<InfoType> <![CDATA[suite_ticket]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<SuiteTicket><![CDATA[asdfasfdasdfasdf]]></SuiteTicket>
</xml>
*/
type EventSuiteTicket struct {
	Event
	SuiteTicket string
}

// https://open.work.weixin.qq.com/api/doc/90001/90143/90642
/*
授权成功通知
<xml>
    <SuiteId><![CDATA[ww4asffe9xxx4c0f4c]]></SuiteId>
    <AuthCode><![CDATA[AUTHCODE]]></AuthCode>
    <InfoType><![CDATA[create_auth]]></InfoType>
    <TimeStamp>1403610513</TimeStamp>
    <State><![CDATA[123]]></State>
</xml>
*/
type EventAuthorized struct {
	Event
	AuthCode string
	State    string
}

/*
取消授权通知
<xml>

	<SuiteId><![CDATA[ww4asffe99e54cxxxx]]></SuiteId>
	<InfoType><![CDATA[cancel_auth]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<AuthCorpId><![CDATA[wxf8b4f85fxx794xxx]]></AuthCorpId>

</xml>
*/
type EventUnauthorized struct {
	Event
	AuthCorpId string
}

/*
变更授权通知
<xml>

	<SuiteId><![CDATA[ww4asffe99exxx0f4c]]></SuiteId>
	<InfoType><![CDATA[change_auth]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>

</xml>
*/
type EventUpdateAuthorized struct {
	Event
	AuthCorpId string
}
