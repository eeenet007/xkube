// search_model.go
package models

import (
	"encoding/json"
	"fmt"
	"log"
	"xkube/common"

	"github.com/beego/beego/v2/client/orm"
)

type Xkb_index struct { //查询时会将大写以下划线分开
	Id         int64  `json:"id"`
	ClusterId  string `json:"clusterId"`
	ResType    string `json:"resType"`
	ResName    string `json:"resName"`
	NameSpace  string `json:"nameSpace"`
	Message    string `json:"message"`
	Createtime string `json:"createtime"`
}

func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Xkb_index))
}

// InsertIndex 插入或更新索引记录
// InsertIndex inserts or updates an index record
// clusterId: 集群ID / Cluster ID
// resType: 资源类型 / Resource type
// resName: 资源名称 / Resource name
// nameSpace: 命名空间 / Namespace
// message: 消息内容 / Message content
// 返回插入记录数和错误信息 / Returns number of inserted records and error information
func InsertIndex(clusterId, resType, resName, nameSpace, message string) (int64, error) {
	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("INSERT INTO xkb_index(cluster_id,res_type,res_name,name_space,message)VALUE('%s','%s','%s','%s','%s') ON DUPLICATE KEY UPDATE message = VALUES(message)", clusterId, resType, resName, nameSpace, message)
	//log.Println(sqlstr)
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Println("[INFO] InsertIndex affected nums: ", num)
		return num, nil
	}
	return 0, err
}

func GetList_Index(keyword string) (xindex []*Xkb_index, count int64) {
	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("select * from xkb_index where res_name like '%%%s%%' OR message LIKE '%%%s%%'", keyword, keyword)
	num, err := o.Raw(sqlstr).QueryRows(&xindex)
	if err != nil {
		log.Println(err)
	}
	return xindex, num
}

func DelById_Index(id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Xkb_index{Id: id})
	return status, err
}

func DeployInsert(clusterId string) (int, error) {
	xList, err := DeployList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] DeployInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.DeployName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,image:%s,port:%s", vv.Labels, vv.ImageUrl, vv.ContainerPort)
		_, _ = InsertIndex(clusterId, "deploy", resName, NameSpace, msg)
	}
	//统计熟练并写入缓存用于首页显示
	_ = common.HSet("count_"+clusterId, "deploy", fmt.Sprintf("%d", len(xList)))

	bodystr, err := json.Marshal(&xList)
	if err == nil {
		_ = common.SetEx("deployList"+clusterId, string(bodystr), 43200)
	}

	return len(xList), nil
}

func StsInsert(clusterId string) (int, error) {
	xList, err := StatefulsetList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] stsInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.StatefulsetName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,image:%s,port:%s", vv.Labels, vv.ImgUrl, vv.ContainerPort)
		_, _ = InsertIndex(clusterId, "sts", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "sts", fmt.Sprintf("%d", len(xList)))
	return len(xList), nil
}

func DsInsert(clusterId string) (int, error) {
	xList, err := DaemonsetList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] dsInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.DaemonsetName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,image:%s,port:%s", vv.Labels, vv.ImgUrl, vv.ContainerPort)
		_, _ = InsertIndex(clusterId, "ds", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "ds", fmt.Sprintf("%d", len(xList)))
	return len(xList), nil
}

func CronjobInsert(clusterId string) (int, error) {
	xList, err := CronjobList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] cronjobInsert-v1 Get List Fail:%s\n", err)
		xList, err = CronjobListBeta1(clusterId, "", "", "", "")
		if err != nil {
			log.Printf("[ERROR] cronjobInsert-Beta1 Get List Fail:%s\n", err)
			return 0, err
		}
	}
	for _, vv := range xList {
		resName := vv.CronjobName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,image:%s,args:%s", vv.Labels, vv.ImgUrl, vv.CmdArgs)
		_, _ = InsertIndex(clusterId, "cronjob", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "cronjob", fmt.Sprintf("%d", len(xList)))
	return len(xList), nil
}

