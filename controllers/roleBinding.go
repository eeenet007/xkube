package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type RoleBindingController struct {
	beego.Controller
}

func (this *RoleBindingController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	List, err := m.RoleBindingList(clusterId, nameSpace)
	msg := "success"
	code := 0
	count := len(List)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &List}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *RoleBindingController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	rbName := this.GetString("rbName")
	yamlStr, _ := m.GetRoleBindingYaml(clusterId, nameSpace, rbName)
	this.Ctx.WriteString(yamlStr)
}
