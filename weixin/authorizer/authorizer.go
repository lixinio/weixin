package authorizer

import (
	"context"
	"errors"
	"fmt"

	"github.com/lixinio/weixin/utils"
)

const (
	WXServerUrl = "https://api.weixin.qq.com" // 微信 api 服务器地址
)

var (
	ErrTokenUpdateForbidden  = errors.New("can NOT refresh&update token in wxopen lite mode")
	ErrJsApiTicketForbidden  = errors.New("can NOT refresh&update jsapi ticket without enable it")
	ErrWxCardTicketForbidden = errors.New("can NOT refresh&update wxcard ticket without enable it")
)

type Authorizer struct {
	ComponentAppid    string
	Appid             string
	Client            *utils.Client
	accessTokenCache  *utils.AccessTokenCache // 用于支持手动刷新Token， Client不对外暴露该对象
	jsApiTicketCache  *utils.AccessTokenCache
	wxCardTicketCache *utils.AccessTokenCache
}

func New(
	cache utils.Cache,
	locker utils.Lock,
	componentAppid, appid string,
	accessTokenGetter RefreshAccessToken,
) *Authorizer {
	accessTokenCache := utils.NewAccessTokenCache(
		newAdapter(componentAppid, appid, accessTokenGetter), cache, locker,
	)
	return &Authorizer{
		ComponentAppid:   componentAppid,
		Appid:            appid,
		Client:           utils.NewClient(WXServerUrl, accessTokenCache),
		accessTokenCache: accessTokenCache,
	}
}

// 不支持刷新Token的简化模式
func NewLite(
	cache utils.Cache,
	locker utils.Lock,
	componentAppid, appid string,
) *Authorizer {
	accessTokenCache := utils.NewAccessTokenCache(
		newAdapter(componentAppid, appid, func(context.Context) (string, int, error) {
			return "", 0, fmt.Errorf(
				"can NOT refresh token in lite mod, appid(%s , %s), %w",
				componentAppid, appid, ErrTokenUpdateForbidden,
			)
		}), cache, locker,
	)
	return &Authorizer{
		ComponentAppid: componentAppid,
		Appid:          appid,
		Client:         utils.NewClient(WXServerUrl, accessTokenCache),
	}
}

func (authorizer *Authorizer) RefreshAccessToken(
	ctx context.Context, expireBefore int,
) (string, error) {
	if authorizer.accessTokenCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrTokenUpdateForbidden,
		)
	}
	return authorizer.accessTokenCache.RefreshAccessToken(ctx, expireBefore)
}

func (authorizer *Authorizer) ClearAccessToken(ctx context.Context) error {
	if authorizer.accessTokenCache == nil {
		return fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrTokenUpdateForbidden,
		)
	}
	return authorizer.accessTokenCache.ClearAccessToken(ctx)
}

func (authorizer *Authorizer) RefreshJsApiTicket(
	ctx context.Context, expireBefore int,
) (string, error) {
	if authorizer.jsApiTicketCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrJsApiTicketForbidden,
		)
	}
	return authorizer.jsApiTicketCache.RefreshAccessToken(ctx, expireBefore)
}

func (authorizer *Authorizer) ClearJsApiTicket(ctx context.Context) error {
	if authorizer.jsApiTicketCache == nil {
		return fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrJsApiTicketForbidden,
		)
	}
	return authorizer.jsApiTicketCache.ClearAccessToken(ctx)
}

func (authorizer *Authorizer) RefreshWxCardTicket(
	ctx context.Context, expireBefore int,
) (string, error) {
	if authorizer.wxCardTicketCache == nil {
		return "", fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrWxCardTicketForbidden,
		)
	}
	return authorizer.wxCardTicketCache.RefreshAccessToken(ctx, expireBefore)
}

func (authorizer *Authorizer) ClearWxCardTicket(ctx context.Context) error {
	if authorizer.wxCardTicketCache == nil {
		return fmt.Errorf(
			"authorizer appid : %s,%s, error: %w",
			authorizer.ComponentAppid, authorizer.Appid,
			ErrWxCardTicketForbidden,
		)
	}
	return authorizer.wxCardTicketCache.ClearAccessToken(ctx)
}
