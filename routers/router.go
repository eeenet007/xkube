package routers

import (
	"xkube/controllers"
	admin "xkube/xadmin"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	admin.Run()

	//清理redis缓存 | Clear Redis cache
	beego.Router("/public/cache/Clear", &controllers.PublicIntfController{}, "*:Clear")

	//镜像更新接口 | Image update interface
	beego.Router("/public/cicd/UpdateImage", &controllers.PublicIntfController{}, "*:UpdateImage")

	//app网络检测
	beego.Router("/public/appnet/Check", &controllers.PublicIntfController{}, "*:Check")

	//定时重启deployment 和 sts
	beego.Router("/public/v1/ScheduleRestart", &controllers.PublicIntfController{}, "*:ScheduleRestart")

	//定时任务每天更新一下索引 | Scheduled task to update index daily
	beego.Router("/task/v1/UpdateIndex", &controllers.SearchController{}, "*:UpdateIndex") //更新索引 | Update index

	beego.Router("/search/v1/List", &controllers.SearchController{}, "*:List") //搜索查询 | Search query
	beego.Router("/search/v1/Del", &controllers.SearchController{}, "*:Del")   //删除无用的索引 | Delete unused indexes

	beego.Router("/xkube/backup/v1/List", &controllers.BackupController{}, "*:List")
	//beego.Router("/xkube/backup/v1/Del", &controllers.BackupController{}, "*:Del")
	beego.Router("/xkube/backup/v1/View", &controllers.BackupController{}, "*:View")
	beego.Router("/xkube/backup/v1/Backup", &controllers.BackupController{}, "*:Backup")
	beego.Router("/xkube/backup/v1/Recover", &controllers.BackupController{}, "*:Recover")

	//获取克隆任务列表 | Get clone task list
	beego.Router("/xkube/clone/v1/List", &controllers.CloneController{}, "*:List")

	//集群管理相关路由 | Cluster management related routes
	beego.Router("/xkube/cluster/v1/List", &controllers.ClusterController{}, "*:List")     //集群列表 | Cluster list
	beego.Router("/xkube/cluster/v1/Detail", &controllers.ClusterController{}, "*:Detail") //集群详情 | Cluster details
	beego.Router("/xkube/cluster/v1/Add", &controllers.ClusterController{}, "*:Add")       //添加集群 | Add cluster
	beego.Router("/xkube/cluster/v1/Edit", &controllers.ClusterController{}, "*:Edit")     //编辑集群 | Edit cluster
	beego.Router("/xkube/cluster/v1/Update", &controllers.ClusterController{}, "*:Update") //更新集群 | Update cluster
	beego.Router("/xkube/cluster/v1/Del", &controllers.ClusterController{}, "*:Del")       //删除集群 | Delete cluster
	beego.Router("/xkube/cluster/v1/Count", &controllers.ClusterController{}, "*:Count")   //集群统计 | Cluster count

	//应用名称管理相关路由 | App name management related routes
	beego.Router("/xkube/appname/v1/List", &controllers.AppnameController{}, "*:List")     //应用名称列表 | App name list
	beego.Router("/xkube/appname/v1/Add", &controllers.AppnameController{}, "*:Add")       //添加应用名称 | Add app name
	beego.Router("/xkube/appname/v1/Edit", &controllers.AppnameController{}, "*:Edit")     //编辑应用名称 | Edit app name
	beego.Router("/xkube/appname/v1/Update", &controllers.AppnameController{}, "*:Update") //更新应用名称 | Update app name
	beego.Router("/xkube/appname/v1/Del", &controllers.AppnameController{}, "*:Del")       //删除应用名称 | Delete app name

	//节点管理相关路由 | Node management related routes
	beego.Router("/xkube/node/v1/List", &controllers.NodeController{}, "*:List")                   //节点列表 | Node list
	beego.Router("/xkube/node/v1/PoolList", &controllers.NodeController{}, "*:PoolList")           //节点池列表 | Node pool list
	beego.Router("/xkube/node/v1/Detail", &controllers.NodeController{}, "*:Detail")               //节点详情 | Node details
	beego.Router("/xkube/node/v1/Yaml", &controllers.NodeController{}, "*:Yaml")                   //节点YAML | Node YAML
	beego.Router("/xkube/node/v1/Unschedulable", &controllers.NodeController{}, "*:Unschedulable") //节点可调度性设置 | Node schedulable setting
	beego.Router("/xkube/node/v1/Drain", &controllers.NodeController{}, "*:Drain")                 //节点排空 | Node drain
	beego.Router("/xkube/node/v1/Del", &controllers.NodeController{}, "*:Del")
	beego.Router("/xkube/node/v1/ModifyByYaml", &controllers.NodeController{}, "*:ModifyByYaml")
	beego.Router("/xkube/node/v1/GetLabels", &controllers.NodeController{}, "get:GetLabels")
	beego.Router("/xkube/node/v1/UpdateLabels", &controllers.NodeController{}, "post:UpdateLabels")
	beego.Router("/xkube/node/v1/GetTaint", &controllers.NodeController{}, "get:GetTaint")
	beego.Router("/xkube/node/v1/UpdateTaint", &controllers.NodeController{}, "post:UpdateTaint")
	beego.Router("/xkube/node/v1/GetAllocated", &controllers.NodeController{}, "get:GetAllocated")

	//命名空间管理相关路由 | Namespace management related routes
	beego.Router("/xkube/ns/v1/List", &controllers.NsController{}, "*:List")                 //命名空间列表 | Namespace list
	beego.Router("/xkube/ns/v1/Detail", &controllers.NsController{}, "*:Detail")             //命名空间详情 | Namespace details
	beego.Router("/xkube/ns/v1/Create", &controllers.NsController{}, "*:Create")             //创建命名空间 | Create namespace
	beego.Router("/xkube/ns/v1/Del", &controllers.NsController{}, "*:Del")                   //删除命名空间 | Delete namespace
	beego.Router("/xkube/ns/v1/ModifyByYaml", &controllers.NsController{}, "*:ModifyByYaml") //通过YAML修改命名空间 | Modify namespace by YAML
	beego.Router("/xkube/ns/v1/Yaml", &controllers.NsController{}, "*:Yaml")                 //命名空间YAML | Namespace YAML
	beego.Router("/xkube/ns/v1/LimitRange", &controllers.NsController{}, "*:LimitRange")     //命名空间资源限制 | Namespace resource limits

	//YAML应用相关路由 | YAML apply related routes
	beego.Router("/xkube/apply/v1/ApplyYaml", &controllers.ApplyYamlController{}, "*:ApplyYaml")       //应用YAML | Apply YAML
	beego.Router("/xkube/apply/v1/CreateByYaml", &controllers.ApplyYamlController{}, "*:CreateByYaml") //通过YAML创建 | Create by YAML

	//Deployment管理相关路由 | Deployment management related routes
	beego.Router("/xkube/deploy/v1/List", &controllers.DeployController{}, "*:List")                     //Deployment列表 | Deployment list
	beego.Router("/xkube/deploy/v1/Detail", &controllers.DeployController{}, "*:Detail")                 //Deployment详情 | Deployment details
	beego.Router("/xkube/deploy/v1/Create", &controllers.DeployController{}, "*:Create")                 //创建Deployment | Create Deployment
	beego.Router("/xkube/deploy/v1/Modify", &controllers.DeployController{}, "*:Modify")                 //修改Deployment | Modify Deployment
	beego.Router("/xkube/deploy/v1/ModifyByYaml", &controllers.DeployController{}, "*:ModifyByYaml")     //通过YAML修改Deployment | Modify Deployment by YAML
	beego.Router("/xkube/deploy/v1/Del", &controllers.DeployController{}, "*:Del")                       //删除Deployment | Delete Deployment
	beego.Router("/xkube/deploy/v1/Yaml", &controllers.DeployController{}, "*:Yaml")                     //Deployment YAML | Deployment YAML
	beego.Router("/xkube/deploy/v1/ReplicasetYaml", &controllers.DeployController{}, "*:ReplicasetYaml") //副本集YAML | ReplicaSet YAML
	beego.Router("/xkube/deploy/v1/ReplicasetList", &controllers.DeployController{}, "*:ReplicasetList") //副本集列表 | ReplicaSet list
	beego.Router("/xkube/deploy/v1/RollBack", &controllers.DeployController{}, "*:RollBack")             //回滚Deployment | Rollback Deployment
	beego.Router("/xkube/deploy/v1/Restart", &controllers.DeployController{}, "*:Restart")               //重启Deployment | Restart Deployment
	beego.Router("/xkube/deploy/v1/Labels", &controllers.DeployController{}, "*:Labels")                 //Deployment标签 | Deployment labels
	beego.Router("/xkube/deploy/v1/Image", &controllers.DeployController{}, "*:Image")                   //Deployment镜像 | Deployment image
	beego.Router("/xkube/deploy/v1/Replicas", &controllers.DeployController{}, "*:Replicas")             //Deployment副本数 | Deployment replicas
	beego.Router("/xkube/deploy/v1/Strategy", &controllers.DeployController{}, "*:Strategy")             //Deployment策略 | Deployment strategy
	beego.Router("/xkube/deploy/v1/Clone", &controllers.DeployController{}, "*:Clone")                   //克隆Deployment | Clone Deployment
	beego.Router("/xkube/deploy/v1/Check", &controllers.DeployController{}, "*:Check")                   //检查Deployment | Check Deployment
	beego.Router("/xkube/deploy/v1/Host", &controllers.DeployController{}, "get:Host")                   //
	beego.Router("/xkube/deploy/v1/Resource", &controllers.DeployController{}, "get:Resource")           //
	beego.Router("/xkube/deploy/v1/Probe", &controllers.DeployController{}, "get:Probe")                 //
	beego.Router("/xkube/deploy/v1/Env", &controllers.DeployController{}, "get:Env")
	beego.Router("/xkube/deploy/v1/Lifecycle", &controllers.DeployController{}, "get:Lifecycle")
	beego.Router("/xkube/deploy/v1/UpdateHost", &controllers.DeployController{}, "post:UpdateHost")         //
	beego.Router("/xkube/deploy/v1/UpdateResource", &controllers.DeployController{}, "post:UpdateResource") //
	beego.Router("/xkube/deploy/v1/UpdateProbe", &controllers.DeployController{}, "post:UpdateProbe")       //
	beego.Router("/xkube/deploy/v1/UpdateEnv", &controllers.DeployController{}, "post:UpdateEnv")
	beego.Router("/xkube/deploy/v1/UpdateLifecycle", &controllers.DeployController{}, "post:UpdateLifecycle")
	beego.Router("/xkube/deploy/v1/GetTolerations", &controllers.DeployController{}, "get:GetTolerations")
	beego.Router("/xkube/deploy/v1/UpdateTolerations", &controllers.DeployController{}, "post:UpdateTolerations")
	beego.Router("/xkube/deploy/v1/GetNodeAffinity", &controllers.DeployController{}, "get:GetNodeAffinity")
	beego.Router("/xkube/deploy/v1/UpdateNodeAffinity", &controllers.DeployController{}, "post:UpdateNodeAffinity")
	beego.Router("/xkube/deploy/v1/GetPodAffinity", &controllers.DeployController{}, "get:GetPodAffinity")
	beego.Router("/xkube/deploy/v1/UpdatePodAffinity", &controllers.DeployController{}, "post:UpdatePodAffinity")

	//StatefulSet管理相关路由 | StatefulSet management related routes
	beego.Router("/xkube/sts/v1/List", &controllers.StatefulsetController{}, "*:List")     //StatefulSet列表 | StatefulSet list
	beego.Router("/xkube/sts/v1/Detail", &controllers.StatefulsetController{}, "*:Detail") //StatefulSet详情 | StatefulSet details
	beego.Router("/xkube/sts/v1/Create", &controllers.StatefulsetController{}, "*:Create") //创建StatefulSet | Create StatefulSet
	beego.Router("/xkube/sts/v1/Modify", &controllers.StatefulsetController{}, "*:Modify") //修改StatefulSet | Modify StatefulSet
	beego.Router("/xkube/sts/v1/Del", &controllers.StatefulsetController{}, "*:Del")       //删除StatefulSet | Delete StatefulSet
	beego.Router("/xkube/sts/v1/Yaml", &controllers.StatefulsetController{}, "*:Yaml")     //StatefulSet YAML | StatefulSet YAML
	beego.Router("/xkube/sts/v1/ModifyByYaml", &controllers.StatefulsetController{}, "*:ModifyByYaml")
	beego.Router("/xkube/sts/v1/RollBack", &controllers.StatefulsetController{}, "*:RollBack") //回滚StatefulSet | Rollback StatefulSet
	beego.Router("/xkube/sts/v1/Restart", &controllers.StatefulsetController{}, "*:Restart")   //重启StatefulSet | Restart StatefulSet
	beego.Router("/xkube/sts/v1/Clone", &controllers.StatefulsetController{}, "*:Clone")       //克隆StatefulSet | Clone StatefulSet
	beego.Router("/xkube/sts/v1/Labels", &controllers.StatefulsetController{}, "*:Labels")     //StatefulSet标签 | StatefulSet labels
	beego.Router("/xkube/sts/v1/Image", &controllers.StatefulsetController{}, "*:Image")       //StatefulSet镜像 | StatefulSet image
	beego.Router("/xkube/sts/v1/Replicas", &controllers.StatefulsetController{}, "*:Replicas") //StatefulSet副本数 | StatefulSet replicas
	beego.Router("/xkube/sts/v1/GetHost", &controllers.StatefulsetController{}, "get:GetHost")
	beego.Router("/xkube/sts/v1/UpdateHost", &controllers.StatefulsetController{}, "post:UpdateHost")
	beego.Router("/xkube/sts/v1/GetEnv", &controllers.StatefulsetController{}, "get:GetEnv")
	beego.Router("/xkube/sts/v1/UpdateEnv", &controllers.StatefulsetController{}, "post:UpdateEnv")
	beego.Router("/xkube/sts/v1/GetProbe", &controllers.StatefulsetController{}, "get:GetProbe")
	beego.Router("/xkube/sts/v1/UpdateProbe", &controllers.StatefulsetController{}, "post:UpdateProbe")
	beego.Router("/xkube/sts/v1/GetResource", &controllers.StatefulsetController{}, "get:GetResource")
	beego.Router("/xkube/sts/v1/UpdateResource", &controllers.StatefulsetController{}, "post:UpdateResource")
	beego.Router("/xkube/sts/v1/GetLifecycle", &controllers.StatefulsetController{}, "get:GetLifecycle")
	beego.Router("/xkube/sts/v1/UpdateLifecycle", &controllers.StatefulsetController{}, "post:UpdateLifecycle")
	beego.Router("/xkube/sts/v1/UpdateNodeAffinity", &controllers.StatefulsetController{}, "post:UpdateNodeAffinity")
	beego.Router("/xkube/sts/v1/GetNodeAffinity", &controllers.StatefulsetController{}, "get:GetNodeAffinity")

	//DaemonSet管理相关路由 | DaemonSet management related routes
	beego.Router("/xkube/ds/v1/List", &controllers.DaemonsetController{}, "*:List")                 //DaemonSet列表 | DaemonSet list
	beego.Router("/xkube/ds/v1/Detail", &controllers.DaemonsetController{}, "*:Detail")             //DaemonSet详情 | DaemonSet details
	beego.Router("/xkube/ds/v1/Create", &controllers.DaemonsetController{}, "*:Create")             //创建DaemonSet | Create DaemonSet
	beego.Router("/xkube/ds/v1/Modify", &controllers.DaemonsetController{}, "*:Modify")             //修改DaemonSet | Modify DaemonSet
	beego.Router("/xkube/ds/v1/Del", &controllers.DaemonsetController{}, "*:Del")                   //删除DaemonSet | Delete DaemonSet
	beego.Router("/xkube/ds/v1/Yaml", &controllers.DaemonsetController{}, "*:Yaml")                 //DaemonSet YAML | DaemonSet YAML
	beego.Router("/xkube/ds/v1/ModifyByYaml", &controllers.DaemonsetController{}, "*:ModifyByYaml") //通过YAML修改DaemonSet | Modify DaemonSet by YAML

	//Job管理相关路由 | Job management related routes
	beego.Router("/xkube/job/v1/List", &controllers.JobController{}, "*:List")     //Job列表 | Job list
	beego.Router("/xkube/job/v1/Detail", &controllers.JobController{}, "*:Detail") //Job详情 | Job details
	beego.Router("/xkube/job/v1/Create", &controllers.JobController{}, "*:Create") //创建Job | Create Job
	beego.Router("/xkube/job/v1/Modify", &controllers.JobController{}, "*:Modify") //修改Job | Modify Job
	beego.Router("/xkube/job/v1/Del", &controllers.JobController{}, "*:Del")       //删除Job | Delete Job
	beego.Router("/xkube/job/v1/Yaml", &controllers.JobController{}, "*:Yaml")     //Job YAML | Job YAML
	beego.Router("/xkube/job/v1/Log", &controllers.JobController{}, "*:Log")       //Job日志 | Job logs
	beego.Router("/xkube/job/v1/ModifyByYaml", &controllers.JobController{}, "*:ModifyByYaml")
	//CronJob管理相关路由 | CronJob management related routes
	beego.Router("/xkube/cronjob/v1/List", &controllers.CronjobController{}, "*:List")                 //CronJob列表 | CronJob list
	beego.Router("/xkube/cronjob/v1/Detail", &controllers.CronjobController{}, "*:Detail")             //CronJob详情 | CronJob details
	beego.Router("/xkube/cronjob/v1/Create", &controllers.CronjobController{}, "*:Create")             //创建CronJob | Create CronJob
	beego.Router("/xkube/cronjob/v1/Modify", &controllers.CronjobController{}, "*:Modify")             //修改CronJob | Modify CronJob
	beego.Router("/xkube/cronjob/v1/ModifyByYaml", &controllers.CronjobController{}, "*:ModifyByYaml") //通过YAML修改CronJob | Modify CronJob by YAML
	beego.Router("/xkube/cronjob/v1/Del", &controllers.CronjobController{}, "*:Del")                   //删除CronJob | Delete CronJob
	beego.Router("/xkube/cronjob/v1/Yaml", &controllers.CronjobController{}, "*:Yaml")                 //CronJob YAML | CronJob YAML
	beego.Router("/xkube/cronjob/v1/Labels", &controllers.CronjobController{}, "*:Labels")             //CronJob标签 | CronJob labels
	beego.Router("/xkube/cronjob/v1/Clone", &controllers.CronjobController{}, "*:Clone")               //克隆CronJob | Clone CronJob
	beego.Router("/xkube/cronjob/v1/Run", &controllers.CronjobController{}, "*:Run")                   //运行CronJob | Run CronJob

	//CronJob beta1版本管理相关路由 | CronJob beta1 version management related routes
	beego.Router("/xkube/cronjob/beta1/List", &controllers.CronjobBeta1Controller{}, "*:List")                 //CronJob beta1列表 | CronJob beta1 list
	beego.Router("/xkube/cronjob/beta1/Detail", &controllers.CronjobBeta1Controller{}, "*:Detail")             //CronJob beta1详情 | CronJob beta1 details
	beego.Router("/xkube/cronjob/beta1/Create", &controllers.CronjobBeta1Controller{}, "*:Create")             //创建CronJob beta1 | Create CronJob beta1
	beego.Router("/xkube/cronjob/beta1/Modify", &controllers.CronjobBeta1Controller{}, "*:Modify")             //修改CronJob beta1 | Modify CronJob beta1
	beego.Router("/xkube/cronjob/beta1/ModifyByYaml", &controllers.CronjobBeta1Controller{}, "*:ModifyByYaml") //通过YAML修改CronJob beta1 | Modify CronJob beta1 by YAML
	beego.Router("/xkube/cronjob/beta1/Del", &controllers.CronjobBeta1Controller{}, "*:Del")                   //删除CronJob beta1 | Delete CronJob beta1
	beego.Router("/xkube/cronjob/beta1/Yaml", &controllers.CronjobBeta1Controller{}, "*:Yaml")                 //CronJob beta1 YAML | CronJob beta1 YAML
	beego.Router("/xkube/cronjob/beta1/Labels", &controllers.CronjobBeta1Controller{}, "*:Labels")             //CronJob beta1标签 | CronJob beta1 labels
	beego.Router("/xkube/cronjob/beta1/Clone", &controllers.CronjobBeta1Controller{}, "*:Clone")               //克隆CronJob beta1 | Clone CronJob beta1
	beego.Router("/xkube/cronjob/beta1/Run", &controllers.CronjobBeta1Controller{}, "*:Run")                   //运行CronJob beta1 | Run CronJob beta1

	//Pod管理相关路由 | Pod management related routes
	beego.Router("/xkube/pod/v1/List", &controllers.PodController{}, "*:List")                   //Pod列表 | Pod list
	beego.Router("/xkube/pod/v1/ContainerList", &controllers.PodController{}, "*:ContainerList") //容器列表 | Container list
	beego.Router("/xkube/pod/v1/Detail", &controllers.PodController{}, "*:Detail")               //Pod详情 | Pod details
	beego.Router("/xkube/pod/v1/Log", &controllers.PodController{}, "*:Log")                     //Pod日志 | Pod logs
	beego.Router("/xkube/pod/v1/Del", &controllers.PodController{}, "*:Del")                     //删除Pod | Delete Pod
	beego.Router("/xkube/pod/v1/Yaml", &controllers.PodController{}, "*:Yaml")                   //Pod YAML | Pod YAML
	beego.Router("/xkube/pod/v1/Exec", &controllers.PodController{}, "*:Exec")                   //Pod执行 | Pod exec
	beego.Router("/xkube/pod/v1/Check", &controllers.PodController{}, "*:Check")
	beego.Router("/xkube/pod/v1/ModifyByYaml", &controllers.PodController{}, "*:ModifyByYaml")

	//CRD管理相关路由 | CRD management related routes
	beego.Router("/xkube/cdr/v1/List", &controllers.CdrController{}, "*:List") //CRD列表 | CRD list
	beego.Router("/xkube/cdr/v1/Yaml", &controllers.CdrController{}, "*:Yaml") //CRD YAML | CRD YAML
	beego.Router("/xkube/cdr/v1/Del", &controllers.CdrController{}, "*:Del")   //删除CRD | Delete CRD

	//Service管理相关路由 | Service management related routes
	beego.Router("/xkube/svc/v1/List", &controllers.SvcController{}, "*:List")                 //Service列表 | Service list
	beego.Router("/xkube/svc/v1/Detail", &controllers.SvcController{}, "*:Detail")             //Service详情 | Service details
	beego.Router("/xkube/svc/v1/Create", &controllers.SvcController{}, "*:Create")             //创建Service | Create Service
	beego.Router("/xkube/svc/v1/ModifyByYaml", &controllers.SvcController{}, "*:ModifyByYaml") //通过YAML修改Service | Modify Service by YAML
	beego.Router("/xkube/svc/v1/Del", &controllers.SvcController{}, "*:Del")                   //删除Service | Delete Service
	beego.Router("/xkube/svc/v1/Yaml", &controllers.SvcController{}, "*:Yaml")                 //Service YAML | Service YAML
	beego.Router("/xkube/svc/v1/Clone", &controllers.SvcController{}, "*:Clone")               //克隆Service | Clone Service

	//Ingress管理相关路由 | Ingress management related routes
	beego.Router("/xkube/ing/v1/List", &controllers.IngressController{}, "*:List")                 //Ingress列表 | Ingress list
	beego.Router("/xkube/ing/v1/Detail", &controllers.IngressController{}, "*:Detail")             //Ingress详情 | Ingress details
	beego.Router("/xkube/ing/v1/Create", &controllers.IngressController{}, "*:Create")             //创建Ingress | Create Ingress
	beego.Router("/xkube/ing/v1/ModifyByYaml", &controllers.IngressController{}, "*:ModifyByYaml") //通过YAML修改Ingress | Modify Ingress by YAML
	beego.Router("/xkube/ing/v1/Del", &controllers.IngressController{}, "*:Del")                   //删除Ingress | Delete Ingress
	beego.Router("/xkube/ing/v1/Yaml", &controllers.IngressController{}, "*:Yaml")                 //Ingress YAML | Ingress YAML
	beego.Router("/xkube/ing/v1/GetRule", &controllers.IngressController{}, "get:GetRule")
	beego.Router("/xkube/ing/v1/UpdateRule", &controllers.IngressController{}, "post:UpdateRule")
	beego.Router("/xkube/ing/v1/GetTlsHost", &controllers.IngressController{}, "get:GetTlsHost")
	beego.Router("/xkube/ing/v1/UpdateTlsHost", &controllers.IngressController{}, "post:UpdateTlsHost")

	//Ingress v1beta1版本管理相关路由 | Ingress v1beta1 version management related routes
	beego.Router("/xkube/ing/v1beta1/List", &controllers.IngressV1beta1Controller{}, "*:List")                 //Ingress v1beta1列表 | Ingress v1beta1 list
	beego.Router("/xkube/ing/v1beta1/Detail", &controllers.IngressV1beta1Controller{}, "*:Detail")             //Ingress v1beta1详情 | Ingress v1beta1 details
	beego.Router("/xkube/ing/v1beta1/Create", &controllers.IngressV1beta1Controller{}, "*:Create")             //创建Ingress v1beta1 | Create Ingress v1beta1
	beego.Router("/xkube/ing/v1beta1/ModifyByYaml", &controllers.IngressV1beta1Controller{}, "*:ModifyByYaml") //通过YAML修改Ingress v1beta1 | Modify Ingress v1beta1 by YAML
	beego.Router("/xkube/ing/v1beta1/Del", &controllers.IngressV1beta1Controller{}, "*:Del")                   //删除Ingress v1beta1 | Delete Ingress v1beta1
	beego.Router("/xkube/ing/v1beta1/Yaml", &controllers.IngressV1beta1Controller{}, "*:Yaml")                 //Ingress v1beta1 YAML | Ingress v1beta1 YAML
	beego.Router("/xkube/ing/v1beta1/GetRule", &controllers.IngressV1beta1Controller{}, "get:GetRule")
	beego.Router("/xkube/ing/v1beta1/UpdateRule", &controllers.IngressV1beta1Controller{}, "post:UpdateRule")
	beego.Router("/xkube/ing/v1beta1/GetTlsHost", &controllers.IngressV1beta1Controller{}, "get:GetTlsHost")
	beego.Router("/xkube/ing/v1beta1/UpdateTlsHost", &controllers.IngressV1beta1Controller{}, "post:UpdateTlsHost")

	//GatewayClass相关路由 | GatewayClass related routes
	beego.Router("/xkube/gatewayclass/v1/List", &controllers.GatewayClassController{}, "*:List")     //GatewayClass列表 | GatewayClass list
	beego.Router("/xkube/gatewayclass/v1/Detail", &controllers.GatewayClassController{}, "*:Detail") //GatewayClass详情 | GatewayClass detail
	beego.Router("/xkube/gatewayclass/v1/Yaml", &controllers.GatewayClassController{}, "*:Yaml")     //GatewayClass YAML配置 | GatewayClass YAML config

	//Gateway相关路由 | Gateway related routes
	beego.Router("/xkube/gateway/v1/List", &controllers.GatewayController{}, "*:List")                    //Gateway列表 | Gateway list
	beego.Router("/xkube/gateway/v1/Detail", &controllers.GatewayController{}, "*:Detail")                //Gateway详情 | Gateway detail
	beego.Router("/xkube/gateway/v1/Yaml", &controllers.GatewayController{}, "*:Yaml")                    //Gateway YAML配置 | Gateway YAML config
	beego.Router("/xkube/gateway/v1/Delete", &controllers.GatewayController{}, "*:Delete")                //删除Gateway | Delete Gateway
	beego.Router("/xkube/gateway/v1/Create", &controllers.GatewayController{}, "post:Create")             //创建Gateway | Create Gateway
	beego.Router("/xkube/gateway/v1/UpdateByYaml", &controllers.GatewayController{}, "post:UpdateByYaml") //通过YAML更新Gateway | Update Gateway by YAML

	//HTTPRoute相关路由 | HTTPRoute related routes
	beego.Router("/xkube/httproute/v1/List", &controllers.HTTPRouteController{}, "*:List")                    //HTTPRoute列表 | HTTPRoute list
	beego.Router("/xkube/httproute/v1/Detail", &controllers.HTTPRouteController{}, "*:Detail")                //HTTPRoute详情 | HTTPRoute detail
	beego.Router("/xkube/httproute/v1/Yaml", &controllers.HTTPRouteController{}, "*:Yaml")                    //HTTPRoute YAML配置 | HTTPRoute YAML config
	beego.Router("/xkube/httproute/v1/Delete", &controllers.HTTPRouteController{}, "*:Delete")                //删除HTTPRoute | Delete HTTPRoute
	beego.Router("/xkube/httproute/v1/Create", &controllers.HTTPRouteController{}, "post:Create")             //创建HTTPRoute | Create HTTPRoute
	beego.Router("/xkube/httproute/v1/UpdateByYaml", &controllers.HTTPRouteController{}, "post:UpdateByYaml") //通过YAML更新HTTPRoute | Update HTTPRoute by YAML

	//GRPCRoute相关路由 | GRPCRoute related routes
	beego.Router("/xkube/grpcroute/v1/List", &controllers.GRPCRouteController{}, "*:List")                    //GRPCRoute列表 | GRPCRoute list
	beego.Router("/xkube/grpcroute/v1/Detail", &controllers.GRPCRouteController{}, "*:Detail")                //GRPCRoute详情 | GRPCRoute detail
	beego.Router("/xkube/grpcroute/v1/Yaml", &controllers.GRPCRouteController{}, "*:Yaml")                    //GRPCRoute YAML配置 | GRPCRoute YAML config
	beego.Router("/xkube/grpcroute/v1/Delete", &controllers.GRPCRouteController{}, "*:Delete")                //删除GRPCRoute | Delete GRPCRoute
	beego.Router("/xkube/grpcroute/v1/Create", &controllers.GRPCRouteController{}, "post:Create")             //创建GRPCRoute | Create GRPCRoute
	beego.Router("/xkube/grpcroute/v1/UpdateByYaml", &controllers.GRPCRouteController{}, "post:UpdateByYaml") //通过YAML更新GRPCRoute | Update GRPCRoute by YAML

	//TCPRoute相关路由 | TCPRoute related routes
	beego.Router("/xkube/tcproute/v1/List", &controllers.TCPRouteController{}, "*:List")                    //TCPRoute列表 | TCPRoute list
	beego.Router("/xkube/tcproute/v1/Detail", &controllers.TCPRouteController{}, "*:Detail")                //TCPRoute详情 | TCPRoute detail
	beego.Router("/xkube/tcproute/v1/Yaml", &controllers.TCPRouteController{}, "*:Yaml")                    //TCPRoute YAML配置 | TCPRoute YAML config
	beego.Router("/xkube/tcproute/v1/Delete", &controllers.TCPRouteController{}, "*:Delete")                //删除TCPRoute | Delete TCPRoute
	beego.Router("/xkube/tcproute/v1/Create", &controllers.TCPRouteController{}, "post:Create")             //创建TCPRoute | Create TCPRoute
	beego.Router("/xkube/tcproute/v1/UpdateByYaml", &controllers.TCPRouteController{}, "post:UpdateByYaml") //通过YAML更新TCPRoute | Update TCPRoute by YAML

	//UDPRoute相关路由 | UDPRoute related routes
	beego.Router("/xkube/udproute/v1/List", &controllers.UDPRouteController{}, "*:List")                    //UDPRoute列表 | UDPRoute list
	beego.Router("/xkube/udproute/v1/Detail", &controllers.UDPRouteController{}, "*:Detail")                //UDPRoute详情 | UDPRoute detail
	beego.Router("/xkube/udproute/v1/Yaml", &controllers.UDPRouteController{}, "*:Yaml")                    //UDPRoute YAML配置 | UDPRoute YAML config
	beego.Router("/xkube/udproute/v1/Delete", &controllers.UDPRouteController{}, "*:Delete")                //删除UDPRoute | Delete UDPRoute
	beego.Router("/xkube/udproute/v1/Create", &controllers.UDPRouteController{}, "post:Create")             //创建UDPRoute | Create UDPRoute
	beego.Router("/xkube/udproute/v1/UpdateByYaml", &controllers.UDPRouteController{}, "post:UpdateByYaml") //通过YAML更新UDPRoute | Update UDPRoute by YAML

	//ConfigMap管理相关路由 | ConfigMap management related routes
	beego.Router("/xkube/cm/v1/List", &controllers.ConfigMapController{}, "*:List")                 //ConfigMap列表 | ConfigMap list
	beego.Router("/xkube/cm/v1/Detail", &controllers.ConfigMapController{}, "*:Detail")             //ConfigMap详情 | ConfigMap details
	beego.Router("/xkube/cm/v1/Create", &controllers.ConfigMapController{}, "*:Create")             //创建ConfigMap | Create ConfigMap
	beego.Router("/xkube/cm/v1/ModifyByYaml", &controllers.ConfigMapController{}, "*:ModifyByYaml") //通过YAML修改ConfigMap | Modify ConfigMap by YAML
	beego.Router("/xkube/cm/v1/Del", &controllers.ConfigMapController{}, "*:Del")                   //删除ConfigMap | Delete ConfigMap
	beego.Router("/xkube/cm/v1/Yaml", &controllers.ConfigMapController{}, "*:Yaml")                 //ConfigMap YAML | ConfigMap YAML
	beego.Router("/xkube/cm/v1/Clone", &controllers.ConfigMapController{}, "*:Clone")               //克隆ConfigMap | Clone ConfigMap

	//Secret管理相关路由 | Secret management related routes
	beego.Router("/xkube/secret/v1/List", &controllers.SecretController{}, "*:List")                 //Secret列表 | Secret list
	beego.Router("/xkube/secret/v1/Detail", &controllers.SecretController{}, "*:Detail")             //Secret详情 | Secret details
	beego.Router("/xkube/secret/v1/Create", &controllers.SecretController{}, "*:Create")             //创建Secret | Create Secret
	beego.Router("/xkube/secret/v1/ModifyByYaml", &controllers.SecretController{}, "*:ModifyByYaml") //通过YAML修改Secret | Modify Secret by YAML
	beego.Router("/xkube/secret/v1/Del", &controllers.SecretController{}, "*:Del")                   //删除Secret | Delete Secret
	beego.Router("/xkube/secret/v1/Yaml", &controllers.SecretController{}, "*:Yaml")                 //Secret YAML | Secret YAML
	beego.Router("/xkube/secret/v1/Clone", &controllers.SecretController{}, "*:Clone")               //克隆Secret | Clone Secret

	//StorageClass管理相关路由 | StorageClass management related routes
	beego.Router("/xkube/storageclass/v1/List", &controllers.StorageClassController{}, "*:List")     //StorageClass列表 | StorageClass list
	beego.Router("/xkube/storageclass/v1/Yaml", &controllers.StorageClassController{}, "*:Yaml")     //StorageClass YAML | StorageClass YAML
	beego.Router("/xkube/storageclass/v1/Detail", &controllers.StorageClassController{}, "*:Detail") //StorageClass详情 | StorageClass details

	//PV管理相关路由 | PV management related routes
	beego.Router("/xkube/pv/v1/List", &controllers.PvController{}, "*:List")     //PV列表 | PV list
	beego.Router("/xkube/pv/v1/Yaml", &controllers.PvController{}, "*:Yaml")     //PV YAML | PV YAML
	beego.Router("/xkube/pv/v1/Detail", &controllers.PvController{}, "*:Detail") //PV详情 | PV details

	//PVC管理相关路由 | PVC management related routes
	beego.Router("/xkube/pvc/v1/List", &controllers.PvcController{}, "*:List")     //PVC列表 | PVC list
	beego.Router("/xkube/pvc/v1/Yaml", &controllers.PvcController{}, "*:Yaml")     //PVC YAML | PVC YAML
	beego.Router("/xkube/pvc/v1/Detail", &controllers.PvcController{}, "*:Detail") //PVC详情 | PVC details

	//HPA管理相关路由 | HPA management related routes
	beego.Router("/xkube/hpa/v1/List", &controllers.HpaController{}, "*:List")                 //HPA列表 | HPA list
	beego.Router("/xkube/hpa/v1/Yaml", &controllers.HpaController{}, "*:Yaml")                 //HPA YAML | HPA YAML
	beego.Router("/xkube/hpa/v1/ModifyByYaml", &controllers.HpaController{}, "*:ModifyByYaml") //通过YAML修改HPA | Modify HPA by YAML
	beego.Router("/xkube/hpa/v1/Create", &controllers.HpaController{}, "*:Create")             //创建HPA | Create HPA
	beego.Router("/xkube/hpa/v1/Del", &controllers.HpaController{}, "*:Del")                   //删除HPA | Delete HPA

	//HPA v2beta2版本管理相关路由 | HPA v2beta2 version management related routes
	beego.Router("/xkube/hpa/v2beta2/List", &controllers.HpaV2beta2Controller{}, "*:List")                 //HPA v2beta2列表 | HPA v2beta2 list
	beego.Router("/xkube/hpa/v2beta2/Yaml", &controllers.HpaV2beta2Controller{}, "*:Yaml")                 //HPA v2beta2 YAML | HPA v2beta2 YAML
	beego.Router("/xkube/hpa/v2beta2/ModifyByYaml", &controllers.HpaV2beta2Controller{}, "*:ModifyByYaml") //通过YAML修改HPA v2beta2 | Modify HPA v2beta2 by YAML
	beego.Router("/xkube/hpa/v2beta2/Create", &controllers.HpaV2beta2Controller{}, "*:Create")             //创建HPA v2beta2 | Create HPA v2beta2
	beego.Router("/xkube/hpa/v2beta2/Del", &controllers.HpaV2beta2Controller{}, "*:Del")                   //删除HPA v2beta2 | Delete HPA v2beta2

	//事件管理相关路由 | Event management related routes
	beego.Router("/xkube/event/v1/List", &controllers.EventController{}, "*:List") //事件列表 | Event list

	//WebSocket终端相关路由 | WebSocket terminal related routes
	beego.Handler("/public/pod/terminal/ws", &controllers.TerminalSockjs{}, true)

	//集群角色绑定管理相关路由 | Cluster role binding management related routes
	beego.Router("/xkube/clusterrolebinding/v1/List", &controllers.ClusterRoleBindingController{}, "*:List") //集群角色绑定列表 | Cluster role binding list
	beego.Router("/xkube/clusterrolebinding/v1/Yaml", &controllers.ClusterRoleBindingController{}, "*:Yaml") //集群角色绑定YAML | Cluster role binding YAML

	//集群角色管理相关路由 | Cluster roles management related routes
	beego.Router("/xkube/clusterroles/v1/List", &controllers.ClusterRolesController{}, "*:List") //集群角色列表 | Cluster roles list
	beego.Router("/xkube/clusterroles/v1/Yaml", &controllers.ClusterRolesController{}, "*:Yaml") //集群角色YAML | Cluster roles YAML

	//服务账户管理相关路由 | Service accounts management related routes
	beego.Router("/xkube/serviceaccounts/v1/List", &controllers.ServiceAccountsController{}, "*:List") //服务账户列表 | Service accounts list
	beego.Router("/xkube/serviceaccounts/v1/Yaml", &controllers.ServiceAccountsController{}, "*:Yaml") //服务账户YAML | Service accounts YAML

	//角色管理相关路由 | Roles management related routes
	beego.Router("/xkube/roles/v1/List", &controllers.RolesController{}, "*:List") //角色列表 | Roles list
	beego.Router("/xkube/roles/v1/Yaml", &controllers.RolesController{}, "*:Yaml") //角色YAML | Roles YAML

	//角色绑定管理相关路由 | Role binding management related routes
	beego.Router("/xkube/rolebinding/v1/List", &controllers.RoleBindingController{}, "*:List") //角色绑定列表 | Role binding list
	beego.Router("/xkube/rolebinding/v1/Yaml", &controllers.RoleBindingController{}, "*:Yaml") //角色绑定YAML | Role binding YAML

	//监控指标相关路由 | Metrics related routes
	beego.Router("/xkube/metrics/v1/PodList", &controllers.MetricsController{}, "*:PodList")   //Pod指标列表 | Pod metrics list
	beego.Router("/xkube/metrics/PodUsage", &controllers.MetricsController{}, "*:PodUsage")    //Pod使用情况 | Pod usage
	beego.Router("/xkube/metrics/NodeUsage", &controllers.MetricsController{}, "*:NodeUsage")  //节点使用情况 | Node usage
	beego.Router("/xkube/metrics/v1/NodeList", &controllers.MetricsController{}, "*:NodeList") //节点指标列表 | Node metrics list

	//CICD相关路由 | CICD related routes
	beego.Router("/cicd/v1/List", &controllers.CicdController{}, "*:List")     //CICD列表 | CICD list
	beego.Router("/cicd/v1/Add", &controllers.CicdController{}, "*:Add")       //添加CICD | Add CICD
	beego.Router("/cicd/v1/Update", &controllers.CicdController{}, "*:Update") //更新CICD | Update CICD
	beego.Router("/cicd/v1/Del", &controllers.CicdController{}, "*:Del")       //删除CICD | Delete CICD
	beego.Router("/cicd/v1/Edit", &controllers.CicdController{}, "*:Edit")     //编辑CICD | Edit CICD

	beego.Router("/cicd/v1/AkAdd", &controllers.CicdController{}, "*:AkAdd")   //添加访问密钥 | Add access key
	beego.Router("/cicd/v1/AkList", &controllers.CicdController{}, "*:AkList") //访问密钥列表 | Access key list
	beego.Router("/cicd/v1/AkDel", &controllers.CicdController{}, "*:AkDel")   //删除访问密钥 | Delete access key

	beego.Router("/cicd/v1/JksAdd", &controllers.CicdController{}, "*:JksAdd")
	beego.Router("/cicd/v1/JksList", &controllers.CicdController{}, "*:JksList")
	beego.Router("/cicd/v1/GetJksList", &controllers.CicdController{}, "*:JksList")
	beego.Router("/cicd/v1/JksDel", &controllers.CicdController{}, "*:JksDel")

	beego.Router("/cicd/jks/v1/Run", &controllers.CicdController{}, "*:Run")
	beego.Router("/cicd/jks/v1/BuildList", &controllers.CicdController{}, "*:BuildList")
	beego.Router("/cicd/jks/v1/BuildLog", &controllers.CicdController{}, "*:BuildLog")
	beego.Router("/cicd/jks/v1/BuildState", &controllers.CicdController{}, "*:BuildState") //构建状态

	beego.Router("/cicd/v1/GetCicdInfo", &controllers.CicdController{}, "*:GetCicdInfo")   //获取CICD信息 | Get CICD info
	beego.Router("/cicd/v1/GetPipelines", &controllers.CicdController{}, "*:GetPipelines") //获取流水线 | Get pipelines

	beego.Router("/cicd/v1/ListAppname", &controllers.CicdController{}, "*:ListAppname")                   //应用名称列表 | App name list
	beego.Router("/cicd/v1/PostStatus", &controllers.CicdController{}, "*:PostStatus")                     //提交状态 | Post status
	beego.Router("/cicd/v1/GetAliyunIdList", &controllers.CicdController{}, "*:GetAliyunIdList")           //获取阿里云ID列表 | Get Aliyun ID list
	beego.Router("/cicd/v1/GetOrganizationsByAk", &controllers.CicdController{}, "*:GetOrganizationsByAk") //通过访问密钥获取组织 | Get organizations by access key

	beego.Router("/cicd/pipeline/Start", &controllers.CicdController{}, "*:Start")         //启动流水线 | Start pipeline
	beego.Router("/cicd/pipeline/ListRun", &controllers.CicdController{}, "*:ListRun")     //运行列表 | Run list
	beego.Router("/cicd/pipeline/GetRun", &controllers.CicdController{}, "*:GetRun")       //获取运行信息 | Get run info
	beego.Router("/cicd/pipeline/GetJobLog", &controllers.CicdController{}, "*:GetJobLog") //获取作业日志 | Get job log

	//Wiki相关路由 | Wiki related routes
	beego.Router("/wiki/v1/List", &controllers.WikiController{}, "*:List")           //Wiki列表 | Wiki list
	beego.Router("/wiki/v1/Add", &controllers.WikiController{}, "*:Add")             //添加Wiki | Add Wiki
	beego.Router("/wiki/v1/Update", &controllers.WikiController{}, "*:Update")       //更新Wiki | Update Wiki
	beego.Router("/wiki/v1/Del", &controllers.WikiController{}, "*:Del")             //删除Wiki | Delete Wiki
	beego.Router("/wiki/v1/Read", &controllers.WikiController{}, "*:Read")           //读取Wiki | Read Wiki
	beego.Router("/wiki/v1/ReadEncry", &controllers.WikiController{}, "*:ReadEncry") //读取加密Wiki | Read encrypted Wiki
	beego.Router("/wiki/v1/Upload", &controllers.WikiController{}, "*:Upload")       //上传Wiki | Upload Wiki

	//收藏夹相关路由 | Favorite related routes
	beego.Router("/fav/v1/List", &controllers.FavController{}, "*:List") //收藏列表 | Favorite list
	beego.Router("/fav/v1/Add", &controllers.FavController{}, "*:Add")   //添加收藏 | Add favorite
	beego.Router("/fav/v1/Del", &controllers.FavController{}, "*:Del")   //删除收藏 | Delete favorite
}
