package material_api

// https://github.com/fastwego/wxwork/blob/master/corporation/apis/material/material.go
// Copyright 2021 FastWeGo
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

// Package material 素材管理

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork/agent"
)

const (
	apiUpload    = "/cgi-bin/media/upload"
	apiUploadImg = "/cgi-bin/media/uploadimg"
	apiGet       = "/cgi-bin/media/get"
	apiJssdk     = "/cgi-bin/media/get/jssdk"
)

const (
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeFile  = "file"
)

type MaterialApi struct {
	*utils.Client
}

func NewAgentApi(agent *agent.Agent) *MaterialApi {
	return &MaterialApi{
		Client: agent.Client,
	}
}

type MaterialID struct {
	utils.CommonError
	MediaID   string `json:"media_id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

/*
上传临时素材
See: https://work.weixin.qq.com/api/doc/90000/90135/90253
POST(@media) https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=ACCESS_TOKEN&type=TYPE
*/
func (api *MaterialApi) Upload(filename string, content io.Reader, mediaType string) (result *MaterialID, err error) {
	params := url.Values{}
	params.Add("type", mediaType)
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("media", path.Base(filename))
		if err != nil {
			return
		}
		if _, err = io.Copy(part, content); err != nil {
			return
		}

	}()

	var resp []byte
	resp, err = api.Client.HTTPPost(apiUpload+"?"+params.Encode(), r, m.FormDataContentType())
	if err != nil {
		return
	}

	result = &MaterialID{}
	err = json.Unmarshal(resp, result)
	if err != nil {
		return
	}
	return
}

type MaterialUrl struct {
	utils.CommonError
	URL string `json:"url"`
}

/*
上传图片
See: https://work.weixin.qq.com/api/doc/90000/90135/90256
POST(@media) https://qyapi.weixin.qq.com/cgi-bin/media/uploadimg?access_token=ACCESS_TOKEN
*/
func (api *MaterialApi) UploadImg(filename string, content io.Reader) (url string, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("media", path.Base(filename))
		if err != nil {
			return
		}

		if _, err = io.Copy(part, content); err != nil {
			return
		}

	}()

	var resp []byte
	resp, err = api.Client.HTTPPost(apiUploadImg, r, m.FormDataContentType())
	if err != nil {
		return
	}
	var result MaterialUrl
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return
	}
	return result.URL, nil
}

/*
获取临时素材
See: https://work.weixin.qq.com/api/doc/90000/90135/90254
GET https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=ACCESS_TOKEN&media_id=MEDIA_ID
*/
func (api *MaterialApi) Get(mediaID string) (resp *http.Response, err error) {
	params := url.Values{}
	params.Add("media_id", mediaID)

	resp, err = api.Client.HTTPGetWithParamsRaw(apiGet, params)
	if err != nil {
		return
	}

	ct := utils.ContentType(resp)
	if ct != "application/json" && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent) {
		return resp, nil
	}

	var body []byte
	body, err = utils.ResponseFilter(resp)
	if err == nil {
		// 不应该走到这里来
		err = fmt.Errorf("unknown error %s(%s)", ct, string(body))
	}
	return

}

// /*
// 获取高清语音素材
// See: https://work.weixin.qq.com/api/doc/90000/90135/90255
// GET https://qyapi.weixin.qq.com/cgi-bin/media/get/jssdk?access_token=ACCESS_TOKEN&media_id=MEDIA_ID
// */
// func (api *MaterialApi) Jssdk(params url.Values) (resp *http.Response, err error) {
// 	accessToken, err := ctx.AccessToken.GetAccessTokenHandler(ctx)
// 	if err != nil {
// 		return
// 	}

// 	req, err := http.NewRequest(http.MethodGet, corporation.WXServerUrl+apiJssdk+"?access_token="+accessToken+"&"+params.Encode(), nil)
// 	if err != nil {
// 		return
// 	}

// 	return http.DefaultClient.Do(req)
// }
