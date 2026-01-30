// k8sclient.go
// Package common provides common utilities and functions for the xkube project
// 包 common 为 xkube 项目提供通用工具和函数
package common

import (
	//"io/ioutil"
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientSet creates and returns a kubernetes clientset based on the cluster ID
// ClientSet 根据集群 ID 创建并返回一个 kubernetes clientset
func ClientSet(clusterid string) *kubernetes.Clientset {
	// var kubeconfigPath string

	// switch kubeconfig {
	// case "zx-cluster":
	// 	kubeconfigPath = "./kubeconfig.txt"
	// default:
	// 	kubeconfigPath = "./kubeconfig.txt"
	// }
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// kubeconfigBytes, err := ioutil.ReadFile(kubeconfigPath)
	// if err != nil {
	// 	panic(err.Error())
	// }

	//kubeconfig := string(kubeconfigBytes)
	//config, err := clientcmd.NewClientConfigFromBytes([]byte(kubeconfig))

	// Get kubeconfig string by cluster ID
	// 通过集群 ID 获取 kubeconfig 字符串
	kubeconfigstr, err := GetKubeConfigByClusterId(clusterid)
	kubeconfigBytes := []byte(kubeconfigstr)
	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] NewClientConfigFromBytes err:%s\n", err)
	}

	config, err := clientconfig.ClientConfig()
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] ClientSet ClientConfig err:%s\n", err)
	}
	// Set QPS and Burst limits for the client
	// 设置客户端的 QPS 和 Burst 限制
	config.QPS = 10
	config.Burst = 100
	// Get bearer token by cluster ID for authentication
	// 通过集群 ID 获取用于认证的 bearer token
	bearerToken, err := GetBearerTokenByClusterId(clusterid)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] GetBearerTokenByClusterId err:%s\n", err)
	}
	//log.Println(bearerToken)
	if bearerToken != "" && bearerToken != "null" {
		config.BearerToken = bearerToken
	}

	// create the clientset
	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] ClientSet NewForConfig err:%s\n", err)
	}

	// Duplicate code block - getting bearer token again
	// 重复代码块 - 再次获取 bearer token
	bearerToken, err2 := GetBearerTokenByClusterId(clusterid)
	if err2 != nil {
		//panic(err.Error())
		log.Printf("[ERROR] GetBearerTokenByClusterId err:%s\n", err2)
	}
	//log.Println(bearerToken)
	if bearerToken != "" && bearerToken != "null" {
		config.BearerToken = bearerToken
	}

	return clientset
}

// ClientSetConfig creates and returns both a kubernetes clientset and its configuration
// ClientSetConfig 创建并返回 kubernetes clientset 及其配置
func ClientSetConfig(clusterid string) (*kubernetes.Clientset, *rest.Config) {

	var clientset = &kubernetes.Clientset{}
	var config = &rest.Config{}
	var err2 error

	// Check if cluster ID is empty
	// 检查集群 ID 是否为空
	if clusterid == "" {
		log.Printf("[ERROR] ClientSetConfig ClusterId null:%s\n", clusterid)
		return clientset, config
	}

	// Get kubeconfig string by cluster ID
	// 通过集群 ID 获取 kubeconfig 字符串
	kubeconfigstr, err := GetKubeConfigByClusterId(clusterid)
	if err != nil {
		log.Printf("[ERROR] GetKubeConfigByClusterId error:%s\n", err.Error())
		return clientset, config
	}

	kubeconfigBytes := []byte(kubeconfigstr)
	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] NewClientConfigFromBytes error:%s\n", err.Error())
		return clientset, config
	}

	config, err2 = clientconfig.ClientConfig()
	if err2 != nil {
		//panic(err.Error())
		log.Printf("[ERROR] ClientConfig error:%s\n", err2.Error())
		return clientset, config
	}
	// Set QPS and Burst limits for the client
	// 设置客户端的 QPS 和 Burst 限制
	config.QPS = 10
	config.Burst = 100
	// create the clientset
	// 创建 clientset
	clientset, err2 = kubernetes.NewForConfig(config)
	if err2 != nil {
		//panic(err.Error())
		log.Printf("[ERROR] NewForConfig error:%s\n", err2.Error())
		return clientset, config
	}
	return clientset, config
}

// RestClient creates and returns a REST client and its configuration
// RestClient 创建并返回 REST 客户端及其配置
func RestClient(clusterid string) (*rest.RESTClient, *rest.Config) {
	var kubeconfigPath string
	switch clusterid {
	case "zx-pcauto":
		kubeconfigPath = "./kubeconfig.txt"
	default:
		kubeconfigPath = "./kubeconfig.txt"
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] RestClient BuildConfigFromFlags err:%s\n", err.Error())
	}

	// Set group version for the REST client
	// 为 REST 客户端设置组版本
	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config.GroupVersion = &groupversion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = scheme.Codecs

	// Create the REST client
	// 创建 REST 客户端
	restclient, err := rest.RESTClientFor(config)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] RestClient RESTClientFor err:%s\n", err.Error())
	}
	return restclient, config
}

// DynamicClient creates and returns a dynamic client for Kubernetes
// DynamicClient 创建并返回 Kubernetes 的动态客户端
func DynamicClient(clusterid string) *dynamic.DynamicClient {

	// Get kubeconfig string by cluster ID
	// 通过集群 ID 获取 kubeconfig 字符串
	kubeconfigstr, err := GetKubeConfigByClusterId(clusterid)
	kubeconfigBytes := []byte(kubeconfigstr)
	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] DynamicClient NewClientConfigFromBytes err:%s\n", err.Error())
	}

	config, err := clientconfig.ClientConfig()
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] DynamicClient ClientConfig err:%s\n", err.Error())
	}

	// Create dynamic client
	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] DynamicClient NewForConfig err:%s\n", err.Error())
	}
	return dynamicClient
}

// ClientConfig creates and returns a REST configuration for Kubernetes client
// ClientConfig 创建并返回 Kubernetes 客户端的 REST 配置
func ClientConfig(clusterid string) *rest.Config {

	// Get kubeconfig string by cluster ID
	// 通过集群 ID 获取 kubeconfig 字符串
	kubeconfigstr, err := GetKubeConfigByClusterId(clusterid)
	kubeconfigBytes := []byte(kubeconfigstr)
	clientconfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] ClientConfig NewClientConfigFromBytes err:%s\n", err.Error())
	}

	config, err := clientconfig.ClientConfig()
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] ClientConfig ClientConfig err:%s\n", err.Error())
	}

	// Get bearer token by cluster ID for authentication
	// 通过集群 ID 获取用于认证的 bearer token
	bearerToken, err := GetBearerTokenByClusterId(clusterid)
	if err != nil {
		//panic(err.Error())
		log.Printf("[ERROR] ClientConfig GetBearerTokenByClusterId err:%s\n", err)
	}
	//log.Println(bearerToken)
	if bearerToken != "" && bearerToken != "null" {
		config.BearerToken = bearerToken
	}

	return config
}
