package message_api

import (
	"context"
	"net/url"
	"strconv"

	"github.com/lixinio/weixin/utils"
)

// 订阅消息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/sendMessage.html
// https://developers.weixin.qq.com/minigame/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
// https://developers.weixin.qq.com/doc/offiaccount/Subscription_Messages/api.html#send%E5%8F%91%E9%80%81%E8%AE%A2%E9%98%85%E9%80%9A%E7%9F%A5

const (
	apiSubscribeMsgSend                = "/cgi-bin/message/subscribe/send"
	apiSubscribeMsgBizSend             = "/cgi-bin/message/subscribe/bizsend" //
	apiSubscribeAddTemplate            = "/wxaapi/newtmpl/addtemplate"
	apiSubscribeGetTemplate            = "/wxaapi/newtmpl/gettemplate"
	apiSubscribeDelTemplate            = "/wxaapi/newtmpl/deltemplate"
	apiSubscribeGetCategory            = "/wxaapi/newtmpl/getcategory"
	apiSubscribeGetPubTemplateKeywords = "/wxaapi/newtmpl/getpubtemplatekeywords"
	apiSubscribeGetPubTemplateTitles   = "/wxaapi/newtmpl/getpubtemplatetitles"
)

const (
	// 跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	MiniProgramStateDeveloper = "developer"
	MiniProgramStateTrial     = "trail"
	MiniProgramStateFormal    = "formal"
)

type SendSubscribeMpMessageRequest struct {
	TemplateID       string            `json:"template_id"`       // 所需下发的订阅模板id
	Page             string            `json:"page,omitempty"`    // 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转
	Touser           string            `json:"touser"`            // 接收者（用户）的 openid
	Data             map[string]*Value `json:"data"`              // 模板内容，格式形如 { "key1": { "value": any }, "key2": { "value": any } }的object
	MiniprogramState string            `json:"miniprogram_state"` // 跳转小程序类型：developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Lang             string            `json:"lang,omitempty"`    // 进入小程序查看”的语言类型，支持zh_CN(简体中文)、en_US(英文)、zh_HK(繁体中文)、zh_TW(繁体中文)，默认为zh_CN
}

type Value struct {
	Value string `json:"value"`
}

// 发送(小程序)订阅消息
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/sendMessage.html
func (api MessageApi) SendSubscribeMpMessage(
	ctx context.Context,
	req *SendSubscribeMpMessageRequest,
	payload map[string]string,
) error {
	req.Data = map[string]*Value{}
	for k, v := range payload {
		req.Data[k] = &Value{Value: v}
	}

	if err := api.Client.HTTPPostJson(
		ctx, apiSubscribeMsgSend, req, nil,
	); err != nil {
		return err
	}

	return nil
}

type SendSubscribeOaMessageRequest struct {
	TemplateID  string             `json:"template_id"`    // 所需下发的订阅模板id
	Page        string             `json:"page,omitempty"` // 跳转网页时填写
	Touser      string             `json:"touser"`         // 接收者（用户）的 openid
	MiniProgram *TemplateMessageMp `json:"miniprogram,omitempty"`
	Data        map[string]*Value  `json:"data"` // 模板内容，格式形如 { "key1": { "value": any }, "key2": { "value": any } }的object
}

// 发送(服务号)订阅消息
// https://developers.weixin.qq.com/doc/offiaccount/Subscription_Messages/api.html#send%E5%8F%91%E9%80%81%E8%AE%A2%E9%98%85%E9%80%9A%E7%9F%A5
func (api MessageApi) SendSubscribeOaMessage(
	ctx context.Context,
	req *SendSubscribeOaMessageRequest,
	payload map[string]string,
) error {
	req.Data = map[string]*Value{}
	for k, v := range payload {
		req.Data[k] = &Value{Value: v}
	}

	if err := api.Client.HTTPPostJson(
		ctx, apiSubscribeMsgBizSend, req, nil,
	); err != nil {
		return err
	}

	return nil
}

// 添加模板
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/addMessageTemplate.html
// sceneDesc 官方文档可选， 实际不能为空
// sceneDesc 官方文档可选， 实际不能为空
// sceneDesc 官方文档可选， 实际不能为空
func (api MessageApi) SubscribeAddTemplate(
	ctx context.Context, tid int, kidList []int, sceneDesc string,
) (string, error) {
	resp := &struct {
		utils.WeixinError
		PriTmplID string `json:"priTmplId"`
	}{}
	if err := api.Client.HTTPPostJson(
		ctx, apiSubscribeAddTemplate, map[string]interface{}{
			"tid":       tid,
			"kidList":   kidList,
			"sceneDesc": sceneDesc,
		}, resp,
	); err != nil {
		return "", err
	}

	return resp.PriTmplID, nil
}

