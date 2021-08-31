## 网络配置

将特定域名（例如`lan.kimq.cn`）指向测试机局域网IP（比如`192.168.1.168`)， 无需公网IP

如果80端口冲突，可以nginx代理一下

``` nginx
server{
    listen 80;
    server_name lan.kimq.cn;
    location / {

        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host  $host;
        proxy_set_header X-Nginx-Proxy true;
        proxy_set_header Connection "";
        proxy_pass      http://127.0.0.1:9998 ;
    }
}
```
## 登录测试
> 确保手机和服务器在一个局域网
> 
> 需要在微信中打开


访问 <http://lan.kimq.cn/>， 选择`微信登录`

### [公众号测试平台](https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo)

> 需要先配置 `体验接口权限表` `网页服务` `网页帐号` `网页授权获取用户基本信息`为 `lan.kimq.cn`


### 正式服务号

需要先配置  `设置` `公众号设置` `功能设置` `网页授权域名`为 `lan.kimq.cn` ，微信强制`域名所有权检验`， 所以必须是公网可访问， 并且服务器需要配合响应一个特定的请求


# 回调

> 因为微信的登录必须在微信运行时里面，这通常是一个手机， 手机是不太好配置本地域名映射（如/etc/hosts）， 所以需要在`路由器`或者`公网`将域名（如lan.kimq.cn）指向一个局域网地址
> 
> 如果切换网络导致局域网IP变化又带来诸多不便， 所以条件允许， 建议用ngrok代理到本机

如果要测试回调， 就必须可以公网可访问

在本例中， 实现了
+ 接收文本， 回复相同内容
+ 接收图片， 回复相同内容
+ 接收语音， 回复相同内容
+ 接收视频
+ 接收地址， 回复图文
+ 其他， 具体看代码（examples/official_account/callback.go）