// pod_model.go
package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"xkube/common"

	"encoding/json"

	"github.com/tidwall/gjson"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// Service 用于返回给前端的Service信息结构体 Service information structure for frontend return
type Service struct {
	ServiceName string `json:"serviceName"` // Service名称 Service name
	NameSpace   string `json:"nameSpace"`   // 命名空间 Namespace
	Labels      string `json:"labels"`      // 标签 Labels
	SvcType     string `json:"svcType"`     // Service类型 Service type
	SvcIp       string `json:"svcIp"`       // Service IP地址 Service IP address
	SvcPort     string `json:"svcPort"`     // Service端口 Service port
	LanEndpoint string `json:"lanEndpoint"` // 内部端点 Internal endpoint
	WanEndpoint string `json:"wanEndpoint"` // 外部端点 External endpoint
	CreateTime  string `json:"createTime"`  // 创建时间 Create time
}

// ServiceMain Service主要信息结构体 Main service information structure
type ServiceMain struct {
	ServiceName string            `json:"serviceName"` // Service名称 Service name
	NameSpace   string            `json:"nameSpace"`   // 命名空间 Namespace
	Labels      map[string]string `json:"labels"`      // 标签映射 Label mapping
	SvcType     string            `json:"svcType"`     // Service类型 Service type
	Ports       []ServicePort     `json:"ports"`       // 端口列表 Port list
	IsHeadless  string            `json:"isHeadless"`  // 是否为Headless Service Is headless service
}

// KV 键值对结构体 Key-value pair structure
type KV struct {
	Key   string `json:"key"`   // 键 Key
	Value string `json:"value"` // 值 Value
}

// ServicePort Service端口信息结构体 Service port information structure
type ServicePort struct {
	PortName   string `json:"portName"`   // 端口名称 Port name
	SvcPort    int32  `json:"svcPort"`    // Service端口 Service port
	TargetPort int32  `json:"targetPort"` // 目标端口 Target port
}

