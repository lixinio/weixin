package agent

import (
	"fmt"

	"github.com/lixinio/weixin/utils"
	work "github.com/lixinio/weixin/wxwork"
)

type Config struct {
	AgentId string // 企业（自建）应用ID
	Secret  string // 企业（自建）应用密钥
}

type Agent struct {
	Config *Config
	wxwork *work.WxWork
	Client *utils.Client
}

func New(corp *work.WxWork, cache utils.Cache, config *Config) *Agent {
	instance := &Agent{
		Config: config,
		wxwork: corp,
	}
	instance.Client = corp.NewClient(utils.NewAccessTokenCache(instance, cache, 0))
	return instance
}

// GetAccessToken 接口 weixin.AccessTokenGetter 实现
func (agent *Agent) GetAccessToken() (accessToken string, expiresIn int, err error) {
	accessToken, expiresIn, err = agent.refreshAccessTokenFromWXServer()
	return
}

// GetAccessTokenKey 接口 weixin.AccessTokenGetter 实现
func (agent *Agent) GetAccessTokenKey() string {
	return fmt.Sprintf(
		"access-token:qywx-agent:%s:%s",
		agent.wxwork.Config.Corpid,
		agent.Config.AgentId,
	)
}
