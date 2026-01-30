// pod_model.go
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	//"time"
	"xkube/common"

	//"k8s.io/api"
	//appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// Kubenode k8s节点信息结构体 / Structure for k8s node information
type Kubenode struct {
	NodeName string `json:"nodeName"` //节点名称 / Node name
	NodeIp   string `json:"nodeIp"`   //节点IP / Node IP
	//NodePool    string `json:"nodePool"`    //节点池
	NodeRole      string `json:"nodeRole"`      //读取标签中的角色 / Role from labels
	NodeState     string `json:"nodeState"`     //读取conditions中为true的状态 / Node state from conditions with true status
	NodeInfo      string `json:"nodeInfo"`      //系统/kubelet/containerd/内核版本 / System/kubelet/containerd/kernel versions
	PodCIDR       string `json:"podCIDR"`       //Endpoint/PodCIDR
	Unschedulable string `json:"unschedulable"` //调度状态 / Scheduling status
	Capacity      string `json:"capacity"`      //容量/可分配 / Capacity/Allocatable
	Podnum        string `json:"podNum"`        //显示pod的数量和节点pod的额度 / Show number of pods and node pod quota
	Allocatable   string `json:"allocatable"`   //资源利用率、磁盘、内存、负载、CPU / Resource utilization, disk, memory, load, CPU
	Resusage      string `json:"resusage"`      //读取节点的metric / Read node metrics
	Cpuusage      string `json:"cpuUsage"`      //cpu利用率 / CPU utilization
	Memusage      string `json:"memUsage"`      //内存利用率 / Memory utilization
	CreateTime    string `json:"createTime"`    //创建时间 / Creation time
	//节点排水
	//调度设置
	//资源使用的性能图
	//读取所有conditions
}

// NodeDetails 节点详细信息结构体 / Structure for node detailed information
type NodeDetails struct {
	NodeName                string `json:"nodeName"`                //节点名称 / Node name
	NodeIp                  string `json:"nodeIp"`                  //节点IP / Node IP
	NodeRole                string `json:"nodeRole"`                //读取标签中的角色 / Role from labels
	NodeState               string `json:"nodeState"`               //读取conditions中为true的状态 / Node state from conditions with true status
	Architecture            string `json:"architecture"`            // 架构 / Architecture
	BootID                  string `json:"bootID"`                  // 启动ID / Boot ID
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"` // 容器运行时版本 / Container runtime version
	KernelVersion           string `json:"kernelVersion"`           // 内核版本 / Kernel version
	KubeProxyVersion        string `json:"kubeProxyVersion"`        // KubeProxy版本 / KubeProxy version
	KubeletVersion          string `json:"kubeletVersion"`          // Kubelet版本 / Kubelet version
	MachineID               string `json:"machineID"`               // 机器ID / Machine ID
	OperatingSystem         string `json:"operatingSystem"`         // 操作系统 / Operating system
	OsImage                 string `json:"osImage"`                 // 系统镜像 / OS image
	SystemUUID              string `json:"systemUUID"`              // 系统UUID / System UUID
	PodCIDR                 string `json:"podCIDR"`                 //Endpoint/PodCIDR
	//Capacity                string           `json:"capacity"`      //容量/可分配
	Allocatable   string           `json:"allocatable"`   //资源利用率、磁盘、内存、负载、CPU / Resource utilization, disk, memory, load, CPU
	Unschedulable bool             `json:"unschedulable"` //调度状态 / Scheduling status
	Taints        string           `json:"taints"`        // 污点 / Taints
	Labels        string           `json:"labels"`        // 标签 / Labels
	Annotations   string           `json:"annotations"`   // 注解 / Annotations
	Conditions    []NodeConditions `json:"conditions"`    // 条件 / Conditions
	CreateTime    string           `json:"createTime"`    //创建时间 / Creation time
}

// NodeConditions 节点状态条件结构体 / Structure for node status conditions
type NodeConditions struct {
	Ctype              string `json:"ctype"`              // 条件类型 / Condition type
	Status             string `json:"status"`             // 状态 / Status
	LastHertbeatTime   string `json:"lastHertbeatTime"`   // 上次心跳时间 / Last heartbeat time
	LastTransitionTime string `json:"lastTransitionTime"` // 上次转换时间 / Last transition time
	Reason             string `json:"reason"`             // 原因 / Reason
	Message            string `json:"message"`            // 消息 / Message
}

