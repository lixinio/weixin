package wxopen

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGetTemplateDraftList = "/wxa/gettemplatedraftlist"
	apiAddToTemplate        = "/wxa/addtotemplate"
	apiGetTemplateList      = "/wxa/gettemplatelist"
	apiDeleteTemplate       = "/wxa/deletetemplate"
)

/*
获取代码草稿列表
通过本接口，可以获取草稿箱中所有的草稿（临时代码模板）；草稿是由第三方平台的开发小程序在使用微信开发者工具上传的
See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/code_template/gettemplatedraftlist.html
GET https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=ACCESS_TOKEN
*/
type Draft struct {
	DraftID     int32  `json:"draft_id"`
	UserVersion string `json:"user_version"`
	UserDesc    string `json:"user_desc"`
	CreateTime  int64  `json:"create_time"`
}

func (api *WxOpen) GetTemplateDraftList(ctx context.Context) ([]Draft, error) {
	result := struct {
		utils.WeixinError
		Drafts []Draft `json:"draft_list"`
	}{}
	err := api.Client.HTTPGet(ctx, apiGetTemplateDraftList, &result)
	if err != nil {
		return nil, err
	}
	return result.Drafts, nil
}

/*
将草稿添加到代码模板库
可以通过获取草稿箱中所有的草稿得到草稿 ID；调用本接口可以将临时草稿选为持久的代码模板
See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/code_template/addtotemplate.html
POST https://api.weixin.qq.com/wxa/addtotemplate?access_token=ACCESS_TOKEN
*/
func (api *WxOpen) AddToTemplate(ctx context.Context, draftID int32) error {
	return api.Client.HTTPPostJson(ctx, apiAddToTemplate, map[string]int32{
		"draft_id": draftID,
	}, nil)
}

/*
获取代码模板列表
第三方平台运营者可以登录 open.weixin.qq.com 或者通过将草稿箱的草稿选为代码模板接口，将草稿箱中的某个代码版本添加到代码模板库中
See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/code_template/gettemplatelist.html
GET https://api.weixin.qq.com/wxa/gettemplatelist?access_token=ACCESS_TOKEN
*/
type Template struct {
	TemplateID   int32  `json:"template_id"`
	TemplateType int8   `json:"template_type"`
	UserVersion  string `json:"user_version"`
	UserDesc     string `json:"user_desc"`
	CreateTime   int64  `json:"create_time"`
}

func (api *WxOpen) GetTemplateList(ctx context.Context) ([]Template, error) {
	result := struct {
		utils.WeixinError
		Templates []Template `json:"template_list"`
	}{}
	err := api.Client.HTTPGet(ctx, apiGetTemplateList, &result)
	if err != nil {
		return nil, err
	}
	return result.Templates, nil
}

/*
删除指定代码模板
因为代码模板库的模板数量是有上限的，当达到上限或者有某个模板不再需要时，可以调用本接口删除指定的代码模板。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/ThirdParty/code_template/deletetemplate.html
POST https://api.weixin.qq.com/wxa/deletetemplate?access_token=ACCESS_TOKEN
*/
func (api *WxOpen) DeleteTemplate(ctx context.Context, templateID int32) error {
	return api.Client.HTTPPostJson(ctx, apiDeleteTemplate, map[string]int32{
		"template_id": templateID,
	}, nil)
}
