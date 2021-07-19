package invoice_api

// https://github.com/fastwego/offiaccount/blob/master/apis/invoice/invoice.go
// Copyright 2020 FastWeGo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package invoice 微信发票

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"

	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/weixin/official_account"
)

const (
	apiGetAuthUrl                   = "/card/invoice/getauthurl"
	apiGetAuthData                  = "/card/invoice/getauthdata"
	apiRejectInsert                 = "/card/invoice/rejectinsert"
	apiMakeOutInvoice               = "/card/invoice/makeoutinvoice"
	apiClearOutInvoice              = "/card/invoice/clearoutinvoice"
	apiQueryInvoceInfo              = "/card/invoice/queryinvoceinfo"
	apiSetUrl                       = "/card/invoice/seturl"
	apiPlatformCreateCard           = "/card/invoice/platform/createcard"
	apiPlatformSetpdf               = "/card/invoice/platform/setpdf"
	apiPlatformGetpdf               = "/card/invoice/platform/getpdf"
	apiInsert                       = "/card/invoice/insert"
	apiPlatformUpdateStatus         = "/card/invoice/platform/updatestatus"
	apiReimburseGetInvoiceInfo      = "/card/invoice/reimburse/getinvoiceinfo"
	apiReimburseGetInvoiceBatch     = "/card/invoice/reimburse/getinvoicebatch"
	apiReimburseUpdateInvoiceStatus = "/card/invoice/reimburse/updateinvoicestatus"
	apiReimburseUpdateStatusBatch   = "/card/invoice/reimburse/updatestatusbatch"
	apiGetUserTitleUrl              = "/card/invoice/biz/getusertitleurl"
	apiGetSelectTitleUrl            = "/card/invoice/biz/getselecttitleurl"
	apiScanTitle                    = "/card/invoice/scantitle"
	apiSetbizattr                   = "/card/invoice/setbizattr"
)

type InvoiceApi struct {
	*utils.Client
	OfficialAccount *official_account.OfficialAccount
}

func NewOfficialAccountApi(officialAccount *official_account.OfficialAccount) *InvoiceApi {
	return &InvoiceApi{
		Client:          officialAccount.Client,
		OfficialAccount: officialAccount,
	}
}

// https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html#17
// 商户获取授权链接之前，需要先设置商户的联系方式
type SetbizattrObj struct {
	Phone   string `json:"phone"`
	TimeOut int    `json:"time_out"`
}

func (api *InvoiceApi) SetbizattrRaw(payload []byte, params url.Values) (resp []byte, err error) {
	return api.Client.HTTPPost(
		apiSetbizattr+"?"+params.Encode(),
		bytes.NewReader(payload),
		"application/json;charset=utf-8",
	)
}
func (api *InvoiceApi) SetContact(param *SetbizattrObj) error {
	params := url.Values{}
	params.Add("action", "set_contact")

	payload := &struct {
		Contact *SetbizattrObj `json:"contact"`
	}{
		Contact: param,
	}

	return utils.ApiPostWrapperEx(api.SetbizattrRaw, payload, params, nil)
}
func (api *InvoiceApi) GetContact() (*SetbizattrObj, error) {
	params := url.Values{}
	params.Add("action", "get_contact")

	result := &struct {
		Contact *SetbizattrObj `json:"contact"`
	}{}

	err := utils.ApiPostWrapperEx(api.SetbizattrRaw, "{}", params, result)
	if err != nil {
		return nil, err
	}
	return result.Contact, nil
}

// https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html#9
// 当用户使用type=1的类型的授权页时，可以使用本接口设置授权页上需要用户填写的信息。若使用type=0或type=2类型的授权页，无需调用本接口。本接口为一次性设置，后续除非在需要调整页面字段时才需要再次调用。
// 注意，设置为显示状态的字段均为必填字段，用户若不填写将无法进入后续流程

type AuthFieldObj struct {
	UserField *AuthUserField `json:"user_field,"`
	BizField  *AuthBizField  `json:"biz_field,omitempty"`
}

