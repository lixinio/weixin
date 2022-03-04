package authorizer

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiBindTester   = "/wxa/bind_tester"
	apiMemberAuth   = "/wxa/memberauth"
	apiUnbindTester = "/wxa/unbind_tester"
)

type BindTesterParams struct {
	Wechatid string `json:"wechatid"`
}

/*
绑定微信用户为体验者
第三方平台在帮助旗下授权的小程序提交代码审核之前，可先让小程序运营者体验，体验之前需要将运营者的个人微信号添加到该小程序的体验者名单中。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_AdminManagement/Admin.html
POST https://api.weixin.qq.com/wxa/bind_tester?access_token=ACCESS_TOKEN
*/
func (api *Authorizer) BindTester(ctx context.Context, wechatid string) (string, error) {
	result := struct {
		utils.WeixinError
		Userstr string `json:"userstr"`
	}{}
	params := BindTesterParams{
		Wechatid: wechatid,
	}
	err := api.Client.HTTPPostJson(ctx, apiBindTester, params, &result)
	if err != nil {
		return "", err
	}
	return result.Userstr, nil
}

type MemberAuthParams struct {
	Action string `json:"action"`
}

type MemberAuthResult struct {
	utils.WeixinError
	Members []Member `json:"members"`
}

type Member struct {
	Userstr string `json:"userstr"`
}

/*
获取体验者列表
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_AdminManagement/Admin.html
POST https://api.weixin.qq.com/wxa/memberauth?access_token=TOKEN
*/
func (api *Authorizer) MemberAuth(ctx context.Context) (*MemberAuthResult, error) {
	params := MemberAuthParams{
		Action: "get_experiencer",
	}
	result := MemberAuthResult{}
	err := api.Client.HTTPPostJson(ctx, apiMemberAuth, params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
解除绑定体验者
调用本接口可以将特定微信用户从小程序的体验者列表中解绑。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_AdminManagement/unbind_tester.html
POST https://api.weixin.qq.com/wxa/unbind_tester?access_token=ACCESS_TOKEN
*/

type UnbindTesterParams struct {
	Wechatid string `json:"wechatid,omitempty"`
	Userstr  string `json:"userstr,omitempty"`
}

func (api *Authorizer) UnbindTester(ctx context.Context, wechatid string, userstr string) error {
	params := UnbindTesterParams{
		Wechatid: wechatid,
		Userstr:  userstr,
	}
	return api.Client.HTTPPostJson(ctx, apiUnbindTester, params, nil)
}