// SvcList 获取Service列表 Get Service list
// kubeconfig: 集群配置信息 Cluster configuration information
// namespace: 命名空间 Namespace
// serviceName: Service名称 Service name
// labelsKey: 标签键 Label key
// labelsValue: 标签值 Label value
// 返回Service列表和错误信息 Return service list and error information
func SvcList(kubeconfig, namespace, serviceName string, labelsKey, labelsValue string) ([]Service, error) {
	clientset := common.ClientSet(kubeconfig)
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}
	//var selector labels.Selector

	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue)}
	}

	svcList, err := clientset.CoreV1().Services(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list service error:%v\n", err)
	}

	//fmt.Println("svc count:", len(svcList.Items))
	var bbb = make([]Service, 0)
	for _, svc := range svcList.Items {
		//搜索
		if serviceName != "" {
			if !strings.Contains(svc.Name, serviceName) {
				continue
			}
		}
		//fmt.Printf("name: %s\n", svc.Name)
		var labelsStr string
		for kk, vv := range svc.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		var wanEndPoint string
		svcType := string(svc.Spec.Type)
		if svcType == "LoadBalancer" {
			if svc.Status.LoadBalancer.Ingress != nil {
				wanEndPoint = fmt.Sprintf("%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, svc.Spec.Ports[0].Port)
			}
		}
		lanEndPoint := ""
		var svcPort string
		if len(svc.Spec.Ports) > 0 {
			vsapp := svc.Spec.Selector["app"]
			if vsapp == "" {
				vsapp = svc.Spec.Selector["k8s-app"]
			}
			for _, vv := range svc.Spec.Ports {
				//lanEndPoint += fmt.Sprintf("%s:%s:%s\n", svc.Spec.Selector["app"], svc.Spec.Ports[0].TargetPort.StrVal, svc.Spec.Ports[0].Protocol)
				//svcPort += fmt.Sprintf("%s:", svc.Spec.Ports[0].Port)
				var vport string
				if vv.TargetPort.Type == 0 {
					vport = fmt.Sprintf("%d", vv.TargetPort.IntVal)
				} else {
					vport = vv.TargetPort.StrVal
				}
				lanEndPoint += fmt.Sprintf("%s:%s:%s,", vsapp, vport, vv.Protocol)
				svcPort += fmt.Sprintf("%d,", vv.Port)
			}
			if len(svcPort) > 0 {
				svcPort = svcPort[0 : len(svcPort)-1]
			}
			lanEndPoint = lanEndPoint[0 : len(lanEndPoint)-1]
		}
		Items := &Service{
			ServiceName: svc.Name,
			NameSpace:   svc.Namespace,
			SvcType:     svcType,
			SvcIp:       svc.Spec.ClusterIP,
			Labels:      labelsStr,
			SvcPort:     svcPort,
			LanEndpoint: lanEndPoint,
			WanEndpoint: wanEndPoint,
			CreateTime:  svc.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
	}
	return bbb, err
}

// SvcCreate 创建Service Create Service
// kubeconfig: 集群配置信息 Cluster configuration information
// bodys: 请求体数据 Request body data
// 返回错误信息 Return error information
func SvcCreate(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)
	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	serviceName := gp.Get("serviceName").String()
	nameSpace := gp.Get("nameSpace").String()
	svcType := gp.Get("svcType").String()
	deployName := gp.Get("deployName").String()
	isHeadless := gp.Get("isHeadless").Str

	var labelsMap = make(map[string]string)
	labelsMap["app"] = serviceName
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").Str] = vv.Get("value").Str
	}

	selectApp := serviceName
	if deployName != "" {
		selectApp = deployName
	}

	var serviceType corev1.ServiceType
	switch svcType {
	case "NodePort":
		serviceType = corev1.ServiceTypeNodePort
	case "LoadBalancer":
		serviceType = corev1.ServiceTypeLoadBalancer
	default:
		serviceType = corev1.ServiceTypeClusterIP
	}

	if isHeadless == "on" {
		serviceType = corev1.ServiceTypeClusterIP
	}

	svcInstance := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": selectApp,
			},
			Type: serviceType,
		},
	}
	if isHeadless == "on" {
		svcInstance.Spec.ClusterIP = corev1.ClusterIPNone
	}

	ports := gp.Get("ports").Array()
	var svcPorts = make([]corev1.ServicePort, 0)
	for _, vv := range ports {
		var svcProtocol corev1.Protocol
		if vv.Get("protocol").Str == "UDP" {
			svcProtocol = corev1.ProtocolUDP
		} else {
			svcProtocol = corev1.ProtocolTCP
		}
		svcPort := &corev1.ServicePort{
			Name:       vv.Get("portName").Str,
			Port:       int32(vv.Get("svcPort").Int()),
			Protocol:   svcProtocol,
			TargetPort: intstr.FromInt32(int32(vv.Get("targetPort").Int())),
		}
		svcPorts = append(svcPorts, *svcPort)
	}
	svcInstance.Spec.Ports = svcPorts

	clientset := common.ClientSet(kubeconfig)
	_, err := clientset.CoreV1().Services(nameSpace).Create(context.TODO(), svcInstance, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// SvcYamlCreate 通过YAML创建Service Create Service by YAML
// kubeconfig: 集群配置信息 Cluster configuration information
// yamlData: YAML数据 YAML data
// 返回错误信息 Return error information
func SvcYamlCreate(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	service := &corev1.Service{}
	err = json.Unmarshal(data, service)
	if err != nil {
		return err
	}

	namespace := service.ObjectMeta.Namespace
	serviceName := service.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, serviceName)
	return err
}

// SvcYamlModify 通过YAML修改Service Modify Service by YAML
// kubeconfig: 集群配置信息 Cluster configuration information
// yamlData: YAML数据 YAML data
// 返回错误信息 Return error information
func SvcYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	service := &corev1.Service{}
	err = json.Unmarshal(data, service)
	if err != nil {
		return err
	}

	namespace := service.ObjectMeta.Namespace
	serviceName := service.ObjectMeta.Name
	clientset := common.ClientSet(kubeconfig)
	_, err = clientset.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Println(namespace, serviceName)
	return err
}

