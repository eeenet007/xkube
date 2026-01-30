// grpcroute_model.go
package models

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"xkube/common"

	"github.com/tidwall/gjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayClient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

// GRPCRoute GRPCRoute列表信息结构体
type GRPCRoute struct {
	GRPCRouteName string `json:"grpcrouteName"`
	NameSpace     string `json:"nameSpace"`
	ParentRefs    string `json:"parentRefs"`
	Hostnames     string `json:"hostnames"`
	Rules         string `json:"rules"`
	CreateTime    string `json:"createTime"`
}

// GRPCRouteDetail GRPCRoute详细信息结构体
type GRPCRouteDetail struct {
	GRPCRoute
	Conditions     []GRPCRouteCondition `json:"conditions"`
	ControllerName string               `json:"controllerName"`
	ParentRef      string               `json:"parentRef"`
}

// GRPCRouteCondition GRPCRoute条件信息结构体
type GRPCRouteCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

// GRPCRouteList 获取GRPCRoute列表
func GRPCRouteList(kubeconfig string) ([]GRPCRoute, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	grpcrouteList, err := clientset.GatewayV1().GRPCRoutes("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list grpcroutes: %v\n", err)
		return nil, err
	}
	if len(grpcrouteList.Items) == 0 {
		log.Printf("[WARN] grpcrouteList NoFound RequestedResource")
		return nil, fmt.Errorf("NoFound RequestedResource")
	}

	var grpcroutes []GRPCRoute
	for _, gr := range grpcrouteList.Items {
		// 获取parentRefs信息
		var parentRefs []string
		for _, parentRef := range gr.Spec.ParentRefs {
			parentRefs = append(parentRefs, string(parentRef.Name))
		}
		parentRefsStr := strings.Join(parentRefs, ",")

		// 获取hostnames信息
		var hostnames []string
		for _, hostname := range gr.Spec.Hostnames {
			hostnames = append(hostnames, string(hostname))
		}
		hostnamesStr := strings.Join(hostnames, ",")

		// 获取rules信息
		var rulesStr string
		for _, rule := range gr.Spec.Rules {
			var matchs string
			backendRefs := ""
			// 处理matches
			for _, match := range rule.Matches {
				// method类型
				methods := ""
				headers := ""
				if match.Method != nil {
					if match.Method.Method != nil {
						methods += "Method:" + *match.Method.Method + ","
					}
					if match.Method.Service != nil {
						methods += "Service:" + *match.Method.Service
					}
				}
				// headers类型
				for _, header := range match.Headers {
					//headerType := fmt.Sprintf("%v", header.Type)
					headerName := string(header.Name)
					headerValue := header.Value
					//rules = append(rules, headerType+":"+headerName+":"+headerValue)
					headers += headerName + ":" + headerValue + ","
				}
				if len(headers) > 0 {
					headers = headers[0 : len(headers)-1]
				}
				if headers != "" {
					matchs += fmt.Sprintf("%s,header[%s];", methods, headers)
				} else {
					matchs += fmt.Sprintf("%s;", methods)
				}
			}

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
			if len(matchs) > 0 {
				matchs = matchs[0 : len(matchs)-1]
			}
			if len(backendRefs) > 0 {
				backendRefs = backendRefs[0 : len(backendRefs)-1]
			}
			rulesStr += fmt.Sprintf("[%s]-->[%s]|", matchs, backendRefs)
		}
		if len(rulesStr) > 0 {
			rulesStr = rulesStr[0 : len(rulesStr)-1]
		}

		grpcroute := GRPCRoute{
			GRPCRouteName: gr.Name,
			NameSpace:     gr.Namespace,
			ParentRefs:    parentRefsStr,
			Hostnames:     hostnamesStr,
			Rules:         rulesStr,
			CreateTime:    gr.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		grpcroutes = append(grpcroutes, grpcroute)
	}

	return grpcroutes, nil
}

