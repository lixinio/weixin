package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork_provider"
)

func index(provider *wxwork_provider.WxWorkProvider) http.HandlerFunc {
	html := `
<br/>
<br/>
<form method="get" action="/login_sso">
<button type="submit">网页扫码登录</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	}
}

func login(provider *wxwork_provider.WxWorkProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/login/callback", r.Host)
		url = provider.GetAuthorizeUrl(url, "admin", "state")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func callback(provider *wxwork_provider.WxWorkProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("auth_code")
		user, err := provider.GetLoginInfo(r.Context(), code)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, user.UserInfo.UserID)
	}
}

func main() {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	provider := wxwork_provider.New(redis, redis, &wxwork_provider.Config{
		CorpID:         test.WxWorkProviderCorpID,
		ProviderSecret: test.WxWorkProviderSecret,
	})

	http.HandleFunc("/", index(provider))
	http.HandleFunc("/login_sso", login(provider))
	http.HandleFunc("/login/callback", callback(provider))

	// http.HandleFunc(fmt.Sprintf("/weixin/%s/%d", test.CorpID, test.AgentID), msgCallback(serverApi))
	// http.HandleFunc(fmt.Sprintf("/weixin/%s/%s", test.CorpID, "0"), msgCallback(serverApi))

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
