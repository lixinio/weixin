package wxwork_suite

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCode2Session = "/cgi-bin/service/miniprogram/jscode2session"
)

type MppSession struct {
	utils.WeixinError
	CorpID     string `json:"corpid"`
	UserID     string `json:"userid"`
	SessionKey string `json:"session_key"`
	OpenUserID string `json:"open_userid"`
}

// https://open.work.weixin.qq.com/api/doc/90001/90144/92427
// https://open.work.weixin.qq.com/api/doc/90001/90144/92423
func (suite *WxWorkSuite) Code2Session(ctx context.Context, jsCode string) (*MppSession, error) {
	session := &MppSession{}
	if err := suite.Client.HTTPGetWithParams(ctx, apiCode2Session, func(params url.Values) {
		params.Add("js_code", jsCode)
		params.Add("grant_type", "authorization_code")
	}, session); err != nil {
		return nil, err
	}

	return session, nil
}
