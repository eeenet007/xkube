// Package models provides data models and related functions for the xkube application
// 包models为xkube应用提供数据模型和相关函数
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"time"
	//"strconv"
	//"io/ioutil"
	"log"
	"xkube/common"

	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"

	batchv1 "k8s.io/api/batch/v1" //注意有的地方要用beta1版本才行
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

// Cronjob represents a Kubernetes CronJob resource
// Cronjob表示Kubernetes CronJob资源
type Cronjob struct {
	CronjobName      string `json:"cronjobName"`
	NameSpace        string `json:"nameSpace"`
	Labels           string `json:"labels"`
	Annotations      string `json:"annotations"`
	ImgUrl           string `json:"imgUrl"`
	CmdArgs          string `json:"cmdArgs"`
	Suspend          string `json:"suspend"`
	Schedule         string `json:"schedule"`
	Active           int    `json:"active"`
	CreateTime       string `json:"createTime"`
	LastScheduleTime string `json:"lastScheduleTime"`
}

// CronjobList retrieves a list of CronJobs from a Kubernetes cluster
// CronjobList从Kubernetes集群中检索CronJobs列表
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace to list CronJobs from, or all namespaces if empty
// cronjobName - filter CronJobs by name (substring match)
// labelsKey - filter by label key
// labelsValue - filter by label value
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - 列出CronJobs的命名空间，如果为空则列出所有命名空间
// cronjobName - 按名称过滤CronJobs（子字符串匹配）
// labelsKey - 按标签键过滤
// labelsValue - 按标签值过滤
//
// Returns:
// []Cronjob - list of CronJobs
// error - any error that occurred
//
// 返回:
// []Cronjob - CronJobs列表
// error - 发生的任何错误
func CronjobList(kubeconfig, namespace, cronjobName string, labelsKey, labelsValue string) ([]Cronjob, error) {
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}
	var bbb = make([]Cronjob, 0)

	//设置ListOptions
	// Set up ListOptions
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	//BatchV1 会提示the server could not find the requested resource,改成BatchV1beta1 后正常
	//xList, err := common.ClientSet(kubeconfig).BatchV1beta1().CronJobs(namespace).List(context.TODO(), listOptions) //zx-pcauto ok,yt-pcauto异常：the server could not find the requested resource
	// BatchV1 may report "the server could not find the requested resource", switching to BatchV1beta1 works
	// xList, err := common.ClientSet(kubeconfig).BatchV1beta1().CronJobs(namespace).List(context.TODO(), listOptions) //zx-pcauto ok,yt-pcauto exception: the server could not find the requested resource
	xList, err := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace).List(context.TODO(), listOptions) //yt-pcauto ok,zx-pcauto 异常：the server could not find the requested resource

	if err != nil {
		log.Printf("[ERROR] ListCronjobError err:%v", err)
		return bbb, err
	}

	for _, vv := range xList.Items {
		//搜索
		// Search
		if cronjobName != "" {
			if !strings.Contains(vv.Name, cronjobName) {
				continue
			}
		}

		var labelsStr, imgUrlStr, cmdArgsStr string
		for k1, v1 := range vv.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}
		for _, v2 := range vv.Spec.JobTemplate.Spec.Template.Spec.Containers {
			imgUrlStr += fmt.Sprintf("%s,", v2.Image)
			cmdArgsStr += fmt.Sprintf("%s %s,", strings.Join(v2.Command, " "), strings.Join(v2.Args, " "))
		}

		if len(imgUrlStr) > 0 {
			imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
		}

		if len(cmdArgsStr) > 0 {
			cmdArgsStr = cmdArgsStr[0 : len(cmdArgsStr)-1]
		}

		var lastScheduleTime string
		if vv.Status.LastScheduleTime != nil {
			lastScheduleTime = vv.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
		}
		//vv.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Args

		xItems := &Cronjob{
			CronjobName: vv.Name,
			NameSpace:   vv.Namespace,
			Labels:      labelsStr,
			ImgUrl:      imgUrlStr,
			CmdArgs:     cmdArgsStr,
			Suspend:     fmt.Sprintf("%v", *vv.Spec.Suspend),
			Schedule:    vv.Spec.Schedule,
			Active:      len(vv.Status.Active),
			CreateTime:  vv.CreationTimestamp.Format("2006-01-02 15:04:05"),
			//LastScheduleTime: vv.Status.LastScheduleTime.Format("2006-01-02 15:04:05"),
			LastScheduleTime: lastScheduleTime,
		}
		bbb = append(bbb, *xItems)
	}
	return bbb, err
}