// SvcClone 克隆Service Clone Service
// kubeconfig: 源集群配置信息 Source cluster configuration information
// namespace: 源命名空间 Source namespace
// objname: 源对象名称 Source object name
// target_clusterid: 目标集群ID Target cluster ID
// target_namespace: 目标命名空间 Target namespace
// target_objname: 目标对象名称 Target object name
// 返回错误信息 Return error information
func SvcClone(kubeconfig, namespace, objname, target_clusterid, target_namespace, target_objname string) error {
	//old cluster
	svc, err := common.ClientSet(kubeconfig).CoreV1().Services(namespace).Get(context.TODO(), objname, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if target_namespace == "" && target_objname == "" && target_clusterid == "" {
		return fmt.Errorf("target can't be all empty")
	}

	if target_namespace == "" {
		target_namespace = namespace
	}

	if target_objname == "" {
		target_objname = objname
	}

	if target_clusterid == "" {
		target_clusterid = kubeconfig
	}
	//在新集群创建namespace
	err3 := CreateNsByExist(target_clusterid, target_namespace)
	if err3 != nil {
		return err3
	}

	svc.Name = target_objname
	svc.ResourceVersion = ""
	svc.Namespace = target_namespace
	svc.ObjectMeta.Labels["app"] = target_objname
	svc.Spec.Selector["app"] = target_objname
	svc.Spec.ClusterIP = ""
	svc.Spec.ClusterIPs = []string{}
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		svc.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{}
	}
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer || svc.Spec.Type == corev1.ServiceTypeNodePort {
		for i, _ := range svc.Spec.Ports {
			svc.Spec.Ports[i].NodePort = 0
		}
	}

	//new cluster
	NewSvcClient := common.ClientSet(target_clusterid).CoreV1().Services(target_namespace)
	_, err2 := NewSvcClient.Get(context.TODO(), target_objname, metav1.GetOptions{})
	if errors.IsNotFound(err2) { //没有就创建
		_, err := NewSvcClient.Create(context.TODO(), svc, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if err2 == nil {
		_, err = NewSvcClient.Update(context.TODO(), svc, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	return err2
}

// SvcDetail 获取Service详情 Get Service details
// kubeconfig: 集群配置信息 Cluster configuration information
// namespace: 命名空间 Namespace
// serviceName: Service名称 Service name
// 返回Service详情和错误信息 Return service details and error information
func SvcDetail(kubeconfig, namespace, serviceName string) (*Service, error) {
	svc, err := common.ClientSet(kubeconfig).CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	var labelsStr string
	for kk, vv := range svc.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}
	var wanEndPoint string
	svcType := string(svc.Spec.Type)
	if svcType == "LoadBalancer" {
		wanEndPoint = fmt.Sprintf("%s:%d", svc.Status.LoadBalancer.Ingress[0].IP, svc.Spec.Ports[0].Port)
	}

	lanEndPoint := ""
	var svcPort string
	if len(svc.Spec.Ports) > 0 {
		vsapp := svc.Spec.Selector["app"]
		if vsapp == "" {
			vsapp = svc.Spec.Selector["k8s-app"]
		}
		for _, vv := range svc.Spec.Ports {
			//lanEndPoint += fmt.Sprintf("%s:%s:%s\n", svc.Spec.Selector["app"], svc.Spec.Ports[0].TargetPort.StrVal, svc.Spec.Ports[0].Protocol)
			//svcPort += fmt.Sprintf("%s:", svc.Spec.Ports[0].Port)
			var vport string
			if vv.TargetPort.Type == 0 {
				vport = fmt.Sprintf("%d", vv.TargetPort.IntVal)
			} else {
				vport = vv.TargetPort.StrVal
			}
			lanEndPoint += fmt.Sprintf("%s:%s:%s,", vsapp, vport, vv.Protocol)
			svcPort += fmt.Sprintf("%d,", vv.Port)
		}
		if len(svcPort) > 0 {
			svcPort = svcPort[0 : len(svcPort)-1]
		}
		lanEndPoint = lanEndPoint[0 : len(lanEndPoint)-1]
	}

	return &Service{
		ServiceName: svc.Name,
		NameSpace:   svc.Namespace,
		SvcType:     svcType,
		SvcIp:       svc.Spec.ClusterIP,
		Labels:      labelsStr,
		SvcPort:     svcPort,
		LanEndpoint: lanEndPoint,
		WanEndpoint: wanEndPoint,
		CreateTime:  svc.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetSvcYaml 获取Service的YAML格式 Get Service in YAML format
// kubeconfig: 集群配置信息 Cluster configuration information
// namespace: 命名空间 Namespace
// serviceName: Service名称 Service name
// 返回YAML字符串和错误信息 Return YAML string and error information
func GetSvcYaml(kubeconfig, namespace, serviceName string) (string, error) {
	servicesClient := common.ClientSet(kubeconfig).CoreV1().Services(namespace)
	service, err := servicesClient.Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	serviceUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(service)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(serviceUnstructured)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

// SvcDelete 删除Service Delete Service
// kubeconfig: 集群配置信息 Cluster configuration information
// namespace: 命名空间 Namespace
// serviceName: Service名称 Service name
// 返回错误信息 Return error information
func SvcDelete(kubeconfig, namespace, serviceName string) error {
	//deletePolicy := metav1.DeletePropagationForeground
	if err := common.ClientSet(kubeconfig).CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
