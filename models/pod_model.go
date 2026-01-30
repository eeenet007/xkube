// pod_model.go
// Pod数据模型定义及操作函数 / Pod data model definition and operation functions
package models

import (
	"bytes"
	"context"
	"fmt"

	"io"
	//"io/ioutil"
	"log"
	"os"
	"strings"

	//"time"
	"strconv"
	"xkube/common"

	//"golang.org/x/text/transform"

	//"k8s.io/api"
	"golang.org/x/text/encoding/simplifiedchinese"
	//appsv1 "k8s.io/api/apps/v1"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// Podinfo Pod信息结构体 / Pod information struct
type Podinfo struct {
	PodName   string `json:"podName"` //podname 和podip / pod name and pod ip
	PodIp     string `json:"podIp"`
	NameSpace string `json:"nameSpace"`
	NodeName  string `json:"nodeName"` //节点名称和节点IP / node name and node IP
	HostIp    string `json:"hostIp"`
	PodPhase  string `json:"podPhase"` //容器状态 Running / container status Running
	ImgUrl    string `json:"imgUrl"`
	//PodStatus    PodStatusItem
	RestartCount int32   `json:"restartCount"` //重启次数 / restart count
	Labels       string  `json:"labels"`       //标签 / labels
	ResUsage     string  `json:"resUsage"`     //cpu及内存使用率 / cpu and memory usage
	CpuUsage     float64 `json:"cpuUsage"`
	MemUsage     float64 `json:"memUsage"`
	CreateTime   string  `json:"createTime"`
}

// PodDetails Pod详细信息结构体 / Pod details struct
type PodDetails struct {
	PodName        string          `json:"podName"` //podname 和podip / pod name and pod ip
	PodIp          string          `json:"podIp"`
	NameSpace      string          `json:"nameSpace"`
	NodeName       string          `json:"nodeName"` //节点名称和节点IP / node name and node IP
	HostIp         string          `json:"hostIp"`
	PodPhase       string          `json:"podPhase"`       //容器状态 Running / container status Running
	RestartCount   int32           `json:"restartCount"`   //重启次数 / restart count
	OwnerKind      string          `json:"ownerKind"`      //创建者类型 / creator type
	OwnerName      string          `json:"ownerName"`      //创建者名称 / creator name
	Labels         string          `json:"labels"`         //标签 / labels
	Lastterminated string          `json:"lastTerminated"` //上一次重启的原因及时间 / reason and time of last restart
	Annotations    string          `json:"annotations"`
	CreateTime     string          `json:"createTime"`
	Conditions     []PodConditions `json:"conditions"`
	Volumes        []PodVolumes    `json:"volumes"`
	Containers     []Container     `json:"containers"`
}

// PodConditions Pod状态条件 / Pod status conditions
type PodConditions struct {
	LastTransitionTime string `json:"lastTransitionTime"`
	LastUpdateTime     string `json:"lastUpdateTime"`
	Message            string `json:"message"`
	Reason             string `json:"reason"`
	Status             string `json:"status"`
	Ctype              string `json:"ctype"`
}

// PodVolumes Pod存储卷信息 / Pod volume information
type PodVolumes struct {
	VolumeName string `json:"volumeName"`
	VolumeType string `json:"volumeType"`
	VolumeInfo string `json:"volumeInfo"`
}

// Container 容器信息 / Container information
type Container struct {
	ContainerName  string `json:"containerName"`
	Envs           string `json:"envs"`
	Mounts         string `json:"mounts"`
	ContainerImage string `json:"containerImage"`
	PullPolicy     string `json:"pullPolicy"`
	Ports          string `json:"ports"`
	ResLimits      string `json:"resLimits"`
	ResRequests    string `json:"resRequests"`
	ResUsed        string `json:"resUsed"`
}

// type PodStatusItem struct {
// 	Initialized     string `json:"initialized"`
// 	Ready           string `json:"ready"`
// 	ContainersReady string `json:"containersReady"`
// 	PodScheduled    string `json:"podScheduled"`
// }

// PodList 获取Pod列表 / Get Pod list
// kubeconfig: k8s集群配置 / k8s cluster configuration
// namespace: 命名空间 / namespace
// deployName: 部署名称 / deployment name
// podName: Pod名称 / Pod name
// labelsKey, labelsValue: 标签键值对 / label key-value pair
// nodeName: 节点名称 / node name
func PodList(kubeconfig, namespace, deployName, podName string, labelsKey, labelsValue, nodeName string) ([]Podinfo, error) {

	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	clientset := common.ClientSet(kubeconfig)
	var podList *corev1.PodList
	var err error

	//设置ListOptions / Set ListOptions
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
			//FieldSelector: "status.phase=Running,spec.nodeName=ais-master1",
		}
	}

	if nodeName != "" {
		listOptions = metav1.ListOptions{
			FieldSelector: fmt.Sprintf("status.phase!=Succeeded,spec.nodeName=%s", nodeName),
		}
	}

	podList, err = clientset.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		log.Printf("list pods error:%v\n", err)
	}
	var bbb = make([]Podinfo, 0)
	for _, pod := range podList.Items {
		//搜索 / Search
		if podName != "" {
			if !strings.Contains(pod.Name, podName) {
				//if !strings.HasPrefix(pod.Name,podName)
				continue
			}
		}
		var resUsage string
		var cpuUsage, memUsage float64
		if deployName != "" || nodeName != "" {
			if !strings.HasPrefix(pod.Name, deployName) {
				continue
			}
			cpuValue, memValue, _ := GetPodMetricV2(kubeconfig, pod.ObjectMeta.Namespace, pod.Name)
			resUsage = fmt.Sprintf("cpu:%.3f核,mem:%.3fMB", float64(cpuValue)/1000, float64(memValue)/1024/1024)
			cpuUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(cpuValue)/1000), 64)
			memUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(memValue)/1024/1024), 64)
		}

		var labelsStr, imgurlStr string
		for kk, vv := range pod.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		var containerState = fmt.Sprintf("%v", pod.Status.Phase)

		var restartNum int32
		if len(pod.Status.ContainerStatuses) > 0 {
			imgurlStr = pod.Status.ContainerStatuses[0].Image
			restartNum = pod.Status.ContainerStatuses[0].RestartCount
			started := pod.Status.ContainerStatuses[0].Started
			//if containerState == "Pending" {
			if !*started {
				// if containerState == "Failed" {
				// 	containerState = pod.Status.ContainerStatuses[0].State.Terminated.Reason
				// } else {
				// 	containerState = pod.Status.ContainerStatuses[0].State.Waiting.Reason
				// }
				if pod.Status.ContainerStatuses[0].State.Terminated != nil {
					containerState = pod.Status.ContainerStatuses[0].State.Terminated.Reason
				}
				if pod.Status.ContainerStatuses[0].State.Waiting != nil {
					containerState = pod.Status.ContainerStatuses[0].State.Waiting.Reason
				}
			}
			if pod.ObjectMeta.DeletionTimestamp != nil {
				containerState = "Terminating"
			}
		}

		Items := &Podinfo{
			PodName:      pod.Name,
			PodIp:        pod.Status.PodIP,
			ImgUrl:       imgurlStr,
			NameSpace:    pod.ObjectMeta.Namespace,
			NodeName:     pod.Spec.NodeName,
			HostIp:       pod.Status.HostIP,
			PodPhase:     containerState,
			Labels:       labelsStr,
			RestartCount: restartNum,
			ResUsage:     resUsage,
			MemUsage:     memUsage,
			CpuUsage:     cpuUsage,
			CreateTime:   pod.ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// PodListV2 获取Pod列表V2版本 / Get Pod list V2 version
// kubeconfig: k8s集群配置 / k8s cluster configuration
// namespace: 命名空间 / namespace
// resName: 资源名称 / resource name
// resType: 资源类型 / resource type
// podName: Pod名称 / Pod name
// labelsKey, labelsValue: 标签键值对 / label key-value pair
// nodeName: 节点名称 / node name
func PodListV2(kubeconfig, namespace, resName, resType, podName string, labelsKey, labelsValue, nodeName string) ([]Podinfo, error) {

	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	clientset := common.ClientSet(kubeconfig)
	var podList *corev1.PodList
	var err error

	//设置ListOptions / Set ListOptions
	var listOptions = metav1.ListOptions{
		Limit: 4000,
	}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
			Limit:         4000,
			//FieldSelector: "status.phase=Running,spec.nodeName=ais-master1",
		}
	}
	if nodeName != "" {
		listOptions = metav1.ListOptions{
			FieldSelector: fmt.Sprintf("status.phase!=Succeeded,spec.nodeName=%s", nodeName),
			Limit:         400,
		}
	}

	///job详情页是靠标签 / job details page relies on labels
	if resName != "" {
		if resType == "deploy" {
			deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
			if err != nil {
				log.Printf("PodList Error getting deployment: %s\n", err)
			}
			// 通过Deployment的标签选择器查找对应的Pod / Find corresponding Pods through Deployment's label selector
			labelSelector := metav1.FormatLabelSelector(deployment.Spec.Selector)
			listOptions = metav1.ListOptions{
				LabelSelector: labelSelector,
				Limit:         100,
			}
		}
		if resType == "sts" {
			sts, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
			if err != nil {
				log.Printf("PodListError getting sts: %s\n", err)
			}
			labelSelector := metav1.FormatLabelSelector(sts.Spec.Selector)
			listOptions = metav1.ListOptions{
				LabelSelector: labelSelector,
				Limit:         100,
			}
		}
		if resType == "ds" {
			ds, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
			if err != nil {
				log.Printf("PodListError getting ds: %s\n", err)
			}
			labelSelector := metav1.FormatLabelSelector(ds.Spec.Selector)
			listOptions = metav1.ListOptions{
				LabelSelector: labelSelector,
				Limit:         100,
			}
		}
	}

	podList, err = clientset.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		log.Printf("list pods error:%v\n", err)
	}
	var bbb = make([]Podinfo, 0)
	for _, pod := range podList.Items {
		//搜索 / Search
		if podName != "" {
			if !strings.Contains(pod.Name, podName) {
				//if !strings.HasPrefix(pod.Name,podName)
				continue
			}
		}
		var resUsage string
		var cpuUsage, memUsage float64
		if resName != "" || nodeName != "" {
			//if !strings.HasPrefix(pod.Name, resName) {
			//	continue
			//}
			cpuValue, memValue, _ := GetPodMetricV2(kubeconfig, pod.ObjectMeta.Namespace, pod.Name)
			resUsage = fmt.Sprintf("cpu:%.3f核,mem:%.3fMB", float64(cpuValue)/1000, float64(memValue)/1024/1024)
			cpuUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(cpuValue)/1000), 64)
			memUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(memValue)/1024/1024), 64)
		}

		var labelsStr, imgurlStr string
		for kk, vv := range pod.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		var containerState = fmt.Sprintf("%v", pod.Status.Phase)

		var restartNum int32
		if len(pod.Status.ContainerStatuses) > 0 {
			imgurlStr = pod.Status.ContainerStatuses[0].Image
			restartNum = pod.Status.ContainerStatuses[0].RestartCount
			started := pod.Status.ContainerStatuses[0].Started
			//if containerState == "Pending" {
			if !*started {

				if pod.Status.ContainerStatuses[0].State.Terminated != nil {
					containerState = pod.Status.ContainerStatuses[0].State.Terminated.Reason
				}
				if pod.Status.ContainerStatuses[0].State.Waiting != nil {
					containerState = pod.Status.ContainerStatuses[0].State.Waiting.Reason
				}
			}
			if pod.ObjectMeta.DeletionTimestamp != nil {
				containerState = "Terminating"
			}
		}

		Items := &Podinfo{
			PodName:      pod.Name,
			PodIp:        pod.Status.PodIP,
			ImgUrl:       imgurlStr,
			NameSpace:    pod.ObjectMeta.Namespace,
			NodeName:     pod.Spec.NodeName,
			HostIp:       pod.Status.HostIP,
			PodPhase:     containerState,
			Labels:       labelsStr,
			RestartCount: restartNum,
			ResUsage:     resUsage,
			MemUsage:     memUsage,
			CpuUsage:     cpuUsage,
			CreateTime:   pod.ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// PodCountByNode 统计节点上的Pod数量 / Count the number of Pods on a node
// kubeconfig: k8s集群配置 / k8s cluster configuration
// nodeName: 节点名称 / node name
func PodCountByNode(kubeconfig, nodeName string) (int, error) {
	resp := common.Get("podCount" + kubeconfig + nodeName)
	if resp != "" {
		num, err := strconv.Atoi(resp)
		if err != nil {
			log.Printf("[WARN]PodCountByNode Atoi err:%s\n", err)
		}
		return num, nil
	}
	namespace := corev1.NamespaceAll
	clientset := common.ClientSet(kubeconfig)
	listOptions := metav1.ListOptions{
		FieldSelector: fmt.Sprintf("status.phase!=Succeeded,spec.nodeName=%s", nodeName),
	}
	podList, err := clientset.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		log.Printf("[WARN]PodCountByNode list pods error:%v\n", err)
	}
	count := len(podList.Items)
	_ = common.SetEx("podCount"+kubeconfig+nodeName, fmt.Sprintf("%d", count), 1800)
	return count, err
}