// CronjobDetail retrieves detailed information about a specific CronJob
// CronjobDetail检索特定CronJob的详细信息
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace where the CronJob is located
// cronjobName - the name of the CronJob to retrieve
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - CronJob所在的命名空间
// cronjobName - 要检索的CronJob的名称
//
// Returns:
// *Cronjob - pointer to the Cronjob object with detailed information
// error - any error that occurred
//
// 返回:
// *Cronjob - 包含详细信息的Cronjob对象指针
// error - 发生的任何错误
func CronjobDetail(kubeconfig, namespace, cronjobName string) (*Cronjob, error) {
	//BatchV1 会提示the server could not find the requested resource,改成BatchV1beta1 后正常
	//yt-pcauto:BatchV1,zx-pcauto:BatchV1beta1
	// BatchV1 may report "the server could not find the requested resource", switching to BatchV1beta1 works
	// yt-pcauto:BatchV1,zx-pcauto:BatchV1beta1
	cjobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace)
	cronjob, err := cjobClient.Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		return &Cronjob{}, err
	}

	var labelsStr, imgUrlStr, cmdArgsStr string
	for k1, v1 := range cronjob.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}
	for _, v2 := range cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers {
		imgUrlStr += fmt.Sprintf("%s,", v2.Image)
		cmdArgsStr += fmt.Sprintf("%s %s,", strings.Join(v2.Command, " "), strings.Join(v2.Args, " "))
	}

	var annotationsStr string
	for kk, vv := range cronjob.ObjectMeta.Annotations {
		if strings.Contains(kk, "last-applied-configuration") {
			continue
		}
		annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	if len(imgUrlStr) > 0 {
		imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
	}

	if len(cmdArgsStr) > 0 {
		cmdArgsStr = cmdArgsStr[0 : len(cmdArgsStr)-1]
	}
	var lastScheduleTime string
	if cronjob.Status.LastScheduleTime != nil {
		lastScheduleTime = cronjob.Status.LastScheduleTime.Format("2006-01-02 15:04:05")
	}

	return &Cronjob{
		CronjobName:      cronjob.Name,
		NameSpace:        cronjob.Namespace,
		Labels:           labelsStr,
		Annotations:      annotationsStr,
		ImgUrl:           imgUrlStr,
		CmdArgs:          cmdArgsStr,
		Suspend:          fmt.Sprintf("%v", *cronjob.Spec.Suspend),
		Schedule:         cronjob.Spec.Schedule,
		Active:           len(cronjob.Status.Active),
		CreateTime:       cronjob.CreationTimestamp.Format("2006-01-02 15:04:05"),
		LastScheduleTime: lastScheduleTime,
	}, nil
}

