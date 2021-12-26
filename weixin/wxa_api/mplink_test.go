package wxa_api

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/stretchr/testify/require"
)

func TestGetShowWxaItem(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initWxa(),
		// initAuthorizer(),
	} {
		wxaApi := NewApi(client)

		list, err := wxaApi.GetWxaMplinkForShow(ctx, 0, 0)
		require.Equal(t, nil, err)
		fmt.Println(list)

		b, err := wxaApi.GetShowWxaItem(ctx)
		require.Equal(t, nil, err)
		fmt.Println(b)

		err = wxaApi.UpdateShowWxaItem(ctx, 1, test.OfficialAccountAppid)
		require.Equal(t, nil, err)
	}
}
