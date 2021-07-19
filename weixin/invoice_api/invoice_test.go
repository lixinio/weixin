package invoice_api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
)

func newInvoiceApi() *InvoiceApi {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})

	return NewOfficialAccountApi(officialAccount)
}

func TestInvoiceUploadPdf(t *testing.T) {
	api := newInvoiceApi()

	file, err := os.Open(test.InvoicePdf)
	require.Empty(t, err)
	defer file.Close()

	fi, err := file.Stat()
	require.Empty(t, err)

	mediaID, err := api.PlatformSetpdf("fapiao.pdf", fi.Size(), file)
	require.Equal(t, nil, err)
	fmt.Printf("media id %s\n", mediaID)
}

func TestSetContact(t *testing.T) {
	api := newInvoiceApi()

	setbizattrObj := &SetbizattrObj{
		Phone:   test.InvoicePhone,
		TimeOut: 7200,
	}

	err := api.SetContact(setbizattrObj)
	require.Equal(t, nil, err)

	result, err := api.GetContact()
	require.Equal(t, nil, err)
	require.Equal(t, result.Phone, setbizattrObj.Phone)
	require.Equal(t, result.TimeOut, setbizattrObj.TimeOut)
}

func TestPlatformCreateCard(t *testing.T) {
	api := newInvoiceApi()

	cardID, err := api.PlatformCreateCard(&CreateCardObj{
		Payee: test.InvoicePayee,
		Type:  test.InvoiceType,
		BaseInfo: &CreateCardBaseInfo{
			Title:                test.InvoiceCustomUrlName,
			CustomUrlName:        test.InvoiceCustomUrlName,
			CustomURL:            test.InvoiceCustomURL,
			CustomUrlSubTitle:    test.InvoiceCustomUrlSubTitle,
			PromotionUrlName:     "查看其他",
			PromotionURL:         "https://www.baidu.com",
			PromotionUrlSubTitle: "详情",
			LogoUrl:              "https://mmbiz.qpic.cn/mmbiz_png/5tTVBJAGiap2TWlw0pPpbVtE80xH4sUs4u1aPZOlKHgPNS3sKm1CpJM3aLKd36yLreXqAHenD3q8QU3Hovpjv0g/0",
		},
	})
	require.Equal(t, nil, err)
	fmt.Printf("card id %s\n", cardID)
}

func TestInvoiceInsert(t *testing.T) {
	api := newInvoiceApi()

	billingTime := 0
	{
		layout := "2006-01-02"
		tm, err := time.Parse(layout, "2021-06-23")
		require.Equal(t, nil, err)
		billingTime = int(tm.Unix())
	}

	param := &InvoiceInsertObj{
		OrderID: "1624612433713210184",
		CardID:  "p-mcP1FC6QHZ515goRP3CsXZcXmI",
		Appid:   test.OfficialAccountAppid,
		CardExt: &InvoiceInsertCardExt{
			NonceStr: fmt.Sprintf("%d", time.Now().UnixNano()),
			UserCard: struct {
				InvoiceUserData *InvoiceInsertCardExtUser `json:"invoice_user_data"`
			}{
				InvoiceUserData: &InvoiceInsertCardExtUser{
					Fee:           10,
					Title:         "邱金武",
					BillingTime:   billingTime,
					BillingNO:     "044032000211",
					BillingCode:   "62522141",
					CheckCode:     "85073690672647647833",
					FeeWithoutTax: 9,
					Tax:           1,
					SPdfMediaID:   "71381497449443328",
					// Cashier:               "方婷",
					// Maker:                 "朱芷娆",
					// SellerNumber:          "91440300MA5EKEKU9G",
					// SellerBankAccount:     "中国银行华润城支行 774469529787",
					// SellerAddressAndPhone: "深圳市宝安区西乡街道蚝业社区宝安互联网产业基地A区1栋3B08 0755-26406220",
					// Info: []InvoiceInsertCardExtItem{
					// 	{
					// 		Name:  "*信息技术服务*平台服务费",
					// 		Price: 10,
					// 		Num:   1,
					// 		Unit:  "次",
					// 	},
					// },
				},
			},
		},
	}

	b, _ := json.Marshal(param)
	fmt.Println(string(b))

	result, err := api.Insert(param)
	require.Equal(t, nil, err)
	fmt.Printf("code : %s, openid: %s, unionid: %s\n", result.Code, result.OpenID, result.UnionID)

}

func TestRejectInsert(t *testing.T) {
	api := newInvoiceApi()

	err := api.RejectInsert(&RejectInsertObj{
		OrderID: "1624605258318629788",
		SPappID: "d3gxMTY5NGJiNDI4YTMyZTg4X0jdlhfLZft3pZEI0pLVYp3CRPzlu2kW_06OUzJGyaZ3",
		Reason:  "就是不开",
	})
	require.Equal(t, nil, err)
}

func TestSetAuthField(t *testing.T) {
	api := newInvoiceApi()
	param := &AuthFieldObj{
		UserField: &AuthUserField{
			ShowTitle:    1,
			ShowPhone:    1,
			ShowEmail:    0,
			RequirePhone: 1,
			RequireEmail: 0,
		},
		BizField: &AuthBizField{
			ShowTitle:       1,
			ShowTaxNO:       1,
			ShowAddr:        1,
			ShowPhone:       1,
			ShowBankType:    1,
			ShowBankNO:      1,
			RequireTaxNO:    1,
			RequireAddr:     0,
			RequirePhone:    0,
			RequireBankType: 0,
			RequireBankNO:   0,
		},
	}
	err := api.SetAuthField(param)
	require.Equal(t, nil, err)
}

func TestInvoice(t *testing.T) {
	api := newInvoiceApi()

	// {
	// 	result, err := api.GetAuthData(&AuthDataObj{
	// 		OrderID: "1623930786748654309",
	// 		SPappID: "d3g1OTg2NGE5ZTU3ODIyOWVhX_oiY7-5OuzNHme3fHyMQQWjstgWqHfPcktQ40c-H73D",
	// 	})
	// 	require.Equal(t, nil, err)
	// 	fmt.Println(result.InvoiceStatus)
	// }
	///////////////////////////// debug

	spappID := ""
	{
		result, err := api.SetUrl()
		require.Equal(t, nil, err)
		fmt.Println(result)

		u, err := url.Parse(result)
		require.Equal(t, nil, err)
		m, err := url.ParseQuery(u.RawQuery)
		require.Equal(t, nil, err)
		pappid, ok := m["s_pappid"]
		require.Equal(t, true, ok)
		require.NotEmpty(t, pappid)
		spappID = pappid[0]
		fmt.Printf("s_pappid : %s\n", spappID)
	}

	orderID := fmt.Sprintf("%d", time.Now().UnixNano())
	fmt.Printf("order id %s\n", orderID)
	{
		ticket, _, err := api.OfficialAccount.GetWxCardApiTicket()
		require.Equal(t, nil, err)

		result, err := api.GetAuthUrl(&AuthUrlObj{
			SPappID:   spappID,
			Money:     10,
			Source:    "web",
			OrderID:   orderID,
			Timestamp: time.Now().Unix(),
			Type:      1,
			Ticket:    ticket,
		})
		require.Equal(t, nil, err)
		fmt.Printf("url : %s , appid %s\n", result.AuthURL, result.AppID)
	}

	{
		result, err := api.GetAuthData(&AuthDataObj{
			OrderID: orderID,
			SPappID: spappID,
		})
		require.Equal(t, nil, err)
		fmt.Println(result.InvoiceStatus)
	}
}