type AuthCustomField struct {
	Key       string `json:"key"`
	IsRequire int    `json:"is_require,omitempty"` // 0：否，1：是， 默认为0
	Notice    string `json:"notice"`
}

type AuthUserField struct {
	ShowTitle    int               `json:"show_title"`    // 是否填写抬头，0为否，1为是
	ShowPhone    int               `json:"show_phone"`    // 是否填写电话号码，0为否，1为是
	ShowEmail    int               `json:"show_email"`    // 是否填写邮箱，0为否，1为是
	RequirePhone int               `json:"require_phone"` // 电话号码是否必填,0为否，1为是
	RequireEmail int               `json:"require_email"` // 邮箱是否必填，0位否，1为是
	CustomFields []AuthCustomField `json:"custom_field,omitempty"`
}

type AuthBizField struct {
	ShowTitle       int               `json:"show_title"`        // 是否填写抬头，0为否，1为是
	ShowTaxNO       int               `json:"show_tax_no"`       // 是否填写税号，0为否，1为是
	ShowAddr        int               `json:"show_addr"`         // 是否填写单位地址，0为否，1为是
	ShowPhone       int               `json:"show_phone"`        // 是否填写电话号码，0为否，1为是
	ShowBankType    int               `json:"show_bank_type"`    // 是否填写开户银行，0为否，1为是
	ShowBankNO      int               `json:"show_bank_no"`      // 是否填写银行帐号，0为否，1为是
	RequireTaxNO    int               `json:"require_tax_no"`    // 税号是否必填，0为否，1为是
	RequireAddr     int               `json:"require_addr"`      // 单位地址是否必填，0为否，1为是
	RequirePhone    int               `json:"require_phone"`     // 电话号码是否必填，0为否，1为是
	RequireBankType int               `json:"require_bank_type"` // 开户类型是否必填，0为否，1为是
	RequireBankNO   int               `json:"require_bank_no"`   // 税号是否必填，0为否，1为是
	CustomFields    []AuthCustomField `json:"custom_field,omitempty"`
}

func (api *InvoiceApi) SetAuthField(param *AuthFieldObj) error {
	params := url.Values{}
	params.Add("action", "set_auth_field")

	payload := &struct {
		AuthField *AuthFieldObj `json:"auth_field"`
	}{
		AuthField: param,
	}

	return utils.ApiPostWrapperEx(api.SetbizattrRaw, payload, params, nil)
}

/*
获取自身的开票平台识别码
开票平台可以通过此接口获得本开票平台的预开票url，进而获取s_pappid。开票平台将该s_pappid并透传给商户，商户可以通过该s_pappid参数在微信电子发票方案中标识出为自身提供开票服务的开票平台
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Invoicing_Platform_API_List.html
POST https://api.weixin.qq.com/card/invoice/seturl?access_token={access_token}
*/
func (api *InvoiceApi) SetUrlRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiSetUrl, bytes.NewReader(payload), "application/json;charset=utf-8")
}

func (api *InvoiceApi) SetUrl() (string, error) {
	result := &struct {
		InvoiceUrl string `json:"invoice_url"`
	}{}

	err := utils.ApiPostWrapper(api.SetUrlRaw, "{}", result)
	if err != nil {
		return "", err
	}
	return result.InvoiceUrl, nil
}

type AuthUrlObj struct {
	SPappID     string `json:"s_pappid"`
	OrderID     string `json:"order_id"`
	Money       int    `json:"money"`
	Timestamp   int64  `json:"timestamp"`
	Source      string `json:"source"` // 开票来源，app：app开票，web：微信h5开票，wxa：小程序开发票，wap：普通网页开票
	RedirectURL string `json:"redirect_url,omitempty"`
	Ticket      string `json:"ticket"`
	Type        int    `json:"type"` // 授权类型，0：开票授权，1：填写字段开票授权，2：领票授权
}

type AuthUrlResult struct {
	AuthURL string `json:"auth_url"`
	AppID   string `json:"appid,omitempty"`
}

