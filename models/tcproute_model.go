// tcproute_model.go
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

// TCPRoute TCPRoute列表信息结构体
type TCPRoute struct {
	TCPRouteName string `json:"tcprouteName"`
	NameSpace    string `json:"nameSpace"`
	ParentRefs   string `json:"parentRefs"`
	Rules        string `json:"rules"`
	CreateTime   string `json:"createTime"`
}

// TCPRouteDetail TCPRoute详细信息结构体
type TCPRouteDetail struct {
	TCPRoute
	//ParentRef string `json:"parentRef"`
	Rules string `json:"rules"`
}

// TCPRouteList 获取TCPRoute列表
func TCPRouteList(kubeconfig string) ([]TCPRoute, error) {
	var tcproutes []TCPRoute
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return tcproutes, err
	}

	tcprouteList, err := clientset.GatewayV1alpha2().TCPRoutes("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list tcproutes: %v\n", err)
		return tcproutes, err
	}

	if len(tcprouteList.Items) == 0 {
		log.Printf("[WARN] tcprouteList NoFound RequestedResource")
		return tcproutes, fmt.Errorf("NoFound RequestedResource")
	}

	for _, tr := range tcprouteList.Items {
		// 获取parentRefs信息
		var parentRefs []string
		for _, parentRef := range tr.Spec.ParentRefs {
			ref := string(parentRef.Name)
			if parentRef.SectionName != nil {
				ref += "/" + string(*parentRef.SectionName)
			}
			parentRefs = append(parentRefs, ref)
		}
		parentRefsStr := strings.Join(parentRefs, ",")

		// 获取rules信息
		var rules []string
		for _, rule := range tr.Spec.Rules {
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

		tcproute := TCPRoute{
			TCPRouteName: tr.Name,
			NameSpace:    tr.Namespace,
			ParentRefs:   parentRefsStr,
			Rules:        rulesStr,
			CreateTime:   tr.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		tcproutes = append(tcproutes, tcproute)
	}

	return tcproutes, nil
}

// GetTCPRouteDetail 获取TCPRoute详细信息
func GetTCPRouteDetail(kubeconfig, nameSpace, tcprouteName string) (*TCPRouteDetail, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	tr, err := clientset.GatewayV1alpha2().TCPRoutes(nameSpace).Get(context.TODO(), tcprouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get tcproute %s: %v\n", tcprouteName, err)
		return nil, err
	}

	// 获取parentRefs信息
	var parentRefs []string
	for _, parentRef := range tr.Spec.ParentRefs {
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
	for _, rule := range tr.Spec.Rules {
		// 处理backendRefs
		for _, backendRef := range rule.BackendRefs {
			backend := fmt.Sprintf("%s,%s,%d,%d", *backendRef.Kind, backendRef.Name, *backendRef.Port, *backendRef.Weight)
			rules = append(rules, backend)
		}
	}
	rulesStr := strings.Join(rules, ";")

	// 获取parentRef (第一个)
	// parentRef := ""
	// if len(tr.Spec.ParentRefs) > 0 {
	// 	parentRef = string(tr.Spec.ParentRefs[0].Name)
	// 	if tr.Spec.ParentRefs[0].SectionName != nil {
	// 		parentRef += "/" + string(*tr.Spec.ParentRefs[0].SectionName)
	// 	}
	// }

	tcproute := TCPRoute{
		TCPRouteName: tr.Name,
		NameSpace:    tr.Namespace,
		ParentRefs:   parentRefsStr,
		Rules:        rulesStr,
		CreateTime:   tr.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}

	detail := &TCPRouteDetail{
		TCPRoute: tcproute,
		//ParentRef: parentRef,
		Rules: rulesStr,
	}

	return detail, nil
}

// DeleteTCPRoute 删除TCPRoute
func DeleteTCPRoute(kubeconfig, nameSpace, tcprouteName string) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	err = clientset.GatewayV1alpha2().TCPRoutes(nameSpace).Delete(context.TODO(), tcprouteName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete tcproute %s: %v\n", tcprouteName, err)
		return err
	}

	return nil
}

// CreateTCPRoute 创建TCPRoute
func CreateTCPRoute(kubeconfig string, bodys []byte) error {

	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	tcprouteName := gp.Get("tcprouteName").String()
	gatewayName := gp.Get("gatewayName").String()
	nameSpace := gp.Get("nameSpace").String()
	sectionName := gp.Get("sectionName").String()

	var labelsMap = make(map[string]string)
	for _, vv := range gp.Get("labels").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	rules := make([]gatewayv1alpha2.TCPRouteRule, 0, len(gp.Get("rules").Array()))
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
		var backendRefs []gatewayv1alpha2.BackendRef
		backendRefs = append(backendRefs, backendRef)
		rules = append(rules, gatewayv1alpha2.TCPRouteRule{
			BackendRefs: backendRefs,
		})
	}

	// 构建ParentRef
	parentRef := gatewayv1.ParentReference{
		Name: gatewayv1.ObjectName(gatewayName),
	}
	if sectionName != "" {
		section := gatewayv1.SectionName(sectionName)
		parentRef.SectionName = &section
	}

	// 构建TCPRoute对象
	tcproute := &gatewayv1alpha2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tcprouteName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: gatewayv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{parentRef},
			},
			Rules: rules,
		},
	}

	_, err = clientset.GatewayV1alpha2().TCPRoutes(nameSpace).Create(context.TODO(), tcproute, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create tcproute %s: %v\n", tcprouteName, err)
		return err
	}

	return nil
}

// UpdateTCPRouteByYaml 通过YAML更新TCPRoute
func UpdateTCPRouteByYaml(kubeconfig string, yamlData []byte) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 解析YAML
	tcproute := &gatewayv1alpha2.TCPRoute{}
	if err := yaml.Unmarshal(yamlData, tcproute); err != nil {
		log.Printf("Failed to unmarshal yaml: %v\n", err)
		return err
	}

	tr, err := clientset.GatewayV1alpha2().TCPRoutes(tcproute.Namespace).Get(context.TODO(), tcproute.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get tcproute %s: %v\n", tcproute.Name, err)
		return err
	}

	tcproute.ObjectMeta.ResourceVersion = tr.ObjectMeta.ResourceVersion

	// 更新TCPRoute
	_, err = clientset.GatewayV1alpha2().TCPRoutes(tcproute.Namespace).Update(context.TODO(), tcproute, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update tcproute %s: %v\n", tcproute.Name, err)
		return err
	}

	return nil
}

// GetTCPRouteYaml 获取TCPRoute的YAML配置
func GetTCPRouteYaml(kubeconfig, nameSpace, tcprouteName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return "", err
	}

	tr, err := clientset.GatewayV1alpha2().TCPRoutes(nameSpace).Get(context.TODO(), tcprouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get tcproute %s: %v\n", tcprouteName, err)
		return "", err
	}

	trUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(tr)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(trUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
