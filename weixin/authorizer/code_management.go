package authorizer

import (
	"context"
	"io/ioutil"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCommit      = "/wxa/commit"
	apiGetQrcode   = "/wxa/get_qrcode"
	apiSubmitAudit = "/wxa/submit_audit"
	apiRelease     = "/wxa/release"
)

type AuthorizerApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *AuthorizerApi {
	return &AuthorizerApi{Client: client}
}

/*
上传代码
第三方平台需要先将草稿添加到代码模板库，或者从代码模板库中选取某个代码模板，得到对应的模板 id（template_id）； 然后调用本接口可以为已授权的小程序上传代码。
注意，通过该接口提交代码后，会同时生成体验版。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/commit.html
POST https://api.weixin.qq.com/wxa/commit?access_token=ACCESS_TOKEN
*/
func (api *AuthorizerApi) Commit(ctx context.Context, templateID int32, extJson string, userVersion string, userDesc string) error {
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
func (api *AuthorizerApi) GetQrcode(ctx context.Context, path string) ([]byte, error) {
	resp, err := api.Client.HTTPGetRaw(ctx, apiGetQrcode, func(params url.Values) {
		params.Add("path", path)
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/*
提交审核
在调用上传代码接口为小程序上传代码后，可以调用本接口，将上传的代码提交审核。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/code/submit_audit.html
POST https://api.weixin.qq.com/wxa/submit_audit?access_token=ACCESS_TOKEN
*/
func (api *AuthorizerApi) SubmitAudit(ctx context.Context, auditParams map[string]interface{}) (int32, error) {
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
func (api *AuthorizerApi) Release(ctx context.Context) error {
	return api.Client.HTTPPostJson(ctx, apiRelease, struct{}{}, nil)
}