// KubeNodePool 节点池信息结构体 / Structure for node pool information
type KubeNodePool struct {
	NodePoolId   string `json:"nodePoolId"`   // 节点池ID / Node pool ID
	NodePoolName string `json:"nodePoolName"` // 节点池名称 / Node pool name
	NodeNumber   int64  `json:"nodeNumber"`   // 节点数量 / Number of nodes
	Remarks      string `json:"remarks"`      // 备注 / Remarks
	CreateTime   string `json:"createTime"`   // 创建时间 / Creation time
}

type NodeTaint struct {
	Key    string `json:"key"`
	Effect string `json:"effect"`
	Value  string `json:"value"`
}

// NodeList 获取节点列表 / Get node list
// kubeconfig: 集群认证文件 / Cluster authentication file
// 返回节点列表和错误信息 / Returns node list and error information
func NodeList(kubeconfig string) ([]Kubenode, error) {
	//通过kubeconfig集群认证文件生成一个客户端操作对象clientset
	// Generate a client operation object clientset through the kubeconfig cluster authentication file
	clientset := common.ClientSet(kubeconfig)

	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list nodes error:%v\n", err)
	}
	var NodeArry = make([]Kubenode, 0)
	//fmt.Println("dev node count:", len(nodeList.Items))
	for _, node := range nodeList.Items {

		// var labelsStr string
		// for kk, vv := range node.ObjectMeta.Labels {
		// 	labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		// }
		var nodeStatus = "Ready:True"
		for _, vv := range node.Status.Conditions {
			if vv.Status == "True" || vv.Status == "Unknown" {
				nodeStatus = fmt.Sprintf("%s:%s", vv.Type, vv.Status)
				break
			} else {
				nodeStatus = fmt.Sprintf("%s:%s", vv.Type, vv.Status)
			}
		}

		var nodeRole = node.Labels["node-role.kubernetes.io/work"]
		if nodeRole == "" {
			nodeRole = "Worker"
		}

		// var bbb = make([]NodeConditions, 0)
		// for _, v1 := range node.Status.Conditions {
		// 	xItems := &NodeConditions{
		// 		LastTransitionTime: v1.LastTransitionTime.Format("2006-01-02 15:04:05"),
		// 		LastHertbeatTime:   v1.LastHeartbeatTime.Format("2006-01-02 15:04:05"),
		// 		Message:            v1.Message,
		// 		Reason:             v1.Reason,
		// 		Status:             fmt.Sprintf("%v", v1.Status),
		// 		Ctype:              fmt.Sprintf("%v", v1.Type),
		// 	}
		// 	bbb = append(bbb, *xItems)
		// }

		//resUsage, _ := GetNodeMetricV2(kubeconfig, node.Name)
		cpuValue, memValue, _ := GetNodeMetricV3(kubeconfig, node.Name)

		resUsage := fmt.Sprintf("CPU:%.3f核<br>Mem:%.3fGB", float64(cpuValue)/1000, float64(memValue)/1024/1024/1024)
		//fmt.Printf("cpu:%d,%d\n", cpuValue, node.Status.Capacity.Cpu().Value() )
		//fmt.Printf("mem:%d,%d\n", memValue, node.Status.Capacity.Memory().Value() )
		cpuUsage := fmt.Sprintf("%0.2f", float64(cpuValue)/float64(node.Status.Capacity.Cpu().Value()*1000)*100)
		memUsage := fmt.Sprintf("%0.2f", float64(memValue)/float64(node.Status.Capacity.Memory().Value())*100)

		//imgNum := len(node.Status.Images)
		//tm3 := time.Now().UnixNano() / 1e6
		podNum, _ := PodCountByNode(kubeconfig, node.Name)
		//tm4 := time.Now().UnixNano() / 1e6
		//fmt.Printf("PodCountByNode:%d\n", tm4-tm3)
		//maxPods := node.Status.Capacity["pods"]
		//maxPods.String()
		NodeArry = append(NodeArry, Kubenode{
			NodeName: node.Name,
			NodeIp:   node.Status.Addresses[0].Address,
			//NodePool:   "",
			NodeRole:      nodeRole,
			NodeState:     nodeStatus,
			NodeInfo:      fmt.Sprintf("KernelVersion:%s,OS:%s,KubeletVersion:%s,Container:%s", node.Status.NodeInfo.KernelVersion, node.Status.NodeInfo.OSImage, node.Status.NodeInfo.KubeletVersion, node.Status.NodeInfo.ContainerRuntimeVersion),
			PodCIDR:       fmt.Sprintf("PodCIDR:%s", node.Spec.PodCIDR),
			Unschedulable: fmt.Sprintf("%v", node.Spec.Unschedulable),
			Capacity:      fmt.Sprintf("CPU:%d,Mem:%dMi,Disk:%dGi,Pods:%d", node.Status.Capacity.Cpu().Value(), node.Status.Capacity.Memory().Value()/1024/1024, node.Status.Capacity.StorageEphemeral().Value()/1024/1024/1024, node.Status.Capacity.Pods().Value()),
			Allocatable:   fmt.Sprintf("CPU:%d,Mem:%dMi,Disk:%dGi,Pods:%d", node.Status.Allocatable.Cpu().Value(), node.Status.Allocatable.Memory().Value()/1024/1024, node.Status.Allocatable.StorageEphemeral().Value()/1024/1024/1024, node.Status.Allocatable.Pods().Value()-int64(podNum)),
			//Labels:     labelsStr[0 : len(labelsStr)-1],
			Podnum:     fmt.Sprintf("%d/%d", podNum, node.Status.Capacity.Pods().Value()),
			Resusage:   resUsage,
			Cpuusage:   cpuUsage,
			Memusage:   memUsage,
			CreateTime: node.ObjectMeta.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return NodeArry, err
}

// NodeDetail 获取节点详细信息 / Get node detailed information
// kubeconfig: 集群认证文件 / Cluster authentication file
// nodeName: 节点名称 / Node name
// 返回节点详细信息和错误信息 / Returns node detailed information and error information
func NodeDetail(kubeconfig, nodeName string) (*NodeDetails, error) {
	node, err := common.ClientSet(kubeconfig).CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("list pods error:%v\n", err)
	}
	var labelsStr, annotationsStr, taintsStr string
	for k1, v1 := range node.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("<br>%s:%s", k1, v1)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	for kk, vv := range node.ObjectMeta.Annotations {
		annotationsStr += fmt.Sprintf("<br>%s:%s", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	for _, v2 := range node.Spec.Taints {
		taintsStr += fmt.Sprintf("%s,", v2.String())
	}
	if len(taintsStr) > 0 {
		taintsStr = taintsStr[0 : len(taintsStr)-1]
	}

	var nodeStatus = "Ready:True"
	for _, vv := range node.Status.Conditions {
		if vv.Status == "True" || vv.Status == "Unknown" {
			nodeStatus = fmt.Sprintf("%s:%s", vv.Type, vv.Status)
			break
		} else {
			nodeStatus = fmt.Sprintf("%s:%s", vv.Type, vv.Status)
		}
	}

	var nodeRole = node.Labels["node-role.kubernetes.io/work"]
	if nodeRole == "" {
		nodeRole = "Worker"
	}

	// 创建节点条件切片并预分配容量以提高性能
	conditions := make([]NodeConditions, 0, len(node.Status.Conditions))

	for _, condition := range node.Status.Conditions {
		conditions = append(conditions, NodeConditions{
			LastTransitionTime: condition.LastTransitionTime.Format("2006-01-02 15:04:05"),
			LastHertbeatTime:   condition.LastHeartbeatTime.Format("2006-01-02 15:04:05"),
			Message:            condition.Message,
			Reason:             condition.Reason,
			Status:             string(condition.Status),
			Ctype:              string(condition.Type),
		})
	}
	//node.Spec.Unschedulable
	//node.Spec.Taints[0].String() //污点
	//node.Status.Images
	return &NodeDetails{
		NodeName:                node.Name,
		NodeIp:                  node.Status.Addresses[0].Address,
		NodeRole:                nodeRole,
		NodeState:               nodeStatus,
		Architecture:            node.Status.NodeInfo.Architecture,
		BootID:                  node.Status.NodeInfo.BootID,
		ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
		KernelVersion:           node.Status.NodeInfo.KernelVersion,
		KubeProxyVersion:        node.Status.NodeInfo.KubeProxyVersion,
		KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
		MachineID:               node.Status.NodeInfo.MachineID,
		OperatingSystem:         node.Status.NodeInfo.OperatingSystem,
		OsImage:                 node.Status.NodeInfo.OSImage,
		SystemUUID:              node.Status.NodeInfo.SystemUUID,
		PodCIDR:                 fmt.Sprintf("PodCIDR:%s", node.Spec.PodCIDR),
		//Capacity:                fmt.Sprintf("CPU:%d,Mem:%dMi,Disk:%dGi,Pods:%d", node.Status.Capacity.Cpu().Value(), node.Status.Capacity.Memory().Value()/1024/1024, node.Status.Capacity.StorageEphemeral().Value()/1024/1024/1024, node.Status.Capacity.Pods().Value()),
		Allocatable:   fmt.Sprintf("CPU:%d<br>内存:%dMi<br>磁盘:%dGi<br>Pods:%d", node.Status.Allocatable.Cpu().Value(), node.Status.Allocatable.Memory().Value()/1024/1024, node.Status.Allocatable.StorageEphemeral().Value()/1024/1024/1024, node.Status.Allocatable.Pods().Value()),
		Unschedulable: node.Spec.Unschedulable,
		Taints:        taintsStr,
		Labels:        labelsStr,
		Annotations:   annotationsStr,
		Conditions:    conditions,
		CreateTime:    node.ObjectMeta.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

// NodePoolList 获取节点池列表 / Get node pool list
// kubeconfig: 集群认证文件 / Cluster authentication file
// 返回节点池列表和错误信息 / Returns node pool list and error information
func NodePoolList(kubeconfig string) ([]KubeNodePool, error) {
	//通过kubeconfig集群认证文件生成一个客户端操作对象clientset
	// Generate a client operation object clientset through the kubeconfig cluster authentication file
	clientset := common.ClientSet(kubeconfig)

	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list pods error:%v\n", err)
	}
	var NodeArry = make([]KubeNodePool, 0)
	//fmt.Println("dev node count:", len(nodeList.Items))
	var poolMap = make(map[string]int64)
	for _, node := range nodeList.Items {
		if vv, ok := node.ObjectMeta.Labels["alibabacloud.com/nodepool-id"]; ok {
			poolMap[vv] += 1
		}
	}
	for k1, v1 := range poolMap {
		NodeArry = append(NodeArry, KubeNodePool{
			NodePoolId:   k1,
			NodePoolName: "",
			NodeNumber:   v1,
			Remarks:      "",
			CreateTime:   "",
		})
	}
	return NodeArry, err
}

// GetNodeYaml 获取节点YAML信息 / Get node YAML information
// kubeconfig: 集群认证文件 / Cluster authentication file
// nodeName: 节点名称 / Node name
// 返回节点YAML字符串和错误信息 / Returns node YAML string and error information
func GetNodeYaml(kubeconfig, nodeName string) (string, error) {
	node, err := common.ClientSet(kubeconfig).CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("list pods error:%v\n", err)
	}
	nodeUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(node)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(nodeUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

func NodeYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	node := &corev1.Node{}
	err = json.Unmarshal(data, node)
	if err != nil {
		return err
	}
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	return err
}

// NodeDrain 节点排水操作 / Node drain operation
// kubeconfig: 集群认证文件 / Cluster authentication file
// nodeName: 节点名称 / Node name
// 返回错误信息 / Returns error information
// k8s version > 1.22
// 配置了节点亲和性的无法排除
func NodeDrain(kubeconfig, nodeName string) error {
	clientset := common.ClientSet(kubeconfig)
	// podlist, err := PodList(kubeconfig, "", "", "", "", "", nodeName)
	// if err != nil {
	// 	log.Printf("[ERROR] NodeDrain GetPodList Fail:%s\n", err)
	// 	return err
	// }
	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	}
	podList, err := clientset.CoreV1().Pods("").List(context.Background(), listOptions)
	if err != nil {
		log.Printf("[ERROR] NodeDrain GetPodList Fail:%s\n", err)
		return err
	}
	for _, pod := range podList.Items {
		if pod.ObjectMeta.OwnerReferences != nil {
			var ds bool
			for _, ref := range pod.ObjectMeta.OwnerReferences {
				if ref.Kind == "DaemonSet" {
					ds = true
					continue
				}
			}
			if ds {
				continue
			}
		}
		eviction := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			DeleteOptions: &metav1.DeleteOptions{},
		}
		//如遇policy/v1未注册错误，需降级使用policy/v1beta1版本（K8s <1.22时）1，修改Eviction对象导入路径为policy/v1beta1并调整API调用。
		err = clientset.PolicyV1().Evictions(pod.Namespace).Evict(context.Background(), eviction)
		if err != nil {
			log.Printf("[ERROR] NodeDrain pod Evict Fail:%s\n", err)
			return err
		}
	}
	return nil
}

// NodeDrain 节点排水操作 / Node drain operation
// kubeconfig: 集群认证文件 / Cluster authentication file
// nodeName: 节点名称 / Node name
// 返回错误信息 / Returns error information
// k8s version < 1.22
func NodeDrainV1beta1(kubeconfig, nodeName string) error {
	clientset := common.ClientSet(kubeconfig)
	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	}
	podList, err := clientset.CoreV1().Pods("").List(context.Background(), listOptions)
	if err != nil {
		log.Printf("[ERROR] NodeDrain GetPodList Fail:%s\n", err)
		return err
	}
	for _, pod := range podList.Items {
		if pod.ObjectMeta.OwnerReferences != nil {
			for _, ref := range pod.ObjectMeta.OwnerReferences {
				if ref.Kind == "DaemonSet" {
					continue
				}
			}
		}
		eviction := &policyv1beta1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			DeleteOptions: &metav1.DeleteOptions{},
		}
		//如遇policy/v1未注册错误，需降级使用policy/v1beta1版本（K8s <1.22时）1，修改Eviction对象导入路径为policy/v1beta1并调整API调用。
		err = clientset.PolicyV1beta1().Evictions(pod.Namespace).Evict(context.Background(), eviction)
		if err != nil {
			log.Printf("[ERROR] NodeDrain pod Evict Fail:%s\n", err)
			return err
		}
	}
	return nil
}