// CronjobCreate creates a new CronJob in a Kubernetes cluster
// CronjobCreate在Kubernetes集群中创建新的CronJob
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// bodys - JSON byte array containing the CronJob configuration
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// bodys - 包含CronJob配置的JSON字节数组
//
// Returns:
// error - any error that occurred
//
// 返回:
// error - 发生的任何错误
func CronjobCreate(kubeconfig string, bodys []byte) error {
	gp := gjson.ParseBytes(bodys)

	clusterId := gp.Get("clusterId").String()
	if kubeconfig == "" {
		kubeconfig = clusterId
	}
	cronjobName := gp.Get("cronjobName").String()
	nameSpace := gp.Get("nameSpace").String()
	var pullPolicy corev1.PullPolicy
	imagePullPolicy := gp.Get("imagePullPolicy").String()
	switch imagePullPolicy {
	case "Never":
		pullPolicy = corev1.PullNever
	case "IfNotPresent":
		pullPolicy = corev1.PullIfNotPresent
	default:
		pullPolicy = corev1.PullAlways
	}
	imageUrl := gp.Get("imageUrl").String()

	resourceLimitCheck := gp.Get("resourceLimitCheck").String()
	periodCheck := gp.Get("periodCheck").String()
	taskSetCheck := gp.Get("taskSetCheck").String()

	schedule := gp.Get("schedule").String()
	if schedule == "" {
		schedule = "* * * * *"
	}

	labelsMap := map[string]string{
		"app": cronjobName,
	}
	for _, vv := range gp.Get("lables").Array() {
		labelsMap[vv.Get("key").String()] = vv.Get("value").String()
	}

	cronjob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cronjobName,
			Namespace: nameSpace,
			Labels:    labelsMap,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: schedule, // 定时任务表达式
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{ //v1版本有JobSpec,beta1版本没有
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:            cronjobName,
									Image:           imageUrl,
									ImagePullPolicy: pullPolicy,
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	commandStr := gp.Get("command").Str
	if commandStr != "" {
		commandArry := strings.Split(commandStr, ",")
		cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Command = commandArry
	}
	argsStr := gp.Get("args").Str
	if argsStr != "" {
		argsArry := strings.Split(argsStr, ",")
		cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Args = argsArry
	}

	if resourceLimitCheck == "on" {
		limit_cpu := gp.Get("limit_cpu").Str
		limit_mem := gp.Get("limit_mem").Str
		request_cpu := gp.Get("request_cpu").Str
		request_mem := gp.Get("request_mem").Str
		resReq := &corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(limit_cpu),
				corev1.ResourceMemory: resource.MustParse(limit_mem),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(request_cpu),
				corev1.ResourceMemory: resource.MustParse(request_mem),
			},
		}
		cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Resources = *resReq
	}

	if periodCheck == "on" {
		var postStartArry, preStopArry []string
		postStartArry = strings.Split(gp.Get("postStart").Str, ",")
		preStopArry = strings.Split(gp.Get("preStop").Str, ",")
		lifeCycle := &corev1.Lifecycle{}
		if len(postStartArry) > 0 {
			lifeCycle.PostStart = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: postStartArry,
				},
			}
		}
		if len(preStopArry) > 0 {
			lifeCycle.PreStop = &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: preStopArry,
				},
			}
		}
		cronjob.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Lifecycle = lifeCycle

	}

	if taskSetCheck == "on" {
		backoffLimit := gp.Get("backoffLimit").Int()
		parallelism := gp.Get("parallelism").Int() //保留失败的pod
		completions := gp.Get("completions").Int() //保留成功的pod
		activeDeadlineSeconds := gp.Get("activeDeadlineSeconds").Int()
		failedJobHistoryLimit := gp.Get("failedJobHistoryLimit").Int()         //保留失败的pod
		successfulJobHistoryLimit := gp.Get("successfulJobHistoryLimit").Int() //保留成功的pod
		concurrencyPolicy := gp.Get("concurrencyPolicy").Str

		failedJobHistoryLimitd := int32(failedJobHistoryLimit)
		successfulJobHistoryLimitd := int32(successfulJobHistoryLimit)
		cronjob.Spec.FailedJobsHistoryLimit = &failedJobHistoryLimitd
		cronjob.Spec.SuccessfulJobsHistoryLimit = &successfulJobHistoryLimitd

		var concPolicy batchv1.ConcurrencyPolicy
		switch concurrencyPolicy {
		case "Replace":
			concPolicy = batchv1.ReplaceConcurrent
		case "Forbid":
			concPolicy = batchv1.ForbidConcurrent
		default:
			concPolicy = batchv1.AllowConcurrent
		}
		cronjob.Spec.ConcurrencyPolicy = concPolicy

		parallelismd := int32(parallelism)
		completionsd := int32(completions)
		cronjob.Spec.JobTemplate.Spec.Parallelism = &parallelismd
		cronjob.Spec.JobTemplate.Spec.Completions = &completionsd
		cronjob.Spec.JobTemplate.Spec.ActiveDeadlineSeconds = &activeDeadlineSeconds
		backoffLimitd := int32(backoffLimit)
		cronjob.Spec.JobTemplate.Spec.BackoffLimit = &backoffLimitd

		var rePolicy corev1.RestartPolicy
		restartPolicy := gp.Get("restartPolicy").String()
		if restartPolicy == "OnFailure" {
			rePolicy = corev1.RestartPolicyOnFailure
		} else {
			rePolicy = corev1.RestartPolicyNever
		}
		cronjob.Spec.JobTemplate.Spec.Template.Spec.RestartPolicy = rePolicy
	}

	jobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(nameSpace)
	_, err := jobClient.Create(context.TODO(), cronjob, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// CronjobModify modifies an existing CronJob in a Kubernetes cluster
// CronjobModify修改Kubernetes集群中现有的CronJob
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace where the CronJob is located
// cronjobName - the name of the CronJob to modify
// key - the attribute to modify (schedule, image, suspend)
// value - the new value for the attribute
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - CronJob所在的命名空间
// cronjobName - 要修改的CronJob的名称
// key - 要修改的属性（schedule, image, suspend）
// value - 属性的新值
//
// Returns:
// error - any error that occurred
//
// 返回:
// error - 发生的任何错误
func CronjobModify(kubeconfig, namespace, cronjobName, key, value string) error {
	jobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace)
	jobInstance, err := jobClient.Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	switch key {
	case "schedule":
		jobInstance.Spec.Schedule = value
	case "image":
		jobInstance.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = value
	case "suspend":
		var isSuspend bool
		if value == "true" {
			isSuspend = true
		}
		jobInstance.Spec.Suspend = &isSuspend
	default:
		return fmt.Errorf("noSupport Args")
	}

	_, err = jobClient.Update(context.TODO(), jobInstance, metav1.UpdateOptions{})

	if err != nil {
		return err
	}
	return nil
}

