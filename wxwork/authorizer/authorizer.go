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
	ErrTokenUpdateForbidden      = errors.New("can NOT refresh&update token in wxwork lite mode")
	ErrCorpJsApiTicketForbidden  = errors.New("can NOT refresh&update corp jsapi ticket without enable it")
	ErrAgentJsApiTicketForbidden = errors.New("can NOT refresh&update agent jsapi ticket without enable it")
)

type Authorizer struct {
	SuiteID               string
	CorpID                string
	AgentID               int
	Client                *utils.Client
	accessTokenCache      *utils.AccessTokenCache // 用于支持手动刷新Token， Client不对外暴露该对象
	corpJsApiTicketCache  *utils.AccessTokenCache
	agentJsApiTicketCache *utils.AccessTokenCache
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

func (authorizer *Authorizer) ClearAccessToken() error {
	if authorizer.accessTokenCache == nil {
		return fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrTokenUpdateForbidden,
		)
	}
	return authorizer.accessTokenCache.ClearAccessToken()
}

func (authorizer *Authorizer) RefreshCorpJsApiTicket(expireBefore int) (string, error) {
	if authorizer.corpJsApiTicketCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrCorpJsApiTicketForbidden,
		)
	}
	return authorizer.corpJsApiTicketCache.RefreshAccessToken(expireBefore)
}

func (authorizer *Authorizer) ClearCorpJsApiTicket() error {
	if authorizer.corpJsApiTicketCache == nil {
		return fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrCorpJsApiTicketForbidden,
		)
	}
	return authorizer.corpJsApiTicketCache.ClearAccessToken()
}

func (authorizer *Authorizer) RefreshAgentJsApiTicket(expireBefore int) (string, error) {
	if authorizer.agentJsApiTicketCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrAgentJsApiTicketForbidden,
		)
	}
	return authorizer.agentJsApiTicketCache.RefreshAccessToken(expireBefore)
}

func (authorizer *Authorizer) ClearAgentJsApiTicket() error {
	if authorizer.agentJsApiTicketCache == nil {
		return fmt.Errorf(
			"authorizer appid : %s,%s,%d, error: %w",
			authorizer.SuiteID, authorizer.CorpID, authorizer.AgentID,
			ErrAgentJsApiTicketForbidden,
		)
	}
	return authorizer.agentJsApiTicketCache.ClearAccessToken()
}
