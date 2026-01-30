// ingress_model_v1beta1.go
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"xkube/common"

	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// IngListV1beta1 - 获取Ingress列表 (使用networking/v1beta1 API版本)
// IngListV1beta1 - Get Ingress list (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径
//	namespace: 命名空间
//	ingressName: Ingress名称，用于过滤
//	serviceName: 服务名称，用于过滤
//	labelsKey: 标签键
//	labelsValue: 标签值
//
// Parameters:
//
//	kubeconfig: kubeconfig file path
//	namespace: namespace
//	ingressName: ingress name, used for filtering
//	serviceName: service name, used for filtering
//	labelsKey: label key
//	labelsValue: label value
//
// 返回值:
//
//	[]Ingress: Ingress列表
//	error: 错误信息
//
// Return value:
//
//	[]Ingress: Ingress list
//	error: error message
func IngListV1beta1(kubeconfig, namespace, ingressName, serviceName string, labelsKey, labelsValue string) ([]Ingress, error) {
	clientset := common.ClientSet(kubeconfig)
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue)}
	}

	ingList, err := clientset.NetworkingV1beta1().Ingresses(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list service error:%v\n", err)
	}

	var bbb = make([]Ingress, 0)
	for _, ing := range ingList.Items {
		//搜索
		if ingressName != "" {
			if !strings.Contains(ing.Name, ingressName) {
				continue
			}
		}

		var labelsStr string
		for kk, vv := range ing.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		var rulestr string
		var isNeedService bool
		for _, v1 := range ing.Spec.Rules {
			for _, v2 := range v1.HTTP.Paths {
				var portstr string
				if v2.Backend.ServicePort.IntVal > 0 {
					portstr = fmt.Sprintf("%d", v2.Backend.ServicePort.IntVal)
				}
				rulestr += fmt.Sprintf("%s%s --> %s%s,", v1.Host, v2.Path, v2.Backend.ServiceName, portstr)

				//用serviceName来搜索
				if serviceName != "" && v2.Backend.ServiceName == serviceName {
					isNeedService = true
				}
			}
		}
		//用serviceName来搜索
		if serviceName != "" && !isNeedService {
			continue
		}

		var endpointIp string
		if len(ing.Status.LoadBalancer.Ingress) > 0 {
			endpointIp = ing.Status.LoadBalancer.Ingress[0].IP
		}

		Items := &Ingress{
			IngressName:  ing.Name,
			NameSpace:    ing.Namespace,
			IngressClass: ing.ObjectMeta.Annotations["kubernetes.io/ingress.class"],
			Rules:        rulestr,
			Endpoint:     endpointIp,
			Labels:       labelsStr,
			CreateTime:   ing.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// IngCreateV1beta1 - 创建Ingress (使用networking/v1beta1 API版本)
// IngCreateV1beta1 - Create Ingress (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	bodys: 包含Ingress信息的JSON数据
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	bodys: JSON data containing Ingress information
//
// 返回值:
//
//	error: 错误信息
//
// Return value:
//
//	error: error message
func IngCreateV1beta1(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	ingressName := gp.Get("ingressName").String()
	nameSpace := gp.Get("nameSpace").String()
	ingressHost := gp.Get("ingressHost").String()

	var labelsMap = make(map[string]string)
	labelsMap["app"] = ingressName
	for _, vv := range gp.Get("labels").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	var annotationsMap = make(map[string]string)
	annotationsMap["kubernetes.io/ingress.class"] = "nginx"

	var ingresPaths []networkingv1beta1.HTTPIngressPath
	paths := gp.Get("paths").Array()
	for _, vv := range paths {
		var pathType networkingv1beta1.PathType
		switch vv.Get("pathType").Str {
		case "ImplementationSpecific":
			pathType = networkingv1beta1.PathTypeImplementationSpecific
		case "Exact":
			pathType = networkingv1beta1.PathTypeExact
		default:
			pathType = networkingv1beta1.PathTypePrefix
		}

		ingresPaths = append(ingresPaths, networkingv1beta1.HTTPIngressPath{
			Path:     vv.Get("path").Str,
			PathType: &pathType,
			Backend: networkingv1beta1.IngressBackend{
				ServiceName: vv.Get("serviceName").Str,
				ServicePort: intstr.FromInt32(int32(vv.Get("servicePort").Int())),
			},
		})
	}

	newIngress := &networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ingressName,
			Namespace:   nameSpace,
			Labels:      labelsMap,
			Annotations: annotationsMap,
		},
		Spec: networkingv1beta1.IngressSpec{
			Rules: []networkingv1beta1.IngressRule{
				{
					Host: ingressHost,
					IngressRuleValue: networkingv1beta1.IngressRuleValue{
						HTTP: &networkingv1beta1.HTTPIngressRuleValue{
							Paths: ingresPaths,
						},
					},
				},
			},
		},
	}

	tlsCert := gp.Get("tlsCert").String()
	if tlsCert != "" {
		var ingTls = []networkingv1beta1.IngressTLS{
			{
				Hosts:      []string{ingressHost},
				SecretName: tlsCert,
			},
		}
		newIngress.Spec.TLS = ingTls
	}

	// 创建Service
	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.NetworkingV1beta1().Ingresses(nameSpace).Create(context.TODO(), newIngress, metav1.CreateOptions{})
	return err
}

