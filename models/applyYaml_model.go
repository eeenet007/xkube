// applyYaml_model.go
package models

import (
	"bytes"
	"context"
	"fmt"
	"io"

	//"log"
	"xkube/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	//"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	syaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	sigyaml "sigs.k8s.io/yaml"
)

// GtGVR Get GroupVersionResource from GroupVersionKind
// GtGVR 根据GroupVersionKind获取GroupVersionResource
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// gvk: GroupVersionKind object
// gvk: GroupVersionKind对象
// Returns:
//   - schema.GroupVersionResource: GroupVersionResource object
//   - error: error message
//
// 返回值:
//   - schema.GroupVersionResource: GroupVersionResource对象
//   - error: 错误信息
func GtGVR(kubeconfig string, gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	clientset := common.ClientSet(kubeconfig)
	gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return mapping.Resource, nil
}

// ApplyYaml Apply YAML configuration to Kubernetes cluster
// ApplyYaml 应用YAML配置到Kubernetes集群
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// applyYaml: YAML configuration content to apply
// applyYaml: 要应用的YAML配置内容
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
func ApplyYaml(kubeconfig, applyYaml string) error {
	dynameicclient := common.DynamicClient(kubeconfig)
	d := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(applyYaml), 4096)

	for {
		var rawObj runtime.RawExtension
		err := d.Decode(&rawObj)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decode is err %v", err)
		}

		obj, _, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return fmt.Errorf("rawobj is err%v", err)
		}

		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return fmt.Errorf("tounstructured is err %v", err)
		}

		unstructureObj := &unstructured.Unstructured{Object: unstructuredMap}
		gvr, err := GtGVR(kubeconfig, unstructureObj.GroupVersionKind())
		if err != nil {
			return err
		}

		nameSpace := unstructureObj.GetNamespace()
		if nameSpace == "" {
			nameSpace = "default"
		}

		unstructuredYaml, err := sigyaml.Marshal(unstructureObj)
		if err != nil {
			return fmt.Errorf("unable to marshal resource as yaml: %w", err)
		}
		_, getErr := dynameicclient.Resource(gvr).Namespace(nameSpace).Get(context.Background(), unstructureObj.GetName(), metav1.GetOptions{})
		if getErr != nil {
			_, createErr := dynameicclient.Resource(gvr).Namespace(nameSpace).Create(context.Background(), unstructureObj, metav1.CreateOptions{})
			if createErr != nil {
				return createErr
			}
		}

		force := true
		_, err = dynameicclient.Resource(gvr).
			Namespace(nameSpace).
			Patch(context.Background(),
				unstructureObj.GetName(),
				types.ApplyPatchType,
				unstructuredYaml, metav1.PatchOptions{
					FieldManager: unstructureObj.GetName(),
					Force:        &force,
				})

		if err != nil {
			return fmt.Errorf("unable to patch resource: %w", err)
		}
	}
	return nil

}

// CreateByYaml Create Kubernetes resources from YAML data
// CreateByYaml 从YAML数据创建Kubernetes资源
// kubeconfig: cluster configuration identifier
// kubeconfig: 集群配置信息标识符
// namespace: namespace to create resources in
// namespace: 创建资源的命名空间
// yamlData: YAML data as byte array
// yamlData: YAML数据字节数组
// Returns:
//   - error: error message
//
// 返回值:
//   - error: 错误信息
//
// 测试deploy/cronjob/service ok
func CreateByYaml(kubeconfig, namespace string, yamlData []byte) error {

	clientset := common.ClientSet(kubeconfig)
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlData), 100)
	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := syaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return err
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			return err
		}

		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

		gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
		if err != nil {
			return err
		}

		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			//log.Fatal(err)
			return err
		}

		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(namespace)
			}
			dri = common.DynamicClient(kubeconfig).Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = common.DynamicClient(kubeconfig).Resource(mapping.Resource)
		}

		obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			//log.Fatal(err)
			return err
		}

		fmt.Printf("%s/%s created", obj2.GetKind(), obj2.GetName())
	}
	return nil
}
