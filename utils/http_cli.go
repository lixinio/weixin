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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

var (
	ErrorAccessToken = errors.New("access token error")
	ErrorSystemBusy  = errors.New("system busy")
	UserAgent        = "lixinio/weixin"
)

type WeixinError struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (we WeixinError) Error() string {
	return we.Errmsg
}

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
func (client *Client) HTTPGet(ctx context.Context, uri string) (resp []byte, err error) {
	return client.HTTPGetWithParams(ctx, uri, url.Values{})
}

func (client *Client) HTTPGetWithParams(ctx context.Context, uri string, params url.Values) (resp []byte, err error) {
	newUrl, err := client.applyAccessToken(uri, params)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newUrl, nil)
	if err != nil {
		return
	}

	return client.httpDo(req.WithContext(ctx))
}

// 素材下载， 需要根据Content-Type来判断Body， 可以是json，可能是二进制
func (client *Client) HTTPGetWithParamsRaw(ctx context.Context, uri string, params url.Values) (resp *http.Response, err error) {
	newUrl, err := client.applyAccessToken(uri, params)
	if err != nil {
		return
	}

	// 调用http请求
	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newUrl, nil)
	if err != nil {
		return
	}

	cli := &http.Client{Transport: newTransport()}
	req.Header.Add("User-Agent", client.userAgent)
	return cli.Do(req.WithContext(ctx))
}

//HTTPPost POST 请求
func (client *Client) HTTPUpload(
	ctx context.Context,
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

	return client.httpDo(req.WithContext(ctx))
}

// Upload 上传文件
func (client *Client) Upload(ctx context.Context, uri string, key string, filename string, content io.Reader) (resp []byte, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(key, filepath.Base(filename))
	if err != nil {
		return
	}

	if _, err = io.Copy(part, content); err != nil {
		return
	}
	writer.Close()

	return client.HTTPPost(ctx, uri, body, writer.FormDataContentType())
}

//HTTPPost POST 请求
func (client *Client) HTTPPost(
	ctx context.Context, uri string, payload io.Reader, contentType string,
) (resp []byte, err error) {
	newUrl, err := client.applyAccessToken(uri, url.Values{})
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, client.serverUrl+newUrl, payload)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", contentType)

	return client.httpDo(req.WithContext(ctx))
}

//httpDo 执行 请求
func (client *Client) httpDo(req *http.Request) (resp []byte, err error) {
	req.Header.Add("User-Agent", client.userAgent)

	cli := &http.Client{Transport: newTransport()}
	response, err := cli.Do(req)
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
		err = fmt.Errorf("status %s", response.Status)
		return
	}

	resp, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	errorResponse := WeixinError{}
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
		err = errorResponse
		return
	}

	return
}

/// 工具方法
func (client *Client) ApiGetWrapper(ctx context.Context, urlPath string, paramFunc func(url.Values), result interface{}) error {
	params := url.Values{}
	paramFunc(params)
	resp, err := client.HTTPGetWithParams(ctx, urlPath, params)
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}

func (client *Client) ApiGetNullWrapper(ctx context.Context, urlPath string, result interface{}) error {
	resp, err := client.HTTPGet(ctx, urlPath)
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}

func (client *Client) ApiPostWrapper(ctx context.Context, urlPath string, payload interface{}, result interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.HTTPPost(ctx, urlPath, bytes.NewReader(body), "application/json;charset=utf-8")
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}

func (client *Client) ApiPostWrapperEx(
	context context.Context,
	urlPath string, obj interface{},
	paramFunc func(url.Values), result interface{},
) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	params := url.Values{}
	paramFunc(params)
	resp, err := client.HTTPPost(
		context, urlPath+"?"+params.Encode(),
		bytes.NewReader(body), "application/json;charset=utf-8",
	)
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}
