package menu_api

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiCreateMenu = "/cgi-bin/menu/create"
	apiGetMenu    = "/cgi-bin/menu/get"
	apiDeleteMenu = "/cgi-bin/menu/delete"
)

type MenuApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *MenuApi {
	return &MenuApi{Client: client}
}

// 菜单项目
type MenuItem struct {
	Name string `json:"name"`           // 菜单标题，不超过16个字节，子菜单不超过60个字节
	Type string `json:"type,omitempty"` // 菜单的响应动作类型
	Key  string `json:"key,omitempty"`  // 菜单KEY值，用于消息接口推送，不超过128字节。click等点击类型必须。
	// 网页链接，用户点击菜单可打开链接，不超过1024字节。
	// type为miniprogram时，不支持小程序的老版本客户端将打开本url。view、miniprogram类型必填。
	Uri string `json:"url,omitempty"`
	// 调用新增永久素材接口返回的合法media_id。media_id类型和view_limited类型必须
	MediaID *string `json:"media_id,omitempty"`
	// 发布后获得的合法 article_id，article_id类型和article_view_limited类型必须
	ArticleID  *string     `json:"article_id,omitempty"`
	AppID      string      `json:"appid,omitempty"`      // 小程序的appid（仅认证公众号可配置），miniprogram类型必须
	PagePath   string      `json:"pagepath,omitempty"`   // 小程序的页面路径，miniprogram类型必须
	SubButtons []*MenuItem `json:"sub_button,omitempty"` // 子菜单
}

// 创建自定义菜单
// https://developers.weixin.qq.com/doc/subscription/api/custommenu/api_createcustommenu.html
func (api *MenuApi) CreateMenu(ctx context.Context, menus []*MenuItem) error {
	req := struct {
		Buttons []*MenuItem `json:"button"`
	}{
		Buttons: menus,
	}
	return api.Client.HTTPPostJson(ctx, apiCreateMenu, &req, nil)
}

// 查询自定义菜单信息
// https://developers.weixin.qq.com/doc/subscription/api/custommenu/api_getcurrentselfmenuinfo.html
func (api *MenuApi) GetMenu(ctx context.Context) ([]*MenuItem, error) {
	resp := struct {
		utils.WeixinError
		Menus struct {
			Buttons []*MenuItem `json:"button"`
		} `json:"menu"`
	}{
		Menus: struct {
			Buttons []*MenuItem "json:\"button\""
		}{
			Buttons: []*MenuItem{},
		},
	}

	err := api.Client.HTTPGet(ctx, apiGetMenu, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Menus.Buttons, nil
}

// 删除自定义菜单
// https://developers.weixin.qq.com/doc/subscription/api/custommenu/api_deletemenu.html
func (api *MenuApi) DeleteMenu(ctx context.Context) error {
	return api.Client.HTTPGet(ctx, apiDeleteMenu, nil)
}
