package main

import (
	"fmt"
	"sync"

	"github.com/lixinio/weixin/utils"
)

type TokenCache interface {
	SetAccessToken(string, int) error  // 保存 Access token
	GetAceessToken() (string, error)   // 获取 Access token
	SetPermanentCode(string) error     // 保存 Permanent code
	GetPermanentCode() (string, error) // 获取 PermanentCode code
}

// 获取特定的 服务号&小程序 的 TokenCache
type TokenCacheManager interface {
	GetTokenCache(string, string, int) (TokenCache, error)
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

func (mng *AuthorizerTokenCacheManager) getCacheKey(
	suiteID, corpID string, agentID int,
) string {
	return fmt.Sprintf("suite-authorizer.%s.%s.%d", suiteID, corpID, agentID)
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
	suiteID, corpID string, agentID int,
) (TokenCache, error) {
	key := mng.getCacheKey(suiteID, corpID, agentID)
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

	cache = NewAuthorizerTokenCache(suiteID, corpID, agentID, mng.cache, mng.locker)
	mng.caches[key] = cache
	return cache, nil
}
