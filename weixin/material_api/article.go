package material_api

// 文章管理

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiBatchGetArticles = "/cgi-bin/freepublish/batchget"   // 获取已发布的消息列表
	apiGetArticle       = "/cgi-bin/freepublish/getarticle" // 获取已发布图文信息
)

type Article struct {
	Title  string `json:"title"`  // 标题
	Author string `json:"author"` // 作者名
	Digest string `json:"digest"` // 图文消息的摘要，仅有单图文消息才有摘要，多图文此处为空。如果本字段为没有填写，则默认抓取正文前54个字。
	// 图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS,
	// 涉及图片url必须来源 "上传图文消息内的图片获取URL"接口获取。
	Content            string `json:"content"`               // 外部图片url将被过滤。 图片消息则仅支持纯文本和部分特殊功能标签如商品，商品个数不可超过50个
	ContentSourceUrl   string `json:"content_source_url"`    // 图文消息的原文地址，即点击“阅读原文”后的URL
	ThumbMediaID       string `json:"thumb_media_id"`        // 图文消息的封面图片素材id（必须是永久MediaID）
	ThumbUrl           string `json:"thumb_url"`             // 图文消息的封面图片URL
	NeedOpenComment    int    `json:"need_open_comment"`     // 是否打开评论，0不打开(默认)，1打开
	OnlyFansCanComment int    `json:"only_fans_can_comment"` // 是否粉丝才可评论，0所有人可评论(默认)，1粉丝才可评论
	Url                string `json:"url"`                   // 草稿的临时链接
	IsDeleted          bool   `json:"is_deleted"`            // 该图文是否被删除
}

type ArticleList struct {
	utils.WeixinError
	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`
	Items      []struct {
		ArticleID  string `json:"article_id"`
		UpdateTime int64  `json:"update_time"`
		Content    struct {
			NewsItems []*Article `json:"news_item"`
		} `json:"content"`
	} `json:"item"`
}

/*
获取已发布的消息列表
https://developers.weixin.qq.com/doc/service/api/public/api_freepublish_batchget.html
*/
func (api *MaterialApi) BatchGetArticles(
	ctx context.Context,
	offset, count int,
	noContent bool, // 1 表示不返回content字段，0表示正常返回，默认为0
) (*ArticleList, error) {
	noContentInt := 0
	if noContent {
		noContentInt = 1
	}

	resp := &ArticleList{}
	if err := api.Client.HTTPPostJson(ctx, apiBatchGetArticles, map[string]interface{}{
		"no_content": noContentInt,
		"offset":     offset,
		"count":      count,
	}, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

/*
获取已发布图文信息
https://developers.weixin.qq.com/doc/service/api/public/api_freepublishgetarticle.html
*/
func (api *MaterialApi) GetArticle(
	ctx context.Context, articleID string,
) ([]*Article, error) {
	resp := &struct {
		utils.WeixinError
		NewsItems []*Article `json:"news_item"`
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiGetArticle, map[string]interface{}{
		"article_id": articleID,
	}, resp); err != nil {
		return nil, err
	}

	return resp.NewsItems, nil
}
