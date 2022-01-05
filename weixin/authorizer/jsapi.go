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

const (
	apiGetJSApiTicket = "/cgi-bin/ticket/getticket"
)

// 代公众号使用js sdk
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Before_Develop/Official_Accounts/js_sdk_instructions.html

// 通过config接口注入权限验证配置
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#4
type JsApiCorpConfig struct {
	Url       string `json:"url"`
	NonceStr  string `json:"nonceStr"`
	AppID     string `json:"appId"`
	TimeStamp string `json:"timestamp"`
	Signature string `json:"signature"`
}

// JS-SDK使用权限签名算法
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
func (api *Authorizer) GetJSApiConfig(
	ctx context.Context, url string,
) (*JsApiCorpConfig, error) {
	jsApiTicket, _, err := api.GetJSApiTicket(ctx)
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
		AppID:     api.Appid,
		TimeStamp: timestamp,
		Signature: signature,
	}, nil
}

/*
获取 jsapi_ticket

sapi_ticket是公众号用于调用微信JS接口的临时票据。正常情况下，jsapi_ticket的有效期为7200秒，通过access_token来获取。由于获取jsapi_ticket的api调用次数非常有限，频繁刷新jsapi_ticket会导致api调用受限，影响自身业务，开发者必须在自己的服务全局缓存jsapi_ticket

See: https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62

GET https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
*/
func (api *Authorizer) GetJSApiTicket(
	ctx context.Context,
) (jsapiTicket string, expiresIn int64, err error) {
	return api.getApiTicket(ctx, "jsapi")
}

/*
获取 wxcard_ticket

商户在调用授权页前需要先获取一个7200s过期的授权页ticket，在获取授权页接口中，该ticket作为参数传入，加强安全性。

See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html#1

GET https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=wx_card
*/
func (api *Authorizer) GetWxCardApiTicket(
	ctx context.Context,
) (jsapiTicket string, expiresIn int64, err error) {
	return api.getApiTicket(ctx, "wx_card")
}

func (api *Authorizer) getApiTicket(
	ctx context.Context, tp string,
) (jsapiTicket string, expiresIn int64, err error) {
	jsapiTicketResp := struct {
		utils.WeixinError
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
	}{}

	if err = api.Client.HTTPGetWithParams(ctx, apiGetJSApiTicket, func(params url.Values) {
		params.Add("type", tp)
	}, &jsapiTicketResp); err != nil {
		return "", 0, err
	}

	return jsapiTicketResp.Ticket, jsapiTicketResp.ExpiresIn, nil
}
