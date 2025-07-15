package agent

import (
	"context"
	"fmt"

	"github.com/lixinio/weixin/utils"
	work "github.com/lixinio/weixin/wxwork"
)

type Config struct {
	AgentID int    // 企业（自建）应用ID
	Secret  string // 企业（自建）应用密钥
}

type Agent struct {
	Config *Config
	wxwork *work.WxWork
	Client *utils.Client
}

func New(corp *work.WxWork, cache utils.Cache, locker utils.Lock, config *Config) *Agent {
	instance := &Agent{
		Config: config,
		wxwork: corp,
	}
	instance.Client = corp.NewClient(utils.NewAccessTokenCache(
		newAdapter(
			corp.Config.Corpid,
			config.AgentID,
			instance.refreshAccessTokenFromWXServer,
		),
		cache, locker,
	))
	return instance
}

func NewLite(
	corp *work.WxWork, cache utils.Cache, locker utils.Lock, agentID int,
) *Agent {
	client := corp.NewClient(
		utils.NewAccessTokenCache(
			newAdapter(
				corp.Config.Corpid,
				agentID,
				func(ctx context.Context) (string, int, error) {
					return "", 0, fmt.Errorf(
						"can NOT refresh token in lite mod, corp(%s), agentid(%d), %w",
						corp.Config.Corpid, agentID, ErrTokenUpdateForbidden,
					)
				}),
			cache, locker,
		),
	)
	return &Agent{
		Config: &Config{AgentID: agentID},
		wxwork: corp,
		Client: client,
	}
}

func (agent *Agent) CorpID() string {
	return agent.wxwork.Config.Corpid
}
