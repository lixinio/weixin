package authorizer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxopen"
	"github.com/stretchr/testify/require"
)

func initWxOpen() *wxopen.WxOpen {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxopen := wxopen.New(redis, redis, &wxopen.Config{
		Appid:          test.WxOpenAppid,
		Secret:         test.WxOpenSecret,
		Token:          test.WxOpenToken,
		EncodingAESKey: test.WxOpenEncodingAESKey,
	})
	return wxopen
}

func initAuthorizer() *utils.Client {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxopenOA := NewLite(
		redis, redis,
		test.WxOpenAppid,
		test.WxOpenOAAppid,
	)
	return wxopenOA.Client
}

func TestCommit(t *testing.T) {
	open := initWxOpen()
	templates, err := open.GetTemplateList(context.Background())
	require.Empty(t, err)
	require.NotEmpty(t, templates)

	templateID := templates[len(templates)-1].TemplateID
	require.NotEmpty(t, templateID)

	api := NewApi(initAuthorizer())
	err = api.Commit(context.Background(), templateID, "{}", "test", "test")
	require.Empty(t, err)
}

func TestGetQrcode(t *testing.T) {
	api := NewApi(initAuthorizer())
	qrcode, err := api.GetQrcode(context.Background(), "")
	require.Empty(t, err)
	require.NotEmpty(t, qrcode)
	// save qrcode to temp file
	file, _ := os.CreateTemp("", "qrcode-"+test.WxOpenOAAppid+"-*.png")
	defer file.Close()

	_, _ = io.Copy(file, bytes.NewReader(qrcode))
	fmt.Printf("save qrcode to a temp file: %s\n", file.Name())
}

func TestSubmitAudit(t *testing.T) {
	api := NewApi(initAuthorizer())
	_, err := api.SubmitAudit(context.Background(), map[string]interface{}{})
	require.Empty(t, err)
}

func TestRelease(t *testing.T) {
	api := NewApi(initAuthorizer())
	err := api.Release(context.Background())
	require.Empty(t, err)
}