// 删除模板
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/deleteMessageTemplate.html
func (api MessageApi) SubscribeDelTemplate(
	ctx context.Context, priTmplId string,
) error {
	if err := api.Client.HTTPPostJson(
		ctx, apiSubscribeDelTemplate, map[string]interface{}{
			"priTmplId": priTmplId,
		}, nil,
	); err != nil {
		return err
	}

	return nil
}

type SubscribeTemplateKeyword struct {
	EnumValueList []string `json:"enumValueList"` // 枚举参数的 key
	KeywordCode   string   `json:"keywordCode"`   // 枚举参数值范围列表
}

type SubscribeTemplate struct {
	PriTmplID            string                      `json:"priTmplId"` // 添加至帐号下的模板 id，发送小程序订阅消息时所需
	Title                string                      `json:"title"`
	Content              string                      `json:"content"`
	Example              string                      `json:"example"`
	Type                 int                         `json:"type"`                 // 模版类型，2 为一次性订阅，3 为长期订阅
	KeywordEnumValueList []*SubscribeTemplateKeyword `json:"keywordEnumValueList"` // 枚举参数值范围
}

// 获取个人模板列表
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/getMessageTemplateList.html
func (api MessageApi) SubscribeGetTemplate(
	ctx context.Context,
) ([]*SubscribeTemplate, error) {
	resp := &struct {
		utils.WeixinError
		Data []*SubscribeTemplate `json:"data"`
	}{}
	if err := api.Client.HTTPGet(
		ctx, apiSubscribeGetTemplate, resp,
	); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

type SubscribeCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// 获取类目
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/getCategory.html
func (api MessageApi) SubscribeGetCategory(
	ctx context.Context,
) ([]*SubscribeCategory, error) {
	resp := &struct {
		utils.WeixinError
		Data []*SubscribeCategory `json:"data"`
	}{}
	if err := api.Client.HTTPGet(
		ctx, apiSubscribeGetCategory, resp,
	); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

type SubscribeKeyword struct {
	Kid     int    `json:"kid"` // 关键词 id，选用模板时需要
	Name    string `json:"name"`
	Example string `json:"example"`
	Rule    string `json:"rule"` // 参数类型
}

// 获取关键词列表
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/getPubTemplateKeyWordsById.html
func (api MessageApi) SubscribeGetPubTemplateKeywords(
	ctx context.Context, tid int,
) ([]*SubscribeKeyword, int, error) {
	resp := &struct {
		utils.WeixinError
		Count int                 `json:"count"`
		Data  []*SubscribeKeyword `json:"data"`
	}{}
	if err := api.Client.HTTPGetWithParams(
		ctx, apiSubscribeGetPubTemplateKeywords, func(v url.Values) {
			v.Add("tid", strconv.Itoa(tid))
		}, resp,
	); err != nil {
		return nil, 0, err
	}

	return resp.Data, resp.Count, nil
}

type SubscribePubTemplate struct {
	Tid        int    `json:"tid"` // 模版标题 id
	Title      string `json:"title"`
	Type       int    `json:"type"`       // 模版类型，2 为一次性订阅，3 为长期订阅
	CategoryID string `json:"categoryId"` // 模版所属类目 id (官方文档 类型是 number)
}

// 获取所属类目下的公共模板
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/mp-message-management/subscribe-message/getPubTemplateTitleList.html
func (api MessageApi) SubscribeGetPubTemplateTitles(
	ctx context.Context, ids string, start, limit int,
) ([]*SubscribePubTemplate, int, error) {
	resp := &struct {
		utils.WeixinError
		Count int                     `json:"count"`
		Data  []*SubscribePubTemplate `json:"data"`
	}{}
	if err := api.Client.HTTPGetWithParams(
		ctx, apiSubscribeGetPubTemplateTitles, func(v url.Values) {
			v.Add("ids", ids)
			v.Add("start", strconv.Itoa(start))
			v.Add("limit", strconv.Itoa(limit))
		}, resp,
	); err != nil {
		return nil, 0, err
	}

	return resp.Data, resp.Count, nil
}
