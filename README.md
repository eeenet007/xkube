[TOC]
# 一.项目介绍
  一款永久免费且无任何功能限制的云原生kubernates多集群管理工具，提供PC端和移动端全平台支持，具备完善的集群管理、运维监控和CI/CD发布功能，让企业轻松实现跨集群的统一管理和自动化运维。

# 二.源代码说明

1. 本项目中前端部分提供了源代码，后端部分均采用源代码编译的二进制。
2. 该项目旨在进行关于kubernates、client-go、go语言开发的学习与技术交流。
3. 本人郑重承诺不做任何功能限制，永久免费使用、无广告、不用注册、不用授权、无后门、代码无加密、可二次开发、无法律限制、无附加条款。
4. 关于后端源代码的开源问题，作如下说明:
- 本项目会持续更新迭代，持续发布windows、Linux的二进制程序及相关更新信息。
- 为避免不尊重他人劳动成果、能理解作者辛苦付出并表示支持的，源代码采取有偿提供。 
5. 如果需要定制开发，发邮件至`eeenet@qq.com`。

# 三.源代码有偿提供说明如下：
1. 有偿方式为支付100元人民币，作者提供当前最新版本的源代码。
2. 购买源代码为虚拟物品交易，购买后不支持退换货。
3. 仅为源代码交易，不提供任何售后服务和技术支持，作者会视情况免费提供一些基本的技术指导。
4. 用户在购买源代码以后，在成交后一年内有新的版本迭代、免费提供最新的源代码。
5. 该项目及代码仅作为学习、技术交流、不支持二次销售或包装成商业产品销售。
6. 用户在使用该项目过程中所产生的任何问题和后果，作者不承担任何责任。
7. 有兴趣购买源代码，微信联系作者或发邮件至：`eeenet@qq.com`。

![作者微信](https://gitee.com/eeenet/xkube/raw/master/images/weixin.png)

# 四.安全建议：
1. 不建议将该项目部署环境暴漏在公网。若需部署在公网环境建议采取https、IP白名单、nginx的auth_basic、jwt-token等方式进行二次安全校验。
2. 作者无法保证该项目100%无漏洞，所以在使用过程中请自行进行网络安全的规划。
3. 尽快修改项目中配置的数据库密码、后台默认密码。
4. 作者采用的安全措施有：启用https、IP白名单、在api的调用、app调用的过程中也是采用token的二次校验。
5. 后续会继续完善手机验证码校验、app安全码校验等功能。


# 五.功能特性
- 跨公有云、跨IDC的多k8s集群统一管理平台
- 具备节点、节电池、命名空间、clusterRoleBinding、clusterRoles、RoleBinding、Roles、serviceAccounts的创建、详情、yaml文件查看、删除等功能。
- workload方面支持对deployment、statefuleset、dameset、cronjob、job、pod容器组、cdr自定义资源、hpa伸缩的功能创建、yaml查看修改、删除功能。
- 无状态【deployment]:功能具备yaml在线编辑、yaml下载,在线修改升级策略,在线标签修改,在线重启,镜像更新，查看关联的pod对象、查看关联的service、ingress，创建hpa对deployment进行自动伸缩容，可以在线操作回滚到指定的镜像版本，查看该deployment的相关事件，在线查看日志，ssh终端登录关联的pod。
- 有状态【statefulset】：功能和deployment的类似，除有少量差异之外，相关功能基本一致。
- 守护进程【daemonset】的功能主要是对deamonset的相关信息进行分类查看。
- 任务【job】：具备日志、信息、事件、状态的查询功能。
- 定时任务【cronjob】：在线通过图形库界面进行创建、yaml文件进行创建、对计划任务的在线修改，在线更改状态等功能。
- 容器组【pod】：具备ssh登录、日志查看、实时查看pod的内存、cpu使用情况【需k8s环境安装metric-beat】等功能。
- 扩缩容【hpa】：在线图形化操作，根据pod的cpu、内存使用情况、定义pod的扩缩容。
- 自定义资源【cdr】:自定义资源的信息查看
- yaml操作：可以通过在线的各种deployment、service、ingress、cronjob等yaml文件模板来进行资源的创建。
- 服务【service】:支持通过yaml来进行创建，对service的yaml配置查看和修改、支持图形操作创建。
- 路由【ingres】：支持通过yaml来进行创建，对service的yaml配置查看和修改、支持图形操作创建，目前只支持nginx-ingress。
- 配置【configmap】：支持图形化、yaml配置的查看、创建、修改和删除。
- 保密字典【secret】:支持图形化、yaml配置的查看、创建、修改和删除。
- 存储声明【pvc】:pvc的yaml查看、创建、修改、删除。
- 存储卷【pv】:存储卷信息的yaml查看、创建、修改、删除。
- 存储类【storageclass】:存储类信息的yaml查看、创建、修改、删除。
- 事件信息：查看当前集群中发生的事件信息。
- 应用集：按照资源的标签appname=myapp进行划分，将该项目所涉及的资源整合到统一界面便于管理。
- 权限管理：按照角色进行权限划分：超级管理员、普通管理员、只读等角色，并按照集群进行授权，只有授权的了对应集群权限的用户才能访问该集群的资源。
- 文档中心：markdown格式的文档编辑器、用于运维文档记录。
- CICD:支持对接阿里云流水线，可以试下流水线的运行、信息查看、流程查看、日志查看等功能。
- aws的eks集群管理：eks的管理需要进行通过aws进行认证才能进行管理，可以通过bearerToken实现对接管理。

# 六.安装部署

## 1.docker-compose部署

将docker-compose拷贝到服务器后，进入docker-compose目录，执行启动命令：`docker-compose up -d` 然后浏览器打开：http://ip:8080/
输入用户名:admin,密码:admin 进行登录，登录后添加k8s配置即可进行管理。停止命令：`docker-compose down`
作者运行环境版本参考：
- docker版本：version 27.3.1
- docker-compose版本： version v2.32.4

## 2.k8s环境部署

