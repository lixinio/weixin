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

package wxwork_suite

import "encoding/xml"

type EventExternalBase struct {
	XMLName    xml.Name `xml:"xml"`
	SuiteID    string   `xml:"SuiteId"`    // 第三方应用ID
	AuthCorpID string   `xml:"AuthCorpId"` // 授权企业的CorpID
	InfoType   string   `xml:"InfoType"`   // 固定为change_external_contact
	TimeStamp  int64    `xml:"TimeStamp"`  // 时间戳
	ChangeType string   `xml:"ChangeType"` //
}

type EventExternalContact struct {
	EventExternalBase
	UserID         string `xml:"UserID"`         // 企业服务人员的UserID
	ExternalUserID string `xml:"ExternalUserID"` // 外部联系人的userid，注意不是企业成员的账号
}

// https://developer.work.weixin.qq.com/document/path/92277#%E5%AE%A2%E6%88%B7%E6%8E%A5%E6%9B%BF%E5%A4%B1%E8%B4%A5%E4%BA%8B%E4%BB%B6
const (
	EventSubTypeAddExternalContact     = "add_external_contact"      // 添加企业客户事件
	EventSubTypeEditExternalContact    = "edit_external_contact"     // 编辑企业客户事件
	EventSubTypeAddHalfExternalContact = "add_half_external_contact" // 外部联系人免验证添加成员事件
	EventSubTypeDelExternalContact     = "del_external_contact"      // 删除企业客户事件
	EventSubTypeDelFollowUser          = "del_follow_user"           // 删除跟进成员事件
	EventSubTypeCustomerRefused        = "transfer_fail"             // 客户接替失败事件
)

const (
	EventSubTypeExternalChatCreate             = "create"        // 客户群创建事件
	EventSubTypeExternalChatUpdate             = "update"        // 客户群变更事件
	EventSubTypeExternalChatDismiss            = "dismiss"       // 客户群解散事件
	EventSubTypeExternalChatUpdateAddMember    = "add_member"    // 成员入群
	EventSubTypeExternalChatUpdateDelMember    = "del_member"    // 成员退群
	EventSubTypeExternalChatUpdateChangeOwner  = "change_owner"  // 群主变更
	EventSubTypeExternalChatUpdateChangeName   = "change_name"   // 群名变更
	EventSubTypeExternalChatUpdateChangeNotice = "change_notice" // 群公告变更
)

const (
	EventSubTypeExternalTagCreate  = "create"  // 企业客户标签创建事件
	EventSubTypeExternalTagUpdate  = "update"  // 企业客户标签变更事件
	EventSubTypeExternalTagDelete  = "delete"  // 企业客户标签删除事件
	EventSubTypeExternalTagShuffle = "shuffle" // 企业客户标签重排事件
)

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_contact]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChangeType><![CDATA[add_external_contact]]></ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<ExternalUserID><![CDATA[woAJ2GCAAAXtWyujaWJHDDGi0mACH71w]]></ExternalUserID>
	<State><![CDATA[teststate]]></State>
	<WelcomeCode><![CDATA[WELCOMECODE]]></WelcomeCode>

</xml>
*/
type EventAddExternalContact struct {
	EventExternalContact
	State       string `xml:"State"`       // 添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
	WelcomeCode string `xml:"WelcomeCode"` // 欢迎语code，可用于发送欢迎语
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_contact]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChangeType><![CDATA[edit_external_contact]]></ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<ExternalUserID><![CDATA[woAJ2GCAAAXtWyujaWJHDDGi0mACH71w]]></ExternalUserID>

</xml>
*/
type EventEditExternalContact struct {
	EventExternalContact
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_contact]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChangeType><![CDATA[add_half_external_contact]]></ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<ExternalUserID><![CDATA[woAJ2GCAAAXtWyujaWJHDDGi0mACH71w]]></ExternalUserID>
	<State><![CDATA[teststate]]></State>
	<WelcomeCode><![CDATA[WELCOMECODE]]></WelcomeCode>

</xml>
*/
type EventAddHalfExternalContact struct {
	EventExternalContact
	State       string `xml:"State"`       // 添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
	WelcomeCode string `xml:"WelcomeCode"` // 欢迎语code，可用于发送欢迎语
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_contact]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChangeType><![CDATA[del_external_contact]]></ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<ExternalUserID><![CDATA[woAJ2GCAAAXtWyujaWJHDDGi0mACH71w]]></ExternalUserID>
	<Source><![CDATA[DELETE_BY_TRANSFER]]></Source>

</xml>
*/
type EventDelExternalContact struct {
	EventExternalContact
	// 删除客户的操作来源，DELETE_BY_TRANSFER表示此客户是因在职继承自动被转接成员删除
	Source string `xml:"Source"`
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_contact]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChangeType><![CDATA[del_follow_user]]></ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<ExternalUserID><![CDATA[woAJ2GCAAAXtWyujaWJHDDGi0mACH71w]]></ExternalUserID>

