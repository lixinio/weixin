package official_account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
从微信服务器获取新的 AccessToken
See: https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
See: https://github.com/fastwego/offiaccount/blob/master/client.go#L268
*/
func (officialAccount *OfficialAccount) refreshAccessTokenFromWXServer() (accessToken string, expiresIn int, err error) {
	params := url.Values{}
	params.Add("appid", officialAccount.Config.Appid)
	params.Add("secret", officialAccount.Config.Secret)
	params.Add("grant_type", "client_credential")
	url := WXServerUrl + "/cgi-bin/token?" + params.Encode()

	response, err := http.Get(url)
	if err != nil {
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("GET %s RETURN %s", url, response.Status)
		return
	}

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	var result = struct {
		AccessToken string  `json:"access_token"`
		ExpiresIn   int     `json:"expires_in"`
		Errcode     float64 `json:"errcode"`
		Errmsg      string  `json:"errmsg"`
	}{}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		err = fmt.Errorf("unmarshal error %s", string(resp))
		return
	}

	if result.AccessToken == "" {
		err = fmt.Errorf("%s", string(resp))
		return
	}

	return result.AccessToken, result.ExpiresIn, nil
}
