package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type RolesController struct {
	beego.Controller
}

func (this *RolesController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	List, err := m.RolesList(clusterId, nameSpace)
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

func (this *RolesController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	rolesName := this.GetString("rolesName")
	yamlStr, _ := m.GetRolesYaml(clusterId, nameSpace, rolesName)
	this.Ctx.WriteString(yamlStr)
}
