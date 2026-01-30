package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"xkube/common"

	"k8s.io/apimachinery/pkg/api/errors"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"

	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/util/retry"
)

type LabelsKv struct {
	// Key represents the label key | Key表示标签的键
	Key string `json:"key"`
	// Value represents the label value | Value表示标签的值
	Value string `json:"value"`
}

type EnvKv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type HostKv struct {
	Domain string `json:"domain"`
	Ip     string `json:"ip"`
}

type EnvSt struct {
	ContainerId   int     `json:"containerId"`
	ContainerName string  `json:"containerName"`
	Envs          []EnvKv `json:"envs"`
}

type LifecycleSt struct {
	ContainerId             int    `json:"containerId"`
	ContainerName           string `json:"containerName"`
	PreStop_type            string `json:"preStop_type"`
	PreStop_execCommand     string `json:"preStop_execCommand"`
	PreStop_httpGetPath     string `json:"preStop_httpGetPath"`
	PreStop_httpGetPort     int32  `json:"preStop_httpGetPort"`
	PreStop_httpGetScheme   string `json:"preStop_httpGetScheme"`
	PreStop_tcpSocketPort   int32  `json:"preStop_tcpSocketPort"`
	PostStart_type          string `json:"postStart_type"`
	PostStart_execCommand   string `json:"postStart_execCommand"`
	PostStart_httpGetPath   string `json:"postStart_httpGetPath"`
	PostStart_httpGetPort   int32  `json:"postStart_httpGetPort"`
	PostStart_httpGetScheme string `json:"postStart_httpGetScheme"`
	PostStart_tcpSocketPort int32  `json:"postStart_tcpSocketPort"`
}

type ResourceSt struct {
	ContainerId   int    `json:"containerId"`
	ContainerName string `json:"containerName"`
	Request_cpu   string `json:"request_cpu"`
	Request_mem   string `json:"request_mem"`
	Limit_cpu     string `json:"limit_cpu"`
	Limit_mem     string `json:"limit_mem"`
}

type ProbeST struct {
	ContainerId                   int    `json:"containerId"`
	ContainerName                 string `json:"containerName"`
	Readiness_checkType           string `json:"readiness_checkType"`
	Readiness_path                string `json:"readiness_path"`
	Readiness_httpPort            int32  `json:"readiness_httpPort"`
	Readiness_tcpPort             int32  `json:"readiness_tcpPort"`
	Readiness_cmd                 string `json:"readiness_cmd"`
	Readiness_initialDelaySeconds int32  `json:"readiness_initialDelaySeconds"`
	Readiness_periodSeconds       int32  `json:"readiness_periodSeconds"`
	Readiness_successThreshold    int32  `json:"readiness_successThreshold"`
	Readiness_failureThreshold    int32  `json:"readiness_failureThreshold"`
	Readiness_timeoutSeconds      int32  `json:"readiness_timeoutSeconds"`
	Liveness_checkType            string `json:"liveness_checkType"`
	Liveness_path                 string `json:"liveness_path"`
	Liveness_httpPort             int32  `json:"liveness_httpPort"`
	Liveness_tcpPort              int32  `json:"liveness_tcpPort"`
	Liveness_cmd                  string `json:"liveness_cmd"`
	Liveness_initialDelaySeconds  int32  `json:"liveness_initialDelaySeconds"`
	Liveness_periodSeconds        int32  `json:"liveness_periodSeconds"`
	Liveness_successThreshold     int32  `json:"liveness_successThreshold"`
	Liveness_failureThreshold     int32  `json:"liveness_failureThreshold"`
	Liveness_timeoutSeconds       int32  `json:"liveness_timeoutSeconds"`
}

type ImageKv struct {
	// ContainerName represents the name of the container | ContainerName表示容器的名称
	ContainerName string `json:"containerName"`
	// ContainerId represents the ID of the container | ContainerId表示容器的ID
	ContainerId int `json:"containerId"`
	// Image represents the image of the container | Image表示容器的镜像
	Image string `json:"image"`
}

type Deploy struct {
	// DeployName represents the name of the deployment | DeployName表示Deployment的名称
	DeployName string `json:"deployName" form:"deployName"`
	// NameSpace represents the namespace of the deployment | NameSpace表示Deployment的命名空间
	NameSpace string `json:"nameSpace" form:"nameSpace"`
	//CreationTimestamp    time.Time `json:"created_time"`
	// RevisionHistoryLimit represents the number of old ReplicaSets to retain | RevisionHistoryLimit表示保留的旧ReplicaSets数量
	RevisionHistoryLimit int32 `json:"revisionHistoryLimit" form:"historyVersionLimit"`
	// Replicas represents the desired number of replicas | Replicas表示期望的副本数量
	Replicas int32 `json:"replicas" form:"replicas"`
	// AvailableReplicas represents the number of available replicas | AvailableReplicas表示可用的副本数量
	AvailableReplicas int32 `json:"availableReplicas"`
	// PodNumber represents the pod number in format "available/replicas" | PodNumber表示Pod数量，格式为"可用/副本"
	PodNumber string `json:"podNumber"`
	// Labels represents the labels of the deployment | Labels表示Deployment的标签
	Labels string `json:"labels" form:"labels"`
	// ContainerName represents the name of the container | ContainerName表示容器名称
	ContainerName string `form:"containerName"`
	// ImageUrl represents the image url of the container | ImageUrl表示容器的镜像地址
	ImageUrl string `json:"imageUrl" form:"imageUrl"`
	// ContainerPortName represents the name of the container port | ContainerPortName表示容器端口名称
	ContainerPortName string `json:"-" form:"containerPortName"`
	// ContainerPort represents the container port | ContainerPort表示容器端口
	ContainerPort string `form:"containerPort"`
	// HostPort represents the host port | HostPort表示主机端口
	HostPort int32 `form:"hostPort"`
	// Resource             string    `json:"resource"`
	// StrategyType         string    `json:"strategyType"`
	// CreateTime represents the creation time of the deployment | CreateTime表示Deployment的创建时间
	CreateTime string `json:"createTime"` //创建时间
}

type DeployDetails struct {
	// DeployName represents the name of the deployment | DeployName表示Deployment的名称
	DeployName string `json:"deployName"`
	// NameSpace represents the namespace of the deployment | NameSpace表示Deployment的命名空间
	NameSpace string `json:"nameSpace"`
	// Strategy represents the deployment strategy | Strategy表示Deployment策略
	Strategy string `json:"strategy"`
	// StrategyRollingUpdate represents the rolling update strategy details | StrategyRollingUpdate表示滚动更新策略详情
	StrategyRollingUpdate string `json:"strategyRollingUpdate"`
	// PodNumber represents the number of pods | PodNumber表示Pod数量
	PodNumber string `json:"podNumber"`
	// Selector represents the selector of the deployment | Selector表示Deployment的选择器
	Selector string `json:"selector"`
	// ImageUrl represents the image url of the container | ImageUrl表示容器的镜像地址
	ImageUrl string `json:"imageUrl"`
	// Ports represents the ports of the container | Ports表示容器的端口
	Ports string `json:"ports"`
	// Annotations represents the annotations of the deployment | Annotations表示Deployment的注解
	Annotations string `json:"annotations"`
	// Status represents the status of the deployment | Status表示Deployment的状态
	Status string `json:"status"`
	// Labels represents the labels of the deployment | Labels表示Deployment的标签
	Labels string `json:"labels"`
	// CreateTime represents the creation time of the deployment | CreateTime表示Deployment的创建时间
	CreateTime string `json:"createTime"`
	// Conditions represents the conditions of the deployment | Conditions表示Deployment的条件
	Conditions []StatusConditions `json:"conditions"`
	// Replicasets represents the replica sets of the deployment | Replicasets表示Deployment的副本集
	Replicasets []Replicaset `json:"replicasets"`
}

type DeployStatus struct {
	// AvailableReplicas represents the number of available replicas | AvailableReplicas表示可用的副本数量
	AvailableReplicas string `json:"availableReplicas"`
	// ObservedGeneration represents the observed generation | ObservedGeneration表示观察到的版本
	ObservedGeneration string `json:"observedGeneration"`
	// ReadyReplicas represents the number of ready replicas | ReadyReplicas表示就绪的副本数量
	ReadyReplicas string `json:"readyReplicas"`
	// Replicas represents the number of replicas | Replicas表示副本数量
	Replicas string `json:"replicas"`
	// UpdatedReplicas represents the number of updated replicas | UpdatedReplicas表示已更新的副本数量
	UpdatedReplicas string `json:"updatedReplicas"`
	// Conditions represents the conditions of the deployment | Conditions表示Deployment的条件
	Conditions []StatusConditions
}

type StatusConditions struct {
	// LastTransitionTime represents the last transition time | LastTransitionTime表示最后转换时间
	LastTransitionTime string `json:"lastTransitionTime"`
	// LastUpdateTime represents the last update time | LastUpdateTime表示最后更新时间
	LastUpdateTime string `json:"lastUpdateTime"`
	// Message represents the message of the condition | Message表示条件的消息
	Message string `json:"message"`
	// Reason represents the reason of the condition | Reason表示条件的原因
	Reason string `json:"reason"`
	// Status represents the status of the condition | Status表示条件的状态
	Status string `json:"status"`
	// Ctype represents the type of the condition | Ctype表示条件的类型
	Ctype string `json:"ctype"`
}

type Replicaset struct {
	// ReplicasetName represents the name of the replica set | ReplicasetName表示副本集的名称
	ReplicasetName string `json:"replicasetName"`
	// ImageUrl represents the image url of the container | ImageUrl表示容器的镜像地址
	ImageUrl string `json:"imageUrl"`
	// CreateTime represents the creation time of the replica set | CreateTime表示副本集的创建时间
	CreateTime string `json:"createTime"`
}

type StrategyST struct {
	// StrategyType represents the type of deployment strategy | StrategyType表示Deployment策略类型
	StrategyType string `json:"strategyType"`
	// MinReadySeconds represents the minimum ready seconds | MinReadySeconds表示最小就绪秒数
	MinReadySeconds int32 `json:"minReadySeconds"`
	// MaxSurge represents the maximum surge | MaxSurge表示最大激增
	MaxSurge string `json:"maxSurge"`
	// MaxUnavailable represents the maximum unavailable | MaxUnavailable表示最大不可用
	MaxUnavailable string `json:"maxUnavailable"`
}

type TolerationSt struct {
	Key               string `json:"key"`
	Operator          string `json:"operator"`
	Value             string `json:"value"`
	Effect            string `json:"effect"`
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty"`
}

// affinity
type Expression struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Value    []string `json:"value"`
}

// affinity
type MatchExpression struct {
	Expressions []Expression `json:"expressions"`
}

// affinity
type Preference struct {
	Weight      int32        `json:"weight"`
	Expressions []Expression `json:"expressions"`
}

// affinity
type DeployNodeAffinitySt struct {
	NodeNames         string            `json:"nodeNames"`
	NodeSelector      []LabelsKv        `json:"nodeSelector"`
	RequiredAffinity  []MatchExpression `json:"requiredAffinity"`
	PreferredAffinity []Preference      `json:"preferredAffinity"`
}

