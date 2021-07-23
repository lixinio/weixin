package utils

import (
	"time"
)

const (
	defaultExpireBefore     int           = 300                    // 缺省提前5分钟
	defaultLockTimeout      time.Duration = 60 * time.Second       // 加锁的时长， 过期自动释放锁
	defaultLockRetryTimeout time.Duration = 200 * time.Millisecond // 加锁失败再次重试的休眠时长
	defaultLockRetryTime    time.Duration = 60 * time.Second       // 加锁失败重试的总时长
)

// AccessTokenGetter 获取token接口， 微信服务号/小程序 企业微信应用， 服务商等都实现该接口
type AccessTokenGetter interface {
	GetAccessToken() (string, int, error) // 直接获取token， 不做任何缓存
	GetAccessTokenKey() string            // 获取token的key
	GetAccessTokenLockKey() string        // 获取token的lock key
}

// AccessTokenCache token缓存对象
type AccessTokenCache struct {
	cache             Cache             // 用来缓存token的容器
	accessTokenLock   Lock              // 避免刷新token冲突
	accessTokenGetter AccessTokenGetter // 获取token对象
	expireBefore      int               // 提前多少秒重刷
}

func NewAccessTokenCache(
	accessTokenGetter AccessTokenGetter,
	cache Cache,
	locker Lock,
	expireBefore int,
) *AccessTokenCache {
	if expireBefore == 0 {
		expireBefore = defaultExpireBefore
	}
	return &AccessTokenCache{
		accessTokenGetter: accessTokenGetter,
		accessTokenLock:   locker,
		cache:             cache,
		expireBefore:      expireBefore,
	}
}

// GetAccessToken 刷新token， 优先从缓存获取
func (atc *AccessTokenCache) GetAccessToken() (accessToken string, err error) {
	accessToken, err = atc.getCachedAccessToken()
	if err == nil && accessToken != "" {
		// 直接从缓存获取
		return accessToken, nil
	} else if err != nil {
		// 出错了， 直接报错， 而不是用不缓存的Token， 因为获取Token有次数限制
		return "", err
	}

	//加上lock，是为了防止在并发获取token时，cache刚好失效，导致从服务器上获取到不同token
	lockKey := atc.accessTokenGetter.GetAccessTokenLockKey()
	locked, err := atc.accessTokenLock.LockTimeout(
		lockKey, defaultLockTimeout, defaultLockRetryTime, defaultLockRetryTimeout,
	)
	if err != nil || !locked {
		// 出错或者加锁失败
		return "", err
	}
	defer atc.accessTokenLock.UnLock(lockKey)

	// 是不是别人已经获取到Token了
	accessToken, err = atc.getCachedAccessToken()
	if err == nil && accessToken != "" {
		return accessToken, nil
	} else if err != nil {
		return "", err
	}

	// 直接从服务器刷新
	return atc.refreshAccessToken()
}

func (atc *AccessTokenCache) getCachedAccessToken() (accessToken string, err error) {
	accessTokenCacheKey := atc.accessTokenGetter.GetAccessTokenKey()
	exist := false
	exist, err = atc.cache.Get(accessTokenCacheKey, &accessToken)
	if err == nil {
		if exist {
			// 存在内容
			return
		}
		// 不存在， 返回空token
		return "", nil
	}
	return "", err
}

func (atc *AccessTokenCache) refreshAccessToken() (accessToken string, err error) {
	// 从服务器获取Token
	expiresIn := 0
	accessToken, expiresIn, err = atc.accessTokenGetter.GetAccessToken()
	if err != nil {
		// 失败
		return
	}

	// 减去提前刷新的时间
	expires := expiresIn - atc.expireBefore // 秒
	accessTokenCacheKey := atc.accessTokenGetter.GetAccessTokenKey()
	err = atc.cache.Set(accessTokenCacheKey, accessToken, time.Duration(expires)*time.Second)
	if err != nil {
		// 如果存到缓存失败， token依然是可用的，
		// 因为如果缓存出了问题， 下次刷新Token也会失败， 不会导致token配额用尽
		return accessToken, nil
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
