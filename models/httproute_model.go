// httproute_model.go
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"xkube/common"

	"github.com/tidwall/gjson"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	yamlutil "k8s.io/apimachinery/pkg/util/yaml"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayClient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

// HTTPRoute HTTPRoute列表信息结构体
type HTTPRoute struct {
	HTTPRouteName string `json:"httprouteName"`
	NameSpace     string `json:"nameSpace"`
	ParentRefs    string `json:"parentRefs"`
	Hostnames     string `json:"hostnames"`
	Rules         string `json:"rules"`
	//BackendRefs   string `json:"backendRefs"`
	CreateTime string `json:"createTime"`
}

// HTTPRouteDetail HTTPRoute详细信息结构体
type HTTPRouteDetail struct {
	HTTPRouteName  string `json:"httprouteName"`
	NameSpace      string `json:"nameSpace"`
	ParentRefs     string `json:"parentRefs"`
	Hostnames      string `json:"hostnames"`
	CreateTime     string `json:"createTime"`
	ControllerName string `json:"controllerName"`
	//ParentRef      string `json:"parentRef"`
	//BackendRefs    string               `json:"backendRefs"`
	Conditions []HTTPRouteCondition `json:"conditions"`
	Rules      []RouteRules         `json:"rules"`
}

type RouteRules struct {
	MatchPath   string `json:"matchPath"`
	Headers     string `json:"headers"`
	BackendRefs string `json:"backendRefs"`
}

// HTTPRouteCondition HTTPRoute条件信息结构体
type HTTPRouteCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

// HTTPRouteList 获取HTTPRoute列表
func HTTPRouteList(kubeconfig string) ([]HTTPRoute, error) {

	var httproutes []HTTPRoute
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return httproutes, err
	}

	httprouteList, err := clientset.GatewayV1().HTTPRoutes("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list httproutes: %v\n", err)
		return httproutes, err
	}

	if len(httprouteList.Items) == 0 {
		log.Printf("[WARN] httproutes NoFound RequestedResource")
		return httproutes, fmt.Errorf("NoFound RequestedResource")
	}

	for _, hr := range httprouteList.Items {
		// 获取parentRefs信息
		var parentRefs []string
		for _, parentRef := range hr.Spec.ParentRefs {
			parentRefs = append(parentRefs, string(parentRef.Name))
		}
		parentRefsStr := strings.Join(parentRefs, ",")

		// 获取hostnames信息
		var hostnames []string
		for _, hostname := range hr.Spec.Hostnames {
			hostnames = append(hostnames, string(hostname))
		}
		hostnamesStr := strings.Join(hostnames, ",")

		// 获取rules信息
		//var rules []string
		var rulesStr string
		for _, rule := range hr.Spec.Rules {
			var paths string
			var headers string
			// 处理matches
			for _, match := range rule.Matches {
				// path类型
				if match.Path != nil {
					//pathType := ""
					//if match.Path.Type != nil {
					//	pathType = string(*match.Path.Type)
					//}
					pathValue := ""
					if match.Path.Value != nil {
						pathValue = *match.Path.Value
					}
					//rules = append(rules, pathType+":"+pathValue)
					paths += fmt.Sprintf("%s,", pathValue)
				}
				// headers类型
				for _, header := range match.Headers {
					//headerType := ""
					// 使用fmt.Sprintf处理HeaderMatchType转换
					//headerType = fmt.Sprintf("%v", *header.Type)
					headerName := string(header.Name)
					headerValue := header.Value
					//rules = append(rules, headerType+":"+headerName+":"+headerValue)
					headers += fmt.Sprintf("%s:%s,", headerName, headerValue)
				}
			}
			if len(paths) > 0 {
				paths = paths[0 : len(paths)-1]
			}
			if len(headers) > 0 {
				headers = headers[0 : len(headers)-1]
			}

			// 处理backendRefs
			var backendRefs string
			for _, backendRef := range rule.BackendRefs {
				//kind := "Service" // 默认为Service
				//if backendRef.Kind != nil {
				//	kind = string(*backendRef.Kind)
				//}
				service := string(backendRef.Name)
				port := ""
				if backendRef.Port != nil {
					port = fmt.Sprintf("%d", *backendRef.Port)
				}
				weight := "1" // 默认权重为1
				if backendRef.Weight != nil {
					weight = fmt.Sprintf("%d", *backendRef.Weight)
				}
				//backendInfo := fmt.Sprintf("%s:%s,port:%s,weight:%s", kind, service, port, weight)
				backendRefs += fmt.Sprintf("%s:%s,weight:%s;", service, port, weight)
				//rules = append(rules, backendInfo)
			}
			if len(backendRefs) > 0 {
				backendRefs = backendRefs[0 : len(backendRefs)-1]
			}
			if headers != "" {
				rulesStr += fmt.Sprintf("[%s]-->header[%s]-->[%s]|", paths, headers, backendRefs)
			} else {
				rulesStr += fmt.Sprintf("[%s]-->[%s]|", paths, backendRefs)
			}
		}
		//rulesStr := strings.Join(rules, ",")
		if len(rulesStr) > 0 {
			rulesStr = rulesStr[0 : len(rulesStr)-1]
		}
		//log.Println(rulesStr)
		httproute := HTTPRoute{
			HTTPRouteName: hr.Name,
			NameSpace:     hr.Namespace,
			ParentRefs:    parentRefsStr,
			Hostnames:     hostnamesStr,
			Rules:         rulesStr,
			//BackendRefs:   backendRefsStr,
			CreateTime: hr.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		httproutes = append(httproutes, httproute)
	}
	return httproutes, nil
}

