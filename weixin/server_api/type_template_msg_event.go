package server_api

// https://github.com/fastwego/offiaccount/blob/master/type/type_event/type_template_msg_event.go

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

const (
	EventTypeTemplateSendJobFinish = "TEMPLATESENDJOBFINISH" // 模版消息发送任务完成
	EventTypeMassSendJobFinish     = "MASSSENDJOBFINISH"     // 群发消息任务完成
)

/*
<xml>

	<ToUserName><![CDATA[gh_7f083739789a]]></ToUserName>
	<FromUserName><![CDATA[oia2TjuEGTNoeX76QEjQNrcURxG8]]></FromUserName>
	<CreateTime>1395658920</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[TEMPLATESENDJOBFINISH]]></Event>
	<MsgID>200163836</MsgID>
	<Status><![CDATA[success]]></Status>

</xml>
*/
type EventTemplateSendJobFinish struct {
	Event
	MsgID  int64
	Status string
}

// https://developers.weixin.qq.com/doc/service/guide/product/message/Batch_Sends.html
/*
<xml>
	<ToUserName><![CDATA[gh_4d00ed8d6399]]></ToUserName>
	<FromUserName><![CDATA[oV5CrjpxgaGXNHIQigzNlgLTnwic]]></FromUserName>
	<CreateTime>1481013459</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[MASSSENDJOBFINISH]]></Event>
	<MsgID>1000001625</MsgID>
	<Status><![CDATA[err(30003)]]></Status>
	<TotalCount>0</TotalCount>
	<FilterCount>0</FilterCount>
	<SentCount>0</SentCount>
	<ErrorCount>0</ErrorCount>
	<CopyrightCheckResult>
		<Count>2</Count>
		<ResultList>
			<item>
				<ArticleIdx>1</ArticleIdx>
				<UserDeclareState>0</UserDeclareState>
				<AuditState>2</AuditState>
				<OriginalArticleUrl><![CDATA[Url_1]]></OriginalArticleUrl>
				<OriginalArticleType>1</OriginalArticleType>
				<CanReprint>1</CanReprint>
				<NeedReplaceContent>1</NeedReplaceContent>
				<NeedShowReprintSource>1</NeedShowReprintSource>
			</item>
			<item>
				<ArticleIdx>2</ArticleIdx>
				<UserDeclareState>0</UserDeclareState>
				<AuditState>2</AuditState>
				<OriginalArticleUrl><![CDATA[Url_2]]></OriginalArticleUrl>
				<OriginalArticleType>1</OriginalArticleType>
				<CanReprint>1</CanReprint>
				<NeedReplaceContent>1</NeedReplaceContent>
				<NeedShowReprintSource>1</NeedShowReprintSource>
			</item>
		</ResultList>
		<CheckState>2</CheckState>
	</CopyrightCheckResult>
	<ArticleUrlResult>
		<Count>1</Count>
		<ResultList>
			<item>
				<ArticleIdx>1</ArticleIdx>
				<ArticleUrl><![CDATA[Url]]></ArticleUrl>
			</item>
		</ResultList>
	</ArticleUrlResult>
</xml>
*/
type EventMassSendJobFinish struct {
	Event
	MsgID                int64
	Status               string
	TotalCount           int32
	FilterCount          int32
	SentCount            int32
	ErrorCount           int32
	CopyrightCheckResult struct {
		Count      int32
		ResultList struct {
			Items []struct {
				ArticleIdx            int32
				UserDeclareState      int32
				AuditState            int32
				OriginalArticleUrl    string
				OriginalArticleType   int32
				CanReprint            int32
				NeedReplaceContent    int32
				NeedShowReprintSource int32
			} `xml:"item"`
		}
	}
	ArticleUrlResult struct {
		Count      int32
		ResultList struct {
			Items []struct {
				ArticleIdx int32
				ArticleUrl string
			} `xml:"item"`
		}
	}
}
