package wxopen

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatePreAuthCode(t *testing.T) {
	open := initWxOpen()
	code, expiresIn, err := open.CreatePreAuthCode(context.Background())
	fmt.Println(code, expiresIn)
	require.Empty(t, err)
}

func TestGetAuthorizerList(t *testing.T) {
	open := initWxOpen()
	details, err := open.GetAuthorizerList(context.Background(), 0, 10)
	require.Empty(t, err)
	for _, detail := range details {
		fmt.Printf(
			"%s %d %s\n",
			detail.AuthorizerAppid,
			detail.AuthTime,
			detail.AuthorizerRefreshToken,
		)
	}
}
