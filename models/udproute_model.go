// udproute_model.go
package models

import (
	"context"
	"fmt"
	"log"
	"strings"
	"xkube/common"

	"github.com/tidwall/gjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gatewayClient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

// UDPRoute UDPRoute列表信息结构体
type UDPRoute struct {
	UDPRouteName string `json:"udprouteName"`
	NameSpace    string `json:"nameSpace"`
	ParentRefs   string `json:"parentRefs"`
	Rules        string `json:"rules"`
	CreateTime   string `json:"createTime"`
}

// UDPRouteDetail UDPRoute详细信息结构体
type UDPRouteDetail struct {
	UDPRoute
	//ParentRef string `json:"parentRef"`
	Rules string `json:"rules"`
}

// UDPRouteList 获取UDPRoute列表
func UDPRouteList(kubeconfig string) ([]UDPRoute, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	udprouteList, err := clientset.GatewayV1alpha2().UDPRoutes("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list udproutes: %v\n", err)
		return nil, err
	}

	if len(udprouteList.Items) == 0 {
		log.Printf("[WARN] udprouteList NoFound RequestedResource")
		return nil, fmt.Errorf("NoFound RequestedResource")
	}

	var udproutes []UDPRoute
	for _, ur := range udprouteList.Items {
		// 获取parentRefs信息
		var parentRefs []string
		for _, parentRef := range ur.Spec.ParentRefs {
			ref := string(parentRef.Name)
			if parentRef.SectionName != nil {
				ref += "/" + string(*parentRef.SectionName)
			}
			parentRefs = append(parentRefs, ref)
		}
		parentRefsStr := strings.Join(parentRefs, ",")

		// 获取rules信息
		var rules []string
		for _, rule := range ur.Spec.Rules {
			// 处理backendRefs
			var backendRefs []string
			for _, backendRef := range rule.BackendRefs {
				backend := fmt.Sprintf("%s:%d,weight:%d", backendRef.Name, *backendRef.Port, *backendRef.Weight)
				backendRefs = append(backendRefs, backend)
			}
			backendRefStr := strings.Join(backendRefs, ";")
			rules = append(rules, backendRefStr)
		}
		rulesStr := strings.Join(rules, "|")

		udproute := UDPRoute{
			UDPRouteName: ur.Name,
			NameSpace:    ur.Namespace,
			ParentRefs:   parentRefsStr,
			Rules:        rulesStr,
			CreateTime:   ur.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		udproutes = append(udproutes, udproute)
	}

	return udproutes, nil
}

// GetUDPRouteDetail 获取UDPRoute详细信息
func GetUDPRouteDetail(kubeconfig, nameSpace, udprouteName string) (*UDPRouteDetail, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	ur, err := clientset.GatewayV1alpha2().UDPRoutes(nameSpace).Get(context.TODO(), udprouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get udproute %s: %v\n", udprouteName, err)
		return nil, err
	}

	// 获取parentRefs信息
	var parentRefs []string
	for _, parentRef := range ur.Spec.ParentRefs {
		var ref string
		var sectionName string
		//ref := string(parentRef.Name)
		if parentRef.SectionName != nil {
			//ref += "/" + string(*parentRef.SectionName)
			sectionName = string(*parentRef.SectionName)
		}
		ref = fmt.Sprintf("group:%s,kind:%s,name:%s,sectionName:%s", *parentRef.Group, *parentRef.Kind, string(parentRef.Name), sectionName)
		parentRefs = append(parentRefs, ref)
	}
	parentRefsStr := strings.Join(parentRefs, ";")
	parentRefsStr = parentRefsStr[0 : len(parentRefsStr)-1]

	// 获取rules信息
	var rules []string
	for _, rule := range ur.Spec.Rules {
		// 处理backendRefs
		for _, backendRef := range rule.BackendRefs {
			backend := fmt.Sprintf("%s,%s,%d,%d", *backendRef.Kind, backendRef.Name, *backendRef.Port, *backendRef.Weight)
			rules = append(rules, backend)
		}
	}
	rulesStr := strings.Join(rules, ";")

	// 获取parentRef (第一个)
	// parentRef := ""
	// if len(ur.Spec.ParentRefs) > 0 {
	// 	parentRef = string(ur.Spec.ParentRefs[0].Name)
	// 	if ur.Spec.ParentRefs[0].SectionName != nil {
	// 		parentRef += "/" + string(*ur.Spec.ParentRefs[0].SectionName)
	// 	}
	// }

	udproute := UDPRoute{
		UDPRouteName: ur.Name,
		NameSpace:    ur.Namespace,
		ParentRefs:   parentRefsStr,
		Rules:        rulesStr,
		CreateTime:   ur.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}

	detail := &UDPRouteDetail{
		UDPRoute: udproute,
		//ParentRef: parentRef,
		Rules: rulesStr,
	}

	return detail, nil
}

// DeleteUDPRoute 删除UDPRoute
func DeleteUDPRoute(kubeconfig, nameSpace, udprouteName string) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	err = clientset.GatewayV1alpha2().UDPRoutes(nameSpace).Delete(context.TODO(), udprouteName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete udproute %s: %v\n", udprouteName, err)
		return err
	}

	return nil
}