// GetGRPCRouteDetail 获取GRPCRoute详细信息
func GetGRPCRouteDetail(kubeconfig, nameSpace, grpcrouteName string) (*GRPCRouteDetail, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	gr, err := clientset.GatewayV1().GRPCRoutes(nameSpace).Get(context.TODO(), grpcrouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get grpcroute %s: %v\n", grpcrouteName, err)
		return nil, err
	}

	// 获取parentRefs信息
	var parentRefs []string
	for _, parentRef := range gr.Spec.ParentRefs {
		parentRefs = append(parentRefs, string(parentRef.Name))
	}
	parentRefsStr := strings.Join(parentRefs, ",")

	// 获取hostnames信息
	var hostnames []string
	for _, hostname := range gr.Spec.Hostnames {
		hostnames = append(hostnames, string(hostname))
	}
	hostnamesStr := strings.Join(hostnames, ",")

	// 获取rules信息
	//var rules []string
	var rulesStr string
	for _, rule := range gr.Spec.Rules {
		var matchs string
		backendRefs := ""
		// 处理matches
		for _, match := range rule.Matches {
			// method类型
			methods := ""
			headers := ""
			if match.Method != nil {
				if match.Method.Method != nil {
					methods += "Method:" + *match.Method.Method + ","
				}
				if match.Method.Service != nil {
					methods += "Service:" + *match.Method.Service
				}
			}
			// headers类型
			for _, header := range match.Headers {
				//headerType := fmt.Sprintf("%v", header.Type)
				headerName := string(header.Name)
				headerValue := header.Value
				//rules = append(rules, headerType+":"+headerName+":"+headerValue)
				headers += headerName + ":" + headerValue + ","
			}
			if len(headers) > 0 {
				headers = headers[0 : len(headers)-1]
			}
			if headers != "" {
				matchs += fmt.Sprintf("%s,header[%s];", methods, headers)
			} else {
				matchs += fmt.Sprintf("%s;", methods)
			}
		}

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
		if len(matchs) > 0 {
			matchs = matchs[0 : len(matchs)-1]
		}
		if len(backendRefs) > 0 {
			backendRefs = backendRefs[0 : len(backendRefs)-1]
		}
		rulesStr += fmt.Sprintf("%s-->%s|", matchs, backendRefs)
	}

	grpcroute := GRPCRoute{
		GRPCRouteName: gr.Name,
		NameSpace:     gr.Namespace,
		ParentRefs:    parentRefsStr,
		Hostnames:     hostnamesStr,
		Rules:         rulesStr,
		CreateTime:    gr.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}

	// 处理 conditions
	var conditions []GRPCRouteCondition
	// 注意：GRPCRoute的Conditions在Parents中，而不是直接在Status中
	if len(gr.Status.Parents) > 0 {
		for _, parent := range gr.Status.Parents {
			for _, cond := range parent.Conditions {
				condition := GRPCRouteCondition{
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
	for _, vv := range gr.Status.Parents {
		controllerName = append(controllerName, string(vv.ControllerName))
	}
	controllerNameStr := strings.Join(controllerName, ",")

	// 获取parentRef (第一个)
	parentRef := ""
	if len(gr.Spec.ParentRefs) > 0 {
		parentRef = string(gr.Spec.ParentRefs[0].Name)
	}

	detail := &GRPCRouteDetail{
		GRPCRoute:      grpcroute,
		Conditions:     conditions,
		ControllerName: controllerNameStr,
		ParentRef:      parentRef,
	}

	return detail, nil
}

// DeleteGRPCRoute 删除GRPCRoute
func DeleteGRPCRoute(kubeconfig, nameSpace, grpcrouteName string) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	err = clientset.GatewayV1().GRPCRoutes(nameSpace).Delete(context.TODO(), grpcrouteName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete grpcroute %s: %v\n", grpcrouteName, err)
		return err
	}

	return nil
}

// CreateGRPCRoute 创建GRPCRoute
func CreateGRPCRoute(kubeconfig string, bodys []byte) error {

	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	grpcrouteName := gp.Get("grpcrouteName").String()
	gatewayName := gp.Get("gatewayName").String()
	nameSpace := gp.Get("nameSpace").String()
	hostNames := gp.Get("hostNames").String()

	var labelsMap = make(map[string]string)
	for _, vv := range gp.Get("labels").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	var rules []gatewayv1.GRPCRouteRule
	for _, vv := range gp.Get("rules").Array() {
		var matches []gatewayv1.GRPCRouteMatch
		var headers []gatewayv1.GRPCHeaderMatch

		// 处理header匹配规则
		kv := strings.Split(vv.Get("header").Str, ":")
		if len(kv) >= 2 {
			headerType := gatewayv1.GRPCHeaderMatchExact
			headers = append(headers, gatewayv1.GRPCHeaderMatch{
				Type:  &headerType,
				Name:  gatewayv1.GRPCHeaderName(kv[0]),
				Value: kv[1],
			})
		}

		// 处理方法规则
		serviceStr := vv.Get("service").Str
		methodStr := vv.Get("method").Str

		var serviceValue *string
		if serviceStr != "" {
			serviceValue = &serviceStr
		}

		var methodValue *string
		if methodStr != "" {
			methodValue = &methodStr
		}
		methodType := gatewayv1.GRPCMethodMatchExact
		method := gatewayv1.GRPCMethodMatch{
			Type:    &methodType,
			Service: serviceValue,
			Method:  methodValue,
		}
		matches = append(matches, gatewayv1.GRPCRouteMatch{
			Method:  &method,
			Headers: headers,
		})

		// 处理后端服务
		var backendRefs []gatewayv1.GRPCBackendRef
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

			backendRefs = append(backendRefs, gatewayv1.GRPCBackendRef{
				BackendRef: backendRef,
			})
		}

		rules = append(rules, gatewayv1.GRPCRouteRule{
			Matches:     matches,
			BackendRefs: backendRefs,
		})
	}
	//构建域名
	var hostnames []gatewayv1.Hostname
	for _, h := range strings.Split(hostNames, ",") {
		if h != "" {
			hostname := gatewayv1.Hostname(h)
			hostnames = append(hostnames, hostname)
		}
	}

	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 构建GRPCRoute对象
	grpcroute := &gatewayv1.GRPCRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      grpcrouteName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: gatewayv1.GRPCRouteSpec{
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

	_, err = clientset.GatewayV1().GRPCRoutes(nameSpace).Create(context.TODO(), grpcroute, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create grpcroute %s: %v\n", grpcrouteName, err)
		return err
	}

	return nil
}

// UpdateGRPCRouteByYaml 通过YAML更新GRPCRoute
func UpdateGRPCRouteByYaml(kubeconfig string, yamlData []byte) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 解析YAML
	grpcroute := &gatewayv1.GRPCRoute{}
	if err := yaml.Unmarshal(yamlData, grpcroute); err != nil {
		log.Printf("Failed to unmarshal yaml: %v\n", err)
		return err
	}

	gr, err := clientset.GatewayV1().GRPCRoutes(grpcroute.Namespace).Get(context.TODO(), grpcroute.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get grpcroute %s: %v\n", grpcroute.Name, err)
		return err
	}
	grpcroute.ObjectMeta.ResourceVersion = gr.ObjectMeta.ResourceVersion
	// 更新GRPCRoute
	_, err = clientset.GatewayV1().GRPCRoutes(grpcroute.Namespace).Update(context.TODO(), grpcroute, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update grpcroute %s: %v\n", grpcroute.Name, err)
		return err
	}

	return nil
}

// GetGRPCRouteYaml 获取GRPCRoute的YAML配置
func GetGRPCRouteYaml(kubeconfig, nameSpace, grpcrouteName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return "", err
	}

	gr, err := clientset.GatewayV1().GRPCRoutes(nameSpace).Get(context.TODO(), grpcrouteName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get grpcroute %s: %v\n", grpcrouteName, err)
		return "", err
	}

	grUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(gr)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(grUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
