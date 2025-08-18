package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/lixinio/weixin/utils"
)

type TokenCache interface {
	SetAccessToken(context.Context, string, int) error // 保存 Access token
	GetAccessToken(context.Context) (string, error)    // 获取 Access token
	SetRefreshToken(context.Context, string) error     // 保存 Refresh token
	GetRefreshToken(context.Context) (string, error)   // 获取 Refresh token
}

// 获取特定的 服务号&小程序 的 TokenCache
type TokenCacheManager interface {
	GetTokenCache(string, string) (TokenCache, error)
}

type AuthorizerTokenCacheManager struct {
	lock         *sync.RWMutex
	caches       map[string]*AuthorizerTokenCache
	cache        utils.Cache
	locker       utils.Lock
	expireBefore int
}

func NewAuthorizerRefreshTokenManager(
	cache utils.Cache,
	locker utils.Lock,
) TokenCacheManager {
	return &AuthorizerTokenCacheManager{
		lock:         &sync.RWMutex{},
		caches:       make(map[string]*AuthorizerTokenCache),
		locker:       locker,
		cache:        cache,
		expireBefore: 0,
	}
}

func (mng *AuthorizerTokenCacheManager) getCacheKey(componentAppid, appid string) string {
	return fmt.Sprintf("authorizer.%s.%s", componentAppid, appid)
}

func (mng *AuthorizerTokenCacheManager) getCacheWithLock(key string) *AuthorizerTokenCache {
	mng.lock.RLock()
	defer mng.lock.RUnlock()
	return mng.getCache(key)
}

func (mng *AuthorizerTokenCacheManager) getCache(key string) *AuthorizerTokenCache {
	if v, ok := mng.caches[key]; ok {
		return v
	}
	return nil
}

func (mng *AuthorizerTokenCacheManager) GetTokenCache(
	componentAppid, appid string,
) (TokenCache, error) {
	key := mng.getCacheKey(componentAppid, appid)
	cache := mng.getCacheWithLock(key)
	if cache != nil {
		return cache, nil
	}

	// 加锁
	mng.lock.Lock()
	defer mng.lock.Unlock()

	// 再判断
	cache = mng.getCache(key)
	if cache != nil {
		return cache, nil
	}

	cache = NewAuthorizerTokenCache(componentAppid, appid, mng.cache, mng.locker)
	mng.caches[key] = cache
	return cache, nil
}
