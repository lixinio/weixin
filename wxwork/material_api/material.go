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
	"context"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/lixinio/weixin/utils"
)

const (
	apiUpload    = "/cgi-bin/media/upload"
	apiUploadImg = "/cgi-bin/media/uploadimg"
	apiGet       = "/cgi-bin/media/get"
	apiJssdk     = "/cgi-bin/media/get/jssdk"
	// 上传附件资源
	apiUploadTempMedia = "/cgi-bin/media/upload_attachment"
)

const (
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeFile  = "file"
)

// 附件类型，不同的附件类型用于不同的场景。1：朋友圈；2:商品图册
type AttachmentType int

const (
	AttachmentTypeMoment AttachmentType = 1
	AttachmentTypeGoods  AttachmentType = 2
)

type MaterialApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *MaterialApi {
	return &MaterialApi{Client: client}
}

type MaterialID struct {
	utils.WeixinError
	MediaID   string `json:"media_id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

/*
上传临时素材
See: https://work.weixin.qq.com/api/doc/90000/90135/90253
POST(@media) https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=ACCESS_TOKEN&type=TYPE
*/
func (api *MaterialApi) Upload(
	ctx context.Context,
	filename string,
	content io.Reader,
	mediaType string,
) (result *MaterialID, err error) {
	result = &MaterialID{}
	if err := api.Client.HttpFile(
		ctx, apiUpload, "media", filename, content, func(params url.Values) {
			params.Add("type", mediaType)
		}, result,
	); err != nil {
		return nil, err
	}
	return result, nil
}

type MaterialUrl struct {
	utils.WeixinError
	URL string `json:"url"`
}

/*
上传图片
See: https://work.weixin.qq.com/api/doc/90000/90135/90256
POST(@media) https://qyapi.weixin.qq.com/cgi-bin/media/uploadimg?access_token=ACCESS_TOKEN
*/
func (api *MaterialApi) UploadImg(
	ctx context.Context,
	filename string,
	content io.Reader,
) (url string, err error) {
	result := &MaterialUrl{}
	if err := api.Client.HttpFile(
		ctx, apiUploadImg, "media", filename, content, nil, result,
	); err != nil {
		return "", err
	}
	return result.URL, nil
}

/*
获取临时素材
See: https://work.weixin.qq.com/api/doc/90000/90135/90254
GET https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=ACCESS_TOKEN&media_id=MEDIA_ID
*/
func (api *MaterialApi) Get(ctx context.Context, mediaID string) ([]byte, error) {
	resp, err := api.Client.HTTPGetRaw(ctx, apiGet, func(params url.Values) {
		params.Add("media_id", mediaID)
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *MaterialApi) Save(ctx context.Context, mediaID string, saver io.Writer) error {
	resp, err := api.Client.HTTPGetRaw(ctx, apiGet, func(params url.Values) {
		params.Add("media_id", mediaID)
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(saver, resp.Body)
	if err != nil {
		return err
	}
	return nil
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

type AttachmentID struct {
	utils.WeixinError
	MediaID   string `json:"media_id"`
	Type      string `json:"type"`
	CreatedAt int    `json:"created_at"`
}

// 上传附件资源
func (api *MaterialApi) UploadAttachment(
	ctx context.Context,
	filename string,
	content io.Reader,
	mediaType string,
	aType AttachmentType,
) (result *AttachmentID, err error) {
	result = &AttachmentID{}
	if err := api.Client.HttpFile(
		ctx, apiUploadTempMedia, "media", filename, content, func(params url.Values) {
			params.Add("media_type", mediaType)
			params.Add("attachment_type", strconv.Itoa(int(aType)))
		}, result,
	); err != nil {
		return nil, err
	}

	return result, nil
}
