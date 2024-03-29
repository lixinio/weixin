// 内容安全检测

package content_check

import (
	"context"
	"errors"
	"net/http"

	"github.com/lixinio/weixin/utils"
)

const (
	apiMsgSecCheck  = "/wxa/msg_sec_check"
	msgCheckVersion = 2

	apiImgSecCheck    = "/wxa/img_sec_check"
	imgCheckFieldName = "img_check"
	imgCheckFileName  = "img_check"

	SensitiveImgErrCode = 87014
)

// ContentCheckApi 内容检测api
type ContentCheckApi struct {
	*utils.Client
}

// MsgCheckResult 文本检测结果
type MsgCheckResult struct {
	utils.WeixinError
	TraceID int64 // 唯一请求标识，标记单次请求
	Result  struct {
		Suggest string // 建议, 有risky、pass、review三种值
		Label   int64  // 命中标签枚举值, 如100 正常; 10001 广告 ...
	} // 综合结果
	Detail []struct {
		Strategy string
		ErrCode  int64
		Suggest  string
		Label    int64
		Prob     int    // 0-100，代表置信度，越高代表越有可能属于当前返回的标签（label）
		Keyword  string // 命中的自定义关键词
	} // 详细检测结果
}

func NewApi(client *utils.Client) *ContentCheckApi {
	return &ContentCheckApi{client}
}

// CheckMsg 过滤敏感信息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/sec-check/security.msgSecCheck.html
func (api *ContentCheckApi) CheckMsg(
	ctx context.Context,
	openid string,
	scene int,
	content string,
	nickname string,
	title string,
	signature string,
) (*MsgCheckResult, error) {
	result := &MsgCheckResult{}
	payload := struct {
		Version   int    `json:"version"`
		OpenID    string `json:"openid"`
		Scene     int    `json:"scene"`
		Content   string `json:"content"`
		Nickname  string `json:"nickname"`
		Title     string `json:"title"`
		Signature string `json:"signature"`
	}{
		Version:   msgCheckVersion,
		OpenID:    openid,
		Scene:     scene,
		Content:   content,
		Nickname:  nickname,
		Title:     title,
		Signature: signature,
	}

	if err := api.Client.HTTPPostJson(ctx, apiMsgSecCheck, payload, result); err != nil {
		return nil, err
	}

	return result, nil
}

// CheckImg 过滤敏感图片
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/sec-check/security.imgSecCheck.html
func (api *ContentCheckApi) CheckImg(
	ctx context.Context,
	imgURL string,
) (sensitive bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	if err != nil {
		return true, err
	}

	imgResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return true, err
	}

	defer imgResp.Body.Close()

	weixinErr := &utils.WeixinError{}
	err = api.Client.HttpFile(
		ctx, apiImgSecCheck, "media", imgCheckFileName, imgResp.Body, nil, weixinErr,
	)
	if err != nil {
		if errors.As(err, &weixinErr) && weixinErr.ErrCode == SensitiveImgErrCode {
			return true, nil
		}
		return true, err
	}

	return false, nil
}
