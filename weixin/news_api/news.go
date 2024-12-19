package news_api

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	// 获取已成功发布的图文消息列表
	apiBatchGet = "/cgi-bin/freepublish/batchget"
	// 获取素材列表
	apiBatchGetMaterial = "/cgi-bin/material/batchget_material"
)

type NewsApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *NewsApi {
	return &NewsApi{
		Client: client,
	}
}

type Content struct {
	NewsItems []*struct {
		Title              string `json:"title"`
		Author             string `json:"author"`
		Digest             string `json:"digest"`
		Content            string `json:"content"`
		ContentSourceUrl   string `json:"content_source_url"`
		ThumbMediaId       string `json:"thumb_media_id"`
		ShowCoverPic       int8   `json:"show_cover_pic"`
		NeedOpenComment    int8   `json:"need_open_comment"`
		OnlyFansCanComment int8   `json:"only_fans_can_comment"`
		Url                string `json:"url"`
		IsDeleted          bool   `json:"is_deleted"`
	} `json:"news_item"`
}

/*
获取已成功发布的消息列表
https://developers.weixin.qq.com/doc/offiaccount/Publish/Get_publication_records.html
*/

type BatchGetRequest struct {
	Offset    int  `json:"offset"`
	Count     int  `json:"count"`
	NoContent int8 `json:"no_content"`
}

type BatchGetResponse struct {
	utils.WeixinError
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Items      []*struct {
		ArticleId  string   `json:"article_id"`
		Content    *Content `json:"content"`
		UpdateTime int64    `json:"update_time"`
	} `json:"item"`
}

func (api *NewsApi) BatchGet(
	ctx context.Context,
	offset int,
	count int,
	noContent int8,
) (*BatchGetResponse, error) {
	params := &BatchGetRequest{
		Offset:    offset,
		Count:     count,
		NoContent: noContent,
	}

	result := &BatchGetResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiBatchGet, params, result); err != nil {
		return nil, err
	}

	return result, nil
}

type MaterialType string

const (
	MartialTypeImage = "image" // 图片
	MartialTypeVoice = "voice" // 语音
	MartialTypeVideo = "video" // 视频
	// MartialTypeNews  = "news"  // 图文（现在可能已无法从该接口获取，具体功能转移到了草稿与发布）
)

type BatchGetMaterialRequest struct {
	Type   MaterialType `json:"type"`
	Offset int          `json:"offset"`
	Count  int          `json:"count"`
}

type BatchGetMaterialResponse struct {
	utils.WeixinError
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Items      []*struct {
		MediaId    string `json:"media_id"`
		Name       string `json:"name"`
		UpdateTime int64  `json:"update_time"`
		Url        string `json:"url"`
	} `json:"item"`
}

func (api *NewsApi) BatchGetMaterial(
	ctx context.Context,
	materialType MaterialType,
	offset, count int,
) (*BatchGetMaterialResponse, error) {
	params := &BatchGetMaterialRequest{
		Type:   materialType,
		Offset: offset,
		Count:  count,
	}

	result := &BatchGetMaterialResponse{}
	if err := api.Client.HTTPPostJson(ctx, apiBatchGetMaterial, params, result); err != nil {
		return nil, err
	}

	return result, nil
}
