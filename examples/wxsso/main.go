package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/web_sso"
)

func index(oa *web_sso.WebSSO) http.HandlerFunc {
	html := `
<form method="get" action="/login">
    <button type="submit">微信登录</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	}
}

func login(oa *web_sso.WebSSO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/login/callback", r.Host)
		url = oa.GetAuthorizeUrl(url, "state")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func callback(oa *web_sso.WebSSO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		snsAccessToken, err := oa.GetSnsAccessToken(r.Context(), code)
		if err != nil {
			fmt.Println(err)
			utils.HttpAbort(w, http.StatusForbidden)
			return
		}
		fmt.Println(snsAccessToken.Scope, snsAccessToken.AccessToken, snsAccessToken.Openid)

		user_info, err := oa.GetUserInfo(
			r.Context(),
			snsAccessToken.AccessToken,
			snsAccessToken.Openid,
			"",
		)
		if err != nil {
			fmt.Println(err)
			utils.HttpAbort(w, http.StatusForbidden)
			return
		}
		fmt.Println(user_info.Nickname, user_info.Openid, user_info.Headimgurl)

		snsAccessToken2, err := oa.RefreshSnsToken(r.Context(), snsAccessToken.RefreshToken)
		if err != nil {
			fmt.Println(err)
			utils.HttpAbort(w, http.StatusForbidden)
			return
		}
		fmt.Println(snsAccessToken2.Scope, snsAccessToken2.AccessToken, snsAccessToken2.Openid)

		fmt.Fprintf(w, user_info.Nickname)
	}
}

func main() {
	sso := web_sso.New(&web_sso.Config{
		Appid:  test.WxSSOAppid,
		Secret: test.WxSSOSecret,
	})

	http.HandleFunc("/", index(sso))
	http.HandleFunc("/login", login(sso))
	http.HandleFunc("/login/callback", callback(sso))

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
