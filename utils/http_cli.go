package utils

// https://github.com/fastwego/wxwork/blob/master/corporation/client.go
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

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	ErrorAccessToken = errors.New("access token error")
	ErrorSystemBusy  = errors.New("system busy")
	UserAgent        = "lixinio/weixin"
)

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func ContentType(resp *http.Response) string {
	ct := resp.Header.Get("Content-Type")
	if len(ct) > 0 {
		return filterFlags(ct)
	}
	return ""
}

/*
HttpClient 用于向微信接口发送请求
*/
type Client struct {
	serverUrl        string
	userAgent        string
	accessTokenCache *AccessTokenCache
}

func NewClient(serverUrl string, accessTokenCache *AccessTokenCache) *Client {
	return &Client{
		serverUrl:        serverUrl,
		userAgent:        UserAgent,
		accessTokenCache: accessTokenCache,
	}
}

// HTTPGet GET 请求
func (client *Client) HTTPGet(uri string) (resp []byte, err error) {
	return client.HTTPGetWithParams(uri, url.Values{})
}

func (client *Client) HTTPGetWithParams(uri string, params url.Values) (resp []byte, err error) {
	newUrl, err := client.applyAccessToken(uri, params)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newUrl, nil)
	if err != nil {
		return
	}

	return client.httpDo(req)
}

// 素材下载， 需要根据Content-Type来判断Body， 可以是json，可能是二进制
func (client *Client) HTTPGetWithParamsRaw(uri string, params url.Values) (resp *http.Response, err error) {
	newUrl, err := client.applyAccessToken(uri, params)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newUrl, nil)
	req.Header.Add("User-Agent", client.userAgent)
	return http.DefaultClient.Do(req)
}

//HTTPPost POST 请求
func (client *Client) HTTPUpload(
	uri string,
	payload io.Reader,
	key, filename string,
	length int64,
) (resp []byte, err error) {
	// 头部大小
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)
	_, err = bodyWriter.CreateFormFile(key, path.Base(filename))
	if err != nil {
		return
	}
	// 尾部
	closeBuffer := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", bodyWriter.Boundary()))

	newUrl, err := client.applyAccessToken(uri, url.Values{})
	if err != nil {
		return
	}

	reader := io.MultiReader(bodyBuffer, payload, closeBuffer)
	req, err := http.NewRequest(http.MethodPost, client.serverUrl+newUrl, ioutil.NopCloser(reader))
	if err != nil {
		return
	}
	req.TransferEncoding = []string{"identity"}
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	req.ContentLength = length + int64(closeBuffer.Len()) + int64(bodyBuffer.Len())

	return client.httpDo(req)
}

//HTTPPost POST 请求
func (client *Client) HTTPPost(uri string, payload io.Reader, contentType string) (resp []byte, err error) {
	buf := &bytes.Buffer{}
	tee := io.TeeReader(payload, buf)
	reader1, _ := ioutil.ReadAll(tee)
	fmt.Println("debug: ", string(reader1))

	newUrl, err := client.applyAccessToken(uri, url.Values{})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, client.serverUrl+newUrl, buf)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", contentType)

	return client.httpDo(req)
}

//httpDo 执行 请求
func (client *Client) httpDo(req *http.Request) (resp []byte, err error) {
	req.Header.Add("User-Agent", client.userAgent)

	// if client.Ctx.Corporation.Logger != nil {
	// 	client.Ctx.Corporation.Logger.Printf("%s %s Headers %v", req.Method, req.URL.String(), req.Header)
	// }

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	resp, err = ResponseFilter(response)

	// 发现 access_token 过期
	if err == ErrorAccessToken {
		// 通知到位后 access_token 会被刷新，那么可以 retry 了
		var accessToken string
		accessToken, err = client.accessTokenCache.GetAccessToken()
		if err != nil {
			return
		}

		// 换新
		q := req.URL.Query()
		q.Set("access_token", accessToken)
		req.URL.RawQuery = q.Encode()

		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer response.Body.Close()

		resp, err = ResponseFilter(response)
	}

	// -1 系统繁忙，此时请开发者稍候再试
	// 重试一次
	if err == ErrorSystemBusy {

		// if client.Ctx.Corporation.Logger != nil {
		// 	client.Ctx.Corporation.Logger.Printf("%v : retry %s %s Headers %v", ErrorSystemBusy, req.Method, req.URL.String(), req.Header)
		// }

		response, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		defer response.Body.Close()

		resp, err = ResponseFilter(response)
	}

	return
}

/*
在请求地址上附加上 access_token
*/
func (client *Client) applyAccessToken(oldUrl string, params url.Values) (newUrl string, err error) {
	accessToken, err := client.accessTokenCache.GetAccessToken()
	if err != nil {
		return
	}
	params.Add("access_token", accessToken)
	if strings.Contains(oldUrl, "?") {
		newUrl = oldUrl + "&" + params.Encode()
	} else {
		newUrl = oldUrl + "?" + params.Encode()
	}
	return
}

/*
筛查微信 api 服务器响应，判断以下错误：

- http 状态码 不为 200

- 接口响应错误码 errcode 不为 0
*/
func ResponseFilter(response *http.Response) (resp []byte, err error) {
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("Status %s", response.Status)
		return
	}

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	errorResponse := struct {
		Errcode int64  `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	err = json.Unmarshal(resp, &errorResponse)
	if err != nil {
		return
	}

	if errorResponse.Errcode == 40014 {
		err = ErrorAccessToken
		return
	}

	//  -1	系统繁忙，此时请开发者稍候再试
	if errorResponse.Errcode == -1 {
		err = ErrorSystemBusy
		return
	}

	if errorResponse.Errcode != 0 {
		err = errors.New(string(resp))
		return
	}
	return
}
