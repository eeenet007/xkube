// gateway_model.go
package models

import (
	"context"
	//"encoding/json"
	"fmt"
	"log"

	//"strings"
	//"time"
	"xkube/common"

	"github.com/tidwall/gjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayClient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
	"sigs.k8s.io/yaml"
)

// Gateway 网关列表信息结构体
type Gateway struct {
	GatewayName      string `json:"gatewayName"`
	NameSpace        string `json:"nameSpace"`
	GatewayClassName string `json:"gatewayClassName"`
	Address          string `json:"address"`
	Programmed       string `json:"programmed"`
	AttachedRoutes   int32  `json:"attachedRoutes"`
	CreateTime       string `json:"createTime"`
	//Listeners        []Listener `json:"listeners"`
	Listenerstr string `json:"listenerstr"`
}

// Listener 监听器信息结构体
type Listener struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	Port          int32  `json:"port"`
	AllowedRoutes string `json:"allowedRoutes"`
	HostName      string `json:"hostName"`
	Tls           string `json:"tls"`
}

// GatewayDetail 网关详细信息结构体
type GatewayDetail struct {
	Gateway
	Listeners  []Listener         `json:"listeners"`
	Conditions []GatewayCondition `json:"conditions"`
	//SupportedKinds []SupportedKind    `json:"supportedKinds"`
	StatusListener []StatusListener `json:"statusListener"`
}

// 监听状态信息结构体
type StatusListener struct {
	ListenersName  string             `json:"listenersName"`
	AttachedRoutes int32              `json:"attachedRoutes"`
	Conditions     []GatewayCondition `json:"conditions"`
	SupportedKinds string             `json:"supportedKinds"`
}

// GatewayCondition 网关条件信息结构体
type GatewayCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

// SupportedKind 支持的路由类型信息结构体
type SupportedKind struct {
	Group *string `json:"group,omitempty"`
	Kind  string  `json:"kind"`
}

// GatewayList 获取Gateway列表
func GatewayList(kubeconfig string) ([]Gateway, error) {
	var gateways []Gateway
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return gateways, err
	}

	gatewayList, err := clientset.GatewayV1().Gateways("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list gateways: %v\n", err)
		return gateways, err
	}

	if len(gatewayList.Items) == 0 {
		log.Printf("[WARN] gatewayList NoFound RequestedResource")
		return nil, fmt.Errorf("NoFound RequestedResource")
	}

	for _, gw := range gatewayList.Items {
		// 获取 Programmed 状态
		programmed := "Unknown"
		for _, condition := range gw.Status.Conditions {
			if condition.Type == string(gatewayv1.GatewayConditionProgrammed) {
				programmed = string(condition.Status)
				break
			}
		}

		// 获取地址
		address := ""
		if len(gw.Status.Addresses) > 0 {
			address = gw.Status.Addresses[0].Value
		}

		var listeners string
		for _, listener := range gw.Spec.Listeners {
			listeners += fmt.Sprintf("%s:%s:%d,", listener.Name, string(listener.Protocol), listener.Port)
		}
		if len(listeners) > 0 {
			listeners = listeners[0 : len(listeners)-1]
		}

		// 计算附加的路由数量
		attachedRoutes := int32(0)
		if len(gw.Status.Listeners) > 0 {
			for _, listenerStatus := range gw.Status.Listeners {
				attachedRoutes += listenerStatus.AttachedRoutes
			}
		}

		gateway := Gateway{
			GatewayName:      gw.Name,
			NameSpace:        gw.Namespace,
			GatewayClassName: string(gw.Spec.GatewayClassName),
			Address:          address,
			Programmed:       programmed,
			AttachedRoutes:   attachedRoutes,
			CreateTime:       gw.CreationTimestamp.Format("2006-01-02 15:04:05"),
			Listenerstr:      listeners,
		}
		gateways = append(gateways, gateway)
	}

	return gateways, nil
}

