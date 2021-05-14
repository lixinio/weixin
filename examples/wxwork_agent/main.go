package main

import (
	"fmt"
	"net/url"

	"github.com/lixinio/weixin/redis"
	"github.com/lixinio/weixin/wxwork"
	"github.com/lixinio/weixin/wxwork/agent"
	"github.com/lixinio/weixin/wxwork/department_api"
	"github.com/lixinio/weixin/wxwork/tag_api"
	"github.com/lixinio/weixin/wxwork/user_api"
)

func main() {
	cache := redis.NewRedis(&redis.Config{RedisUrl: "redis://127.0.0.1:6379/1"})
	corp := wxwork.New(&wxwork.Config{
		Corpid: "wx247d4bc469342dc4",
	})
	agent := agent.New(corp, cache, &agent.Config{
		AgentId: "20",
		Secret:  "G9x8iHpoQMJ8ynDgcplAvwiF4qWF1tRJ3gMVShXZ1Ks",
	})

	userApi := user_api.NewAgentApi(agent)

	params := url.Values{}
	params.Add("department_id", "1")
	b, e := userApi.SimpleList(params)
	fmt.Println(string(b), " ", e)

	departmentApi := department_api.NewAgentApi(agent)
	params = url.Values{}
	params.Add("id", "1")
	b, e = departmentApi.List(params)
	fmt.Println(string(b), " ", e)

	tagApi := tag_api.NewAgentApi(agent)
	b, e = tagApi.List()
	fmt.Println(string(b), " ", e)

	params = url.Values{}
	params.Add("agentid", agent.Config.AgentId)
	b, e = agent.AgentGet(params)
	fmt.Println(string(b), " ", e)
}
