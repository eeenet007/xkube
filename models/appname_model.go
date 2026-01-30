// appname_model.go
package models

import (
	//"bytes"
	"errors"
	"fmt"

	//"io/ioutil"
	"log"
	"time"

	//"xkube/common"

	//"github.com/beego/beego/v2/core/validation"

	"github.com/beego/beego/v2/client/orm"
)

type XkbAppname struct {
	// Id application name primary key
	// Id 应用名称主键
	Id int64 `json:"id"`
	// Appname application name
	// Appname 应用名称
	Appname string `json:"appname"`
	// Remarks remarks or description of the application
	// Remarks 应用备注或描述
	Remarks string `json:"remarks"`
	// Createtime application name creation time
	// Createtime 应用名称创建时间
	Createtime string `json:"createtime"`
}

func init() {
	//orm.Debug = true
	orm.RegisterModel(new(XkbAppname))
}

// GetList_Appname Get application name list with pagination and filtering capabilities
// GetList_Appname 获取带分页和过滤功能的应用名称列表
// appname: application name for filtering, empty means no filtering
// appname: 用于过滤的应用名称，空表示不过滤
// page: page number (starting from 1)
// page: 页码（从1开始）
// page_size: number of items per page
// page_size: 每页条目数
// Returns:
//   - []XkbAppname: application name list
//   - int64: total count of application names matching the criteria
//
// 返回值:
//   - []XkbAppname: 应用名称列表
//   - int64: 符合条件的应用名称总数
func GetList_Appname(appname string, page, page_size int64) (dps []XkbAppname, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(XkbAppname))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if appname != "" {
		//cond = cond.And("appname", appname)
		cond = cond.And("appname__contains", appname)
	}

	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&dps)
	count, _ = qs.Count()
	return dps, count
}

// Add_Appname Add a new application name
// Add_Appname 添加新应用名称
// u: application name information to be added
// u: 要添加的应用名称信息
// Returns:
//   - int64: inserted record ID
//   - error: error message
//
// 返回值:
//   - int64: 插入记录的ID
//   - error: 错误信息
func Add_Appname(u *XkbAppname) (int64, error) {
	o := orm.NewOrm()
	Ds := new(XkbAppname)
	Ds.Appname = u.Appname
	Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
	Ds.Remarks = u.Remarks
	num, err := o.Insert(Ds)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// DelById_Appname Delete an application name by ID
// DelById_Appname 根据ID删除应用名称
// id: application name ID to be deleted
// id: 要删除的应用名称ID
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func DelById_Appname(id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&XkbAppname{Id: id})
	return status, err
}

// Edit_Appname Edit application name field value by ID
// Edit_Appname 根据ID编辑应用名称字段值
// id: application name ID
// id: 应用名称ID
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
func Edit_Appname(id, key, value string) (int64, error) {
	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("UPDATE xkb_appname SET %s = '%s' WHERE id = %s", key, value, id)
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("[INFO] EditAppname affected nums:%d\n", num)
		return num, nil
	}
	return 0, err
}

// Update_Appname Update application name information
// Update_Appname 更新应用名称信息
// u: application name information to be updated
// u: 要更新的应用名称信息
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func Update_Appname(u *XkbAppname) (int64, error) {
	o := orm.NewOrm()
	op := make(orm.Params)

	op["Appname"] = u.Appname
	op["Remarks"] = u.Remarks

	if len(op) == 0 {
		return 0, errors.New("appname update field is empty")
	}
	var table XkbAppname
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(op)
	return num, err
}
