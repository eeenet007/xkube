// pod_model.go
package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	//"time"
	//"runtime"
	"xkube/common"

	//"k8s.io/api"
	//appsv1 "k8s.io/api/apps/v1"
	"encoding/json"

	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"

	//networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
)

type Configmap struct {
	// ConfigmapName name of the ConfigMap
	// ConfigmapName ConfigMap的名称
	ConfigmapName string `json:"configmapName"`
	// NameSpace namespace where the ConfigMap is located
	// NameSpace ConfigMap所在的命名空间
	NameSpace string `json:"nameSpace"`
	// Labels labels of the ConfigMap
	// Labels ConfigMap的标签
	Labels string `json:"labels"`
	// CreateTime creation time of the ConfigMap
	// CreateTime ConfigMap的创建时间
	CreateTime string `json:"createTime"`
	// Data key-value data stored in the ConfigMap
	// Data ConfigMap中存储的键值对数据
	Data []Cmkv `json:"data"`
}

// Cmkv key-value pair structure for ConfigMap data
// Cmkv ConfigMap数据的键值对结构
type Cmkv struct {
	// Key key of the ConfigMap data
	// Key ConfigMap数据的键
	Key string `json:"key"`
	// Value value of the ConfigMap data
	// Value ConfigMap数据的值
	Value string `json:"value"`
}

