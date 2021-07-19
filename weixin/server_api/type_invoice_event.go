package server_api

// EventAuthorizeInvoice https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html#6
type EventAuthorizeInvoice struct {
	Event
	SuccOrderId    string // 授权成功的订单号，与失败订单号两者必显示其一
	FailOrderId    string // 授权失败的订单号，与成功订单号两者必显示其一
	AuthorizeAppId string // 获取授权页链接的AppId
	Source         string // 授权来源，web：公众号开票，app：app开票，wxa：小程序开票，wap：h5开票
}
