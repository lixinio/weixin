package authorizer

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGetAccountBasicInfo   = "/cgi-bin/account/getaccountbasicinfo"
	apiCheckWxVerifyNickname = "/cgi-bin/wxverify/checkwxverifynickname"
	apiWxaSetNickname        = "/wxa/setnickname"
	apiWxaQueryNickName      = "/wxa/api_wxa_querynickname"
	apiModifyHeadImage       = "/cgi-bin/account/modifyheadimage"
	apiModifySignature       = "/cgi-bin/account/modifysignature"
)

// 获取基本信息
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/Mini_Program_Information_Settings.html
type MpAccountBasicInfo struct {
	utils.WeixinError
	AppID             string   `json:"appid"`
	AccountType       int      `json:"account_type"`    // 帐号类型（1：订阅号，2：服务号，3：小程序）
	PrincipalType     int      `json:"principal_type"`  // 主体类型 0	个人 1	企业 2	媒体 3	政府 4	其他组织
	PrincipalName     string   `json:"principal_name"`  // 主体名称
	Credential        string   `json:"credential"`      // 主体标识
	RealnameStatus    int      `json:"realname_status"` // 实名认证状态 1	实名验证成功 2	实名验证中 3	实名验证失败
	Nickname          string   `json:"nickname"`
	RegisteredCountry int      `json:"registered_country"` // 注册国家
	WxVerifyInfo      struct { // 微信认证信息
		QualificationVerify   bool  `json:"qualification_verify"`     // 是否资质认证，若是，拥有微信认证相关的权限。
		NamingVerify          bool  `json:"naming_verify"`            // 是否名称认证
		AnnualReview          bool  `json:"annual_review"`            // 是否需要年审（qualification_verify == true 时才有该字段）
		AnnualReviewBeginTime int64 `json:"annual_review_begin_time"` // 年审开始时间，时间戳（qualification_verify == true 时才有该字段）
		AnnualReviewEndTime   int64 `json:"annual_review_end_time"`   // 年审截止时间，时间戳（qualification_verify == true 时才有该字段）
	} `json:"wx_verify_info"`
	SignatureInfo struct { // 功能介绍信息
		Signature       string `json:"signature"`         // 功能介绍
		ModifyUsedCount int    `json:"modify_used_count"` // 功能介绍已使用修改次数（本月）
		ModifyQuota     int    `json:"modify_quota"`      // 功能介绍修改次数总额度（本月）
	} `json:"signature_info"`
	HeadImageInfo struct {
		HeadImageUrl    string `json:"head_image_url"`    // 头像 url
		ModifyUsedCount int    `json:"modify_used_count"` // 头像已使用修改次数（本年）
		ModifyQuota     int    `json:"modify_quota"`      // 头像修改次数总额度（本年）
	} `json:"head_image_info"`
	NicknameInfo struct { // 头像信息
		Nickname        string `json:"nickname"`          // 小程序名称
		ModifyUsedCount int    `json:"modify_used_count"` // 小程序名称已使用修改次数（本年）
		ModifyQuota     int    `json:"modify_quota"`      // 小程序名称修改次数总额度（本年）
	} `json:"nickname_info"`
}

func (api *Authorizer) GetAccountBasicInfo(ctx context.Context) (*MpAccountBasicInfo, error) {
	var result MpAccountBasicInfo
	if err := api.Client.HTTPPost(
		ctx, apiGetAccountBasicInfo, "{}", func(params url.Values) {}, &result, "",
	); err != nil {
		return nil, err
	}

	return &result, nil
}

// 微信认证名称检测
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/wxverify_checknickname.html
type WxVerifyNicknameResult struct {
	utils.WeixinError
	HitCondition bool   `json:"hit_condition"` // 是否命中关键字策略。若命中，可以选填关键字材料
	Wording      string `json:"wording"`       // 命中关键字的说明描述
}

func (api *Authorizer) CheckWxVerifyNickname(
	ctx context.Context, nickName string,
) (*WxVerifyNicknameResult, error) {
	var result WxVerifyNicknameResult
	params := map[string]string{
		"nick_name": nickName,
	}

	if err := api.Client.HTTPPostJson(ctx, apiCheckWxVerifyNickname, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// 设置名称
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/setnickname.html
type WxaSetNicknameParams struct {
	NickName           string `json:"nick_name"`                      // 昵称，不支持包含“小程序”关键字的昵称
	IdCard             string `json:"id_card,omitempty"`              // 身份证照片 mediaid
	License            string `json:"license,omitempty"`              // 组织机构代码证或营业执照 mediaid
	NamingOtherStuff_1 string `json:"naming_other_stuff_1,omitempty"` // 其他证明材料 mediaid
	NamingOtherStuff_2 string `json:"naming_other_stuff_2,omitempty"` // 其他证明材料 mediaid
	NamingOtherStuff_3 string `json:"naming_other_stuff_3,omitempty"` // 其他证明材料 mediaid
	NamingOtherStuff_4 string `json:"naming_other_stuff_4,omitempty"` // 其他证明材料 mediaid
	NamingOtherStuff_5 string `json:"naming_other_stuff_5,omitempty"` // 其他证明材料 mediaid
}

type WxaSetNicknameResult struct {
	utils.WeixinError
	Wording string `json:"wording"`  // 材料说明
	AuditID int    `json:"audit_id"` // 审核单 id，通过用于查询改名审核状态
}

func (api *Authorizer) WxaSetNickname(
	ctx context.Context,
	param *WxaSetNicknameParams,
) (*WxaSetNicknameResult, error) {
	var result WxaSetNicknameResult
	if err := api.Client.HTTPPostJson(ctx, apiWxaSetNickname, param, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// 查询改名审核状态
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/api_wxa_querynickname.html
type WxaQueryNicknameResult struct {
	utils.WeixinError
	Nickname   string `json:"nickname"`    // 审核昵称
	AuditStat  int    `json:"audit_stat"`  // 审核状态，1：审核中，2：审核失败，3：审核成功
	FailReason string `json:"fail_reason"` // 失败原因
	CreateTime int64  `json:"create_time"` // 审核提交时间
	AuditTime  int64  `json:"audit_time"`  // 审核完成时间
}

func (api *Authorizer) WxaQueryNickName(
	ctx context.Context, auditID int,
) (*WxaQueryNicknameResult, error) {
	var result WxaQueryNicknameResult
	params := map[string]int{
		"audit_id": auditID,
	}

	if err := api.Client.HTTPPostJson(ctx, apiWxaQueryNickName, params, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// 修改头像
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/modifyheadimage.html
type ModifyHeadImage struct {
	HeadImgMediaID string `json:"head_img_media_id"` // 头像素材 media_id
	X1             string `json:"x1"`                // 裁剪框左上角 x 坐标（取值范围：[0, 1]）
	Y1             string `json:"y1"`                // 裁剪框左上角 y 坐标（取值范围：[0, 1]）
	X2             string `json:"x2"`                // 裁剪框右下角 x 坐标（取值范围：[0, 1]）
	Y2             string `json:"y2"`                // 裁剪框右下角 y 坐标（取值范围：[0, 1]）
}

func (api *Authorizer) ModifyHeadImage(
	ctx context.Context, param *ModifyHeadImage,
) error {
	return api.Client.HTTPPostJson(ctx, apiModifyHeadImage, param, nil)
}

// 修改功能介绍
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/modifysignature.html
func (api *Authorizer) ModifySignature(
	ctx context.Context, signature string,
) error {
	params := map[string]string{
		"signature": signature,
	}
	return api.Client.HTTPPostJson(ctx, apiModifySignature, params, nil)
}
