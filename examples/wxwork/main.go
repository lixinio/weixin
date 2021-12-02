package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/lixinio/weixin/wxwork/server_api"
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
<br/>
<br/>
<br/>
<form method="get" action="/jsapi/oa">
<button type="submit">OA审批</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
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
		user, err := agent.GetUserInfo(r.Context(), code)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, user.UserID)
	}
}

func main() {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agent.New(corp, redis, redis, &agent.Config{
		AgentID: test.AgentID,
		Secret:  test.AgentSecret,
	})

	serverApi := server_api.NewApi(
		test.AgentID,
		test.AgentToken,
		test.AgentEncodingAESKey,
	)

	http.HandleFunc("/", index(agent))
	http.HandleFunc("/login", login(agent, false))
	http.HandleFunc("/login_sso", login(agent, true))
	http.HandleFunc("/login/callback", callback(agent))
	http.HandleFunc("/jsapi/oa", jsapiOA(agent))

	// static
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 域名校验
	for _, ver := range test.WxWorkAgentDomainVer {
		http.HandleFunc(
			fmt.Sprintf("/%s", ver.Filename),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, ver.Content)
			},
		)
	}

	http.HandleFunc(fmt.Sprintf("/weixin/%s/%d", test.CorpID, test.AgentID), msgCallback(serverApi))
	http.HandleFunc(fmt.Sprintf("/weixin/%s/%s", test.CorpID, "0"), msgCallback(serverApi))

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
