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
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/yaml"
	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/util/intstr"
)

type Secret struct {
	SecretName string  `json:"secretName"`
	NameSpace  string  `json:"nameSpace"`
	Labels     string  `json:"labels"`
	SecretType string  `json:"secretType"`
	CreateTime string  `json:"createTime"`
	Data       []Seckv `json:"data"`
}

type Seckv struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SecretList 根据指定条件列出Secret资源
// 参数:
//   - kubeconfig: Kubernetes配置信息，用于连接集群
//   - namespace: 命名空间，如果为空则查询所有命名空间
//   - secretName: Secret名称，用于模糊匹配筛选
//   - labelsKey: 标签键，用于标签筛选
//   - labelsValue: 标签值，用于标签筛选
//
// 返回值:
//   - []Secret: Secret列表
//   - error: 错误信息
func SecretList(kubeconfig, namespace, secretName, labelsKey, labelsValue string) ([]Secret, error) {
	clientset := common.ClientSet(kubeconfig)
	// 如果命名空间为空，则查询所有命名空间
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	// 构造标签筛选条件
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	// 执行Secret查询
	secretList, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}

	// 构造返回结果
	var bbb = make([]Secret, 0)
	for _, sec := range secretList.Items {
		// 根据secret名称进行模糊匹配筛选
		if secretName != "" {
			if !strings.Contains(sec.Name, secretName) {
				continue
			}
		}

		// 构造标签字符串
		var labelsStr string
		for kk, vv := range sec.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		Items := &Secret{
			SecretName: sec.Name,
			NameSpace:  sec.Namespace,
			Labels:     labelsStr,
			SecretType: string(sec.Type),
			CreateTime: sec.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// SecretCreate creates a new Secret resource in the Kubernetes cluster
// Parameters:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - bodys: JSON byte array containing the secret configuration details
//
// Returns:
//   - error: Error information if the operation fails, nil otherwise
func SecretCreate(kubeconfig string, bodys []byte) error {
	// Parse the request body to extract secret configuration
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	secretName := gp.Get("secretName").String()
	nameSpace := gp.Get("nameSpace").String()

	secretType := gp.Get("secretType").String()

	// Initialize labels map with default app label
	var labelsMap = make(map[string]string)
	labelsMap["app"] = secretName
	// Add custom labels from request
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	// Extract secret data key-value pairs
	var dataMap = make(map[string]string)
	for _, vv := range gp.Get("secrets").Array() {
		dataMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	// Determine secret type based on input
	var secType corev1.SecretType

	switch secretType {
	case "dockerconfigjson":
		secType = corev1.SecretTypeDockerConfigJson
	case "tls":
		secType = corev1.SecretTypeTLS
	default:
		secType = corev1.SecretTypeOpaque
	}

	// Create Secret object with provided configuration
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		StringData: dataMap,
		Type:       secType,
	}

	// Connect to Kubernetes cluster and create the Secret
	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.CoreV1().Secrets(nameSpace).Create(context.TODO(), newSecret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// SecretYamlCreate creates a new Secret resource in the Kubernetes cluster based on YAML configuration
// Parameters:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - yamlData: YAML byte array containing the secret configuration
//
// Returns:
//   - error: Error information if the operation fails, nil otherwise
func SecretYamlCreate(kubeconfig string, yamlData []byte) error {
	// Convert YAML data to JSON format
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}

	// Unmarshal JSON data into Secret struct
	secrets := &corev1.Secret{}
	err = json.Unmarshal(data, secrets)
	if err != nil {
		return err
	}

	// Extract namespace and secret name from the Secret object
	namespace := secrets.ObjectMeta.Namespace
	secretsName := secrets.ObjectMeta.Name

	// Connect to Kubernetes cluster and create the Secret
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secrets, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Log the created secret's namespace and name
	fmt.Println(namespace, secretsName)
	return err
}

// SecretYamlModify updates an existing Secret resource in the Kubernetes cluster based on YAML configuration
// Parameters:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - yamlData: YAML byte array containing the updated secret configuration
//
// Returns:
//   - error: Error information if the operation fails, nil otherwise
func SecretYamlModify(kubeconfig string, yamlData []byte) error {
	// Convert YAML data to JSON format
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}

	// Unmarshal JSON data into Secret struct
	secrets := &corev1.Secret{}
	err = json.Unmarshal(data, secrets)
	if err != nil {
		return err
	}

	// Extract namespace from the Secret object
	namespace := secrets.ObjectMeta.Namespace
	//secretsName := secrets.ObjectMeta.Name

	// Connect to Kubernetes cluster and update the Secret
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secrets, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return err
}

// SecretUpdate updates an existing Secret resource in the Kubernetes cluster
// Parameters:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - namespace: The namespace where the secret is located
//   - secret: Pointer to the Secret object containing updated information
//
// Returns:
//   - error: Error information if the operation fails, nil otherwise
func SecretUpdate(kubeconfig, namespace string, secret *Secret) error {
	// Connect to Kubernetes cluster
	clientset := common.ClientSet(kubeconfig)

	// Retry on conflict to handle concurrent updates
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Get the existing secret from the cluster
		secret, err := clientset.CoreV1().Secrets("default").Get(context.Background(), "my-secret", metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		// Update the secret data
		secret.StringData["password"] = "new-password"

		// Apply the update to the cluster
		result, err := clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
		fmt.Printf("Updated secret %q.\n", result.GetObjectMeta().GetName())
		return err
	})

	// Handle retry errors
	if err != nil {
		panic(err)
	}
	return err
}

// SecretClone clones a Secret resource from one Kubernetes cluster/namespace to another cluster/namespace
// Parameters:
//   - kubeconfig: Source Kubernetes configuration for connecting to the cluster
//   - namespace: Source namespace where the secret is located
//   - objname: Name of the secret to be cloned
//   - target_clusterid: Target cluster ID where the secret will be cloned to
//   - target_namespace: Target namespace where the secret will be cloned to
//   - target_objname: Target name for the cloned secret
//
// Returns:
//   - error: Error information if the operation fails, nil otherwise
func SecretClone(kubeconfig, namespace, objname, target_clusterid, target_namespace, target_objname string) error {
	// Get the secret from the source cluster
	sec, err := common.ClientSet(kubeconfig).CoreV1().Secrets(namespace).Get(context.TODO(), objname, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Validate that at least one target parameter is provided
	if target_namespace == "" && target_objname == "" && target_clusterid == "" {
		return fmt.Errorf("target can't be all empty")
	}

	// Set default values for target parameters if not provided
	if target_namespace == "" {
		target_namespace = namespace
	}

	if target_objname == "" {
		target_objname = objname
	}

	if target_clusterid == "" {
		target_clusterid = kubeconfig
	}

	// Create target namespace if it doesn't exist
	err3 := CreateNsByExist(target_clusterid, target_namespace)
	if err3 != nil {
		return err3
	}

	// Update secret metadata for the target environment
	sec.Name = target_objname
	sec.ResourceVersion = ""
	sec.Namespace = target_namespace
	sec.ObjectMeta.Labels["app"] = target_objname

	// Get client for the target cluster
	NewSecretClient := common.ClientSet(target_clusterid).CoreV1().Secrets(target_namespace)
	_, err2 := NewSecretClient.Get(context.TODO(), target_objname, metav1.GetOptions{})

	// If secret doesn't exist in target, create it
	if errors.IsNotFound(err2) { //没有就创建
		_, err := NewSecretClient.Create(context.TODO(), sec, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	// If secret exists in target, update it
	if err2 == nil {
		_, err = NewSecretClient.Update(context.TODO(), sec, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	return err2
}

// SecretDetail gets detailed information of a Secret resource from the Kubernetes cluster
// SecretDetail 从Kubernetes集群中获取Secret资源的详细信息
// Parameters:
// 参数:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - kubeconfig: 用于连接到集群的Kubernetes配置信息
//   - namespace: The namespace where the secret is located
//   - namespace: secret所在的命名空间
//   - secretName: Name of the secret to retrieve
//   - secretName: 要获取的secret名称
//
// Returns:
// 返回值:
//   - *Secret: Pointer to the Secret object containing detailed information
//   - *Secret: 指向包含详细信息的Secret对象的指针
//   - error: Error information if the operation fails, nil otherwise
//   - error: 如果操作失败返回错误信息，否则为nil
func SecretDetail(kubeconfig, namespace, secretName string) (*Secret, error) {
	// Get the secret from the cluster
	// 从集群中获取secret
	sec, err := common.ClientSet(kubeconfig).CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	// Construct labels string from the secret's labels
	// 从secret的标签构造标签字符串
	var labelsStr string
	for kk, vv := range sec.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	// Extract secret data key-value pairs
	// 提取secret数据键值对
	var bbb = make([]Seckv, 0)
	for k1, v1 := range sec.Data {
		Items := &Seckv{
			Key:   k1,
			Value: string(v1),
		}
		bbb = append(bbb, *Items)
	}

	// Return the constructed Secret object
	// 返回构造的Secret对象
	return &Secret{
		SecretName: sec.Name,
		NameSpace:  sec.Namespace,
		Labels:     labelsStr,
		SecretType: string(sec.Type),
		CreateTime: sec.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Data:       bbb,
	}, nil
	//bbb = append(bbb, *Items)
}

// GetSecretYaml retrieves a Secret resource from the Kubernetes cluster and returns it as a YAML string
// GetSecretYaml 从Kubernetes集群中获取Secret资源并将其作为YAML字符串返回
// Parameters:
// 参数:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - kubeconfig: 用于连接到集群的Kubernetes配置信息
//   - namespace: The namespace where the secret is located
//   - namespace: secret所在的命名空间
//   - secretName: Name of the secret to retrieve
//   - secretName: 要获取的secret名称
//
// Returns:
// 返回值:
//   - string: Secret resource in YAML format
//   - string: YAML格式的Secret资源
//   - error: Error information if the operation fails, nil otherwise
//   - error: 如果操作失败返回错误信息，否则为nil
func GetSecretYaml(kubeconfig, namespace, secretName string) (string, error) {
	// Get Secret client for the specified cluster and namespace
	// 获取指定集群和命名空间的Secret客户端
	cmClient := common.ClientSet(kubeconfig).CoreV1().Secrets(namespace)

	// Retrieve the secret from the cluster
	// 从集群中获取secret
	secret, err := cmClient.Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	// Convert secret to unstructured data format
	// 将secret转换为非结构化数据格式
	cmUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(secret)
	if err != nil {
		return "", err
	}

	// Marshal the unstructured data into YAML format
	// 将非结构化数据编组为YAML格式
	yamlBytes, err := yaml.Marshal(cmUnstructured)
	if err != nil {
		return "", err
	}

	// Return the YAML string representation of the secret
	// 返回secret的YAML字符串表示
	return string(yamlBytes), nil
}

// SecretDelete deletes a Secret resource from the Kubernetes cluster
// SecretDelete 从Kubernetes集群中删除Secret资源
// Parameters:
// 参数:
//   - kubeconfig: Kubernetes configuration information for connecting to the cluster
//   - kubeconfig: 用于连接到集群的Kubernetes配置信息
//   - namespace: The namespace where the secret is located
//   - namespace: secret所在的命名空间
//   - secretName: Name of the secret to delete
//   - secretName: 要删除的secret名称
//
// Returns:
// 返回值:
//   - error: Error information if the operation fails, nil otherwise
//   - error: 如果操作失败返回错误信息，否则为nil
func SecretDelete(kubeconfig, namespace, secretName string) error {
	// Connect to Kubernetes cluster
	// 连接到Kubernetes集群
	clientset := common.ClientSet(kubeconfig)

	// Delete the secret from the specified namespace
	// 从指定命名空间中删除secret
	err := clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
