package message_api

// 模板消息

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiSetIndustry           = "/cgi-bin/template/api_set_industry"
	apiGetIndustry           = "/cgi-bin/template/get_industry"
	apiTemplateSend          = "/cgi-bin/message/template/send"
	apiAddTemplate           = "/cgi-bin/template/api_add_template"
	apiGetAllPrivateTemplate = "/cgi-bin/template/get_all_private_template"
	apiDelPrivateTemplate    = "/cgi-bin/template/del_private_template"
)

type MessageApi struct{ *utils.Client }

func NewApi(client *utils.Client) *MessageApi {
	return &MessageApi{Client: client}
}

type TemplateMessageMp struct {
	AppID    string `json:"appid"`              // 所需跳转到的小程序appid（该小程序 appid 必须与发模板消息的公众号是绑定关联关系）
	PagePath string `json:"pagepath,omitempty"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），要求该小程序已发布
}

type TemplateMessageData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"` // 模板内容字体颜色，不填默认为黑色
}

type TemplateMessage struct {
	ToUser      string                          `json:"touser,omitempty"`        // 接收者openid
	TemplateID  string                          `json:"template_id"`             // 模板ID
	URL         string                          `json:"url,omitempty"`           // 模板跳转链接（海外帐号没有跳转能力）
	MiniProgram *TemplateMessageMp              `json:"miniprogram,omitempty"`   // 跳小程序所需数据，不需跳小程序可不用传该数据
	ClientMsgID string                          `json:"client_msg_id,omitempty"` // 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果。若无防重入需求，可不填
	Datas       map[string]*TemplateMessageData `json:"data"`                    // 模板数据
}

/*
发送模板消息
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) SendTemplateMessage(
	ctx context.Context, msg *TemplateMessage,
) (int64, error) {
	resp := &struct {
		utils.WeixinError
		MsgID int64 `json:"msgid"`
	}{}
	if err := api.Client.HTTPPostJson(ctx, apiTemplateSend, msg, resp); err != nil {
		return 0, err
	}
	return resp.MsgID, nil
}

/*
设置所属行业
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) SetIndustry(
	ctx context.Context, industryID1, industryID2 string,
) error {
	if err := api.Client.HTTPPostJson(ctx, apiSetIndustry, map[string]string{
		"industry_id1": industryID1,
		"industry_id2": industryID2,
	}, nil); err != nil {
		return err
	}
	return nil
}

type Industry struct {
	FirstClass  string `json:"first_class"`
	SecondClass string `json:"second_class"`
}

type IndustryInfo struct {
	utils.WeixinError
	PrimaryIndustry   *Industry `json:"primary_industry"`
	SecondaryIndustry *Industry `json:"secondary_industry"`
}

/*
获取设置的行业信息
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) GetIndustry(ctx context.Context) (*IndustryInfo, error) {
	resp := &IndustryInfo{}
	if err := api.Client.HTTPGet(ctx, apiGetIndustry, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

/*
获得模板ID(添加模板)
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) AddTemplate(
	ctx context.Context,
	templateIdShort string,
	keywords ...string,
) (string, error) {
	resp := &struct {
		utils.WeixinError
		TemplateID string `json:"template_id"`
	}{}

	if err := api.Client.HTTPPostJson(ctx, apiAddTemplate, map[string]interface{}{
		"template_id_short": templateIdShort,
		"keyword_name_list": keywords,
	}, resp); err != nil {
		return "", err
	}
	return resp.TemplateID, nil
}

/*
删除模板
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) DelPrivateTemplate(
	ctx context.Context, templateID string,
) error {
	if err := api.Client.HTTPPostJson(ctx, apiDelPrivateTemplate, map[string]string{
		"template_id": templateID,
	}, nil); err != nil {
		return err
	}
	return nil
}

type PrivateTemplate struct {
	TemplateID      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

/*
获取模板列表
https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
*/
func (api *MessageApi) GetAllPrivateTemplate(
	ctx context.Context,
) ([]*PrivateTemplate, error) {
	resp := &struct {
		utils.WeixinError
		TemplateList []*PrivateTemplate `json:"template_list"`
	}{}
	if err := api.Client.HTTPGet(ctx, apiGetAllPrivateTemplate, resp); err != nil {
		return nil, err
	}
	return resp.TemplateList, nil
}