</xml>
*/
type EventDelExternalContactFollowUser struct {
	EventExternalContact
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_contact]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChangeType><![CDATA[transfer_fail]]></ChangeType>
	<FailReason><![CDATA[customer_refused]]></FailReason>
	<UserID><![CDATA[zhangsan]]></UserID>
	<ExternalUserID><![CDATA[woAJ2GCAAAXtWyujaWJHDDGi0mACH71w]]></ExternalUserID>

</xml>
*/
type EventRefuseExternalContactCustomer struct {
	EventExternalContact
	// 接替失败的原因, customer_refused-客户拒绝， customer_limit_exceed-接替成员的客户数达到上限
	FailReason string `xml:"FailReason"`
}

type EventExternalChat struct {
	EventExternalBase
	ChatID string `xml:"ChatId"` // 群ID
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_chat]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChatId><![CDATA[CHAT_ID]]></ChatId>
	<ChangeType><![CDATA[create]]></ChangeType>

</xml>
*/
// 客户群创建事件
type EventCreateExternalChat struct {
	EventExternalChat
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_chat]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChatId><![CDATA[CHAT_ID]]></ChatId>
	<ChangeType><![CDATA[update]]></ChangeType>
	<UpdateDetail><![CDATA[add_member]]></UpdateDetail>
	<JoinScene>1</JoinScene>
	<QuitScene>0</QuitScene>
	<MemChangeCnt>10</MemChangeCnt>
	<MemChangeList>
		<Item>Jack</Item>
		<Item>Rose</Item>
	</MemChangeList>
	<LastMemVer>9c3f97c2ada667dfb5f6d03308d963e1</LastMemVer>
	<CurMemVer>71217227bbd112ecfe3a49c482195cb4</CurMemVer>

</xml>
*/
// 客户群变更事件
type EventUpdateExternalChat struct {
	EventExternalChat
	UpdateDetail  string `xml:"UpdateDetail"`
	JoinScene     int    `xml:"JoinScene"`
	QuitScene     int    `xml:"QuitScene"`
	MemChangeList []struct {
		Item string `xml:"Item"`
	} `xml:"MemChangeList"`
	MemChangeCnt int    `xml:"MemChangeCnt"`
	LastMemVer   string `xml:"LastMemVer"`
	CurMemVer    string `xml:"CurMemVer"`
}

/*
<xml>
	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_chat]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<ChatId><![CDATA[CHAT_ID]]></ChatId>
	<ChangeType><![CDATA[dismiss]]></ChangeType>
</xml>
*/
// 客户群解散事件
type EventDismissExternalChat struct {
	EventExternalChat
}

type EventExternalTag struct {
	EventExternalBase
	ID string `xml:"Id"` // 标签或标签组的ID
}

/*
<xml>

	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_tag]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<Id><![CDATA[TAG_ID]]></Id>
	<TagType><![CDATA[tag]]></TagType>
	<ChangeType><![CDATA[create]]></ChangeType>

</xml>
*/
// 企业客户标签创建事件
type EventCreateExternalTag struct {
	EventExternalTag
	TagType string `xml:"TagType"` // 创建标签时，此项为tag，创建标签组时，此项为tag_group
}

/*
<xml>
	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_tag]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<Id><![CDATA[TAG_ID]]></Id>
	<TagType><![CDATA[tag]]></TagType>
	<ChangeType><![CDATA[update]]></ChangeType>
</xml>
*/
// 企业客户标签变更事件
type EventUpdateExternalTag struct {
	EventExternalTag
}

/*
<xml>
	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_tag]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<Id><![CDATA[TAG_ID]]></Id>
	<TagType><![CDATA[tag]]></TagType>
	<ChangeType><![CDATA[delete]]></ChangeType>
</xml>
*/
// 企业客户标签删除事件
type EventDeleteExternalTag struct {
	EventExternalTag
	TagType string `xml:"TagType"` // 创建标签时，此项为tag，创建标签组时，此项为tag_group
}

/*
<xml>
	<SuiteId><![CDATA[ww4asffe99e54c0f4c]]></SuiteId>
	<AuthCorpId><![CDATA[wxf8b4f85f3a794e77]]></AuthCorpId>
	<InfoType><![CDATA[change_external_tag]]></InfoType>
	<TimeStamp>1403610513</TimeStamp>
	<Id><![CDATA[TAG_ID]]></Id>
	<ChangeType><![CDATA[shuffle]]></ChangeType>
</xml>
*/
// 企业客户标签重排事件
type EventShuffleExternalTag struct {
	EventExternalTag
}
