package message_api

import (
	"context"
	"testing"
)

func TestCustomerMessage(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*messageItem{
		initOfficialAccount(),
		initAuthorizer(),
	} {
		messageApi := NewApi(client.Client)
		messageApi.SendCustomTextMessage(ctx, client.OpenID, "发多了开发")
	}
}