/*
获取授权页链接
本接口供商户调用。商户通过本接口传入订单号、开票平台标识等参数，获取授权页的链接。在微信中向用户展示授权页，当用户点击了授权页上的“领取发票”/“申请开票”按钮后，即完成了订单号与该用户的授权关系绑定，后续开票平台可凭此订单号发起将发票卡券插入用户卡包的请求，微信也将据此授权关系校验是否放行插卡请求
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html
POST https://api.weixin.qq.com/card/invoice/getauthurl?access_token={access_token}
*/
func (api *InvoiceApi) GetAuthUrlRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiGetAuthUrl, bytes.NewReader(payload), "application/json;charset=utf-8")
}

func (api *InvoiceApi) GetAuthUrl(param *AuthUrlObj) (*AuthUrlResult, error) {
	var result AuthUrlResult
	err := utils.ApiPostWrapper(api.GetAuthUrlRaw, param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type AuthDataObj struct {
	OrderID string `json:"order_id"`
	SPappID string `json:"s_pappid"`
}

type CustomField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// 若用户填入的是个人抬头
type AuthDataUserField struct {
	Title        string        `json:"title"`
	Phone        string        `json:"phone"`
	Email        string        `json:"email"`
	CustomFields []CustomField `json:"custom_field"`
}

// 若用户填入的是单位抬头
type AuthDataBizField struct {
	Title        string        `json:"title"`
	Phone        string        `json:"phone"`
	TaxNO        string        `json:"tax_no"`
	Addr         string        `json:"addr"`
	BankType     string        `json:"bank_type"`
	BankNO       string        `json:"bank_no"`
	CustomFields []CustomField `json:"custom_field"`
}

type AuthDataResult struct {
	// "auth success" / "reject insert" / "never auth"
	InvoiceStatus string `json:"invoice_status"` // 订单授权状态，当errcode为0时会出现
	AuthTime      int64  `json:"auth_time"`
	UserAuthInfo  struct {
		UserField *AuthDataUserField `json:"user_field,omitempty"`
		BizField  *AuthDataBizField  `json:"biz_field,omitempty"`
	} `json:"user_auth_info"`
}

/*
查询授权完成状态
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html
POST https://api.weixin.qq.com/card/invoice/getauthdata?access_token={access_token}
*/
func (api *InvoiceApi) GetAuthDataRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiGetAuthData, bytes.NewReader(payload), "application/json;charset=utf-8")
}

func (api *InvoiceApi) GetAuthData(param *AuthDataObj) (*AuthDataResult, error) {
	var result AuthDataResult
	err := utils.ApiPostWrapper(api.GetAuthDataRaw, param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type RejectInsertObj struct {
	OrderID string `json:"order_id"`
	SPappID string `json:"s_pappid"`
	Reason  string `json:"reason"`        // 商家解释拒绝开票的原因，如重复开票，抬头无效、已退货无法开票等
	URL     string `json:"url,omitempty"` // 跳转链接，引导用户进行下一步处理，如重新发起开票、重新填写抬头、展示订单情况等
}

/*
拒绝开票
户完成授权后，商户若发现用户提交信息错误、或者发生了退款时，可以调用该接口拒绝开票并告知用户。拒绝开票后，该订单无法向用户再次开票。已经拒绝开票的订单，无法再次使用，如果要重新开票，需使用新的order_id，获取授权链接，让用户再次授权
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html
POST https://api.weixin.qq.com/card/invoice/rejectinsert?access_token={access_token}
*/
func (api *InvoiceApi) RejectInsertRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiRejectInsert, bytes.NewReader(payload), "application/json;charset=utf-8")
}

func (api *InvoiceApi) RejectInsert(param *RejectInsertObj) error {
	return utils.ApiPostWrapper(api.RejectInsertRaw, param, nil)
}