// DeployList returns a list of deployments based on the provided filters
// DeployList根据提供的过滤条件返回Deployment列表
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace to list deployments in | namespace: 要列出Deployment的命名空间
//   - deployName: the name of the deployment to filter by | deployName: 用于过滤的Deployment名称
//   - labelsKey: the label key to filter by | labelsKey: 用于过滤的标签键
//   - labelsValue: the label value to filter by | labelsValue: 用于过滤的标签值
//
// Returns:
//   - []Deploy: a list of deployments | []Deploy: Deployment列表
//   - error: any error that occurred | error: 发生的任何错误
func DeployList(kubeconfig, namespace, deployName string, labelsKey, labelsValue string) ([]Deploy, error) {

	//通过kubeconfig集群认证文件生成一个客户端操作对象clientset
	clientset := common.ClientSet(kubeconfig)

	//创建一个deployment资源的接口对象DeploymentClient，用于操作指定名称空间的deployment资源
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}
	DeploymentClient := clientset.AppsV1().Deployments(namespace)
	//调用接口对象DeploymentClient中的Get方法，获取相应的deployment资源数据	deploymentInstance,err:=DeploymentClient.Get(context.TODO(),deploymentName,metaV1.GetOptions{})

	//设置ListOptions
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}
	var bbb = make([]Deploy, 0)
	deplist, err := DeploymentClient.List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list deployment error, err:%v", err)
		return bbb, err
	}

	for _, dep := range deplist.Items {
		//搜索
		if deployName != "" {
			if !strings.Contains(dep.Name, deployName) {
				continue
			}
		}
		//用于判断是否有appanme标签，并根据标签记录对应的资源到应用下
		//var appName string

		var labelsStr, imgUrlStr, containerNameStr, ContainerPortNameStr, ContainerPortStr string
		for kk, vv := range dep.ObjectMeta.Labels {
			// if kk == "appname" {
			// 	appName = vv
			// }
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		for _, v2 := range dep.Spec.Template.Spec.Containers {
			containerNameStr += fmt.Sprintf("%s,", v2.Name)
			imgUrlStr += fmt.Sprintf("%s,", v2.Image)
			if len(v2.Ports) > 0 {
				ContainerPortNameStr += fmt.Sprintf("%s,", v2.Ports[0].Name)
				ContainerPortStr += fmt.Sprintf("%d,", v2.Ports[0].ContainerPort)
			}
		}
		if len(containerNameStr) > 0 {
			containerNameStr = containerNameStr[0 : len(containerNameStr)-1]
		}
		if len(imgUrlStr) > 0 {
			imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
		}
		if len(ContainerPortNameStr) > 0 {
			ContainerPortNameStr = ContainerPortNameStr[0 : len(ContainerPortNameStr)-1]
		}
		if len(ContainerPortStr) > 0 {
			ContainerPortStr = ContainerPortStr[0 : len(ContainerPortStr)-1]
		}

		depItems := &Deploy{
			DeployName: dep.Name,
			NameSpace:  dep.Namespace,
			//CreationTimestamp:    dep.CreationTimestamp.Time,
			RevisionHistoryLimit: *dep.Spec.RevisionHistoryLimit,
			Replicas:             *dep.Spec.Replicas,
			AvailableReplicas:    dep.Status.AvailableReplicas,
			PodNumber:            fmt.Sprintf("%d/%d", dep.Status.AvailableReplicas, *dep.Spec.Replicas),
			Labels:               labelsStr,
			ContainerName:        containerNameStr,
			ImageUrl:             imgUrlStr,
			ContainerPortName:    ContainerPortNameStr,
			ContainerPort:        ContainerPortStr,
			//HostPort:             dep.Spec.Template.Spec.Containers[0].Ports[0].HostPort,
			// Resource:			dep.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().
			// StrategyType:
			CreateTime: dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *depItems)
	}
	return bbb, err
}

// DeployDetail returns detailed information about a specific deployment
// DeployDetail返回特定Deployment的详细信息
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//
// Returns:
//   - *DeployDetails: detailed information about the deployment | *DeployDetails: Deployment的详细信息
//   - error: any error that occurred | error: 发生的任何错误
func DeployDetail(kubeconfig, namespace, deployName string) (*DeployDetails, error) {
	DeploymentClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	dep, err := DeploymentClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		return &DeployDetails{}, err
	}

	var imgUrlStr, portStr string
	if len(dep.Spec.Template.Spec.Containers) > 0 {
		for _, vc1 := range dep.Spec.Template.Spec.Containers {
			imgUrlStr += fmt.Sprintf("%s,", vc1.Image)
			if len(vc1.Ports) > 0 {
				for _, vp1 := range vc1.Ports {
					portStr += fmt.Sprintf("%d,", vp1.ContainerPort)
				}
			}
		}
		imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
		if len(portStr) > 0 {
			portStr = portStr[0 : len(portStr)-1]
		}
	}

	//return返回指针类型的deployment结构体类型的实例
	var selectorStr, labelsStr, annotationsStr string
	for kk, vv := range dep.Spec.Selector.MatchLabels {
		selectorStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(selectorStr) > 0 {
		selectorStr = selectorStr[0 : len(selectorStr)-1]
	}

	for kk, vv := range dep.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	for kk, vv := range dep.ObjectMeta.Annotations {
		if strings.Contains(kk, "last-applied-configuration") {
			continue
		}
		annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	var bbb = make([]StatusConditions, 0)

	for _, v1 := range dep.Status.Conditions {
		xItems := &StatusConditions{
			LastTransitionTime: v1.LastTransitionTime.Format("2006-01-02 15:04:05"),
			LastUpdateTime:     v1.LastUpdateTime.Format("2006-01-02 15:04:05"),
			Message:            v1.Message,
			Reason:             v1.Reason,
			Status:             fmt.Sprintf("%v", v1.Status),
			Ctype:              fmt.Sprintf("%v", v1.Type),
		}
		bbb = append(bbb, *xItems)
	}

	var maxSurge, maxUnavailable = "0", "0"
	if dep.Spec.Strategy.Type == "RollingUpdate" { //RollingUpdate or Recreate
		maxSurge = dep.Spec.Strategy.RollingUpdate.MaxSurge.String()
		maxUnavailable = dep.Spec.Strategy.RollingUpdate.MaxUnavailable.String()
	}

	//replicaset list
	selector, err := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	if err != nil {
		panic(err.Error())
	}
	rsList, err := common.ClientSet(kubeconfig).AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		panic(err.Error())
	}
	var repcs = make([]Replicaset, 0)
	for _, rs := range rsList.Items {
		xItems := &Replicaset{
			ReplicasetName: rs.Name,
			ImageUrl:       rs.Spec.Template.Spec.Containers[0].Image,
			CreateTime:     rs.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		repcs = append(repcs, *xItems)
	}

	return &DeployDetails{
		DeployName:            dep.Name,
		NameSpace:             dep.Namespace,
		Strategy:              string(dep.Spec.Strategy.Type),
		StrategyRollingUpdate: fmt.Sprintf("超过期望的Pod数量:%s,不可用Pod最大数量:%s", maxSurge, maxUnavailable),
		PodNumber:             fmt.Sprintf("%d", *dep.Spec.Replicas),
		Selector:              selectorStr,
		ImageUrl:              imgUrlStr,
		Ports:                 portStr,
		Annotations:           annotationsStr,
		Status:                fmt.Sprintf("就绪:%d个,已更新:%d个,可用:%d个", dep.Status.ReadyReplicas, dep.Status.UpdatedReplicas, dep.Status.AvailableReplicas),
		Labels:                labelsStr,
		CreateTime:            dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Conditions:            bbb,
		Replicasets:           repcs,
	}, nil
}

// DeployCreate creates a new deployment based on the provided configuration
// DeployCreate根据提供的配置创建新的Deployment
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - bodys: the JSON byte array containing deployment configuration | bodys: 包含Deployment配置的JSON字节数组
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
func DeployCreate(kubeconfig string, bodys []byte) (string, error) {

	gp := gjson.ParseBytes(bodys)

	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	deployName := gp.Get("deployName").String()
	nameSpace := gp.Get("nameSpace").String()
	containerPort := gp.Get("containerPort").Int()
	replicas := gp.Get("replicas").Int()
	var pullPolicy corev1.PullPolicy
	imagePullPolicy := gp.Get("imagePullPolicy").String()
	switch imagePullPolicy {
	case "Never":
		pullPolicy = corev1.PullNever
	case "IfNotPresent":
		pullPolicy = corev1.PullIfNotPresent
	default:
		pullPolicy = corev1.PullAlways
	}
	imageUrl := gp.Get("imageUrl").String()

	resourceLimitCheck := gp.Get("resourceLimitCheck").String()
	healthCheck := gp.Get("healthCheck").String()
	periodCheck := gp.Get("periodCheck").String()
	variableCheck := gp.Get("variableCheck").String()
	serviceCheck := gp.Get("serviceCheck").String()

	livenessProbe := gp.Get("livenessProbe").String()
	readinessProbe := gp.Get("readinessProbe").String()
	readiness_checkType := gp.Get("readiness_checkType").String()
	liveness_checkType := gp.Get("liveness_checkType").String()

	//var labelsMap = make(map[string]string)
	labelsMap := map[string]string{
		"app": deployName,
	}
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").String()] = vv.Get("value").String()
	}

	deployInstance := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(int32(replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labelsMap,
			},
			RevisionHistoryLimit: int32Ptr(10),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					//Name: "golang-pod",
					Labels: labelsMap,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:            deployName,
							Image:           imageUrl,
							ImagePullPolicy: pullPolicy,
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									//Name:     ""     ,
									ContainerPort: int32(containerPort),
									//HostPort:      deployment.HostPort,
									Protocol: corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	if resourceLimitCheck == "on" {
		limit_cpu := gp.Get("limit_cpu").Str
		limit_mem := gp.Get("limit_mem").Str
		request_cpu := gp.Get("request_cpu").Str
		request_mem := gp.Get("request_mem").Str
		resReq := &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(limit_cpu),
				corev1.ResourceMemory: resource.MustParse(limit_mem),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(request_cpu),
				corev1.ResourceMemory: resource.MustParse(request_mem),
			},
		}
		deployInstance.Spec.Template.Spec.Containers[0].Resources = *resReq
	}

	if healthCheck == "on" {
		if readinessProbe == "on" {
			readiness_path := gp.Get("readiness_path").String()
			readiness_httpPort := int32(gp.Get("readiness_httpPort").Int())
			readiness_tcpPort := int32(gp.Get("readiness_tcpPort").Int())
			readiness_cmd := gp.Get("readiness_cmd").String()
			cmdStr := strings.Split(readiness_cmd, " ")
			readiness_initialDelaySeconds := gp.Get("readiness_initialDelaySeconds").Int()
			readiness_periodSeconds := gp.Get("readiness_periodSeconds").Int()
			readiness_successThreshold := gp.Get("readiness_successThreshold").Int()
			readiness_failureThreshold := gp.Get("readiness_failureThreshold").Int()
			readiness_timeoutSeconds := gp.Get("readiness_timeoutSeconds").Int()

			httpGet := &corev1.HTTPGetAction{}
			tcpSocket := &corev1.TCPSocketAction{}
			execCmd := &corev1.ExecAction{}
			probeHandler := &corev1.ProbeHandler{}
			if readiness_checkType == "HTTP" {
				httpGet = &corev1.HTTPGetAction{
					Path: readiness_path,
					Port: intstr.FromInt32(readiness_httpPort),
				}
				probeHandler.HTTPGet = httpGet
			}
			if readiness_checkType == "TCP" {
				tcpSocket = &corev1.TCPSocketAction{
					Port: intstr.FromInt32(readiness_tcpPort),
				}
				probeHandler.TCPSocket = tcpSocket
			}
			if readiness_checkType == "CMD" {
				execCmd = &corev1.ExecAction{
					Command: cmdStr,
				}
				probeHandler.Exec = execCmd
			}

			probe := &corev1.Probe{
				InitialDelaySeconds: int32(readiness_initialDelaySeconds),
				PeriodSeconds:       int32(readiness_periodSeconds),
				TimeoutSeconds:      int32(readiness_timeoutSeconds),
				SuccessThreshold:    int32(readiness_successThreshold),
				FailureThreshold:    int32(readiness_failureThreshold),
				ProbeHandler:        *probeHandler,
			}
			deployInstance.Spec.Template.Spec.Containers[0].ReadinessProbe = probe
		}

		if livenessProbe == "on" {
			liveness_path := gp.Get("liveness_path").String()
			liveness_httpPort := int32(gp.Get("liveness_httpPort").Int())
			liveness_tcpPort := int32(gp.Get("liveness_tcpPort").Int())
			liveness_cmd := gp.Get("liveness_cmd").String()
			cmdStr := strings.Split(liveness_cmd, " ")
			liveness_initialDelaySeconds := gp.Get("liveness_initialDelaySeconds").Int()
			liveness_periodSeconds := gp.Get("liveness_periodSeconds").Int()
			liveness_successThreshold := gp.Get("liveness_successThreshold").Int()
			liveness_failureThreshold := gp.Get("liveness_failureThreshold").Int()
			liveness_timeoutSeconds := gp.Get("liveness_timeoutSeconds").Int()

			httpGet := &corev1.HTTPGetAction{}
			tcpSocket := &corev1.TCPSocketAction{}
			execCmd := &corev1.ExecAction{}
			probeHandler := &corev1.ProbeHandler{}
			if liveness_checkType == "HTTP" {
				httpGet = &corev1.HTTPGetAction{
					Path: liveness_path,
					Port: intstr.FromInt32(liveness_httpPort),
				}
				probeHandler.HTTPGet = httpGet
			}
			if liveness_checkType == "TCP" {
				tcpSocket = &corev1.TCPSocketAction{
					Port: intstr.FromInt32(liveness_tcpPort),
				}
				probeHandler.TCPSocket = tcpSocket
			}
			if liveness_checkType == "CMD" {
				execCmd = &corev1.ExecAction{
					Command: cmdStr,
				}
				probeHandler.Exec = execCmd
			}

			probe := &corev1.Probe{
				InitialDelaySeconds: int32(liveness_initialDelaySeconds),
				PeriodSeconds:       int32(liveness_periodSeconds),
				TimeoutSeconds:      int32(liveness_timeoutSeconds),
				SuccessThreshold:    int32(liveness_successThreshold), //必须为1
				FailureThreshold:    int32(liveness_failureThreshold),
				ProbeHandler:        *probeHandler,
			}
			deployInstance.Spec.Template.Spec.Containers[0].LivenessProbe = probe
		}
	}

	if periodCheck == "on" {
		var postStartArry, preStopArry []string
		postStartStr := gp.Get("postStart").Str
		preStopStr := gp.Get("preStop").Str

		lifeCycle := &corev1.Lifecycle{}
		if postStartStr != "" {
			postStartArry = strings.Split(postStartStr, ",")
			lifeCycle.PostStart = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: postStartArry,
				},
			}
		}
		if preStopStr != "" {
			preStopArry = strings.Split(preStopStr, ",")
			lifeCycle.PreStop = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: preStopArry,
				},
			}
		}

		deployInstance.Spec.Template.Spec.Containers[0].Lifecycle = lifeCycle
	}

	if variableCheck == "on" {
		var envsArry []corev1.EnvVar
		for _, vv := range gp.Get("envs").Array() {
			envsArry = append(envsArry, corev1.EnvVar{
				Name:  vv.Get("key").String(),
				Value: vv.Get("value").String(),
			})
		}
		deployInstance.Spec.Template.Spec.Containers[0].Env = envsArry
	}

	if serviceCheck == "on" {
		var svcPortArry []ServicePort
		for _, vv := range gp.Get("svc_ports").Array() {
			svcPortArry = append(svcPortArry, ServicePort{
				PortName:   vv.Get("portName").String(),
				SvcPort:    int32(vv.Get("svcPort").Int()),
				TargetPort: int32(vv.Get("targetPort").Int()),
			})
		}
		service := ServiceMain{
			ServiceName: deployName,
			NameSpace:   nameSpace,
			SvcType:     gp.Get("svcType").String(),
			Ports:       svcPortArry,
			Labels:      labelsMap,
		}
		data, err := json.Marshal(service)
		if err != nil {
			return deployName, err
		}
		_ = SvcCreate(kubeconfig, data)
	}

	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(nameSpace)
	//调用DeployClient接口中的create方法，创建deployment资源
	_, err := DeployClient.Create(context.TODO(), deployInstance, metav1.CreateOptions{})
	if err != nil {
		return deployName, err
	}
	return deployName, nil
}

