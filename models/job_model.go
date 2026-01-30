// Package models provides data models and related functions for Kubernetes resources
// 包models提供了Kubernetes资源的数据模型和相关函数
package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"xkube/common"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// Job represents basic information about a Kubernetes Job
// Job表示Kubernetes Job的基本信息
type Job struct {
	JobName   string `json:"jobName"`   // Job名称
	NameSpace string `json:"nameSpace"` // 命名空间
	Status    string `json:"status"`    // 状态
	//PodStatus     string `json:"podStatus"`     // Pod状态
	Labels        string `json:"labels"`        // 标签
	ContainerName string `json:"containerName"` // 容器名称
	ImgUrl        string `json:"imgUrl"`        // 镜像地址
	CreateTime    string `json:"createTime"`    // 创建时间
	CompleteTime  string `json:"completeTime"`  // 完成时间
}

// JobDetails represents detailed information about a Kubernetes Job
// JobDetails表示Kubernetes Job的详细信息
type JobDetails struct {
	JobName       string          `json:"jobName"`       // Job名称
	NameSpace     string          `json:"nameSpace"`     // 命名空间
	Status        string          `json:"status"`        // 状态
	PodStatus     string          `json:"podStatus"`     // Pod状态
	Labels        string          `json:"labels"`        // 标签
	Annotations   string          `json:"annotations"`   // 注解
	ContainerName string          `json:"containerName"` // 容器名称
	ImgUrl        string          `json:"imgUrl"`        // 镜像地址
	CreateTime    string          `json:"createTime"`    // 创建时间
	CompleteTime  string          `json:"completeTime"`  // 完成时间
	Conditions    []JobConditions `json:"conditions"`    // 条件状态列表
}

// JobConditions represents the condition status of a Job
// JobConditions表示Job的条件状态
type JobConditions struct {
	Ctype              string `json:"ctype"`              // 条件类型
	Status             string `json:"status"`             // 状态
	Reason             string `json:"reason"`             // 原因
	Message            string `json:"message"`            // 消息
	LastTransitionTime string `json:"lastTransitionTime"` // 最后转换时间
	LastProbeTime      string `json:"lastProbeTime"`      // 最后探测时间
}

// NewJob represents the parameters needed to create a new Job
// NewJob表示创建新Job所需的参数
type NewJob struct {
	JobName       string            `json:"jobName"`       // Job名称
	NameSpace     string            `json:"nameSpace"`     // 命名空间
	Labels        map[string]string `json:"labels"`        // 标签
	ContainerName string            `json:"containerName"` // 容器名称
	ImgUrl        string            `json:"imgUrl"`        // 镜像地址
	BackoffLimit  *int32            `json:"backoffLimit"`  // 允许执行失败的次数，默认值是6，0表示不允许执行失败；
	Command       []string          `json:"command"`       // 命令
	Args          []string          `json:"args"`          // 参数
	//RestartPolicy string `json:"restartPolicy"` //Nerver，则失败后会创建新的Pod，如果是OnFailed，则会重启Pod
	//activeDeadlineSeconds   int64 //Job的超时时间，一旦一个Job运行的时间超出该限制，则Job失败，所有运行中的Pod会被结束并删除；
	//ttlSecondsAfterFinished int64 //当设置了ttlSecondsAfterFinished参数，job完成或者失败后都会在ttlSecondsAfterFinished所设置的时间后被清理；
}