// GetGatewayDetail 获取Gateway详细信息
func GetGatewayDetail(kubeconfig, nameSpace, gatewayName string) (*GatewayDetail, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return nil, err
	}

	gw, err := clientset.GatewayV1().Gateways(nameSpace).Get(context.TODO(), gatewayName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get gateway %s: %v\n", gatewayName, err)
		return nil, err
	}

	// 获取 Programmed 状态
	programmed := "Unknown"
	for _, condition := range gw.Status.Conditions {
		if condition.Type == string(gatewayv1.GatewayConditionProgrammed) {
			programmed = string(condition.Status)
			break
		}
	}

	// 获取地址
	address := ""
	if len(gw.Status.Addresses) > 0 {
		address = gw.Status.Addresses[0].Value
	}

	// 获取监听器信息
	var listeners []Listener
	for _, listener := range gw.Spec.Listeners {

		var hostName string
		if listener.Hostname != nil {
			hostName = string(*listener.Hostname)
		}
		var allowdRoutes string
		if listener.AllowedRoutes != nil {
			var kinds string
			for _, kk := range listener.AllowedRoutes.Kinds {
				kinds += fmt.Sprintf("%s/%s,", string(*kk.Group), string(kk.Kind))
			}

			var fromstr string
			if listener.AllowedRoutes.Namespaces.From != nil {
				fromstr = string(*listener.AllowedRoutes.Namespaces.From)
			}

			var selectorstr string
			if listener.AllowedRoutes.Namespaces.Selector != nil {
				for k1, v1 := range listener.AllowedRoutes.Namespaces.Selector.MatchLabels {
					selectorstr += fmt.Sprintf("%s:%s;", k1, v1)
				}
				if len(selectorstr) > 0 {
					selectorstr = selectorstr[0 : len(selectorstr)-1]
				}
			}
			allowdRoutes = fmt.Sprintf("kinds:%s,namespaceFrom:%s,selectorMatchLabels:%s", kinds, fromstr, selectorstr)
		}
		var tlsstr string
		if listener.TLS != nil {
			for _, v2 := range listener.TLS.CertificateRefs {
				tlsstr += fmt.Sprintf("kind:%s;group:%s;name:%s;namespace:%s,", string(*v2.Kind), string(*v2.Group), v2.Name, string(*v2.Namespace))
			}
			tlsstr = tlsstr[0 : len(tlsstr)-1]
		}

		listenerInfo := Listener{
			Name:          string(listener.Name),
			Protocol:      string(listener.Protocol),
			Port:          int32(listener.Port),
			AllowedRoutes: allowdRoutes,
			HostName:      hostName,
			Tls:           tlsstr,
		}
		listeners = append(listeners, listenerInfo)
	}

	// 计算附加的路由数量
	attachedRoutes := int32(0)
	if len(gw.Status.Listeners) > 0 {
		for _, listenerStatus := range gw.Status.Listeners {
			attachedRoutes += listenerStatus.AttachedRoutes
		}
	}

	gateway := Gateway{
		GatewayName:      gw.Name,
		NameSpace:        gw.Namespace,
		GatewayClassName: string(gw.Spec.GatewayClassName),
		Address:          address,
		Programmed:       programmed,
		AttachedRoutes:   attachedRoutes,
		CreateTime:       gw.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}

	// 处理 conditions
	var conditions []GatewayCondition
	for _, cond := range gw.Status.Conditions {
		condition := GatewayCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			Reason:             cond.Reason,
			Message:            cond.Message,
			LastTransitionTime: cond.LastTransitionTime.Format("2006-01-02 15:04:05"),
		}
		conditions = append(conditions, condition)
	}

	// 处理 StatusListener
	var statusListener []StatusListener
	if len(gw.Status.Listeners) > 0 {
		for _, listenerStatus := range gw.Status.Listeners {

			var supportedKinds string
			for _, supportedKind := range listenerStatus.SupportedKinds {
				var groupStr string
				if supportedKind.Group != nil {
					groupStr = string(*supportedKind.Group)
				}
				supportedKinds += fmt.Sprintf("group:%s/kind:%s,", groupStr, string(supportedKind.Kind))
			}
			if len(supportedKinds) > 0 {
				supportedKinds = supportedKinds[0 : len(supportedKinds)-1]
			}

			var statusConditions []GatewayCondition
			for _, cond := range listenerStatus.Conditions {
				condition := GatewayCondition{
					Type:               string(cond.Type),
					Status:             string(cond.Status),
					Reason:             cond.Reason,
					Message:            cond.Message,
					LastTransitionTime: cond.LastTransitionTime.Format("2006-01-02 15:04:05"),
				}
				statusConditions = append(statusConditions, condition)
			}

			statusListener = append(statusListener, StatusListener{
				ListenersName:  string(listenerStatus.Name),
				AttachedRoutes: listenerStatus.AttachedRoutes,
				Conditions:     statusConditions,
				SupportedKinds: supportedKinds,
			})
		}
	}

	detail := &GatewayDetail{
		Gateway:        gateway,
		Listeners:      listeners,
		Conditions:     conditions,
		StatusListener: statusListener,
	}

	return detail, nil
}

