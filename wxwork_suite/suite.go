package wxwork_suite

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl      = "https://qyapi.weixin.qq.com" // 微信 api 服务器地址
	accessTokenKey   = "suite_access_token"
	apiGetSuiteToken = "/cgi-bin/service/get_suite_token"
)

// https://open.work.weixin.qq.com/api/doc/90001/90143/90600
/*
由于第三方服务商可能托管了大量的企业，其安全问题造成的影响会更加严重，故API中除了合法来源IP校验之外，还额外增加了suite_ticket作为安全凭证。
获取suite_access_token时，需要suite_ticket参数。suite_ticket由企业微信后台定时推送给“指令回调URL”，每十分钟更新一次，见推送suite_ticket。
suite_ticket实际有效期为30分钟，可以容错连续两次获取suite_ticket失败的情况，但是请永远使用最新接收到的suite_ticket。
*/
const ticketExpiresIn = 60*30 - 30 // 12h (43200)

var (
	ErrTicketUpdateForbidden = errors.New("can NOT update ticket in wxwork suite lite mode")
	ErrTokenUpdateForbidden  = errors.New("can NOT refresh&update token in wxwork suite lite mode")
)

type Config struct {
	SuiteID        string // 企业服务商应用ID
	SuiteSecret    string // 企业服务商应用密钥
	Token          string
	EncodingAESKey string
}

type WxWorkSuite struct {
	Config           *Config
	Client           *utils.Client
	ticketCache      *utils.AccessTokenCache
	accessTokenCache *utils.AccessTokenCache
}

func New(
	cache utils.Cache,
	locker utils.Lock,
	config *Config,
	tokenRefreshHandler utils.TokenRefreshHandler, // 刷新callback
) *WxWorkSuite {
	ticketCache := utils.NewAccessTokenCache(newTicketAdapter(config.SuiteID), cache, locker)
	accessTokenCache := utils.NewAccessTokenCache(
		newAccessTokenAdaptor(config, ticketCache),
		cache, locker,
		utils.CacheClientTokenOptWithExpireBefore(tokenRefreshHandler),
	)
	instance := &WxWorkSuite{
		Config:           config,
		Client:           utils.NewClient(WXServerUrl, accessTokenCache),
		ticketCache:      ticketCache,
		accessTokenCache: accessTokenCache,
	}
	instance.Client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	return instance
}

func NewLite(
	cache utils.Cache,
	locker utils.Lock,
	suiteID string,
) *WxWorkSuite {
	config := &Config{SuiteID: suiteID}
	instance := &WxWorkSuite{
		Config: config,
		Client: utils.NewClient(WXServerUrl, utils.NewAccessTokenCache(
			newAccessTokenAdaptor(config, nil), cache, locker,
		)),
	}
	instance.Client.UpdateAccessTokenKey(accessTokenKey) // token的名称不一样
	return instance
}

// 当收到EventComponentVerifyTicket时， 用于更新ticket到cache
func (suite *WxWorkSuite) UpdateTicket(
	ctx context.Context, token string,
) error {
	if suite.ticketCache == nil {
		return fmt.Errorf(
			"wxcorp suite appid : %s, error: %w", suite.Config.SuiteID, ErrTicketUpdateForbidden,
		)
	}
	_, err := suite.ticketCache.UpdateAccessToken(ctx, token, ticketExpiresIn)
	return err
}

func (suite *WxWorkSuite) RefreshAccessToken(
	ctx context.Context, expireBefore int,
) (string, error) {
	if suite.accessTokenCache == nil {
		return "", fmt.Errorf(
			"wxopen appid : %s, error: %w", suite.Config.SuiteID, ErrTokenUpdateForbidden,
		)
	}
	return suite.accessTokenCache.RefreshAccessToken(ctx, expireBefore)
}

func (suite *WxWorkSuite) ClearAccessToken(ctx context.Context) error {
	if suite.accessTokenCache == nil {
		return fmt.Errorf(
			"suiteid : %s, error: %w",
			suite.Config.SuiteID, ErrTokenUpdateForbidden,
		)
	}
	return suite.accessTokenCache.ClearAccessToken(ctx)
}
