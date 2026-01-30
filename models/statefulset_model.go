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

	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/yaml"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/util/retry"
)

// Statefulset - 状态集模型结构体(Statefulset - StatefulSet model struct)
type Statefulset struct {
	StatefulsetName string `json:"statefulsetName" form:"statefulsetName"` // Statefulset名称(Statefulset name)
	NameSpace       string `json:"nameSpace" form:"nameSpace"`             // 命名空间(Namespace)
	//CreationTimestamp    time.Time `json:"created_time"`
	RevisionHistoryLimit int32  `json:"revisionHistoryLimit" form:"historyVersionLimit"` // 历史版本限制(Revision history limit)
	Replicas             int32  `json:"replicas" form:"replicas"`                        // 副本数(Replicas)
	AvailableReplicas    int32  `json:"availableReplicas"`                               // 可用副本数(Available replicas)
	PodNumber            string `json:"podNumber"`                                       // Pod数量(Pod number)
	Labels               string `json:"labels" form:"labels"`                            // 标签(Labels)
	ContainerName        string `form:"containerName"`                                   // 容器名称(Container name)
	ImgUrl               string `json:"imgUrl" form:"imgUrl"`                            // 镜像地址(Image URL)
	ContainerPortName    string `json:"-" form:"containerPortName"`                      // 容器端口名称(Container port name)
	ContainerPort        string `form:"containerPort"`                                   // 容器端口(Container port)
	//HostPort             int32  `form:"hostPort"`
	// Resource             string    `json:"resource"`
	// StrategyType         string    `json:"strategyType"`
	CreateTime string `json:"createTime"` //创建时间(Creation time)
}

// StsDetails - 状态集详情结构体(StsDetails - StatefulSet details struct)
type StsDetails struct {
	StatefulsetName string `json:"statefulsetName"` // Statefulset名称(Statefulset name)
	NameSpace       string `json:"nameSpace"`       // 命名空间(Namespace)
	Strategy        string `json:"strategy"`        // 策略(Strategy)
	//StrategyRollingUpdate string `json:"strategyRollingUpdate"`
	//PodNumber             string `json:"podNumber"`
	Selector string `json:"selector"` // 选择器(Selector)
	//ImageUrl              string `json:"imageUrl"`
	Annotations string          `json:"annotations"` // 注解(Annotations)
	Status      string          `json:"status"`      // 状态(Status)
	Labels      string          `json:"labels"`      // 标签(Labels)
	CreateTime  string          `json:"createTime"`  //创建时间(Creation time)
	Replicasets []StsReplicaset `json:"replicasets"` // 副本集列表(Replica set list)
}

// StsReplicaset - 状态集副本集结构体(StsReplicaset - StatefulSet replica set struct)
type StsReplicaset struct {
	ReplicasetName string `json:"replicasetName"` // 副本集名称(Replica set name)
	ImageUrl       string `json:"imageUrl"`       // 镜像地址(Image URL)
	CreateTime     string `json:"createTime"`     // 创建时间(Creation time)
}

// StatefulsetList - 获取状态集列表(StatefulsetList - Get statefulset list)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulsetName - 状态集名称(Statefulset name)
// labelsKey - 标签键(Labels key)
// labelsValue - 标签值(Labels value)
// 返回值: 状态集列表和错误信息(Return: statefulset list and error information)
func StatefulsetList(kubeconfig, namespace, statefulsetName, labelsKey, labelsValue string) ([]Statefulset, error) {
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	var bbb = make([]Statefulset, 0)
	xList, err := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("[ERROR] ListStsError err:%v", err)
		return bbb, err
	}

	for _, vv := range xList.Items {
		var labelsStr, imgUrlStr, containerNameStr, ContainerPortNameStr, ContainerPortStr string

		//搜索(Search)
		if statefulsetName != "" {
			if !strings.Contains(vv.Name, statefulsetName) {
				continue
			}
		}

		for k1, v1 := range vv.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		for _, v2 := range vv.Spec.Template.Spec.Containers {
			containerNameStr += fmt.Sprintf("%s,", v2.Image)
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
		xItems := &Statefulset{
			StatefulsetName:      vv.Name,                                                          // 状态集名称(Statefulset name)
			NameSpace:            vv.Namespace,                                                     // 命名空间(Namespace)
			RevisionHistoryLimit: *vv.Spec.RevisionHistoryLimit,                                    // 历史版本限制(Revision history limit)
			Replicas:             *vv.Spec.Replicas,                                                // 副本数(Replicas)
			AvailableReplicas:    vv.Status.AvailableReplicas,                                      // 可用副本数(Available replicas)
			PodNumber:            fmt.Sprintf("%d/%d", vv.Status.ReadyReplicas, *vv.Spec.Replicas), // Pod数量(Pod number)
			Labels:               labelsStr,                                                        // 标签(Labels)

			ContainerName: containerNameStr, // 容器名称(Container name)
			ImgUrl:        imgUrlStr,        // 镜像地址(Image URL)

			ContainerPortName: ContainerPortNameStr, // 容器端口名称(Container port name)
			ContainerPort:     ContainerPortStr,     // 容器端口(Container port)
			//HostPort:          vv.Spec.Template.Spec.Containers[0].Ports[0].HostPort,

			CreateTime: vv.CreationTimestamp.Format("2006-01-02 15:04:05"), // 创建时间(Creation time)
		}
		bbb = append(bbb, *xItems)
	}
	return bbb, err
}

