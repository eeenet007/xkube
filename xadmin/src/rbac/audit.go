// audit.go
package rbac

import (
	//"fmt"
	"github.com/tidwall/gjson"

	m "xkube/xadmin/src/models"
)

type AuditController struct {
	CommonController
}

func (ac *AuditController) List() {
	login_user := ac.GetString("login_user")
	starttime := ac.GetString("starttime")
	endtime := ac.GetString("endtime")
	id, _ := ac.GetInt64("id")
	page, _ := ac.GetInt64("page")
	page_size, _ := ac.GetInt64("limit")
	if ac.Ctx.Request.Method == "POST" {
		gp := gjson.ParseBytes(ac.Ctx.Input.RequestBody)
		login_user = gp.Get("login_user").String()
		starttime = gp.Get("starttime").String()
		endtime = gp.Get("endtime").String()
		page = gp.Get("page").Int()
		page_size = gp.Get("limit").Int()
	}
	datas, count := m.GetList_LogAudit(login_user, starttime, endtime, id, page, page_size)
	ac.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	ac.ServeJSON()
}
