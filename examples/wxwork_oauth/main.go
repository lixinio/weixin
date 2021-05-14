package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
)

func index(agent *agent.Agent) http.HandlerFunc {
	html := `
<form method="get" action="/login">
    <button type="submit">登 录</button>
</form>
<br/>
<br/>
<br/>
<br/>
<form method="get" action="/login_sso">
<button type="submit">单点登录</button>
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
	cache := redis.NewRedis(&redis.Config{RedisUrl: "redis://127.0.0.1:6379/1"})
	corp := wxwork.New(&wxwork.Config{
		Corpid: "wx247d4bc469342dc4",
	})
	agent := agent.New(corp, cache, &agent.Config{
		AgentId: "20",
		Secret:  "G9x8iHpoQMJ8ynDgcplAvwiF4qWF1tRJ3gMVShXZ1Ks",
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