// StatefulsetDetail - 获取状态集详细信息(StatefulsetDetail - Get statefulset detailed information)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulset - 状态集名称(Statefulset name)
// 返回值: 状态集详情和错误信息(Return: statefulset details and error information)
func StatefulsetDetail(kubeconfig, namespace, statefulset string) (*StsDetails, error) {
	sts, err := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulset, metav1.GetOptions{})
	if err != nil {
		return &StsDetails{}, err
	}

	var selectorStr, labelsStr, annotationsStr string
	for kk, vv := range sts.Spec.Selector.MatchLabels {
		selectorStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(selectorStr) > 0 {
		selectorStr = selectorStr[0 : len(selectorStr)-1]
	}

	for kk, vv := range sts.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	for kk, vv := range sts.ObjectMeta.Annotations {
		if strings.Contains(kk, "last-applied-configuration") {
			continue
		}
		annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	//replicaset list
	selector, err := metav1.LabelSelectorAsSelector(sts.Spec.Selector)
	if err != nil {
		panic(err.Error())
	}
	rsList, err := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		panic(err.Error())
	}
	var repcs = make([]StsReplicaset, 0)
	for _, rs := range rsList.Items {
		xItems := &StsReplicaset{
			ReplicasetName: rs.Name,                                            // 副本集名称(Replica set name)
			ImageUrl:       rs.Spec.Template.Spec.Containers[0].Image,          // 镜像地址(Image URL)
			CreateTime:     rs.CreationTimestamp.Format("2006-01-02 15:04:05"), // 创建时间(Creation time)
		}
		repcs = append(repcs, *xItems)
	}

	return &StsDetails{
		StatefulsetName: sts.Name,                                                                                                                 // 状态集名称(Statefulset name)
		NameSpace:       sts.Namespace,                                                                                                            // 命名空间(Namespace)
		Strategy:        string(sts.Spec.UpdateStrategy.Type),                                                                                     // 策略(Strategy)
		Selector:        selectorStr,                                                                                                              // 选择器(Selector)
		Annotations:     annotationsStr,                                                                                                           // 注解(Annotations)
		Status:          fmt.Sprintf("就绪:%d个,已更新:%d个,可用:%d个", sts.Status.ReadyReplicas, sts.Status.UpdatedReplicas, sts.Status.AvailableReplicas), // 状态(Status)
		Labels:          labelsStr,                                                                                                                // 标签(Labels)
		CreateTime:      sts.CreationTimestamp.Format("2006-01-02 15:04:05"),                                                                      // 创建时间(Creation time)
		Replicasets:     repcs,                                                                                                                    // 副本集列表(Replica set list)
	}, nil
}

// StatefulsetCreate - 创建状态集(StatefulsetCreate - Create statefulset)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// bodys - 请求体数据(Request body data)
// 返回值: 错误信息(Return: error information)
func StatefulsetCreate(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	statefulsetName := gp.Get("statefulsetName").String()
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
		"app": statefulsetName,
	}
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").String()] = vv.Get("value").String()
	}

	//ContainerPortInt, _ := strconv.Atoi(statefulset.ContainerPort)
	stsInstance := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{ //定义deployment资源的apiVersion、kind类型等(Define deployment resource apiVersion, kind type, etc.)
			Kind:       "Statefulset",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{ //定义控制器的metadata：名称、namespace、labels等(Define controller metadata: name, namespace, labels, etc.)
			Name:      statefulsetName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: appsv1.StatefulSetSpec{ //定义deployment控制器的spec字段(Define deployment controller spec field)
			Replicas: int32Ptr(int32(replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labelsMap,
			},
			RevisionHistoryLimit: int32Ptr(10),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelsMap,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:            statefulsetName,
							Image:           imageUrl,
							ImagePullPolicy: pullPolicy,
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									//Name:          "",
									ContainerPort: int32(containerPort),
									//HostPort:      statefulset.HostPort,
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
		stsInstance.Spec.Template.Spec.Containers[0].Resources = *resReq
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
			stsInstance.Spec.Template.Spec.Containers[0].ReadinessProbe = probe
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
				//deployInstance.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Path = readiness_path
				//deployInstance.Spec.Template.Spec.Containers[0].ReadinessProbe.HTTPGet.Port = intstr.FromString(readiness_httpPort)
				httpGet = &corev1.HTTPGetAction{
					Path: liveness_path,
					Port: intstr.FromInt32(liveness_httpPort),
				}
				probeHandler.HTTPGet = httpGet
			}
			if liveness_checkType == "TCP" {
				//deployInstance.Spec.Template.Spec.Containers[0].ReadinessProbe.TCPSocket.Port =
				//probe.TCPSocket.Port = intstr.FromString(readiness_tcpPort)
				tcpSocket = &corev1.TCPSocketAction{
					Port: intstr.FromInt32(liveness_tcpPort),
				}
				probeHandler.TCPSocket = tcpSocket
			}
			if liveness_checkType == "CMD" {
				//deployInstance.Spec.Template.Spec.Containers[0].ReadinessProbe.Exec.Command = cmdStr
				execCmd = &corev1.ExecAction{
					Command: cmdStr,
				}
				probeHandler.Exec = execCmd
			}

			probe := &corev1.Probe{
				InitialDelaySeconds: int32(liveness_initialDelaySeconds),
				PeriodSeconds:       int32(liveness_periodSeconds),
				TimeoutSeconds:      int32(liveness_timeoutSeconds),
				SuccessThreshold:    int32(liveness_successThreshold), //必须为1(Must be 1)
				FailureThreshold:    int32(liveness_failureThreshold),
				ProbeHandler:        *probeHandler,
			}
			stsInstance.Spec.Template.Spec.Containers[0].LivenessProbe = probe
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
		stsInstance.Spec.Template.Spec.Containers[0].Lifecycle = lifeCycle
	}

	if variableCheck == "on" {
		var envsArry []corev1.EnvVar
		for _, vv := range gp.Get("envs").Array() {
			envsArry = append(envsArry, corev1.EnvVar{
				Name:  vv.Get("key").String(),
				Value: vv.Get("value").String(),
			})
		}
		stsInstance.Spec.Template.Spec.Containers[0].Env = envsArry
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
			ServiceName: statefulsetName,
			NameSpace:   nameSpace,
			SvcType:     gp.Get("svcType").String(),
			Ports:       svcPortArry,
			Labels:      labelsMap,
			IsHeadless:  "on",
		}
		data, err := json.Marshal(service)
		if err != nil {
			return err
		}
		_ = SvcCreate(kubeconfig, data)
	}

	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(nameSpace)
	_, err := stsClient.Create(context.TODO(), stsInstance, metav1.CreateOptions{})
	return err
}