// IngYamlCreateV1beta1 - 通过YAML创建Ingress (使用networking/v1beta1 API版本)
// IngYamlCreateV1beta1 - Create Ingress via YAML (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	yamlData: 包含Ingress信息的YAML数据
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	yamlData: YAML data containing Ingress information
//
// 返回值:
//
//	error: 错误信息
//
// Return value:
//
//	error: error message
func IngYamlCreateV1beta1(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	ingress := &networkingv1beta1.Ingress{}
	err = json.Unmarshal(data, ingress)
	if err != nil {
		return err
	}

	namespace := ingress.ObjectMeta.Namespace
	ingressName := ingress.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.NetworkingV1beta1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, ingressName)
	return err
}

// IngYamlModifyV1beta1 - 通过YAML修改Ingress (使用networking/v1beta1 API版本)
// IngYamlModifyV1beta1 - Modify Ingress via YAML (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	yamlData: 包含Ingress信息的YAML数据
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	yamlData: YAML data containing Ingress information
//
// 返回值:
//
//	error: 错误信息
//
// Return value:
//
//	error: error message
func IngYamlModifyV1beta1(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	ingress := &networkingv1beta1.Ingress{}
	err = json.Unmarshal(data, ingress)
	if err != nil {
		return err
	}

	namespace := ingress.ObjectMeta.Namespace
	ingressName := ingress.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.NetworkingV1beta1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, ingressName)
	return err
}

