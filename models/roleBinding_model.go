// roleBinding_model.go
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

// RoleBinding 结构体表示Kubernetes中的RoleBinding资源
// RoleBinding struct represents the RoleBinding resource in Kubernetes
type RoleBinding struct {
	RbName      string       `json:"rbName"`
	NameSpace   string       `json:"nameSpace"`
	Labels      string       `json:"labels"`
	Annotations string       `json:"annotations"`
	CreateTime  string       `json:"createTime"`
	RoleRef     string       `json:"roleRef"`
	Subjects    []RbSubjects `json:"subjects"`
}

// RbSubjects 结构体表示RoleBinding中的Subjects字段
// RbSubjects struct represents the Subjects field in RoleBinding
type RbSubjects struct {
	Kind      string `json:"kind"`
	NameSpace string `json:"nameSpace"`
	Name      string `json:"name"`
}

// RoleBindingList 获取指定命名空间中的RoleBinding列表
// RoleBindingList gets the list of RoleBindings in the specified namespace
// kubeconfig: Kubernetes配置信息 / Kubernetes configuration information
// nameSpace: 命名空间，空字符串表示所有命名空间 / Namespace, empty string means all namespaces
// 返回RoleBinding列表和错误信息 / Returns RoleBinding list and error information
func RoleBindingList(kubeconfig, nameSpace string) ([]RoleBinding, error) {
	xList, err := common.ClientSet(kubeconfig).RbacV1().RoleBindings(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}
	var bbb = make([]RoleBinding, 0)
	for _, rb := range xList.Items {
		var labelsStr string
		for kk, vv := range rb.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		var annotationsStr string
		for k2, v2 := range rb.ObjectMeta.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}
		if len(annotationsStr) > 0 {
			annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
		}

		roleRef := fmt.Sprintf("apiGroup:%s,kind:%s,name:%s", rb.RoleRef.APIGroup, rb.RoleRef.Kind, rb.RoleRef.Name)

		var subjectsArry = make([]RbSubjects, 0)
		if rb.Subjects != nil {
			for _, v3 := range rb.Subjects {
				subjectsArry = append(subjectsArry, RbSubjects{
					Kind:      v3.Kind,
					NameSpace: v3.Namespace,
					Name:      v3.Name,
				})
			}
		}

		Items := &RoleBinding{
			RbName:      rb.Name,
			NameSpace:   rb.Namespace,
			Labels:      labelsStr,
			Annotations: annotationsStr,
			RoleRef:     roleRef,
			CreateTime:  rb.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Subjects:    subjectsArry,
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// GetRoleBindingYaml 获取指定RoleBinding的YAML配置
// GetRoleBindingYaml gets the YAML configuration of the specified RoleBinding
// kubeconfig: Kubernetes配置信息 / Kubernetes configuration information
// nameSpace: 命名空间 / Namespace
// rbName: RoleBinding名称 / RoleBinding name
// 返回YAML字符串和错误信息 / Returns YAML string and error information
func GetRoleBindingYaml(kubeconfig, nameSpace, rbName string) (string, error) {
	roleBinding, err := common.ClientSet(kubeconfig).RbacV1().RoleBindings(nameSpace).Get(context.TODO(), rbName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	rbUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(roleBinding)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(rbUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
