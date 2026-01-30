
####  xkube是一款开源的kubernetes集群管理平台，支持PC端和APP，具备丰富的集群管理、运维监控和CI/CD发布功能，让企业轻松实现跨集群的统一管理和自动化运维。

### [官方网址](https://xkube.eeenet.net/)   |   [APP下载](https://xkube.eeenet.net/app.html)   |   [演示地址](https://xkube.eeenet.net/)   |   [文档地址](https://xkube.eeenet.net/doc/xkube.html)  |  [意见反馈](https://xkube.eeenet.net/feedback.html) 

### QQ交流群：664180281
> 关注公众号了解xkube最新动态

<img src="https://www.eeenet.net/static/myimg/sou_xkube.png" width="100" height="100" alt="关注公众号" title="公众号">


# 一、快速上手

### 1.下载编译

```
go clone https://gitee.com/eeenet/xkube.git
cd xkube
go mod tidy
go build main.go
```

windows 下编译成Linux可执行文件:
```
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -ldflags "-s -w -X 'main.BuildTime=2026-01-23T12:21:25' -X 'main.GoVersion=go version go1.25.4 windows/amd64'" -o xkube-v4.0 main.go
```


### 2.修改配置
> 编辑conf/app.conf,将配置文件中的mysql和redis配置更改成本地环境，

[配置文件各参数说明文档](https://xkube.eeenet.net/doc/page/read.html?id=6)

```
#mysql配置信息
db_host = 127.0.0.1
db_port = 3307
db_user = root
db_pass = root#123
db_name = db_xkube
db_type = mysql

#redis配置信息
redisDb = "127.0.0.1:6379"
redisPasswd = "redis#123"

SessionProviderConfig = "127.0.0.1:6379,100,redis#123"

```

### 3.导入sql文件
> 创建mysql数据库db_xkube并导入 docker-compose/db_xkube.sql，mysql版本为8.0

```
mysql -uroot -p
CREATE DATABASE db_xkube CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
use db_xkube
source db_xkube.sql

```

### 4.启动
> 执行命令：./xkube 

# 二、docker-compose部署

```
 git clone https://gitee.com/eeenet/xkube.git 
 cd xkube/docker-compose
 #启动命令
 docker-compose up -d  或者 docker compose up -d
 #停止
 docker-compose down 或 docker compose down
```



# 三、目录结构说明
> 介绍各目录功能及作用

|目录|功能|
|--|--|
|common|公共类库文件,主要是redis类、client-go类库文件、mysql类库文件及公共函数|
|conf|存放配置文件、各参数说明参考:[参数说明](https://xkube.eeenet.net/doc/page/read.html?id=6)|
|controllers|控制器目录，主要处理http请求和响应的逻辑及数据解析,然后再调用models里的函数进行数据处理|
|models|主要进行数据库的数据操作、或处理k8s api的数据|
|routers|api 路由|
|views|前端页面,html文件放在views\front\page下|
|xadmin|后台管理的包，主要处理后台的权限控制、用户管理、角色管理、日志审计等功能|
|main.go|入口文件|



## xkube插件gitee源码地址
*  [k8s管理助手APP源码](https://gitee.com/eeenet/xkubeApp)
*  [k8s快速部署安装包](https://gitee.com/eeenet/k8s-install)
*  [xkube阿里云云效插件](https://gitee.com/eeenet/XkubeAliyunFlow)
*  [xkube服务检测插件](https://gitee.com/eeenet/xkube-check)


## 帮助文档
- [xkube文档](https://xkube.eeenet.net/doc/xkube.html)
- [容器化部署](https://xkube.eeenet.net/doc/page/read.html?id=4)
- [Linux服务器部署](https://xkube.eeenet.net/doc/page/read.html?id=5)
- [k8s部署](https://xkube.eeenet.net/doc/page/read.html?id=3)
- [部署xkube时,nginx反向代理配置](https://xkube.eeenet.net/doc/page/read.html?id=10)
- [免费ssl证书部署](https://xkube.eeenet.net/doc/page/read.html?id=23)
- [免费ssl证书续签](https://github.com/justjavac/certbot-dns-aliyun)
- [xkube-v3.8使用文档](https://xkube.eeenet.net/doc/page/read.html?id=38)
- [go语言安装部署](https://xkube.eeenet.net/doc/page/read.html?id=39)
- [常见问题](https://xkube.eeenet.net/doc/xkube.html#/page/list.html?xcolumn=%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98)
- [历史版本下载地址](https://github.com/eeenet007/xkube/releases)
## 感谢：
- [layui](https://layui.dev)
- [beego]( https://github.com/beego/beego)
- [layuimini](https://gitee.com/zhongshaofa/layuimini)
- [editor.md](https://github.com/pandao/editor.md)
