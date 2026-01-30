// clusterModel.go
package models

import (
	//"bytes"
	"errors"
	"fmt"
	"strings"

	//"io/ioutil"
	//"encoding/json"
	"log"
	"time"

	"xkube/common"

	//"github.com/beego/beego/v2/core/validation"

	"github.com/beego/beego/v2/client/orm"
)

type Xkb_cluster struct { //查询时会将大写以下划线分开
	// Id cluster primary key
	// Id 集群主键
	Id int64 `json:"id"`
	// ClusterId cluster identifier
	// ClusterId 集群标识符
	ClusterId string `json:"cluster_id"`
	// ClusterName cluster name
	// ClusterName 集群名称
	ClusterName string `json:"cluster_name"`
	// IdcName IDC name where the cluster is located
	// IdcName 集群所在的IDC名称
	IdcName string `json:"idc_name"`
	// Status cluster status, 1 for active, 0 for inactive
	// Status 集群状态，1为活跃，0为非活跃
	Status int64 `json:"status"`
	// KubeVersion Kubernetes version of the cluster
	// KubeVersion 集群的Kubernetes版本
	KubeVersion string `json:"kube_version"`
	// KubeConfig Kubernetes configuration information
	// KubeConfig Kubernetes配置信息
	KubeConfig string `json:"kube_config"`
	// BearerToken authentication token for accessing the cluster
	// BearerToken 访问集群的认证令牌
	BearerToken string `json:"bearer_token"`
	// Remarks remarks or description of the cluster
	// Remarks 集群备注或描述
	Remarks string `json:"remarks"`
	// LanSlbip internal load balancer IP
	// LanSlbip 内网负载均衡IP
	LanSlbip string `json:"lan_slbip"`
	// WanSlbip external load balancer IP
	// WanSlbip 外网负载均衡IP
	WanSlbip string `json:"wan_slbip"`
	// Createtime cluster creation time
	// Createtime 集群创建时间
	Createtime string `json:"createtime"`
}

type ClusterCount struct {
	// ClusterId cluster identifier
	// ClusterId 集群标识符
	ClusterId string `json:"clusterId"`
	// Node number of nodes in the cluster
	// Node 集群中的节点数
	Node string `json:"node"`
	// Deploy number of deployments in the cluster
	// Deploy 集群中的部署数
	Deploy string `json:"deploy"`
	// Pod number of pods in the cluster
	// Pod 集群中的Pod数
	Pod string `json:"pod"`
	// Service number of services in the cluster
	// Service 集群中的服务数
	Service string `json:"service"`
	// Cronjob number of cron jobs in the cluster
	// Cronjob 集群中的定时任务数
	Cronjob string `json:"cronjob"`
	// Configmap number of configmaps in the cluster
	// Configmap 集群中的配置映射数
	Configmap string `json:"configmap"`
	// Secret number of secrets in the cluster
	// Secret 集群中的密钥数
	Secret string `json:"secret"`
}

func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Xkb_cluster))
}

// GetList_Cluster Get cluster list with pagination and filtering capabilities
// GetList_Cluster 获取带分页和过滤功能的集群列表
// id: cluster ID for filtering, empty means no filtering
// id: 用于过滤的集群ID，空表示不过滤
// clusterId: cluster identifier for filtering, empty means no filtering
// clusterId: 用于过滤的集群标识符，空表示不过滤
// page: page number (starting from 1)
// page: 页码（从1开始）
// page_size: number of items per page
// page_size: 每页条目数
// Returns:
//   - []Xkb_cluster: cluster list
//   - int64: total count of clusters matching the criteria
//
// 返回值:
//   - []Xkb_cluster: 集群列表
//   - int64: 符合条件的集群总数
func GetList_Cluster(id, clusterId string, page, page_size int64) (dps []Xkb_cluster, count int64) {
	o := orm.NewOrm()
	//qs := o.QueryTable(new(Xkb_cluster))
	qs := o.QueryTable(new(Xkb_cluster)) //不显示两个字段
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if id != "" {
		cond = cond.And("id", id)
	}

	if clusterId != "" {
		cond = cond.And("clusterId", clusterId)
	}

	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&dps, "id", "cluster_id", "cluster_name", "idc_name", "status", "kube_version", "remarks", "lan_slbip", "wan_slbip", "createtime")
	//qs.All(&xfr)
	count, _ = qs.Count()
	return dps, count
}

// GetDetail_Cluster Get detailed information of a cluster by ID
// GetDetail_Cluster 根据ID获取集群详细信息
// id: cluster ID
// id: 集群ID
// Returns:
//   - Xkb_cluster: cluster detailed information
//   - error: error message
//
// 返回值:
//   - Xkb_cluster: 集群详细信息
//   - error: 错误信息
func GetDetail_Cluster(id int64) (Xkb_cluster, error) {
	dps := Xkb_cluster{Id: id}
	o := orm.NewOrm()
	err := o.Read(&dps)
	return dps, err
}

