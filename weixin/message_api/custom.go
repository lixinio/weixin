package message_api

// 客服消息

import "context"

const (
	apiCustomSend = "/cgi-bin/message/custom/send"
)

type MessageHeader struct {
	ToUser  string `json:"touser,omitempty"`
	MsgType string `json:"msgtype"`
}

type TextMessage struct {
	*MessageHeader
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

/*
发送客服消息（文本）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
func (api *MessageApi) SendCustomTextMessage(
	ctx context.Context, openID, content string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &TextMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "text",
		},
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}, nil)
}

/*
发送客服消息（图片）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type ImageMessage struct {
	*MessageHeader
	Image struct {
		MediaID string `json:"media_id"`
	} `json:"image"`
}

func (api *MessageApi) SendCustomImageMessage(
	ctx context.Context, openID, mediaID string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &ImageMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "image",
		},
		Image: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, nil)
}

/*
发送客服消息（语音）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type VoiceMessage struct {
	*MessageHeader
	Voice struct {
		MediaID string `json:"media_id"`
	} `json:"voice"`
}

func (api *MessageApi) SendCustomVoiceMessage(
	ctx context.Context, openID, mediaID string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &VoiceMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "voice",
		},
		Voice: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, nil)
}

/*
发送客服消息（视频）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type Video struct {
	MediaID      string `json:"media_id"`
	ThumbMediaID string `json:"thumb_media_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type VideoMessage struct {
	*MessageHeader
	Video *Video `json:"video"`
}

func (api *MessageApi) SendCustomVideoMessage(
	ctx context.Context, openID string, video *Video,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &VideoMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "video",
		},
		Video: video,
	}, nil)
}

/*
发送客服消息（音乐）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type Music struct {
	MusicURL     string `json:"musicurl"`
	HqMusicURL   string `json:"hqmusicurl"`
	ThumbMediaID string `json:"thumb_media_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
}

type MusicMessage struct {
	*MessageHeader
	Music *Music `json:"music"`
}

func (api *MessageApi) SendCustomMusicMessage(
	ctx context.Context, openID string, music *Music,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &MusicMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "music",
		},
		Music: music,
	}, nil)
}

/*
发送客服消息（图文消息（点击跳转到外链））
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type Article struct {
	PicURL      string `json:"picurl"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type NewsMessage struct {
	*MessageHeader
	News struct {
		Articles []*Article `json:"articles"`
	} `json:"news"`
}

func (api *MessageApi) SendCustomNewsMessage(
	ctx context.Context, openID string, articles []*Article,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &NewsMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "news",
		},
		News: struct {
			Articles []*Article `json:"articles"`
		}{
			Articles: articles,
		},
	}, nil)
}

/*
发送客服消息（图文消息（点击跳转到图文消息页面））
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type MpnewsMessage struct {
	*MessageHeader
	Mpnews struct {
		MediaID string `json:"media_id"`
	} `json:"mpnews"`
}

func (api *MessageApi) SendCustomMpnewsMessage(
	ctx context.Context, openID, mediaID string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &MpnewsMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "mpnews",
		},
		Mpnews: struct {
			MediaID string `json:"media_id"`
		}{
			MediaID: mediaID,
		},
	}, nil)
}

/*
发送客服消息（图文消息（点击跳转到图文消息页面））
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type MpnewsArticleMessage struct {
	*MessageHeader
	MpnewsArticle struct {
		ArticleID string `json:"article_id"`
	} `json:"mpnewsarticle"`
}

func (api *MessageApi) SendCustomMpnewsArticleMessage(
	ctx context.Context, openID, articleID string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &MpnewsArticleMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "mpnewsarticle",
		},
		MpnewsArticle: struct {
			ArticleID string `json:"article_id"`
		}{
			ArticleID: articleID,
		},
	}, nil)
}

/*
发送客服消息（菜单）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type Menu struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type MenuMessage struct {
	*MessageHeader
	Menu struct {
		HeadContent string  `json:"head_content"`
		TailContent string  `json:"tail_content"`
		Menus       []*Menu `json:"list"`
	} `json:"msgmenu"`
}

func (api *MessageApi) SendCustomMenuMessage(
	ctx context.Context, openID, head, tail string, menus []*Menu,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &MenuMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "msgmenu",
		},
		Menu: struct {
			HeadContent string  `json:"head_content"`
			TailContent string  `json:"tail_content"`
			Menus       []*Menu `json:"list"`
		}{
			HeadContent: head,
			TailContent: tail,
			Menus:       menus,
		},
	}, nil)
}

/*
发送客服消息（卡券）
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/
type WxcardMessage struct {
	*MessageHeader
	Wxcard struct {
		CardID string `json:"card_id"`
	} `json:"wxcard"`
}

func (api *MessageApi) SendCustomWxcardMessage(
	ctx context.Context, openID, cardID string,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &WxcardMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "wxcard",
		},
		Wxcard: struct {
			CardID string `json:"card_id"`
		}{
			CardID: cardID,
		},
	}, nil)
}

/*
发送客服消息（小程序卡片（要求小程序与公众号已关联））
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Service_Center_messages.html#7
*/

type MiniProgram struct {
	Title        string `json:"title"`
	AppID        string `json:"appid"`
	PagePath     string `json:"pagepath"`
	ThumbMediaID string `json:"thumb_media_id"`
}

type MiniProgrampageMessage struct {
	*MessageHeader
	MiniProgram *MiniProgram `json:"miniprogrampage"`
}

func (api *MessageApi) SendCustomMiniProgramMessage(
	ctx context.Context, openID string, mp *MiniProgram,
) error {
	return api.Client.HTTPPostJson(ctx, apiCustomSend, &MiniProgrampageMessage{
		MessageHeader: &MessageHeader{
			ToUser:  openID,
			MsgType: "miniprogrampage",
		},
		MiniProgram: mp,
	}, nil)
}