上传k8s_deploy.yaml到服务器,修改k8s_deploy.yaml中的mysql和redis的配置，以及ingress的域名，修改以后执行命令：`kubectl apply -f k8s_deploy.yaml`,部署完成以后将域名解析到k8s的ingress入口IP后，在浏览器打开刚解析的域名。用户名和密码均是admin,登录以后会提示：读取集群列表失败、读取集群信息失败。这是因为还没有添加k8s集群，添加后就正常了。

> mysql配置更改:

```bash
	#数据库配置时需要导入docker-compose/xkube.sql
    db_host = mysql_server #更改成mysql的IP
    db_port = 3306			#更改成mysql的端口
    db_user = db_xkube		#更改成mysql的用户名
    db_pass = xkube_DB_123456		#更改成mysql的密码
    db_name = db_xkube		#更改成mysql的数据库名
    db_type = mysql
```

>   #redis配置更改

```bash
    redisDb = "redis_server:6379"	#更改成redis的IP和端口
    redisPasswd = "Redis_123456"	#更改成redis的密码

    SessionProviderConfig = "redis_server:6379,100,Redis_123456" #更改成redis的IP和端口、密码
```

> #域名配置更改

```yaml
spec:
  ingressClassName: nginx
  rules:
  - host: xkube.xxxx.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: xkube-admin
            port: 
              number: 80
```

## 3.服务器上安装部署

修改conf/app.conf中的mysql、redis配置后再进行启动

### 第一步：修改mysql配置

	安装完mysql以后，将docker-compose/xkube.sql导入到数据库中，然后修改conf/app.conf中的如下配置：
	db_host = mysql_server   #mysql的IP
	db_port = 3306	#MySQL的端口
	db_user = db_xkube	#mysql的用户名
	db_pass = xkube_DB_123456 	#mysqld 密码
	db_name = db_xkube	#数据库名
	
### 第二步：修改redis配置

	安装完redis以后：然后修改conf/app.conf中的如下配置：
	redisDb = "192.168.1.115:6379"	#redis的IP和端口
	redisPasswd = "Redis_123456"	#redis配置设置的密码
	SessionProviderConfig = "192.168.1.115:6379,100,Redis_123456"	#将redis的IP、端口、密码进行更改，100这个数字保留即可。
	
### 第三步：启动服务

	修改完以上配置后：Linux 环境下进入目录直接执行：nohub ./xkube & ，windows环境下执行./xkube.exe即可启动
	就可以通过http://ip:8080/index 进行访问了。

### 第四步：前端若有nginx反向代理需增加如下配置：
    location ~^/xkube/pod/terminal/ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Origin "";
    }

### 第五步：后台设置
    1. 登录后台后，默认用户名和密码都是admin，在k8s列表里进行添加需要管理的集群。
    2. 当有多个集群时，可以为自己设置一个常用集群，该设置会记录到cookie，下一次登录会继续管理该集群。
    3. 完成以上两步以后就可以在线管理k8s 了。

### 第六步：后台主面板功能说明

