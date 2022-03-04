package authorizer

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAccountBasicInfo(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	info, err := api.GetAccountBasicInfo(ctx)
	require.Equal(t, err, nil)
	fmt.Print(info)
}

func TestGetCategory(t *testing.T) {
	api := initAuthorizer()
	ctx := context.Background()

	info, err := api.GetCategory(ctx)
	require.Equal(t, err, nil)
	fmt.Print(info)
}
