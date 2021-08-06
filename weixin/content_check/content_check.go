// 内容安全检测

package content_check

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/official_account"
	"net/http"
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
	OfficialAccount *official_account.OfficialAccount
}

// MsgCheckResult 文本检测结果
type MsgCheckResult struct {
	ErrCode int64
	ErrMsg  string
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

// ImgCheckResult 图片检测结果
type ImgCheckResult struct {
	ErrCode int64
	ErrMsg  string
}

func NewOfficialAccountApi(officialAccount *official_account.OfficialAccount) *ContentCheckApi {
	return &ContentCheckApi{
		officialAccount.Client,
		officialAccount,
	}
}

// CheckMsg 过滤敏感信息
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/sec-check/security.msgSecCheck.html
func (api *ContentCheckApi) CheckMsg(ctx context.Context, openid string, scene int, content string, nickname string, title string, signature string) (*MsgCheckResult, error) {
	result := MsgCheckResult{}
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

	err := api.Client.ApiPostWrapper(ctx, apiMsgSecCheck, payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CheckImg 过滤敏感图片
// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/sec-check/security.imgSecCheck.html
func (api *ContentCheckApi) CheckImg(ctx context.Context, imgURL string) (sensitive bool, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imgURL, nil)
	if err != nil {
		return true, err
	}

	client := http.Client{}
	imgResp, err := client.Do(req)
	if err != nil {
		return true, err
	}

	resp, err := api.Client.Upload(ctx, apiImgSecCheck, imgCheckFieldName, imgCheckFileName, imgResp.Body)
	if err != nil {
		weixinErr := utils.WeixinError{}
		if errors.As(err, &weixinErr) && weixinErr.Errcode == SensitiveImgErrCode {
			return true, nil
		}
		return true, err
	}

	result := ImgCheckResult{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return true, err
	}

	return false, nil
}
