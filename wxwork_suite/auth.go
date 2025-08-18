package wxwork_suite

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiInstall            = "https://open.work.weixin.qq.com/3rdapp/install"
	apiGetPreAuthCode     = "/cgi-bin/service/get_pre_auth_code"
	apiSetSessionInfo     = "/cgi-bin/service/set_session_info"
	apiGetPermanentCode   = "/cgi-bin/service/get_permanent_code"
	apiGetAuthInfo        = "/cgi-bin/service/get_auth_info"
	apiGetPermanentCodeV2 = "/cgi-bin/service/v2/get_permanent_code"
	apiGetAuthInfoV2      = "/cgi-bin/service/v2/get_auth_info"
	apiGetCorpToken       = "/cgi-bin/service/get_corp_token"
	apiGetAdminList       = "/cgi-bin/service/get_admin_list"
)

type PreAuthCode struct {
	utils.WeixinError
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

// 从服务商网站发起
// https://open.work.weixin.qq.com/api/doc/90001/90143/90597
func (suite *WxWorkSuite) GetInstallUrl(
	redirectUri, preAuthCode, state string,
) (authorizeUrl string) {
	params := url.Values{}
	params.Add("suite_id", suite.Config.SuiteID)
	params.Add("redirect_uri", redirectUri)
	params.Add("pre_auth_code", preAuthCode)
	params.Add("state", state)
	return apiInstall + "?" + params.Encode()
}

// 获取预授权码
// https://open.work.weixin.qq.com/api/doc/90001/90143/90601
func (suite *WxWorkSuite) GetPreAuthCode(ctx context.Context) (*PreAuthCode, error) {
	result := &PreAuthCode{}
	if err := suite.Client.HTTPGet(ctx, apiGetPreAuthCode, result); err != nil {
		return nil, err
	}
	return result, nil
}

type SessionInfo struct {
	PreAuthCode string `json:"pre_auth_code"`
	SessionInfo struct {
		AppID    []int `json:"appid,omitempty"`
		AuthType int   `json:"auth_type"`
	} `json:"session_info"`
}

// 设置授权配置
// 授权类型：0 正式授权， 1 测试授权。 默认值为0。注意，请确保应用在正式发布后的授权类型为“正式授权”
// https://open.work.weixin.qq.com/api/doc/90001/90143/90602
func (suite *WxWorkSuite) SetSessionInfo(
	ctx context.Context,
	preAuthCode string,
	authType int,
) error {
	return suite.Client.HTTPPostJson(ctx, apiSetSessionInfo, &SessionInfo{
		PreAuthCode: preAuthCode,
		SessionInfo: struct {
			AppID    []int `json:"appid,omitempty"`
			AuthType int   `json:"auth_type"`
		}{
			AuthType: authType,
		},
	}, nil)
}

// v2 如果需要获取 CorpWxQrcode ， 单独调用接口
// get_app_qrcode
// https://developer.work.weixin.qq.com/document/path/95430
type AuthCorpInfo struct {
	CorpID            string `json:"corpid"`
	CorpName          string `json:"corp_name"`
	CorpFullName      string `json:"corp_full_name"`
	CorpType          string `json:"corp_type"`
	CorpWxQrcode      string `json:"corp_wxqrcode"` // v2 不支持
	CorpScale         string `json:"corp_scale"`
	CorpIndustry      string `json:"corp_industry"`
	CorpSubIndustry   string `json:"corp_sub_industry"`
	CorpSquareLogoUrl string `json:"corp_square_logo_url"`
	CorpUserMax       int    `json:"corp_user_max"`
	CorpAgentMax      int    `json:"corp_agent_max"`
	VerifiedEndEime   int    `json:"verified_end_time"`
	SubjectType       int    `json:"subject_type"`
}

type AuthUserInfo struct {
	UserID     string `json:"userid"`
	OpenUserID string `json:"open_userid"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar"`
}

type RegisterCodeInfo struct {
	RegisterCode string `json:"register_code"`
	TemplateID   string `json:"template_id"`
	State        string `json:"state"`
}

type DealerCorpInfo struct {
	CorpID   string `json:"corpid"`
	CorpName string `json:"corp_name"`
}

type AgentInfo struct {
	AgentID         int    `json:"agentid"`
	Name            string `json:"name"`
	RoundLogoUrl    string `json:"round_logo_url"`
	SquareLogoUrl   string `json:"square_logo_url"`
	AppID           int    `json:"appid"`
	AuthMode        int    `json:"auth_mode"`
	IsCustomizedApp bool   `json:"is_customized_app"`
	SharedFrom      struct {
		CorpID string `json:"corpid"`
	} `json:"shared_from"`
}

type PermanentCodeInfo struct {
	utils.WeixinError
	AccessToken      string            `json:"access_token"`
	ExpiresIn        int               `json:"expires_in"`
	PermanentCode    string            `json:"permanent_code"`
	DealerCorpInfo   *DealerCorpInfo   `json:"dealer_corp_info"`
	AuthCorpInfo     *AuthCorpInfo     `json:"auth_corp_info"`
	AuthUserInfo     *AuthUserInfo     `json:"auth_user_info"`
	RegisterCodeInfo *RegisterCodeInfo `json:"register_code_info"`
	AuthInfo         struct {
		Agents []AgentInfo `json:"agent"`
	} `json:"auth_info"`
}

// 获取企业永久授权码
// https://open.work.weixin.qq.com/api/doc/90001/90143/90603
func (suite *WxWorkSuite) GetPermanentCode(
	ctx context.Context, authCode string,
) (*PermanentCodeInfo, error) {
	result := &PermanentCodeInfo{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetPermanentCode, map[string]string{
		"auth_code": authCode,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

type PermanentCodeInfoV2 struct {
	utils.WeixinError
	PermanentCode string `json:"permanent_code"`
	State         string `json:"state"`
	AuthCorpInfo  *struct {
		CorpID   string `json:"corpid"`
		CorpName string `json:"corp_name"`
	} `json:"auth_corp_info"`
	AuthUserInfo     *AuthUserInfo     `json:"auth_user_info"`
	RegisterCodeInfo *RegisterCodeInfo `json:"register_code_info"`
}

// 获取企业永久授权码
// https://developer.work.weixin.qq.com/document/path/100776
func (suite *WxWorkSuite) GetPermanentCodeV2(
	ctx context.Context, authCode string,
) (*PermanentCodeInfoV2, error) {
	result := &PermanentCodeInfoV2{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetPermanentCodeV2, map[string]string{
		"auth_code": authCode,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

type AuthInfo struct {
	utils.WeixinError
	DealerCorpInfo *DealerCorpInfo `json:"dealer_corp_info"`
	AuthCorpInfo   *AuthCorpInfo   `json:"auth_corp_info"`
	AuthInfo       struct {
		Agents []AgentInfo `json:"agent"`
	} `json:"auth_info"`
}

// 获取企业授权信息
// https://open.work.weixin.qq.com/api/doc/90001/90143/90604
func (suite *WxWorkSuite) GetAuthInfo(
	ctx context.Context, authCorpID, permanentCode string,
) (*AuthInfo, error) {
	result := &AuthInfo{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetAuthInfo, map[string]string{
		"auth_corpid":    authCorpID,
		"permanent_code": permanentCode,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取企业授权信息
// https://developer.work.weixin.qq.com/document/path/100779
func (suite *WxWorkSuite) GetAuthInfoV2(
	ctx context.Context, authCorpID, permanentCode string,
) (*AuthInfo, error) {
	result := &AuthInfo{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetAuthInfoV2, map[string]string{
		"auth_corpid":    authCorpID,
		"permanent_code": permanentCode,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

// 获取企业凭证
// https://open.work.weixin.qq.com/api/doc/90001/90143/90605
func (suite *WxWorkSuite) GetCorpToken(
	ctx context.Context, authCorpID, permanentCode string,
) (*utils.TokenResponse, error) {
	result := &utils.TokenResponse{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetCorpToken, map[string]string{
		"auth_corpid":    authCorpID,
		"permanent_code": permanentCode,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}

type AgentAdminUser struct {
	UserID     string `json:"userid"`
	OpenUserID string `json:"open_userid"`
	AuthType   int    `json:"auth_type"`
}

type AgentAdmin struct {
	utils.WeixinError
	Admin []AgentAdminUser `json:"admin"`
}

// 获取应用的管理员列表
// https://open.work.weixin.qq.com/api/doc/90001/90143/90606
func (suite *WxWorkSuite) GetAdminList(
	ctx context.Context, authCorpID string, agentid int,
) (*AgentAdmin, error) {
	result := &AgentAdmin{}
	if err := suite.Client.HTTPPostJson(ctx, apiGetAdminList, map[string]interface{}{
		"auth_corpid": authCorpID,
		"agentid":     agentid,
	}, result); err != nil {
		return nil, err
	}
	return result, nil
}
