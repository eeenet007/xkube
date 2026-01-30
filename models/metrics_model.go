// metrics_model.go
// Package models provides data models and related functions for Kubernetes metrics
// 包 models 提供 Kubernetes 指标的数据模型和相关函数
package models

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"time"
	"xkube/common"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NodeMetirc represents metrics data for a Kubernetes node
// NodeMetirc 表示 Kubernetes 节点的指标数据
type NodeMetirc struct {
	Timestamp string  `json:"timestamp"` // Metric timestamp 指标时间戳
	NodeName  string  `json:"nodeName"`  // Node name 节点名称
	Cpu       float64 `json:"cpu"`       // CPU usage CPU使用量
	Mem       float64 `json:"mem"`       // Memory usage 内存使用量
}

// PodMetirc represents metrics data for a Kubernetes pod
// PodMetirc 表示 Kubernetes Pod 的指标数据
// type PodMetirc2 struct {
// 	Timestamp string `json:"timestamp"` // Metric timestamp 指标时间戳
// 	PodName   string `json:"podName"`   // Pod name Pod名称
// 	Cpu       string `json:"cpu"`       // CPU usage CPU使用量
// 	Mem       string `json:"mem"`       // Memory usage 内存使用量
// 	NameSpace string `json:"nameSpace"` //namespace 命名空间
// }

type PodMetirc struct {
	Timestamp string  `json:"timestamp"` // Metric timestamp 指标时间戳
	PodName   string  `json:"podName"`   // Pod name Pod名称
	Cpu       float64 `json:"cpu"`       // CPU usage CPU使用量
	Mem       float64 `json:"mem"`       // Memory usage 内存使用量
	NameSpace string  `json:"nameSpace"` //namespace 命名空间
}

// ContainerMetirc represents metrics data for a container
// ContainerMetirc 表示容器的指标数据
type ContainerMetirc struct {
	Timestamp     string `json:"timestamp"`     // Metric timestamp 指标时间戳
	ContainerName string `json:"containerName"` // Container name 容器名称
	Cpu           string `json:"cpu"`           // CPU usage CPU使用量
	Mem           string `json:"mem"`           // Memory usage 内存使用量
}

// GetNodeMetricList retrieves metrics for all nodes in the Kubernetes cluster
// GetNodeMetricList 获取 Kubernetes 集群中所有节点的指标
func GetNodeMetricList(kubeconfig string) ([]NodeMetirc, error) {
	config := common.ClientConfig(kubeconfig)
	var bbb = make([]NodeMetirc, 0)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Printf("[ERROR] GetNodeMetricList err:%s\n", err)
		return bbb, err
	}
	nodeMetricsList, err := metricsClient.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("[ERROR] NodeMetricList err:%s\n", err)
		return bbb, err
	}
	for _, vv := range nodeMetricsList.Items {
		bbb = append(bbb, NodeMetirc{
			Timestamp: vv.CreationTimestamp.Format("15:04:05"), //2006-01-02 15:04:05
			NodeName:  vv.Name,
			//Cpu:       fmt.Sprintf("%.3f", float64(vv.Usage.Cpu().MilliValue())/1000),
			//Mem:       fmt.Sprintf("%.3fGi", float64(vv.Usage.Memory().Value())/1024/1024/1024),
			Cpu: math.Round(float64(vv.Usage.Cpu().MilliValue())/1000*1000) / 1000,
			Mem: math.Round(float64(vv.Usage.Memory().Value())/1024/1024/1024*1000) / 1000,
		})
	}
	return bbb, nil
}

// GetPodMetricList retrieves metrics for all pods in a specific namespace
// GetPodMetricList 获取指定命名空间中所有 Pod 的指标
func GetPodMetricList(kubeconfig, nameSpace, sorts string) ([]PodMetirc, error) {
	config := common.ClientConfig(kubeconfig)
	if nameSpace == "" {
		nameSpace = corev1.NamespaceAll
	}

	var bbb = make([]PodMetirc, 0)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Printf("[ERROR] GetPodMetricList err:%s\n", err)
		return bbb, err
	}
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(nameSpace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("[ERROR] PodMetricList err:%s\n", err)
		return bbb, err
	}
	for _, vv := range podMetricsList.Items {
		var cpuUsage, memoryUsage int64
		for _, v1 := range vv.Containers {
			cpud := v1.Usage.Cpu().MilliValue()
			memd := v1.Usage.Memory().Value()
			cpuUsage += cpud
			memoryUsage += memd
		}

		bbb = append(bbb, PodMetirc{
			Timestamp: vv.CreationTimestamp.Format("15:04:05"), //2006-01-02 15:04:05
			PodName:   vv.Name,
			Cpu:       math.Round(float64(cpuUsage)/1000*1000) / 1000,
			Mem:       math.Round(float64(memoryUsage)/1024/1024*1000) / 1000,
			NameSpace: vv.Namespace,
		})
	}
	if sorts == "cpu" {
		sort.Slice(bbb, func(i, j int) bool {
			return bbb[i].Cpu > bbb[j].Cpu // 降序排序
		})
	} else {
		// 使用sort.Slice进行降序排序
		sort.Slice(bbb, func(i, j int) bool {
			return bbb[i].Mem > bbb[j].Mem // 降序排序
		})
	}

	return bbb, nil
}

