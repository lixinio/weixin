package authorizer

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/stretchr/testify/require"
)

func TestWxaMpLinkGet(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	info, err := api.WxaMpLinkGet(ctx)
	require.Equal(t, err, nil)
	fmt.Print(info)
}

func TestWxaMpUnLink(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	err := api.WxaMpUnLink(ctx, test.WxOpenOAAppid)
	require.Equal(t, err, nil)
}

func TestWxaMpLink(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	err := api.WxaMpLink(ctx, test.WxOpenOAAppid, "0", "0")
	require.Equal(t, err, nil)
}
