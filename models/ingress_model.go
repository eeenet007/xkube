// ingress_model.go
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
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// Ingress - Ingress结构体，用于返回Ingress列表信息
// Ingress - Ingress struct for returning Ingress list information
type Ingress struct {
	IngressName  string `json:"ingressName"`  // Ingress名称 // Ingress name
	NameSpace    string `json:"nameSpace"`    // 命名空间 // Namespace
	Labels       string `json:"labels"`       // 标签 // Labels
	IngressClass string `json:"ingressClass"` // Ingress类 // Ingress class
	Rules        string `json:"rules"`        // 规则 // Rules
	Endpoint     string `json:"endpoint"`     //内部端点 // Internal endpoint
	CreateTime   string `json:"createTime"`   // 创建时间 // Create time
}

// IngressMain - Ingress主要信息结构体
// IngressMain - Ingress main information struct
type IngressRule struct {
	Host  string        `json:"host"`  // 主机名 // Host name
	Paths []IngHostPath `json:"paths"` // 路径列表 // Paths list
}

type IngressTlsHost struct {
	Tls   string   `json:"tls"`   // 证书名称
	Hosts []string `json:"hosts"` // 域名数组
}

// IngressDetail - Ingress详细信息结构体
// IngressDetail - Ingress detail information struct
type IngressDetail struct {
	IngressName string      `json:"ingressName"` // Ingress名称 // Ingress name
	NameSpace   string      `json:"nameSpace"`   // 命名空间 // Namespace
	Annotations string      `json:"annotations"` // 注解 // Annotations
	Endpoint    string      `json:"endpoint"`    // 端点 // Endpoint
	Labels      string      `json:"labels"`      // 标签 // Labels
	CreateTime  string      `json:"createTime"`  // 创建时间 // Create time
	Rules       []RulesPath `json:"rules"`       // 规则列表 // Rules list
}

// RulesPath - 规则路径结构体
// RulesPath - Rule path struct
type RulesPath struct {
	Host        string `json:"host"`        // 主机名 // Host name
	PathType    string `json:"pathType"`    // 路径类型 // Path type
	Path        string `json:"path"`        // 路径 // Path
	ServiceName string `json:"serviceName"` // 服务名称 // Service name
	ServicePort string `json:"servicePort"` // 服务端口 // Service port
}

// type KV struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }

// IngHostPath - Ingress主机路径结构体
// IngHostPath - Ingress host path struct
type IngHostPath struct {
	Path        string `json:"path"`        // 路径 // Path
	PathType    string `json:"pathType"`    // 路径类型 // Path type
	ServiceName string `json:"serviceName"` // 服务名称 // Service name
	ServicePort int32  `json:"servicePort"` // 服务端口 // Service port
}

// IngList - 获取Ingress列表
// IngList - Get Ingress list
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
func IngList(kubeconfig, namespace, ingressName, serviceName string, labelsKey, labelsValue string) ([]Ingress, error) {
	clientset := common.ClientSet(kubeconfig)
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue)}
	}

	ingList, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), listOptions)
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
				if v2.Backend.Service.Port.Number > 0 {
					portstr = fmt.Sprintf(":%d", v2.Backend.Service.Port.Number)
				}
				rulestr += fmt.Sprintf("%s%s --> %s%s,", v1.Host, v2.Path, v2.Backend.Service.Name, portstr)

				//用serviceName来搜索
				if serviceName != "" && v2.Backend.Service.Name == serviceName {
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

// IngCreate - 创建Ingress
// IngCreate - Create Ingress
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
func IngCreate(kubeconfig string, bodys []byte) error {
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

	var ingresPaths []networkingv1.HTTPIngressPath
	paths := gp.Get("paths").Array()
	for _, vv := range paths {
		var pathType networkingv1.PathType
		switch vv.Get("pathType").Str {
		case "ImplementationSpecific":
			pathType = networkingv1.PathTypeImplementationSpecific
		case "Exact":
			pathType = networkingv1.PathTypeExact
		default:
			pathType = networkingv1.PathTypePrefix
		}

		ingresPaths = append(ingresPaths, networkingv1.HTTPIngressPath{
			Path:     vv.Get("path").Str,
			PathType: &pathType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: vv.Get("serviceName").Str,
					Port: networkingv1.ServiceBackendPort{
						Number: int32(vv.Get("servicePort").Int()),
					},
				},
			},
		})
	}

	newIngress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ingressName,
			Namespace:   nameSpace,
			Labels:      labelsMap,
			Annotations: annotationsMap,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: ingressHost,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: ingresPaths,
						},
					},
				},
			},
		},
	}

	tlsCert := gp.Get("tlsCert").String()
	if tlsCert != "" {
		var ingTls = []networkingv1.IngressTLS{
			{
				Hosts:      []string{ingressHost},
				SecretName: tlsCert,
			},
		}
		newIngress.Spec.TLS = ingTls
	}

	// 创建Service
	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.NetworkingV1().Ingresses(nameSpace).Create(context.TODO(), newIngress, metav1.CreateOptions{})
	return err
}

