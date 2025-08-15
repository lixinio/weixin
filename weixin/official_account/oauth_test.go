package official_account

import (
	"context"
	"fmt"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/stretchr/testify/require"
)

func TestTicket(t *testing.T) {
	ctx := context.TODO()
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := New(redis, redis, &Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})
	officialAccount.EnableJSApiTicketCache(redis, redis)
	officialAccount.EnableWxCardTicketCache(redis, redis)

	ticket, err := officialAccount.GetJSApiTicket(ctx)
	require.Equal(t, nil, err)
	fmt.Println(ticket)

	ticket, err = officialAccount.GetWxCardApiTicket(ctx)
	require.Equal(t, nil, err)
	fmt.Println(ticket)
}
