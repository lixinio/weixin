package authorizer

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/lixinio/weixin/utils"
)

// 企业微信自建应用和服务商接口一样
// wxwork/authorizer/jsapi.go
// wxwork/agent/jsapi.go

const (
	apiGetAgentJSApiTicket = "/cgi-bin/ticket/get"
	apiGetCorpJSApiTicket  = "/cgi-bin/get_jsapi_ticket"
)

/*
获取企业的jsapi_ticket
https://work.weixin.qq.com/api/doc/90001/90144/90539#%E8%8E%B7%E5%8F%96%E4%BC%81%E4%B8%9A%E7%9A%84jsapi_ticket
https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=ACCESS_TOKEN
*/
func (authorizer *Authorizer) getCorpJSApiTicket(
	ctx context.Context,
) (jsapiTicket string, expiresIn int64, err error) {
	jsapiTicketResp := struct {
		utils.WeixinError
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	if err = authorizer.Client.HTTPGet(
		ctx, apiGetCorpJSApiTicket, &jsapiTicketResp,
	); err != nil {
		return "", 0, err
	}

	return jsapiTicketResp.Ticket, jsapiTicketResp.ExpiresIn, nil
}

func (authorizer *Authorizer) GetCorpJSApiTicket(
	ctx context.Context,
) (jsapiTicket string, err error) {
	if authorizer.corpJsApiTicketCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrCorpJsApiTicketForbidden,
		)
	}
	return authorizer.corpJsApiTicketCache.GetAccessToken()
}

type JsApiCorpConfig struct {
	Url       string `json:"url"`
	NonceStr  string `json:"nonceStr"`
	AppID     string `json:"appid"`
	TimeStamp string `json:"timestamp"`
	Signature string `json:"signature"`
}

// https://work.weixin.qq.com/api/doc/90000/90136/90506
func (authorizer *Authorizer) GetCorpJSApiConfig(
	ctx context.Context, url string,
) (*JsApiCorpConfig, error) {
	jsApiTicket, err := authorizer.GetCorpJSApiTicket(ctx)
	if err != nil {
		return nil, err
	}

	nonceStr := utils.GetRandString(6)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	plain := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s",
		jsApiTicket, nonceStr, timestamp, url,
	)
	signature := fmt.Sprintf("%x", sha1.Sum([]byte(plain)))

	return &JsApiCorpConfig{
		Url:       url,
		NonceStr:  nonceStr,
		AppID:     authorizer.CorpID,
		TimeStamp: timestamp,
		Signature: signature,
	}, nil
}

/*
获取应用的jsapi_ticket
https://work.weixin.qq.com/api/doc/90001/90144/90539#%E8%8E%B7%E5%8F%96%E5%BA%94%E7%94%A8%E7%9A%84jsapi_ticket
https://qyapi.weixin.qq.com/cgi-bin/ticket/get?access_token=ACCESS_TOKEN&type=agent_config
*/
func (authorizer *Authorizer) getAgentJSApiTicket(
	ctx context.Context,
) (jsapiTicket string, expiresIn int64, err error) {
	return authorizer.getAgentApiTicket(ctx, "agent_config")
}

func (authorizer *Authorizer) GetAgentJSApiTicket(
	ctx context.Context,
) (jsapiTicket string, err error) {
	if authorizer.agentJsApiTicketCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrAgentJsApiTicketForbidden,
		)
	}
	return authorizer.agentJsApiTicketCache.GetAccessToken()
}

type JsApiAgentConfig struct {
	Url       string `json:"url"`
	NonceStr  string `json:"nonceStr"`
	CorpID    string `json:"corpid"`
	AgentID   int    `json:"agentid"`
	TimeStamp string `json:"timestamp"`
	Signature string `json:"signature"`
}

// https://work.weixin.qq.com/api/doc/90000/90136/90506
func (authorizer *Authorizer) GetAgentJSApiConfig(
	ctx context.Context, url string,
) (*JsApiAgentConfig, error) {
	jsApiTicket, err := authorizer.GetAgentJSApiTicket(ctx)
	if err != nil {
		return nil, err
	}

	nonceStr := utils.GetRandString(6)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	plain := fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s",
		jsApiTicket, nonceStr, timestamp, url,
	)
	signature := fmt.Sprintf("%x", sha1.Sum([]byte(plain)))

	return &JsApiAgentConfig{
		Url:       url,
		NonceStr:  nonceStr,
		CorpID:    authorizer.CorpID,
		AgentID:   authorizer.AgentID,
		TimeStamp: timestamp,
		Signature: signature,
	}, nil
}

func (authorizer *Authorizer) getAgentApiTicket(
	ctx context.Context, tp string,
) (jsapiTicket string, expiresIn int64, err error) {
	jsapiTicketResp := struct {
		utils.WeixinError
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	if err = authorizer.Client.HTTPGetWithParams(
		ctx, apiGetAgentJSApiTicket, func(params url.Values) {
			params.Add("type", tp)
		}, &jsapiTicketResp,
	); err != nil {
		return "", 0, err
	}

	return jsapiTicketResp.Ticket, jsapiTicketResp.ExpiresIn, nil
}