// IngYamlCreate - 通过YAML创建Ingress
// IngYamlCreate - Create Ingress via YAML
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
func IngYamlCreate(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	ingress := &networkingv1.Ingress{}
	err = json.Unmarshal(data, ingress)
	if err != nil {
		return err
	}

	namespace := ingress.ObjectMeta.Namespace
	ingressName := ingress.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, ingressName)
	return err
}

// IngYamlModify - 通过YAML修改Ingress
// IngYamlModify - Modify Ingress via YAML
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
func IngYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	ingress := &networkingv1.Ingress{}
	err = json.Unmarshal(data, ingress)
	if err != nil {
		return err
	}

	namespace := ingress.ObjectMeta.Namespace
	ingressName := ingress.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, ingressName)
	return err
}

// IngUpdate - 更新Ingress
// IngUpdate - Update Ingress
// 参数:
//
//	kubeconfig: kubeconfig配置文件路径或集群ID
//	newIng: 新的Ingress信息
//
// Parameters:
//
//	kubeconfig: kubeconfig file path or cluster ID
//	newIng: new Ingress information
func IngUpdate(kubeconfig string, newIng *Ingress) {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(newIng.NameSpace)
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

// IngDetail - 获取Ingress详细信息
// IngDetail - Get Ingress detailed information
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
func IngDetail(kubeconfig, namespace, ingressName string) (*IngressDetail, error) {
	ing, err := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		//log.Println(err)
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
			if v2.Backend.Service.Port.Number > 0 {
				portstr = fmt.Sprintf(":%d", v2.Backend.Service.Port.Number)
			}
			rulestr += fmt.Sprintf("%s%s-->%s%s,", v1.Host, v2.Path, v2.Backend.Service.Name, portstr)

			xItems := &RulesPath{
				Host:        v1.Host,
				PathType:    fmt.Sprintf("%v", *v2.PathType),
				Path:        v2.Path,
				ServiceName: v2.Backend.Service.Name,
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

// GetIngYaml - 获取Ingress的YAML格式信息
// GetIngYaml - Get Ingress information in YAML format
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
func GetIngYaml(kubeconfig, namespace, ingressName string) (string, error) {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace)
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

// IngDelete - 删除Ingress
// IngDelete - Delete Ingress
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
func IngDelete(kubeconfig, namespace, ingressName string) error {
	//deletePolicy := metav1.DeletePropagationForeground
	if err := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func GetIngRule(kubeconfig, namespace, ingressName string) ([]IngressRule, error) {
	var bbb = make([]IngressRule, 0)
	ing, err := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return bbb, err
	}
	for _, v1 := range ing.Spec.Rules {
		var ppp = make([]IngHostPath, 0)
		for _, v2 := range v1.HTTP.Paths {
			ppp = append(ppp, IngHostPath{
				PathType:    fmt.Sprintf("%v", *v2.PathType),
				Path:        v2.Path,
				ServiceName: v2.Backend.Service.Name,
				ServicePort: v2.Backend.Service.Port.Number,
			})
		}
		bbb = append(bbb, IngressRule{
			Host:  v1.Host,
			Paths: ppp,
		})
	}
	return bbb, nil
}

func UpdateIngRule(kubeconfig, namespace, ingressName string, ingRule []IngressRule) error {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace)
	ing, err := ingClient.Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var ingRuls = make([]networkingv1.IngressRule, 0)
	for _, v1 := range ingRule {
		var ingresPaths = make([]networkingv1.HTTPIngressPath, 0)
		for _, v2 := range v1.Paths {
			var pathType networkingv1.PathType
			switch v2.PathType {
			case "ImplementationSpecific":
				pathType = networkingv1.PathTypeImplementationSpecific
			case "Exact":
				pathType = networkingv1.PathTypeExact
			default:
				pathType = networkingv1.PathTypePrefix
			}
			ingresPaths = append(ingresPaths, networkingv1.HTTPIngressPath{
				Path:     v2.Path,
				PathType: &pathType,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: v2.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Number: v2.ServicePort,
						},
					},
				},
			})
		}

		ingRuls = append(ingRuls, networkingv1.IngressRule{
			Host: v1.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: ingresPaths,
				},
			},
		})
	}
	ing.Spec.Rules = ingRuls
	_, err1 := ingClient.Update(context.TODO(), ing, metav1.UpdateOptions{})
	return err1
}

func GetIngTlsHost(kubeconfig, namespace, ingressName string) ([]IngressTlsHost, error) {
	var bbb = make([]IngressTlsHost, 0)
	ing, err := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
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

func UpdateIngTlsHost(kubeconfig, namespace, ingressName string, tlsHost []IngressTlsHost) error {
	ingClient := common.ClientSet(kubeconfig).NetworkingV1().Ingresses(namespace)
	ing, err := ingClient.Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	var ingTls = make([]networkingv1.IngressTLS, 0)
	for _, vv := range tlsHost {
		ingTls = append(ingTls, networkingv1.IngressTLS{
			Hosts:      vv.Hosts,
			SecretName: vv.Tls,
		})
	}
	ing.Spec.TLS = ingTls
	_, err1 := ingClient.Update(context.TODO(), ing, metav1.UpdateOptions{})
	return err1
}
