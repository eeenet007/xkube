package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"xkube/common"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

type Daemonset struct {
	// DaemonsetName - Daemonset资源的名称 / Name of the Daemonset resource
	DaemonsetName string `json:"daemonsetName"`
	// NameSpace - Daemonset所属的命名空间 / Namespace to which the Daemonset belongs
	NameSpace string `json:"nameSpace" form:"nameSpace"`
	// RevisionHistoryLimit - 保留的历史版本数量 / Number of historical versions to retain
	RevisionHistoryLimit int32 `json:"revisionHistoryLimit" form:"historyVersionLimit"`
	// PodNumber - Pod数量状态信息 / Pod quantity status information
	PodNumber string `json:"podNumber"`
	// Labels - 资源标签 / Resource labels
	Labels string `json:"labels" form:"labels"`

	// Strategy - 更新策略 / Update strategy
	Strategy string `json:"strategy"`
	// Selector - 标签选择器 / Label selector
	Selector string `json:"selector"`
	// Annotations - 资源注解 / Resource annotations
	Annotations string `json:"annotations"`
	// Status - 资源状态 / Resource status
	Status string `json:"status"`

	// ContainerName - 容器名称 / Container name
	ContainerName string `form:"containerName"`
	// ImgUrl - 镜像地址 / Image URL
	ImgUrl string `json:"imgUrl" form:"imgUrl"`
	// ContainerPortName - 容器端口名称 / Container port name
	ContainerPortName string `json:"-" form:"containerPortName"`
	// ContainerPort - 容器端口 / Container port
	ContainerPort string `form:"containerPort"`
	// CreateTime - 创建时间 / Creation time
	CreateTime string `json:"createTime"` //创建时间
}

// DaemonsetList - 获取Daemonset列表 / Get Daemonset list
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// namespace - 命名空间，如果为空则查询所有命名空间 / Namespace, if empty, query all namespaces
// daemonsetName - Daemonset名称 / Daemonset name
// labelsKey - 标签键 / Label key
// labelsValue - 标签值 / Label value
// []Daemonset - 返回Daemonset列表 / Return Daemonset list
// error - 错误信息 / Error information
func DaemonsetList(kubeconfig, namespace, daemonsetName string, labelsKey, labelsValue string) ([]Daemonset, error) {
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}
	var bbb = make([]Daemonset, 0)
	xList, err := common.ClientSet(kubeconfig).AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("[ERROR] ListStsError err:%v", err)
		return bbb, err
	}

	for _, vv := range xList.Items {
		var labelsStr, imgUrlStr, containerNameStr, ContainerPortNameStr, ContainerPortStr string
		for k1, v1 := range vv.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		for _, v2 := range vv.Spec.Template.Spec.Containers {
			containerNameStr += fmt.Sprintf("%s,", v2.Name)
			imgUrlStr += fmt.Sprintf("%s,", v2.Image)
			if len(v2.Ports) > 0 {
				ContainerPortNameStr += fmt.Sprintf("%s,", v2.Ports[0].Name)
				ContainerPortStr += fmt.Sprintf("%d,", v2.Ports[0].ContainerPort)
			}
		}
		if len(containerNameStr) > 0 {
			containerNameStr = containerNameStr[0 : len(containerNameStr)-1]
		}
		if len(imgUrlStr) > 0 {
			imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
		}
		if len(ContainerPortNameStr) > 0 {
			ContainerPortNameStr = ContainerPortNameStr[0 : len(ContainerPortNameStr)-1]
		}
		if len(ContainerPortStr) > 0 {
			ContainerPortStr = ContainerPortStr[0 : len(ContainerPortStr)-1]
		}
		xItems := &Daemonset{
			DaemonsetName:        vv.Name,
			NameSpace:            vv.Namespace,
			RevisionHistoryLimit: *vv.Spec.RevisionHistoryLimit,
			PodNumber:            fmt.Sprintf("%d/%d", vv.Status.NumberReady, vv.Status.CurrentNumberScheduled),
			Labels:               labelsStr,

			ContainerName: containerNameStr,
			ImgUrl:        imgUrlStr,

			ContainerPortName: ContainerPortNameStr,
			ContainerPort:     ContainerPortStr,
			//HostPort:          vv.Spec.Template.Spec.Containers[0].Ports[0].HostPort,

			CreateTime: vv.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *xItems)
	}
	return bbb, err
}

