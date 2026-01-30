// clusterModel.go
package models

import (
	//"bytes"
	//"errors"
	//"fmt"
	//"strings"

	//"io/ioutil"
	//"encoding/json"
	//"log"
	"time"

	//"xkube/common"

	//"github.com/beego/beego/v2/core/validation"
	//xm "xkube/xadmin/src/models"

	"github.com/beego/beego/v2/client/orm"
	//"github.com/beego/beego/v2/server/web/context"
)

type Xkb_clone_log struct { //查询时会将大写以下划线分开
	// Id clone log primary key
	// Id 克隆日志主键
	Id int64 `json:"id"`
	// Clusterid source cluster identifier
	// Clusterid 源集群标识符
	Clusterid string `json:"clusterid"`
	// Namespace source namespace
	// Namespace 源命名空间
	Namespace string `json:"namespace"`
	// Restype resource type
	// Restype 资源类型
	Restype string `json:"restype"`
	// Objname source object name
	// Objname 源对象名称
	Objname string `json:"objname"`
	// TargetClusterid target cluster identifier
	// TargetClusterid 目标集群标识符
	TargetClusterid string `json:"target_clusterid"`
	// TargetNamespace target namespace
	// TargetNamespace 目标命名空间
	TargetNamespace string `json:"target_namespace"`
	// TargetObjname target object name
	// TargetObjname 目标对象名称
	TargetObjname string `json:"target_objname"`
	// Status clone status
	// Status 克隆状态
	Status string `json:"status"`
	// Result clone result
	// Result 克隆结果
	Result string `json:"result"`
	// Remarks remarks or description
	// Remarks 备注或描述
	Remarks string `json:"remarks"`
	// User user who performed the clone operation
	// User 执行克隆操作的用户
	User string `json:"user"`
	// Createtime clone operation creation time
	// Createtime 克隆操作创建时间
	Createtime string `json:"createtime"`
}

func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Xkb_clone_log))
}

// GetList_Clone Get clone log list with pagination
// GetList_Clone 获取带分页的克隆日志列表
// page: page number (starting from 1)
// page: 页码（从1开始）
// page_size: number of items per page
// page_size: 每页条目数
// Returns:
//   - []Xkb_clone_log: clone log list
//   - int64: total count of clone logs
//
// 返回值:
//   - []Xkb_clone_log: 克隆日志列表
//   - int64: 克隆日志总数
func GetList_Clone(page, page_size int64) (dps []Xkb_clone_log, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_clone_log))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-createtime").All(&dps)
	count, _ = qs.Count()
	return dps, count
}

// Add_Clone Add a new clone log entry
// Add_Clone 添加新的克隆日志条目
// u: clone log information to be added
// u: 要添加的克隆日志信息
// Returns:
//   - int64: inserted record ID
//   - error: error message
//
// 返回值:
//   - int64: 插入记录的ID
//   - error: 错误信息
func Add_Clone(u *Xkb_clone_log) (int64, error) {
	o := orm.NewOrm()
	Ds := new(Xkb_clone_log)
	Ds = u
	Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
	num, err := o.Insert(Ds)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// DelById_Clone Delete a clone log entry by ID
// DelById_Clone 根据ID删除克隆日志条目
// id: clone log ID to be deleted
// id: 要删除的克隆日志ID
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func DelById_Clone(id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Xkb_clone_log{Id: id})
	return status, err
}
