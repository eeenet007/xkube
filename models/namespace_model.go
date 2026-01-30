// namespace_model.go
// Package models provides data models and related functions for Kubernetes resources
// 包 models 提供 Kubernetes 资源的数据模型和相关函数
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"xkube/common"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// NameSpace represents a Kubernetes namespace with its properties
// NameSpace 表示一个包含其属性的 Kubernetes 命名空间
type NameSpace struct {
	NameSpace  string `json:"nameSpace"`  // Namespace name 命名空间名称
	Status     string `json:"status"`     // Namespace status 状态
	Labels     string `json:"labels"`     // Namespace labels 标签
	CreateTime string `json:"createTime"` // Creation time 创建时间
}

// LimitRes represents resource limits for a namespace
// LimitRes 表示命名空间的资源限制
type LimitRes struct {
	LimitCpu   string `json:"limitCpu"`   // CPU limit CPU限制
	LimitMem   string `json:"limitMem"`   // Memory limit 内存限制
	RequestMem string `json:"requestMem"` // Memory request 内存请求
	RequestCpu string `json:"requestCpu"` // CPU request CPU请求
}

// NsList retrieves a list of namespaces from the Kubernetes cluster
// NsList 从 Kubernetes 集群中获取命名空间列表
// func NsList(kubeconfig string) ([]corev1.Namespace, error) {
func NsList(kubeconfig string) ([]NameSpace, error) {
	// Generate a clientset object through the kubeconfig cluster authentication file
	// 通过 kubeconfig 集群认证文件生成一个客户端操作对象 clientset
	clientset := common.ClientSet(kubeconfig)

	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list pods error:%v\n", err)
	}
	var NsArry = make([]NameSpace, 0)
	//fmt.Println("dev node count:", len(namespaceList.Items))
	for _, ns := range namespaceList.Items {
		var labelsStr string
		for kk, vv := range ns.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		NsArry = append(NsArry, NameSpace{
			NameSpace:  ns.Name,
			Status:     fmt.Sprintf("%v", ns.Status.Phase),
			Labels:     labelsStr,
			CreateTime: ns.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	//return namespaceList.Items, err
	return NsArry, err
}

// NsDetail retrieves detailed information about a specific namespace
// NsDetail 获取指定命名空间的详细信息
func NsDetail(kubeconfig, namespaceName string) (*corev1.Namespace, error) {
	clientset := common.ClientSet(kubeconfig)
	namespaceInfo, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		return namespaceInfo, err
	}
	return namespaceInfo, nil
}

// NsYaml retrieves the YAML representation of a specific namespace
// NsYaml 获取指定命名空间的 YAML 表示
func NsYaml(kubeconfig, namespaceName string) (string, error) {
	clientset := common.ClientSet(kubeconfig)
	namespaceInfo, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	nsUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(namespaceInfo)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(nsUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// YamlUpdate updates a namespace based on YAML data
// YamlUpdate 根据 YAML 数据更新命名空间
func NsYamlUpdate(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	ns := &corev1.Namespace{}
	err = json.Unmarshal(data, ns)
	if err != nil {
		return err
	}

	//nsName := ns.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	newns, err := clientset.CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(newns)
	return err
}

// CreateNs creates a new namespace with specified labels
// CreateNs 创建一个带有指定标签的新命名空间
func CreateNs(kubeconfig, namespaceName string, labesMap map[string]string) error {
	clientset := common.ClientSet(kubeconfig)
	var namespace corev1.Namespace
	namespace.Name = namespaceName
	namespace.ObjectMeta.Labels = labesMap
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DeleteNs deletes a namespace and its associated LimitRanges
// DeleteNs 删除命名空间及其关联的 LimitRanges
func DeleteNs(kubeconfig, namespaceName, force string) error {
	clientset := common.ClientSet(kubeconfig)

	if force != "1" {
		podlist, _ := PodListV2(kubeconfig, namespaceName, "", "", "", "", "", "")
		svclist, _ := SvcList(kubeconfig, namespaceName, "", "", "")
		cmlist, _ := CmList(kubeconfig, namespaceName, "", "", "")
		if len(podlist) > 0 || len(svclist) > 0 || len(cmlist) > 0 {
			return fmt.Errorf("[WARN]There are still resources in this namespace, please delete them first")
		}
	}

	err := clientset.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	// Delete LimitRanges
	// 删除 LimitRanges
	err = clientset.CoreV1().LimitRanges(namespaceName).Delete(context.TODO(), "limits", metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// CreateLimitRange creates resource limits for a namespace
// CreateLimitRange 为命名空间创建资源限制
func CreateLimitRange(kubeconfig, namespaceName string, lrs LimitRes) error {
	memoryQuantity, err := resource.ParseQuantity(lrs.LimitMem)
	if err != nil {
		return err
	}
	cpuQuantity, err := resource.ParseQuantity(lrs.LimitCpu)
	if err != nil {
		return err
	}
	requestMemoryQuantity, err := resource.ParseQuantity(lrs.RequestMem)
	if err != nil {
		return err
	}
	requestCPUQuantity, err := resource.ParseQuantity(lrs.RequestCpu)
	if err != nil {
		return err
	}

	// Create LimitRanges object
	// 创建 LimitRanges 对象
	limitRanges := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "limits",
			Namespace: namespaceName,
		},
		Spec: corev1.LimitRangeSpec{
			Limits: []corev1.LimitRangeItem{
				{
					Type: corev1.LimitTypeContainer,
					Max: corev1.ResourceList{
						corev1.ResourceMemory: memoryQuantity,
						corev1.ResourceCPU:    cpuQuantity,
					},
					Default: corev1.ResourceList{
						corev1.ResourceMemory: memoryQuantity,
						corev1.ResourceCPU:    cpuQuantity,
					},
					DefaultRequest: corev1.ResourceList{
						corev1.ResourceMemory: requestMemoryQuantity,
						corev1.ResourceCPU:    requestCPUQuantity,
					},
				},
			},
		},
	}
	_, err = common.ClientSet(kubeconfig).CoreV1().LimitRanges(namespaceName).Create(context.TODO(), limitRanges, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// UpdateLimitRange updates resource limits for a namespace
// UpdateLimitRange 更新命名空间的资源限制
func UpdateLimitRange(kubeconfig, namespaceName string, lrs LimitRes) error {

	//log.Println(lrs)
	lr, err := common.ClientSet(kubeconfig).CoreV1().LimitRanges(namespaceName).Get(context.TODO(), "limits", metav1.GetOptions{})
	if err != nil {
		return err
	}

	lr.Spec.Limits = []corev1.LimitRangeItem{
		{
			Type: corev1.LimitTypeContainer,
			Max: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(lrs.LimitMem),
				corev1.ResourceCPU:    resource.MustParse(lrs.LimitCpu),
			},
			Default: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(lrs.LimitMem),
				corev1.ResourceCPU:    resource.MustParse(lrs.LimitCpu),
			},
			DefaultRequest: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(lrs.RequestMem),
				corev1.ResourceCPU:    resource.MustParse(lrs.RequestCpu),
			},
		},
	}

	//log.Printf("%d,%d,%d,%d", limitMem.Value(), limitCpu.MilliValue(), requestMem.Value(), requestCpu.MilliValue())
	_, err = common.ClientSet(kubeconfig).CoreV1().LimitRanges(namespaceName).Update(context.TODO(), lr, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetLimitRange retrieves resource limits for a namespace
// GetLimitRange 获取命名空间的资源限制
func GetLimitRange(kubeconfig, namespaceName string) (*LimitRes, error) {
	lr, err := common.ClientSet(kubeconfig).CoreV1().LimitRanges(namespaceName).Get(context.TODO(), "limits", metav1.GetOptions{})
	if err != nil {
		return &LimitRes{}, err
	}
	return &LimitRes{
		LimitMem:   lr.Spec.Limits[0].Default.Memory().String(),
		LimitCpu:   lr.Spec.Limits[0].Default.Cpu().String(),
		RequestMem: lr.Spec.Limits[0].DefaultRequest.Memory().String(),
		RequestCpu: lr.Spec.Limits[0].DefaultRequest.Cpu().String(),
	}, nil
}

// CreateNsByExist checks if a namespace exists and creates it if it doesn't
// CreateNsByExist 检查命名空间是否存在，如果不存在则创建
func CreateNsByExist(kubeconfig, namespaceName string) error {
	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err == nil {
		// Namespace already exists
		// 命名空间已存在
		return nil
	}
	var namespace corev1.Namespace
	namespace.Name = namespaceName
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &namespace, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
