package agent

import (
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
	instance.Client = corp.NewClient(utils.NewAccessTokenCache(instance, cache, locker, 0))
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
		"qywx_%s_%d.access_token",
		agent.wxwork.Config.Corpid,
		agent.Config.AgentID,
	)
}

// GetAccessTokenLockKey 接口 weixin.AccessTokenGetter 实现
func (agent *Agent) GetAccessTokenLockKey() string {
	return fmt.Sprintf(
		"qywx_%s_%d.access_token.lock",
		agent.wxwork.Config.Corpid,
		agent.Config.AgentID,
	)
}
