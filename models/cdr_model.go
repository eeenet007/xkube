// storageclass_model.go
package models

import (
	"context"
	"log"
	"xkube/common"

	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	//apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

type CustomResourceDefinition struct {
	CdrName    string `json:"cdrName"`
	CdrKind    string `json:"cdrKind"`
	ApiGroup   string `json:"apiGroup"`
	ApiVersion string `json:"apiVersion"`
	Scope      string `json:"scope"`
}

// CdrList List all Custom Resource Definitions (CRD) in the specified kubeconfig cluster
// CdrList 获取指定kubeconfig集群中的所有自定义资源定义(CRD)列表
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// Returns:
//   - []CustomResourceDefinition: CRD list
//   - error: error message
//
// 返回值:
//   - []CustomResourceDefinition: CRD列表
//   - error: 错误信息
func CdrList(kubeconfig string) ([]CustomResourceDefinition, error) {
	config := common.ClientConfig(kubeconfig)
	var bbb = make([]CustomResourceDefinition, 0)
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		return bbb, err
	}
	crds, err := clientset.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("list crd error, err:%v\n", err)
	}

	// Iterate through the CRD list and extract key information
	// 遍历CRD列表，提取关键信息
	for _, vv := range crds.Items {
		Items := &CustomResourceDefinition{
			CdrName:    vv.ObjectMeta.Name,
			CdrKind:    vv.Spec.Names.Kind,
			ApiGroup:   vv.Spec.Group,
			ApiVersion: vv.Spec.Versions[0].Name,
			Scope:      string(vv.Spec.Scope),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// GetCdrYaml Get the YAML representation of a specified Custom Resource Definition
// GetCdrYaml 获取指定自定义资源定义的YAML表示
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// cdrName: name of the Custom Resource Definition to retrieve
// cdrName: 要获取的自定义资源定义名称
// Returns:
//   - string: YAML representation of the CRD
//   - error: error message
//
// 返回值:
//   - string: CRD的YAML表示
//   - error: 错误信息
func GetCdrYaml(kubeconfig, cdrName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	crdclientset, err := clientset.NewForConfig(config)
	if err != nil {
		return "", err
	}
	cdr, err := crdclientset.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), cdrName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	cdrUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cdr)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(cdrUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// CdrDelete Delete a specified Custom Resource Definition
// CdrDelete 删除指定的自定义资源定义
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// crdName: name of the Custom Resource Definition to delete
// crdName: 要删除的自定义资源定义名称
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CdrDelete(kubeconfig, crdName string) error {
	config := common.ClientConfig(kubeconfig)
	crdclientset, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	err = crdclientset.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), crdName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