// CheckYamlSyntax checks if the provided YAML data has valid syntax
// CheckYamlSyntax检查提供的YAML数据是否有有效语法
// Parameters:
//   - yamlData: the YAML data to check | yamlData: 要检查的YAML数据
//
// Returns:
//   - error: any error that occurred during syntax checking | error: 语法检查期间发生的任何错误
//
// 能检测yaml语法是否正确，不能检测deploy的语法是否正确
func CheckYamlSyntax(yamlData []byte) error {
	var data interface{}
	err := yamlutil.UnmarshalStrict(yamlData, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %v", err)
	}
	return nil
}

// DeployYamlCreateV2 creates a deployment from YAML data
// DeployYamlCreateV2从YAML数据创建Deployment
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - yamlData: the YAML data containing deployment configuration | yamlData: 包含Deployment配置的YAML数据
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
//
// 测试ok 未用
func DeployYamlCreateV2(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	deployment := &appsv1.Deployment{}
	err = json.Unmarshal(data, deployment)
	if err != nil {
		return err
	}
	//cluster := deployment.ObjectMeta.ClusterName
	namespace := deployment.ObjectMeta.Namespace
	deploymentName := deployment.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	deploy, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, deploymentName)
	fmt.Println(deploy)
	return err
}

// DeployClone clones a deployment to a target cluster/namespace/name
// DeployClone将Deployment克隆到目标集群/命名空间/名称
// Parameters:
//   - kubeconfig: the source kubeconfig file path or cluster id | kubeconfig: 源kubeconfig文件路径或集群ID
//   - namespace: the source namespace | namespace: 源命名空间
//   - deployName: the source deployment name | deployName: 源Deployment名称
//   - target_clusterid: the target cluster id | target_clusterid: 目标集群ID
//   - target_namespace: the target namespace | target_namespace: 目标命名空间
//   - target_objname: the target deployment name | target_objname: 目标Deployment名称
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
func DeployClone(kubeconfig, namespace, deployName, target_clusterid, target_namespace, target_objname string) error {
	//old cluster
	//log.Println(kubeconfig, namespace, deployName, target_clusterid, target_namespace, target_objname)
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if target_namespace == "" && target_objname == "" && target_clusterid == "" {
		return fmt.Errorf("target can't be all empty")
	}

	if target_namespace == "" {
		target_namespace = namespace
	}

	if target_objname == "" {
		target_objname = deployName
	}

	if target_clusterid == "" {
		target_clusterid = kubeconfig
	}

	//在新集群创建namespace
	err3 := CreateNsByExist(target_clusterid, target_namespace)
	if err3 != nil {
		return err3
	}

	deployInstance.Name = target_objname
	deployInstance.Namespace = target_namespace
	deployInstance.ResourceVersion = ""
	deployInstance.ObjectMeta.Labels["app"] = target_objname
	deployInstance.Spec.Selector.MatchLabels["app"] = target_objname
	deployInstance.Spec.Template.Labels["app"] = target_objname

	//new cluster
	NewDeployClient := common.ClientSet(target_clusterid).AppsV1().Deployments(target_namespace)
	_, err2 := NewDeployClient.Get(context.TODO(), target_objname, metav1.GetOptions{})
	if errors.IsNotFound(err2) { //没有就创建
		_, err := NewDeployClient.Create(context.TODO(), deployInstance, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if err2 == nil {
		_, err = NewDeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	return err2
}

// DeployYamlModify modifies a deployment using YAML data
// DeployYamlModify使用YAML数据修改Deployment
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - yamlData: the YAML data containing updated deployment configuration | yamlData: 包含更新的Deployment配置的YAML数据
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
func DeployYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	deployment := &appsv1.Deployment{}
	err = json.Unmarshal(data, deployment)
	if err != nil {
		return err
	}
	//cluster := deployment.ObjectMeta.ClusterName
	namespace := deployment.ObjectMeta.Namespace
	//deploymentName := deployment.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	return err
}

// DeployModify modifies an existing deployment
// DeployModify修改现有的Deployment
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - newDeploy: the deployment object with updated values | newDeploy: 包含更新值的Deployment对象
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
func DeployModify(kubeconfig string, newDeploy *Deploy) error {
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(newDeploy.NameSpace)
	deployInstance, err := DeployClient.Get(context.TODO(), newDeploy.DeployName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ContainerPortInt, _ := strconv.Atoi(newDeploy.ContainerPort)
	deployInstance.Spec.Replicas = &newDeploy.Replicas
	deployInstance.Spec.RevisionHistoryLimit = &newDeploy.RevisionHistoryLimit
	deployInstance.Spec.Template.Spec.Containers[0].Image = newDeploy.ImageUrl
	deployInstance.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(ContainerPortInt)
	//deployInstance.Spec.Template.Spec.Containers[0].Ports[0].HostPort = newDeploy.HostPort
	//调用DeploymentClient接口对象中的update方法，来更新deployment资源
	_, err = DeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// DeployLabels gets or updates the labels of a deployment
// DeployLabels获取或更新Deployment的标签
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - labelsMap: the map of labels to set (for POST method) | labelsMap: 要设置的标签映射(用于POST方法)
//
// Returns:
//   - []LabelsKv: the list of labels (for GET method) | []LabelsKv: 标签列表(用于GET方法)
//   - error: any error that occurred | error: 发生的任何错误
func DeployLabels(kubeconfig, namespace, deployName string, method string, labelsMap map[string]string) ([]LabelsKv, error) {
	var bbb = make([]LabelsKv, 0)
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployLabels Get error:%s\n", err)
		return bbb, err
	}
	if method == "POST" {
		deployInstance.ObjectMeta.Labels = labelsMap
		_, err = DeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployLabels Update error:%s\n", err)
			return bbb, err
		}
		return bbb, nil
	} else {
		if deployInstance.ObjectMeta.Labels != nil {
			for kk, vv := range deployInstance.ObjectMeta.Labels {
				bbb = append(bbb, LabelsKv{
					Key:   kk,
					Value: vv,
				})
			}
		}
		return bbb, nil
	}
}

// DeployAnnotations gets a specific annotation from a deployment
// DeployAnnotations从Deployment获取特定注解
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - key: the annotation key to retrieve | key: 要检索的注解键
//
// Returns:
//   - string: the value of the annotation, or empty string if not found | string: 注解的值，如果未找到则为空字符串
func DeployAnnotations(kubeconfig, namespace, deployName, key string) string {
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployAnnotations Get error:%s\n", err)
		return ""
	}
	if deployInstance.ObjectMeta.Annotations != nil {
		if vv, ok := deployInstance.ObjectMeta.Annotations[key]; ok {
			return vv
		}
	}
	return ""
}

// DeployImage gets or updates the images of a deployment's containers
// DeployImage获取或更新Deployment容器的镜像
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - newImages: the list of new images to set (for POST method) | newImages: 要设置的新镜像列表(用于POST方法)
//
// Returns:
//   - []ImageKv: the list of current images (for GET method) | []ImageKv: 当前镜像列表(用于GET方法)
//   - error: any error that occurred | error: 发生的任何错误
func DeployImage(kubeconfig, namespace, deployName string, method string, newImages []ImageKv) ([]ImageKv, error) {
	var bbb = make([]ImageKv, 0)
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] Deployment Get error:%s\n", err)
		return bbb, err
	}
	if method == "POST" {
		for k1, v1 := range deployInstance.Spec.Template.Spec.Containers {
			for _, v2 := range newImages {
				log.Printf("[DEUBG] id:%d,name:%s,image:%s\n", k1, v1.Name, v2.Image)
				if v1.Name == v2.ContainerName {
					deployInstance.Spec.Template.Spec.Containers[k1].Image = v2.Image
				}
			}
		}
		_, err = DeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] Deployment Update image error:%s\n", err)
			return bbb, err
		}
		return bbb, nil
	} else {
		for kk, vv := range deployInstance.Spec.Template.Spec.Containers {
			bbb = append(bbb, ImageKv{
				ContainerName: vv.Name,
				ContainerId:   kk,
				Image:         vv.Image,
			})
		}
		return bbb, nil
	}
}

