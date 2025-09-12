package material_api

// Package material 素材管理

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/lixinio/weixin/utils"
)

const (
	apiUpload         = "/cgi-bin/media/upload"               // 新增临时素材
	apiGet            = "/cgi-bin/media/get"                  // 获取临时素材
	apiUploadImg      = "/cgi-bin/media/uploadimg"            // 新增永久素材(图片)
	apiUploadMaterial = "/cgi-bin/material/add_material"      // 新增永久素材
	apiGetMaterial    = "/cgi-bin/material/get_material"      // 获取永久素材
	apiDeleteMaterial = "/cgi-bin/material/del_material"      // 删除永久素材
	apiCountMaterial  = "/cgi-bin/material/get_materialcount" // 获取素材总数
	apiListMaterial   = "/cgi-bin/material/batchget_material" // 获取素材列表
	apiJssdk          = "/cgi-bin/media/get/jssdk"            // 高清语音素材获取接口
	apiAddNews        = "/cgi-bin/material/add_news"          // 新增永久素材(news)
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVoice MediaType = "voice"
	MediaTypeVideo MediaType = "video"
	MediaTypeThumb MediaType = "thumb"
	MediaTypeNews  MediaType = "news"
)

type MaterialApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *MaterialApi {
	return &MaterialApi{Client: client}
}

/*
上传临时素材
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
*/
type MediaID struct {
	utils.WeixinError
	MediaID   string `json:"media_id"`
	Type      string `json:"type"`
	CreatedAt int64  `json:"created_at"`
}