// StatefulsetClone - 克隆状态集(StatefulsetClone - Clone statefulset)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 源命名空间(Source namespace)
// objname - 源对象名称(Source object name)
// target_clusterid - 目标集群ID(Target cluster ID)
// target_namespace - 目标命名空间(Target namespace)
// target_objname - 目标对象名称(Target object name)
// 返回值: 错误信息(Return: error information)
func StatefulsetClone(kubeconfig, namespace, objname, target_clusterid, target_namespace, target_objname string) error {
	stsInstance, err := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).Get(context.TODO(), objname, metav1.GetOptions{})
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
		target_objname = objname
	}

	if target_clusterid == "" {
		target_clusterid = kubeconfig
	}

	//在新集群创建namespace(Create namespace in new cluster)
	err3 := CreateNsByExist(target_clusterid, target_namespace)
	if err3 != nil {
		return err3
	}

	stsInstance.Name = target_objname
	stsInstance.Namespace = target_namespace
	stsInstance.ResourceVersion = ""
	stsInstance.ObjectMeta.Labels["app"] = target_objname
	stsInstance.Spec.Selector.MatchLabels["app"] = target_objname
	stsInstance.Spec.Template.Labels["app"] = target_objname
	//new cluster
	NewClient := common.ClientSet(target_clusterid).AppsV1().StatefulSets(target_namespace)
	_, err2 := NewClient.Get(context.TODO(), target_objname, metav1.GetOptions{})
	if errors.IsNotFound(err2) { //没有就创建(Create if not exists)
		_, err := NewClient.Create(context.TODO(), stsInstance, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if err2 == nil {
		_, err = NewClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	return err2
}

// StatefulsetModify - 修改状态集(StatefulsetModify - Modify statefulset)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// newStatefulset - 新的状态集信息(New statefulset information)
// 返回值: 错误信息(Return: error information)
func StatefulsetModify(kubeconfig string, newStatefulset *Statefulset) error {
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(newStatefulset.NameSpace)
	stsInstance, err := stsClient.Get(context.TODO(), newStatefulset.StatefulsetName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ContainerPortInt, _ := strconv.Atoi(newStatefulset.ContainerPort)
	stsInstance.Spec.Replicas = &newStatefulset.Replicas
	stsInstance.Spec.RevisionHistoryLimit = &newStatefulset.RevisionHistoryLimit
	stsInstance.Spec.Template.Spec.Containers[0].Image = newStatefulset.ImgUrl
	stsInstance.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(ContainerPortInt)
	//stsInstance.Spec.Template.Spec.Containers[0].Ports[0].HostPort = newStatefulset.HostPort
	//调用StatefulsetmentClient接口对象中的update方法，来更新deployment资源(Call the update method in the StatefulsetmentClient interface object to update the deployment resource)
	_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil

}

func StatefulsetYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	sts := &appsv1.StatefulSet{}
	err = json.Unmarshal(data, sts)
	if err != nil {
		return err
	}
	namespace := sts.ObjectMeta.Namespace
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.AppsV1().StatefulSets(namespace).Update(context.TODO(), sts, metav1.UpdateOptions{})
	return err
}

// StatefulsetDel - 删除状态集(StatefulsetDel - Delete statefulset)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulset - 状态集名称(Statefulset name)
// 返回值: 错误信息(Return: error information)
func StatefulsetDel(kubeconfig, namespace, statefulset string) error {
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	err := stsClient.Delete(context.TODO(), statefulset, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetStatefulsetYaml - 获取状态集YAML定义(GetStatefulsetYaml - Get statefulset YAML definition)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulset - 状态集名称(Statefulset name)
// 返回值: YAML字符串和错误信息(Return: YAML string and error information)
func GetStatefulsetYaml(kubeconfig, namespace, statefulset string) (string, error) {

	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulset, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	stsUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(stsInstance)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(stsUnstructured)
	if err != nil {
		return "", err
	}
	//fmt.Println(string(yamlBytes))
	return string(yamlBytes), nil

}

// RollBack
// StatefulsetRollBack - 状态集回滚(StatefulsetRollBack - Statefulset rollback)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulsetName - 状态集名称(Statefulset name)
// replicaSet - 副本集名称(Replica set name)
// 返回值: 错误信息(Return: error information)
func StatefulsetRollBack(kubeconfig, namespace, statefulsetName, replicaSet string) error {

	sts, err := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	repset, err1 := common.ClientSet(kubeconfig).AppsV1().ReplicaSets(namespace).Get(context.TODO(), replicaSet, metav1.GetOptions{})
	if err1 != nil {
		return err1
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		sts.Spec.Template = repset.Spec.Template
		_, updateErr := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).Update(context.TODO(), sts, metav1.UpdateOptions{})
		return updateErr
	})
	return retryErr
}

// Restart
// StatefulsetRestart - 重启状态集(StatefulsetRestart - Restart statefulset)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulsetName - 状态集名称(Statefulset name)
// 返回值: 错误信息(Return: error information)
func StatefulsetRestart(kubeconfig, namespace, statefulsetName string) error {
	patchOpt := metav1.PatchOptions{FieldManager: "kubectl-rollout"}
	patchInfo := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
	_, err := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace).Patch(context.TODO(), statefulsetName, types.StrategicMergePatchType, []byte(patchInfo), patchOpt)
	return err
}

// StatefulsetImage - 更新状态集镜像(StatefulsetImage - Update statefulset image)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulsetName - 状态集名称(Statefulset name)
// method - 请求方法(Request method)
// newImages - 新镜像列表(New image list)
// 返回值: 镜像列表和错误信息(Return: image list and error information)
func StatefulsetImage(kubeconfig, namespace, statefulsetName string, method string, newImages []ImageKv) ([]ImageKv, error) {
	var bbb = make([]ImageKv, 0)
	dClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	sts, err := dClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] Statefulset Get error:%s\n", err)
		return bbb, err
	}

	if method == "POST" {
		for k1, v1 := range sts.Spec.Template.Spec.Containers {
			for _, v2 := range newImages {
				log.Printf("[DEUBG] id:%d,name:%s,image:%s\n", k1, v1.Name, v2.Image)
				if v1.Name == v2.ContainerName {
					sts.Spec.Template.Spec.Containers[k1].Image = v2.Image
				}
			}
		}
		_, err = dClient.Update(context.TODO(), sts, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] Statefulset Update image error:%s\n", err)
			return bbb, err
		}
		return bbb, nil
	} else {
		for kk, vv := range sts.Spec.Template.Spec.Containers {
			bbb = append(bbb, ImageKv{
				ContainerName: vv.Name,
				ContainerId:   kk,
				Image:         vv.Image,
			})
		}
		return bbb, nil
	}
}

