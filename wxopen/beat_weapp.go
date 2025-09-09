package wxopen

// 试用小程序
import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiFastRegisterCorpWeapp     = "/cgi-bin/component/fastregisterweapp"     // 快速注册企业小程序
	apiFastRegisterPersonalWeapp = "/wxa/component/fastregisterpersonalweapp" // 快速注册个人小程序
	apiFastRegisterBetaWeapp     = "/wxa/component/fastregisterbetaweapp"     // 注册试用小程序
	apiVerifyBetaWeapp           = "/wxa/verifybetaweapp"                     // 试用小程序快速转正
	apiSetBetaWeappNickname      = "/wxa/setbetaweappnickname"                // 修改试用小程序名称
)

// 创建试用小程序
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/beta_Mini_Programs/fastregister.html
type FastRegisterBetaWeappResult struct {
	utils.WeixinError
	UniqueID     string `json:"unique_id"`     // 该请求的唯一标识符，用于关联微信用户和后面产生的appid
	AuthorizeUrl string `json:"authorize_url"` // 用户授权确认url，需将该url发送给用户，用户进入授权页面完成授权方可创建小程序
}

func (api *WxOpen) FastRegisterBetaWeapp(
	ctx context.Context, name, openid string,
) (*FastRegisterBetaWeappResult, error) {
	payload := map[string]string{
		"name":   name,
		"openid": openid,
	}

	result := FastRegisterBetaWeappResult{}
	err := api.Client.HTTPPostJson(ctx, apiFastRegisterBetaWeapp, payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

/**
<xml>
    <AppId><![CDATA[第三方平台appid]]></AppId>
    <CreateTime>1535442403</CreateTime>
    <InfoType><![CDATA[notify_third_fastregisterbetaapp]]></InfoType>
    <appid>创建小程序appid<appid>
    <status>0</status>
    <msg>OK</msg>
    <info>
    <unique_id><![CDATA[unique_id]]></unique_id>
    <name><![CDATA[小程序名称]]></name>
    </info>
</xml>
**/

// 试用小程序快速认证
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/beta_Mini_Programs/fastverify.html
type BetaWeappVerifyInfo struct {
	EnterpriseName     string `json:"enterprise_name"`      // 企业名（需与工商部门登记信息一致)；如果是“无主体名称个体工商户”则填“个体户+法人姓名”，例如“个体户张三”
	Code               string `json:"code"`                 // 企业代码
	CodeType           string `json:"code_type"`            // 企业代码类型 1：统一社会信用代码（18 位） 2：组织机构代码（9 位 xxxxxxxx-x） 3：营业执照注册号(15 位)
	LegalPersonaWechat string `json:"legal_persona_wechat"` // 法人微信号
	LegalPersonaName   string `json:"legal_persona_name"`   // 法人姓名（绑定银行卡）
	LegalPersonaIdcard string `json:"legal_persona_idcard"` // 法人身份证号
	ComponentPhone     string `json:"component_phone"`      // 第三方联系电话
}

func (api *WxOpen) VerifyBetaWeapp(
	ctx context.Context, info *BetaWeappVerifyInfo,
) error {
	payload := map[string]*BetaWeappVerifyInfo{
		"verify_info": info,
	}
	return api.Client.HTTPPostJson(ctx, apiVerifyBetaWeapp, payload, nil)
}

// 修改试用小程序名称
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/beta_Mini_Programs/fastmodify.html
// 该接口仅适用于试用小程序，不适用于已经完成认证的普通小程序，即转正之后不可以调用。
// 待小程序转正之后，需要服务商调用setnickname的接口重新设置名称，因为认证后不会自动去除“的试用小程序后缀”。
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/setnickname.html
// 如果需要修改试用小程序的头像，请调用【设置头像】进行修改，且不占用正式小程序修改头像的次数quota。
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/modifyheadimage.html
func (api *WxOpen) SetBetaWeappNickname(
	ctx context.Context, name string,
) error {
	payload := map[string]string{
		"name": name,
	}
	return api.Client.HTTPPostJson(ctx, apiSetBetaWeappNickname, payload, nil)
}

// 快速注册个人小程序
// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/register-management/fast-registration-ind/fastRegisterPersonalMp.html
type FastRegisterPersonalWeappResult struct {
	utils.WeixinError
	TaskID       string `json:"taskid"`        // 任务id，后面query查询需要用到
	AuthorizeUrl string `json:"authorize_url"` // 用户授权确认url，需将该url发送给用户，用户进入授权页面完成授权方可创建小程序
	Status       int64  `json:"status"`        // 任务的状态
}

func (api *WxOpen) CreateFastRegisterPersonalWeapp(
	ctx context.Context, idname, wxuser, componentPhone string,
) (*FastRegisterPersonalWeappResult, error) {
	payload := map[string]string{
		"idname":          idname,
		"wxuser":          wxuser,
		"component_phone": componentPhone,
	}

	result := FastRegisterPersonalWeappResult{}
	err := api.Client.HTTPPost(
		ctx, apiFastRegisterPersonalWeapp, payload,
		func(u url.Values) {
			u.Add("action", "create")
		},
		&result, "",
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (api *WxOpen) QueryFastRegisterPersonalWeapp(
	ctx context.Context, taskid string,
) (*FastRegisterPersonalWeappResult, error) {
	payload := map[string]string{
		"taskid": taskid,
	}

	result := FastRegisterPersonalWeappResult{}
	err := api.Client.HTTPPost(
		ctx, apiFastRegisterPersonalWeapp, payload,
		func(u url.Values) {
			u.Add("action", "query")
		},
		&result, "",
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/register-management/fast-registration-ent/registerMiniprogram.html
// 快速注册企业小程序
type FastRegisterPersonalWeappRequest struct {
	utils.WeixinError
	Name               string `json:"name"`                 // 企业名（需与工商部门登记信息一致）；如果是“无主体名称个体工商户”则填“个体户+法人姓名”，例如“个体户张三”
	Code               string `json:"code"`                 // 企业代码
	CodeType           int64  `json:"code_type"`            // 企业代码类型 1：统一社会信用代码（18 位） 2：组织机构代码（9 位 xxxxxxxx-x） 3：营业执照注册号(15 位)
	LegalPersonaWechat string `json:"legal_persona_wechat"` // 法人微信号
	LegalPersonaName   string `json:"legal_persona_name"`   // 法人姓名（绑定银行卡）
	ComponentPhone     string `json:"component_phone"`      // 第三方联系电话
}

func (api *WxOpen) CreateFastRegisterCorpWeapp(
	ctx context.Context, req *FastRegisterPersonalWeappRequest,
) error {
	err := api.Client.HTTPPost(
		ctx, apiFastRegisterPersonalWeapp, req,
		func(u url.Values) {
			u.Add("action", "create")
		}, nil, "",
	)
	if err != nil {
		return err
	}

	return nil
}

/*
"name": "tencent", // 企业名
"legal_persona_wechat": "123", // 法人微信
"legal_persona_name": "pony" // 法人姓名
*/
func (api *WxOpen) QueryFastRegisterCorpWeapp(
	ctx context.Context, name, legalPersonaWechat, legalPersonaName string,
) (*FastRegisterPersonalWeappResult, error) {
	payload := map[string]string{
		"name":                 name,
		"legal_persona_wechat": legalPersonaWechat,
		"legal_persona_name":   legalPersonaName,
	}

	result := FastRegisterPersonalWeappResult{}
	err := api.Client.HTTPPost(
		ctx, apiFastRegisterPersonalWeapp, payload,
		func(u url.Values) {
			u.Add("action", "search")
		},
		&result, "",
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
