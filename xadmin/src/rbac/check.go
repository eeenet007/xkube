package rbac

import (
	"log"
	"xkube/common"

	"github.com/beego/beego/v2/client/orm"
)

type T1 struct{}

//service check
func (this *MainController) Check() {
	o := orm.NewOrm()
	sqlstr := "select 1"
	var t1 T1
	err := o.Raw(sqlstr).QueryRow(&t1)
	if err != nil {
		log.Printf("[ERROR] healthy check error:%s\n", err)
		this.Ctx.Output.SetStatus(500)
		this.Rsp(false, "mysqlCheckFail")
	}
	err1 := common.Set("check", "ok")
	if err1 != nil {
		this.Ctx.Output.SetStatus(500)
		this.Rsp(false, "redisSetCheckFail")
	}
	vv := common.Get("check")
	if vv != "ok" {
		this.Ctx.Output.SetStatus(500)
		this.Rsp(false, "redisGetCheckFail")
	}
	this.Rsp(true, "ok")
}
