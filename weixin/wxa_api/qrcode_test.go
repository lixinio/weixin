package wxa_api

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/lixinio/weixin/utils"
	"github.com/stretchr/testify/require"
)

func saveFile(data []byte, fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(data))
	return err
}

func TestGetWxaCodeUnlimit(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initWxa(),
		// initAuthorizer(),
	} {
		wxaApi := NewApi(client)
		b, err := wxaApi.GetWxaCodeUnlimit(ctx, &GetWxaCodeUnlimitRequest{
			Scene:     "test",
			Page:      "modules/usedcar/Showroom/index",
			CheckPath: false,
			Width:     1024,
			AutoColor: true,
			IsHyaline: true,
		})
		require.Equal(t, nil, err)

		err = saveFile(b, "a.png")
		require.Equal(t, nil, err)
	}
}

func TestGetWxaCode(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initWxa(),
		// initAuthorizer(),
	} {
		wxaApi := NewApi(client)
		b, err := wxaApi.GetWxaCode(ctx, &GetWxaCodeRequest{
			Path:      "modules/usedcar/Showroom/index",
			Width:     10240,
			AutoColor: true,
			IsHyaline: true,
		})
		require.Equal(t, nil, err)

		err = saveFile(b, "b.png")
		require.Equal(t, nil, err)
	}
}

func TestCreateWxaQRCode(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initWxa(),
		// initAuthorizer(),
	} {
		wxaApi := NewApi(client)
		b, err := wxaApi.CreateWxaQRCode(ctx, "modules/usedcar/Showroom/index", 0)
		require.Equal(t, nil, err)

		err = saveFile(b, "c.png")
		require.Equal(t, nil, err)
	}
}