// StatefulsetReplicas - 更新状态集副本数(StatefulsetReplicas - Update statefulset replicas)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulsetName - 状态集名称(Statefulset name)
// method - 请求方法(Request method)
// podNumber - Pod数量(Pod number)
// 返回值: Pod数量和错误信息(Return: pod number and error information)
func StatefulsetReplicas(kubeconfig, namespace, statefulsetName string, method string, podNumber int32) (int32, error) {
	dClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	sts, err := dClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] Statefulset Get error:%s\n", err)
		return 0, err
	}

	if method == "POST" {
		sts.Spec.Replicas = &podNumber
		_, err = dClient.Update(context.TODO(), sts, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] Statefulset Replicas Update error:%s\n", err)
			return 0, err
		}
		return podNumber, nil
	} else {
		podnum := sts.Spec.Replicas
		return *podnum, nil
	}
}

// StatefulsetLabels - 更新状态集标签(StatefulsetLabels - Update statefulset labels)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// namespace - 命名空间(Namespace)
// statefulsetName - 状态集名称(Statefulset name)
// method - 请求方法(Request method)
// labelsMap - 标签映射(Labels map)
// 返回值: 标签列表和错误信息(Return: labels list and error information)
func StatefulsetLabels(kubeconfig, namespace, statefulsetName string, method string, labelsMap map[string]string) ([]LabelsKv, error) {
	var bbb = make([]LabelsKv, 0)
	dClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	sts, err := dClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] Statefulset Get error:%s\n", err)
		return bbb, err
	}
	if method == "POST" {
		sts.ObjectMeta.Labels = labelsMap
		_, err = dClient.Update(context.TODO(), sts, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetLabels Update error:%s\n", err)
			return bbb, err
		}
		return bbb, nil
	} else {
		if sts.ObjectMeta.Labels != nil {
			for kk, vv := range sts.ObjectMeta.Labels {
				bbb = append(bbb, LabelsKv{
					Key:   kk,
					Value: vv,
				})
			}
		}
		return bbb, nil
	}
}

