package authorizer

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiQrcodeCreate = "/cgi-bin/qrcode/create"
)

type CreateStrOrcodeReply struct {
	utils.WeixinError

	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	Url           string `json:"url"`
}

func (api *Authorizer) CreateStrQrcode(
	ctx context.Context,
	sceneSstr string,
	expireSeconds int64,
	actionName string,
) (*CreateStrOrcodeReply, error) {
	if len(actionName) == 0 {
		actionName = "QR_STR_SCENE"
	}
	data := map[string]interface{}{
		"action_name": actionName,
		"action_info": map[string]interface{}{
			"scene": map[string]interface{}{
				"scene_str": sceneSstr,
			},
		},
	}
	if expireSeconds > 0 {
		data["expire_seconds"] = expireSeconds
	}
	result := CreateStrOrcodeReply{}
	e := api.Client.HTTPPostJson(ctx, apiQrcodeCreate, data, &result)
	if e != nil {
		return nil, e
	}
	return &result, nil
}
