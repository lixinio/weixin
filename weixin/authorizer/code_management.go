package authorizer

import (
	"context"
	"io"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCommit         = "/wxa/commit"
	apiGetQrcode      = "/wxa/get_qrcode"
	apiSubmitAudit    = "/wxa/submit_audit"
	apiRelease        = "/wxa/release"
	apiGetAuditStatus = "/wxa/get_auditstatus"
)

type AuditResult struct {
	Status     int32  `json:"status"`     // 审核状态
	Reason     string `json:"reason"`     // 当 status = 1 时，返回的拒绝原因; status = 4 时，返回的延后原因
	Screenshot string `json:"screenshot"` // 当 status = 1 时，会返回审核失败的小程序截图示例。用竖线分隔的 media_id 的列表，可通过获取永久素材接口拉取截图内容
}

// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/getAuditStatus.html
// 查询审核单状态
func (api *Authorizer) GetAuditStatus(
	ctx context.Context,
	auditid int32,
) (*AuditResult, error) {
	result := &struct {
		utils.WeixinError
		AuditResult
	}{}

	if err := api.Client.HTTPPostJson(ctx, apiGetAuditStatus, map[string]int32{
		"auditid": auditid,
	}, result); err != nil {
		return nil, err
	}

	return &result.AuditResult, nil
}

/*
上传代码
第三方平台需要先将草稿添加到代码模板库，或者从代码模板库中选取某个代码模板，得到对应的模板 id（template_id）； 然后调用本接口可以为已授权的小程序上传代码。
注意，通过该接口提交代码后，会同时生成体验版。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/commit.html
POST https://api.weixin.qq.com/wxa/commit?access_token=ACCESS_TOKEN
*/
func (api *Authorizer) CodeCommit(
	ctx context.Context,
	templateID int32,
	extJson string,
	userVersion string,
	userDesc string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCommit, map[string]interface{}{
		"template_id":  templateID,
		"ext_json":     extJson,
		"user_version": userVersion,
		"user_desc":    userDesc,
	}, nil)
}

/*
获取体验版二维码
调用本接口可以获取小程序的体验版二维码
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/get_qrcode.html
GET https://api.weixin.qq.com/wxa/get_qrcode?access_token=ACCESS_TOKEN&path=page%2Findex%3Faction%3D1
*/
func (api *Authorizer) GetTestQrcode(ctx context.Context, path string) ([]byte, error) {
	resp, err := api.Client.HTTPGetRaw(ctx, apiGetQrcode, func(params url.Values) {
		params.Add("path", path)
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type AuditParamsItem struct {
	Address     string `json:"address,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Title       string `json:"title,omitempty"`
	FirstID     int    `json:"first_id,omitempty"`
	SecondID    int    `json:"second_id,omitempty"`
	ThirdID     int    `json:"third_id,omitempty"`
	FirstClass  string `json:"first_class,omitempty"`
	SecondClass string `json:"second_class,omitempty"`
	ThirdClass  string `json:"third_class,omitempty"`
}
type PreviewInfo struct {
	VideoIdList []string `json:"video_id_list,omitempty"`
	PicIdList   []string `json:"pic_id_list,omitempty"`
}
type UGCDeclare struct {
	Scene          []int  `json:"scene,omitempty"`
	OtherSceneDesc string `json:"other_scene_desc,omitempty"`
	Method         []int  `json:"method,omitempty"`
	HasAuditTeam   int    `json:"has_audit_team,omitempty"`
	AuditDesc      string `json:"audit_desc,omitempty"`
}
type AuditParams struct {
	ItemList         []AuditParamsItem `json:"item_list,omitempty"`
	PreviewInfo      PreviewInfo       `json:"preview_info,omitempty"`
	VersionDesc      string            `json:"version_desc,omitempty"`
	FeedbackInfo     string            `json:"feedback_info,omitempty"`
	FeedbackStuff    string            `json:"feedback_stuff,omitempty"`
	PrivacyApiNotUse bool              `json:"privacy_api_not_use,omitempty"`
	UGCDecalre       UGCDeclare        `json:"ugc_declare,omitempty"`
}

/*
提交审核
在调用上传代码接口为小程序上传代码后，可以调用本接口，将上传的代码提交审核。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/submit_audit.html
https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/code-management/submitAudit.html
POST https://api.weixin.qq.com/wxa/submit_audit?access_token=ACCESS_TOKEN
*/
func (api *Authorizer) CodeSubmitAudit(
	ctx context.Context,
	auditParams *AuditParams,
) (int32, error) {
	result := struct {
		utils.WeixinError
		AuditID int32 `json:"auditid"`
	}{}
	err := api.Client.HTTPPostJson(ctx, apiSubmitAudit, auditParams, &result)
	if err != nil {
		return -1, err
	}
	return result.AuditID, nil
}

/*
发布已通过审核的小程序
调用本接口可以发布最后一个审核通过的小程序代码版本。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/release.html
POST https://api.weixin.qq.com/wxa/release?access_token=ACCESS_TOKEN
*/
func (api *Authorizer) CodeRelease(ctx context.Context) error {
	return api.Client.HTTPPostJson(ctx, apiRelease, struct{}{}, nil)
}
