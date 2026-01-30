package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type ClusterRoleBindingController struct {
	beego.Controller
}

func (this *ClusterRoleBindingController) List() {
	clusterId := this.GetString("clusterId")
	List, err := m.ClusterRoleBindingList(clusterId)
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

func (this *ClusterRoleBindingController) Yaml() {
	clusterId := this.GetString("clusterId")
	crbName := this.GetString("crbName")
	yamlStr, _ := m.GetClusterRoleBindingYaml(clusterId, crbName)
	this.Ctx.WriteString(yamlStr)
}
