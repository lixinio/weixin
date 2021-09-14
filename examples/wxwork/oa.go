package main

import (
	"html/template"
	"net/http"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/wxwork/agent"
)

// js代码来源
// https://github.com/liyuexi/qywx-third-java/blob/1a80b862f1196822cf4645a2e49fd97a31c71b82/src/main/resources/templates/oa/index.html
func jsapiOA(a *agent.Agent) http.HandlerFunc {
	// 自建应用和服务商是一样的内容
	t, _ := template.ParseFiles("static/oa.html")
	return func(w http.ResponseWriter, r *http.Request) {
		url := "http://" + r.Host + r.RequestURI
		corpConfig, err := a.GetCorpJSApiConfig(r.Context(), url)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}
		agentConfig, err := a.GetAgentJSApiConfig(r.Context(), url)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}

		val := struct {
			Corp    *agent.JsApiCorpConfig
			Agent   *agent.JsApiAgentConfig
			TmplID  string
			ThirdNo string
		}{
			Corp:    corpConfig,
			Agent:   agentConfig,
			TmplID:  test.AgentOATmplID,
			ThirdNo: test.AgentOAThirdNo,
		}
		// t.Execute(os.Stdout, val)
		err = t.Execute(w, val)
		if err != nil {
			httpAbort(w, http.StatusBadRequest)
			return
		}
	}
}
