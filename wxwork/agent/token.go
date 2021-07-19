package agent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/wxwork"
)

/*
从微信服务器获取新的 AccessToken

See: https://developers.weixin.qq.com/doc/corporation/Basic_Information/Get_access_token.html
*/
func (agent *Agent) refreshAccessTokenFromWXServer() (accessToken string, expiresIn int, err error) {
	params := url.Values{}
	params.Add("corpid", agent.wxwork.Config.Corpid)
	params.Add("corpsecret", agent.Config.Secret)
	url := wxwork.QyWXServerUrl + "/cgi-bin/gettoken?" + params.Encode()

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

	var result utils.TokenResponse

	err = json.Unmarshal(resp, &result)
	if err != nil {
		err = fmt.Errorf("Unmarshal error %s", string(resp))
		return
	}

	if result.AccessToken == "" {
		err = fmt.Errorf("%s", string(resp))
		return
	}

	return result.AccessToken, result.ExpiresIn, nil
}
