// clusterRoles_model.go
// ClusterRoles数据模型定义及操作函数 / ClusterRoles data model definition and operation functions
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

// ClusterRoles ClusterRoles信息结构体 / ClusterRoles information struct
type ClusterRoles struct {
	CrName      string    `json:"crName"`      //ClusterRole名称 / ClusterRole name
	Annotations string    `json:"annotations"` //注解 / Annotations
	CreateTime  string    `json:"createTime"`  //创建时间 / Creation time
	Rules       []CrRules `json:"rules"`       //规则列表 / Rules list
}

// CrRules ClusterRole中的规则 / Rule in ClusterRole
type CrRules struct {
	Verbs     string `json:"verbs"`     //动词列表 / Verbs list
	ApiGroups string `json:"apiGroups"` //API组 / API groups
	Resources string `json:"resources"` //资源列表 / Resources list
}

// ClusterRolesList 获取ClusterRoles列表 / Get ClusterRoles list
// kubeconfig: k8s集群配置 / k8s cluster configuration
func ClusterRolesList(kubeconfig string) ([]ClusterRoles, error) {
	xList, err := common.ClientSet(kubeconfig).RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list ClusterRoles error, err:%v\n", err)
	}
	var bbb = make([]ClusterRoles, 0)
	for _, cr := range xList.Items {

		var annotationsStr string
		for k2, v2 := range cr.ObjectMeta.Annotations {
			annotationsStr += fmt.Sprintf("%s:%s,", k2, v2)
		}
		var rulesArry = make([]CrRules, 0)
		if cr.Rules != nil {
			for _, v3 := range cr.Rules {
				var verbs, apiGroups, resources string
				if v3.Verbs != nil {
					verbs = strings.Join(v3.Verbs, ",")
				}
				if v3.APIGroups != nil {
					apiGroups = strings.Join(v3.APIGroups, ",")
				}
				if v3.Resources != nil {
					resources = strings.Join(v3.Resources, ",")
				}

				rulesArry = append(rulesArry, CrRules{
					Verbs:     verbs,
					ApiGroups: apiGroups,
					Resources: resources,
				})
			}
		}

		Items := &ClusterRoles{
			CrName:      cr.Name,
			Annotations: annotationsStr,
			CreateTime:  cr.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Rules:       rulesArry,
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// GetClusterRolesYaml 获取ClusterRole的YAML定义 / Get ClusterRole YAML definition
// kubeconfig: k8s集群配置 / k8s cluster configuration
// crName: ClusterRole名称 / ClusterRole name
func GetClusterRolesYaml(kubeconfig, crName string) (string, error) {
	clusterRoles, err := common.ClientSet(kubeconfig).RbacV1().ClusterRoles().Get(context.TODO(), crName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	crUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(clusterRoles)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(crUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}