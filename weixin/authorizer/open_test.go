package authorizer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/stretchr/testify/require"
)

func TestOpenCreate(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	openid, err := api.WxOpenCreate(ctx, test.WxOpenOAAppid)
	if err == nil {
		fmt.Print(openid)
	} else {
		var wxError *utils.WeixinError
		if !errors.As(err, &wxError) {
			require.Equal(t, err, nil)
		} else if wxError.ErrCode != 89000 { // ErrCode:89000, ErrMsg:"account has bound open
			require.Equal(t, err, nil)
		}
	}
}

func TestOpen(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	openid, err := api.WxOpenGet(ctx, test.WxOpenOAAppid)
	if err != nil {
		var wxError *utils.WeixinError
		if !errors.As(err, &wxError) {
			require.Equal(t, err, nil)
		} else if wxError.ErrCode != 89002 { // ErrCode:89002, ErrMsg:"open not exists
			require.Equal(t, err, nil)
		}
		return
	}

	fmt.Printf("openid : %s\n", openid)

	have, err := api.WxOpenHave(ctx)
	require.Equal(t, err, nil)
	fmt.Println(have)

	fbind := func() {
		err = api.WxOpenBind(ctx, test.WxOpenOAAppid, openid)
		if err != nil {
			var wxError *utils.WeixinError
			if !errors.As(err, &wxError) {
				require.Equal(t, err, nil)
			} else if wxError.ErrCode != 89000 { // ErrCode:89000, ErrMsg:"account has bound open rid
				require.Equal(t, err, nil)
			}
		}
	}

	fbind()

	err = api.WxOpenUnBind(ctx, test.WxOpenOAAppid, openid)
	require.Equal(t, err, nil)

	have, err = api.WxOpenHave(ctx)
	require.Equal(t, err, nil)
	fmt.Println(have)

	fbind()
}

func TestRid(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()
	req, err := api.RidGet(ctx, "622177b1-686aebfa-2d6bb912")
	require.Equal(t, err, nil)
	fmt.Print(req)
}
