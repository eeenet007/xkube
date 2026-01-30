// clusterRoleBinding_model.go
// ClusterRoleBinding数据模型定义及操作函数 / ClusterRoleBinding data model definition and operation functions
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

// ClusterRoleBinding ClusterRoleBinding信息结构体 / ClusterRoleBinding information struct
type ClusterRoleBinding struct {
	CrbName     string        `json:"crbName"`     //ClusterRoleBinding名称 / ClusterRoleBinding name
	Labels      string        `json:"labels"`      //标签 / Labels
	Annotations string        `json:"annotations"` //注解 / Annotations
	CreateTime  string        `json:"createTime"`  //创建时间 / Creation time
	RoleRef     string        `json:"roleRef"`     //角色引用 / Role reference
	Subjects    []CrbSubjects `json:"subjects"`    //主体列表 / Subjects list
}

// CrbSubjects ClusterRoleBinding中的主体 / Subject in ClusterRoleBinding
type CrbSubjects struct {
	Kind      string `json:"kind"`      //主体类型 / Subject kind
	NameSpace string `json:"nameSpace"` //命名空间 / Namespace
	Name      string `json:"name"`      //名称 / Name
}

// ClusterRoleBindingList 获取ClusterRoleBinding列表 / Get ClusterRoleBinding list
// kubeconfig: k8s集群配置 / k8s cluster configuration
func ClusterRoleBindingList(kubeconfig string) ([]ClusterRoleBinding, error) {
	xList, err := common.ClientSet(kubeconfig).RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}
	var bbb = make([]ClusterRoleBinding, 0)
	for _, crb := range xList.Items {
		var labelsStr string
		for kk, vv := range crb.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		var annotationsStr string
		for k2, v2 := range crb.ObjectMeta.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}
		if len(annotationsStr) > 0 {
			annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
		}

		roleRef := fmt.Sprintf("apiGroup:%s,kind:%s,name:%s", crb.RoleRef.APIGroup, crb.RoleRef.Kind, crb.RoleRef.Name)

		var subjectsArry = make([]CrbSubjects, 0)
		if crb.Subjects != nil {
			for _, v3 := range crb.Subjects {
				subjectsArry = append(subjectsArry, CrbSubjects{
					Kind:      v3.Kind,
					NameSpace: v3.Namespace,
					Name:      v3.Name,
				})
			}
		}

		Items := &ClusterRoleBinding{
			CrbName:     crb.Name,
			Labels:      labelsStr,
			Annotations: annotationsStr,
			RoleRef:     roleRef,
			CreateTime:  crb.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Subjects:    subjectsArry,
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// GetClusterRoleBindingYaml 获取ClusterRoleBinding的YAML定义 / Get ClusterRoleBinding YAML definition
// kubeconfig: k8s集群配置 / k8s cluster configuration
// crbName: ClusterRoleBinding名称 / ClusterRoleBinding name
func GetClusterRoleBindingYaml(kubeconfig, crbName string) (string, error) {
	crb, err := common.ClientSet(kubeconfig).RbacV1().ClusterRoleBindings().Get(context.TODO(), crbName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	crbUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(crb)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(crbUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}