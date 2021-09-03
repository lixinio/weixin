package agent

import (
	"context"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCode2Session = "/cgi-bin/miniprogram/jscode2session"
)

type MppSession struct {
	utils.WeixinError
	CorpID     string `json:"corpid"`
	UserID     string `json:"userid"`
	SessionKey string `json:"session_key"`
}

// https://work.weixin.qq.com/api/doc/90000/90136/91507
func (agent *Agent) Code2Session(ctx context.Context, jsCode string) (*MppSession, error) {
	session := &MppSession{}
	if err := agent.Client.HTTPGetWithParams(ctx, apiCode2Session, func(params url.Values) {
		params.Add("js_code", jsCode)
		params.Add("grant_type", "authorization_code")
	}, session); err != nil {
		return nil, err
	}

	return session, nil
}
