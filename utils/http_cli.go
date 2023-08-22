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
	"reflect"
	"strings"
)

const (
	defaultTokenKey = "access_token"   // 默认的access token的参数名称
	userAgent       = "lixinio/weixin" // 自定义user agent
)

type ClientAccessTokenGetter interface {
	GetAccessToken() (string, error)
}

type EmptyClientAccessTokenGetter int

type MultipartWriter func(writer *multipart.Writer) error

func (EmptyClientAccessTokenGetter) GetAccessToken() (string, error) {
	return "", errors.New("can NOT get token from empty client access-token getter")
}

type StaticClientAccessTokenGetter string

func (s StaticClientAccessTokenGetter) GetAccessToken() (string, error) {
	return string(s), nil
}

/*
HttpClient 用于向微信接口发送请求
*/
type Client struct {
	serverUrl         string
	userAgent         string
	accessTokenKey    string
	accessTokenGetter ClientAccessTokenGetter
}

func NewClient(serverUrl string, accessTokenGetter ClientAccessTokenGetter) *Client {
	return &Client{
		serverUrl:         serverUrl,
		userAgent:         userAgent,
		accessTokenKey:    defaultTokenKey,
		accessTokenGetter: accessTokenGetter,
	}
}

func (client *Client) UpdateAccessTokenKey(accessTokenKey string) {
	client.accessTokenKey = accessTokenKey
}

// HTTPGet GET 请求
func (client *Client) HTTPGet(
	ctx context.Context, path string, result interface{},
) (err error) {
	return client.HTTPGetWithParams(ctx, path, nil, result)
}

// HTTPGetWithParams GET 请求， 支持query参数
func (client *Client) HTTPGetWithParams(
	ctx context.Context, path string, querysFunc func(url.Values), result interface{},
) (err error) {
	newPath, err := client.applyAccessToken(path, querysFunc, true)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newPath, nil)
	if err != nil {
		return
	}

	return client.httpDo(ctx, req, result)
}

// 用来 刷新Token 等不用access-token的接口
func (client *Client) HTTPGetToken(
	ctx context.Context, path string, querysFunc func(url.Values), result interface{},
) (err error) {
	newPath, err := client.applyAccessToken(path, querysFunc, false)
	if err != nil {
		return
	}

	// 调用http请求
	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newPath, nil)
	if err != nil {
		return
	}

	resp, err := client.httpDoRaw(ctx, req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if err = doWeixinError(resp, result); err != nil {
		return err
	}

	return nil
}

// 素材下载， 需要根据Content-Type来判断Body， 可以是json，可能是二进制
// HTTPGetRaw 素材下载， 需要根据Content-Type来判断Body， 可以是json，可能是二进制
func (client *Client) HTTPGetRaw(
	ctx context.Context, path string, querysFunc func(url.Values),
) (resp *http.Response, err error) {
	newPath, err := client.applyAccessToken(path, querysFunc, true)
	if err != nil {
		return
	}

	// 调用http请求
	req, err := http.NewRequest(http.MethodGet, client.serverUrl+newPath, nil)
	if err != nil {
		return
	}

	resp, err = client.httpDoRaw(ctx, req)
	if err != nil {
		return nil, err
	}

	// 如果Content-Type 是 Json, 那出错了
	if hasTextContentType(resp) {
		defer resp.Body.Close()
		result := &WeixinError{}
		if err = doWeixinError(resp, result); err != nil {
			return nil, err
		} else {
			// wtf
			panic(fmt.Errorf(
				"request (%s) response invalid json response(%d: %s)",
				req.URL.Path, result.ErrCode, result.ErrMsg,
			))
		}
	}

	return resp, nil
}