type CreateCardBaseInfo struct {
	LogoUrl              string `json:"logo_url"`                          // 发票商家 LOGO ，请参考 新增永久素材
	Title                string `json:"title"`                             // 收款方（显示在列表），上限为 9 个汉字，建议填入商户简称
	CustomUrlName        string `json:"custom_url_name"`                   // 开票平台自定义入口名称，与 custom_url 字段共同使用，长度限制在 5 个汉字内
	CustomURL            string `json:"custom_url"`                        // 开票平台自定义入口跳转外链的地址链接 , 发票外跳的链接会带有发票参数，用于标识是从哪张发票跳出的链接
	CustomUrlSubTitle    string `json:"custom_url_sub_title,omitempty"`    // 显示在入口右侧的 tips ，长度限制在 6 个汉字内
	PromotionUrlName     string `json:"promotion_url_name,omitempty"`      // 营销场景的自定义入口
	PromotionURL         string `json:"promotion_url,omitempty"`           // 入口跳转外链的地址链接，发票外跳的链接会带有发票参数，用于标识是从那张发票跳出的链接
	PromotionUrlSubTitle string `json:"promotion_url_sub_title,omitempty"` // 显示在入口右侧的 tips ，长度限制在 6 个汉字内
}

type CreateCardObj struct {
	BaseInfo *CreateCardBaseInfo `json:"base_info"` // 发票卡券模板基础信息
	Payee    string              `json:"payee"`     // 收款方（开票方）全称，显示在发票详情内。故建议一个收款方对应一个发票卡券模板
	Type     string              `json:"type"`      // 发票类型
}

/*
创建发票卡券模板
通过本接口可以为创建一个商户的发票卡券模板，为该商户配置发票卡券模板上的自定义栏位。创建发票卡券模板生成的card_id将在创建发票卡券时被引用，故创建发票卡券模板是创建发票卡券的基础
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Invoicing_Platform_API_List.html
POST https://api.weixin.qq.com/card/invoice/platform/createcard?access_token={access_token}
*/
func (api *InvoiceApi) PlatformCreateCardRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiPlatformCreateCard, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *InvoiceApi) PlatformCreateCard(param *CreateCardObj) (string, error) {
	payload := struct {
		InvoiceInfo *CreateCardObj `json:"invoice_info"`
	}{
		InvoiceInfo: param,
	}
	result := struct {
		CardID string `json:"card_id"`
	}{}
	err := utils.ApiPostWrapper(api.PlatformCreateCardRaw, payload, &result)
	if err != nil {
		return "", err
	}
	return result.CardID, nil
}

/*
上传PDF
商户或开票平台可以通过该接口上传PDF。PDF上传成功后将获得发票文件的标识，后续可以通过插卡接口将PDF关联到用户的发票卡券上，一并插入到收票用户的卡包中
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Invoicing_Platform_API_List.html
POST https://api.weixin.qq.com/card/invoice/platform/setpdf?access_token={access_token}
*/
func (api *InvoiceApi) PlatformSetpdf(filename string, length int64, content io.Reader) (mediaID string, err error) {

	var resp []byte
	resp, err = api.Client.HTTPUpload(apiPlatformSetpdf, content, "pdf", filename, length)
	if err != nil {
		return
	}

	result := &struct {
		SMediaID string `json:"s_media_id"`
	}{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return
	}
	mediaID = result.SMediaID
	return
}

type InvoiceInsertCardExtItem struct {
	Name  string `json:"name,omitempty"` // 项目的名称
	Price int    `json:"price"`          // 项目的单价
	Num   int    `json:"num,omitempty"`  // 项目的数量
	Unit  string `json:"unit,omitempty"` // 项目的单位，如个
}

type InvoiceInsertCardExt struct {
	NonceStr string `json:"nonce_str"` // 随机字符串，防止重复
	UserCard struct {
		InvoiceUserData *InvoiceInsertCardExtUser `json:"invoice_user_data"` // 用户信息结构体
	} `json:"user_card"` // 用户信息结构体
}

