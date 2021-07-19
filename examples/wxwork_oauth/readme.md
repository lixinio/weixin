## 网络配置

将特定域名指向测试机局域网IP（比如`192.168.1.168`)， 无需公网IP

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

## 网页扫码登录

> 需要先配置 `企业微信授权登录` `Web网页` `授权回调域` 为 `lan.kimq.cn`

访问 <http://lan.kimq.cn/>， 选择`网页扫码登录`

## 企业微信认证
> 在企业微信中打开
>
> 需要先配置 `网页授权及JS-SDK` `可信域名` 为 `lan.kimq.cn`, 无需`域名所有权检验`


访问 <http://lan.kimq.cn/>， 选择`企业微信登录`

> 确保手机和服务器在一个局域网
