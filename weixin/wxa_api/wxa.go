package wxa_api

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGenerateUrlLink    = "/wxa/generate_urllink"
	apiQueryUrlLink       = "/wxa/query_urllink"
	apiGenerateScheme     = "/wxa/generatescheme"
	apiQueryScheme        = "/wxa/queryscheme"
	apiGetUserPhonenumber = "/wxa/business/getuserphonenumber"
)

type WxaApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *WxaApi {
	return &WxaApi{Client: client}
}

/*
获取小程序 URL Link，适用于短信、邮件、网页、微信内等拉起小程序的业务场景。
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/url-link/urllink.generate.html
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Business/urllink.generate.html
*/

type GenerateUrlLinkRequest struct {
	Path       string `json:"path"`
	Query      string `json:"query,omitempty"`
	EnvVersion string `json:"env_version,omitempty"`
	IsExpire   bool   `json:"is_expire"`
	ExpireTime int64  `json:"expire_time,omitempty"`
}

func (api *WxaApi) GenerateUrlLink(
	ctx context.Context, param *GenerateUrlLinkRequest,
) (string, error) {
	result := &struct {
		utils.WeixinError
		UrlLink string `json:"url_link"`
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiGenerateUrlLink, param, result); err != nil {
		return "", err
	}
	return result.UrlLink, nil
}

/**
查询小程序 url_link 配置，及长期有效 quota。
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/url-link/urllink.query.html
**/
type UrlLinkInfo struct {
	UrlLinkInfo struct {
		AppID      string `json:"appid"`
		Path       string `json:"path"`
		Query      string `json:"query,omitempty"`
		ExpireTime int64  `json:"expire_time,omitempty"`
		CreateTime int64  `json:"create_time"`
		EnvVersion string `json:"env_version,omitempty"`
	} `json:"url_link_info"`
	UrlLinkQuota struct {
		LongTimeUsed  int `json:"long_time_used"`
		LongTimeLimit int `json:"long_time_limit"`
	} `json:"url_link_quota"`
}

func (api *WxaApi) GetUrlLink(
	ctx context.Context, urlLink string,
) (*UrlLinkInfo, error) {
	param := &struct {
		UrlLink string `json:"url_link"`
	}{urlLink}

	result := &struct {
		utils.WeixinError
		UrlLinkInfo
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiQueryUrlLink, param, result); err != nil {
		return nil, err
	}
	return &result.UrlLinkInfo, nil
}

/**
获取小程序 scheme 码，适用于短信、邮件、外部网页、微信内等拉起小程序的业务场景。
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/url-scheme/urlscheme.generate.html
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Business/url_scheme.html
**/
type GenerateSchemeRequest struct {
	JumpWxa *struct {
		Path       string `json:"path"`
		Query      string `json:"query,omitempty"`
		EnvVersion string `json:"env_version,omitempty"`
	} `json:"jump_wxa"`
	IsExpire   bool  `json:"is_expire"`
	ExpireTime int64 `json:"expire_time,omitempty"`
}

func (api *WxaApi) GenerateScheme(
	ctx context.Context, param *GenerateSchemeRequest,
) (string, error) {
	result := &struct {
		utils.WeixinError
		OpenLink string `json:"openlink"`
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiGenerateScheme, param, result); err != nil {
		return "", err
	}
	return result.OpenLink, nil
}

/**
查询小程序 scheme 码，及长期有效 quota。
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/url-scheme/urlscheme.query.html
**/
type SchemaInfo struct {
	SchemeInfo struct {
		AppID      string `json:"appid"`
		Path       string `json:"path"`
		Query      string `json:"query,omitempty"`
		ExpireTime int64  `json:"expire_time,omitempty"`
		CreateTime int64  `json:"create_time"`
		EnvVersion string `json:"env_version,omitempty"`
	} `json:"scheme_info"`
	SchemeQuota struct {
		LongTimeUsed  int `json:"long_time_used"`
		LongTimeLimit int `json:"long_time_limit"`
	} `json:"scheme_quota"`
}

func (api *WxaApi) GetSchema(
	ctx context.Context, scheme string,
) (*SchemaInfo, error) {
	param := &struct {
		Scheme string `json:"scheme"`
	}{scheme}

	result := &struct {
		utils.WeixinError
		SchemaInfo
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiQueryScheme, param, result); err != nil {
		return nil, err
	}
	return &result.SchemaInfo, nil
}

type PhoneInfoWatermark struct {
	Timestamp int64  `json:"timestamp"`
	AppID     string `json:"appid"`
}

type PhoneInfo struct {
	PhoneNumber     string              `json:"phoneNumber"`
	PurePhoneNumber string              `json:"purePhoneNumber"`
	CountryCode     string              `json:"countryCode"`
	Watermark       *PhoneInfoWatermark `json:"watermark"`
}

// 该接口用于将code换取用户手机号。 说明，每个code只能使用一次，code的有效期为5min。
// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-info/phone-number/getPhoneNumber.html
func (api *WxaApi) GetPhoneNumber(
	ctx context.Context, code, openid string,
) (*PhoneInfo, error) {
	param := &struct {
		Code   string `json:"code"`
		OpenID string `json:"openid"`
	}{
		Code:   code,
		OpenID: openid,
	}

	result := &struct {
		utils.WeixinError
		PhoneInfo PhoneInfo `json:"phone_info"`
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiGetUserPhonenumber, param, result); err != nil {
		return nil, err
	}
	return &result.PhoneInfo, nil
}
