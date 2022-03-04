package authorizer

import (
	"fmt"
	"testing"
)

func TestGetFastRegisterAuthUrl(t *testing.T) {
	api := initAuthorizer()
	url := api.GetFastRegisterAuthUrl("1", "https://test.lixinchuxing.cn/gateway/component/notify")
	fmt.Print(url)
}