// DaemonsetDetail - 获取Daemonset详细信息 / Get Daemonset detailed information
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// namespace - 命名空间 / Namespace
// deploy - Daemonset名称 / Daemonset name
// *Daemonset - 返回Daemonset详细信息指针 / Return pointer to Daemonset detailed information
// error - 错误信息 / Error information
func DaemonsetDetail(kubeconfig, namespace, deploy string) (*Daemonset, error) {
	//创建一个deployment资源的接口对象DaemonsetmentClient，用于操作指定名称空间的deployment资源
	//Create a DaemonsetmentClient interface object for the deployment resource to operate the deployment resource in the specified namespace
	DaemonsetmentClient := common.ClientSet(kubeconfig).AppsV1().DaemonSets(namespace)
	//调用接口对象DaemonsetmentClient中的Get方法，获取相应的deployment资源数据
	//Call the Get method in the DaemonsetmentClient interface object to obtain the corresponding deployment resource data
	sts, err := DaemonsetmentClient.Get(context.TODO(), deploy, metav1.GetOptions{})
	if err != nil {
		return &Daemonset{}, err
	}
	//return返回指针类型的deployment结构体类型的实例
	//Return an instance of the pointer type of the deployment structure
	var selectorStr, labelsStr, imgUrlStr, containerNameStr, ContainerPortNameStr, ContainerPortStr, annotationsStr string
	for k1, v1 := range sts.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	for kk, vv := range sts.ObjectMeta.Annotations {
		if strings.Contains(kk, "last-applied-configuration") {
			continue
		}
		annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	for kk, vv := range sts.Spec.Selector.MatchLabels {
		selectorStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(selectorStr) > 0 {
		selectorStr = selectorStr[0 : len(selectorStr)-1]
	}

	for _, v2 := range sts.Spec.Template.Spec.Containers {
		containerNameStr += fmt.Sprintf("%s,", v2.Name)
		imgUrlStr += fmt.Sprintf("%s,", v2.Image)
		if len(v2.Ports) > 0 {
			ContainerPortNameStr += fmt.Sprintf("%s,", v2.Ports[0].Name)
			ContainerPortStr += fmt.Sprintf("%d,", v2.Ports[0].ContainerPort)
		}
	}
	if len(containerNameStr) > 0 {
		containerNameStr = containerNameStr[0 : len(containerNameStr)-1]
	}
	if len(imgUrlStr) > 0 {
		imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
	}
	if len(ContainerPortNameStr) > 0 {
		ContainerPortNameStr = ContainerPortNameStr[0 : len(ContainerPortNameStr)-1]
	}
	if len(ContainerPortStr) > 0 {
		ContainerPortStr = ContainerPortStr[0 : len(ContainerPortStr)-1]
	}
	//dep.Spec.Strategy.Type
	//dep.Status.Conditions
	return &Daemonset{
		DaemonsetName:        sts.Name,
		NameSpace:            sts.Namespace,
		RevisionHistoryLimit: *sts.Spec.RevisionHistoryLimit,
		PodNumber:            fmt.Sprintf("%d/%d", sts.Status.NumberReady, sts.Status.CurrentNumberScheduled),
		Labels:               labelsStr[0 : len(labelsStr)-1],
		ContainerName:        containerNameStr,
		ImgUrl:               imgUrlStr,
		ContainerPortName:    ContainerPortNameStr,
		ContainerPort:        ContainerPortStr,
		//HostPort:             sts.Spec.Template.Spec.Containers[0].Ports[0].HostPort,
		Strategy:    string(sts.Spec.UpdateStrategy.Type),
		Selector:    selectorStr,
		Annotations: annotationsStr,
		Status:      fmt.Sprintf("就绪:%d/%d个,已更新:%d个,可用:%d个", sts.Status.NumberReady, sts.Status.CurrentNumberScheduled, sts.Status.UpdatedNumberScheduled, sts.Status.NumberAvailable),
		CreateTime:  sts.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}, nil
}

// DaemonsetCreate - 创建Daemonset资源 / Create Daemonset resource
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// daemonset - Daemonset结构体指针，包含创建所需信息 / Daemonset structure pointer containing information needed for creation
// error - 错误信息 / Error information
func DaemonsetCreate(kubeconfig string, daemonset *Daemonset) error {
	//定义deployment资源的创建清单，类似k8s中通过创建yaml清单来部署集群资源
	//Define the creation manifest for the deployment resource, similar to deploying cluster resources through creating a yaml manifest in k8s
	var labelsMap = make(map[string]string)
	for _, vv := range strings.Split(daemonset.Labels, ",") {
		vArry := strings.Split(vv, ":")
		kk := vArry[0]
		kv := vArry[1]
		labelsMap[kk] = kv
	}
	ContainerPortInt, _ := strconv.Atoi(daemonset.ContainerPort)
	dsInstance := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{ //定义deployment资源的apiVersion、kind类型等
			Kind:       "Daemonset",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{ //定义控制器的metadata：名称、namespace、labels等
			Name:      daemonset.DaemonsetName,
			Namespace: daemonset.NameSpace,
			Labels:    labelsMap,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelsMap,
			},
			RevisionHistoryLimit: &daemonset.RevisionHistoryLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					//Name: "golang-pod",
					Labels: labelsMap,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:            daemonset.ContainerName,
							Image:           daemonset.ImgUrl,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          daemonset.ContainerPortName,
									ContainerPort: int32(ContainerPortInt),
									//HostPort:      daemonset.HostPort,
									Protocol: corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}

	dsClient := common.ClientSet(kubeconfig).AppsV1().DaemonSets(daemonset.NameSpace)
	_, err := dsClient.Create(context.TODO(), dsInstance, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DaemonsetModify - 更新Daemonset资源 / Update Daemonset resource
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// newDaemonset - 新的Daemonset结构体指针，包含更新信息 / New Daemonset structure pointer containing update information
// error - 错误信息 / Error information
func DaemonsetModify(kubeconfig string, newDaemonset *Daemonset) error {
	dsClient := common.ClientSet(kubeconfig).AppsV1().DaemonSets(newDaemonset.NameSpace)
	dsInstance, err := dsClient.Get(context.TODO(), newDaemonset.DaemonsetName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ContainerPortInt, _ := strconv.Atoi(newDaemonset.ContainerPort)
	dsInstance.Spec.Template.Spec.Containers[0].Image = newDaemonset.ImgUrl
	dsInstance.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = int32(ContainerPortInt)
	//stsInstance.Spec.Template.Spec.Containers[0].Ports[0].HostPort = newDaemonset.HostPort
	//调用DaemonsetmentClient接口对象中的update方法，来更新deployment资源
	_, err = dsClient.Update(context.TODO(), dsInstance, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil

}

// DaemonsetDel - 删除Daemonset资源 / Delete Daemonset resource
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// namespace - 命名空间 / Namespace
// daemonset - Daemonset名称 / Daemonset name
// error - 错误信息 / Error information
func DaemonsetDel(kubeconfig, namespace, daemonset string) error {
	stsClient := common.ClientSet(kubeconfig).AppsV1().DaemonSets(namespace)
	err := stsClient.Delete(context.TODO(), daemonset, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetDaemonsetYaml - 获取Daemonset的YAML定义 / Get Daemonset YAML definition
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// namespace - 命名空间 / Namespace
// daemonset - Daemonset名称 / Daemonset name
// string - 返回Daemonset的YAML字符串 / Return Daemonset YAML string
// error - 错误信息 / Error information
func GetDaemonsetYaml(kubeconfig, namespace, daemonset string) (string, error) {
	// 创建Daemonset客户端实例
	// Create Daemonset client instance
	dsClient := common.ClientSet(kubeconfig).AppsV1().DaemonSets(namespace)

	// 获取指定名称的Daemonset实例
	// Get the Daemonset instance with the specified name
	dsInstance, err := dsClient.Get(context.TODO(), daemonset, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// 将Daemonset实例转换为非结构化对象
	// Convert Daemonset instance to unstructured object
	dsUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(dsInstance)
	if err != nil {
		return "", err
	}

	// 将非结构化对象序列化为YAML格式
	// Serialize the unstructured object into YAML format
	yamlBytes, err := yaml.Marshal(dsUnstructured)
	if err != nil {
		return "", err
	}

	// 返回YAML字符串和nil错误
	// Return YAML string and nil error
	return string(yamlBytes), nil
}

// DaemonsetYamlModify - 通过YAML更新Daemonset资源 / Update Daemonset resource via YAML
// kubeconfig - Kubernetes配置信息 / Kubernetes configuration information
// yamlData - YAML格式的Daemonset定义数据 / Daemonset definition data in YAML format
// error - 错误信息 / Error information
func DaemonsetYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	ds := &appsv1.DaemonSet{}
	err = json.Unmarshal(data, ds)
	if err != nil {
		return err
	}
	//cluster := ds.ObjectMeta.ClusterName
	namespace := ds.ObjectMeta.Namespace
	//dsName := ds.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.AppsV1().DaemonSets(namespace).Update(context.TODO(), ds, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	//fmt.Println(namespace, dsName)
	//fmt.Println(newds)
	return err
}
