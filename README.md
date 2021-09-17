# weixin
微信/企业微信/开放平台Golang实现

基于微信官方文档实现：https://developers.weixin.qq.com/doc/

## 特点

### 将各个子接口分成独立的模块
+ 自建和服务商共享大部分接口， 所以如果和agent/suite等强关联， 势必写两份重复代码
+ 减少模块依赖， 用到多少下载多少， 避免`helloworld`都需要下载整个库

## 快速开始

1. `git clone git@github.com:lixinio/weixin.git`
2. `cp ./weixin/test/config.go.example ./weixin/test/config.go`
3. 修改 `./weixin/test/config.go` 中相应的配置为你的配置信息
4. `cd ./weixin/examples/wxopen` 并执行 `go run .`

## 致谢

- [fastwego](https://github.com/fastwego)，部分实现参考了该项目