// NodeUnschedulable 设置节点调度状态 / Set node scheduling status
// kubeconfig: 集群认证文件 / Cluster authentication file
// nodeName: 节点名称 / Node name
// value: 调度状态值 / Scheduling status value
// 返回错误信息 / Returns error information
func NodeUnschedulable(kubeconfig, nodeName string, value bool) error {
	clientset := common.ClientSet(kubeconfig)
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	node.Spec.Unschedulable = value
	_, err = clientset.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	return err
}

// NodeDelete 删除节点 / Delete node
// kubeconfig: 集群认证文件 / Cluster authentication file
// nodeName: 节点名称 / Node name
// 返回错误信息 / Returns error information
func NodeDelete(kubeconfig, nodeName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := common.ClientSet(kubeconfig).CoreV1().Nodes().Delete(context.TODO(), nodeName, metav1.DeleteOptions{
		GracePeriodSeconds: new(int64),
		PropagationPolicy:  &deletePolicy,
	})
	return err
}

func GetNodeLabels(kubeconfig, nodeName string) ([]LabelsKv, error) {
	var bbb = make([]LabelsKv, 0)
	nodeClient := common.ClientSet(kubeconfig).CoreV1().Nodes()
	node, err := nodeClient.Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] GetNodeLabels error:%v\n", err)
		return bbb, err
	}
	if node.ObjectMeta.Labels != nil {
		for kk, vv := range node.ObjectMeta.Labels {
			bbb = append(bbb, LabelsKv{
				Key:   kk,
				Value: vv,
			})
		}
	}
	return bbb, nil
}