// CronjobYamlModify updates a CronJob using YAML data
// CronjobYamlModify使用YAML数据更新CronJob
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// yamlData - byte array containing the YAML representation of the CronJob
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// yamlData - 包含CronJob的YAML表示的字节数组
//
// Returns:
// error - any error that occurred
//
// 返回:
// error - 发生的任何错误
func CronjobYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	cronjob := &batchv1.CronJob{}
	err = json.Unmarshal(data, cronjob)
	if err != nil {
		return err
	}

	namespace := cronjob.ObjectMeta.Namespace
	//cronjobName := cronjob.ObjectMeta.Name
	jobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace)
	_, err = jobClient.Update(context.TODO(), cronjob, metav1.UpdateOptions{})
	return err
}

// CronjobClone clones a CronJob to another namespace or cluster
// CronjobClone将CronJob克隆到另一个命名空间或集群
//
// Parameters:
// kubeconfig - the kubeconfig of the source cluster
// namespace - the namespace of the source CronJob
// objname - the name of the source CronJob
// target_clusterid - the kubeconfig of the target cluster
// target_namespace - the target namespace
// target_objname - the name for the cloned CronJob
//
// 参数:
// kubeconfig - 源集群的kubeconfig
// namespace - 源CronJob的命名空间
// objname - 源CronJob的名称
// target_clusterid - 目标集群的kubeconfig
// target_namespace - 目标命名空间
// target_objname - 克隆的CronJob的名称
//
// Returns:
// error - any error that occurred
//
// 返回:
// error - 发生的任何错误
func CronjobClone(kubeconfig, namespace, objname, target_clusterid, target_namespace, target_objname string) error {
	cronjob, err := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace).Get(context.TODO(), objname, metav1.GetOptions{})
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
	// Create namespace in the new cluster
	err3 := CreateNsByExist(target_clusterid, target_namespace)
	if err3 != nil {
		return err3
	}

	cronjob.Name = target_objname
	cronjob.Namespace = target_namespace
	cronjob.ResourceVersion = ""
	cronjob.ObjectMeta.Labels["app"] = target_objname

	//new cluster
	// 新集群
	NewClient := common.ClientSet(target_clusterid).BatchV1().CronJobs(target_namespace)
	_, err2 := NewClient.Get(context.TODO(), target_objname, metav1.GetOptions{})
	if errors.IsNotFound(err2) { //没有就创建
		// Create if not exists
		_, err := NewClient.Create(context.TODO(), cronjob, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	if err2 == nil {
		_, err = NewClient.Update(context.TODO(), cronjob, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	}
	return err2
}

// CronjobDel deletes a CronJob from a Kubernetes cluster
// CronjobDel从Kubernetes集群中删除CronJob
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace where the CronJob is located
// cronjobName - the name of the CronJob to delete
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - CronJob所在的命名空间
// cronjobName - 要删除的CronJob的名称
//
// Returns:
// error - any error that occurred
//
// 返回:
// error - 发生的任何错误
func CronjobDel(kubeconfig, namespace, cronjobName string) error {
	jobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace)
	err := jobClient.Delete(context.TODO(), cronjobName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetCronjobYaml retrieves the YAML representation of a CronJob
// GetCronjobYaml检索CronJob的YAML表示
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace where the CronJob is located
// cronjobName - the name of the CronJob to retrieve
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - CronJob所在的命名空间
// cronjobName - 要检索的CronJob的名称
//
// Returns:
// string - YAML representation of the CronJob
// error - any error that occurred
//
// 返回:
// string - CronJob的YAML表示
// error - 发生的任何错误
func GetCronjobYaml(kubeconfig, namespace, cronjobName string) (string, error) {

	//yt-pcauto:BatchV1,zx-pcauto:BatchV1beta1
	jobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace)
	jobInstance, err := jobClient.Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	jobUnstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(jobInstance)
	if err != nil {
		return "", err
	}
	yamlBytes, err := yaml.Marshal(jobUnstructured)
	if err != nil {
		return "", err
	}
	//fmt.Println(string(yamlBytes))
	return string(yamlBytes), nil

}

// CronjobLabels manages labels for a CronJob
// CronjobLabels管理CronJob的标签
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace where the CronJob is located
// cronjobName - the name of the CronJob
// method - HTTP method (POST to update labels, other to get labels)
// labelsMap - map of labels to set (used when method is POST)
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - CronJob所在的命名空间
// cronjobName - CronJob的名称
// method - HTTP方法（POST更新标签，其他获取标签）
// labelsMap - 要设置的标签映射（当method为POST时使用）
//
// Returns:
// []LabelsKv - list of labels (when method is not POST)
// error - any error that occurred
//
// 返回:
// []LabelsKv - 标签列表（当method不是POST时）
// error - 发生的任何错误
func CronjobLabels(kubeconfig, namespace, cronjobName string, method string, labelsMap map[string]string) ([]LabelsKv, error) {
	var bbb = make([]LabelsKv, 0)
	jobClient := common.ClientSet(kubeconfig).BatchV1().CronJobs(namespace)
	cronjob, err := jobClient.Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] DeployLabels Get error:%s\n", err)
		return bbb, err
	}
	if method == "POST" {
		cronjob.ObjectMeta.Labels = labelsMap
		_, err = jobClient.Update(context.TODO(), cronjob, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("[ERROR] DeployLabels Update error:%s\n", err)
			return bbb, err
		}
		return bbb, nil
	} else {
		if cronjob.ObjectMeta.Labels != nil {
			for kk, vv := range cronjob.ObjectMeta.Labels {
				bbb = append(bbb, LabelsKv{
					Key:   kk,
					Value: vv,
				})
			}
		}
		return bbb, nil
	}
}

