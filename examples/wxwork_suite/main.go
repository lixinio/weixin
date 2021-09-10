package main

import (
	"fmt"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork/authorizer"
	"github.com/lixinio/weixin/wxwork_suite"
)

func index(suite *wxwork_suite.WxWorkSuite) http.HandlerFunc {
	html := `
<br/>
<br/>
<form method="get" action="/login_sso">
<button type="submit">网页登录</button>
</form>
<br/>
<br/>
<form method="get" action="/install">
<button type="submit">扫码授权</button>
</form>
	`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
	}
}

func login(suite *wxwork_suite.WxWorkSuite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("http://%s/login/callback", r.Host)
		url = suite.GetAuthorizeUrl(url, "snsapi_userinfo", "state")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func callback(suite *wxwork_suite.WxWorkSuite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		user, err := suite.GetUserInfo3rd(r.Context(), code)
		if err != nil {
			panic(err)
		}
		fmt.Printf("corp %s user %s\n", user.CorpID, user.UserID)

		detail, err := suite.GetUserDetail3rd(r.Context(), user.UserTicket)
		if err != nil {
			panic(err)
		}
		fmt.Printf("name %s avatar %s qrcode %s\n", detail.Name, detail.Avatar, detail.QrCode)
		fmt.Fprintf(w, detail.Name)
	}
}

func install(suite *wxwork_suite.WxWorkSuite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		preAuthCode, err := suite.GetPreAuthCode(r.Context())
		if err != nil {
			panic(err)
		}
		err = suite.SetSessionInfo(r.Context(), preAuthCode.PreAuthCode, 1)
		if err != nil {
			panic(err)
		}

		url := fmt.Sprintf("http://%s/install/callback", r.Host)
		url = suite.GetInstallUrl(url, preAuthCode.PreAuthCode, "state")
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func install_callback(
	suite *wxwork_suite.WxWorkSuite, manager TokenCacheManager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("auth_code")
		info, err := suite.GetPermanentCode(r.Context(), code)
		if err != nil {
			panic(err)
		}
		info2, err := suite.GetAuthInfo(r.Context(), info.AuthCorpInfo.CorpID, info.PermanentCode)
		if err != nil {
			panic(err)
		}
		fmt.Printf(
			"corp %s %s %s %s\n",
			info.AuthCorpInfo.CorpID, info.AuthCorpInfo.CorpName,
			info.AuthCorpInfo.CorpFullName, info.AuthCorpInfo.CorpIndustry,
		)
		fmt.Printf(
			"corp %s %s %s\n", info.AuthUserInfo.Name,
			info.AuthUserInfo.UserID, info.AuthUserInfo.OpenUserID,
		)
		fmt.Printf(
			"%d %s\n",
			info2.AuthInfo.Agents[0].AgentID, info2.AuthInfo.Agents[0].Name,
		)

		admins, err := suite.GetAdminList(
			r.Context(),
			info.AuthCorpInfo.CorpID,
			info.AuthInfo.Agents[0].AgentID,
		)
		if err != nil {
			panic(err)
		}
		for _, admin := range admins.Admin {
			fmt.Printf("admin %s %s %d\n", admin.UserID, admin.OpenUserID, admin.AuthType)
		}

		token, err := suite.GetCorpToken(r.Context(), info.AuthCorpInfo.CorpID, info.PermanentCode)
		if err != nil {
			panic(err)
		}
		fmt.Println("access token ", token.AccessToken)

		authorizer, err := manager.GetTokenCache(
			suite.Config.SuiteID, info.AuthCorpInfo.CorpID, info.AuthInfo.Agents[0].AgentID,
		)
		if err != nil {
			panic(err)
		}

		// 存起来
		if err = authorizer.SetAccessToken(token.AccessToken, token.ExpiresIn); err != nil {
			panic(err)
		}
		if err = authorizer.SetPermanentCode(info.PermanentCode); err != nil {
			panic(err)
		}

		fmt.Fprint(w, info.AuthCorpInfo.CorpFullName)
	}
}

func main() {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	suite := wxwork_suite.New(redis, redis, &wxwork_suite.Config{
		SuiteID:        test.WxWorkSuiteID,
		SuiteSecret:    test.WxWorkSuiteSecret,
		Token:          test.WxWorkSuiteToken,
		EncodingAESKey: test.WxWorkSuiteEncodingAESKey,
	})

	// 存储授权后的access token, refresh token
	// 生产环境应该存储到DB
	manager := NewAuthorizerRefreshTokenManager(redis, redis)

	// 初始化测试的授权服务号
	oaTokenCache, err := manager.GetTokenCache(
		suite.Config.SuiteID, test.WxWorkSuiteCorpID, test.WxWorkSuiteAgentID,
	)
	if err != nil {
		panic(err)
	}
	wxworkSuiteAgent := authorizer.New(
		redis, redis,
		suite.Config.SuiteID,
		test.WxWorkSuiteCorpID,
		test.WxWorkSuiteAgentID,
		GetAuthorizerAccessToken(suite, oaTokenCache, test.WxWorkSuiteCorpID),
	)
	// 立即刷新Token
	if _, err = wxworkSuiteAgent.RefreshAccessToken(0); err != nil {
		fmt.Printf("refresh token fail %s", err.Error())
	}

	http.HandleFunc("/", index(suite))
	http.HandleFunc("/login_sso", login(suite))
	http.HandleFunc("/login/callback", callback(suite))
	http.HandleFunc("/install", install(suite))
	http.HandleFunc("/install/callback", install_callback(suite, manager))

	http.HandleFunc(
		fmt.Sprintf("/weixin/%s/%s/data", test.WxWorkProviderCorpID, test.WxWorkSuiteID),
		weixinCallback(suite),
	)
	http.HandleFunc(
		fmt.Sprintf("/weixin/%s/%s/cmd", test.WxWorkProviderCorpID, test.WxWorkSuiteID),
		weixinCallback(suite),
	)

	// 域名校验
	for _, ver := range test.WxWorkSuiteDomainVer {
		http.HandleFunc(
			fmt.Sprintf("/%s", ver.Filename),
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, ver.Content)
			},
		)
	}

	err = http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
