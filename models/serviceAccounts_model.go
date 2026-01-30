// serviceAccounts_model.go
package models

import (
	"context"
	"fmt"

	//"time"
	"log"
	"xkube/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type ServiceAccounts struct {
	SaName           string `json:"saName"`
	NameSpace        string `json:"nameSpace"`
	Labels           string `json:"labels"`
	Annotations      string `json:"annotations"` //注解
	ImagePullSecrets string `json:"imagePullSecrets"`
	CreateTime       string `json:"createTime"`
}

// ServiceAccountsList - 获取服务账户列表 (Get service accounts list)
// kubeconfig - Kubernetes配置信息 (Kubernetes configuration information)
// namespace - 命名空间 (Namespace)
// 返回值: 服务账户列表和错误信息 (Return: service accounts list and error information)
func ServiceAccountsList(kubeconfig, namespace string) ([]ServiceAccounts, error) {
	xList, err := common.ClientSet(kubeconfig).CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list serviceAccounts error, err:%v\n", err)
	}
	var bbb = make([]ServiceAccounts, 0)
	for _, sa := range xList.Items {
		var labelsStr string
		for kk, vv := range sa.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		var annotationsStr string
		for k2, v2 := range sa.ObjectMeta.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}
		if len(annotationsStr) > 0 {
			annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
		}

		var imagePullSecrets string
		if sa.ImagePullSecrets != nil {
			for _, v3 := range sa.ImagePullSecrets {
				imagePullSecrets += fmt.Sprintf("%s,", v3.Name)
			}
		}
		if len(imagePullSecrets) > 0 {
			imagePullSecrets = imagePullSecrets[0 : len(imagePullSecrets)-1]
		}

		Items := &ServiceAccounts{
			SaName:           sa.Name,
			NameSpace:        sa.Namespace,
			Labels:           labelsStr,
			Annotations:      annotationsStr,
			ImagePullSecrets: imagePullSecrets,
			CreateTime:       sa.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// GetServiceAccountsYaml - 获取服务账户YAML定义 (Get service account YAML definition)
// kubeconfig - Kubernetes配置信息 (Kubernetes configuration information)
// namespace - 命名空间 (Namespace)
// saName - 服务账户名称 (Service account name)
// 返回值: YAML字符串和错误信息 (Return: YAML string and error information)
func GetServiceAccountsYaml(kubeconfig, namespace, saName string) (string, error) {
	sa, err := common.ClientSet(kubeconfig).CoreV1().ServiceAccounts(namespace).Get(context.TODO(), saName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	saUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(sa)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(saUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
