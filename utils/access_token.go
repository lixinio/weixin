package utils

import (
	"time"
)

const (
	defaultExpireBefore                   = 300                    // 缺省提前5分钟
	defaultMinusTTL                       = 30                     // 存储token, 减掉的ttl, 避免失效的token仍存在cache中
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
}

type refreshTokenHandler func() (string, int, error)

func NewAccessTokenCache(
	accessTokenGetter AccessTokenGetter,
	cache Cache,
	locker Lock,
) *AccessTokenCache {
	return &AccessTokenCache{
		accessTokenGetter: accessTokenGetter,
		accessTokenLock:   locker,
		cache:             cache,
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

	return atc.updateAccessToken(
		atc.accessTokenGetter.GetAccessToken,
		true,
	)
}

// 清除Token, 某些应用需要在取消授权之后, 立即清除Token
func (atc *AccessTokenCache) ClearAccessToken() error {
	closer, err := atc.lock()
	if err != nil {
		return err
	}
	defer closer()
	return atc.cache.Delete(atc.accessTokenGetter.GetAccessTokenKey())
}

// 强制刷新Token, 为了避免Token到期争抢刷新, 一般会有定时任务在Token过期之前的某个时刻强制刷新
func (atc *AccessTokenCache) RefreshAccessToken(beforeTTL int) (accessToken string, err error) {
	if beforeTTL == 0 {
		beforeTTL = defaultExpireBefore
	}
	ttl, err := atc.cache.TTL(atc.accessTokenGetter.GetAccessTokenKey())
	if err != nil {
		return "", err
	}
	// 如果不存在, 返回-2
	if ttl > beforeTTL {
		// 未更新
		return "", nil
	}

	return atc.updateAccessToken(
		atc.accessTokenGetter.GetAccessToken,
		false,
	)
}

func (atc *AccessTokenCache) lock() (func(), error) {
	lockKey := atc.accessTokenGetter.GetAccessTokenLockKey()
	locked, err := atc.accessTokenLock.LockTimeout(
		lockKey, defaultLockTimeout, defaultLockRetryTime, defaultLockRetryTimeout,
	)
	if err != nil || !locked {
		// 出错或者加锁失败
		return nil, err
	}
	return func() {
		atc.accessTokenLock.UnLock(lockKey)
	}, nil
}

func (atc *AccessTokenCache) updateAccessToken(
	handler refreshTokenHandler,
	checkLatest bool,
) (accessToken string, err error) {
	//加上lock，是为了防止在并发获取token时，cache刚好失效，导致从服务器上获取到不同token
	closer, err := atc.lock()
	if err != nil {
		return "", err
	}
	defer closer()

	if checkLatest {
		// 是不是别人已经获取到Token了
		accessToken, err = atc.getCachedAccessToken()
		if err == nil && accessToken != "" {
			return accessToken, nil
		} else if err != nil {
			return "", err
		}
	}

	// 直接从服务器刷新
	return atc.refreshAccessToken(handler)
}

// 直接从外部更新Token, 并更新缓存
// 比如服务商模式的应用, 或者微信上报的ticket
func (atc *AccessTokenCache) UpdateAccessToken(
	token string,
	expiresIn int,
) (accessToken string, err error) {
	return atc.updateAccessToken(func() (string, int, error) {
		return token, expiresIn, nil
	}, false)
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

func (atc *AccessTokenCache) refreshAccessToken(
	handler refreshTokenHandler,
) (accessToken string, err error) {
	// 从服务器获取Token
	expiresIn := 0
	accessToken, expiresIn, err = handler()
	if err != nil {
		// 失败
		return
	}

	// 减去提前刷新的时间
	expires := expiresIn - defaultMinusTTL // 秒
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
	WeixinError
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
