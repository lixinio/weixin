package utils

import (
	"sync"
	"time"
)

const (
	defaultExpireBefore int = 300 // 缺省提前5分钟
)

// AccessTokenGetter 获取token接口， 微信服务号/小程序 企业微信应用， 服务商等都实现该接口
type AccessTokenGetter interface {
	GetAccessToken() (string, int, error) // 直接获取token， 不做任何缓存
	GetAccessTokenKey() string            // 获取token的key
}

// AccessTokenCache token缓存对象
type AccessTokenCache struct {
	cache             Cache             // 用来缓存token的容器
	accessTokenLock   *sync.Mutex       // 避免刷新token冲突
	accessTokenGetter AccessTokenGetter // 获取token对象
	expireBefore      int               // 提前多少秒重刷
}

func NewAccessTokenCache(
	accessTokenGetter AccessTokenGetter,
	cache Cache,
	expireBefore int,
) *AccessTokenCache {
	if expireBefore == 0 {
		expireBefore = defaultExpireBefore
	}
	return &AccessTokenCache{
		accessTokenGetter: accessTokenGetter,
		accessTokenLock:   new(sync.Mutex),
		cache:             cache,
		expireBefore:      expireBefore,
	}
}

// GetAccessToken 刷新token， 优先从缓存获取
func (atc *AccessTokenCache) GetAccessToken() (accessToken string, err error) {
	//加上lock，是为了防止在并发获取token时，cache刚好失效，导致从微信服务器上获取到不同token
	atc.accessTokenLock.Lock()
	defer atc.accessTokenLock.Unlock()

	accessTokenCacheKey := atc.accessTokenGetter.GetAccessTokenKey()
	val := atc.cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	//cache失效，从微信服务器获取
	expiresIn := 0
	accessToken, expiresIn, err = atc.accessTokenGetter.GetAccessToken()
	if err != nil {
		return
	}

	expires := expiresIn - atc.expireBefore
	err = atc.cache.Set(accessTokenCacheKey, accessToken, time.Duration(expires)*time.Second)
	if err != nil {
		return
	}
	return
}

// TokenResponse 刷新token相应体
type TokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   int     `json:"expires_in"`
	Errcode     float64 `json:"errcode"`
	Errmsg      string  `json:"errmsg"`
}