// DeploymentUpdateImage updates the image of a specific container in a deployment
// DeploymentUpdateImage更新Deployment中特定容器的镜像
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - containerName: the name of the container to update | containerName: 要更新的容器名称
//   - image: the new image to set | image: 要设置的新镜像
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
func DeploymentUpdateImage(kubeconfig, namespace, deployName, containerName, image string) error {
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeploymentUpdateImage Deployment Get error:%s\n", err)
		return err
	}

	for k1, v1 := range deployInstance.Spec.Template.Spec.Containers {
		log.Printf("[DEUBG] DeploymentUpdateImage id:%d,name:%s,image:%s,containerName:%s\n", k1, v1.Name, image, containerName)
		if v1.Name == containerName {
			deployInstance.Spec.Template.Spec.Containers[k1].Image = image
		}
	}
	_, err = DeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("[ERROR] DeploymentUpdateImage Deployment image Update error:%s\n", err)
		return err
	}
	return nil
}

// DeployReplicas gets or updates the number of replicas in a deployment
// DeployReplicas获取或更新Deployment中的副本数量
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - podNumber: the new number of replicas to set (for POST method) | podNumber: 要设置的新副本数量(用于POST方法)
//
// Returns:
//   - int32: the current number of replicas | int32: 当前副本数量
//   - error: any error that occurred | error: 发生的任何错误
func DeployReplicas(kubeconfig, namespace, deployName string, method string, podNumber int32) (int32, error) {
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] Deployreplicas Get error:%s\n", err)
		return 0, err
	}
	if method == "POST" {
		deployInstance.Spec.Replicas = &podNumber
		_, err = DeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] Deployreplicas Update error:%s\n", err)
			return 0, err
		}
		//log.Printf("Deployment %s replica set to %d\n", deployInstance.Name, *deployInstance.Spec.Replicas)
		return podNumber, nil
	} else {
		podnum := deployInstance.Spec.Replicas
		return *podnum, nil
	}
}

// DeployStrategy gets or updates the deployment strategy
// DeployStrategy获取或更新Deployment策略
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - sst: the strategy settings to set (for POST method) | sst: 要设置的策略设置(用于POST方法)
//
// Returns:
//   - StrategyST: the current strategy settings | StrategyST: 当前策略设置
//   - error: any error that occurred | error: 发生的任何错误
func DeployStrategy(kubeconfig, namespace, deployName string, method string, sst StrategyST) (StrategyST, error) {
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := DeployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployLabels Get error:%s\n", err)
		return sst, err
	}
	if method == "POST" {
		if sst.StrategyType == "RollingUpdate" {
			sst.MaxSurge = strings.ReplaceAll(sst.MaxSurge, "%25", "%")
			sst.MaxUnavailable = strings.ReplaceAll(sst.MaxUnavailable, "%25", "%")

			deployInstance.Spec.Strategy.Type = appsv1.RollingUpdateDeploymentStrategyType
			deployInstance.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
				MaxSurge: &intstr.IntOrString{Type: intstr.String, StrVal: sst.MaxSurge},
				//MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: int32(0)},
				MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: sst.MaxUnavailable},
			}
		} else {
			//deployInstance.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType
			deployInstance.Spec.Strategy = appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			}
		}

		deployInstance.Spec.MinReadySeconds = sst.MinReadySeconds

		_, err = DeployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployStrategy Update error:%s\n", err)
			return sst, err
		}
		//log.Printf("Deployment %s replica set to %d\n", deployInstance.Name, *deployInstance.Spec.Replicas)
		return sst, nil
	} else {
		strategyType := deployInstance.Spec.Strategy.Type
		minReadySeconds := deployInstance.Spec.MinReadySeconds
		var maxSurge, maxUnavailable string
		if strategyType == "RollingUpdate" {
			maxSurge = deployInstance.Spec.Strategy.RollingUpdate.MaxSurge.StrVal
			maxUnavailable = deployInstance.Spec.Strategy.RollingUpdate.MaxUnavailable.StrVal
		}

		return StrategyST{
			StrategyType:    fmt.Sprintf("%v", strategyType),
			MinReadySeconds: minReadySeconds,
			MaxSurge:        maxSurge,
			MaxUnavailable:  maxUnavailable,
		}, nil
	}
}

// DeployDel deletes a deployment
// DeployDel删除一个Deployment
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deploy: the name of the deployment to delete | deploy: 要删除的Deployment名称
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
func DeployDel(kubeconfig, namespace, deploy string) error {
	DeployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	err := DeployClient.Delete(context.TODO(), deploy, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetDeployYaml returns the YAML representation of a deployment
// GetDeployYaml返回Deployment的YAML表示
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deploy: the name of the deployment | deploy: Deployment的名称
//
// Returns:
//   - string: the YAML representation of the deployment | string: Deployment的YAML表示
//   - error: any error that occurred | error: 发生的任何错误
func GetDeployYaml(kubeconfig, namespace, deploy string) (string, error) {

	DeploymentClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployment, err := DeploymentClient.Get(context.TODO(), deploy, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	deploymentUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(deployment)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(deploymentUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
	// 将yaml配置文件写入文件
	// err = ioutil.WriteFile(fmt.Sprintf("%s.yaml", deploymentName), yamlBytes, 0644)
	// if err != nil {
	// 	panic(err.Error())
	// }
}

// GetReplicasetYaml returns the YAML representation of a replica set
// GetReplicasetYaml返回副本集的YAML表示
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the replica set | namespace: 副本集的命名空间
//   - replicaSet: the name of the replica set | replicaSet: 副本集的名称
//
// Returns:
//   - string: the YAML representation of the replica set | string: 副本集的YAML表示
//   - error: any error that occurred | error: 发生的任何错误
//
// replicaset
func GetReplicasetYaml(kubeconfig, namespace, replicaSet string) (string, error) {
	repset, err := common.ClientSet(kubeconfig).AppsV1().ReplicaSets(namespace).Get(context.TODO(), replicaSet, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	deploymentUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(repset)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(deploymentUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// GetReplicasetList returns a list of replica sets based on the provided filters
// GetReplicasetList根据提供的过滤条件返回副本集列表
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace to list replica sets in | namespace: 要列出副本集的命名空间
//   - labelsKey: the label key to filter by | labelsKey: 用于过滤的标签键
//   - labelsValue: the label value to filter by | labelsValue: 用于过滤的标签值
//
// Returns:
//   - []Replicaset: a list of replica sets | []Replicaset: 副本集列表
//   - error: any error that occurred | error: 发生的任何错误
//
// replicaset list
func GetReplicasetList(kubeconfig, namespace, labelsKey, labelsValue string) ([]Replicaset, error) {
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	rsList, err := common.ClientSet(kubeconfig).AppsV1().ReplicaSets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("[ERROR] GetReplicaSetList list Fail:%s\n", err)
	}
	var repcs = make([]Replicaset, 0)
	for _, rs := range rsList.Items {
		xItems := &Replicaset{
			ReplicasetName: rs.Name,
			ImageUrl:       rs.Spec.Template.Spec.Containers[0].Image,
			CreateTime:     rs.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		repcs = append(repcs, *xItems)
	}
	return repcs, err
}

// DeployRollBack rolls back a deployment to a specific replica set
// DeployRollBack将Deployment回滚到特定的副本集
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - replicaSet: the name of the replica set to roll back to | replicaSet: 要回滚到的副本集名称
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
//
// RollBack
func DeployRollBack(kubeconfig, namespace, deployName, replicaSet string) error {

	deploy, err := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace).Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	repset, err1 := common.ClientSet(kubeconfig).AppsV1().ReplicaSets(namespace).Get(context.TODO(), replicaSet, metav1.GetOptions{})
	if err1 != nil {
		return err1
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deploy.Spec.Template = repset.Spec.Template
		_, updateErr := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
		return updateErr
	})
	return retryErr
}

// DeployRestart restarts a deployment by updating its annotation
// DeployRestart通过更新注解来重启Deployment
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//
// Returns:
//   - error: any error that occurred | error: 发生的任何错误
//
// Restart
func DeployRestart(kubeconfig, namespace, deployName string) error {
	patchOpt := metav1.PatchOptions{FieldManager: "kubectl-rollout"}
	patchInfo := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
	_, err := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace).Patch(context.TODO(), deployName, types.StrategicMergePatchType, []byte(patchInfo), patchOpt)
	return err
}

// int32Ptr returns a pointer to the provided int32 value
// int32Ptr返回指向提供的int32值的指针
// Parameters:
//   - i: the int32 value to create a pointer for | i: 要创建指针的int32值
//
// Returns:
//   - *int32: a pointer to the provided int32 value | *int32: 指向提供的int32值的指针
func int32Ptr(i int32) *int32 { return &i }

// DeployHost gets or updates the host aliases of a deployment
// DeployHost获取或更新Deployment的主机别名
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - hostAlias: the list of host aliases to set (for POST method) | hostAlias: 要设置的主机别名列表(用于POST方法)
//
// Returns:
//   - []HostKv: the list of current host aliases | []HostKv: 当前主机别名列表
//   - error: any error that occurred | error: 发生的任何错误
func DeployHost(kubeconfig, namespace, deployName string, method string, hostAlias []HostKv) ([]HostKv, error) {
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployHost Get deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
		return hostAlias, err
	}

	// 处理POST请求，更新Deployment的HostAliases配置
	// Handle POST request to update Deployment's HostAliases configuration
	if method == "POST" {
		// 验证输入参数
		// Validate input parameters
		for _, hv := range hostAlias {
			if hv.Ip == "" || hv.Domain == "" {
				return hostAlias, fmt.Errorf("invalid host alias: IP and Domain must not be empty")
			}
		}

		emptyHostAliases := make([]corev1.HostAlias, 0, len(hostAlias))
		for _, vv := range hostAlias {
			emptyHostAliases = append(emptyHostAliases, corev1.HostAlias{
				IP:        vv.Ip,
				Hostnames: strings.Split(vv.Domain, ","),
			})
		}

		deployInstance.Spec.Template.Spec.HostAliases = emptyHostAliases
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployHost Update deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
			return hostAlias, err
		}
		return hostAlias, nil
	}

	// GET 逻辑
	// GET logic
	if deployInstance.Spec.Template.Spec.HostAliases != nil {
		result := make([]HostKv, 0, len(deployInstance.Spec.Template.Spec.HostAliases))
		for _, vv := range deployInstance.Spec.Template.Spec.HostAliases {
			result = append(result, HostKv{
				Ip:     vv.IP,
				Domain: strings.Join(vv.Hostnames, ","),
			})
		}
		return result, nil
	}

	// 返回空切片而不是原hostAlias参数，保持一致性
	// Return empty slice instead of the original hostAlias parameter for consistency
	return []HostKv{}, nil
}

// DeployResource gets or updates the resource limits/requests of a deployment's container
// DeployResource获取或更新Deployment容器的资源限制/请求
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - containerName: the name of the container | containerName: 容器名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - resl: the resource limits/requests to set (for POST method) | resl: 要设置的资源限制/请求(用于POST方法)
//
// Returns:
//   - ResourceSt: the current resource limits/requests | ResourceSt: 当前资源限制/请求
//   - error: any error that occurred | error: 发生的任何错误
func DeployResource(kubeconfig, namespace, deployName, method string, resl ResourceSt) ([]ResourceSt, error) {
	var emptyResouceSt = make([]ResourceSt, 0)
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployHost Get deployment error: %s\n", err)
		return emptyResouceSt, err
	}

	// 处理POST请求，更新Deployment的资源限制配置
	// Handle POST request to update Deployment's resource limits configuration
	if method == "POST" {
		resReq := &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(resl.Limit_cpu),
				corev1.ResourceMemory: resource.MustParse(resl.Limit_mem),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(resl.Request_cpu),
				corev1.ResourceMemory: resource.MustParse(resl.Request_mem),
			},
		}
		deployInstance.Spec.Template.Spec.Containers[resl.ContainerId].Resources = *resReq
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployResLimit Update deployment error: %s\n", err)
			return emptyResouceSt, err
		}
		return emptyResouceSt, nil
	}

	lt := len(deployInstance.Spec.Template.Spec.Containers)
	var ccc = make([]ResourceSt, 0)
	vv := deployInstance.Spec.Template.Spec.Containers
	for i := 0; i < lt; i++ {
		result := ResourceSt{}
		result.ContainerId = i
		result.ContainerName = vv[i].Name
		// 检查是否有资源配置（检查Requests和Limits是否为空）
		// Check if there are resource configurations (check if Requests and Limits are empty)
		if len(vv[i].Resources.Requests) > 0 || len(vv[i].Resources.Limits) > 0 {
			// 安全地获取资源值
			// Safely get resource values
			if vv[i].Resources.Requests != nil {
				if cpu, exists := vv[i].Resources.Requests[corev1.ResourceCPU]; exists {
					result.Request_cpu = cpu.String()
				}
				if mem, exists := vv[i].Resources.Requests[corev1.ResourceMemory]; exists {
					result.Request_mem = mem.String()
				}
			}

			if vv[i].Resources.Limits != nil {
				if cpu, exists := vv[i].Resources.Limits[corev1.ResourceCPU]; exists {
					result.Limit_cpu = cpu.String()
				}
				if mem, exists := vv[i].Resources.Limits[corev1.ResourceMemory]; exists {
					result.Limit_mem = mem.String()
				}
			}
		}
		ccc = append(ccc, result)
	}
	return ccc, nil
}

