package authorizer

import (
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl = "https://qyapi.weixin.qq.com" // 微信 api 服务器地址
)

var (
	ErrTokenUpdateForbidden = errors.New("can NOT refresh&update token in wxwork lite mode")
)

type Authorizer struct {
	SuiteID          string
	CorpID           string
	AgentID          int
	Client           *utils.Client
	accessTokenCache *utils.AccessTokenCache // 用于支持手动刷新Token， Client不对外暴露该对象
}

func New(
	cache utils.Cache,
	locker utils.Lock,
	suiteID, corpID string, agentID int,
	accessTokenGetter RefreshAccessToken,
) *Authorizer {
	accessTokenCache := utils.NewAccessTokenCache(
		newAdapter(suiteID, corpID, agentID, accessTokenGetter), cache, locker,
	)
	return &Authorizer{
		SuiteID:          suiteID,
		CorpID:           corpID,
		AgentID:          agentID,
		Client:           utils.NewClient(WXServerUrl, accessTokenCache),
		accessTokenCache: accessTokenCache,
	}
}

// 不支持刷新Token的简化模式
func NewLite(
	cache utils.Cache,
	locker utils.Lock,
	suiteID, corpID string, agentID int,
) *Authorizer {
	accessTokenCache := utils.NewAccessTokenCache(
		newAdapter(suiteID, corpID, agentID, func() (string, int, error) {
			return "", 0, fmt.Errorf(
				"can NOT refresh token in lite mod, appid(%s , %s , %d), %w",
				suiteID, corpID, agentID, ErrTokenUpdateForbidden,
			)
		}), cache, locker,
	)
	return &Authorizer{
		SuiteID: suiteID,
		CorpID:  corpID,
		AgentID: agentID,
		Client:  utils.NewClient(WXServerUrl, accessTokenCache),
	}
}

func (authorizer *Authorizer) RefreshAccessToken(expireBefore int) (string, error) {
	if authorizer.accessTokenCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrTokenUpdateForbidden,
		)
	}
	return authorizer.accessTokenCache.RefreshAccessToken(expireBefore)
}
