package main

import (
	"fmt"
	"net/url"

	"github.com/lixinio/weixin/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/lixinio/weixin/weixin/user_api"
)

func main() {
	cache := redis.NewRedis(&redis.Config{RedisUrl: "redis://127.0.0.1:6379/1"})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  "wx59864a9e578229ea",
		Secret: "e4a2a478789ccc9378d5e93533689c6b",
	})

	userApi := user_api.NewOfficialAccountApi(officialAccount)

	params := url.Values{}
	b, e := userApi.Get(params)
	fmt.Println(string(b), " ", e)

}