func UpdateNodeLabels(kubeconfig, nodeName string, labelsMap map[string]string) error {
	nodeClient := common.ClientSet(kubeconfig).CoreV1().Nodes()
	node, err := nodeClient.Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR]Get Node error:%v\n", err)
		return err
	}
	node.ObjectMeta.Labels = labelsMap
	_, err = nodeClient.Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("[ERROR] NodeLabels Update error:%s\n", err)
		return err
	}
	return nil

}

func GetNodeTaint(kubeconfig, nodeName string) ([]NodeTaint, error) {
	var bbb = make([]NodeTaint, 0)
	nodeClient := common.ClientSet(kubeconfig).CoreV1().Nodes()
	node, err := nodeClient.Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] GetNodeLabels error:%v\n", err)
		return bbb, err
	}

	if node.Spec.Taints != nil {
		for _, vv := range node.Spec.Taints {
			var effectStr string
			if vv.Effect == corev1.TaintEffectNoExecute {
				effectStr = "NoExecute"
			} else if vv.Effect == corev1.TaintEffectNoSchedule {
				effectStr = "NoSchedule"
			} else {
				effectStr = "PreferNoSchedule"
			}
			bbb = append(bbb, NodeTaint{
				Key:    vv.Key,
				Effect: effectStr,
				Value:  vv.Value,
			})
		}
	}
	return bbb, nil
}