func DeployProbe(kubeconfig, namespace, deployName, method string, probe ProbeST) ([]ProbeST, error) {

	var emptyProbe = make([]ProbeST, 0)

	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployHost Get deployment error: %s\n", err)
		return emptyProbe, err
	}

	if method == "POST" {
		if probe.Readiness_checkType != "" {
			probeHandler := &corev1.ProbeHandler{}
			if probe.Readiness_checkType == "HTTP" {
				probeHandler.HTTPGet = &corev1.HTTPGetAction{
					Path: probe.Readiness_path,
					Port: intstr.FromInt32(probe.Readiness_httpPort),
				}
			}
			if probe.Readiness_checkType == "TCP" {
				probeHandler.TCPSocket = &corev1.TCPSocketAction{
					Port: intstr.FromInt32(probe.Readiness_tcpPort),
				}
			}
			if probe.Readiness_checkType == "CMD" {
				probeHandler.Exec = &corev1.ExecAction{
					Command: strings.Split(probe.Readiness_cmd, " "),
				}
			}

			readinessProbe := &corev1.Probe{
				InitialDelaySeconds: int32(probe.Readiness_initialDelaySeconds),
				PeriodSeconds:       int32(probe.Readiness_periodSeconds),
				TimeoutSeconds:      int32(probe.Readiness_timeoutSeconds),
				SuccessThreshold:    int32(probe.Readiness_successThreshold),
				FailureThreshold:    int32(probe.Readiness_failureThreshold),
				ProbeHandler:        *probeHandler,
			}
			deployInstance.Spec.Template.Spec.Containers[probe.ContainerId].ReadinessProbe = readinessProbe
		}

		if probe.Liveness_checkType != "" {
			probeHandler := &corev1.ProbeHandler{}
			if probe.Liveness_checkType == "HTTP" {
				probeHandler.HTTPGet = &corev1.HTTPGetAction{
					Path: probe.Liveness_path,
					Port: intstr.FromInt32(probe.Liveness_httpPort),
				}
			}
			if probe.Liveness_checkType == "TCP" {
				probeHandler.TCPSocket = &corev1.TCPSocketAction{
					Port: intstr.FromInt32(probe.Liveness_tcpPort),
				}
			}
			if probe.Liveness_checkType == "CMD" {
				probeHandler.Exec = &corev1.ExecAction{
					Command: strings.Split(probe.Liveness_cmd, " "),
				}
			}

			livenessProbe := &corev1.Probe{
				InitialDelaySeconds: int32(probe.Liveness_initialDelaySeconds),
				PeriodSeconds:       int32(probe.Liveness_periodSeconds),
				TimeoutSeconds:      int32(probe.Liveness_timeoutSeconds),
				SuccessThreshold:    int32(probe.Liveness_successThreshold), //必须为1
				FailureThreshold:    int32(probe.Liveness_failureThreshold),
				ProbeHandler:        *probeHandler,
			}
			deployInstance.Spec.Template.Spec.Containers[probe.ContainerId].LivenessProbe = livenessProbe
		}
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployProbe Update deployment error: %s\n", err)
			return emptyProbe, err
		}
		return emptyProbe, nil
	}

	lt := len(deployInstance.Spec.Template.Spec.Containers)
	var ccc = make([]ProbeST, 0)
	vv := deployInstance.Spec.Template.Spec.Containers
	// GET
	for i := 0; i < lt; i++ {
		probe := ProbeST{}
		probe.ContainerId = i
		probe.ContainerName = vv[i].Name
		if vv[i].ReadinessProbe != nil {
			var readiness_checkType, readiness_path, readiness_cmd string
			var readiness_httpPort, readiness_tcpPort int32
			if vv[i].ReadinessProbe.HTTPGet != nil {
				readiness_checkType = "HTTP"
				readiness_path = vv[i].ReadinessProbe.HTTPGet.Path
				readiness_httpPort = vv[i].ReadinessProbe.HTTPGet.Port.IntVal
			}
			if vv[i].ReadinessProbe.TCPSocket != nil {
				readiness_checkType = "TCP"
				readiness_tcpPort = vv[i].ReadinessProbe.TCPSocket.Port.IntVal
			}
			if vv[i].ReadinessProbe.Exec != nil {
				readiness_checkType = "CMD"
				readiness_cmd = strings.Join(vv[i].ReadinessProbe.Exec.Command, " ")
			}
			probe.Readiness_checkType = readiness_checkType
			probe.Readiness_path = readiness_path
			probe.Readiness_httpPort = readiness_httpPort
			probe.Readiness_tcpPort = readiness_tcpPort
			probe.Readiness_cmd = readiness_cmd
			probe.Readiness_initialDelaySeconds = vv[i].ReadinessProbe.InitialDelaySeconds
			probe.Readiness_periodSeconds = vv[i].ReadinessProbe.PeriodSeconds
			probe.Readiness_successThreshold = vv[i].ReadinessProbe.SuccessThreshold
			probe.Readiness_failureThreshold = vv[i].ReadinessProbe.FailureThreshold
			probe.Readiness_timeoutSeconds = vv[i].ReadinessProbe.TimeoutSeconds
		}
		if vv[i].LivenessProbe != nil {
			var liveness_checkType, liveness_path, liveness_cmd string
			var liveness_httpPort, liveness_tcpPort int32
			if vv[i].LivenessProbe.HTTPGet != nil {
				liveness_checkType = "HTTP"
				liveness_path = vv[i].LivenessProbe.HTTPGet.Path
				liveness_httpPort = vv[i].LivenessProbe.HTTPGet.Port.IntVal
			}
			if vv[i].LivenessProbe.TCPSocket != nil {
				liveness_checkType = "TCP"
				liveness_tcpPort = vv[i].LivenessProbe.TCPSocket.Port.IntVal
			}
			if vv[i].LivenessProbe.Exec != nil {
				liveness_checkType = "CMD"
				liveness_cmd = strings.Join(vv[i].LivenessProbe.Exec.Command, " ")
			}
			probe.Liveness_checkType = liveness_checkType
			probe.Liveness_path = liveness_path
			probe.Liveness_httpPort = liveness_httpPort
			probe.Liveness_tcpPort = liveness_tcpPort
			probe.Liveness_cmd = liveness_cmd
			probe.Liveness_initialDelaySeconds = vv[i].LivenessProbe.InitialDelaySeconds
			probe.Liveness_periodSeconds = vv[i].LivenessProbe.PeriodSeconds
			probe.Liveness_successThreshold = vv[i].LivenessProbe.SuccessThreshold
			probe.Liveness_failureThreshold = vv[i].LivenessProbe.FailureThreshold
			probe.Liveness_timeoutSeconds = vv[i].LivenessProbe.TimeoutSeconds
		}
		ccc = append(ccc, probe)
	}
	return ccc, nil
}