// GetPodMetric retrieves metrics for a specific pod
// GetPodMetric 获取指定 Pod 的指标
func GetPodMetric(kubeconfig, namespace, podName string) string {
	config := common.ClientConfig(kubeconfig)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Println(err)
	}
	metrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[WARN] GetPodMetric Fail:%s\n", err)
		return fmt.Sprintf(`{"cpu":0.0,"mem":0.0,"time":"%s"}`, time.Now().Format("15:04:05"))
	}
	var cpuUsage, memoryUsage int64
	for _, vv := range metrics.Containers {
		cpud := vv.Usage.Cpu().MilliValue()
		memd := vv.Usage.Memory().Value()
		cpuUsage += cpud
		memoryUsage += memd
	}
	//metrics.CreationTimestamp.Unix()
	return fmt.Sprintf(`{"cpu":%.3f,"mem":%.3f,"time":"%s"}`, float64(cpuUsage), float64(memoryUsage)/1024/1024, time.Now().Format("15:04:05"))
}

// GetContainerMetric retrieves metrics for containers in a specific pod
// GetContainerMetric 获取指定 Pod 中容器的指标
func GetContainerMetric(kubeconfig, namespace, podName string) ([]ContainerMetirc, error) {
	var bbb = make([]ContainerMetirc, 0)
	config := common.ClientConfig(kubeconfig)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return bbb, err
	}
	metrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] GetContainerMetric Fail:%s\n", err)
		return bbb, err
		//return fmt.Sprintf(`{"cpu":0.0,"mem":0.0,"time":"%s"}`, time.Now().Format("15:04:05"))
	}
	//var cpuUsage, memoryUsage int64
	for _, vv := range metrics.Containers {
		cpud := vv.Usage.Cpu().MilliValue()
		memd := vv.Usage.Memory().Value()
		dt := time.Now().Format("15:04:05")
		bbb = append(bbb, ContainerMetirc{
			Timestamp:     dt,
			ContainerName: vv.Name,
			Cpu:           fmt.Sprintf("%.3f", float64(cpud)/1000),
			Mem:           fmt.Sprintf("%.3fMi", float64(memd)/1024/1024),
		})
	}
	return bbb, nil
}

// GetNodeMetric retrieves metrics for a specific node
// GetNodeMetric 获取指定节点的指标
func GetNodeMetric(kubeconfig, nodeName string) string {
	config := common.ClientConfig(kubeconfig)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Println(err)
	}
	metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] GetNodeMetric Fail:%s\n", err)
		return fmt.Sprintf(`{"cpu":0.0,"mem":0.0,"time":"%s"}`, time.Now().Format("15:04:05"))
	}
	//return fmt.Sprintf(`{"cpu":%.3f,"mem":%.3f,"time":"%s"}`, float64(cpuUsage), float64(memoryUsage)/1024/1024, time.Now().Format("15:04:05"))
	return fmt.Sprintf(`{"cpu":%.3f,"mem":%.3f,"time":"%s"}`, float64(metrics.Usage.Cpu().MilliValue())/1000, float64(metrics.Usage.Memory().Value())/1024/1024/1024, time.Now().Format("15:04:05"))
}

// GetNodeMetricV3 retrieves raw CPU and memory metrics for a specific node
// GetNodeMetricV3 获取指定节点的原始 CPU 和内存指标
func GetNodeMetricV3(kubeconfig, nodeName string) (int64, int64, error) {
	config := common.ClientConfig(kubeconfig)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	metrics, err := metricsClient.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] GetNodeMetric Fail:%s\n", err)
		return 0, 0, err
	}
	//return fmt.Sprintf(`{"cpu":%.3f,"mem":%.3f,"time":"%s"}`, float64(cpuUsage), float64(memoryUsage)/1024/1024, time.Now().Format("15:04:05"))
	return metrics.Usage.Cpu().MilliValue(), metrics.Usage.Memory().Value(), nil
}

// GetPodMetricV2 retrieves raw CPU and memory metrics for a specific pod
// GetPodMetricV2 获取指定 Pod 的原始 CPU 和内存指标
func GetPodMetricV2(kubeconfig, namespace, podName string) (int64, int64, error) {
	config := common.ClientConfig(kubeconfig)
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	metrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] GetPodMetric Fail:%s\n", err)
		return 0, 0, err
	}
	var cpuUsage, memoryUsage int64
	for _, vv := range metrics.Containers {
		cpud := vv.Usage.Cpu().MilliValue()
		memd := vv.Usage.Memory().Value()
		cpuUsage += cpud
		memoryUsage += memd
	}
	return cpuUsage, memoryUsage, nil
}
