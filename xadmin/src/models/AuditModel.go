// AuditModel.go
package models

import (
	"fmt"
	"log"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// Xkb_audit 审计日志结构体 / Audit log structure
type Xkb_audit struct {
	Id         int64  `json:"id"`         // Id 日志ID / Log ID
	Login_user string `json:"login_user"` // Login_user 登录用户 / Login user
	User_ip    string `json:"user_ip"`    // User_ip 用户IP地址 / User IP address
	Action     string `json:"action"`     // Action 操作行为 / Action performed
	Status     string `json:"status"`     // Status 操作状态 / Action status
	Message    string `json:"message"`    // Message 操作信息 / Action message
	Createtime string `json:"createtime"` // Createtime 创建时间 / Creation time
}

func init() {
	orm.RegisterModel(new(Xkb_audit))
}

// InsertLogAudit 插入日志到审计表 / Insert log to audit table
// user_ip: 用户IP地址 / User IP address
// login_user: 登录用户 / Login user
// action: 操作行为 / Action performed
// status: 操作状态 / Action status
// message: 操作信息 / Action message
// error: 返回错误信息 / Return error message
func InsertLogAudit(user_ip, login_user, action, status, message string) error {
	o := orm.NewOrm()
	//createtime := time.Now().Format("2006-01-02 15:04:05")
	sqlstr := fmt.Sprintf("INSERT INTO xkb_audit(user_ip, login_user,action,status,message)VALUE('%s','%s','%s','%s','%s')", user_ip, login_user, action, status, message)
	//主配置 / Main configuration
	res, err := o.Raw(sqlstr).Exec()
	if err != nil {
		return err
	}
	num, _ := res.RowsAffected()
	log.Printf("[INFO] InsertLogAudit result:%d\n", num)
	return nil
}

func InsertLogAuditV2(this *beego.Controller, action, status, message string) error {
	xff := this.Ctx.Request.Header.Get("x-forwarded-for")
	remoteAddr := this.Ctx.Request.RemoteAddr
	user_ip := fmt.Sprintf("%s,%s", xff, remoteAddr)

	userinfo := this.GetSession("userinfo").(map[string]interface{})
	login_user := userinfo["username"].(string)

	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("INSERT INTO xkb_audit(user_ip, login_user,action,status,message)VALUE('%s','%s','%s','%s','%s')", user_ip, login_user, action, status, message)
	res, err := o.Raw(sqlstr).Exec()
	if err != nil {
		return err
	}
	num, err2 := res.RowsAffected()
	if err2 != nil {
		log.Printf("[INFO] InsertLogAudit result:%d,%v\n", num, err2)
	}
	return err2
}

// GetList_LogAudit 获取审计日志列表 / Get audit log list
// login_user: 登录用户 / Login user
// starttime: 开始时间 / Start time
// endtime: 结束时间 / End time
// id: 日志ID / Log ID
// page: 页码 / Page number
// page_size: 每页数量 / Page size
// []Xkb_audit: 返回审计日志列表 / Return audit log list
// count: 返回总数量 / Return total count
func GetList_LogAudit(login_user, starttime, endtime string, id, page, page_size int64) (xps []Xkb_audit, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_audit))
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
	if login_user != "" {
		cond = cond.And("login_user", login_user)
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
