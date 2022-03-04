package wxopen

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/stretchr/testify/require"
)

func TestFastRegisterBetaWeapp(t *testing.T) {
	open := initWxOpen()
	ctx := context.Background()
	result, err := open.FastRegisterBetaWeapp(ctx, "叮当当", test.WxOpenOAOpenID)
	require.Empty(t, err)
	fmt.Println(result)
}