func DeployEnv(kubeconfig, namespace, deployName string, method string, env EnvSt) ([]EnvSt, error) {
	var emptyEnvSt = make([]EnvSt, 0)
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployEnv Get deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
		return emptyEnvSt, err
	}

	if method == "POST" {
		for _, hv := range env.Envs {
			if hv.Key == "" || hv.Value == "" {
				return emptyEnvSt, fmt.Errorf("invalid env: Key and Value must not be empty")
			}
		}

		var envsArry []corev1.EnvVar
		for _, vv := range env.Envs {
			envsArry = append(envsArry, corev1.EnvVar{
				Name:  vv.Key,
				Value: vv.Value,
			})
		}

		deployInstance.Spec.Template.Spec.Containers[env.ContainerId].Env = envsArry
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployEnv Update deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
			return emptyEnvSt, err
		}
		return emptyEnvSt, nil
	}

	lt := len(deployInstance.Spec.Template.Spec.Containers)
	var ccc = make([]EnvSt, 0)
	vv := deployInstance.Spec.Template.Spec.Containers

	for i := 0; i < lt; i++ {
		result := EnvSt{}
		result.ContainerId = i
		result.ContainerName = vv[i].Name
		var ev []EnvKv
		if vv[i].Env != nil {
			for _, vv := range vv[i].Env {
				ev = append(ev, EnvKv{
					Key:   vv.Name,
					Value: vv.Value,
				})
			}
		}
		result.Envs = ev
		ccc = append(ccc, result)
	}
	return ccc, nil
}

func DeployLifecycle(kubeconfig, namespace, deployName string, method string, lct LifecycleSt) ([]LifecycleSt, error) {
	var emptyLifecycleSt = make([]LifecycleSt, 0)
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployLifecycle Get deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
		return emptyLifecycleSt, err
	}

	if method == "POST" {
		lifeCycle := &corev1.Lifecycle{}
		if lct.PostStart_type == "HTTP" {
			var schemes corev1.URIScheme
			if lct.PostStart_httpGetScheme == "https" || lct.PostStart_httpGetScheme == "HTTPS" {
				schemes = corev1.URISchemeHTTPS
			} else {
				schemes = corev1.URISchemeHTTP
			}
			lifeCycle.PostStart = &corev1.LifecycleHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   lct.PostStart_httpGetPath,
					Port:   intstr.FromInt32(lct.PostStart_httpGetPort),
					Scheme: schemes,
				},
			}
		}

		if lct.PostStart_type == "CMD" {
			lifeCycle.PostStart = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: strings.Split(lct.PostStart_execCommand, " "),
				},
			}
		}

		if lct.PostStart_type == "TCP" {
			lifeCycle.PostStart = &corev1.LifecycleHandler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt32(lct.PostStart_tcpSocketPort),
				},
			}
		}

		if lct.PreStop_type == "HTTP" {
			var schemes corev1.URIScheme
			if lct.PreStop_httpGetScheme == "https" || lct.PreStop_httpGetScheme == "HTTPS" {
				schemes = corev1.URISchemeHTTPS
			} else {
				schemes = corev1.URISchemeHTTP
			}
			lifeCycle.PreStop = &corev1.LifecycleHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   lct.PreStop_httpGetPath,
					Port:   intstr.FromInt32(lct.PreStop_httpGetPort),
					Scheme: schemes,
				},
			}
		}

		if lct.PreStop_type == "CMD" {
			lifeCycle.PreStop = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: strings.Split(lct.PreStop_execCommand, " "),
				},
			}
		}

		if lct.PreStop_type == "TCP" {
			lifeCycle.PreStop = &corev1.LifecycleHandler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.FromInt32(lct.PreStop_tcpSocketPort),
				},
			}
		}

		deployInstance.Spec.Template.Spec.Containers[lct.ContainerId].Lifecycle = lifeCycle
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployLifecycle Update deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
			return emptyLifecycleSt, err
		}
		return emptyLifecycleSt, nil
	}

	lt := len(deployInstance.Spec.Template.Spec.Containers)
	var ccc = make([]LifecycleSt, 0)
	vv := deployInstance.Spec.Template.Spec.Containers
	for i := 0; i < lt; i++ {
		lcts := LifecycleSt{}
		lcts.ContainerId = i
		lcts.ContainerName = vv[i].Name

		// 检查Lifecycle是否为nil
		if vv[i].Lifecycle == nil {
			ccc = append(ccc, lcts)
			continue
		}

		if vv[i].Lifecycle.PostStart != nil {
			if vv[i].Lifecycle.PostStart.HTTPGet != nil {
				lcts.PostStart_type = "HTTP"
				lcts.PostStart_httpGetPath = vv[i].Lifecycle.PostStart.HTTPGet.Path
				lcts.PostStart_httpGetPort = vv[i].Lifecycle.PostStart.HTTPGet.Port.IntVal
				if vv[i].Lifecycle.PostStart.HTTPGet.Scheme != "" {
					lcts.PostStart_httpGetScheme = string(vv[i].Lifecycle.PostStart.HTTPGet.Scheme)
				}
			}
			if vv[i].Lifecycle.PostStart.TCPSocket != nil {
				lcts.PostStart_type = "TCP"
				lcts.PostStart_tcpSocketPort = vv[i].Lifecycle.PostStart.TCPSocket.Port.IntVal
			}
			if vv[i].Lifecycle.PostStart.Exec != nil {
				lcts.PostStart_type = "CMD"
				lcts.PostStart_execCommand = strings.Join(vv[i].Lifecycle.PostStart.Exec.Command, " ")
			}
		}
		if vv[i].Lifecycle.PreStop != nil {
			if vv[i].Lifecycle.PreStop.HTTPGet != nil {
				lcts.PreStop_type = "HTTP"
				lcts.PreStop_httpGetPath = vv[i].Lifecycle.PreStop.HTTPGet.Path
				lcts.PreStop_httpGetPort = vv[i].Lifecycle.PreStop.HTTPGet.Port.IntVal
				if vv[i].Lifecycle.PreStop.HTTPGet.Scheme != "" {
					lcts.PreStop_httpGetScheme = string(vv[i].Lifecycle.PreStop.HTTPGet.Scheme)
				}
			}
			if vv[i].Lifecycle.PreStop.TCPSocket != nil {
				lcts.PreStop_type = "TCP"
				lcts.PreStop_tcpSocketPort = vv[i].Lifecycle.PreStop.TCPSocket.Port.IntVal
			}
			if vv[i].Lifecycle.PreStop.Exec != nil {
				lcts.PreStop_type = "CMD"
				lcts.PreStop_execCommand = strings.Join(vv[i].Lifecycle.PreStop.Exec.Command, " ")
			}
		}
		ccc = append(ccc, lcts)
	}
	return ccc, nil
}

// DeployTolerations gets or updates the tolerations of a deployment
// DeployTolerations获取或更新Deployment的容忍配置
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the deployment | namespace: Deployment的命名空间
//   - deployName: the name of the deployment | deployName: Deployment的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - tolerations: the list of tolerations to set (for POST method) | tolerations: 要设置的容忍配置列表(用于POST方法)
//
// Returns:
//   - []TolerationSt: the list of current tolerations | []TolerationSt: 当前容忍配置列表
//   - error: any error that occurred | error: 发生的任何错误
func DeployTolerations(kubeconfig, namespace, deployName string, method string, tolerations []TolerationSt) ([]TolerationSt, error) {
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployTolerations Get deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
		return tolerations, err
	}

	// 处理POST请求，更新Deployment的Tolerations配置
	// Handle POST request to update Deployment's Tolerations configuration
	if method == "POST" {
		// 验证输入参数
		// Validate input parameters
		for _, tv := range tolerations {
			if tv.Effect != "" {
				if tv.Effect != "NoSchedule" && tv.Effect != "NoExecute" && tv.Effect != "PreferNoSchedule" {
					return tolerations, fmt.Errorf("invalid toleration effect: must be one of NoSchedule, NoExecute, PreferNoSchedule")
				}
			}

			if tv.Operator != "" {
				if tv.Operator != "Equal" && tv.Operator != "Exists" {
					return tolerations, fmt.Errorf("invalid toleration operator: must be Equal or Exists")
				}
			}
		}

		emptyTolerations := make([]corev1.Toleration, 0, len(tolerations))
		for _, vv := range tolerations {
			toleration := corev1.Toleration{
				Key:      vv.Key,
				Operator: corev1.TolerationOperator(vv.Operator),
				Value:    vv.Value,
				Effect:   corev1.TaintEffect(vv.Effect),
			}

			//if vv.TolerationSeconds != nil {
			if *vv.TolerationSeconds > 0 {
				toleration.TolerationSeconds = vv.TolerationSeconds
			}

			emptyTolerations = append(emptyTolerations, toleration)
		}

		deployInstance.Spec.Template.Spec.Tolerations = emptyTolerations
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployTolerations Update deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
			return tolerations, err
		}
		return tolerations, nil
	}

	// GET 逻辑
	// GET logic
	if deployInstance.Spec.Template.Spec.Tolerations != nil {
		result := make([]TolerationSt, 0, len(deployInstance.Spec.Template.Spec.Tolerations))
		for _, vv := range deployInstance.Spec.Template.Spec.Tolerations {
			toleration := TolerationSt{
				Key:      vv.Key,
				Operator: string(vv.Operator),
				Value:    vv.Value,
				Effect:   string(vv.Effect),
			}

			if vv.TolerationSeconds != nil {
				seconds := *vv.TolerationSeconds
				toleration.TolerationSeconds = &seconds
			}

			result = append(result, toleration)
		}
		return result, nil
	}

	// 返回空切片而不是原tolerations参数，保持一致性
	// Return empty slice instead of the original tolerations parameter for consistency
	return []TolerationSt{}, nil
}