// JobList retrieves a list of Jobs based on specified criteria
// JobList根据指定条件检索Job列表
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - namespace: 命名空间，如果为空则查询所有命名空间
//   - jobName: Job名称，用于过滤
//   - labelsKey: 标签键
//   - labelsValue: 标签值
//
// Returns:
//   - []Job: Job列表
//   - error: 错误信息
func JobList(kubeconfig, namespace, jobName, labelsKey, labelsValue string) ([]Job, error) {
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	//设置ListOptions
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	var bbb = make([]Job, 0)
	xList, err := common.ClientSet(kubeconfig).BatchV1().Jobs(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Printf("[ERROR] ListJobsError err:%v", err)
		return bbb, err
	}

	for _, vv := range xList.Items {
		//搜索
		if jobName != "" {
			if !strings.Contains(vv.Name, jobName) {
				continue
			}
		}

		var labelsStr, imgUrlStr, containerNameStr string
		for k1, v1 := range vv.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		for _, v2 := range vv.Spec.Template.Spec.Containers {
			containerNameStr += fmt.Sprintf("%s,", v2.Name)
			imgUrlStr += fmt.Sprintf("%s,", v2.Image)
		}
		if len(containerNameStr) > 0 {
			containerNameStr = containerNameStr[0 : len(containerNameStr)-1]
		}
		if len(imgUrlStr) > 0 {
			imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
		}

		var xStatus = "Unknown"
		if vv.Status.Active > 0 {
			xStatus = "Active"
		} else if vv.Status.Succeeded > 0 {
			xStatus = "Succeeded"
		} else if vv.Status.Failed > 0 {
			xStatus = "Failed"
		}

		var completeTime = ""
		if vv.Status.Succeeded > 0 {
			completeTime = vv.Status.CompletionTime.Format("2006-01-02 15:04:05")
		}

		xItems := &Job{
			JobName:   vv.Name,
			NameSpace: vv.Namespace,
			Status:    xStatus,
			//PodStatus:     fmt.Sprintf("活跃:%d<br>成功:%d<br>失败:%d", vv.Status.Active, vv.Status.Succeeded, vv.Status.Failed),
			Labels:        labelsStr,
			ContainerName: containerNameStr,
			ImgUrl:        imgUrlStr,
			CreateTime:    vv.Status.StartTime.Format("2006-01-02 15:04:05"),
			CompleteTime:  completeTime,
		}
		bbb = append(bbb, *xItems)
	}
	return bbb, err
}

// JobListv2 retrieves a list of Jobs with additional CronJob filtering
// JobListv2检索Job列表，支持额外的CronJob过滤条件
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - namespace: 命名空间，如果为空则查询所有命名空间
//   - cronjobName: CronJob名称，用于过滤由特定CronJob创建的Job
//   - jobName: Job名称，用于过滤
//   - labelsKey: 标签键
//   - labelsValue: 标签值
//
// Returns:
//   - []Job: Job列表
//   - error: 错误信息
func JobListv2(kubeconfig, namespace, cronjobName, jobName, labelsKey, labelsValue string) ([]Job, error) {
	if namespace == "" {
		//namespace = corev1.NamespaceDefault
		namespace = corev1.NamespaceAll
	}

	//设置ListOptions
	var listOptions = metav1.ListOptions{}
	if labelsKey != "" && labelsValue != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", labelsKey, labelsValue),
		}
	}

	var bbb = make([]Job, 0)
	xList, err := common.ClientSet(kubeconfig).BatchV1().Jobs(namespace).List(context.TODO(), listOptions)

	if err != nil {
		log.Printf("[ERROR] ListJobsError err:%v", err)
		return bbb, err
	}

	for _, vv := range xList.Items {
		//搜索
		if jobName != "" {
			if !strings.Contains(vv.Name, jobName) {
				continue
			}
		}
		//指定cronjob查询
		if cronjobName != "" {
			var isok bool
			for _, ref := range vv.ObjectMeta.OwnerReferences {
				if ref.Name == cronjobName {
					isok = true
				} else {
					continue
				}
			}
			if !isok {
				continue
			}
		}

		var labelsStr, imgUrlStr, containerNameStr string
		for k1, v1 := range vv.ObjectMeta.Labels {
			labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
		}
		if len(labelsStr) > 0 {
			labelsStr = labelsStr[0 : len(labelsStr)-1]
		}

		for _, v2 := range vv.Spec.Template.Spec.Containers {
			containerNameStr += fmt.Sprintf("%s,", v2.Name)
			imgUrlStr += fmt.Sprintf("%s,", v2.Image)
		}
		if len(containerNameStr) > 0 {
			containerNameStr = containerNameStr[0 : len(containerNameStr)-1]
		}
		if len(imgUrlStr) > 0 {
			imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
		}
		var xStatus = "Unknown"
		if vv.Status.Active > 0 {
			xStatus = "Active"
		} else if vv.Status.Succeeded > 0 {
			xStatus = "Succeeded"
		} else if vv.Status.Failed > 0 {
			xStatus = "Failed"
		}

		var completeTime = ""
		if vv.Status.Succeeded > 0 && vv.Status.CompletionTime != nil {
			completeTime = vv.Status.CompletionTime.Format("2006-01-02 15:04:05")
		}

		xItems := &Job{
			JobName:   vv.Name,
			NameSpace: vv.Namespace,
			Status:    xStatus,
			//PodStatus:     fmt.Sprintf("活跃:%d<br>成功:%d<br>失败:%d", vv.Status.Active, vv.Status.Succeeded, vv.Status.Failed),
			Labels:        labelsStr,
			ContainerName: containerNameStr,
			ImgUrl:        imgUrlStr,
			CreateTime:    vv.Status.StartTime.Format("2006-01-02 15:04:05"),
			CompleteTime:  completeTime,
		}
		bbb = append(bbb, *xItems)
	}
	return bbb, err
}

