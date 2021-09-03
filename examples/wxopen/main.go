package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/authorizer"
	"github.com/lixinio/weixin/weixin/server_api"
	"github.com/lixinio/weixin/wxopen"
)

func httpAbort(w http.ResponseWriter, code int) {
	w.WriteHeader(http.StatusBadRequest)
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func index(wo *wxopen.WxOpen) http.HandlerFunc {
	html := `
<form method="get" action="/auth">
    <button type="submit">PC扫码授权</button>
</form>
<form method="get" action="/auth/mobile">
    <button type="submit">H5扫码授权</button>
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

func main() {
	// 缓存
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	// 存储授权后的access token, refresh token
	// 生产环境应该存储到DB
	manager := NewAuthorizerRefreshTokenManager(redis, redis)
	// 开放平台
	wxopen := wxopen.New(redis, redis, &wxopen.Config{
		Appid:          test.WxOpenAppid,
		Secret:         test.WxOpenSecret,
		Token:          test.WxOpenToken,
		EncodingAESKey: test.WxOpenEncodingAESKey,
	})

	// 初始化测试的授权服务号
	oaTokenCache, err := manager.GetTokenCache(wxopen.Config.Appid, test.WxOpenOAAppid)
	if err != nil {
		panic(err)
	}
	wxopenOA := authorizer.New(
		redis, redis,
		wxopen.Config.Appid,
		test.WxOpenOAAppid,
		GetAuthorizerAccessToken(wxopen, oaTokenCache, test.WxOpenOAAppid),
	)
	// 刷新Token, 如果第一次运行， 可能会失败， 因为没有wxopen主动推送的ticket
	wxopenOA.RefreshAccessToken()

	// server
	serverApi := server_api.NewApi(
		wxopen.Config.Appid, // 这里必须是服务商的AppID
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
		fmt.Sprintf("/gateway/component/%s/notify", test.WxOpenAppid),
		weixinCallback(wxopen),
	)
	http.HandleFunc(
		fmt.Sprintf(
			"/gateway/component/%s/authorizer/%s/callback",
			test.WxOpenAppid,
			test.WxOpenOAAppid,
		),
		authorizerCallback(serverApi),
	)
	http.HandleFunc("/", index(wxopen))
	http.HandleFunc("/auth", auth(wxopen, false))
	http.HandleFunc("/auth/mobile", auth(wxopen, true))
	http.HandleFunc("/auth/callback", callback(wxopen, manager))

	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