// CreateUDPRoute 创建UDPRoute
func CreateUDPRoute(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	udprouteName := gp.Get("udprouteName").String()
	gatewayName := gp.Get("gatewayName").String()
	nameSpace := gp.Get("nameSpace").String()
	sectionName := gp.Get("sectionName").String()

	labelsMap := make(map[string]string)
	for _, vv := range gp.Get("labels").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	rules := make([]gatewayv1alpha2.UDPRouteRule, 0, len(gp.Get("rules").Array()))
	for _, vv := range gp.Get("rules").Array() {
		// 处理后端服务
		serviceName := vv.Get("serviceName").Str
		port := gatewayv1.PortNumber(vv.Get("port").Int())
		weightVal := int32(vv.Get("weight").Int())

		backendRef := gatewayv1alpha2.BackendRef{
			BackendObjectReference: gatewayv1alpha2.BackendObjectReference{
				Name: gatewayv1alpha2.ObjectName(serviceName),
				Port: &port,
			},
			Weight: &weightVal,
		}

		rule := gatewayv1alpha2.UDPRouteRule{
			BackendRefs: []gatewayv1alpha2.BackendRef{backendRef},
		}
		rules = append(rules, rule)
	}

	// 构建ParentRef
	parentRef := gatewayv1.ParentReference{
		Name: gatewayv1.ObjectName(gatewayName),
	}
	if sectionName != "" {
		section := gatewayv1.SectionName(sectionName)
		parentRef.SectionName = &section
	}

	// 构建UDPRoute对象
	udproute := &gatewayv1alpha2.UDPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      udprouteName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: gatewayv1alpha2.UDPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{parentRef},
			},
			Rules: rules,
		},
	}

	_, err = clientset.GatewayV1alpha2().UDPRoutes(nameSpace).Create(context.TODO(), udproute, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create udproute %s: %v\n", udprouteName, err)
		return err
	}

	return nil
}

// UpdateUDPRouteByYaml 通过YAML更新UDPRoute
func UpdateUDPRouteByYaml(kubeconfig string, yamlData []byte) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 解析YAML
	udproute := &gatewayv1alpha2.UDPRoute{}
	if err := yaml.Unmarshal(yamlData, udproute); err != nil {
		log.Printf("Failed to unmarshal yaml: %v\n", err)
		return err
	}

	ur, err := clientset.GatewayV1alpha2().UDPRoutes(udproute.Namespace).Get(context.TODO(), udproute.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get udproute %s: %v\n", udproute.Name, err)
		return err
	}

	udproute.ObjectMeta.ResourceVersion = ur.ObjectMeta.ResourceVersion

	// 更新UDPRoute
	_, err = clientset.GatewayV1alpha2().UDPRoutes(udproute.Namespace).Update(context.TODO(), udproute, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update udproute %s: %v\n", udproute.Name, err)
		return err
	}

	return nil
}

// GetUDPRouteYaml 获取UDPRoute的YAML配置
func GetUDPRouteYaml(kubeconfig, nameSpace, udprouteName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return "", err
	}

	ur, err := clientset.GatewayV1alpha2().UDPRoutes(nameSpace).Get(context.TODO(), udprouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get udproute %s: %v\n", udprouteName, err)
		return "", err
	}

	urUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(ur)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(urUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