type InvoiceInsertCardExtUser struct {
	Fee                   int                        `json:"fee"`                                // 发票的金额，以分为单位
	Title                 string                     `json:"title"`                              // 发票的抬头
	BillingTime           int                        `json:"billing_time"`                       // 发票的开票时间，为10位时间戳（utc+8）
	BillingNO             string                     `json:"billing_no"`                         // 发票的发票代码
	BillingCode           string                     `json:"billing_code"`                       // 发票的发票号码
	CheckCode             string                     `json:"check_code"`                         // 校验码，发票pdf右上角，开票日期下的校验码
	FeeWithoutTax         int                        `json:"fee_without_tax"`                    // 不含税金额，以分为单位
	Tax                   int                        `json:"tax"`                                // 税额，以分为单位
	SPdfMediaID           string                     `json:"s_pdf_media_id"`                     // 发票pdf文件上传到微信发票平台后，会生成一个发票s_media_id，该s_media_id可以直接用于关联发票PDF和发票卡券。发票上传参考“ 3 上传PDF ”一节
	STripPdfMediaID       string                     `json:"s_trip_pdf_media_id,omitempty"`      // 其它消费附件的PDF，如行程单、水单等，PDF上传方式参考“ 3 上传PDF ”一节
	BuyerNumber           string                     `json:"buyer_number,omitempty"`             // 购买方纳税人识别号
	BuyerAddressAndPhone  string                     `json:"buyer_address_and_phone,omitempty"`  // 购买方地址、电话
	BuyerBankAccount      string                     `json:"buyer_bank_account,omitempty"`       // 购买方开户行及账号
	SellerNumber          string                     `json:"seller_number,omitempty"`            // 销售方纳税人识别号
	SellerAddressAndPhone string                     `json:"seller_address_and_phone,omitempty"` // 销售方地址、电话
	SellerBankAccount     string                     `json:"seller_bank_account,omitempty"`      // 销售方开户行及账号
	Remarks               string                     `json:"remarks,omitempty"`                  // 备注，发票右下角初
	Cashier               string                     `json:"cashier,omitempty"`                  // 收款人，发票左下角处
	Maker                 string                     `json:"maker,omitempty"`                    // 开票人，发票下方处
	Info                  []InvoiceInsertCardExtItem `json:"info,omitempty"`                     //	商品详情结构
}

type InvoiceInsertObj struct {
	OrderID string                `json:"order_id"` // 发票order_id，既商户给用户授权开票的订单号
	CardID  string                `json:"card_id"`  // 发票card_id
	Appid   string                `json:"appid"`    // 该订单号授权时使用的appid，一般为商户appid
	CardExt *InvoiceInsertCardExt `json:"card_ext"` // 发票具体内容
}

type InvoiceInsertResult struct {
	Code    string `json:"code"`    // 发票code
	OpenID  string `json:"openid"`  // 获得发票用户的openid
	UnionID string `json:"unionid"` // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段
}

