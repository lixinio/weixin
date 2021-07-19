package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
)

func index(agent *agent.Agent) http.HandlerFunc {
	html := `
<form method="get" action="/login">
    <button type="submit">企业微信登录</button>
</form>
<br/>
<br/>
<br/>
<br/>
<form method="get" action="/login_sso">
<button type="submit">网页扫码登录</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, html)
	}
}

func login(agent *agent.Agent, sso bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/login/callback", r.Host)

		if !sso {
			url = agent.GetAuthorizeUrl(url, "state")
		} else {
			url = agent.GetSSOAuthorizeUrl(url, "state")
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func callback(agent *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		user, err := agent.GetUserInfo(code)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, user.UserID)
	}
}

func main() {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, cache, &agent.Config{
		AgentId: test.AgentID,
		Secret:  test.AgentSecret,
	})

	http.HandleFunc("/", index(agent))
	http.HandleFunc("/login", login(agent, false))
	http.HandleFunc("/login_sso", login(agent, true))
	http.HandleFunc("/login/callback", callback(agent))

	err := http.ListenAndServe(":9998", nil)
	if err != nil {
		fmt.Println(err)
	}
}