// https://kubernetes.io/zh-cn/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
func DeployNodeAffinity(kubeconfig, namespace, deployName string, method string, affinitySt DeployNodeAffinitySt) (DeployNodeAffinitySt, error) {
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployNodeAffinity Get deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
		return affinitySt, err
	}

	if method == "POST" {
		//log.Println(affinitySt)
		// 设置 NodeName（如果提供了节点名称）
		if affinitySt.NodeNames != "" {
			deployInstance.Spec.Template.Spec.NodeName = affinitySt.NodeNames
		} else {
			deployInstance.Spec.Template.Spec.NodeName = ""
		}

		// 设置 NodeSelector
		nodeSelector := make(map[string]string)
		for _, selector := range affinitySt.NodeSelector {
			nodeSelector[selector.Key] = selector.Value
		}
		deployInstance.Spec.Template.Spec.NodeSelector = nodeSelector

		// 构建 RequiredDuringSchedulingIgnoredDuringExecution
		var requiredNodeSelectorTerms []corev1.NodeSelectorTerm
		if len(affinitySt.RequiredAffinity) > 0 {
			for _, matchExpression := range affinitySt.RequiredAffinity {
				var nodeSelectorRequirements []corev1.NodeSelectorRequirement
				for _, expression := range matchExpression.Expressions {
					if expression.Operator == "Exists" || expression.Operator == "DoesNotExist" {
						nodeSelectorRequirements = append(nodeSelectorRequirements, corev1.NodeSelectorRequirement{
							Key:      expression.Key,
							Operator: corev1.NodeSelectorOperator(expression.Operator),
						})
					} else {
						nodeSelectorRequirements = append(nodeSelectorRequirements, corev1.NodeSelectorRequirement{
							Key:      expression.Key,
							Operator: corev1.NodeSelectorOperator(expression.Operator),
							Values:   expression.Value,
						})
					}
				}
				requiredNodeSelectorTerms = append(requiredNodeSelectorTerms, corev1.NodeSelectorTerm{
					MatchExpressions: nodeSelectorRequirements,
				})
			}
		}

		// 构建 PreferredDuringSchedulingIgnoredDuringExecution
		var preferredSchedulingTerm []corev1.PreferredSchedulingTerm
		if len(affinitySt.PreferredAffinity) > 0 {
			for _, preference := range affinitySt.PreferredAffinity {
				var nodeSelectorRequirements []corev1.NodeSelectorRequirement
				for _, expression := range preference.Expressions {
					//log.Println(expression.Operator)
					if expression.Operator == "Exists" || expression.Operator == "DoesNotExist" {
						nodeSelectorRequirements = append(nodeSelectorRequirements, corev1.NodeSelectorRequirement{
							Key:      expression.Key,
							Operator: corev1.NodeSelectorOperator(expression.Operator),
						})
					} else {
						nodeSelectorRequirements = append(nodeSelectorRequirements, corev1.NodeSelectorRequirement{
							Key:      expression.Key,
							Operator: corev1.NodeSelectorOperator(expression.Operator),
							Values:   expression.Value,
						})
					}
				}
				preferredSchedulingTerm = append(preferredSchedulingTerm, corev1.PreferredSchedulingTerm{
					Weight:     preference.Weight,
					Preference: corev1.NodeSelectorTerm{MatchExpressions: nodeSelectorRequirements},
				})
			}
		}

		// 设置亲和性配置
		if deployInstance.Spec.Template.Spec.Affinity == nil {
			deployInstance.Spec.Template.Spec.Affinity = &corev1.Affinity{}
		}

		if deployInstance.Spec.Template.Spec.Affinity.NodeAffinity == nil {
			deployInstance.Spec.Template.Spec.Affinity.NodeAffinity = &corev1.NodeAffinity{}
		}

		// 根据是否有required规则设置RequiredDuringSchedulingIgnoredDuringExecution
		if len(requiredNodeSelectorTerms) > 0 {
			deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{
				NodeSelectorTerms: requiredNodeSelectorTerms,
			}
		} else {
			deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = nil
		}

		// 根据是否有preferred规则设置PreferredDuringSchedulingIgnoredDuringExecution
		if len(preferredSchedulingTerm) > 0 {
			deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = preferredSchedulingTerm
		} else {
			deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = nil
		}

		// 更新Deployment
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployNodeAffinity Update deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
			return affinitySt, err
		}

		// 返回更新后的亲和性配置
		return affinitySt, nil
	}

	// GET logic
	var nodeName string
	if deployInstance.Spec.Template.Spec.NodeName != "" {
		nodeName = deployInstance.Spec.Template.Spec.NodeName
	}

	var labelsArry = make([]LabelsKv, 0)
	if deployInstance.Spec.Template.Spec.NodeSelector != nil {
		for kk, vv := range deployInstance.Spec.Template.Spec.NodeSelector {
			labelsArry = append(labelsArry, LabelsKv{
				Key:   kk,
				Value: vv,
			})
		}
	}

	var matchExpressions = make([]MatchExpression, 0)
	if deployInstance.Spec.Template.Spec.Affinity != nil &&
		deployInstance.Spec.Template.Spec.Affinity.NodeAffinity != nil &&
		deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		for _, vv := range deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
			var expressions = make([]Expression, 0)
			for _, v1 := range vv.MatchExpressions {
				xValues := v1.Values
				if string(v1.Operator) == "Exists" || string(v1.Operator) == "DosNotExist" {
					xValues = []string{}
				}
				expressions = append(expressions, Expression{
					Key:      v1.Key,
					Operator: string(v1.Operator),
					Value:    xValues,
				})
			}
			matchExpressions = append(matchExpressions, MatchExpression{
				Expressions: expressions,
			})
		}
	}

	var preferences = make([]Preference, 0)
	if deployInstance.Spec.Template.Spec.Affinity != nil &&
		deployInstance.Spec.Template.Spec.Affinity.NodeAffinity != nil &&
		len(deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution) > 0 {
		for _, vv := range deployInstance.Spec.Template.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
			var expressions = make([]Expression, 0)
			for _, v1 := range vv.Preference.MatchExpressions {
				xValues := v1.Values
				if string(v1.Operator) == "Exists" || string(v1.Operator) == "DosNotExist" {
					xValues = []string{}
				}
				expressions = append(expressions, Expression{
					Key:      v1.Key,
					Operator: string(v1.Operator),
					Value:    xValues,
				})
			}
			preferences = append(preferences, Preference{
				Weight:      vv.Weight,
				Expressions: expressions,
			})
		}
	}

	// 返回当前的亲和性配置
	return DeployNodeAffinitySt{
		NodeNames:         nodeName,
		NodeSelector:      labelsArry,
		RequiredAffinity:  matchExpressions,
		PreferredAffinity: preferences,
	}, nil
}

// Pod affinity/anti-affinity
type PodAffinityReq struct {
	LabelSelector  []Expression `json:"labelSelector"`
	TopologyKey    string       `json:"topologyKey"`
	NameSpaces     []string     `json:"nameSpaces"`
	MatchLableKeys []string     `json:"matchLableKeys"`
}

type PodAffinityTerm struct {
	Weight         int32        `json:"weight"`
	LabelSelector  []Expression `json:"labelSelector"`
	TopologyKey    string       `json:"topologyKey"`
	NameSpaces     []string     `json:"nameSpaces"`
	MatchLableKeys []string     `json:"matchLableKeys"`
}

type PodAffinitySt struct {
	Required  []PodAffinityReq  `json:"required"`
	Preferred []PodAffinityTerm `json:"preferred"`
}

type PodAntiAffinitySt struct {
	Required  []PodAffinityReq  `json:"required"`
	Preferred []PodAffinityTerm `json:"preferred"`
}

type DeployPodAffinitySt struct {
	PodAffinity     PodAffinitySt     `json:"podAffinity"`
	PodAntiAffinity PodAntiAffinitySt `json:"podAntiAffinity"`
}