// GetHTTPRouteDetail 获取HTTPRoute详细信息
func GetHTTPRouteDetail(kubeconfig, nameSpace, httprouteName string) (*HTTPRouteDetail, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	hr, err := clientset.GatewayV1().HTTPRoutes(nameSpace).Get(context.TODO(), httprouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get httproute %s: %v\n", httprouteName, err)
		return nil, err
	}

	// 获取parentRefs信息
	var parentRefs []string
	for _, parentRef := range hr.Spec.ParentRefs {
		parentRefs = append(parentRefs, string(parentRef.Name))
	}
	parentRefsStr := strings.Join(parentRefs, ",")

	// 获取hostnames信息
	var hostnames []string
	for _, hostname := range hr.Spec.Hostnames {
		hostnames = append(hostnames, string(hostname))
	}
	hostnamesStr := strings.Join(hostnames, ",")

	// 获取rules信息
	var rules []RouteRules
	for _, rule := range hr.Spec.Rules {
		var paths string
		var headers string
		var backendRefs string
		// 处理matches
		for _, match := range rule.Matches {
			// path类型
			if match.Path != nil {
				pathType := ""
				if match.Path.Type != nil {
					pathType = string(*match.Path.Type)
				}
				pathValue := ""
				if match.Path.Value != nil {
					pathValue = *match.Path.Value
				}
				paths += fmt.Sprintf("%s:%s,", pathType, pathValue)
			}
			// headers类型
			for _, header := range match.Headers {
				headerType := ""
				// 使用fmt.Sprintf处理HeaderMatchType转换
				headerType = fmt.Sprintf("%v", *header.Type)
				headerName := string(header.Name)
				headerValue := header.Value
				headers += fmt.Sprintf("%s:%s:%s,", headerType, headerName, headerValue)
			}

			// 处理backendRefs
			for _, backendRef := range rule.BackendRefs {
				kind := "Service" // 默认为Service
				if backendRef.Kind != nil {
					kind = string(*backendRef.Kind)
				}
				service := string(backendRef.Name)
				port := ""
				if backendRef.Port != nil {
					port = fmt.Sprintf("%d", *backendRef.Port)
				}
				weight := "1" // 默认权重为1
				if backendRef.Weight != nil {
					weight = fmt.Sprintf("%d", *backendRef.Weight)
				}
				//backendInfo := fmt.Sprintf("%s:%s,port:%s,weight:%s", kind, service, port, weight)
				backendRefs += fmt.Sprintf("kind:%s,%s:%s,weight:%s;", kind, service, port, weight)
				//rules = append(rules, backendInfo)
			}

		}
		if len(paths) > 0 {
			paths = paths[0 : len(paths)-1]
		}
		if len(headers) > 0 {
			headers = headers[0 : len(headers)-1]
		}
		if len(backendRefs) > 0 {
			backendRefs = backendRefs[0 : len(backendRefs)-1]
		}
		rules = append(rules, RouteRules{
			MatchPath:   paths,
			Headers:     headers,
			BackendRefs: backendRefs,
		})
	}
	// 处理 conditions
	var conditions []HTTPRouteCondition
	// 注意：HTTPRoute的Conditions在ParentRefStatuses中，而不是直接在Status中
	if len(hr.Status.Parents) > 0 {
		for _, parent := range hr.Status.Parents {
			for _, cond := range parent.Conditions {
				condition := HTTPRouteCondition{
					Type:               string(cond.Type),
					Status:             string(cond.Status),
					Reason:             cond.Reason,
					Message:            cond.Message,
					LastTransitionTime: cond.LastTransitionTime.Format("2006-01-02 15:04:05"),
				}
				conditions = append(conditions, condition)
			}
		}
	}

	// 获取controllerName
	var controllerName []string
	for _, vv := range hr.Status.Parents {
		controllerName = append(controllerName, string(vv.ControllerName))
	}
	controllerNameStr := strings.Join(controllerName, ",")

	detail := &HTTPRouteDetail{
		HTTPRouteName:  hr.Name,
		NameSpace:      hr.Namespace,
		ParentRefs:     parentRefsStr,
		Hostnames:      hostnamesStr,
		CreateTime:     hr.CreationTimestamp.Format("2006-01-02 15:04:05"),
		ControllerName: controllerNameStr,
		Conditions:     conditions,
		Rules:          rules,
	}

	return detail, nil
}

