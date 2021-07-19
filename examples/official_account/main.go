package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/lixinio/weixin/weixin/server_api"
)

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func index(oa *official_account.OfficialAccount) http.HandlerFunc {
	html := `
<form method="get" action="/login">
    <button type="submit">微信登录</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, html)
	}
}

func login(oa *official_account.OfficialAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/login/callback", r.Host)
		url = oa.GetAuthorizeUrl(url, official_account.ScopeSnsapiUserinfo, "state")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func callback(oa *official_account.OfficialAccount) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		snsAccessToken, err := oa.GetSnsAccessToken(code)
		if err != nil {
			fmt.Println(err)
			httpAbort(w, http.StatusForbidden)
			return
		}

		user_info, err := oa.GetUserInfo(snsAccessToken.AccessToken, snsAccessToken.Openid, "")
		if err != nil {
			fmt.Println(err)
			httpAbort(w, http.StatusForbidden)
			return
		}

		fmt.Fprintf(w, user_info.Nickname)
	}
}

func main() {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})
	serverApi := server_api.NewOfficialAccountApi(
		test.OfficialAccountToken,
		test.OfficialAccountAESKey,
		officialAccount,
	)

	http.HandleFunc("/", index(officialAccount))
	http.HandleFunc("/login", login(officialAccount))
	http.HandleFunc("/login/callback", callback(officialAccount))
	http.HandleFunc(fmt.Sprintf("/%s", test.OfficialAccountAuthKey), func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, test.OfficialAccountAuthValue)
	})
	http.HandleFunc(fmt.Sprintf("/weixin/%s", test.OfficialAccountAppid), weixinCallback(serverApi))

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}

}