// StatefulsetHost gets or updates the host aliases of a statefulset
// StatefulsetHost获取或更新StatefulSet的主机别名
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the statefulset | namespace: StatefulSet的命名空间
//   - statefulsetName: the name of the statefulset | statefulsetName: StatefulSet的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - hostAlias: the list of host aliases to set (for POST method) | hostAlias: 要设置的主机别名列表(用于POST方法)
//
// Returns:
//   - []HostKv: the list of current host aliases | []HostKv: 当前主机别名列表
//   - error: any error that occurred | error: 发生的任何错误
func StatefulsetHost(kubeconfig, namespace, statefulsetName string, method string, hostAlias []HostKv) ([]HostKv, error) {
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] StatefulsetHost Get statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
		return hostAlias, err
	}

	// Handle POST request to update StatefulSet's HostAliases configuration
	if method == "POST" {
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

		stsInstance.Spec.Template.Spec.HostAliases = emptyHostAliases
		_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetHost Update statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
			return hostAlias, err
		}
		return hostAlias, nil
	}

	// GET logic
	if stsInstance.Spec.Template.Spec.HostAliases != nil {
		result := make([]HostKv, 0, len(stsInstance.Spec.Template.Spec.HostAliases))
		for _, vv := range stsInstance.Spec.Template.Spec.HostAliases {
			result = append(result, HostKv{
				Ip:     vv.IP,
				Domain: strings.Join(vv.Hostnames, ","),
			})
		}
		return result, nil
	}

	// Return empty slice instead of the original hostAlias parameter for consistency
	return []HostKv{}, nil
}

// StatefulsetEnv gets or updates the environment variables of a statefulset's container
// StatefulsetEnv获取或更新StatefulSet容器的环境变量
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the statefulset | namespace: StatefulSet的命名空间
//   - statefulsetName: the name of the statefulset | statefulsetName: StatefulSet的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - env: the environment variables to set (for POST method) | env: 要设置的环境变量(用于POST方法)
//
// Returns:
//   - []EnvSt: the list of current environment variables for all containers | []EnvSt: 所有容器的当前环境变量列表
//   - error: any error that occurred | error: 发生的任何错误
func StatefulsetEnv(kubeconfig, namespace, statefulsetName string, method string, env EnvSt) ([]EnvSt, error) {
	var emptyEnvSt = make([]EnvSt, 0)
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] StatefulsetEnv Get statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
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

		stsInstance.Spec.Template.Spec.Containers[env.ContainerId].Env = envsArry
		_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetEnv Update statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
			return emptyEnvSt, err
		}
		return emptyEnvSt, nil
	}

	// GET logic
	lt := len(stsInstance.Spec.Template.Spec.Containers)
	var ccc = make([]EnvSt, 0)
	vv := stsInstance.Spec.Template.Spec.Containers

	for i := 0; i < lt; i++ {
		result := EnvSt{}
		result.ContainerId = i
		result.ContainerName = vv[i].Name
		var ev []EnvKv
		if vv[i].Env != nil {
			for _, envVar := range vv[i].Env {
				ev = append(ev, EnvKv{
					Key:   envVar.Name,
					Value: envVar.Value,
				})
			}
		}
		result.Envs = ev
		ccc = append(ccc, result)
	}
	return ccc, nil
}

// StatefulsetResource gets or updates the resource limits/requests of a statefulset's container
// StatefulsetResource获取或更新StatefulSet容器的资源限制/请求
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the statefulset | namespace: StatefulSet的命名空间
//   - statefulsetName: the name of the statefulset | statefulsetName: StatefulSet的名称
//   - containerName: the name of the container | containerName: 容器名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - resl: the resource limits/requests to set (for POST method) | resl: 要设置的资源限制/请求(用于POST方法)
//
// Returns:
//   - ResourceSt: the current resource limits/requests | ResourceSt: 当前资源限制/请求
//   - error: any error that occurred | error: 发生的任何错误
func StatefulsetResource(kubeconfig, namespace, statefulsetName, method string, resl ResourceSt) ([]ResourceSt, error) {
	var emptyResouceSt = make([]ResourceSt, 0)
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] StatefulsetResource Get statefulset error: %s\n", err)
		return emptyResouceSt, err
	}

	// Handle POST request to update StatefulSet's resource limits configuration
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
		stsInstance.Spec.Template.Spec.Containers[resl.ContainerId].Resources = *resReq
		_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetResource Update statefulset error: %s\n", err)
			return emptyResouceSt, err
		}
		return emptyResouceSt, nil
	}

	lt := len(stsInstance.Spec.Template.Spec.Containers)
	var ccc = make([]ResourceSt, 0)
	vv := stsInstance.Spec.Template.Spec.Containers
	for i := 0; i < lt; i++ {
		result := ResourceSt{}
		result.ContainerId = i
		result.ContainerName = vv[i].Name
		// Check if there are resource configurations (check if Requests and Limits are empty)
		if len(vv[i].Resources.Requests) > 0 || len(vv[i].Resources.Limits) > 0 {
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

// StatefulsetLifecycle gets or updates the lifecycle configuration of a statefulset's container
// StatefulsetLifecycle获取或更新StatefulSet容器的生命周期配置
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the statefulset | namespace: StatefulSet的命名空间
//   - statefulsetName: the name of the statefulset | statefulsetName: StatefulSet的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - lct: the lifecycle configuration to set (for POST method) | lct: 要设置的生命周期配置(用于POST方法)
//
// Returns:
//   - []LifecycleSt: the current lifecycle configuration for all containers | []LifecycleSt: 所有容器的当前生命周期配置
//   - error: any error that occurred | error: 发生的任何错误
func StatefulsetLifecycle(kubeconfig, namespace, statefulsetName string, method string, lct LifecycleSt) ([]LifecycleSt, error) {
	var emptyLifecycleSt = make([]LifecycleSt, 0)
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] StatefulsetLifecycle Get statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
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

		stsInstance.Spec.Template.Spec.Containers[lct.ContainerId].Lifecycle = lifeCycle
		_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetLifecycle Update statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
			return emptyLifecycleSt, err
		}
		return emptyLifecycleSt, nil
	}

	lt := len(stsInstance.Spec.Template.Spec.Containers)
	var ccc = make([]LifecycleSt, 0)
	vv := stsInstance.Spec.Template.Spec.Containers
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