// JobDetail retrieves detailed information about a specific Job
// JobDetail检索特定Job的详细信息
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - namespace: 命名空间
//   - jobName: Job名称
//
// Returns:
//   - *JobDetails: Job详细信息
//   - error: 错误信息
func JobDetail(kubeconfig, namespace, jobName string) (*JobDetails, error) {
	//调用接口对象JobmentClient中的Get方法，获取相应的deployment资源数据
	job, err := common.ClientSet(kubeconfig).BatchV1().Jobs(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return &JobDetails{}, err
	}

	var labelsStr, imgUrlStr, containerNameStr, annotationsStr string
	for k1, v1 := range job.ObjectMeta.Labels {
		labelsStr += fmt.Sprintf("%s:%s,", k1, v1)
	}
	if len(labelsStr) > 0 {
		labelsStr = labelsStr[0 : len(labelsStr)-1]
	}

	for kk, vv := range job.ObjectMeta.Annotations {
		annotationsStr += fmt.Sprintf("%s:%s,", kk, vv)
	}
	if len(annotationsStr) > 0 {
		annotationsStr = annotationsStr[0 : len(annotationsStr)-1]
	}

	for _, v2 := range job.Spec.Template.Spec.Containers {
		containerNameStr += fmt.Sprintf("%s,", v2.Image)
		imgUrlStr += fmt.Sprintf("%s,", v2.Image)
	}
	if len(containerNameStr) > 0 {
		containerNameStr = containerNameStr[0 : len(containerNameStr)-1]
	}
	if len(imgUrlStr) > 0 {
		imgUrlStr = imgUrlStr[0 : len(imgUrlStr)-1]
	}

	var xStatus = "未知"
	if job.Status.Active > 0 {
		xStatus = "运行中"
	} else if job.Status.Succeeded > 0 {
		xStatus = "成功"
	} else if job.Status.Failed > 0 {
		xStatus = "失败"
	}

	var bbb = make([]JobConditions, 0)
	for _, v1 := range job.Status.Conditions {
		xItems := &JobConditions{
			LastTransitionTime: v1.LastTransitionTime.Format("2006-01-02 15:04:05"),
			LastProbeTime:      v1.LastProbeTime.Format("2006-01-02 15:04:05"),
			Message:            v1.Message,
			Reason:             v1.Reason,
			Status:             fmt.Sprintf("%v", v1.Status),
			Ctype:              fmt.Sprintf("%v", v1.Type),
		}
		bbb = append(bbb, *xItems)
	}
	var completeTime string
	if job.Status.CompletionTime != nil {
		completeTime = job.Status.CompletionTime.Format("2006-01-02 15:04:05")
	}

	return &JobDetails{
		JobName:       job.Name,
		NameSpace:     job.Namespace,
		Status:        xStatus,
		PodStatus:     fmt.Sprintf("活跃:%d,成功:%d,失败:%d", job.Status.Active, job.Status.Succeeded, job.Status.Failed),
		Labels:        labelsStr,
		Annotations:   annotationsStr,
		ContainerName: containerNameStr,
		ImgUrl:        imgUrlStr,
		CreateTime:    job.Status.StartTime.Format("2006-01-02 15:04:05"),
		CompleteTime:  completeTime,
		Conditions:    bbb,
	}, nil
}