func (api *MaterialApi) UploadMedia(
	ctx context.Context,
	filename string,
	length int64,
	content io.Reader,
	mediaType MediaType,
) (result *MediaID, err error) {
	result = &MediaID{}
	if err := api.Client.HTTPUpload(
		ctx, apiUpload, content, "media", filename, length, func(params url.Values) {
			params.Add("type", string(mediaType))
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
获取临时素材
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_temporary_materials.html
*/
func (api *MaterialApi) GetMedia(ctx context.Context, mediaID string) ([]byte, error) {
	resp, err := api.Client.HTTPGetRaw(ctx, apiGet, func(params url.Values) {
		params.Add("media_id", mediaID)
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *MaterialApi) SaveMedia(ctx context.Context, mediaID string, saver io.Writer) error {
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

/*
新增永久素材(图片)
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Adding_Permanent_Assets.html
*/
func (api *MaterialApi) UploadImg(
	ctx context.Context,
	filename string,
	length int64,
	content io.Reader,
) (url string, err error) {
	result := &MaterialUrl{}
	if err := api.Client.HTTPUpload(
		ctx, apiUploadImg, content, "media", filename, length, nil, result,
	); err != nil {
		return "", err
	}

	return result.URL, nil
}

/*
新增永久素材
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/New_temporary_materials.html
*/
type MaterialID struct {
	utils.WeixinError
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}

func (api *MaterialApi) UploadMaterial(
	ctx context.Context,
	filename string,
	length int64,
	content io.Reader,
	mediaType MediaType,
) (result *MaterialID, err error) {
	result = &MaterialID{}
	if err := api.Client.HTTPUpload(
		ctx, apiUploadMaterial, content, "media", filename, length, func(params url.Values) {
			params.Add("type", string(mediaType))
		}, result,
	); err != nil {
		return nil, err
	}

	return result, nil
}

func (api *MaterialApi) UploadVideoMaterial(
	ctx context.Context,
	filename, title, introduction string,
	length int64,
	content io.Reader,
) (result *MaterialID, err error) {
	result = &MaterialID{}
	multipartWriter := func(writer *multipart.Writer) error {
		fw, err := writer.CreateFormField("description")
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(fw)
		if err = encoder.Encode(map[string]string{
			"title":        title,
			"introduction": introduction,
		}); err != nil {
			return err
		}

		return nil
	}

	// multipartWriter,
	if err = api.Client.HTTPUpload(
		ctx, apiUploadMaterial, content, "media", filename, length, func(params url.Values) {
			params.Add("type", string(MediaTypeVideo))
		}, result, multipartWriter,
	); err != nil {
		return nil, err
	}

	return result, nil
}

/*
获取永久素材
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Getting_Permanent_Assets.html
*/
func (api *MaterialApi) GetMaterial(ctx context.Context, mediaID string) ([]byte, error) {
	resp, err := api.Client.HTTPPostDownload(ctx, apiGetMaterial, map[string]string{
		"media_id": mediaID,
	}, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (api *MaterialApi) SaveMaterial(ctx context.Context, mediaID string, saver io.Writer) error {
	resp, err := api.Client.HTTPPostDownload(ctx, apiGetMaterial, map[string]string{
		"media_id": mediaID,
	}, nil)
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

// 下载 视频消息素材：
type VideoMaterial struct {
	utils.WeixinError
	Title       string `json:"title"`
	Description string `json:"description"`
	DownloadUrl string `json:"down_url"`
}

func (api *MaterialApi) GetVideoMaterial(
	ctx context.Context, mediaID string,
) (*VideoMaterial, error) {
	resp := &VideoMaterial{}
	if err := api.Client.HTTPPostJson(ctx, apiGetMaterial, map[string]string{
		"media_id": mediaID,
	}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// https://developers.weixin.qq.com/doc/offiaccount/Comments_management/Image_Comments_Management_Interface.html
// 下载 图文消息素材
type MpnewsMaterialItem struct {
	Title              string `json:"title"`
	ThumbMediaID       string `json:"thumb_media_id"`
	ShowCoverPic       int    `json:"show_cover_pic"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	Content            string `json:"content"`
	URL                string `json:"url"`
	ContentSourceURL   string `json:"content_source_url"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}

func (api *MaterialApi) GetMpnewsMaterial(
	ctx context.Context, mediaID string,
) ([]*MpnewsMaterialItem, error) {
	resp := &struct {
		utils.WeixinError
		Items []*MpnewsMaterialItem `json:"news_item"`
	}{}

	if err := api.Client.HTTPPostJson(ctx, apiGetMaterial, map[string]string{
		"media_id": mediaID,
	}, resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

/*
删除永久素材
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Deleting_Permanent_Assets.html
*/
func (api *MaterialApi) DeleteMaterial(
	ctx context.Context, mediaID string,
) error {
	if err := api.Client.HTTPPostJson(ctx, apiDeleteMaterial, map[string]string{
		"media_id": mediaID,
	}, nil); err != nil {
		return err
	}
	return nil
}

/*
获取素材总数
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_the_total_of_all_materials.html
*/
type MaterialStatistics struct {
	utils.WeixinError
	VoiceCount int `json:"voice_count"`
	VideoCount int `json:"video_count"`
	ImageCount int `json:"image_count"`
	NewsCount  int `json:"news_count"`
}

func (api *MaterialApi) GetMaterialStatistics(
	ctx context.Context,
) (*MaterialStatistics, error) {
	resp := &MaterialStatistics{}
	if err := api.Client.HTTPGet(ctx, apiCountMaterial, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

/*
获取素材列表
See: https://developers.weixin.qq.com/doc/offiaccount/Asset_Management/Get_materials_list.html
*/
type Material struct {
	MediaID    string `json:"media_id"`
	Name       string `json:"name"`
	UpdateTime int64  `json:"update_time"`
	URL        string `json:"url"`
}

type MaterialList struct {
	utils.WeixinError
	TotalCount int         `json:"total_count"`
	ItemCount  int         `json:"item_count"`
	Items      []*Material `json:"item"`
}

func (api *MaterialApi) ListMaterial(
	ctx context.Context, tp MediaType, offset, count int,
) (*MaterialList, error) {
	resp := &MaterialList{}
	if err := api.Client.HTTPPostJson(ctx, apiListMaterial, map[string]interface{}{
		"type":   string(tp),
		"offset": offset,
		"count":  count,
	}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// 永久图文消息素材列表
type MpnewsMaterial struct {
	MediaID string `json:"media_id"`
	Content struct {
		NewsItems []*MpnewsMaterialItem `json:"news_item"`
	} `json:"content"`
}

type MpnewsMaterialList struct {
	utils.WeixinError
	TotalCount int               `json:"total_count"`
	ItemCount  int               `json:"item_count"`
	Items      []*MpnewsMaterial `json:"item"`
}

func (api *MaterialApi) ListMpnewsMaterial(
	ctx context.Context, offset, count int,
) (*MpnewsMaterialList, error) {
	resp := &MpnewsMaterialList{}
	if err := api.Client.HTTPPostJson(ctx, apiListMaterial, map[string]interface{}{
		"type":   string(MediaTypeNews),
		"offset": offset,
		"count":  count,
	}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (api *MaterialApi) AddMpnewsMaterial(
	ctx context.Context, articles []*MpnewsMaterialItem,
) (string, error) {
	resp := &struct {
		utils.WeixinError
		MediaID string `json:"media_id"`
	}{}

	if err := api.Client.HTTPPostJson(ctx, apiAddNews, map[string]interface{}{
		"articles": articles,
	}, resp); err != nil {
		return "", err
	}
	return resp.MediaID, nil
}