// DeleteHTTPRoute 删除HTTPRoute
func DeleteHTTPRoute(kubeconfig, nameSpace, httprouteName string) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	err = clientset.GatewayV1().HTTPRoutes(nameSpace).Delete(context.TODO(), httprouteName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete httproute %s: %v\n", httprouteName, err)
		return err
	}

	return nil
}

// CreateHTTPRoute 创建HTTPRoute
func CreateHTTPRoute(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	httprouteName := gp.Get("httprouteName").String()
	gatewayName := gp.Get("gatewayName").String()
	nameSpace := gp.Get("nameSpace").String()
	hostNames := gp.Get("hostNames").String()

	var labelsMap = make(map[string]string)
	for _, vv := range gp.Get("labels").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	var rules []gatewayv1.HTTPRouteRule
	for _, vv := range gp.Get("rules").Array() {
		var matches []gatewayv1.HTTPRouteMatch
		var headers []gatewayv1.HTTPHeaderMatch

		// 处理header匹配规则
		kv := strings.Split(vv.Get("header").Str, ":")
		if len(kv) >= 2 {
			headerType := gatewayv1.HeaderMatchExact
			headers = append(headers, gatewayv1.HTTPHeaderMatch{
				Type:  &headerType,
				Name:  gatewayv1.HTTPHeaderName(kv[0]),
				Value: kv[1],
			})
		}

		// 处理路径匹配规则
		pathTypeStr := vv.Get("pathType").Str
		var pathType *gatewayv1.PathMatchType
		if pathTypeStr != "" {
			pt := gatewayv1.PathMatchType(pathTypeStr)
			pathType = &pt
		}

		pathValueStr := vv.Get("pathValue").Str
		var pathValue *string
		if pathValueStr != "" {
			pathValue = &pathValueStr
		}

		pathMatch := gatewayv1.HTTPPathMatch{
			Type:  pathType,
			Value: pathValue,
		}

		matches = append(matches, gatewayv1.HTTPRouteMatch{
			Headers: headers,
			Path:    &pathMatch,
		})

		// 处理后端服务
		var backendRefs []gatewayv1.HTTPBackendRef
		kv2 := strings.Split(vv.Get("svcPort").Str, ":")
		if len(kv2) >= 2 {
			svcName := kv2[0]
			port, err := strconv.ParseInt(kv2[1], 10, 32)
			if err != nil {
				log.Printf("Failed to parse port: %v\n", err)
				return err
			}

			portNum := gatewayv1.PortNumber(port)
			weightVal := int32(vv.Get("weight").Int())

			backendRef := gatewayv1.BackendRef{
				BackendObjectReference: gatewayv1.BackendObjectReference{
					Name: gatewayv1.ObjectName(svcName),
					Port: &portNum,
				},
				Weight: &weightVal,
			}

			backendRefs = append(backendRefs, gatewayv1.HTTPBackendRef{
				BackendRef: backendRef,
			})
		}

		rules = append(rules, gatewayv1.HTTPRouteRule{
			Matches:     matches,
			BackendRefs: backendRefs,
		})
	}

	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 构建HTTPRoute对象
	var hostnames []gatewayv1.Hostname
	for _, h := range strings.Split(hostNames, ",") {
		if h != "" {
			hostname := gatewayv1.Hostname(h)
			hostnames = append(hostnames, hostname)
		}
	}

	httproute := &gatewayv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      httprouteName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: gatewayv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{
					{
						Name: gatewayv1.ObjectName(gatewayName),
					},
				},
			},
			Hostnames: hostnames,
			Rules:     rules,
		},
	}

	_, err = clientset.GatewayV1().HTTPRoutes(nameSpace).Create(context.TODO(), httproute, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create httproute %s: %v\n", httprouteName, err)
		return err
	}

	return nil
}

