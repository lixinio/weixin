package qrcode_api

import (
	"context"
	"io"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGetWxaCodeUnlimit = "/wxa/getwxacodeunlimit"
	apiGetWxaCode        = "/wxa/getwxacode"
	apiCreateWxaQrcode   = "/cgi-bin/wxaapp/createwxaqrcode"
	apiCreateQrcode      = "/cgi-bin/qrcode/create" // 生成带参数的二维码 (服务号)
)

type QrcodeApi struct {
	*utils.Client
}

func NewApi(client *utils.Client) *QrcodeApi {
	return &QrcodeApi{Client: client}
}

/*
获取小程序码，适用于需要的码数量极多的业务场景。通过该接口生成的小程序码，永久有效，数量暂无限制。
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/qrcode/getwxacodeunlimit.html
*/

type GetWxaCodeUnlimitRequest struct {
	Scene      string     `json:"scene"`
	Page       string     `json:"page,omitempty"` // 页面 page，例如 pages/index/index，根路径前不要填加 /，不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
	CheckPath  bool       `json:"check_path"`
	EnvVersion string     `json:"env_version,omitempty"`
	Width      int        `json:"width,omitempty"` // 二维码的宽度，单位 px，最小 280px，最大 1280px
	AutoColor  bool       `json:"auto_color,omitempty"`
	IsHyaline  bool       `json:"is_hyaline,omitempty"`
	LineColor  *LineColor `json:"line_color,omitempty"`
}

func (api *QrcodeApi) GetWxaCodeUnlimit(
	ctx context.Context, param *GetWxaCodeUnlimitRequest, // content io.Writer,
) ([]byte, error) {
	resp, err := api.Client.HTTPPostDownload(ctx, apiGetWxaCodeUnlimit, param, nil)
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

type LineColor struct {
	R int32 `json:"r"`
	G int32 `json:"g"`
	B int32 `json:"b"`
}

/*
获取小程序码，适用于需要的码数量较少的业务场景。通过该接口生成的小程序码，永久有效，有数量限制
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.get.html
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/qrcode/getwxacode.html
*/
type GetWxaCodeRequest struct {
	// 扫码进入的小程序页面路径，最大长度 128 字节，不能为空；对于小游戏，可以只传入 query 部分，来实现传参效果，如：传入 "?foo=bar"，即可在 wx.getLaunchOptionsSync 接口中的 query 参数获取到 {foo:"bar"}。
	Path       string     `json:"path,omitempty"`
	Width      int        `json:"width,omitempty"` // 二维码的宽度，单位 px，最小 280px，最大 1280px
	AutoColor  bool       `json:"auto_color,omitempty"`
	IsHyaline  bool       `json:"is_hyaline,omitempty"`
	LineColor  *LineColor `json:"line_color,omitempty"`
	EnvVersion string     `json:"env_version"`
}

func (api *QrcodeApi) GetWxaCode(
	ctx context.Context, param *GetWxaCodeRequest,
) ([]byte, error) {
	resp, err := api.Client.HTTPPostDownload(ctx, apiGetWxaCode, param, nil)
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

/*
获取小程序二维码，适用于需要的码数量较少的业务场景。通过该接口生成的小程序码，永久有效，有数量限制
https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.createQRCode.html
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/qrcode/createwxaqrcode.html
*/
func (api *QrcodeApi) CreateWxaQRCode(
	ctx context.Context, path string, width int,
) ([]byte, error) {
	param := &struct {
		Path  string `json:"path"`
		Width int    `json:"width,omitempty"`
	}{path, width}
	resp, err := api.Client.HTTPPostDownload(ctx, apiCreateWxaQrcode, param, nil)
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

type QrcodeResp struct {
	utils.WeixinError
	Ticket        string `json:"ticket"`         // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	ExpireSeconds int32  `json:"expire_seconds"` // 该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）。
	Url           string `json:"url"`            // 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}

// 生成带参数的二维码（服务号）
// https://developers.weixin.qq.com/doc/service/api/qrcode/qrcodes/api_createqrcode.html
func (api *QrcodeApi) CreateQRCode(
	ctx context.Context,
	actionName string,
	sceneID int32, sceneStr string,
	expireSeconds int32,
) (*QrcodeResp, error) {
	param := &struct {
		ExpireSeconds int32  `json:"expire_seconds,omitempty"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene struct {
				SceneID  int32  `json:"scene_id,omitempty"`
				SceneStr string `json:"scene_str,omitempty"`
			} `json:"scene"`
		} `json:"action_info"`
	}{
		ExpireSeconds: expireSeconds,
		ActionName:    actionName,
		ActionInfo: struct {
			Scene struct {
				SceneID  int32  "json:\"scene_id,omitempty\""
				SceneStr string "json:\"scene_str,omitempty\""
			} "json:\"scene\""
		}{
			Scene: struct {
				SceneID  int32  "json:\"scene_id,omitempty\""
				SceneStr string "json:\"scene_str,omitempty\""
			}{
				SceneID:  sceneID,
				SceneStr: sceneStr,
			},
		},
	}

	resp := &QrcodeResp{}

	err := api.Client.HTTPPostJson(ctx, apiCreateQrcode, param, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
