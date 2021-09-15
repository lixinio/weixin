package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/weixin/server_api"
	"github.com/lixinio/weixin/wxopen"
)

func index(wo *wxopen.WxOpen) http.HandlerFunc {
	html := `
<form method="get" action="/auth">
    <button type="submit">PC扫码授权</button>
</form>
<br/>
<br/>
<br/>
<form method="get" action="/auth/mobile">
    <button type="submit">H5扫码授权</button>
</form>
<br/>
<br/>
<br/>
<form method="get" action="/login">
    <button type="submit">H5登录</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	}
}

func auth(wo *wxopen.WxOpen, mobile bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/auth/callback", r.Host)
		preAuthCode, _, err := wo.CreatePreAuthCode(r.Context())
		if err != nil {
			panic(err)
		}
		redirectUri := ""
		if !mobile {
			redirectUri = wo.GetComponentLoginPage(preAuthCode, url, wxopen.AuthTypeAll, "")
		} else {
			redirectUri = wo.GetComponentLoginH5Page(preAuthCode, url, wxopen.AuthTypeAll, "")
		}

		http.Redirect(w, r, redirectUri, http.StatusFound)
	}
}

func callback(wo *wxopen.WxOpen, manager TokenCacheManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("auth_code")
		authorizationInfo, err := wo.QueryAuth(r.Context(), code)
		if err != nil {
			panic(err)
		}

		authorizer, err := manager.GetTokenCache(
			wo.Config.Appid, authorizationInfo.AuthorizerAppid,
		)
		if err != nil {
			panic(err)
		}

		// 保存refresh token
		err = authorizer.SetRefreshToken(authorizationInfo.AuthorizerRefreshToken)
		if err != nil {
			panic(err)
		}

		// 保存 access token
		err = authorizer.SetAccessToken(
			authorizationInfo.AuthorizerAccessToken,
			authorizationInfo.ExpiresIn,
		)
		if err != nil {
			panic(err)
		}

		fmt.Printf("appid %s %s\n", wo.Config.Appid, authorizationInfo.AuthorizerAppid)
		for i, f := range authorizationInfo.FuncInfo {
			fmt.Printf("func info %d, %d\n", i, f.FuncscopeCategory.ID)
		}

		// 获取完整信息
		detail, err := wo.GetAuthorizerInfo(r.Context(), authorizationInfo.AuthorizerAppid)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s %s %s %s %s\n",
			detail.AuthorizerInfo.NickName,
			detail.AuthorizerInfo.PrincipalName,
			detail.AuthorizerInfo.UserName,
			detail.AuthorizerInfo.Alias,
			detail.AuthorizerInfo.HeadImg,
		)

		fmt.Fprint(w, "ok")
	}
}

func login(wo *wxopen.WxOpen, wxopenOA *authorizer.Authorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/login/callback", r.Host)
		redirectUri := wo.GetAuthorizeUrl(wxopenOA.Appid, url, "snsapi_userinfo", "test")
		http.Redirect(w, r, redirectUri, http.StatusFound)
	}
}

func loginCallback(wo *wxopen.WxOpen, wxopenOA *authorizer.Authorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		snsAccessToken, err := wo.GetSnsAccessToken(r.Context(), wxopenOA.Appid, code)
		if err != nil {
			fmt.Println(err)
			utils.HttpAbort(w, http.StatusForbidden)
			return
		}
		fmt.Println(snsAccessToken.Scope, snsAccessToken.AccessToken, snsAccessToken.Openid)

		user_info, err := wo.GetUserInfo(
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

		snsAccessToken2, err := wo.RefreshSnsToken(
			r.Context(), wxopenOA.Appid, snsAccessToken.RefreshToken,
		)
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
	// 缓存
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	// 存储授权后的access token, refresh token
	// 生产环境应该存储到DB
	manager := NewAuthorizerRefreshTokenManager(redis, redis)
	// 开放平台
	wxopenApi := wxopen.New(redis, redis, &wxopen.Config{
		Appid:          test.WxOpenAppid,
		Secret:         test.WxOpenSecret,
		Token:          test.WxOpenToken,
		EncodingAESKey: test.WxOpenEncodingAESKey,
	})

	// 初始化测试的授权服务号
	oaTokenCache, err := manager.GetTokenCache(wxopenApi.Config.Appid, test.WxOpenOAAppid)
	if err != nil {
		panic(err)
	}
	wxopenOA := authorizer.New(
		redis, redis,
		wxopenApi.Config.Appid,
		test.WxOpenOAAppid,
		GetAuthorizerAccessToken(wxopenApi, oaTokenCache, test.WxOpenOAAppid),
	)

	// server
	serverApi := server_api.NewApi(
		wxopenApi.Config.Appid, // 这里必须是服务商的AppID
		test.WxOpenToken,
		test.WxOpenEncodingAESKey,
		wxopenOA.Client,
	)

	// 域名校验
	for _, ver := range test.WxOpenDomainVer {
		http.HandleFunc(
			fmt.Sprintf("/%s", ver.Filename),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, ver.Content)
			},
		)
	}

	http.HandleFunc(
		fmt.Sprintf("/gateway/component/%s/notify", wxopenApi.Config.Appid),
		weixinCallback(wxopenApi, wxopen.ReleaseAppIDS),
	)
	http.HandleFunc(
		fmt.Sprintf(
			"/gateway/component/%s/authorizer/%s/callback",
			wxopenApi.Config.Appid,
			wxopenOA.Appid,
		),
		authorizerCallback(serverApi),
	)

	// 全网发布
	for appid := range wxopen.ReleaseAppIDS {
		http.HandleFunc(
			fmt.Sprintf(
				"/gateway/component/%s/authorizer/%s/callback",
				wxopenApi.Config.Appid, appid,
			),
			releaseCallback(wxopenApi, serverApi),
		)
	}

	http.HandleFunc("/", index(wxopenApi))
	http.HandleFunc("/auth", auth(wxopenApi, false))
	http.HandleFunc("/auth/mobile", auth(wxopenApi, true))
	http.HandleFunc("/auth/callback", callback(wxopenApi, manager))
	http.HandleFunc("/login", login(wxopenApi, wxopenOA))
	http.HandleFunc("/login/callback", loginCallback(wxopenApi, wxopenOA))

	// 刷新Token
	go RefreshWxOpenToken(wxopenApi)
	go RefreshAuthorizerToken([]*authorizer.Authorizer{wxopenOA})

	err = http.ListenAndServe(":5001", nil)
	if err != nil {
		fmt.Println(err)
	}
}