// UpdateHTTPRouteByYaml 通过YAML更新HTTPRoute
func UpdateHTTPRouteByYaml_1(kubeconfig string, yamlData []byte) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 解析YAML
	httproute := &gatewayv1.HTTPRoute{}
	if err := yaml.Unmarshal(yamlData, httproute); err != nil {
		log.Printf("Failed to unmarshal yaml: %v\n", err)
		return err
	}

	// 更新HTTPRoute
	_, err = clientset.GatewayV1().HTTPRoutes(httproute.Namespace).Update(context.TODO(), httproute, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update httproute %s: %v\n", httproute.Name, err)
		return err
	}

	return nil
}

func UpdateHTTPRouteByYaml(kubeconfig string, yamlData []byte) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	httproute := &gatewayv1.HTTPRoute{}
	err = json.Unmarshal(data, httproute)
	if err != nil {
		return err
	}

	hr, err := clientset.GatewayV1().HTTPRoutes(httproute.Namespace).Get(context.TODO(), httproute.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get httproute %s: %v\n", httproute.Name, err)
		return err
	}
	httproute.ObjectMeta.ResourceVersion = hr.ObjectMeta.ResourceVersion
	// 更新HTTPRoute
	_, err = clientset.GatewayV1().HTTPRoutes(httproute.Namespace).Update(context.TODO(), httproute, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update httproute %s: %v\n", httproute.Name, err)
		return err
	}
	return nil
}

// GetHTTPRouteYaml 获取HTTPRoute的YAML配置
func GetHTTPRouteYaml(kubeconfig, nameSpace, httprouteName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return "", err
	}

	hr, err := clientset.GatewayV1().HTTPRoutes(nameSpace).Get(context.TODO(), httprouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get httproute %s: %v\n", httprouteName, err)
		return "", err
	}

	hrUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(hr)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(hrUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