// DeleteGateway 删除Gateway
func DeleteGateway(kubeconfig, nameSpace, gatewayName string) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	err = clientset.GatewayV1().Gateways(nameSpace).Delete(context.TODO(), gatewayName, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Failed to delete gateway %s: %v\n", gatewayName, err)
		return err
	}

	return nil
}

// CreateGateway 创建Gateway
func CreateGateway(kubeconfig string, bodys []byte) error {

	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	gatewayName := gp.Get("gatewayName").String()
	gatewayClassName := gp.Get("gatewayClassName").String()
	nameSpace := gp.Get("nameSpace").String()

	var labelsMap = make(map[string]string)
	for _, vv := range gp.Get("labels").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	var listener []gatewayv1.Listener
	for _, v1 := range gp.Get("listeners").Array() {
		listener = append(listener, gatewayv1.Listener{
			Name:     gatewayv1.SectionName(v1.Get("name").Str),
			Port:     gatewayv1.PortNumber(int32(v1.Get("port").Int())),
			Protocol: gatewayv1.ProtocolType(v1.Get("protocol").Str),
		})
	}

	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 构建监听器支持的路由类型
	// var supportedRouteKinds []gatewayv1.RouteGroupKind
	// for _, kind := range kinds {
	// 	supportedRouteKinds = append(supportedRouteKinds, gatewayv1.RouteGroupKind{
	// 		Group: (*gatewayv1.Group)(&gatewayv1.GroupVersion.Group),
	// 		Kind:  gatewayv1.Kind(kind),
	// 	})
	// }

	// 构建Gateway对象
	gateway := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gatewayName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: gatewayv1.GatewaySpec{
			GatewayClassName: gatewayv1.ObjectName(gatewayClassName),
			Listeners:        listener,
		},
	}

	_, err = clientset.GatewayV1().Gateways(nameSpace).Create(context.TODO(), gateway, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Failed to create gateway %s: %v\n", gatewayName, err)
		return err
	}

	return nil
}

// UpdateGatewayByYaml 通过YAML更新Gateway
func UpdateGatewayByYaml(kubeconfig string, yamlData []byte) error {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return err
	}

	// 解析YAML
	gateway := &gatewayv1.Gateway{}
	if err := yaml.Unmarshal(yamlData, gateway); err != nil {
		log.Printf("Failed to unmarshal yaml: %v\n", err)
		return err
	}

	gw, err := clientset.GatewayV1().Gateways(gateway.Namespace).Get(context.TODO(), gateway.Name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get gateway %s: %v\n", gateway.Name, err)
		return err
	}

	gw.ObjectMeta.ResourceVersion = gateway.ObjectMeta.ResourceVersion

	// 更新Gateway
	_, err = clientset.GatewayV1().Gateways(gateway.Namespace).Update(context.TODO(), gateway, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Failed to update gateway %s: %v\n", gateway.Name, err)
		return err
	}

	return nil
}

// GetGatewayYaml 获取Gateway的YAML配置
func GetGatewayYaml(kubeconfig, nameSpace, gatewayName string) (string, error) {
	config := common.ClientConfig(kubeconfig)
	clientset, err := gatewayClient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create gateway client: %v\n", err)
		return "", err
	}

	gw, err := clientset.GatewayV1().Gateways(nameSpace).Get(context.TODO(), gatewayName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Failed to get gateway %s: %v\n", gatewayName, err)
		return "", err
	}

	gwUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(gw)
	if err != nil {
		return "", err
	}

	yamlBytes, err := yaml.Marshal(gwUnstructured)
	if err != nil {
		return "", err
	}

	return string(yamlBytes), nil
}
