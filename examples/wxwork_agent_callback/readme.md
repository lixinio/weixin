# 网络配置

由于请求来自腾讯服务器， 所以需要公网可访问的域名

为了方便本地调试， 可以考虑用`ngrok`

需要通讯录的回调配置和普通应用的配置保持一致（`Token` / `EncodingAESKey`）

回调Url

+ 普通应用： http://test.kimq.cn/weixin/$CORP_ID/$AGENT_ID
+ 通讯录应用： http://test.kimq.cn/weixin/$CORP_ID/0

> 按需要更换域名
> 
> 按需要更换`$CORP_ID`和`$AGENT_ID`

## 菜单
本实例支持菜单接口回调， 具体菜单配置可见单元测试 `wxwork/agent/app_test.go`