// JobCreate creates a new Job based on the provided parameters
// JobCreate根据提供的参数创建一个新的Job
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - njob: 新Job的参数
//
// Returns:
//   - error: 错误信息
func JobCreate(kubeconfig string, njob *NewJob) error {
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      njob.JobName,
			Namespace: njob.NameSpace,
			Labels:    njob.Labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    njob.JobName,
							Image:   njob.ImgUrl,
							Command: njob.Command,
							Args:    njob.Args,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
			BackoffLimit: njob.BackoffLimit,
		},
	}
	jobClient := common.ClientSet(kubeconfig).BatchV1().Jobs(njob.NameSpace)
	_, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// JobModify modifies an existing Job
// JobModify修改现有的Job
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - newJob: 要修改的Job信息
//
// Returns:
//   - error: 错误信息
func JobModify(kubeconfig string, newJob *Job) error {
	jobClient := common.ClientSet(kubeconfig).BatchV1().Jobs(newJob.NameSpace)
	jobInstance, err := jobClient.Get(context.TODO(), newJob.JobName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	jobInstance.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyNever

	_, err = jobClient.Update(context.TODO(), jobInstance, metav1.UpdateOptions{})

	if err != nil {
		return err
	}

	return nil

}

// JobDel deletes a Job
// JobDel删除一个Job
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - namespace: 命名空间
//   - jobName: Job名称
//
// Returns:
//   - error: 错误信息
func JobDel(kubeconfig, namespace, jobName string) error {
	jobClient := common.ClientSet(kubeconfig).BatchV1().Jobs(namespace)
	err := jobClient.Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetJobYaml retrieves the YAML representation of a Job
// GetJobYaml获取Job的YAML表示
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - namespace: 命名空间
//   - jobName: Job名称
//
// Returns:
//   - string: Job的YAML字符串表示
//   - error: 错误信息
func GetJobYaml(kubeconfig, namespace, jobName string) (string, error) {

	jobClient := common.ClientSet(kubeconfig).BatchV1().Jobs(namespace)
	jobInstance, err := jobClient.Get(context.TODO(), jobName, metav1.GetOptions{})
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

// http://127.0.0.1:8001/xkube/job/v1/Log?clusterId=yt-pcauto&nameSpace=zx-app&resType=job&resName=sync-globalssi-28362542
// JobLog retrieves logs for a Job
// JobLog获取Job的日志
// Parameters:
//   - kubeconfig: Kubernetes配置文件路径
//   - nameSpace: 命名空间
//   - resName: 资源名称
//   - encode: 编码格式
//
// Returns:
//   - string: 日志内容或错误信息
func JobLog(kubeconfig, nameSpace, resName, encode string) string {
	clientset := common.ClientSet(kubeconfig)
	xjob, err := JobList(kubeconfig, nameSpace, "", "job-name", resName)
	if err != nil {
		return fmt.Sprintf("getJobListError:%s", err)
	}
	if len(xjob) <= 0 {
		return "getJobListFail"
	}
	//log.Println(xjob)
	jobName := xjob[0].JobName
	//设置ListOptions
	var listOptions = metav1.ListOptions{}
	if jobName != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", jobName),
		}
	}
	var podName string
	podxList, err := clientset.CoreV1().Pods(nameSpace).List(context.Background(), listOptions)
	if err != nil {
		return fmt.Sprintf("getPodListError:%s", err)
	}
	//log.Println(podxList.Items)
	podNumber := len(podxList.Items)
	if podNumber > 0 {
		podName = podxList.Items[podNumber-1].Name
	} else {
		return fmt.Sprintf("%s:NoPod", jobName)
	}
	//log.Println(podName)
	logStr := PodLog(kubeconfig, nameSpace, podName, "", 0, encode)
	return logStr
}

func JobYamlModify(kubeconfig string, yamlData []byte) error {
	data, err := yamlutil.ToJSON(yamlData)
	if err != nil {
		return err
	}
	job := &batchv1.Job{}
	err = json.Unmarshal(data, job)
	if err != nil {
		return err
	}

	namespace := job.ObjectMeta.Namespace
	//cronjobName := cronjob.ObjectMeta.Name
	jobClient := common.ClientSet(kubeconfig).BatchV1().Jobs(namespace)
	_, err = jobClient.Update(context.TODO(), job, metav1.UpdateOptions{})
	return err
}
