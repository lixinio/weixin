package authorizer

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

type DomainAction string

const (
	DomainActionAdd    DomainAction = "add"
	DomainActionDelete DomainAction = "delete"
	DomainActionSet    DomainAction = "set"
	DomainActionGet    DomainAction = "get"
)

const (
	apiModifyDomain                = "/wxa/modify_domain"
	apiSetWebViewDomain            = "/wxa/setwebviewdomain"
	apiModifyDomainDirectly        = "/wxa/modify_domain_directly"        // 快速配置小程序服务器域名
	apiSetWebViewDomainDirectly    = "/wxa/setwebviewdomain_directly"     // 快速配置小程序业务域名
	apiGetEffectiveDomain          = "/wxa/get_effective_domain"          // 获取发布后生效服务器域名列表
	apiGetEffectiveWebviewDomain   = "/wxa/get_effective_webviewdomain"   // 获取发布后生效业务域名列表
	apiGetWebviewDomainConfirmFile = "/wxa/get_webviewdomain_confirmfile" // 获取业务域名校验文件
)

type ServerDomain struct {
	RequestDomain   []string `json:"requestdomain"`
	WsrequestDomain []string `json:"wsrequestdomain"`
	UploadDomain    []string `json:"uploaddomain"`
	DownloadDomain  []string `json:"downloaddomain"`
}

/*
设置服务器域名
授权给第三方的小程序，其服务器域名只可以为在第三方平台账号中配置的小程序服务器域名，
当小程序通过第三方平台发布代码上线后，小程序原先自己配置的服务器域名将被删除，只保留第三方平台的域名，
所以第三方平台在代替小程序发布代码之前，需要调用接口为小程序添加第三方平台自身的域名。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/Server_Address_Configuration.html
POST https://api.weixin.qq.com/wxa/modify_domain?access_token=ACCESS_TOKEN
*/
func (api *Authorizer) ModifyDomain(
	ctx context.Context,
	action DomainAction,
	domains *ServerDomain,
) (*ServerDomain, error) {
	params := struct {
		Action DomainAction `json:"action"` // 可选值 add delete set get
		*ServerDomain
	}{
		Action:       action,
		ServerDomain: domains,
	}

	result := struct {
		utils.WeixinError
		ServerDomain
	}{}

	err := api.Client.HTTPPostJson(ctx, apiModifyDomain, &params, &result)
	if err != nil {
		return nil, err
	}
	return &result.ServerDomain, nil
}

/*
快速配置小程序服务器域名
https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/domain-management/modifyServerDomainDirectly.html
*/
func (api *Authorizer) ModifyDomainDirectly(
	ctx context.Context,
	action DomainAction,
	domains *ServerDomain,
) (*ServerDomain, error) {
	params := struct {
		Action DomainAction `json:"action"` // 可选值 add delete set get
		*ServerDomain
	}{
		Action:       action,
		ServerDomain: domains,
	}

	result := struct {
		utils.WeixinError
		ServerDomain
	}{}

	err := api.Client.HTTPPostJson(ctx, apiModifyDomainDirectly, &params, &result)
	if err != nil {
		return nil, err
	}
	return &result.ServerDomain, nil
}

type SetWebViewDomainParams struct {
	Action        DomainAction `json:"action,omitempty"` // 可选值 add delete set get
	WebviewDomain []string     `json:"webviewdomain,omitempty"`
}

/*
设置业务域名
授权给第三方的小程序，其业务域名只可以为在第三方平台账号中配置的小程序业务域名，当小程序通过第三方发布代码上线后，小程序原先自己配置的业务域名将被删除，只保留第三方平台的域名，所以第三方平台在代替小程序发布代码之前，需要调用接口为小程序添加业务域名。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/setwebviewdomain.html
POST https://api.weixin.qq.com/wxa/setwebviewdomain?access_token=ACCESS_TOKEN
*/
func (api *Authorizer) SetWebViewDomain(
	ctx context.Context,
	params *SetWebViewDomainParams,
) ([]string, error) {
	result := struct {
		utils.WeixinError
		Domains []string `json:"webviewdomain"`
	}{}

	err := api.Client.HTTPPostJson(ctx, apiSetWebViewDomain, params, &result)
	if err != nil {
		return nil, err
	}

	return result.Domains, nil
}

/*
快速配置小程序业务域名
https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/domain-management/modifyJumpDomainDirectly.html
*/
func (api *Authorizer) SetWebViewDomainDirectly(
	ctx context.Context,
	params *SetWebViewDomainParams,
) ([]string, error) {
	result := struct {
		utils.WeixinError
		Domains []string `json:"webviewdomain"`
	}{}

	err := api.Client.HTTPPostJson(ctx, apiSetWebViewDomainDirectly, params, &result)
	if err != nil {
		return nil, err
	}

	return result.Domains, nil
}

type GetEffectiveDomainResp struct {
	utils.WeixinError
	MpDomain        *ServerDomain `json:"mp_domain"`        // 通过公众平台配置的服务器域名列表
	ThirdDmain      *ServerDomain `json:"third_domain"`     // 通过第三方平台接口modify_domain 配置的服务器域名
	DirectDomain    *ServerDomain `json:"direct_domain"`    // 通过“modify_domain_directly”接口配置的服务器域名列表
	EffectiveDomain *ServerDomain `json:"effective_domain"` // 最后提交代码或者发布上线后生效的域名列表
}

/*
获取发布后生效服务器域名列表
https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/domain-management/getEffectiveServerDomain.html
*/
func (api *Authorizer) GetEffectiveDomain(
	ctx context.Context,
) (*GetEffectiveDomainResp, error) {
	result := &GetEffectiveDomainResp{}

	err := api.Client.HTTPPostJson(ctx, apiGetEffectiveDomain, nil, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type GetEffectiveDomainWebviewResp struct {
	utils.WeixinError
	MpDomain        []string `json:"mp_webviewdomain"`        // 通过公众平台配置的服务器域名列表
	ThirdDmain      []string `json:"third_webviewdomain"`     // 通过第三方平台接口modifyJumpDomain 配置的业务域名
	DirectDomain    []string `json:"direct_webviewdomain"`    // 通过modifyJumpDomainDirectly接口配置的业务域名列表
	EffectiveDomain []string `json:"effective_webviewdomain"` // 最后提交代码或者发布上线后生效的域名列表
}

/*
获取发布后生效业务域名列表
https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/domain-management/getEffectiveJumpDomain.html
*/
func (api *Authorizer) GetEffectiveWebviewDomain(
	ctx context.Context,
) (*GetEffectiveDomainWebviewResp, error) {
	result := &GetEffectiveDomainWebviewResp{}

	err := api.Client.HTTPPostJson(ctx, apiGetEffectiveWebviewDomain, nil, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (api *Authorizer) GetWebviewDomainConfirmFile(
	ctx context.Context,
) (string, string, error) {
	result := struct {
		utils.WeixinError
		FileName    string `json:"file_name"`
		FileContent string `json:"file_content"`
	}{}

	err := api.Client.HTTPPostJson(ctx, apiGetWebviewDomainConfirmFile, nil, &result)
	if err != nil {
		return "", "", err
	}

	return result.FileName, result.FileContent, nil
}