// StatefulsetProbe gets or updates the probe configuration of a statefulset's container
// StatefulsetProbe获取或更新StatefulSet容器的探针配置
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the statefulset | namespace: StatefulSet的命名空间
//   - statefulsetName: the name of the statefulset | statefulsetName: StatefulSet的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - probe: the probe configuration to set (for POST method) | probe: 要设置的探针配置(用于POST方法)
//
// Returns:
//   - []ProbeST: the current probe configuration for all containers | []ProbeST: 所有容器的当前探针配置
//   - error: any error that occurred | error: 发生的任何错误
func StatefulsetProbe(kubeconfig, namespace, statefulsetName string, method string, probe ProbeST) ([]ProbeST, error) {
	var emptyProbeSt = make([]ProbeST, 0)
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] StatefulsetProbe Get statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
		return emptyProbeSt, err
	}

	if method == "POST" {
		// Update readiness probe
		if probe.Readiness_checkType == "HTTP" {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: probe.Readiness_path,
						Port: intstr.FromInt32(probe.Readiness_httpPort),
					},
				},
				InitialDelaySeconds: probe.Readiness_initialDelaySeconds,
				PeriodSeconds:       probe.Readiness_periodSeconds,
				SuccessThreshold:    probe.Readiness_successThreshold,
				FailureThreshold:    probe.Readiness_failureThreshold,
				TimeoutSeconds:      probe.Readiness_timeoutSeconds,
			}
		} else if probe.Readiness_checkType == "TCP" {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.FromInt32(probe.Readiness_tcpPort),
					},
				},
				InitialDelaySeconds: probe.Readiness_initialDelaySeconds,
				PeriodSeconds:       probe.Readiness_periodSeconds,
				SuccessThreshold:    probe.Readiness_successThreshold,
				FailureThreshold:    probe.Readiness_failureThreshold,
				TimeoutSeconds:      probe.Readiness_timeoutSeconds,
			}
		} else if probe.Readiness_checkType == "CMD" {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].ReadinessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: strings.Split(probe.Readiness_cmd, " "),
					},
				},
				InitialDelaySeconds: probe.Readiness_initialDelaySeconds,
				PeriodSeconds:       probe.Readiness_periodSeconds,
				SuccessThreshold:    probe.Readiness_successThreshold,
				FailureThreshold:    probe.Readiness_failureThreshold,
				TimeoutSeconds:      probe.Readiness_timeoutSeconds,
			}
		} else {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].ReadinessProbe = nil
		}

		// Update liveness probe
		if probe.Liveness_checkType == "HTTP" {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: probe.Liveness_path,
						Port: intstr.FromInt32(probe.Liveness_httpPort),
					},
				},
				InitialDelaySeconds: probe.Liveness_initialDelaySeconds,
				PeriodSeconds:       probe.Liveness_periodSeconds,
				SuccessThreshold:    probe.Liveness_successThreshold,
				FailureThreshold:    probe.Liveness_failureThreshold,
				TimeoutSeconds:      probe.Liveness_timeoutSeconds,
			}
		} else if probe.Liveness_checkType == "TCP" {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.FromInt32(probe.Liveness_tcpPort),
					},
				},
				InitialDelaySeconds: probe.Liveness_initialDelaySeconds,
				PeriodSeconds:       probe.Liveness_periodSeconds,
				SuccessThreshold:    probe.Liveness_successThreshold,
				FailureThreshold:    probe.Liveness_failureThreshold,
				TimeoutSeconds:      probe.Liveness_timeoutSeconds,
			}
		} else if probe.Liveness_checkType == "CMD" {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].LivenessProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					Exec: &corev1.ExecAction{
						Command: strings.Split(probe.Liveness_cmd, " "),
					},
				},
				InitialDelaySeconds: probe.Liveness_initialDelaySeconds,
				PeriodSeconds:       probe.Liveness_periodSeconds,
				SuccessThreshold:    probe.Liveness_successThreshold,
				FailureThreshold:    probe.Liveness_failureThreshold,
				TimeoutSeconds:      probe.Liveness_timeoutSeconds,
			}
		} else {
			stsInstance.Spec.Template.Spec.Containers[probe.ContainerId].LivenessProbe = nil
		}

		_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetProbe Update statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
			return emptyProbeSt, err
		}
		return emptyProbeSt, nil
	}

	lt := len(stsInstance.Spec.Template.Spec.Containers)
	var ccc = make([]ProbeST, 0)
	vv := stsInstance.Spec.Template.Spec.Containers

	for i := 0; i < lt; i++ {
		result := ProbeST{}
		result.ContainerId = i
		result.ContainerName = vv[i].Name

		// Get readiness probe
		if vv[i].ReadinessProbe != nil {
			if vv[i].ReadinessProbe.HTTPGet != nil {
				result.Readiness_checkType = "HTTP"
				result.Readiness_path = vv[i].ReadinessProbe.HTTPGet.Path
				result.Readiness_httpPort = vv[i].ReadinessProbe.HTTPGet.Port.IntVal
			}
			if vv[i].ReadinessProbe.TCPSocket != nil {
				result.Readiness_checkType = "TCP"
				result.Readiness_tcpPort = vv[i].ReadinessProbe.TCPSocket.Port.IntVal
			}
			if vv[i].ReadinessProbe.Exec != nil {
				result.Readiness_checkType = "CMD"
				result.Readiness_cmd = strings.Join(vv[i].ReadinessProbe.Exec.Command, " ")
			}
			if vv[i].ReadinessProbe.InitialDelaySeconds != 0 {
				result.Readiness_initialDelaySeconds = vv[i].ReadinessProbe.InitialDelaySeconds
			}
			if vv[i].ReadinessProbe.PeriodSeconds != 0 {
				result.Readiness_periodSeconds = vv[i].ReadinessProbe.PeriodSeconds
			}
			if vv[i].ReadinessProbe.SuccessThreshold != 0 {
				result.Readiness_successThreshold = vv[i].ReadinessProbe.SuccessThreshold
			}
			if vv[i].ReadinessProbe.FailureThreshold != 0 {
				result.Readiness_failureThreshold = vv[i].ReadinessProbe.FailureThreshold
			}
			if vv[i].ReadinessProbe.TimeoutSeconds != 0 {
				result.Readiness_timeoutSeconds = vv[i].ReadinessProbe.TimeoutSeconds
			}
		}

		// Get liveness probe
		if vv[i].LivenessProbe != nil {
			if vv[i].LivenessProbe.HTTPGet != nil {
				result.Liveness_checkType = "HTTP"
				result.Liveness_path = vv[i].LivenessProbe.HTTPGet.Path
				result.Liveness_httpPort = vv[i].LivenessProbe.HTTPGet.Port.IntVal
			}
			if vv[i].LivenessProbe.TCPSocket != nil {
				result.Liveness_checkType = "TCP"
				result.Liveness_tcpPort = vv[i].LivenessProbe.TCPSocket.Port.IntVal
			}
			if vv[i].LivenessProbe.Exec != nil {
				result.Liveness_checkType = "CMD"
				result.Liveness_cmd = strings.Join(vv[i].LivenessProbe.Exec.Command, " ")
			}
			if vv[i].LivenessProbe.InitialDelaySeconds != 0 {
				result.Liveness_initialDelaySeconds = vv[i].LivenessProbe.InitialDelaySeconds
			}
			if vv[i].LivenessProbe.PeriodSeconds != 0 {
				result.Liveness_periodSeconds = vv[i].LivenessProbe.PeriodSeconds
			}
			if vv[i].LivenessProbe.SuccessThreshold != 0 {
				result.Liveness_successThreshold = vv[i].LivenessProbe.SuccessThreshold
			}
			if vv[i].LivenessProbe.FailureThreshold != 0 {
				result.Liveness_failureThreshold = vv[i].LivenessProbe.FailureThreshold
			}
			if vv[i].LivenessProbe.TimeoutSeconds != 0 {
				result.Liveness_timeoutSeconds = vv[i].LivenessProbe.TimeoutSeconds
			}
		}

		ccc = append(ccc, result)
	}
	return ccc, nil
}

