// hpa_model_v2beta2.go
package models

import (
	"context"
	"fmt"
	"strings"

	//"time"
	"encoding/json"
	"log"
	"xkube/common"

	"github.com/tidwall/gjson"

	//v1 "k8s.io/api/core/v1"
	//autoscalingv1 "k8s.io/api/autoscaling/v1"
	//autoscalingv2 "k8s.io/api/autoscaling/v2"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"

	//"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/retry"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

//AutoscalingV2beta2 适合zx-pcauto的版本 1.20.11-aliyun.1 【1.25以下用】
//AutoscalingV2适合yt-pcauto的版本 1.28.1 【1.25以上用】

/*
HpaListV2beta2 - 获取HPA列表 (Get HPA list)
kubeconfig: 集群配置信息 (Cluster configuration)
nameSpace: 命名空间，如果为空则查询所有命名空间 (Namespace, if empty, query all namespaces)
hpaName: HPA名称，用于模糊搜索 (HPA name for fuzzy search)
targetName: 目标对象名称 (Target object name)
*/
func HpaListV2beta2(kubeconfig, nameSpace, hpaName, targetName string) ([]Hpa, error) {
	clientset := common.ClientSet(kubeconfig)
	if nameSpace == "" {
		//namespace = corev1.NamespaceDefault
		nameSpace = corev1.NamespaceAll
	}

	//设置ListOptions (Set ListOptions)
	var listOptions = metav1.ListOptions{}
	// if hpaName != "" {
	// 	listOptions = metav1.ListOptions{
	// 		FieldSelector: fmt.Sprintf("metadata.name=%s", hpaName),
	// 	}
	// }

	//hpas, err := clientset.AutoscalingV1().HorizontalPodAutoscalers(nameSpace).List(context.TODO(), metav1.ListOptions{}) //无vv.Spec.Metrics
	//hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(nameSpace).List(context.TODO(), metav1.ListOptions{}) //zx-pcauto会提示：list haps error, err:the server could not find the requested resource
	//hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(nameSpace).List(context.TODO(), listOptions)
	hpas, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(nameSpace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list haps error, err:%v\n", err)
	}

	var bbb = make([]Hpa, 0)
	for _, vv := range hpas.Items {

		//搜索 (Search)
		if hpaName != "" {
			if !strings.Contains(vv.Name, hpaName) {
				continue
			}
		}

		//根据对象名称来查询 (Query by object name)
		if targetName != "" {
			if vv.Spec.ScaleTargetRef.Name != targetName {
				continue
			}
		}

		var metricsStr string
		for _, v1 := range vv.Spec.Metrics {
			metricsStr += fmt.Sprintf("%s:%d%%,", v1.Resource.Name, *v1.Resource.Target.AverageUtilization)
		}
		if len(metricsStr) > 0 {
			metricsStr = metricsStr[0 : len(metricsStr)-1]
		}
		Items := &Hpa{
			HpaName:         vv.Name,
			NameSpace:       vv.Namespace,
			TargetKind:      vv.Spec.ScaleTargetRef.Kind,
			TargetName:      vv.Spec.ScaleTargetRef.Name,
			MinReplicas:     *vv.Spec.MinReplicas,
			MaxReplicas:     vv.Spec.MaxReplicas,
			CurrentReplicas: vv.Status.CurrentReplicas,
			Metrics:         metricsStr,
			CreateTime:      vv.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

/*
GetHpaYamlV2beta2 - 获取HPA的YAML配置 (Get HPA YAML configuration)
kubeconfig: 集群配置信息 (Cluster configuration)
nameSpace: 命名空间 (Namespace)
hpaName: HPA名称 (HPA name)
*/
func GetHpaYamlV2beta2(kubeconfig, nameSpace, hpaName string) (string, error) {

	hpa, err := common.ClientSet(kubeconfig).AutoscalingV2beta2().HorizontalPodAutoscalers(nameSpace).Get(context.TODO(), hpaName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	hpaUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(hpa)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(hpaUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

/*
HpaCreateV2beta2 - 创建HPA (Create HPA)
kubeconfig: 集群配置信息 (Cluster configuration)
bodys: 请求体数据 (Request body data)
*/
func HpaCreateV2beta2(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	hpaName := gp.Get("hpaName").String()
	nameSpace := gp.Get("nameSpace").String()

	var labelsMap = make(map[string]string)
	labelsMap["app"] = hpaName
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	hpaKind := gp.Get("hpaKind").String()
	objName := gp.Get("objName").String()
	cpuUsage := gp.Get("cpuUsage").Int()
	memUsage := gp.Get("memUsage").Int()
	maxPods := gp.Get("maxPods").Int()
	minPods := gp.Get("minPods").Int()
	if minPods < 1 {
		minPods = 1
	}
	if maxPods < 1 {
		maxPods = 2
	}
	if cpuUsage < 1 {
		cpuUsage = 80
	}

	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hpaName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: autoscalingv2beta2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2beta2.CrossVersionObjectReference{
				Kind:       hpaKind,
				Name:       objName,
				APIVersion: "apps/v1",
			},
			MinReplicas: int32Ptr(int32(minPods)),
			MaxReplicas: int32(maxPods),
			//TargetCPUUtilizationPercentage: int32Ptr(int32(cpuUsage)), v1版本
			Metrics: []autoscalingv2beta2.MetricSpec{
				{
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: "memory",
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: int32Ptr(int32(memUsage)),
						},
					},
					Type: "Resource",
				},
				{
					Resource: &autoscalingv2beta2.ResourceMetricSource{
						Name: "cpu",
						Target: autoscalingv2beta2.MetricTarget{
							Type:               autoscalingv2beta2.UtilizationMetricType,
							AverageUtilization: int32Ptr(int32(cpuUsage)),
						},
					},
					Type: "Resource",
				},
			},
		},
	}
	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(nameSpace).Create(context.TODO(), hpa, metav1.CreateOptions{})
	return err
}

/*
HpaDeleteV2beta2 - 删除HPA (Delete HPA)
kubeconfig: 集群配置信息 (Cluster configuration)
nameSpace: 命名空间 (Namespace)
hpaName: HPA名称 (HPA name)
*/
func HpaDeleteV2beta2(kubeconfig, nameSpace, hpaName string) error {
	err := common.ClientSet(kubeconfig).AutoscalingV2beta2().HorizontalPodAutoscalers(nameSpace).Delete(context.TODO(), hpaName, metav1.DeleteOptions{})
	return err
}

/*
HpaYamlModifyV2beta2 - 通过YAML修改HPA (Modify HPA via YAML)
kubeconfig: 集群配置信息 (Cluster configuration)
yamlData: YAML格式的HPA配置数据 (YAML formatted HPA configuration data)
*/
func HpaYamlModifyV2beta2(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	hpa := &autoscalingv2beta2.HorizontalPodAutoscaler{}
	err = json.Unmarshal(data, hpa)
	if err != nil {
		return err
	}
	namespace := hpa.ObjectMeta.Namespace
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Update(context.TODO(), hpa, metav1.UpdateOptions{})
	return err
}