/*
将电子发票卡券插入用户卡包
本接口由开票平台或自建平台商户调用。对用户已经授权过的开票请求，开票平台可以使用本接口将发票制成发票卡券放入用户的微信卡包中
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Invoicing_Platform_API_List.html
POST https://api.weixin.qq.com/card/invoice/insert?access_token={access_token}
*/
func (api *InvoiceApi) InsertRaw(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiInsert, bytes.NewReader(payload), "application/json;charset=utf-8")
}
func (api *InvoiceApi) Insert(param *InvoiceInsertObj) (*InvoiceInsertResult, error) {
	var result InvoiceInsertResult
	err := utils.ApiPostWrapper(api.PlatformCreateCardRaw, param, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
统一开票接口-开具蓝票
对于使用微信电子发票开票接入能力的商户，在公众号后台选择任何一家开票平台的套餐，都可以使用本接口实现电子发票的开具
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html
POST https://api.weixin.qq.com/card/invoice/makeoutinvoice?access_token={access_token}
*/
func (api *InvoiceApi) MakeOutInvoice(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiMakeOutInvoice, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
统一开票接口-发票冲红
对于使用微信电子发票开票接入能力的商户，在公众号后台选择任何一家开票平台的套餐，都可以使用本接口实现电子发票的冲红
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html
POST https://api.weixin.qq.com/card/invoice/clearoutinvoice?access_token={access_token}
*/
func (api *InvoiceApi) ClearOutInvoice(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiClearOutInvoice, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
统一开票接口-查询已开发票
对于使用微信电子发票开票接入能力的商户，在公众号后台选择任何一家开票平台的套餐，都可以使用本接口实现已开具电子发票的查询
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Vendor_API_List.html
POST https://api.weixin.qq.com/card/invoice/queryinvoceinfo?access_token={access_token}
*/
func (api *InvoiceApi) QueryInvoceInfo(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiQueryInvoceInfo, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
查询已上传的PDF文件
用于供发票PDF的上传方查询已经上传的发票或消费凭证PDF
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Invoicing_Platform_API_List.html
POST https://api.weixin.qq.com/card/invoice/platform/getpdf?action=get_url&access_token={access_token}
*/
func (api *InvoiceApi) PlatformGetpdf(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiPlatformGetpdf, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
更新发票卡券状态
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Invoicing_Platform_API_List.html
POST https://api.weixin.qq.com/card/invoice/platform/updatestatus?access_token={access_token}
*/
func (api *InvoiceApi) PlatformUpdateStatus(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiPlatformUpdateStatus, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
查询报销发票信息
通过该接口查询电子发票的结构化信息，并获取发票PDF文件
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Reimburser_API_List.html
POST https://api.weixin.qq.com/card/invoice/reimburse/getinvoiceinfo?access_token={access_token}
*/
func (api *InvoiceApi) ReimburseGetInvoiceInfo(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiReimburseGetInvoiceInfo, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
批量查询报销发票信息
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Reimburser_API_List.html
POST https://api.weixin.qq.com/card/invoice/reimburse/getinvoicebatch?access_token={access_token}
*/
func (api *InvoiceApi) ReimburseGetInvoiceBatch(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiReimburseGetInvoiceBatch, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
报销方更新发票状态
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Reimburser_API_List.html
POST https://api.weixin.qq.com/card/invoice/reimburse/updateinvoicestatus?access_token={access_token}
*/
func (api *InvoiceApi) ReimburseUpdateInvoiceStatus(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiReimburseUpdateInvoiceStatus, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
报销方批量更新发票状态
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/E_Invoice/Reimburser_API_List.html
POST https://api.weixin.qq.com/card/invoice/reimburse/updatestatusbatch?access_token={access_token}
*/
func (api *InvoiceApi) ReimburseUpdateStatusBatch(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiReimburseUpdateStatusBatch, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
将发票抬头信息录入到用户微信中
调用接口，获取添加存储发票抬头信息的链接，将链接发给微信用户，用户确认后将保存该信息
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Quick_issuing/Interface_Instructions.html
POST https://api.weixin.qq.com/card/invoice/biz/getusertitleurl?access_token={access_token
*/
func (api *InvoiceApi) GetUserTitleUrl(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiGetUserTitleUrl, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取用户抬头（方式一）:获取商户专属二维码，立在收银台
商户调用接口，获取链接，将链接转成二维码，用户扫码，可以选择抬头发给商户
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Quick_issuing/Interface_Instructions.html
POST https://api.weixin.qq.com/card/invoice/biz/getselecttitleurl?access_token={access_token}
*/
func (api *InvoiceApi) GetSelectTitleUrl(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiGetSelectTitleUrl, bytes.NewReader(payload), "application/json;charset=utf-8")
}

/*
获取用户抬头（方式二）:商户扫描用户的发票抬头二维码
商户扫用户“我的—个人信息—我的发票抬头”里面的抬头二维码后，通过调用本接口，可以获取用户抬头信息
See: https://developers.weixin.qq.com/doc/offiaccount/WeChat_Invoice/Quick_issuing/Interface_Instructions.html
POST https://api.weixin.qq.com/card/invoice/scantitle?access_token={access_token}
*/
func (api *InvoiceApi) ScanTitle(payload []byte) (resp []byte, err error) {
	return api.Client.HTTPPost(apiScanTitle, bytes.NewReader(payload), "application/json;charset=utf-8")
}