// 生成二维码， 需要根据Content-Type来判断Body， 可以是json，可能是二进制
// HTTPGetRaw 素材下载， 需要根据Content-Type来判断Body， 可以是json，可能是二进制
// 例如 https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func (client *Client) HTTPPostDownload(
	ctx context.Context, path string,
	body interface{}, querysFunc func(url.Values),
) (resp *http.Response, err error) {
	newPath, err := client.applyAccessToken(path, querysFunc, true)
	if err != nil {
		return
	}

	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(body)
	if err != nil {
		return nil, err
	}

	// 调用http请求
	req, err := http.NewRequest(http.MethodPost, client.serverUrl+newPath, payload)
	if err != nil {
		return
	}

	resp, err = client.httpDoRaw(ctx, req)
	if err != nil {
		return nil, err
	}

	// 如果Content-Type 是 Json, 那出错了
	if hasTextContentType(resp) {
		defer resp.Body.Close()
		result := &WeixinError{}
		if err = doWeixinError(resp, result); err != nil {
			return nil, err
		} else {
			// wtf
			panic(fmt.Errorf(
				"request (%s) response invalid json response(%d: %s)",
				req.URL.Path, result.ErrCode, result.ErrMsg,
			))
		}
	}

	return resp, nil
}

// HTTPPost POST 请求, 一次性上传, 优先使用 HttpFile
// 发票上传接口不支持分块(Go Http Client库缺省的方式)
// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Transfer-Encoding
func (client *Client) HTTPUpload(
	ctx context.Context, uri string, payload io.Reader,
	key, filename string, length int64,
	querysFunc func(url.Values), result interface{},
	multipartWriters ...MultipartWriter,
) error {
	// 头部大小
	bodyBuffer := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bodyBuffer)

	for _, multipartWriter := range multipartWriters {
		if err := multipartWriter(bodyWriter); err != nil {
			return err
		}
	}

	_, err := bodyWriter.CreateFormFile(key, path.Base(filename))
	if err != nil {
		return err
	}
	// 尾部
	closeBuffer := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", bodyWriter.Boundary()))

	newUrl, err := client.applyAccessToken(uri, querysFunc, true)
	if err != nil {
		return err
	}

	reader := io.MultiReader(bodyBuffer, payload, closeBuffer)
	req, err := http.NewRequest(http.MethodPost, client.serverUrl+newUrl, ioutil.NopCloser(reader))
	if err != nil {
		return err
	}
	req.TransferEncoding = []string{"identity"}
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	req.ContentLength = length + int64(closeBuffer.Len()) + int64(bodyBuffer.Len())

	return client.httpDo(ctx, req, result)
}

// Upload 上传文件
// HttpFile 上传文件, 适合没有什么定制的文件上传
func (client *Client) HttpFile(
	ctx context.Context, path, key, filename string,
	content io.Reader, querysFunc func(url.Values), result interface{},
	multipartWriters ...MultipartWriter,
) (err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		for _, multipartWriter := range multipartWriters {
			if err = multipartWriter(m); err != nil {
				return
			}
		}

		part, err := m.CreateFormFile(key, filepath.Base(filename))
		if err != nil {
			return
		}
		if _, err = io.Copy(part, content); err != nil {
			return
		}
	}()

	return client.HTTPPostRaw(ctx, path, r, querysFunc, result, m.FormDataContentType(), true)
}

//HTTPPost POST 请求
func (client *Client) HTTPPostJson(
	ctx context.Context, path string, body interface{}, result interface{},
) (err error) {
	return client.HTTPPost(ctx, path, body, nil, result, "")
}

//HTTPPost POST 请求(json, 文件上传)
func (client *Client) HTTPPost(
	ctx context.Context, path string, body interface{},
	querysFunc func(url.Values), result interface{}, contentType string,
) (err error) {
	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(body)
	if err != nil {
		return err
	}

	return client.HTTPPostRaw(ctx, path, payload, querysFunc, result, contentType, true)
}

//HTTPPost POST 请求(无需access-token认证)
func (client *Client) HTTPPostToken(
	ctx context.Context, path string, body interface{}, result interface{},
) (err error) {
	payload := new(bytes.Buffer)
	err = json.NewEncoder(payload).Encode(body)
	if err != nil {
		return err
	}

	return client.HTTPPostRaw(ctx, path, payload, nil, result, "", false)
}