func PodInsert(clusterId string) (int, error) {
	xList, err := PodListV2(clusterId, "", "", "", "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] podInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.PodName
		NameSpace := vv.NameSpace
		dt, _ := PodDetail(clusterId, NameSpace, resName)
		var ctName, ctImage, ctVolumes, ctEnv, ctPort string
		for _, vc := range dt.Containers {
			ctName += vc.ContainerName + ","
			ctImage += vc.ContainerImage + ","
			ctVolumes += vc.Mounts + ","
			ctEnv += vc.Envs + ","
			ctPort += vc.Ports + ","
		}
		//msg := fmt.Sprintf("labels:%s,node:%s,podip:%s,nodeIp:%s,container:%s image:%s volumes:%s env:%s port:%s", vv.Labels, vv.NodeName, vv.PodIp, vv.HostIp, ctName, ctImage, ctVolumes, ctEnv, ctPort)
		msg := fmt.Sprintf("labels:%s,node:%s,podip:%s,nodeIp:%s", vv.Labels, vv.NodeName, vv.PodIp, vv.HostIp)
		_, _ = InsertIndex(clusterId, "pod", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "pod", fmt.Sprintf("%d", len(xList)))

	//将全量列表写入缓存
	bodystr, err := json.Marshal(&xList)
	if err == nil {
		_ = common.SetEx("podList"+clusterId, string(bodystr), 86400)
	}

	return len(xList), nil
}

func SvcInsert(clusterId string) (int, error) {
	xList, err := SvcList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] svcInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.ServiceName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,svcIp:%s,svcPort:%s,lanEndpoint:%s,wanEndpoint:%s", vv.Labels, vv.SvcIp, vv.SvcPort, vv.LanEndpoint, vv.WanEndpoint)
		_, _ = InsertIndex(clusterId, "service", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "service", fmt.Sprintf("%d", len(xList)))

	//全量列表缓存到redis
	bodystr, err := json.Marshal(&xList)
	if err == nil {
		_ = common.SetEx("svcList"+clusterId, string(bodystr), 600)
	}
	return len(xList), nil
}

func IngressInsert(clusterId string) (int, error) {
	xList, err := IngList(clusterId, "", "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] ingressInsert-v1 Get List Fail:%s\n", err)
		xList, err = IngListV1beta1(clusterId, "", "", "", "", "")
		if err != nil {
			log.Printf("[ERROR] ingressInsert-V1beta1 Get List Fail:%s\n", err)
			return 0, err
		}
	}
	for _, vv := range xList {
		resName := vv.IngressName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,endpoint:%s,rules:%s", vv.Labels, vv.Endpoint, vv.Rules)
		_, _ = InsertIndex(clusterId, "ingress", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "ingress", fmt.Sprintf("%d", len(xList)))
	return len(xList), nil
}

func CmInsert(clusterId string) (int, error) {
	xList, err := CmList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] cmInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.ConfigmapName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s", vv.Labels)
		_, _ = InsertIndex(clusterId, "configmap", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "configmap", fmt.Sprintf("%d", len(xList)))

	//将全量列表写入缓存
	bodystr, err := json.Marshal(&xList)
	if err == nil {
		_ = common.SetEx("cmList"+clusterId, string(bodystr), 86400)
	}
	return len(xList), nil
}

func SecretInsert(clusterId string) (int, error) {
	xList, err := SecretList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] secretInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.SecretName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s", vv.Labels)
		_, _ = InsertIndex(clusterId, "secret", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "secret", fmt.Sprintf("%d", len(xList)))
	return len(xList), nil
}

func PvcInsert(clusterId string) (int, error) {
	xList, err := PersistentVolumeClaimList(clusterId, "", "", "", "")
	if err != nil {
		log.Printf("[ERROR] pvcInsert Get List Fail:%s\n", err)
		return 0, err
	}
	for _, vv := range xList {
		resName := vv.PvcName
		NameSpace := vv.NameSpace
		msg := fmt.Sprintf("labels:%s,capacity:%s,volumeName:%s", vv.Labels, vv.Capacity, vv.VolumeName)
		_, _ = InsertIndex(clusterId, "pvc", resName, NameSpace, msg)
	}
	_ = common.HSet("count_"+clusterId, "pvc", fmt.Sprintf("%d", len(xList)))
	return len(xList), nil
}