func UpdateNodeTaint(kubeconfig, nodeName string, taintArry []NodeTaint) error {
	nodeClient := common.ClientSet(kubeconfig).CoreV1().Nodes()
	node, err := nodeClient.Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR]Get Node error:%v\n", err)
		return err
	}
	//corev1.Taint{}
	//corev1.TaintEffect
	var taints []corev1.Taint
	for _, vv := range taintArry {
		var effectStr corev1.TaintEffect
		if vv.Effect == "NoSchedule" {
			effectStr = corev1.TaintEffectNoSchedule
		} else if vv.Effect == "NoExecute" {
			effectStr = corev1.TaintEffectNoExecute
		} else if vv.Effect == "PreferNoSchedule" {
			effectStr = corev1.TaintEffectPreferNoSchedule
		}
		taints = append(taints, corev1.Taint{
			Key:    vv.Key,
			Value:  vv.Value,
			Effect: effectStr,
		})
	}

	node.Spec.Taints = taints
	_, err = nodeClient.Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("[ERROR] NodeLabels Update error:%s\n", err)
		return err
	}
	return nil

}

// NodeAllocatedResource represents resource allocation information for a node
type NodeAllocatedResource struct {
	Resource       string `json:"resource"`
	Requests       string `json:"requests"`
	Limits         string `json:"limits"`
	RequestPercent string `json:"requestPercent"`
	LimitPercent   string `json:"limitPercent"`
}