// CronjobRun manually triggers a CronJob to run immediately
// CronjobRun手动触发CronJob立即运行
//
// Parameters:
// kubeconfig - the kubeconfig to use for connecting to the cluster
// namespace - the namespace where the CronJob is located
// cronjobName - the name of the CronJob to trigger
//
// 参数:
// kubeconfig - 用于连接集群的kubeconfig
// namespace - CronJob所在的命名空间
// cronjobName - 要触发的CronJob的名称
//
// Returns:
// error - any error that occurred
//
// 返回:
// error - 发生的任何错误
func CronjobRun(kubeconfig, namespace, cronjobName string) error {
	clientset := common.ClientSet(kubeconfig)
	cronJob, err := clientset.BatchV1().CronJobs(namespace).Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 创建Job的名字
	// Create Job name
	jobName := fmt.Sprintf("%s-%s", cronJob.Name, time.Now().Format("20060102150405"))
	// 创建Job
	// Create Job
	newJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: cronJob.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: cronJob.Spec.JobTemplate.Spec.Template,
		},
	}

	// 设置Job的OwnerReference为CronJob，以便于跟踪
	// Set Job's OwnerReference to CronJob for tracking
	newJob.OwnerReferences = []metav1.OwnerReference{
		*metav1.NewControllerRef(cronJob, batchv1.SchemeGroupVersion.WithKind("CronJob")),
	}

	// 创建Job
	// Create Job
	_, err = clientset.BatchV1().Jobs(cronJob.Namespace).Create(context.TODO(), newJob, metav1.CreateOptions{})
	return err
}