// PodDetail 获取Pod详细信息 / Get Pod details
// kubeconfig: k8s集群配置 / k8s cluster configuration
// nameSpace: 命名空间 / namespace
// podName: Pod名称 / Pod name
func PodDetail(kubeconfig, nameSpace, podName string) (*PodDetails, error) {
	pod, err := common.ClientSet(kubeconfig).CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	var labelsStr, annotationsStr string
	for kk, vv := range pod.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	for kk, vv := range pod.ObjectMeta.Annotations {
		annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	var bbb = make([]PodConditions, 0)
	for _, v1 := range pod.Status.Conditions {
		xItems := &PodConditions{
			LastTransitionTime: v1.LastTransitionTime.Format("2006-01-02 15:04:05"),
			LastUpdateTime:     v1.LastTransitionTime.Format("2006-01-02 15:04:05"),
			Message:            v1.Message,
			Reason:             v1.Reason,
			Status:             fmt.Sprintf("%v", v1.Status),
			Ctype:              fmt.Sprintf("%v", v1.Type),
		}
		bbb = append(bbb, *xItems)
	}

	var containerState = fmt.Sprintf("%v", pod.Status.Phase)

	var restartNum int32
	if len(pod.Status.ContainerStatuses) > 0 {
		restartNum = pod.Status.ContainerStatuses[0].RestartCount
		if containerState == "Pending" {
			if pod.Status.ContainerStatuses[0].State.Waiting != nil {
				containerState = pod.Status.ContainerStatuses[0].State.Waiting.Reason
			}
		}
	}

	var vvv = make([]PodVolumes, 0)
	for _, v1 := range pod.Spec.Volumes {
		//log.Println(v1.VolumeSource)
		var volumeType, volumeInfo string
		if v1.ConfigMap != nil {
			volumeType = "configMap"
			volumeInfo = fmt.Sprintf("name:%s,defaultMode:%d", v1.ConfigMap.Name, *v1.ConfigMap.DefaultMode)
		}

		if v1.Secret != nil {
			volumeType = "secret"
			volumeInfo = fmt.Sprintf("name:%s,defaultMode:%d", v1.Secret.SecretName, *v1.Secret.DefaultMode)
		}

		if v1.PersistentVolumeClaim != nil {
			volumeType = "persistentVolumeClaim"
			volumeInfo = fmt.Sprintf("claimName:%s", v1.PersistentVolumeClaim.ClaimName)
		}

		if v1.EmptyDir != nil {
			volumeType = "emptyDir"
			volumeInfo = fmt.Sprintf("emptyDirMedium:%s", v1.EmptyDir.Medium)
		}

		if v1.HostPath != nil {
			volumeType = "hostPath"
			volumeInfo = fmt.Sprintf("path:%s,type:%s", v1.HostPath.Path, *v1.HostPath.Type)
		}

		xItems := &PodVolumes{
			VolumeName: v1.Name,
			VolumeType: volumeType,
			VolumeInfo: volumeInfo,
		}
		vvv = append(vvv, *xItems)
	}

	var ccc = make([]Container, 0)

	var ContainerResUsed = []ContainerMetirc{}
	if containerState == "Running" {
		ContainerResUsed, _ = GetContainerMetric(kubeconfig, nameSpace, podName)
	}

	for _, v1 := range pod.Spec.Containers {
		var envsStr, mountsStr, PortsStr string
		for _, ve := range v1.Env {
			if ve.Value != "" {
				envsStr += fmt.Sprintf("%s:%s,", ve.Name, ve.Value)
			}
			if ve.ValueFrom != nil {
				if ve.ValueFrom.SecretKeyRef != nil {
					envsStr += fmt.Sprintf("%s:Secret:%s,", ve.Name, ve.ValueFrom.SecretKeyRef.Name)
				}
				if ve.ValueFrom.ConfigMapKeyRef != nil {
					envsStr += fmt.Sprintf("%s:Configmap:%s,", ve.Name, ve.ValueFrom.ConfigMapKeyRef.Name)
				}
			}

		}
		//log.Println(envsStr)
		//log.Println(len(v1.Env))
		//if len(v1.Env) > 0 {
		if len(envsStr) > 0 {
			envsStr = envsStr[0 : len(envsStr)-1]
		}

		for _, v3 := range v1.VolumeMounts {
			mountsStr += fmt.Sprintf("%s:%s:%v,", v3.Name, v3.MountPath, v3.ReadOnly)
		}
		if len(v1.VolumeMounts) > 0 {
			mountsStr = mountsStr[0 : len(mountsStr)-1]
		}

		for _, v4 := range v1.Ports {
			PortsStr += fmt.Sprintf("%d:%s:%v,", v4.ContainerPort, v4.Name, v4.Protocol)
		}
		if len(v1.Ports) > 0 {
			PortsStr = PortsStr[0 : len(PortsStr)-1]
		}

		var containerResUedstr string
		for _, vv1 := range ContainerResUsed {
			if vv1.ContainerName == v1.Name {
				containerResUedstr = fmt.Sprintf("CPU:%s,Memory:%s", vv1.Cpu, vv1.Mem)
				break
			}
		}

		xItems := &Container{
			ContainerName:  v1.Name,
			Envs:           envsStr,
			Mounts:         mountsStr,
			ContainerImage: v1.Image,
			PullPolicy:     fmt.Sprintf("%v", v1.ImagePullPolicy),
			Ports:          PortsStr,
			ResLimits:      fmt.Sprintf("CPU:%s,Memory:%s", v1.Resources.Limits.Cpu().String(), v1.Resources.Limits.Memory().String()),
			ResRequests:    fmt.Sprintf("CPU:%s,Memory:%s", v1.Resources.Requests.Cpu().String(), v1.Resources.Requests.Memory().String()),
			ResUsed:        containerResUedstr,
		}
		ccc = append(ccc, *xItems)
	}

	var ownerKind, ownerName string
	if pod.OwnerReferences != nil {
		ownerKind = pod.OwnerReferences[0].Kind
		ownerName = pod.OwnerReferences[0].Name
	}
	var lastterminatedStr string
	if len(pod.Status.ContainerStatuses) > 0 && pod.Status.ContainerStatuses[0].LastTerminationState.Terminated != nil {
		sTime := pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.StartedAt.String()
		eTime := pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.FinishedAt.String()
		reason := pod.Status.ContainerStatuses[0].LastTerminationState.Terminated.Reason
		lastterminatedStr = fmt.Sprintf("StartedAt:%s<br>FinishedAt:%s<br>reason:%s", sTime, eTime, reason)
	}

	return &PodDetails{
		PodName:        pod.Name,
		PodIp:          pod.Status.PodIP,
		NameSpace:      pod.ObjectMeta.Namespace,
		NodeName:       pod.Spec.NodeName,
		HostIp:         pod.Status.HostIP,
		PodPhase:       containerState,
		RestartCount:   restartNum,
		OwnerKind:      ownerKind,
		OwnerName:      ownerName,
		Labels:         labelsStr,
		Lastterminated: lastterminatedStr,
		Annotations:    annotationsStr,
		CreateTime:     pod.ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Conditions:     bbb,
		Volumes:        vvv,
		Containers:     ccc,
	}, err

}

// PodContainerList 获取Pod容器列表 / Get Pod container list
// kubeconfig: k8s集群配置 / k8s cluster configuration
// nameSpace: 命名空间 / namespace
// podName: Pod名称 / Pod name
func PodContainerList(kubeconfig, nameSpace, podName string) ([]Container, error) {
	var ccc = make([]Container, 0)
	pod, err := common.ClientSet(kubeconfig).CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return ccc, err
	}
	for _, v1 := range pod.Spec.Containers {
		xItems := &Container{
			ContainerName:  v1.Name,
			Envs:           "",
			Mounts:         "",
			ContainerImage: v1.Image,
			PullPolicy:     fmt.Sprintf("%v", v1.ImagePullPolicy),
			Ports:          "",
			ResLimits:      fmt.Sprintf("CPU:%s,Memory:%s", v1.Resources.Limits.Cpu().String(), v1.Resources.Limits.Memory().String()),
			ResRequests:    fmt.Sprintf("CPU:%s,Memory:%s", v1.Resources.Requests.Cpu().String(), v1.Resources.Requests.Memory().String()),
			ResUsed:        "",
		}
		ccc = append(ccc, *xItems)
	}
	return ccc, nil
}