//HTTPPostRaw POST 请求, 不做内容的序列化， 适合特殊的文件上传
func (client *Client) HTTPPostRaw(
	ctx context.Context, path string, payload io.Reader, querysFunc func(url.Values),
	result interface{}, contentType string, auth bool,
) (err error) {
	newPath, err := client.applyAccessToken(path, querysFunc, auth)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, client.serverUrl+newPath, payload)
	if err != nil {
		return
	}

	if contentType == "" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
	} else {
		req.Header.Add("Content-Type", contentType)
	}
	return client.httpDo(ctx, req, result)
}

//httpDo httpDoRaw加上结果反序列化， 适合返回json的普通请求
func (client *Client) httpDo(
	ctx context.Context, req *http.Request, result interface{},
) (err error) {
	response, err := client.httpDoRaw(ctx, req)
	if err != nil {
		return
	}

	defer response.Body.Close()
	weixinResult := result
	if result == nil {
		// 如果上层并不关心实际的响应, 就简单的判断腾讯的Code
		weixinResult = &WeixinError{}
	}

	if err = doWeixinError(response, weixinResult); err != nil {
		return err
	}

	return nil
}

func hasTextContentType(resp *http.Response) bool {
	ct := resp.Header.Get("Content-Type")
	if len(ct) > 0 {
		// https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_temporary_materials.html
		// 素材管理 /获取临时素材 是  "text/plain"
		return strings.HasPrefix(ct, "application/json") ||
			strings.HasPrefix(ct, "text/plain")
	}
	return false
}

func doWeixinError(response *http.Response, result interface{}) error {
	// 直接从body反序列化， 无需先读取到内存
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return err
	}

	we, ok := result.(WeixinErrorInterface)
	if !ok {
		panic(fmt.Errorf(
			"request payload (%s) not implement weixin error interface",
			reflect.TypeOf(result).String(),
		))
	}

	wxCode := we.WeixinErrorCode()
	if wxCode == 0 {
		return nil
	}

	if wxCode == 40014 {
		// 不合法的access_token
		// https://open.work.weixin.qq.com/devtool/query?e=40014
		return fmt.Errorf("error get token %s, error %w", we.WeixinErrorMessage(), ErrorAccessToken)
	} else if wxCode == -1 {
		//  -1	系统繁忙，服务器暂不可用，建议稍候重试。建议重试次数不超过3次。
		// https://open.work.weixin.qq.com/devtool/query?e=40014
		return fmt.Errorf("error get token %s, error %w", we.WeixinErrorMessage(), ErrorSystemBusy)
	} else {
		return we.GetWeixinError()
	}
}

//httpDoRaw 执行具体的请求发送， 处理认证， user-agent, trace, 判断http code等细节
// 不做结果反序列化， 考虑文件下载
func (client *Client) httpDoRaw(
	ctx context.Context, req *http.Request,
) (resp *http.Response, err error) {
	req.Header.Add("User-Agent", client.userAgent)

	cli := http.DefaultClient
	if ctx != context.TODO() {
		req = req.WithContext(ctx)
		cli = &http.Client{Transport: newTransport()}
	}

	resp, err = cli.Do(req)
	if err != nil {
		return nil, err
	}

	// 根据规范，有些接口返回20x，这里暂不考虑
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		resp.Body.Close()
		err = fmt.Errorf("status %s", resp.Status)
		resp = nil
		return
	}

	return resp, nil
}

/*
在请求地址上附加上 access_token
*/
func (client *Client) applyAccessToken(
	oldUrl string, querysFunc func(url.Values), auth bool,
) (newUrl string, err error) {
	querys := url.Values{}
	// 客户自定义
	if querysFunc != nil {
		querysFunc(querys)
	}

	// 认证
	if auth {
		accessToken, err := client.accessTokenGetter.GetAccessToken()
		if err != nil {
			return "", err
		}

		querys.Add(client.accessTokenKey, accessToken)
	} else if len(querys) == 0 {
		return oldUrl, nil
	}

	if strings.Contains(oldUrl, "?") {
		newUrl = oldUrl + "&" + querys.Encode()
	} else {
		newUrl = oldUrl + "?" + querys.Encode()
	}
	return
}
