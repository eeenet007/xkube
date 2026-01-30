// Package models provides data models and related functionality for the xkube admin system.
// 包models为xkube管理系统提供数据模型及相关功能。
package models

import (
	//"errors"
	"log"
	"time"

	//. "xkube/xadmin/src/lib"
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/tidwall/gjson"
	//"github.com/beego/beego/v2/core/validation"
	//beego "github.com/beego/beego/v2/server/web"
)

// UserCluster represents the relationship between a user and a cluster.
// UserCluster 表示用户与集群之间的关系。
type UserCluster struct {
	Id         int64  `json:"id"`         // Id is the unique identifier of the user-cluster relationship. Id是用户-集群关系的唯一标识符。
	Username   string `json:"username"`   // Username is the name of the user. Username是用户的名称。
	ClusterId  string `json:"clusterId"`  // ClusterId is the identifier of the cluster. ClusterId是集群的标识符。
	Createtime string `json:"createtime"` // Createtime is the creation time of the relationship. Createtime是关系的创建时间。
}

// MyCluster represents cluster information for a user.
// MyCluster 表示用户的集群信息。
type MyCluster struct {
	Username    string `json:"username"`     // Username is the name of the user. Username是用户的名称。
	ClusterId   string `json:"cluster_id"`   // ClusterId is the identifier of the cluster. ClusterId是集群的标识符。
	ClusterName string `json:"cluster_name"` // ClusterName is the name of the cluster. ClusterName是集群的名称。
	KubeVersion string `json:"kube_version"` // KubeVersion is the Kubernetes version of the cluster. KubeVersion是集群的Kubernetes版本。
}

// init registers the UserCluster model with the ORM framework.
// init 将UserCluster模型注册到ORM框架中。
func init() {
	//orm.Debug = true
	orm.RegisterModel(new(UserCluster))
}

// GetUserClusterList retrieves a paginated list of user-cluster relationships based on filters.
// GetUserClusterList 根据过滤条件检索用户-集群关系的分页列表。
// username: filter by username, 筛选用户名
// clusterId: filter by cluster ID, 筛选集群ID
// page: page number, 页码
// page_size: number of items per page, 每页项目数
// Returns: list of UserCluster and total count, 返回: UserCluster列表和总数
func GetUserClusterList(username, clusterId string, page int64, page_size int64) (list []UserCluster, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserCluster))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if username != "" {
		cond = cond.And("username", username)
	}

	if clusterId != "" {
		cond = cond.And("cluster_id", clusterId)
	}

	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&list)
	count, _ = qs.Count()
	return list, count

}

// GetUserClusterListV2 retrieves a list of user-cluster relationships with more flexible filtering options using a map.
// GetUserClusterListV2 使用映射提供更灵活的筛选选项来检索用户-集群关系列表。
// whereMap: map containing filter conditions, whereMap: 包含筛选条件的映射
// Returns: list of UserCluster and total count, 返回: UserCluster列表和总数
func GetUserClusterListV2(whereMap map[string]gjson.Result) (list []UserCluster, count int64) {
	o := orm.NewOrm()
	cond := orm.NewCondition()

	slb := new(UserCluster)
	qs := o.QueryTable(slb)

	var page, limit int64
	if pg, ok := whereMap["page"]; ok {
		page = pg.Int()
		delete(whereMap, "page")
	}
	if pg, ok := whereMap["limit"]; ok {
		limit = pg.Int()
		delete(whereMap, "limit")
	}
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}

	for kk, vv := range whereMap {
		if vv.String() != "" {
			cond = cond.And(kk+"__contains", vv.String())
		}
	}
	qs = qs.SetCond(cond)
	qs.Limit(limit, offset).OrderBy("-id").All(&list)

	count, _ = qs.Count()
	return list, count
}

// GetMyClusterList retrieves the list of clusters that a user has access to.
// GetMyClusterList 检索用户有权访问的集群列表。
// username: the username to query clusters for, username: 要查询集群的用户名
// Returns: slice of MyCluster pointers and an error if any, 返回: MyCluster指针切片和可能的错误
func GetMyClusterList(username string) ([]*MyCluster, error) {
	sqlstr := fmt.Sprintf("SELECT uc.username,xc.cluster_id,xc.cluster_name,xc.kube_version FROM user_cluster uc JOIN xkb_cluster xc ON uc.cluster_id = xc.cluster_id WHERE uc.username ='%s'", username)

	if username == "admin" {
		sqlstr = "SELECT xc.cluster_id,xc.cluster_name,xc.kube_version FROM xkb_cluster xc"
	} else {
		_, num := GetUserClusterList(username, "all", 1, 20)
		if num > 0 {
			sqlstr = "SELECT xc.cluster_id,xc.cluster_name,xc.kube_version FROM xkb_cluster xc"
		}
	}
	//log.Println(sqlstr)
	o := orm.NewOrm()
	var stus []*MyCluster
	_, err := o.Raw(sqlstr).QueryRows(&stus)
	if err != nil {
		log.Println(err)
		return stus, err
	}
	//log.Println(stus)
	return stus, nil
}

// AddUserCluster adds a new user-cluster relationship.
// AddUserCluster 添加新的用户-集群关系。
// u: pointer to UserCluster object to be added, u: 指向要添加的UserCluster对象的指针
// Returns: inserted ID and error if any, 返回: 插入的ID和可能的错误
func AddUserCluster(u *UserCluster) (int64, error) {
	o := orm.NewOrm()
	uc := new(UserCluster)
	uc.Username = u.Username
	uc.ClusterId = u.ClusterId
	uc.Createtime = time.Now().Format("2006-01-02 15:04:05")

	_, num := GetUserClusterList(u.Username, "all", 1, 20)
	if num > 0 {
		return 0, fmt.Errorf("already have permission")
	}
	id, err := o.Insert(uc)
	return id, err
}

// DelUserClusterById deletes a user-cluster relationship by its ID.
// DelUserClusterById 根据ID删除用户-集群关系。
// Id: the ID of the user-cluster relationship to delete, Id: 要删除的用户-集群关系的ID
// Returns: status code and error if any, 返回: 状态码和可能的错误
func DelUserClusterById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&UserCluster{Id: Id})
	return status, err
}