// CmList 获取指定集群和命名空间中的ConfigMap列表，支持按名称和标签筛选
// CmList Get the list of ConfigMaps in the specified cluster and namespace, support filtering by name and labels
func CmList(kubeconfig, namespace, configmapName string, labelsKey, labelsValue string) ([]Configmap, error) {
	clientset := common.ClientSet(kubeconfig)
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	//设置ListOptions
	//Set ListOptions
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	cmList, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}

	var bbb = make([]Configmap, 0)
	for _, cm := range cmList.Items {
		//搜索
		//Search
		if configmapName != "" {
			if !strings.Contains(cm.Name, configmapName) {
				continue
			}
		}
		var labelsStr string
		for kk, vv := range cm.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		Items := &Configmap{
			ConfigmapName: cm.Name,
			NameSpace:     cm.Namespace,
			Labels:        labelsStr,
			CreateTime:    cm.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// CmCreate Create a ConfigMap from JSON data
// CmCreate 从JSON数据创建ConfigMap
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// bodys: JSON data containing ConfigMap configuration
// bodys: 包含ConfigMap配置的JSON数据
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CmCreate(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	configmapName := gp.Get("configmapName").String()
	nameSpace := gp.Get("nameSpace").String()

	var labelsMap = make(map[string]string)
	labelsMap["app"] = configmapName
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	var dataMap = make(map[string]string)
	for _, vv := range gp.Get("configmaps").Array() {
		dataMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Data: dataMap,
	}

	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.CoreV1().ConfigMaps(nameSpace).Create(context.TODO(), cm, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// CmYamlCreate Create a ConfigMap from YAML data
// CmYamlCreate 从YAML数据创建ConfigMap
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// yamlData: YAML data containing ConfigMap configuration
// yamlData: 包含ConfigMap配置的YAML数据
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CmYamlCreate(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	configmap := &corev1.ConfigMap{}
	err = json.Unmarshal(data, configmap)
	if err != nil {
		return err
	}

	namespace := configmap.ObjectMeta.Namespace
	confignameName := configmap.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, confignameName)
	return err
}

// CmYamlModify Modify a ConfigMap from YAML data
// CmYamlModify 从YAML数据修改ConfigMap
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// yamlData: YAML data containing updated ConfigMap configuration
// yamlData: 包含更新的ConfigMap配置的YAML数据
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CmYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	configmap := &corev1.ConfigMap{}
	err = json.Unmarshal(data, configmap)
	if err != nil {
		return err
	}

	namespace := configmap.ObjectMeta.Namespace
	//confignameName := configmap.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configmap, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	//fmt.Println(namespace, confignameName)
	return err
}

// CmUpdate Update a ConfigMap
// CmUpdate 更新ConfigMap
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// namespace: namespace where the ConfigMap is located
// namespace: ConfigMap所在的命名空间
// configmap: ConfigMap object with updated data
// configmap: 包含更新数据的ConfigMap对象
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CmUpdate(kubeconfig, namespace string, configmap *Configmap) error {
	clientset := common.ClientSet(kubeconfig)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-configmap",
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
	cm.Data["key1"] = "new-value1"
	result, err := clientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated configmap %q.\n", result.GetObjectMeta().GetName())
	return err
}

// CmClone Clone a ConfigMap from one cluster/namespace to another
// CmClone 将ConfigMap从一个集群/命名空间克隆到另一个
// kubeconfig: source cluster configuration identifier
// kubeconfig: 源集群配置信息标识符
// namespace: source namespace
// namespace: 源命名空间
// objname: source ConfigMap name
// objname: 源ConfigMap名称
// target_clusterid: target cluster identifier
// target_clusterid: 目标集群标识符
// target_namespace: target namespace
// target_namespace: 目标命名空间
// target_objname: target ConfigMap name
// target_objname: 目标ConfigMap名称
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CmClone(kubeconfig, namespace, objname, target_clusterid, target_namespace, target_objname string) error {
	//old cluster
	cm, err := common.ClientSet(kubeconfig).CoreV1().ConfigMaps(namespace).Get(context.TODO(), objname, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if target_namespace == "" && target_objname == "" && target_clusterid == "" {
		return fmt.Errorf("target can't be all empty")
	}

	if target_namespace == "" {
		target_namespace = namespace
	}

	if target_objname == "" {
		target_objname = objname
	}

	if target_clusterid == "" {
		target_clusterid = kubeconfig
	}

	//在新集群创建namespace
	err3 := CreateNsByExist(target_clusterid, target_namespace)
	if err3 != nil {
		return err3
	}

	cm.Name = target_objname
	cm.ResourceVersion = ""
	cm.Namespace = target_namespace
	cm.ObjectMeta.Labels["app"] = target_objname

	//new cluster
	NewCmClient := common.ClientSet(target_clusterid).CoreV1().ConfigMaps(target_namespace)
	_, err2 := NewCmClient.Get(context.TODO(), target_objname, metav1.GetOptions{})
	if errors.IsNotFound(err2) { //没有就创建
		_, err := NewCmClient.Create(context.TODO(), cm, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if err2 == nil {
		_, err = NewCmClient.Update(context.TODO(), cm, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	return err2
}

// CmDetail Get detailed information of a ConfigMap
// CmDetail 获取ConfigMap的详细信息
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// namespace: namespace where the ConfigMap is located
// namespace: ConfigMap所在的命名空间
// configmapName: name of the ConfigMap
// configmapName: ConfigMap的名称
// Returns:
//   - *Configmap: ConfigMap detailed information
//   - error: error message
//
// 返回值:
//   - *Configmap: ConfigMap详细信息
//   - error: 错误信息
func CmDetail(kubeconfig, namespace, configmapName string) (*Configmap, error) {
	cm, err := common.ClientSet(kubeconfig).CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	var labelsStr string
	for kk, vv := range cm.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}
	var bbb = make([]Cmkv, 0)
	for k1, v1 := range cm.Data {
		Items := &Cmkv{
			Key:   k1,
			Value: v1,
		}
		bbb = append(bbb, *Items)
	}
	return &Configmap{
		ConfigmapName: cm.Name,
		NameSpace:     cm.Namespace,
		Labels:        labelsStr,
		CreateTime:    cm.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Data:          bbb,
	}, nil
}

// GetCmYaml Get YAML representation of a ConfigMap
// GetCmYaml 获取ConfigMap的YAML表示
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// namespace: namespace where the ConfigMap is located
// namespace: ConfigMap所在的命名空间
// configmapName: name of the ConfigMap
// configmapName: ConfigMap的名称
// Returns:
//   - string: YAML representation of the ConfigMap
//   - error: error message
//
// 返回值:
//   - string: ConfigMap的YAML表示
//   - error: 错误信息
func GetCmYaml(kubeconfig, namespace, configmapName string) (string, error) {
	cmClient := common.ClientSet(kubeconfig).CoreV1().ConfigMaps(namespace)
	configmap, err := cmClient.Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	cmUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(configmap)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(cmUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// CmDelete Delete a ConfigMap
// CmDelete 删除ConfigMap
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// namespace: namespace where the ConfigMap is located
// namespace: ConfigMap所在的命名空间
// configmapName: name of the ConfigMap to delete
// configmapName: 要删除的ConfigMap名称
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func CmDelete(kubeconfig, namespace, configmapName string) error {
	clientset := common.ClientSet(kubeconfig)
	err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configmapName, metav1.DeleteOptions{})
	return err
}
