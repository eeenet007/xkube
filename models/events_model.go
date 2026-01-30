// events_model.go
// Event数据模型定义及操作函数 / Event data model definition and operation functions
package models

import (
	"context"
	"fmt"

	"log"
	"sort"
	"strings"

	//"time"
	"xkube/common"

	//v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/util/retry"
	//"sigs.k8s.io/yaml"
)

// Events 事件信息结构体 / Event information struct
type Events struct {
	EventType  string `json:"eventType"` //事件类型 warnning / Event type: warning
	Kind       string `json:"kind"`      //对象类型 pod or deployment / Object type: pod or deployment
	ObjName    string `json:"objName"`   //对象：pod的名称 / Object: pod name
	Message    string `json:"message"`   //信息 / Message
	Reason     string `json:"reason"`    //原因、内容 / Reason, content
	CreateTime string `json:"createTime"` //创建时间 / Creation time
}

// EventList 获取事件列表 / Get event list
// kubeconfig: k8s集群配置 / k8s cluster configuration
// namespace: 命名空间 / namespace
// kind: 对象类型 / object type
// objName: 对象名称 / object name
// limitd: 限制数量 / limit count
func EventList(kubeconfig, namespace, kind, objName string, limitd int64) ([]Events, error) {
	clientset := common.ClientSet(kubeconfig)
	var err error
	var events *corev1.EventList

	var listOptions = metav1.ListOptions{}
	if limitd > 0 {
		listOptions = metav1.ListOptions{
			Limit: limitd,
		}
	}

	if kind != "" && objName != "" {
		listOptions = metav1.ListOptions{
			FieldSelector: fmt.Sprintf("involvedObject.name=%s", objName),
			TypeMeta:      metav1.TypeMeta{Kind: kind},
		}
		events, err = clientset.CoreV1().Events(namespace).List(context.TODO(), listOptions)
	} else if kind != "" && objName == "" {
		listOptions = metav1.ListOptions{
			TypeMeta: metav1.TypeMeta{Kind: kind},
		}
		events, err = clientset.CoreV1().Events(namespace).List(context.TODO(), listOptions)
	} else {
		events, err = clientset.CoreV1().Events(namespace).List(context.TODO(), listOptions)
	}
	if err != nil {
		log.Printf("list deployment error, err:%v\n", err)
	}

	lessFunc := func(i, j int) bool {
		return events.Items[j].CreationTimestamp.Before(&events.Items[i].CreationTimestamp)
	}
	sort.SliceStable(events.Items, lessFunc)

	var bbb = make([]Events, 0)
	for _, item := range events.Items {

		if objName != "" && kind == "" {
			if !strings.Contains(item.InvolvedObject.Name, objName) {
				continue
			}
		}

		Items := &Events{
			EventType:  item.Type,
			Kind:       item.InvolvedObject.Kind,
			ObjName:    item.InvolvedObject.Name,
			Message:    item.Message,
			Reason:     item.Reason,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
		//log.Println(item)
	}

	return bbb, err
}

// EventListV2 获取事件列表V2版本 / Get event list version 2
// kubeconfig: k8s集群配置 / k8s cluster configuration
// namespace: 命名空间 / namespace
// kind: 对象类型 / object type
// objName: 对象名称 / object name
// limitd: 限制数量 / limit count
func EventListV2(kubeconfig, namespace, kind, objName string, limitd int64) ([]Events, error) {
	clientset := common.ClientSet(kubeconfig)
	var err error
	var events *corev1.EventList

	//now := time.Now()
	//dayAgo := now.Add(-24 * time.Hour)

	var listOptions = metav1.ListOptions{
		Limit: limitd,
		//FieldSelector: fmt.Sprintf("lastTimestamp>=%s", dayAgo.Format(time.RFC3339)), //会提示field label not supported: lastTimestamp / Will prompt field label not supported: lastTimestamp
	}
	events, err = clientset.CoreV1().Events(namespace).List(context.TODO(), listOptions)
	if err != nil {
		log.Printf("list event error, err:%v\n", err)
	}

	lessFunc := func(i, j int) bool {
		return events.Items[j].CreationTimestamp.Before(&events.Items[i].CreationTimestamp)
	}
	sort.SliceStable(events.Items, lessFunc)

	var bbb = make([]Events, 0)
	for _, item := range events.Items {

		if objName != "" {
			if !strings.Contains(item.InvolvedObject.Name, objName) {
				continue
			}
		}
		if objName == "" && kind != "" {
			if item.InvolvedObject.Kind != kind {
				continue
			}
		}
		Items := &Events{
			EventType:  item.Type,
			Kind:       item.InvolvedObject.Kind,
			ObjName:    item.InvolvedObject.Name,
			Message:    item.Message,
			Reason:     item.Reason,
			CreateTime: item.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		bbb = append(bbb, *Items)
		//log.Println(item)
	}

	return bbb, err
}