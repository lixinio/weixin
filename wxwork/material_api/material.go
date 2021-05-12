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
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"

	"github.com/lixinio/weixin"
	"github.com/lixinio/weixin/wxwork/agent"
)

const (
	apiUpload    = "/cgi-bin/media/upload"
	apiUploadImg = "/cgi-bin/media/uploadimg"
	apiGet       = "/cgi-bin/media/get"
	apiJssdk     = "/cgi-bin/media/get/jssdk"
)

type MaterialApi struct {
	*weixin.Client
}

func NewAgentApi(agent *agent.Agent) *MaterialApi {
	return &MaterialApi{
		Client: agent.Client,
	}
}

/*
上传临时素材
See: https://work.weixin.qq.com/api/doc/90000/90135/90253
POST(@media) https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=ACCESS_TOKEN&type=TYPE
*/
func (api *MaterialApi) Upload(media string, params url.Values) (resp []byte, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("media", path.Base(media))
		if err != nil {
			return
		}
		file, err := os.Open(media)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}

	}()
	return api.Client.HTTPPost(apiUpload+"?"+params.Encode(), r, m.FormDataContentType())
}

/*
上传图片
See: https://work.weixin.qq.com/api/doc/90000/90135/90256
POST(@media) https://qyapi.weixin.qq.com/cgi-bin/media/uploadimg?access_token=ACCESS_TOKEN
*/
func (api *MaterialApi) UploadImg(media string) (resp []byte, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("media", path.Base(media))
		if err != nil {
			return
		}
		file, err := os.Open(media)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}

	}()
	return api.Client.HTTPPost(apiUploadImg, r, m.FormDataContentType())
}

/*
获取临时素材
See: https://work.weixin.qq.com/api/doc/90000/90135/90254
GET https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=ACCESS_TOKEN&media_id=MEDIA_ID
*/
// func (api *MaterialApi) Get(params url.Values, header http.Header) (resp *http.Response, err error) {

// 	return api.Client.HTTPGetWithParams(apiGet + "?" + params.Encode())
// 	accessToken, err := ctx.AccessToken.GetAccessTokenHandler(ctx)
// 	if err != nil {
// 		return
// 	}

// 	req, err := http.NewRequest(http.MethodGet, corporation.WXServerUrl+apiGet+"?access_token="+accessToken+"&"+params.Encode(), nil)
// 	if err != nil {
// 		return
// 	}

// 	req.Header = header

// 	return http.DefaultClient.Do(req)
// }

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