type AffinitySt struct {
	NodeNames         string            `json:"nodeNames"`
	NodeSelector      []LabelsKv        `json:"nodeSelector"`
	RequiredAffinity  []MatchExpression `json:"requiredAffinity"`
	PreferredAffinity []Preference      `json:"preferredAffinity"`
}

// StatefulsetNodeAffinity gets or updates the node affinity configuration of a statefulset
// StatefulsetNodeAffinity获取或更新StatefulSet的节点亲和性配置
// Parameters:
//   - kubeconfig: the kubeconfig file path or cluster id | kubeconfig: kubeconfig文件路径或集群ID
//   - namespace: the namespace of the statefulset | namespace: StatefulSet的命名空间
//   - statefulsetName: the name of the statefulset | statefulsetName: StatefulSet的名称
//   - method: the HTTP method ("GET" or "POST") | method: HTTP方法("GET"或"POST")
//   - affinity: the affinity configuration to set (for POST method) | affinity: 要设置的亲和性配置(用于POST方法)
//
// Returns:
//   - AffinitySt: the current affinity configuration | AffinitySt: 当前亲和性配置
//   - error: any error that occurred | error: 发生的任何错误
func StatefulsetNodeAffinity(kubeconfig, namespace, statefulsetName string, method string, affinity AffinitySt) (*AffinitySt, error) {
	var emptyAffinitySt = &AffinitySt{}
	stsClient := common.ClientSet(kubeconfig).AppsV1().StatefulSets(namespace)
	stsInstance, err := stsClient.Get(context.TODO(), statefulsetName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] StatefulsetNodeAffinity Get statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
		return emptyAffinitySt, err
	}

	if method == "POST" {
		// 设置nodeName
		if affinity.NodeNames != "" {
			stsInstance.Spec.Template.Spec.NodeName = affinity.NodeNames
		} else {
			stsInstance.Spec.Template.Spec.NodeName = ""
		}

		// 设置nodeSelector
		if len(affinity.NodeSelector) > 0 {
			nodeSelector := make(map[string]string)
			for _, label := range affinity.NodeSelector {
				nodeSelector[label.Key] = label.Value
			}
			stsInstance.Spec.Template.Spec.NodeSelector = nodeSelector
		} else {
			stsInstance.Spec.Template.Spec.NodeSelector = make(map[string]string)
		}

		// 构建nodeAffinity
		nodeAffinity := &corev1.NodeAffinity{}

		// 设置requiredDuringSchedulingIgnoredDuringExecution
		if len(affinity.RequiredAffinity) > 0 {
			var nodeSelectorTerms []corev1.NodeSelectorTerm
			for _, term := range affinity.RequiredAffinity {
				var matchExpressions []corev1.NodeSelectorRequirement
				for _, expr := range term.Expressions {
					if expr.Operator == "Exists" || expr.Operator == "DoesNotExist" {
						//log.Println(expr.Operator)
						matchExpressions = append(matchExpressions, corev1.NodeSelectorRequirement{
							Key:      expr.Key,
							Operator: corev1.NodeSelectorOperator(expr.Operator),
						})
					} else {
						matchExpressions = append(matchExpressions, corev1.NodeSelectorRequirement{
							Key:      expr.Key,
							Operator: corev1.NodeSelectorOperator(expr.Operator),
							Values:   expr.Value,
						})
					}
				}
				nodeSelectorTerms = append(nodeSelectorTerms, corev1.NodeSelectorTerm{
					MatchExpressions: matchExpressions,
				})
			}
			nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &corev1.NodeSelector{
				NodeSelectorTerms: nodeSelectorTerms,
			}
		}

		// 设置preferredDuringSchedulingIgnoredDuringExecution
		if len(affinity.PreferredAffinity) > 0 {
			var preferredSchedulingTerms []corev1.PreferredSchedulingTerm
			for _, term := range affinity.PreferredAffinity {
				var matchExpressions []corev1.NodeSelectorRequirement
				for _, expr := range term.Expressions {
					if expr.Operator == "Exists" || expr.Operator == "DoesNotExist" {
						//log.Println(expr.Operator)
						matchExpressions = append(matchExpressions, corev1.NodeSelectorRequirement{
							Key:      expr.Key,
							Operator: corev1.NodeSelectorOperator(expr.Operator),
						})
					} else {
						matchExpressions = append(matchExpressions, corev1.NodeSelectorRequirement{
							Key:      expr.Key,
							Operator: corev1.NodeSelectorOperator(expr.Operator),
							Values:   expr.Value,
						})
					}
				}
				preferredSchedulingTerms = append(preferredSchedulingTerms, corev1.PreferredSchedulingTerm{
					Weight: term.Weight,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: matchExpressions,
					},
				})
			}
			nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution = preferredSchedulingTerms
		}

		if stsInstance.Spec.Template.Spec.Affinity == nil {
			stsInstance.Spec.Template.Spec.Affinity = &corev1.Affinity{}
		}
		stsInstance.Spec.Template.Spec.Affinity.NodeAffinity = nodeAffinity

		_, err = stsClient.Update(context.TODO(), stsInstance, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] StatefulsetNodeAffinity Update statefulset error: %s, namespace: %s, statefulsetName: %s\n", err, namespace, statefulsetName)
			return emptyAffinitySt, err
		}
		return emptyAffinitySt, nil
	}

	// GET logic
	result := &AffinitySt{}

	// 获取nodeName
	if stsInstance.Spec.Template.Spec.NodeName != "" {
		result.NodeNames = stsInstance.Spec.Template.Spec.NodeName
	}

	// 获取nodeSelector
	if stsInstance.Spec.Template.Spec.NodeSelector != nil {
		for k, v := range stsInstance.Spec.Template.Spec.NodeSelector {
			result.NodeSelector = append(result.NodeSelector, LabelsKv{Key: k, Value: v})
		}
	}

	// 获取nodeAffinity
	if stsInstance.Spec.Template.Spec.Affinity != nil && stsInstance.Spec.Template.Spec.Affinity.NodeAffinity != nil {
		nodeAffinity := stsInstance.Spec.Template.Spec.Affinity.NodeAffinity

		// 获取requiredDuringSchedulingIgnoredDuringExecution
		if nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
			for _, term := range nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms {
				var expressions []Expression
				for _, expr := range term.MatchExpressions {
					xValues := expr.Values
					if string(expr.Operator) == "Exists" || string(expr.Operator) == "DosNotExist" {
						xValues = []string{}
					}
					expressions = append(expressions, Expression{
						Key:      expr.Key,
						Operator: string(expr.Operator),
						Value:    xValues,
					})
				}
				result.RequiredAffinity = append(result.RequiredAffinity, MatchExpression{Expressions: expressions})
			}
		}

		// 获取preferredDuringSchedulingIgnoredDuringExecution
		if nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
			for _, term := range nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
				var expressions []Expression
				for _, expr := range term.Preference.MatchExpressions {
					xValues := expr.Values
					if string(expr.Operator) == "Exists" || string(expr.Operator) == "DosNotExist" {
						xValues = []string{}
					}
					expressions = append(expressions, Expression{
						Key:      expr.Key,
						Operator: string(expr.Operator),
						Value:    xValues,
					})
				}
				result.PreferredAffinity = append(result.PreferredAffinity, Preference{
					Weight:      term.Weight,
					Expressions: expressions,
				})
			}
		}
	}

	return result, nil
}
