// gatewayclass_model.go
package models

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"xkube/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayClient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

type GatewayClass struct {
	GatewayClassName string `json:"gatewayClassName"`
	ControllerName   string `json:"controllerName"`
	Accepted         string `json:"accepted"`
	Description      string `json:"description"`
	Age              string `json:"age"`
	CreateTime       string `json:"createTime"`
}

type GatewayClassDetail struct {
	GatewayClass
	Conditions []GatewayClassCondition `json:"conditions"`
}

type GatewayClassCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

// GatewayClassList 获取GatewayClass列表
func GatewayClassList(kubeconfig string) ([]GatewayClass, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	gatewayClassList, err := clientset.GatewayV1().GatewayClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list gateway classes: %v\n", err)
		return nil, err
	}

	var gatewayClasses []GatewayClass
	for _, gc := range gatewayClassList.Items {
		accepted := "Unknown"
		// 检查是否有 Accepted 状态
		for _, condition := range gc.Status.Conditions {
			if condition.Type == string(gatewayv1.GatewayClassConditionStatusAccepted) {
				accepted = string(condition.Status)
				break
			}
		}

		description := ""
		if gc.Spec.Description != nil {
			description = *gc.Spec.Description
		}

		// 计算 Age
		age := ""
		if !gc.CreationTimestamp.IsZero() {
			duration := time.Since(gc.CreationTimestamp.Time)
			age = formatDuration(duration)
		}

		gatewayClass := GatewayClass{
			GatewayClassName: gc.Name,
			ControllerName:   string(gc.Spec.ControllerName),
			Accepted:         accepted,
			Description:      description,
			Age:              age,
			CreateTime:       gc.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		gatewayClasses = append(gatewayClasses, gatewayClass)
	}

	return gatewayClasses, nil
}

// GetGatewayClassDetail 获取GatewayClass详细信息
func GetGatewayClassDetail(kubeconfig, gatewayClassName string) (*GatewayClassDetail, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	gc, err := clientset.GatewayV1().GatewayClasses().Get(context.TODO(), gatewayClassName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get gateway class %s: %v\n", gatewayClassName, err)
		return nil, err
	}

	accepted := "Unknown"
	// 检查是否有 Accepted 状态
	for _, condition := range gc.Status.Conditions {
		if condition.Type == string(gatewayv1.GatewayClassConditionStatusAccepted) {
			accepted = string(condition.Status)
			break
		}
	}

	description := ""
	if gc.Spec.Description != nil {
		description = *gc.Spec.Description
	}

	// 计算 Age
	age := ""
	if !gc.CreationTimestamp.IsZero() {
		duration := time.Since(gc.CreationTimestamp.Time)
		age = formatDuration(duration)
	}

	gatewayClass := GatewayClass{
		GatewayClassName: gc.Name,
		ControllerName:   string(gc.Spec.ControllerName),
		Accepted:         accepted,
		Description:      description,
		Age:              age,
		CreateTime:       gc.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}

	// 处理 conditions
	var conditions []GatewayClassCondition
	for _, cond := range gc.Status.Conditions {
		condition := GatewayClassCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastTransitionTime: cond.LastTransitionTime.Format("2006-01-02 15:04:05"),
		}
		conditions = append(conditions, condition)
	}

	detail := &GatewayClassDetail{
		GatewayClass: gatewayClass,
		Conditions:   conditions,
	}

	return detail, nil
}

// GetGatewayClassYaml 获取GatewayClass的YAML配置
func GetGatewayClassYaml(kubeconfig, gatewayClassName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return "", err
	}

	gc, err := clientset.GatewayV1().GatewayClasses().Get(context.TODO(), gatewayClassName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get gateway class %s: %v\n", gatewayClassName, err)
		return "", err
	}

	gcUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(gc)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(gcUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	var parts []string

	// 天
	days := int(d.Hours()) / 24
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
		d -= time.Duration(days*24) * time.Hour
	}

	// 小时
	hours := int(d.Hours())
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
		d -= time.Duration(hours) * time.Hour
	}

	// 分钟
	minutes := int(d.Minutes())
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
		d -= time.Duration(minutes) * time.Minute
	}

	// 秒
	seconds := int(d.Seconds())
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	// 最多显示两个单位
	if len(parts) > 2 {
		parts = parts[:2]
	}

	return strings.Join(parts, "")
}
