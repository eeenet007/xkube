// backup_model.go
package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	//beego "github.com/beego/beego/v2/server/web"
)

type Xkb_backup struct {
	Id         int64  `json:"id"`
	ClusterId  string `json:"clusterId"`
	ResType    string `json:"resType"`
	ResName    string `json:"resName"`
	NameSpace  string `json:"nameSpace"`
	Content    string `json:"content"`
	Remarks    string `json:"remarks"`
	Createtime string `json:"createtime"`
}

func init() {
	orm.RegisterModel(new(Xkb_backup))
}

func InsertBackup(clusterId, nameSpace, resName, resType, content, remarks string) error {
	o := orm.NewOrm()
	//createtime := time.Now().Format("2006-01-02 15:04:05")
	content = strings.ReplaceAll(content, "'", "\\'")
	sqlstr := fmt.Sprintf("INSERT INTO xkb_backup(cluster_id, name_space,res_name,res_type,content,remarks)VALUE('%s','%s','%s','%s','%s','%s')", clusterId, nameSpace, resName, resType, content, remarks)
	res, err := o.Raw(sqlstr).Exec()
	if err != nil {
		return err
	}
	num, _ := res.RowsAffected()
	log.Printf("[INFO] InsertBackup result:%d\n", num)
	return nil
}

// GetList_Backup 获取备份列表 / Get backup list
// login_user: 登录用户 / Login user
// starttime: 开始时间 / Start time
// endtime: 结束时间 / End time
// id: 日志ID / Log ID
// page: 页码 / Page number
// page_size: 每页数量 / Page size
// []Xkb_backup: 返回备份列表 / Return backup list
// count: 返回总数量 / Return total count
func GetList_Backup(clusterId, nameSpace, resName, resType, starttime, endtime string, id, page, page_size int64) (xps []Xkb_backup, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_backup))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if id > 0 {
		cond = cond.And("id", id)
	}
	if clusterId != "" {
		cond = cond.And("cluster_id", clusterId)
	}
	if nameSpace != "" {
		cond = cond.And("name_space", nameSpace)
	}
	if resType != "" {
		cond = cond.And("res_type", resType)
	}
	if resName != "" {
		cond = cond.And("res_name", resName)
	}
	if starttime != "" {
		cond = cond.And("createtime__gte", starttime)
	}

	if endtime != "" {
		cond = cond.And("createtime__lte", endtime)
	}

	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&xps)
	//qs.All(&xfr)
	count, _ = qs.Count()
	return xps, count
}

func GetBackupById(id int64) (bak Xkb_backup, err error) {
	o := orm.NewOrm()
	table := new(Xkb_backup)
	err = o.QueryTable(table).Filter("Id", id).One(&bak)
	if err != nil {
		return bak, err
	}
	return bak, nil
}
