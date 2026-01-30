// roles_model.go
package models

import (
	"context"
	"fmt"
	"strings"

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

type Roles struct {
	RolesName   string       `json:"rolesName"`
	NameSpace   string       `json:"nameSpace"`
	Labels      string       `json:"labels"`
	Annotations string       `json:"annotations"`
	CreateTime  string       `json:"createTime"`
	Rules       []RolesRules `json:"rules"`
}

type RolesRules struct {
	Verbs         string `json:"verbs"`
	ApiGroups     string `json:"apiGroups"`
	Resources     string `json:"resources"`
	ResourceNames string `json:"resourceNames"`
}

// RolesList 查询指定命名空间中的Role资源列表
// kubeconfig: Kubernetes配置文件路径，用于创建客户端
// nameSpace: 命名空间，指定要查询的Role所在的命名空间
// 返回值1: Role列表，包含Roles结构体的切片
// 返回值2: 错误信息，如果查询过程中出现错误则返回错误详情
func RolesList(kubeconfig, nameSpace string) ([]Roles, error) {
	xList, err := common.ClientSet(kubeconfig).RbacV1().Roles(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list roles error, err:%v\n", err)
	}
	var bbb = make([]Roles, 0)

	// 遍历查询到的Role资源，提取关键信息并转换为Roles结构体
	for _, rs := range xList.Items {
		var labelsStr string
		// 将Role的标签信息转换为字符串格式 key:value,key2:value2
		for kk, vv := range rs.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		var annotationsStr string
		// 将Role的注解信息转换为字符串格式 key:value,key2:value2
		for k2, v2 := range rs.ObjectMeta.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}
		if len(annotationsStr) > 0 {
			annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
		}

		var rulesArry = make([]RolesRules, 0)
		// 处理Role的规则信息，将规则中的动词、API组、资源等信息提取出来
		if rs.Rules != nil {
			for _, v3 := range rs.Rules {
				var verbs, apiGroups, resources, resourceNames string
				if v3.Verbs != nil {
					verbs = strings.Join(v3.Verbs, ",")
				}
				if v3.APIGroups != nil {
					apiGroups = strings.Join(v3.APIGroups, ",")
				}
				if v3.Resources != nil {
					resources = strings.Join(v3.Resources, ",")
				}
				if v3.ResourceNames != nil {
					resourceNames = strings.Join(v3.ResourceNames, ",")
				}

				rulesArry = append(rulesArry, RolesRules{
					Verbs:         verbs,
					ApiGroups:     apiGroups,
					Resources:     resources,
					ResourceNames: resourceNames,
				})
			}
		}

		Items := &Roles{
			RolesName:   rs.Name,
			NameSpace:   rs.Namespace,
			Labels:      labelsStr,
			Annotations: annotationsStr,
			CreateTime:  rs.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Rules:       rulesArry,
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// GetRolesYaml returns the YAML representation of the specified Role resource
//   - kubeconfig: Kubernetes configuration file path used to create the client
//   - nameSpace: Namespace where the Role to be retrieved is located
//   - rolesName: Name of the Role resource to retrieve
//
// Returns:
//   - string: YAML string representation of the Role resource
//   - error: Error information if an error occurs during retrieval
func GetRolesYaml(kubeconfig, nameSpace, rolesName string) (string, error) {
	roles, err := common.ClientSet(kubeconfig).RbacV1().Roles(nameSpace).Get(context.TODO(), rolesName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	rolesUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(roles)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(rolesUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