// GetPodYaml 获取Pod的YAML定义 / Get Pod YAML definition
// kubeconfig: k8s集群配置 / k8s cluster configuration
// nameSpace: 命名空间 / namespace
// podName: Pod名称 / Pod name
func GetPodYaml(kubeconfig, nameSpace, podName string) (string, error) {
	pod, err := common.ClientSet(kubeconfig).CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	podUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(podUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// PodDel 删除Pod / Delete Pod
// kubeconfig: k8s集群配置 / k8s cluster configuration
// nameSpace: 命名空间 / namespace
// podName: Pod名称 / Pod name
func PodDel(kubeconfig, nameSpace, podName string) error {
	err := common.ClientSet(kubeconfig).CoreV1().Pods(nameSpace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// PodLog 获取Pod日志 / Get Pod logs
// kubeconfig: k8s集群配置 / k8s cluster configuration
// nameSpace: 命名空间 / namespace
// podName: Pod名称 / Pod name
// container: 容器名称 / container name
// logLine: 日志行数 / number of log lines
// encode: 编码格式 / encoding format
func PodLog(kubeconfig, nameSpace, podName, container string, logLine int64, encode string) string {
	clientset := common.ClientSet(kubeconfig)
	pod, err := clientset.CoreV1().Pods(nameSpace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get pod %q: %v\n", podName, err)
		//os.Exit(1)
		return "error"
	}

	var logOptions = &corev1.PodLogOptions{}
	logOptions.Follow = false //持续输出 / continuous output
	logOptions.Timestamps = false
	//var line int64 = 50
	if logLine > 0 {
		//line := logLine
		logOptions.TailLines = &logLine
	}

	var containerPt *string = &container
	if container == "" {
		if len(pod.Spec.Containers) == 0 {
			return "no containers found in pod"
		}
		containerPt = &pod.Spec.Containers[0].Name
	}
	logOptions.Container = *containerPt

	podLogs, err := clientset.CoreV1().Pods(nameSpace).GetLogs(podName, logOptions).Stream(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get logs for pod %q: %v\n", podName, err)
		//os.Exit(1)
		return "error"
	}
	defer podLogs.Close()

	//1单次输出 / 1 single output
	if encode == "utf8" {

		utf8Reader := simplifiedchinese.GB18030.NewDecoder().Reader(podLogs)
		//utf8Reader := transform.NewReader(bytes, simplifiedchinese.GBK.NewDecoder())
		body, err2 := io.ReadAll(utf8Reader)
		if err2 != nil {
			log.Printf("Failed to read pod logs: %v", err2)
			return ""
		}
		return string(body)
	} else {
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			log.Printf("Failed to copy pod logs: %v", err)
			return ""
		}
		return buf.String()
	}
}

func PodYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	pod := &corev1.Pod{}
	err = json.Unmarshal(data, pod)
	if err != nil {
		return err
	}
	namespace := pod.ObjectMeta.Namespace
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	return err
}
