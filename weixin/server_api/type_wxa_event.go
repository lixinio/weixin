package server_api

// 只有服务商才有
const (
	EventTypeWxaNickNameAudit  = "wxa_nickname_audit"  // 名称审核结果事件推送
	EventTypeWxaCategoryAudit  = "wxa_category_audit"  // 类目审核结果事件推送
	EventTypeWeappAuditSuccess = "weapp_audit_success" // 代码审核结果推送 审核通过
	EventTypeWeappAuditFail    = "weapp_audit_fail"    // 代码审核结果推送 审核不通过
	EventTypeWeappAuditDelay   = "weapp_audit_delay"   // 代码审核结果推送 审核延后
)

// 名称审核结果事件推送
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/wxa_nickname_audit.html
/*
<xml>
	<ToUserName><![CDATA[gh_fxxxxxxxa4b2]]></ToUserName>
	<FromUserName><![CDATA[odxxxxM-xxxxxxxx-trm4a7apsU8]]></FromUserName>
	<CreateTime>1488800000</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[wxa_nickname_audit]]></Event>
	<ret>2</ret>
	<nickname>昵称</nickname>
	<reason>驳回原因</reason>
</xml>
*/
type EventWxaNickNameAudit struct {
	Event
	Ret      int    `xml:"ret"`      // 审核结果 2：失败，3：成功
	Nickname string `xml:"nickname"` // 需要更改的昵称
	Reason   string `xml:"reason"`   // 审核失败的驳回原因
}

// 代码审核结果推送
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/audit_event.html

// 审核通过
/*
<xml>
  <ToUserName><![CDATA[gh_fb9688c2a4b2]]></ToUserName>
  <FromUserName><![CDATA[od1P50M-fNQI5Gcq-trm4a7apsU8]]></FromUserName>
  <CreateTime>1488856741</CreateTime>
  <MsgType><![CDATA[event]]></MsgType>
  <Event><![CDATA[weapp_audit_success]]></Event>
  <SuccTime>1488856741</SuccTime>
</xml>
*/
type EventWeappAuditSuccess struct {
	Event
	SuccTime int
}

// 审核不通过
/*
<xml>
  <ToUserName><![CDATA[gh_fb9688c2a4b2]]></ToUserName>
  <FromUserName><![CDATA[od1P50M-fNQI5Gcq-trm4a7apsU8]]></FromUserName>
  <CreateTime>1488856591</CreateTime>
  <MsgType><![CDATA[event]]></MsgType>
  <Event><![CDATA[weapp_audit_fail]]></Event>
  <Reason><![CDATA[1:账号信息不符合规范:<br>(1):包含色情因素<br>2:服务类目"金融业-保险_"与你提交代码审核时设置的功能页面内容不一致:<br>(1):功能页面设置的部分标签不属于所选的服务类目范围。<br>(2):功能页面设置的部分标签与该页面内容不相关。<br>]]></Reason>
  <FailTime>1488856591</FailTime>
  <ScreenShot>xxx|yyy|zzz</ScreenShot>
</xml>
*/
type EventWeappAuditFail struct {
	Event
	Reason     string
	FailTime   int
	ScreenShot string
}

// 审核延后
/*
<xml>
  <ToUserName><![CDATA[gh_fb9688c2a4b2]]></ToUserName>
  <FromUserName><![CDATA[od1P50M-fNQI5Gcq-trm4a7apsU8]]></FromUserName>
  <CreateTime>1488856591</CreateTime>
  <MsgType><![CDATA[event]]></MsgType>
  <Event><![CDATA[weapp_audit_delay]]></Event>
  <Reason><![CDATA[为了更好的服务小程序，您的服务商正在进行提审系统的优化，可能会导致审核时效的增长，请耐心等待]]></Reason>
  <DelayTime>1488856591</DelayTime>
</xml>
*/
type EventWeappAuditDelay struct {
	Event
	Reason    string
	DelayTime int
}

// 类目审核结果事件推送
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/category/wxa_category_audit.html
/*
<xml>
	<ToUserName><![CDATA[gh_fxxxxxxxa4b2]]></ToUserName>
	<FromUserName><![CDATA[odxxxxM-xxxxxxxx-trm4a7apsU8]]></FromUserName>
	<CreateTime>1488800000</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[wxa_category_audit]]></Event>
	<ret>2</ret>
	<first>一级类目id</nickname>
	<second>二级类目id</reason>
      <reason>驳回原因</reason>
</xml>
*/
type EventWxaCategoryAudit struct {
	Event
	Reason string `xml:"reason"` // 审核失败的驳回原因
	Ret    int    `xml:"ret"`    // 审核结果 2.驳回，3通过
	First  string `xml:"first"`
	Second string `xml:"second"`
}