// NodeAllocated returns allocated resources information for a node similar to kubectl describe node
func NodeAllocated(kubeconfig, nodeName string) (map[string]interface{}, error) {
	clientset := common.ClientSet(kubeconfig)
	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Get all pods on this node
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return nil, err
	}

	// Calculate total requests and limits
	var cpuRequests, cpuLimits, memoryRequests, memoryLimits int64

	for _, pod := range pods.Items {
		// Skip failed or succeeded pods
		if pod.Status.Phase == corev1.PodFailed || pod.Status.Phase == corev1.PodSucceeded {
			continue
		}

		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				if cpu, exists := container.Resources.Requests[corev1.ResourceCPU]; exists {
					cpuRequests += cpu.MilliValue()
				}
				if memory, exists := container.Resources.Requests[corev1.ResourceMemory]; exists {
					memoryRequests += memory.Value()
				}
			}

			if container.Resources.Limits != nil {
				if cpu, exists := container.Resources.Limits[corev1.ResourceCPU]; exists {
					cpuLimits += cpu.MilliValue()
				}
				if memory, exists := container.Resources.Limits[corev1.ResourceMemory]; exists {
					memoryLimits += memory.Value()
				}
			}
		}
	}

	// Get node capacity
	nodeCPU := node.Status.Capacity[corev1.ResourceCPU]
	nodeMemory := node.Status.Capacity[corev1.ResourceMemory]

	// Calculate percentages
	cpuRequestPercent := float64(cpuRequests) / float64(nodeCPU.MilliValue()) * 100
	cpuLimitPercent := float64(cpuLimits) / float64(nodeCPU.MilliValue()) * 100
	memoryRequestPercent := float64(memoryRequests) / float64(nodeMemory.Value()) * 100
	memoryLimitPercent := float64(memoryLimits) / float64(nodeMemory.Value()) * 100

	// Prepare response data
	data := make([]NodeAllocatedResource, 2)

	// CPU data
	data[0] = NodeAllocatedResource{
		Resource:       "CPU",
		Requests:       fmt.Sprintf("%dm", cpuRequests),
		Limits:         fmt.Sprintf("%dm", cpuLimits),
		RequestPercent: fmt.Sprintf("%.2f%%", cpuRequestPercent),
		LimitPercent:   fmt.Sprintf("%.2f%%", cpuLimitPercent),
	}

	// Memory data
	data[1] = NodeAllocatedResource{
		Resource:       "Memory",
		Requests:       fmt.Sprintf("%dMi", memoryRequests/1024/1024),
		Limits:         fmt.Sprintf("%dMi", memoryLimits/1024/1024),
		RequestPercent: fmt.Sprintf("%.2f%%", memoryRequestPercent),
		LimitPercent:   fmt.Sprintf("%.2f%%", memoryLimitPercent),
	}

	response := map[string]interface{}{
		"code":  0,
		"msg":   "success",
		"count": 2,
		"data":  data,
	}

	return response, nil
}