// IngUpdateV1beta1 - 更新Ingress (使用networking/v1beta1 API版本)
// IngUpdateV1beta1 - Update Ingress (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	newIng: 新的Ingress信息
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	newIng: new Ingress information
func IngUpdateV1beta1(kubeconfig string, newIng *Ingress) {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(newIng.NameSpace)
	ingress, err := ingClient.Get(context.TODO(), newIng.IngressName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	//ingress.Spec.Rules[0].Host = "new-example.com"
	_, updateErr := ingClient.Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if updateErr != nil {

		// 处理错误
	}
}

// func IngDetail_old(kubeconfig, namespace, ingressName string) (*networkingv1.Ingress, error) {
// 	ingress, err := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return ingress, nil
// }

// IngDetailV1beta1 - 获取Ingress详细信息 (使用networking/v1beta1 API版本)
// IngDetailV1beta1 - Get Ingress detailed information (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	namespace: 命名空间
//	ingressName: Ingress名称
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	namespace: namespace
//	ingressName: Ingress name
//
// 返回值:
//
//	*IngressDetail: Ingress详细信息
//	error: 错误信息
//
// Return value:
//
//	*IngressDetail: Ingress detail information
//	error: error message
func IngDetailV1beta1(kubeconfig, namespace, ingressName string) (*IngressDetail, error) {
	ing, err := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return &IngressDetail{}, err
	}

	var labelsStr string
	for kk, vv := range ing.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	var annotationsStr string
	for kk, vv := range ing.ObjectMeta.Annotations {
		if kk != "kubectl.kubernetes.io/last-applied-configuration" {
			annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	var rulestr string
	var bbb = make([]RulesPath, 0)
	for _, v1 := range ing.Spec.Rules {
		for _, v2 := range v1.HTTP.Paths {
			var portstr string
			if v2.Backend.ServicePort.IntVal > 0 {
				portstr = fmt.Sprintf(":%d", v2.Backend.ServicePort.IntVal)
			}
			rulestr += fmt.Sprintf("%s%s-->%s%s,", v1.Host, v2.Path, v2.Backend.ServiceName, portstr)
			xItems := &RulesPath{
				Host:        v1.Host,
				PathType:    "",
				Path:        v2.Path,
				ServiceName: v2.Backend.ServiceName,
				ServicePort: portstr,
			}
			bbb = append(bbb, *xItems)

		}
	}

	var endpointIp string
	if len(ing.Status.LoadBalancer.Ingress) > 0 {
		endpointIp = ing.Status.LoadBalancer.Ingress[0].IP
	}

	return &IngressDetail{
		IngressName: ing.Name,
		NameSpace:   ing.Namespace,
		Annotations: annotationsStr,
		Endpoint:    endpointIp,
		Labels:      labelsStr,
		CreateTime:  ing.CreationTimestamp.Format("2006-01-02 15:04:05"),
		Rules:       bbb,
	}, nil
}

// GetIngYamlV1beta1 - 获取Ingress的YAML格式信息 (使用networking/v1beta1 API版本)
// GetIngYamlV1beta1 - Get Ingress information in YAML format (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	namespace: 命名空间
//	ingressName: Ingress名称
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	namespace: namespace
//	ingressName: Ingress name
//
// 返回值:
//
//	string: YAML格式的Ingress信息
//	error: 错误信息
//
// Return value:
//
//	string: Ingress information in YAML format
//	error: error message
func GetIngYamlV1beta1(kubeconfig, namespace, ingressName string) (string, error) {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace)
	ingress, err := ingClient.Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	ingresseUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ingress)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(ingresseUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// IngDeleteV1beta1 - 删除Ingress (使用networking/v1beta1 API版本)
// IngDeleteV1beta1 - Delete Ingress (using networking/v1beta1 API version)
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	namespace: 命名空间
//	ingressName: Ingress名称
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	namespace: namespace
//	ingressName: Ingress name
//
// 返回值:
//
//	error: 错误信息
//
// Return value:
//
//	error: error message
//
// 删除Service
func IngDeleteV1beta1(kubeconfig, namespace, ingressName string) error {
	//deletePolicy := metav1.DeletePropagationForeground
	if err := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func GetIngRuleV1beta1(kubeconfig, namespace, ingressName string) ([]IngressRule, error) {
	var bbb = make([]IngressRule, 0)
	ing, err := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return bbb, err
	}
	for _, v1 := range ing.Spec.Rules {
		var ppp = make([]IngHostPath, 0)
		for _, v2 := range v1.HTTP.Paths {
			ppp = append(ppp, IngHostPath{
				PathType:    fmt.Sprintf("%v", *v2.PathType),
				Path:        v2.Path,
				ServiceName: v2.Backend.ServiceName,
				ServicePort: v2.Backend.ServicePort.IntVal,
			})
		}
		bbb = append(bbb, IngressRule{
			Host:  v1.Host,
			Paths: ppp,
		})
	}
	return bbb, nil
}

func UpdateIngRuleV1beta1(kubeconfig, namespace, ingressName string, ingRule []IngressRule) error {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace)
	ing, err := ingClient.Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var ingRuls = make([]networkingv1beta1.IngressRule, 0)
	for _, v1 := range ingRule {
		var ingresPaths = make([]networkingv1beta1.HTTPIngressPath, 0)
		for _, v2 := range v1.Paths {
			var pathType networkingv1beta1.PathType
			switch v2.PathType {
			case "ImplementationSpecific":
				pathType = networkingv1beta1.PathTypeImplementationSpecific
			case "Exact":
				pathType = networkingv1beta1.PathTypeExact
			default:
				pathType = networkingv1beta1.PathTypePrefix
			}
			ingresPaths = append(ingresPaths, networkingv1beta1.HTTPIngressPath{
				Path:     v2.Path,
				PathType: &pathType,
				Backend: networkingv1beta1.IngressBackend{
					ServiceName: v2.ServiceName,
					ServicePort: intstr.FromInt32(v2.ServicePort),
				},
			})
		}

		ingRuls = append(ingRuls, networkingv1beta1.IngressRule{
			Host: v1.Host,
			IngressRuleValue: networkingv1beta1.IngressRuleValue{
				HTTP: &networkingv1beta1.HTTPIngressRuleValue{
					Paths: ingresPaths,
				},
			},
		})
	}
	ing.Spec.Rules = ingRuls
	_, err1 := ingClient.Update(context.TODO(), ing, metav1.UpdateOptions{})
	return err1
}

func GetIngTlsHostV1beta1(kubeconfig, namespace, ingressName string) ([]IngressTlsHost, error) {
	var bbb = make([]IngressTlsHost, 0)
	ing, err := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return bbb, err
	}
	for _, v1 := range ing.Spec.TLS {
		bbb = append(bbb, IngressTlsHost{
			Tls:   v1.SecretName,
			Hosts: v1.Hosts,
		})
	}
	return bbb, nil
}

func UpdateIngTlsHostV1beta1(kubeconfig, namespace, ingressName string, tlsHost []IngressTlsHost) error {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1beta1().Ingresses(namespace)
	ing, err := ingClient.Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var ingTls = make([]networkingv1beta1.IngressTLS, 0)
	for _, vv := range tlsHost {
		ingTls = append(ingTls, networkingv1beta1.IngressTLS{
			Hosts:      vv.Hosts,
			SecretName: vv.Tls,
		})
	}
	ing.Spec.TLS = ingTls
	_, err1 := ingClient.Update(context.TODO(), ing, metav1.UpdateOptions{})
	return err1
}
