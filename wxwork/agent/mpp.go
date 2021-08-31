package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	apiCode2Session = "/cgi-bin/miniprogram/jscode2session"
)

type MppSession struct {
	CorpID     string `json:"corpid"`
	UserID     string `json:"userid"`
	SessionKey string `json:"session_key"`
}

// https://work.weixin.qq.com/api/doc/90000/90136/91507
func (agent *Agent) Code2Session(ctx context.Context, jsCode string) (*MppSession, error) {
	params := url.Values{}
	params.Add("js_code", jsCode)
	params.Add("grant_type", "authorization_code")

	body, err := agent.Client.HTTPGetWithParams(ctx, apiCode2Session, params)
	if err != nil {
		return nil, err
	}

	session := &MppSession{}
	err = json.Unmarshal(body, session)
	if err != nil {
		return nil, fmt.Errorf("%s", string(body))
	}

	return session, nil
}
