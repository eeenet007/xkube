// storageclass_model.go
package models

import (
	"context"
	"fmt"

	//"time"
	"log"
	"xkube/common"

	//v1 "k8s.io/api/core/v1"
	//corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/util/retry"
	"sigs.k8s.io/yaml"
)

// Storageclass - 存储类模型结构体(Storageclass - Storage class model struct)
type Storageclass struct {
	StorageclassName string `json:"storageclassName"` // 存储类名称(Storage class name)
	Provisioner      string `json:"provisioner"`      // 提供者(Provisioner)
	Parameters       string `json:"parameters"`       // 参数(Parameters)
	ReclaimPolicy    string `json:"reclaimPolicy"`    // 回收策略(Reclaim policy)
	Labels           string `json:"labels"`           // 标签(Labels)
	Annotations      string `json:"annotations"`      // 注解(Annotations)
	CreateTime       string `json:"createTime"`       // 创建时间(Creation time)
}

// StorageclassList - 获取存储类列表(StorageclassList - Get storage class list)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// 返回值：存储类列表和错误信息(Return: storage class list and error information)
func StorageclassList(kubeconfig string) ([]Storageclass, error) {
	clientset := common.ClientSet(kubeconfig)                                                                  // 获取Kubernetes客户端(Get Kubernetes client)
	storageclassList, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{}) // 列出所有存储类(List all storage classes)
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}
	var bbb = make([]Storageclass, 0)
	// 遍历存储类列表(Iterate through storage class list)
	for _, sc := range storageclassList.Items {
		var parms string
		// 处理参数(Processing parameters)
		for kk, vv := range sc.Parameters {
			parms += fmt.Sprintf("%s:%s,", kk, vv)
		}

		var labelsStr string
		// 处理标签(Processing labels)
		for kk, vv := range sc.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		var annotationsStr string
		// 处理注解(Processing annotations)
		for k2, v2 := range sc.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}

		// 构造存储类对象(Construct storage class object)
		Items := &Storageclass{
			StorageclassName: sc.Name,                                            // 存储类名称(Storage class name)
			Provisioner:      sc.Provisioner,                                     // 提供者(Provisioner)
			Parameters:       parms,                                              // 参数(Parameters)
			ReclaimPolicy:    string(*sc.ReclaimPolicy),                          // 回收策略(Reclaim policy)
			Labels:           labelsStr,                                          // 标签(Labels)
			Annotations:      annotationsStr,                                     // 注解(Annotations)
			CreateTime:       sc.CreationTimestamp.Format("2006-01-02 15:04:05"), // 创建时间(Creation time)
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// StorageclassDetail - 获取存储类详细信息(StorageclassDetail - Get storage class detailed information)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// storageclassName - 存储类名称(Storage class name)
// 返回值：存储类详细信息和错误信息(Return: storage class detailed information and error information)
func StorageclassDetail(kubeconfig, storageclassName string) (*Storageclass, error) {
	scClient := common.ClientSet(kubeconfig).StorageV1().StorageClasses()          // 获取存储类客户端(Storage class client)
	sc, err := scClient.Get(context.TODO(), storageclassName, metav1.GetOptions{}) // 获取指定存储类(Get specified storage class)
	var parms string
	// 处理参数(Processing parameters)
	for kk, vv := range sc.Parameters {
		parms += fmt.Sprintf("%s:%s,", kk, vv)
	}
	var labelsStr string
	// 处理标签(Processing labels)
	for kk, vv := range sc.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	var annotationsStr string
	// 处理注解(Processing annotations)
	for k2, v2 := range sc.Annotations {
		annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
	}
	// 返回存储类详细信息(Return storage class detailed information)
	return &Storageclass{
		StorageclassName: sc.Name,                                            // 存储类名称(Storage class name)
		Provisioner:      sc.Provisioner,                                     // 提供者(Provisioner)
		Parameters:       parms,                                              // 参数(Parameters)
		ReclaimPolicy:    string(*sc.ReclaimPolicy),                          // 回收策略(Reclaim policy)
		Labels:           labelsStr,                                          // 标签(Labels)
		Annotations:      annotationsStr,                                     // 注解(Annotations)
		CreateTime:       sc.CreationTimestamp.Format("2006-01-02 15:04:05"), // 创建时间(Creation time)
	}, err
}

// GetStorageclassYaml - 获取存储类YAML定义(GetStorageclassYaml - Get storage class YAML definition)
// kubeconfig - Kubernetes配置信息(Kubernetes configuration information)
// storageclassName - 存储类名称(Storage class name)
// 返回值：存储类YAML字符串和错误信息(Return: storage class YAML string and error information)
func GetStorageclassYaml(kubeconfig, storageclassName string) (string, error) {
	scClient := common.ClientSet(kubeconfig).StorageV1().StorageClasses()                    // 获取存储类客户端(Storage class client)
	storageclass, err := scClient.Get(context.TODO(), storageclassName, metav1.GetOptions{}) // 获取指定存储类(Get specified storage class)
	if err != nil {
		return "", err
	}
	scUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(storageclass) // 转换为非结构化对象(Convert to unstructured object)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(scUnstructured) // 将对象序列化为YAML(serialize object to YAML)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil // 返回YAML字符串(Return YAML string)
}
