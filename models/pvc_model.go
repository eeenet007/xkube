// storageclass_model.go
package models

import (
	"context"
	"fmt"

	//"time"
	"log"
	"strings"
	"xkube/common"

	//v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/util/retry"
	"sigs.k8s.io/yaml"
)

// PersistentVolumeClaim 定义了 PersistentVolumeClaim 资源的结构体，用于存储 PVC 的基本信息
type PersistentVolumeClaim struct {
	PvcName      string `json:"pvcName"`      // PVC 名称
	NameSpace    string `json:"nameSpace"`    // 所在命名空间
	Capacity     string `json:"capacity"`     // 容量大小
	AccessMode   string `json:"accessMode"`   // 访问模式
	StorageClass string `json:"storageClass"` // 存储类
	Status       string `json:"status"`       // 状态
	VolumeName   string `json:"volumeName"`   // 关联的 PV 名称
	Labels       string `json:"labels"`       // 标签
	Annotations  string `json:"annotations"`  // 注解
	CreateTime   string `json:"createTime"`   // 创建时间
}

// PersistentVolumeClaimList 获取指定命名空间下的 PersistentVolumeClaim 列表
// 参数:
//   - kubeconfig: Kubernetes 配置文件路径
//   - namespace: 命名空间，如果为空则查询所有命名空间
//   - pvcName: PVC 名称，用于过滤
//   - labelsKey: 标签键，用于标签过滤
//   - labelsValue: 标签值，用于标签过滤
//
// 返回值:
//   - []PersistentVolumeClaim: PVC 列表
//   - error: 错误信息
func PersistentVolumeClaimList(kubeconfig, namespace, pvcName, labelsKey, labelsValue string) ([]PersistentVolumeClaim, error) {
	// 创建 Kubernetes 客户端
	clientset := common.ClientSet(kubeconfig)
	if namespace == "" {
		// 如果命名空间为空，则查询所有命名空间
		namespace = corev1.NamespaceAll
	}

	// 设置 ListOptions，支持标签过滤
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	// 查询 PVC 列表
	pvcList, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}

	// 构建返回结果
	var bbb = make([]PersistentVolumeClaim, 0)
	for _, pvc := range pvcList.Items {
		// 根据名称过滤
		if pvcName != "" {
			if !strings.Contains(pvc.Name, pvcName) {
				continue
			}
		}

		// 构建标签字符串
		var labelsStr string
		for kk, vv := range pvc.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}

		// 构建注解字符串
		var annotationsStr string
		for k2, v2 := range pvc.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}

		// 构建访问模式字符串
		var accessModeStr string
		for _, v3 := range pvc.Spec.AccessModes {
			accessModeStr += string(v3)
		}

		// 获取存储类名称
		var storageClass string
		if pvc.Spec.StorageClassName != nil {
			storageClass = *pvc.Spec.StorageClassName
		}

		// 构建 PersistentVolumeClaim 对象
		Items := &PersistentVolumeClaim{
			PvcName:      pvc.Name,
			NameSpace:    pvc.Namespace,
			Capacity:     pvc.Status.Capacity.Storage().String(),
			AccessMode:   accessModeStr,
			StorageClass: storageClass,
			Status:       string(pvc.Status.Phase),
			VolumeName:   pvc.Spec.VolumeName,
			Labels:       labelsStr,
			Annotations:  annotationsStr,
			CreateTime:   pvc.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// PersistentVolumeClaimDetail 获取指定 PersistentVolumeClaim 的详细信息
// 参数:
//   - kubeconfig: Kubernetes 配置文件路径
//   - namespace: PVC 所在命名空间
//   - pvName: PVC 名称
//
// 返回值:
//   - *PersistentVolumeClaim: PVC 详细信息
//   - error: 错误信息
func PersistentVolumeClaimDetail(kubeconfig, namespace, pvName string) (*PersistentVolumeClaim, error) {
	// 创建 PVC 客户端
	pvClient := common.ClientSet(kubeconfig).CoreV1().PersistentVolumeClaims(namespace)
	// 获取 PVC 详细信息
	pvc, err := pvClient.Get(context.TODO(), pvName, metav1.GetOptions{})

	// 构建标签字符串
	var labelsStr string
	for kk, vv := range pvc.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}

	// 构建注解字符串
	var annotationsStr string
	for k2, v2 := range pvc.Annotations {
		annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
	}

	// 构建访问模式字符串
	var accessModeStr string
	for _, v3 := range pvc.Spec.AccessModes {
		accessModeStr += string(v3)
	}

	// 获取存储类名称
	var storageClass string
	if pvc.Spec.StorageClassName != nil {
		storageClass = *pvc.Spec.StorageClassName
	}

	// 返回 PVC 详细信息
	return &PersistentVolumeClaim{
		PvcName:      pvc.Name,
		NameSpace:    pvc.Namespace,
		Capacity:     pvc.Status.Capacity.Storage().String(),
		AccessMode:   accessModeStr,
		StorageClass: storageClass,
		Status:       string(pvc.Status.Phase),
		VolumeName:   pvc.Spec.VolumeName,
		Labels:       labelsStr,
		Annotations:  annotationsStr,
		CreateTime:   pvc.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}, err
}

// GetPersistentVolumeClaimYaml 获取指定 PersistentVolumeClaim 的 YAML 格式配置
// 参数:
//   - kubeconfig: Kubernetes 配置文件路径
//   - namespace: PVC 所在命名空间
//   - pvName: PVC 名称
//
// 返回值:
//   - string: PVC 的 YAML 格式配置
//   - error: 错误信息
func GetPersistentVolumeClaimYaml(kubeconfig, namespace, pvName string) (string, error) {
	// 创建 PVC 客户端
	pvClient := common.ClientSet(kubeconfig).CoreV1().PersistentVolumeClaims(namespace)
	// 获取 PVC 对象
	pvc, err := pvClient.Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// 转换为非结构化对象
	pvcUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pvc)
	if err != nil {
		return "", err
	}
	// 转换为 YAML 格式
	yamlBytes, err := yaml.Marshal(pvcUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil

}