func DeployPodAffinity(kubeconfig, namespace, deployName string, method string, affinitySt DeployPodAffinitySt) (DeployPodAffinitySt, error) {
	deployClient := common.ClientSet(kubeconfig).AppsV1().Deployments(namespace)
	deployInstance, err := deployClient.Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployPodAffinity Get deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
		return affinitySt, err
	}

	if method == "POST" {
		// 构建 PodAffinity
		var podAffinityTerms []corev1.PodAffinityTerm
		for _, term := range affinitySt.PodAffinity.Required {
			labelSelector := &metav1.LabelSelector{}

			var matchExpressions []metav1.LabelSelectorRequirement
			var matchLabels map[string]string

			// 处理 LabelSelector
			for _, selector := range term.LabelSelector {
				if selector.Operator == "In" || selector.Operator == "NotIn" ||
					selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
					// 使用 MatchExpressions
					if selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
						matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
							Key:      selector.Key,
							Operator: metav1.LabelSelectorOperator(selector.Operator),
						})
					} else {
						matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
							Key:      selector.Key,
							Operator: metav1.LabelSelectorOperator(selector.Operator),
							Values:   selector.Value,
						})
					}
				} else {
					// 使用 MatchLabels
					if matchLabels == nil {
						matchLabels = make(map[string]string)
					}
					if len(selector.Value) > 0 {
						matchLabels[selector.Key] = selector.Value[0]
					}
				}
			}

			// 构建 LabelSelector
			if len(matchExpressions) > 0 {
				labelSelector.MatchExpressions = matchExpressions
			}
			if len(matchLabels) > 0 {
				if labelSelector.MatchLabels == nil {
					labelSelector.MatchLabels = matchLabels
				} else {
					for k, v := range matchLabels {
						labelSelector.MatchLabels[k] = v
					}
				}
			}

			// 构建 PodAffinityTerm
			podAffinityTerm := corev1.PodAffinityTerm{
				LabelSelector: labelSelector,
				TopologyKey:   term.TopologyKey,
			}

			// 如果有 NameSpaces，添加到 term
			if len(term.NameSpaces) > 0 {
				podAffinityTerm.Namespaces = term.NameSpaces
			}

			// 如果有 MatchLabelKeys，添加到 term
			if len(term.MatchLableKeys) > 0 {
				podAffinityTerm.MatchLabelKeys = term.MatchLableKeys
			}

			podAffinityTerms = append(podAffinityTerms, podAffinityTerm)
		}

		var weightedPodAffinityTerms []corev1.WeightedPodAffinityTerm
		for _, pref := range affinitySt.PodAffinity.Preferred {
			if len(pref.LabelSelector) > 0 {
				labelSelector := &metav1.LabelSelector{}

				var matchExpressions []metav1.LabelSelectorRequirement
				var matchLabels map[string]string

				// 处理 LabelSelector
				for _, selector := range pref.LabelSelector {
					if selector.Operator == "In" || selector.Operator == "NotIn" ||
						selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
						// 使用 MatchExpressions
						if selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
							matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
								Key:      selector.Key,
								Operator: metav1.LabelSelectorOperator(selector.Operator),
							})
						} else {
							matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
								Key:      selector.Key,
								Operator: metav1.LabelSelectorOperator(selector.Operator),
								Values:   selector.Value,
							})
						}
					} else {
						// 使用 MatchLabels
						if matchLabels == nil {
							matchLabels = make(map[string]string)
						}
						if len(selector.Value) > 0 {
							matchLabels[selector.Key] = selector.Value[0]
						}
					}
				}

				// 构建 LabelSelector
				if len(matchExpressions) > 0 {
					labelSelector.MatchExpressions = matchExpressions
				}
				if len(matchLabels) > 0 {
					if labelSelector.MatchLabels == nil {
						labelSelector.MatchLabels = matchLabels
					} else {
						for k, v := range matchLabels {
							labelSelector.MatchLabels[k] = v
						}
					}
				}

				affinityTerm := corev1.PodAffinityTerm{
					LabelSelector: labelSelector,
					TopologyKey:   pref.TopologyKey,
				}

				// 如果有 NameSpaces，添加到 term
				if len(pref.NameSpaces) > 0 {
					affinityTerm.Namespaces = pref.NameSpaces
				}

				// 如果有 MatchLabelKeys，添加到 term
				if len(pref.MatchLableKeys) > 0 {
					affinityTerm.MatchLabelKeys = pref.MatchLableKeys
				}

				weightedTerm := corev1.WeightedPodAffinityTerm{
					Weight:          pref.Weight,
					PodAffinityTerm: affinityTerm,
				}

				weightedPodAffinityTerms = append(weightedPodAffinityTerms, weightedTerm)
			}
		}

		// 构建 PodAntiAffinity
		var podAntiAffinityTerms []corev1.PodAffinityTerm
		for _, term := range affinitySt.PodAntiAffinity.Required {
			labelSelector := &metav1.LabelSelector{}

			var matchExpressions []metav1.LabelSelectorRequirement
			var matchLabels map[string]string

			// 处理 LabelSelector
			for _, selector := range term.LabelSelector {
				if selector.Operator == "In" || selector.Operator == "NotIn" ||
					selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
					// 使用 MatchExpressions
					if selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
						matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
							Key:      selector.Key,
							Operator: metav1.LabelSelectorOperator(selector.Operator),
						})
					} else {
						matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
							Key:      selector.Key,
							Operator: metav1.LabelSelectorOperator(selector.Operator),
							Values:   selector.Value,
						})
					}
				} else {
					// 使用 MatchLabels
					if matchLabels == nil {
						matchLabels = make(map[string]string)
					}
					if len(selector.Value) > 0 {
						matchLabels[selector.Key] = selector.Value[0]
					}
				}
			}

			// 构建 LabelSelector
			if len(matchExpressions) > 0 {
				labelSelector.MatchExpressions = matchExpressions
			}
			if len(matchLabels) > 0 {
				if labelSelector.MatchLabels == nil {
					labelSelector.MatchLabels = matchLabels
				} else {
					for k, v := range matchLabels {
						labelSelector.MatchLabels[k] = v
					}
				}
			}

			// 构建 PodAffinityTerm
			podAntiAffinityTerm := corev1.PodAffinityTerm{
				LabelSelector: labelSelector,
				TopologyKey:   term.TopologyKey,
			}

			// 如果有 NameSpaces，添加到 term
			if len(term.NameSpaces) > 0 {
				podAntiAffinityTerm.Namespaces = term.NameSpaces
			}

			// 如果有 MatchLabelKeys，添加到 term
			if len(term.MatchLableKeys) > 0 {
				podAntiAffinityTerm.MatchLabelKeys = term.MatchLableKeys
			}

			podAntiAffinityTerms = append(podAntiAffinityTerms, podAntiAffinityTerm)
		}

		var weightedPodAntiAffinityTerms []corev1.WeightedPodAffinityTerm
		for _, pref := range affinitySt.PodAntiAffinity.Preferred {
			if len(pref.LabelSelector) > 0 {
				labelSelector := &metav1.LabelSelector{}

				var matchExpressions []metav1.LabelSelectorRequirement
				var matchLabels map[string]string

				// 处理 LabelSelector
				for _, selector := range pref.LabelSelector {
					if selector.Operator == "In" || selector.Operator == "NotIn" ||
						selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
						// 使用 MatchExpressions
						if selector.Operator == "Exists" || selector.Operator == "DoesNotExist" {
							matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
								Key:      selector.Key,
								Operator: metav1.LabelSelectorOperator(selector.Operator),
							})
						} else {
							matchExpressions = append(matchExpressions, metav1.LabelSelectorRequirement{
								Key:      selector.Key,
								Operator: metav1.LabelSelectorOperator(selector.Operator),
								Values:   selector.Value,
							})
						}
					} else {
						// 使用 MatchLabels
						if matchLabels == nil {
							matchLabels = make(map[string]string)
						}
						if len(selector.Value) > 0 {
							matchLabels[selector.Key] = selector.Value[0]
						}
					}
				}

				// 构建 LabelSelector
				if len(matchExpressions) > 0 {
					labelSelector.MatchExpressions = matchExpressions
				}
				if len(matchLabels) > 0 {
					if labelSelector.MatchLabels == nil {
						labelSelector.MatchLabels = matchLabels
					} else {
						for k, v := range matchLabels {
							labelSelector.MatchLabels[k] = v
						}
					}
				}

				affinityTerm := corev1.PodAffinityTerm{
					LabelSelector: labelSelector,
					TopologyKey:   pref.TopologyKey,
				}

				// 如果有 NameSpaces，添加到 term
				if len(pref.NameSpaces) > 0 {
					affinityTerm.Namespaces = pref.NameSpaces
				}

				// 如果有 MatchLabelKeys，添加到 term
				if len(pref.MatchLableKeys) > 0 {
					affinityTerm.MatchLabelKeys = pref.MatchLableKeys
				}

				weightedTerm := corev1.WeightedPodAffinityTerm{
					Weight:          pref.Weight,
					PodAffinityTerm: affinityTerm,
				}

				weightedPodAntiAffinityTerms = append(weightedPodAntiAffinityTerms, weightedTerm)
			}
		}

		// 设置亲和性配置
		if deployInstance.Spec.Template.Spec.Affinity == nil {
			deployInstance.Spec.Template.Spec.Affinity = &corev1.Affinity{}
		}

		// 设置 PodAffinity
		if len(podAffinityTerms) > 0 || len(weightedPodAffinityTerms) > 0 {
			deployInstance.Spec.Template.Spec.Affinity.PodAffinity = &corev1.PodAffinity{}

			if len(podAffinityTerms) > 0 {
				deployInstance.Spec.Template.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = podAffinityTerms
			}

			if len(weightedPodAffinityTerms) > 0 {
				deployInstance.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = weightedPodAffinityTerms
			}
		} else {
			deployInstance.Spec.Template.Spec.Affinity.PodAffinity = nil
		}

		// 设置 PodAntiAffinity
		if len(podAntiAffinityTerms) > 0 || len(weightedPodAntiAffinityTerms) > 0 {
			deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity = &corev1.PodAntiAffinity{}

			if len(podAntiAffinityTerms) > 0 {
				deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = podAntiAffinityTerms
			}

			if len(weightedPodAntiAffinityTerms) > 0 {
				deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = weightedPodAntiAffinityTerms
			}
		} else {
			deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity = nil
		}

		// 更新Deployment
		_, err = deployClient.Update(context.TODO(), deployInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployPodAffinity Update deployment error: %s, namespace: %s, deployName: %s\n", err, namespace, deployName)
			return affinitySt, err
		}

		// 返回更新后的亲和性配置
		return affinitySt, nil
	}

	// GET logic
	var podAffinityRequired []PodAffinityReq
	var podAffinityPreferred []PodAffinityTerm
	var podAntiAffinityRequired []PodAffinityReq
	var podAntiAffinityPreferred []PodAffinityTerm

	if deployInstance.Spec.Template.Spec.Affinity != nil {
		// 获取 PodAffinity
		if deployInstance.Spec.Template.Spec.Affinity.PodAffinity != nil {
			// Required
			if deployInstance.Spec.Template.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
				for _, term := range deployInstance.Spec.Template.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
					var labelSelectors []Expression

					if term.LabelSelector != nil && term.LabelSelector.MatchExpressions != nil {
						for _, v := range term.LabelSelector.MatchExpressions {
							labelSelectors = append(labelSelectors, Expression{
								Key:      v.Key,
								Operator: string(v.Operator),
								Value:    v.Values,
							})
						}
					}

					if term.LabelSelector != nil && term.LabelSelector.MatchLabels != nil {
						for kk, vv := range term.LabelSelector.MatchLabels {
							labelSelectors = append(labelSelectors, Expression{
								Key:      kk,
								Operator: "",
								Value:    []string{vv},
							})
						}
					}

					podAffinityRequired = append(podAffinityRequired, PodAffinityReq{
						LabelSelector:  labelSelectors,
						TopologyKey:    term.TopologyKey,
						NameSpaces:     term.Namespaces,
						MatchLableKeys: term.MatchLabelKeys,
					})
				}
			}

			// Preferred
			if deployInstance.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
				for _, term := range deployInstance.Spec.Template.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
					var labelSelectors []Expression
					if term.PodAffinityTerm.LabelSelector != nil && term.PodAffinityTerm.LabelSelector.MatchExpressions != nil {
						for _, v := range term.PodAffinityTerm.LabelSelector.MatchExpressions {
							labelSelectors = append(labelSelectors, Expression{
								Key:      v.Key,
								Operator: string(v.Operator),
								Value:    v.Values,
							})
						}
					}
					if term.PodAffinityTerm.LabelSelector != nil && term.PodAffinityTerm.LabelSelector.MatchLabels != nil {
						for kk, vv := range term.PodAffinityTerm.LabelSelector.MatchLabels {
							labelSelectors = append(labelSelectors, Expression{
								Key:      kk,
								Operator: "",
								Value:    []string{vv},
							})
						}
					}
					podAffinityPreferred = append(podAffinityPreferred, PodAffinityTerm{
						Weight:         term.Weight,
						LabelSelector:  labelSelectors,
						TopologyKey:    term.PodAffinityTerm.TopologyKey,
						NameSpaces:     term.PodAffinityTerm.Namespaces,
						MatchLableKeys: term.PodAffinityTerm.MatchLabelKeys,
					})
				}
			}
		}

		// 获取 PodAntiAffinity
		if deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity != nil {
			// Required
			if deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
				for _, term := range deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
					var labelSelectors []Expression
					if term.LabelSelector != nil && term.LabelSelector.MatchExpressions != nil {
						for _, v := range term.LabelSelector.MatchExpressions {
							labelSelectors = append(labelSelectors, Expression{
								Key:      v.Key,
								Operator: string(v.Operator),
								Value:    v.Values,
							})
						}
					}

					if term.LabelSelector != nil && term.LabelSelector.MatchLabels != nil {
						for kk, vv := range term.LabelSelector.MatchLabels {
							labelSelectors = append(labelSelectors, Expression{
								Key:      kk,
								Operator: "",
								Value:    []string{vv},
							})
						}
					}

					podAntiAffinityRequired = append(podAntiAffinityRequired, PodAffinityReq{
						LabelSelector:  labelSelectors,
						TopologyKey:    term.TopologyKey,
						NameSpaces:     term.Namespaces,
						MatchLableKeys: term.MatchLabelKeys,
					})
				}
			}

			// Preferred
			if deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
				for _, term := range deployInstance.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
					var labelSelectors []Expression
					if term.PodAffinityTerm.LabelSelector != nil && term.PodAffinityTerm.LabelSelector.MatchExpressions != nil {
						for _, v := range term.PodAffinityTerm.LabelSelector.MatchExpressions {
							labelSelectors = append(labelSelectors, Expression{
								Key:      v.Key,
								Operator: string(v.Operator),
								Value:    v.Values,
							})
						}
					}

					if term.PodAffinityTerm.LabelSelector != nil && term.PodAffinityTerm.LabelSelector.MatchLabels != nil {
						for kk, vv := range term.PodAffinityTerm.LabelSelector.MatchLabels {
							labelSelectors = append(labelSelectors, Expression{
								Key:      kk,
								Operator: "",
								Value:    []string{vv},
							})
						}
					}

					podAntiAffinityPreferred = append(podAntiAffinityPreferred, PodAffinityTerm{
						Weight:         term.Weight,
						LabelSelector:  labelSelectors,
						TopologyKey:    term.PodAffinityTerm.TopologyKey,
						NameSpaces:     term.PodAffinityTerm.Namespaces,
						MatchLableKeys: term.PodAffinityTerm.MatchLabelKeys,
					})
				}
			}
		}
	}

	// 返回当前的亲和性配置
	return DeployPodAffinitySt{
		PodAffinity: PodAffinitySt{
			Required:  podAffinityRequired,
			Preferred: podAffinityPreferred,
		},
		PodAntiAffinity: PodAntiAffinitySt{
			Required:  podAntiAffinityRequired,
			Preferred: podAntiAffinityPreferred,
		},
	}, nil
}
