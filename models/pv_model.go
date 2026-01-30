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

// PersistentVolume 定义了 PersistentVolume 资源的结构体，用于存储 PV 的基本信息
type PersistentVolume struct {
	PvName          string `json:"pvName"`          // PV 名称
	Capacity        string `json:"capacity"`        // 容量大小
	AccessMode      string `json:"accessMode"`      // 访问模式
	PvReclaimPolicy string `json:"pvReclaimPolicy"` // 回收策略
	StorageClass    string `json:"storageClass"`    // 存储类
	Status          string `json:"status"`          // 状态
	ClaimRef        string `json:"claimRef"`        // 绑定的 PVC 引用
	Labels          string `json:"labels"`          // 标签
	Annotations     string `json:"annotations"`     // 注解
	CreateTime      string `json:"createTime"`      // 创建时间
}

// PersistentVolumeList 获取所有 PersistentVolume 的列表
// 参数:
//   - kubeconfig: Kubernetes 配置文件路径
//
// 返回值:
//   - []PersistentVolume: PV 列表
//   - error: 错误信息
func PersistentVolumeList(kubeconfig string) ([]PersistentVolume, error) {
	// 构建返回结果
	var bbb = make([]PersistentVolume, 0)
	// 创建 Kubernetes 客户端
	clientset := common.ClientSet(kubeconfig)
	// 查询所有 PV 列表
	pvList, err := clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
		return bbb, err
	}

	for _, pv := range pvList.Items {
		// 构建标签字符串
		var labelsStr string
		for kk, vv := range pv.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}

		// 构建注解字符串
		var annotationsStr string
		for k2, v2 := range pv.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}

		// 构建访问模式字符串
		var accessModeStr string
		for _, v3 := range pv.Spec.AccessModes {
			accessModeStr += string(v3)
		}

		// 构建 PersistentVolume 对象
		Items := &PersistentVolume{
			PvName:          pv.Name,
			Capacity:        pv.Spec.Capacity.Storage().String(),
			AccessMode:      accessModeStr,
			PvReclaimPolicy: string(pv.Spec.PersistentVolumeReclaimPolicy),
			StorageClass:    pv.Spec.StorageClassName,
			Status:          string(pv.Status.Phase),
			ClaimRef:        fmt.Sprintf("命名空间:%s<br>名称:%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name),
			Labels:          labelsStr,
			Annotations:     annotationsStr,
			CreateTime:      pv.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, nil
}

// PersistentVolumeDetail 获取指定 PersistentVolume 的详细信息
// 参数:
//   - kubeconfig: Kubernetes 配置文件路径
//   - pvName: PV 名称
//
// 返回值:
//   - *PersistentVolume: PV 详细信息
//   - error: 错误信息
func PersistentVolumeDetail(kubeconfig, pvName string) (*PersistentVolume, error) {
	// 创建 PV 客户端
	pvClient := common.ClientSet(kubeconfig).CoreV1().PersistentVolumes()
	// 获取指定 PV 的详细信息
	pv, err := pvClient.Get(context.TODO(), pvName, metav1.GetOptions{})

	// 构建标签字符串
	var labelsStr string
	for kk, vv := range pv.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}

	// 构建注解字符串
	var annotationsStr string
	for k2, v2 := range pv.Annotations {
		annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
	}

	// 构建访问模式字符串
	var accessModeStr string
	for _, v3 := range pv.Spec.AccessModes {
		accessModeStr += string(v3)
	}

	// 返回 PV 详细信息
	return &PersistentVolume{
		PvName:          pv.Name,
		Capacity:        pv.Spec.Capacity.Storage().String(),
		AccessMode:      accessModeStr,
		PvReclaimPolicy: string(pv.Spec.PersistentVolumeReclaimPolicy),
		StorageClass:    pv.Spec.StorageClassName,
		Status:          string(pv.Status.Phase),
		ClaimRef:        fmt.Sprintf("命名空间:%s<br>名称:%s", pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name),
		Labels:          labelsStr,
		Annotations:     annotationsStr,
		CreateTime:      pv.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}, err
}

// GetPersistentVolumeYaml 获取指定 PersistentVolume 的 YAML 格式配置
// 参数:
//   - kubeconfig: Kubernetes 配置文件路径
//   - pvName: PV 名称
//
// 返回值:
//   - string: PV 的 YAML 格式配置
//   - error: 错误信息
func GetPersistentVolumeYaml(kubeconfig, pvName string) (string, error) {
	// 创建 PV 客户端
	pvClient := common.ClientSet(kubeconfig).CoreV1().PersistentVolumes()
	// 获取 PV 对象
	persistentVolume, err := pvClient.Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// 转换为非结构化对象
	pvUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(persistentVolume)
	if err != nil {
		return "", err
	}
	// 转换为 YAML 格式
	yamlBytes, err := yaml.Marshal(pvUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
