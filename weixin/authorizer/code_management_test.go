package authorizer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/lixinio/weixin/test"
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

func initAuthorizer() *Authorizer {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	wxopenOA := NewLite(
		redis, redis,
		test.WxOpenAppid,
		test.WxOpenOAAppid,
	)
	return wxopenOA
}

func TestCommit(t *testing.T) {
	open := initWxOpen()
	templates, err := open.GetTemplateList(context.Background())
	require.Empty(t, err)
	require.NotEmpty(t, templates)

	templateID := templates[len(templates)-1].TemplateID
	require.NotEmpty(t, templateID)

	api := initAuthorizer()
	err = api.CodeCommit(context.Background(), templateID, "{}", "test", "test")
	require.Empty(t, err)
}

func TestGetQrcode(t *testing.T) {
	api := initAuthorizer()
	qrcode, err := api.GetTestQrcode(context.Background(), "")
	require.Empty(t, err)
	require.NotEmpty(t, qrcode)
	// save qrcode to temp file
	file, _ := os.Create("/tmp/qrcode-" + test.WxOpenOAAppid + ".png")
	defer file.Close()

	_, _ = io.Copy(file, bytes.NewReader(qrcode))
	fmt.Printf("save qrcode to a temp file: %s\n", file.Name())
}

func TestSubmitAudit(t *testing.T) {
	api := initAuthorizer()
	_, err := api.CodeSubmitAudit(context.Background(), &AuditParams{})
	require.Empty(t, err)
}

func TestRelease(t *testing.T) {
	api := initAuthorizer()
	err := api.CodeRelease(context.Background())
	require.Empty(t, err)
}
