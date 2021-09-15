package main

import (
	"context"
	"html/template"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/authorizer"
	oa_api "github.com/lixinio/weixin/wxwork/oa"
)

func initOA(ctx context.Context, authorizer *authorizer.Authorizer) string {
	oaApi := oa_api.NewApi(authorizer.Client)
	tmplID, err := oaApi.CopyTemplate(ctx, test.WxWorkSuiteOaTmplID)
	if err != nil {
		panic(err)
	}
	return tmplID
}

// js代码来源
// https://github.com/liyuexi/qywx-third-java/blob/1a80b862f1196822cf4645a2e49fd97a31c71b82/src/main/resources/templates/oa/index.html
func jsapiOA(a *authorizer.Authorizer) http.HandlerFunc {
	// 自建应用和服务商是一样的内容
	t, _ := template.ParseFiles("static/oa.html")
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://" + r.Host + r.RequestURI
		corpConfig, err := a.GetCorpJSApiConfig(r.Context(), url)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return
		}
		agentConfig, err := a.GetAgentJSApiConfig(r.Context(), url)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return
		}

		val := struct {
			Corp    *authorizer.JsApiCorpConfig
			Agent   *authorizer.JsApiAgentConfig
			TmplID  string
			ThirdNo string
		}{
			Corp:    corpConfig,
			Agent:   agentConfig,
			TmplID:  test.WxWorkSuiteOaTmplID,
			ThirdNo: test.WxWorkSuiteOaThirdNo,
		}
		// t.Execute(os.Stdout, val)
		err = t.Execute(w, val)
		if err != nil {
			utils.HttpAbortBadRequest(w)
			return
		}
	}
}
