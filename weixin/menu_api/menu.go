package menu_api

import (
	"context"
	"github.com/lixinio/weixin/utils"
)

const (
	apiCreateMenu = "/cgi-bin/menu/create"
)

type Button struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Key        string    `json:"key,omitempty"`
	Url        string    `json:"url,omitempty"`
	AppID      string    `json:"appid,omitempty"`
	Pagepath   string    `json:"pagepath,omitempty"`
	SubButtons []*Button `json:"sub_button,omitempty"`
}

type Menu struct {
	Buttons []*Button `json:"button"`
}

type MenuApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *MenuApi {
	return &MenuApi{Client: client}
}

func (api *MenuApi) Create(ctx context.Context, menu *Menu) error {
	return api.Client.HTTPPostJson(ctx, apiCreateMenu, menu, nil)
}