![基本信息](https://gitee.com/eeenet/xkube/raw/master/images/0-1.png)


# 七.管理平台使用说明文档
> 使用说明文档中的截图非最新版截图，文档更新会有滞后
## 1.k8s管理
### 1.1.登录后首页板块说明
1.快捷入口：略
2.资源统计：该部分统计是当有人访问过某个集群对应的资源列表的同时，就会将对应资源的数量记录到redis，然后在通过接口读取出来，该处资源数量非实时更新，当有资源变动以后，需要访问一次对应的资源列表，此处的数值才会更新。
3.最新文章：略
4.集群事件：该处的数据是当用户设置了常用集群以后，读取当前k8s集群中的事件。

### 1.2.添加k8s集群配置

| 字段 | 描述 | 是否必须 |
| --- | --- | --- |
| 集群ID | 集群的简短描述，支持的英文数字下划线或横线，例如:xkube-a1 | 是 |
| k8s版本 | xkube会通过该版本号采用不同的api版本来读取资源，务必正确 | 是 |
| 集群名称 | 中文描述 | 是 |
| 内网IP入口 | http服务的内网IP | 否 |
| 公网入口IP | http服务的公网入口IP | 否 |
| IDC | 集群所属IDC或公有云的描述，例如：aliyun | 否 |
| 备注 | 该集群的简短描述 | 否 |
| 状态 | 记录该集群的状态 | 否 |
| BearerToken | 此功能主要是针对需要进行额外的认证才能读取api的情况,例如aws的eks，eks的bearerToken创建参考“1.5.aws的eks怎样创建不过期token” | 否 |
| kubeconfig | 调用k8s api必须的配置文件，通常在/root/.kube/config获取 | 是 |


### 1.3.设为常用集群
> 功能描述：在多个k8s集群的情况下，每个人可能管理的k8s集群会不同，可以通过设置该功能，将信息记录到cookie,下一次登录以后默认就可以管理设置好的集群。如果出现每个k8s都需要经常管理的情况，可以直接通过在资源列表右上角的集群进行无缝切换。
### 1.4.wiki
> 功能描述：wiki功能模块采用markedown格式来进行记录，可以将一些资源的维护或信息记录到wiki文档中。
### 1.5.aws的eks怎样创建不过期token

> 第一步：kubectl -n kube-system create serviceaccount kubeconfig-sa

> 第二步：kubectl create clusterrolebinding kubeconfig-sa-binding --clusterrole=cluster-admin --serviceaccount=kube-system:kubeconfig-sa

> 第三步：vim编写secret.yaml

```
apiVersion: v1
kind: Secret
metadata:
  name: kubeconfig-sa-token
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: kubeconfig-sa
type: kubernetes.io/service-account-token
```
> 第四步：kubectl apply -f secret.yaml -n kube-system

> 第五步：TOKEN=$(kubectl -n kube-system get secret kubeconfig-sa-token -o jsonpath='{.data.token}'| base64 --decode)

> 第六步：输出token的值：echo $TOKEN 

## 2.集群信息
### 2.1.节点池
> 功能描述：该功能主要针对阿里云的 容器服务ACK的功能，ack里有一个节电池，主要是将ack的节电池列表显示出来。
### 2.2.节点
#### 2.2.1.节点详情
> 节点详情中包含几个主要功能：1.节点的基本信息 2.内存和cpu的当前使用情况[需要k8s安装metric-server] 3.该节点的容器列表 4.yaml查看、节点排水、调度设置、将该节点移除k8s集群 详情参考如下图：
![基本信息](https://gitee.com/eeenet/xkube/raw/master/images/2.2.1-1.png)
![内存cpu监控](https://gitee.com/eeenet/xkube/raw/master/images/2.2.1-2.png)
![节点容器列表](https://gitee.com/eeenet/xkube/raw/master/images/2.2.1-3.png)

### 2.3.命名空间
#### 2.3.1.添加命名空间
> 图形界面创建命名空间：

![创建命名空间](https://gitee.com/eeenet/xkube/raw/master/images/2.3.1-1.png)
> yaml创建命名空间：

```
apiVersion: v1
kind: Namespace
metadata:
  labels:
    appname: my-app
  name: my-namespace
```
#### 2.3.2.命名空间资源限制
> 资源限制：该限制是将该命名空间下的资源总的cpu和内存资源限制，设置时注意格式:cpu:1核=1000m,可以设置为1、100m、500m，内存格式：200Mi,1Gi

![资源限制](https://gitee.com/eeenet/xkube/raw/master/images/2.3.2-1.png)

#### 2.3.3.命名空间yaml编辑
> 命名空间yaml文件的修改，在建立以后，一般是修改标签。

![yaml编辑](https://gitee.com/eeenet/xkube/raw/master/images/2.3.3-1.png)

### 2.4.clusterRoleBinding
> 该功能主要是将集群资源clusterRoleBinding列出来,主要是作为信息的查看。平时运维很少用到

![clusterRoleBinding](https://gitee.com/eeenet/xkube/raw/master/images/2.4-1.png)

### 2.5.clusterRoles
> 该功能主要是将集群资源clusterRoles列出来,主要是作为信息的查看。平时运维很少用到

![clusterRoles](https://gitee.com/eeenet/xkube/raw/master/images/2.5-1.png)

### 2.6.RoleBinding
> 该功能主要是将集群资源RoleBinding列出来,主要是作为信息的查看。平时运维很少用到

![RoleBinding](https://gitee.com/eeenet/xkube/raw/master/images/2.6-1.png)

### 2.7.Roles
> 该功能主要是将集群资源clusterRoleBinding列出来,主要是作为信息的查看。平时在操作权限时会用到

![Roles](https://gitee.com/eeenet/xkube/raw/master/images/2.7-1.png)

### 2.8.ServiceAccount
> 该功能主要是将集群资源ServiceAccount列出来,主要是作为信息的查看。平时在操作权限时会用到

![ServiceAccount](https://gitee.com/eeenet/xkube/raw/master/images/2.8-1.png)


## 3.工作负载
### 3.1.无状态[deployment]
#### 3.1.1.deployment列表
> 该列表页具备在顶部搜索信息框内实现关键字、命名空间对deployment的搜索，同时可以在右上角进行k8s集群的切换。表格右上方支持对表格的列筛选、表格导出、打印功能。表格的最后一列有对deployment进行操作的快捷键，详情页/cicd/yaml编辑可以参考下图，其中CICD是采用本地流水线或调用阿里云流水线API来实现持续集成和持续部署[前提需要在本地或阿里云进行配置好流水线]，重启按钮是将该deploy下的pod重启，删除是将该deployment删除掉相当于命令：kubectl delete

![列表页](https://gitee.com/eeenet/xkube/raw/master/images/3.1.1-1.png)
![详情页](https://gitee.com/eeenet/xkube/raw/master/images/3.1.1-2.png)
![CICD](https://gitee.com/eeenet/xkube/raw/master/images/3.1.1-3.png)
![YAML编辑](https://gitee.com/eeenet/xkube/raw/master/images/3.1.1-4.png)

#### 3.1.2.界面创建deployment
> deployment的创建必须填的参数：集群选择、名称、镜像地址默认读取当前设置的常用k8s集群
1.命名空间：可以为空，为空时命名空间为default
2.标签：可以为空，可以填写多个key-value值，多个时点击后边的加号按钮增加输入框。
附加配置部分：
1.资源限制：当勾选资源限制复选框时，会显示资源限制的输入框,按照格式输入即可,1核=1000m,1GiB=1024Mi,1Gi即为1GB的内存,200Mi=200MB内存
2.健康检查：当勾选健康检查时显示健康检查的输入框,就绪检查readinessProbe,检测类型：URL检测、端口存活检测、执行命令。路径[path]:检查的URL，端口[port]：探测的端口，延迟探测[initialDelaySeconds]:启动pod后延迟多少秒开始检测,探测频率[periodSeconds]:每隔多少秒检测一次,正常阀值[successThreshold]:检测成功多少次算成功,失败阀值[failureThreshold]:失败多少次算失败，超时[timeoutSeconds]:检测的超时时间
3.生命周期:启动后执行[postStart]:在pod启动前先执行命令，停止前执行[preStop]:在pod停止前执行命令
4.环境变量[env]：以key-value形式添加一个或多个变量。
5.service:在创建deployment的同时创建同名的service服务

![deploy创建](https://gitee.com/eeenet/xkube/raw/master/images/3.1.2-1.png)

#### 3.1.3.yaml创建deployment
> 通过yaml配置创建deployment，自己提供配置或选择已有模板进行修改，参考下图：

![yaml创建](https://gitee.com/eeenet/xkube/raw/master/images/3.1.3-1.png)

#### 3.1.4.副本伸缩
> 通过更改副本的数量，来更改pod的数量,参考如下配置,当deployment有设置hpa时，可能导致pod数量冲突，更改无法生效。

```
spec:
  replicas: 2
```
![pod副本修改](https://gitee.com/eeenet/xkube/raw/master/images/3.1.4-1.png)
#### 3.1.5.重新部署
>该功能主要是将服务进行重启,将原有pod销毁，重新创建新的pod.和列表页的重启是同一功能
命令参考：kubectl rollout restart -n test-namespace deployment test-app

![重新部署](https://gitee.com/eeenet/xkube/raw/master/images/3.1.5-1.png)

#### 3.1.6.升级策略
> 升级策略的调整主要调整如下几个参数：

```
  minReadySeconds: 10  #最小准备时间
  strategy:
    rollingUpdate:
      maxSurge: 25%     #超过pod的数量或百分比
      maxUnavailable: 25%       #最大不可用pod的数量或百分比
    type: RollingUpdate         #策略类型：滚动更新RollingUpdate 和替换升级：Recreate
```
![策略调整](https://gitee.com/eeenet/xkube/raw/master/images/3.1.6-1.png)

#### 3.1.7.标签修改
> 增删改标签labels。点击加号增加key-value输入框，点击减号删除一对标签。key和value只支持英文字母和数字及/,-,_。其他特殊字符不支持。
```
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  creationTimestamp: "2023-12-19T09:52:56Z"
  generation: 2
  labels:
    app: test-kk
    appname: test-kk
```
![标签修改](https://gitee.com/eeenet/xkube/raw/master/images/3.1.7-1.png)

#### 3.1.8.镜像更新
> 镜像更新主要是更新参数中的image和name。id不用更改
```
    spec:
      containers:
      - image: nginx:latest
        imagePullPolicy: Always
        name: test-kk
```
![镜像升级](https://gitee.com/eeenet/xkube/raw/master/images/3.1.8-1.png)

#### 3.1.9.创建HPA
>hpa配置参考如下,主要作用在于自动对资源进行扩缩容

```
spec:
  maxReplicas: 5    #最大pod数量
  metrics:
  - resource:
      name: memory
      target:
        averageUtilization: 90  #当内存资源90%时触发扩容事件
        type: Utilization
    type: Resource
  - resource:
      name: cpu
      target:
        averageUtilization: 90  #当cpu资源达到90%时触发扩容事件
        type: Utilization
    type: Resource
  minReplicas: 2    #最小pod数量，
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: test-kk
```
![HPA列表](https://gitee.com/eeenet/xkube/raw/master/images/3.1.9-1.png)
![创建HPA](https://gitee.com/eeenet/xkube/raw/master/images/3.1.9-2.png)

#### 3.1.10.回滚配置
>通过该点击下图中按钮"回滚到此版本"将deployment回滚到之前的版本,yaml详情按钮可以查看历史版本的配置。

```
spec:
  progressDeadlineSeconds: 600
  replicas: 2
  revisionHistoryLimit: 10  #此参数控制保留配置的历史版本数量
```
![历史版本](https://gitee.com/eeenet/xkube/raw/master/images/3.1.10-1.png)


### 3.2.有状态[statefulset]
>有状态的所有管理功能和deployment的管理功能一致，可参考3.1.无状态的帮助文档

![sts列表](https://gitee.com/eeenet/xkube/raw/master/images/3.2.1-1.png)

### 3.3.守护进程集[daemonset]
>daemonset部分功能会少很多，各个按钮的含义和3.1.无状态的功能点一致。具体功能参考3.1.无状态

![ds列表](https://gitee.com/eeenet/xkube/raw/master/images/3.3.1-1.png)

### 3.4.任务[job]
> job主要是由cronjob定时任务触发创建的任务，可以通过job详情查看job具体的执行情况，并通过job关联的pod查看执行结果日志以及yaml的配置信息。

![job列表](https://gitee.com/eeenet/xkube/raw/master/images/3.4.1-1.png)
![job详情](https://gitee.com/eeenet/xkube/raw/master/images/3.4.1-2.png)

### 3.5.定时任务[cronjob]
#### 3.5.1.cronjob列表
> 列表中的 暂停：true时表示该任务已停止运行，false时会定期执行，计划：表示执行运行任务的时间点，参考Linux的crontab定时任务语法，活跃：当任务在运行中是值为1表示，否则值为0.列表中最后一列的启动和暂停按钮用于启动或停止任务，对应修改“suspend:”的值，修改后的值对应列表的暂停这一列的值。CICD：主要用于通过调用流水线来更新镜像。

![cronjob列表](https://gitee.com/eeenet/xkube/raw/master/images/3.5.1-1.png)

#### 3.5.2.界面创建cronjob
> 必填的几个参数：当前集群、命名空间、名称、镜像地址、定时规则。各项参数对应的yaml配置,参考3.5.4.cronjob的yaml语法解释

![cronjob创建](https://gitee.com/eeenet/xkube/raw/master/images/3.5.2-1.png)

#### 3.5.3.定时任务时间设置语法
> 和Linux crontab的语法相同：* * * * * 分别对应分钟（0-59）、小时（0-23）、日期（1-31）、月份（1-12）、星期（0-7）,
*：星号代表所有可能的值；
,：逗号用于指定多个值；
-：横线用于表示一个范围值；
/：斜杠用于表示重复的频率。
示例：
\* * * * * :每分钟执行一次；
0 * * * * :每小时执行一次；
0 0 * * * :每天午夜执行一次；
0 0 * * 1 :每周一午夜执行一次；
0 0 1 * * :每月1号午夜执行一次。
3,15 * * * * :每小时的第3和第15分钟执行一次
3,15 8-11 * * * :每天的8点到11点的第3和第15分钟执行一次
3,15 8-11 */2 * * :每隔两天的上午8点到11点的第3和第15分钟执行一次
3,15 8-11 * * 1 :每周一上午8点到11点的第3和第15分钟执行一次
30 21 * * * :每晚的21:30执行一次
45 4 1,10,22 * * :每月1、10、22日的4 : 45执行一次
10 1 * * 6,0 :每周六、周日的1 : 10执行一次
0,30 18-23 * * * :每天18 : 00至23 : 00之间每隔30分钟执行一次
0 23 * * 6 :每星期六的晚上23: 00 pm执行一次
* */1 * * * :每一小时执行一次
* 23-7/1 * * * :每天晚上23点到第二天7点之间，每隔一小时执行一次

![cron语法建](https://gitee.com/eeenet/xkube/raw/master/images/3.5.3-1.png)

#### 3.5.4.cronjob的yaml语法解释
```yaml
#curl用法
apiVersion: batch/v1
kind: CronJob
metadata:
  name: testapp-task-test   #cronjob名称
  labels:
    app: testapp-task-test  #标签
    appname: myapp          #标签
  namespace: test-app  #命名空间
spec:
  schedule: "*/1 * * * *"
  suspend: true #是否暂停
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: testapp-task-test #containers名称
            image: xxx-registry-vpc.cn-shenzhen.cr.aliyuncs.com/xxx/curl:latest #镜像地址
            imagePullPolicy: IfNotPresent   #镜像策略
            args:   #启动时执行的命令及参数
            - "/bin/sh"
            - "-ec"
            - "curl -o /dev/stdout --max-time 300 --connect-timeout 10 -H 'host:baidu.com' http://svcName.test-app.svc.cluster.local:8001/test.jsp"
          restartPolicy: OnFailure  #重启策略,仅支持Never或OnFailure
```

### 3.6.容器组[pod]
#### 3.6.1.pod列表
> pod列表默认读取所有命名空间的pod。pod列表中的值非实时刷新，最新的值需要手动刷新页面。针对pod的监控信息、pod日志查看、pod终端可以在操作这一列点击进入

![pod列表](https://gitee.com/eeenet/xkube/raw/master/images/3.6.1-1.png)

### 3.6.2.查看pod详情
> pod详情页将pod 的基本信息、状态、存储挂载、事件等信息展示出来，参考下图：

![pod详情](https://gitee.com/eeenet/xkube/raw/master/images/3.6.2-1.png)
![pod详情](https://gitee.com/eeenet/xkube/raw/master/images/3.6.2-2.png)

#### 3.6.3.查看pod日志
> pod的日志查看命令等同：kubectl log pod xxxx -n namespace查看pod的日志。

![pod日志](https://gitee.com/eeenet/xkube/raw/master/images/3.6.3-1.png)

#### 3.6.4.登录容器终端
> 在pod列表中的终端按钮、或者详情中的终端按钮进入pod终端。pod终端是基于websocket实现，在安装配置时可以参考：快速入手的第四步。

![pod终端](https://gitee.com/eeenet/xkube/raw/master/images/3.6.4-1.png)

#### 3.6.5.查看pod的资源使用情况
> pod的资源查看是通过点击列表，操作这一列中的第二个图标按钮进行查看。每个两秒读取一次pod的内存和cpu使用信息并画曲线图。该功能需要k8s安装了metric-server才行。可以通过kubectl top pod来验证k8s是否有安装metric-server。

![pod资源查看](https://gitee.com/eeenet/xkube/raw/master/images/3.6.5-1.png)

#### 3.6.6.pod yaml创建
> 例子如下：
```yaml
apiVersion: v1
kind: Pod
metadata:
  namespace: zx-app #命名空间
  name: ops-tools   #pod名称
spec:
  containers:
  - name: ops-tools #容器名称
    image: xxx-registry.cn-shenzhen.cr.aliyuncs.com/public/ops-tools:202401012203   #镜像地址
    command: ["sleep","36000"] #执行命令
```


### 3.7.自定义资源[cdr]
> 自定义资源列表，主要针对一些非k8s自身资源，一些自主开发、或插件的云原生开发的资源信息打印出来。

![cdr列表](https://gitee.com/eeenet/xkube/raw/master/images/3.7.1-1.png)

### 3.8.自动伸缩[hpa]
> hpa是用来配置deployment的自动扩缩容的功能、根据cpu和内存的使用阀值自动扩缩pod的数量，具体创建可以参考3.1.9部分

![hpa列表](https://gitee.com/eeenet/xkube/raw/master/images/3.7.1-1.png)

### 3.9.yaml创建资源
> 该功能对应的命令：kubectl appy -f xxx.yaml的功能，主要是将yaml配置复制到输入框，选择对应的集群进行创建。

![yaml创建](https://gitee.com/eeenet/xkube/raw/master/images/3.9.1-1.png)

## 4.网络
### 4.1.服务[service]
#### 4.1.1.sevice列表
> 列表中的ip端口：是将service的IP和映射的端口组合，多个端口的逗号隔开了。内部端点：是selector 的app或k8s-app的值+targetport+protocol来显示的，有多个端口的逗号隔开。外部端点：当type: LoadBalancer 时才会显示的对外提供访问的IP和端口。外部端点这一列默认隐藏了

![svc列表](https://gitee.com/eeenet/xkube/raw/master/images/4.1.1-1.png)

#### 4.1.2.service详情
> service详情 除了将service的基本信息显示以外，还会将关联的deployment、pod、statefulset、关联的事件均显示在这里。参考下图

![svc详情](https://gitee.com/eeenet/xkube/raw/master/images/4.1.2-1.png)
![svc详情](https://gitee.com/eeenet/xkube/raw/master/images/4.1.2-2.png)

#### 4.1.3.界面创建service
> svc名称：service的名称，建议和关联的deployment或statefulset名称一致，关联:deployment或statefulset的名称，
svc类型：通常为ClusterIp,当有负载均衡时为LoadBlancer。NodePort时，访问入口是节点IP和端口。当创建无头服务时需要勾选“实例间服务发现”的复选框。
svc端口部分：
名称：分别是端口名称
服务端口：对外提供服务的端口
容器端口：deployment暴露出来的端口，以及选择协议
如果有多个端口时可以点击后边的+号。建议给service加上appname的标签，便于应用集将信息关联起来。例如：appname:test-app,当应用集是test-app时，可以关联所有标签是：appname:test-app的资源。

![svc创建](https://gitee.com/eeenet/xkube/raw/master/images/4.1.2-1.png)

#### 4.1.3.yaml创建service

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: testapp-readonly   #标签
  name: testapp-readonly    #service名称
  namespace: test-app       #命名空间
spec:
  ports:
  - name: http          #端口名称
    port: 8001          #对外提供服务的端口
    protocol: TCP       #协议
    targetPort: 8001    #deployment暴漏的端口
  selector:
    app: testapp-readonly   #关联的deployment的名称
  sessionAffinity: None
  type: ClusterIP           #服务类型
```

### 4.2.路由[ingress]
#### 4.2.1.ingress列表

![ingress列表](https://gitee.com/eeenet/xkube/raw/master/images/4.2.1-1.png)
> test1

#### 4.2.2.界面创ingress
> 通过界面创建ingress,目前该功能只支持配置单一域名，支持多个后端backend,支持配置证书，标签。当需要配置多个路径不同后端时，店家+号增加输入框。
选择正确的集群、命名空间以后，输入ingress的名称和 域名。
路径：默认是根路径，
匹配规则：一般选择Prefix前缀匹配。exact:准确匹配，ImplementationSpecific：匹配方法取决于 IngressClass。
服务：service的名称
端口：service的端口

![ingress创建](https://gitee.com/eeenet/xkube/raw/master/images/4.2.2-1.png)

#### 4.2.3.ingress详情
> 详情主要展示了ingress的基本信息，路由规则、以及事件信息，参考下图：

![ingress详情](https://gitee.com/eeenet/xkube/raw/master/images/4.2.3-1.png)

#### 4.2.4.yaml创建ingress
> yaml创建可以参考如下配置

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  labels:
    app: test-app-ing
  name: test-app-ing
spec:
  rules:
    - host: test.xxx.com.cn #域名
      http:
        paths:
          - backend:
              serviceName: test-app-nginx   #service的名称
              servicePort: 80
            path: /
    - host: test2.xxx.com.cn    #域名
      http:
        paths:
          - backend:
              serviceName: test1-app        #service的名称
              servicePort: 80
            path: /test1    #路径，不同路径走不同的后端时，可以添加backend
          - backend:
              serviceName: test-all-app     #service的名称
              servicePort: 80
            path: /             #路径，一般为根路径
  tls:
    - hosts:
        - test1.xxx.com.cn      #域名1
        - test2.xxx.com.cn      #域名2
      secretName: xxx-ssl-cert  #将证书和私钥存到secret,这是secret的名称,需要域名1和域名2都可以使用该证书

```


## 5.配置管理
### 5.1.配置项[configmap]
#### 5.1.1.configmap列表
>configmap列表参考下图

![configmap列表](https://gitee.com/eeenet/xkube/raw/master/images/5.1.1-1.png)

#### 5.1.2.configmap详情
> 详情页将configmap的基本信息展示出来，此处的键值信息不可修改，修改需要通过yaml编辑修改。

![configmap详情](https://gitee.com/eeenet/xkube/raw/master/images/5.1.2-1.png)

#### 5.1.3.界面创configmap
> 选择集群、命名空间、输入configmap名称，输入键值对即可，多个键值对，可以点击加号增加输入框。也支持多个标签添加。

![configmap创建](https://gitee.com/eeenet/xkube/raw/master/images/5.1.3-1.png)

#### 5.1.4.yaml创建configmap
>示例一 以单一键值对，key-value的形式存储

```yaml
apiVersion: v1
kind: ConfigMap
data:
  kubernetes-pods: mini
  sync_podMonitor: "false"
  sync_serviceMonitor: "false"
metadata:
  name: coredns
  namespace: kube-system
  resourceVersion: "7141240"
  uid: cec6dc70-a28c-4e36-83bb-8fecb2f0a7c4
```
> 示例二 以整体配置文件的形式

```yaml
apiVersion: v1
kind: ConfigMap
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready

        kubernetes cluster.local in-addr.arpa ip6.arpa {

          pods verified
          fallthrough in-addr.arpa ip6.arpa
          ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
metadata:
  name: coredns
  namespace: kube-system
  resourceVersion: "7141240"
  uid: cec6dc70-a28c-4e36-83bb-8fecb2f0a7c4
```

### 5.2.保密字典[secret]
#### 5.2.1.secret列表
>secret列表列表参考下图

![secret列表](https://gitee.com/eeenet/xkube/raw/master/images/5.2.1-1.png)

#### 5.2.2.secret详情
> 详情页将secret的基本信息展示出来，此处的键值信息不可修改，修改需要通过yaml编辑修改。

![secret详情](https://gitee.com/eeenet/xkube/raw/master/images/5.2.2-1.png)

#### 5.2.3.界面创secret
> 选择集群、命名空间、输入secret名称,
类型：Opaquue 通常为此类型，在此类型下可以输入多个键值对即可，也可以输入文件类型的，文件名或文件内容。镜像仓库登录密钥：改类型下输入的信息参考secret详情中的值，直接将认证信息粘贴进去即可。TLS证书，将crt或pem格式的证书文件粘贴到cert内容，私钥key文件粘贴到Key内。
标签可以通过点击加号添加多个。

![界面创secret](https://gitee.com/eeenet/xkube/raw/master/images/5.2.3-1.png)

#### 5.2.4.yaml创建secret
> Opaque类型的yaml参考

```yaml
apiVersion: v1
kind: Secret
data:
  ALICLOUD_AK: xxxTFRBdddlNck
  ALICLOUD_AS: xxxxR2N0SccclBWMG
  GIT_TOKEN: xxxxaHbbbbbZ
metadata:
  labels:
    app.kubernetes.io/instance: alertmanager2
  name: alertmanager2
  namespace: ops-app
type: Opaque

```

> 镜像仓库登录密钥类型的yaml参考

```yaml
apiVersion: v1
kind: Secret
data:
  .dockerconfigjson: ewoJImF1dGhzIjogewoJCSJ3h2WTJGc1FIQmpZWFYwYnpJd01qQTZaRlZaUUVwTldtNXBRalJ6U21oYSIKCQl9Cgl9Cn0K
metadata:
  labels:
    app: test-app
  name: pull-secrets
  namespace: test-app
type: kubernetes.io/dockerconfigjson

```

> TLS证书类型的yaml参考

```yaml
apiVersion: v1
kind: Secret
data:
  tls.crt: LS0tLS1CRUdJTiBDRycW5rY2tTemlZU1F0amlwSWNKREVIc1hvCjRIQT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==
  tls.key: LS0tLS1Wk92L1Jxc0tnSzNBQ2lzS1FML3JwSUptcXl5Vys0TWxCa2oxVUE9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=
metadata:
  labels:
    app: xxx-tls-secret
  name: xxx-tls-secret
  namespace: argocd
type: kubernetes.io/tls

```

## 6.存储
### 6.1.存储申明[Persistent Volume claim]
#### 6.1.1.pvc列表
> pvc列表功能参考下图:

![pvc列表](https://gitee.com/eeenet/xkube/raw/master/images/6.1.1-1.png)
#### 6.1.2.pvc详情
> pvc详情功能参考下图:

![pvc详情](https://gitee.com/eeenet/xkube/raw/master/images/6.1.2-1.png)
#### 6.1.3.pvc的yaml讲解

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: my-pvc      #pvc名称
  namespace: dev   # 命名空间
spec:
  accessModes:  #访问模式
    - ReadWriteMany #以read-write挂载到多个节点
  resources:
    requests:
      storage: 5Gi  # PVC允许申请的大小

```


### 6.2.存储卷[Persistent Volume]
#### 6.2.1.pv列表
> pv列表功能参考下图:
![pv列表](https://gitee.com/eeenet/xkube/raw/master/images/6.2.1-1.png)

#### 6.2.2.pv详情
> pv详情功能参考下图:

![pv详情](https://gitee.com/eeenet/xkube/raw/master/images/6.2.2-1.png)
#### 6.2.2.pv的yaml讲解

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nginx-pv    # pv的名称
spec:
  accessModes:      # 访问模式
    - ReadWriteMany # PV以read-write挂载到多个节点
  capacity:  # 容量
    storage: 2Gi    # pv可用的大小   
  nfs:
    path: /nfs/data/     # NFS的挂载路径
    server: 10.0.0.16    # NFS服务器地址
```


### 6.3.存储类[storageclass]
#### 6.3.1.存储类列表
> 存储类列表功能参考下图:

![存储类列表](https://gitee.com/eeenet/xkube/raw/master/images/6.3.1-1.png)

#### 6.3.2.存储类详情
> 存储类详情功能参考下图:

![存储类详情](https://gitee.com/eeenet/xkube/raw/master/images/6.3.2-1.png)
#### 6.3.2.存储类的yaml讲解


## 7.运维管理
### 7.1.事件中心
#### 7.1.1.事件查询
> 将当前集群的所有事件查询统一到此处，便于对于集群的事件的整体掌握。

![事件查询列表](https://gitee.com/eeenet/xkube/raw/master/images/7.1.1-1.png)
### 7.2.应用集
#### 7.2.1.应用集添加
> 应用集是指k8s中资源的集合，将同一个域名下的各种deployment/service/configmap/等资源关联一起便于信息查询。该关联是通过在资源中的标签:appname关联的，相同的appname关联到一个集合下。
参考如下配置：

![appname标签](https://gitee.com/eeenet/xkube/raw/master/images/7.2.1-1.png)

#### 7.2.2.应用维护wiki建立
> 可以将每个应用名创建一个独立的栏目，该栏目下存放该应用所有的维护文档。

### 7.3.CICD
#### 7.3.1.CICD列表
>CICD 目前引用阿里云的流水线，将流水线的功能关联到k8s集群，便于进行管理，在创建cicd之前需要在阿里云有开通并配置流水线。然后在阿里云帐号设置 这里添加上阿里云帐号的ak信息，并将流水线的组织ID填写上。

![CICD列表](https://gitee.com/eeenet/xkube/raw/master/images/7.3.1-1.png)

#### 7.3.2.阿里云帐号设置
> 添加界面各字段说明如下：
> - aliyunId:建议用阿里云的主账号名称，或者用一串能表示阿里云帐号意思的英文或数字的字符串。例如：aliyun001
> - accesskeyId: 登录阿里云后-点击右上角的头像--》权限与安全--》AccessKey管理--》创建AccessKey,就会生成AccessKey Id 和 AccessKey Secret了。
> - accesskeySecret: 同上

![添加阿里云帐号](https://gitee.com/eeenet/xkube/raw/master/images/7.3.2-1.png)
![添加阿里云帐号](https://gitee.com/eeenet/xkube/raw/master/images/7.3.2-2.png)

#### 7.3.3.添加阿里云流水线
> 各添加字段说明如下：
> - 名称：流水线名称,建议和deployment的名称一致,
> - 应用名：类似为流水线定义一个分组,通过该分组可以查询该组下所有的流水线。
> - 当前集群:该流水线当前发布的集群，如过名称和deployment一致，集群和命名空间也选择一致。当存在同一个流水线发布到多个k8s集群时，通过集群名称和命名空间区分。
> - 命名空间：当前流水线发布的命名空间，建议和deployment的命名空间一致
> - 阿里云帐号：选择该流水线所在的阿里云帐号。需先在运维-CICD-列表的左上方-阿里云帐号设置这里添加。
> - 组织ID：流水线里会存在一个企业，选择正确的企业下的流水线。
> -流水线ID：例如：https://flow.aliyun.com/pipelines/2679891/current中的2679891

![添加流水线](https://gitee.com/eeenet/xkube/raw/master/images/7.3.3-1.png)

#### 7.3.4.CICD详情

![流水线详情](https://gitee.com/eeenet/xkube/raw/master/images/7.3.4-1.png)

#### 7.3.5.怎样查看阿里云流水线组织ID
> 组织ID不需要手动获取，在添加阿里云帐号时，可以通过组织ID输入框旁边的按钮点击获取，一方面可以验证ak信息是否正确，另一方面能获取到该帐号下的所有组织名称和组织ID。

![添加阿里云帐号](https://gitee.com/eeenet/xkube/raw/master/images/7.3.5-1.png)

#### 7.3.6.怎样查看阿里云流水线ID
> 在流水线的运行界面和编辑界面时，可以复制浏览器中的URL，例如：https://flow.aliyun.com/pipelines/2908037/current ，其中的2908307就是流水线的ID。

#### 7.3.7.怎样在阿里云创建AK
> 参考下图

![阿里云创建AK](https://gitee.com/eeenet/xkube/raw/master/images/7.3.7-1.png)


## 8.权限管理
### 8.1.管理员
#### 8.1.1.管理员列表及修改
> 管理员的添加、删除、修改、角色授权、均在此页面内完成。
![管理员列表](https://gitee.com/eeenet/xkube/raw/master/images/8.1.1-1.png)
![管理员修改](https://gitee.com/eeenet/xkube/raw/master/images/8.1.1-2.png)

### 8.2.角色列表

![角色列表](https://gitee.com/eeenet/xkube/raw/master/images/8.2.1-1.png)

#### 8.2.1.角色添加功能授权
> 在角色列表点击授权列表、弹出框会弹出该角色赋予的授权的URL路径。点击添加会弹出目录结构框，然后勾选路径并选择授权给角色。取消授权：先复选框勾选路径，然后点击取消授权。默认已对不同角色设置好权限，如果不是太熟悉，建议不要在角色中删除授权。
参考下图：

![角色添加授权](https://gitee.com/eeenet/xkube/raw/master/images/8.2.1-2.png)
![角色添加授权2](https://gitee.com/eeenet/xkube/raw/master/images/8.2.1-3.png)

#### 8.2.2.角色赋予用户
> 在角色列表、点击用户列表在弹出框中操作，参考下图：

![用户赋予角色](https://gitee.com/eeenet/xkube/raw/master/images/8.2.2-1.png)

### 8.3.目录分组
#### 8.3.1.目录分组列表
> 目录分组的功能主要是将不同功能类型的URL进行分组，便于查看。

![目录分组列表](https://gitee.com/eeenet/xkube/raw/master/images/8.3.1-1.png)

### 8.4.目录结构
#### 8.4.1.目录结构的添加
> 该页面主要是针对有新开发的功能或路径时，需要添加到此处，并给予角色授权。自定义开时才用到。默认已配置好路径，建议不要删除里面的路径。

![目录列表](https://gitee.com/eeenet/xkube/raw/master/images/8.4.1-1.png)

#### 8.4.2.目录授权给角色
> 该功能和8.2.1的功能一致，只不过可以单独在此页面进行角色的功能授权。

### 8.5.集群授权
#### 8.5.1.集群授权列表
> 默认新建用户没有任何集群的权限，如果需要给用户授予某个集群的访问权限，在此处添加，如果选哟授予所有集群的权限，选择全部集群来授权。某个用户需要授权多个集群时，添加多条记录即可。

![集群授权](https://gitee.com/eeenet/xkube/raw/master/images/8.5.1-1.png)

## 9.文档中心
> 文档中心主要作用是用来针对服务、应用等编写维护文档。文档中心采用markdown的语法来进行记录。对于比较重要的信息可以采用加密的方式。

### 9.1.栏目列表
> 栏目是由添加或修改文章时定义，不支持单独添加，当两篇文章有相同的栏目名称时，同属一个栏目。一般栏目主要起到一个分组作用，就是将相同类型的文章划分到一个分组，便于查找。例如：应用集--wiki，一般认为一个应用则为一个栏目，然后该应用所有的维护文档可以放到这个应用名的栏目下。

![栏目列表](https://gitee.com/eeenet/xkube/raw/master/images/9.1.1-1.png)

### 9.2.文章列表
#### 9.2.1文章列表
![文章列表](https://gitee.com/eeenet/xkube/raw/master/images/9.2.1-1.png)

#### 9.2.2.文章添加
> 添加文章注意事项：
> 1. 栏目：如果需要将文章放到已存在的栏目下，输入的栏目名称必须和已有的名称一致。否则会产生一个新的栏目。
> 2. 标题：标题格式为不包含特殊字符的中英文即可。
> 3. 密码：如果需要给阅读该文章加个密码，可以在此处设置一个阅读密码。如果需要取消该密码，在编辑时将该输入框置空后提交。
> 4. 文章内容：编辑器采用markedown的格式进行编辑，具体的语法可以参考文章列表-顶部的markdown语法示例。或者参考9.2.5的URL。

![文章添加](https://gitee.com/eeenet/xkube/raw/master/images/9.2.2-1.png)

#### 9.2.3.文章修改
> 1. 一种是没有加密的文章，直接编辑修改,栏目、标题、文章内容都可以修改。如果需要加上米，直接在阅读密码输入框输入密码后点击提交。
> 2. 另一种是修改加密文章，修改前需要先输入密码，如果密码失败会刷新页面重新提示输入密码。输入正确密码以后，可以修改栏目、标题、文章，密码是以"******"的形式显示，表示有密码，如果需要修改密码，直接在输入框输入新的密码以后点击提交。

#### 9.2.4.文章加密
> 添加文章时在阅读密码处设置一个密码即可。由于默认是所有用户都能阅读文章，而某些文章中的内容不希望一部分用户看到，这时候可以加一个密码，只有自己或者部门内的人知道这个密码才能查看和修改。如果不希望用户查看文章，还可以通过给该用户设定特定的角色，并在角色权限赋予中不添加文章查看权限即可。
> #### 如果不知道密码怎么办？
> 答：可以在数据库中将该文章的authkey的这个字段置空，然后再编辑重新设置一个密码。
#### 9.2.5.markdown语法
> 参考URL：https://pandao.github.io/editor.md/