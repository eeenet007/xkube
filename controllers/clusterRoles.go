package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type ClusterRolesController struct {
	beego.Controller
}

func (this *ClusterRolesController) List() {
	clusterId := this.GetString("clusterId")
	List, err := m.ClusterRolesList(clusterId)
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

func (this *ClusterRolesController) Yaml() {
	clusterId := this.GetString("clusterId")
	crName := this.GetString("crName")
	yamlStr, _ := m.GetClusterRolesYaml(clusterId, crName)
	this.Ctx.WriteString(yamlStr)
}
