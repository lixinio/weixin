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

package server_api

import (
	"encoding/xml"
	"strconv"
	"time"
)

const (
	EventTypeTaskCardClick = "taskcard_click" // 任务卡片事件
)

type Event struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
}

func (event *Event) Reply() *ReplyMessage {
	return &ReplyMessage{
		XMLName:      event.XMLName,
		ToUserName:   CDATA(event.FromUserName),
		FromUserName: CDATA(event.ToUserName),
		MsgType:      CDATA(event.MsgType),
		CreateTime:   strconv.FormatInt(time.Now().Unix(), 10),
	}
}

/*
<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[taskcard_click]]></Event>
<EventKey><![CDATA[key111]]></EventKey>
<TaskId><![CDATA[taskid111]]></TaskId >
<AgentId>1</AgentId>
</xml>
*/
type EventTaskCardClick struct {
	Event
	EventKey string `xml:"EventKey"`
	TaskID   string `xml:"TaskId"`
	AgentID  int    `xml:"AgentId"`
}