// Add_Cluster Add a new cluster
// Add_Cluster 添加新集群
// u: cluster information to be added
// u: 要添加的集群信息
// Returns:
//   - int64: inserted record ID
//   - error: error message
//
// 返回值:
//   - int64: 插入记录的ID
//   - error: 错误信息
func Add_Cluster(u *Xkb_cluster) (int64, error) {
	o := orm.NewOrm()
	Ds := new(Xkb_cluster)
	Ds.ClusterId = u.ClusterId
	Ds.ClusterName = u.ClusterName
	Ds.IdcName = u.IdcName
	//Ds.Resource = u.Resource
	Ds.KubeVersion = u.KubeVersion
	Ds.KubeConfig = u.KubeConfig
	Ds.BearerToken = u.BearerToken
	Ds.LanSlbip = u.LanSlbip
	Ds.WanSlbip = u.WanSlbip
	Ds.Status = u.Status
	Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
	Ds.Remarks = u.Remarks
	num, err := o.Insert(Ds)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// DelById_Cluster Delete a cluster by ID
// DelById_Cluster 根据ID删除集群
// id: cluster ID to be deleted
// id: 要删除的集群ID
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func DelById_Cluster(id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Xkb_cluster{Id: id})
	return status, err
}

// Edit_Cluster Edit cluster field value by ID
// Edit_Cluster 根据ID编辑集群字段值
// id: cluster ID
// id: 集群ID
// key: field name to be updated
// key: 要更新的字段名
// value: new value for the field
// value: 字段的新值
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func Edit_Cluster(id, key, value string) (int64, error) {
	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("UPDATE xkb_cluster SET %s = '%s' WHERE id = %s", key, value, id)
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("[INFO] EditCluster affected nums:%d\n", num)
		return num, nil
	}
	return 0, err
}

// Update_Cluster Update cluster information
// Update_Cluster 更新集群信息
// u: cluster information to be updated
// u: 要更新的集群信息
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func Update_Cluster(u *Xkb_cluster) (int64, error) {
	o := orm.NewOrm()
	op := make(orm.Params)

	op["ClusterId"] = u.ClusterId
	op["ClusterName"] = u.ClusterName
	op["IdcName"] = u.IdcName
	op["KubeVersion"] = u.KubeVersion
	op["KubeConfig"] = u.KubeConfig
	op["BearerToken"] = u.BearerToken
	op["LanSlbip"] = u.LanSlbip
	op["WanSlbip"] = u.WanSlbip
	op["Remarks"] = u.Remarks
	op["Status"] = u.Status
	//op["Updatetime"] = time.Now().Format("2006-01-02 15:04:05")

	if len(op) == 0 {
		return 0, errors.New("Xkb_cluster update field is empty")
	}
	var table Xkb_cluster
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(op)
	//if err != nil {
	delete(common.ClusterMap, u.ClusterId)
	delete(common.ClusterTokenMap, u.ClusterId)
	delete(common.ClusterVersionMap, u.ClusterId)
	//}
	return num, err
}

// CountWorkLoad Get workload statistics for all clusters
// CountWorkLoad 获取所有集群的工作负载统计信息
// Returns:
//   - []ClusterCount: workload statistics for all clusters
//   - error: error message
//
// 返回值:
//   - []ClusterCount: 所有集群的工作负载统计信息
//   - error: 错误信息
func CountWorkLoad() ([]ClusterCount, error) {
	//clusterArry := common.HKeys("count_*")
	clusterArry := common.Keys("count_*")
	var countArry = make([]ClusterCount, 0)
	for _, vv := range clusterArry {
		vmap := map[string]string{
			"node":      "0",
			"deploy":    "0",
			"pod":       "0",
			"service":   "0",
			"cronjob":   "0",
			"configmap": "0",
			"secret":    "0",
		}
		resp := common.HGetAll(vv)
		for _, v1 := range resp {
			vmap[v1.Key] = v1.Value
		}

		// err := json.Unmarshal([]byte(fmt.Sprintf("%v", resp)), &cc)
		// if err != nil {
		// 	log.Printf("[ERROR] Unmarshal %s fail:%s\n", vv, err)
		// 	continue
		// }
		//fmt.Println(&cc)
		//log.Println(vv, strings.TrimLeft(vv, "count_"))
		//fmt.Println(strings.Replace("count_ops-cluster", "count_", "", 1))
		countArry = append(countArry, ClusterCount{
			ClusterId: strings.Replace(vv, "count_", "", 1),
			Node:      vmap["node"],
			Deploy:    vmap["deploy"],
			Pod:       vmap["pod"],
			Service:   vmap["service"],
			Cronjob:   vmap["cronjob"],
			Configmap: vmap["configmap"],
			Secret:    vmap["secret"],
		})
	}
	return countArry, nil
}
