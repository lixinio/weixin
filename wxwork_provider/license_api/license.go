package license_api

import "github.com/lixinio/weixin/wxwork_provider"

// LicenseApi 接口应用许可相关api
type LicenseApi struct {
	*wxwork_provider.WxWorkProvider
}

func NewApi(provider *wxwork_provider.WxWorkProvider) *LicenseApi {
	return &LicenseApi{provider}
}
