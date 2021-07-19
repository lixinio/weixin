package open

// https://github.com/fastwego/wxopen/blob/master/wxopen.go

// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
微信开放平台 SDK
See: https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/Third_party_platform_appid.html
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
从微信服务器获取新的 ComponentAccessToken

See: https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
*/
func (open *Open) refreshAccessTokenFromWXServer() (accessToken string, expiresIn int, err error) {
	params := map[string]string{
		"component_appid":         open.Config.ComponentAppid,
		"component_appsecret":     open.Config.ComponentSecret,
		"component_verify_ticket": open.component_verify_ticket_getter(open.Config.ComponentAppid),
	}
	payload, err := json.Marshal(params)
	if err != nil {
		return
	}

	/**
	POST 数据示例：
	{
	  "component_appid":  "appid_value" ,
	  "component_appsecret":  "appsecret_value",
	  "component_verify_ticket": "ticket_value"
	}
	*/
	url := WXServerUrl + "/cgi-bin/component/api_component_token"

	response, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(payload))
	if err != nil {
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s RETURN %s", url, response.Status)
		return
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	/**
	返回结果示例：
	{
	  "component_access_token": "61W3mEpU66027wgNZ_MhGHNQDHnFATkDa9-2llqrMBjUwxRSNPbVsMmyD-yq8wZETSoE5NQgecigDrSHkPtIYA",
	  "expires_in": 7200
	}
	*/
	var result = struct {
		AccessToken string  `json:"component_access_token"`
		ExpiresIn   int     `json:"expires_in"`
		Errcode     float64 `json:"errcode"`
		Errmsg      string  `json:"errmsg"`
	}{}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		err = fmt.Errorf("Unmarshal error %s", string(resp))
		return
	}

	if result.AccessToken == "" {
		err = fmt.Errorf("%s", string(resp))
		return
	}

	return result.AccessToken, result.ExpiresIn, nil
}
