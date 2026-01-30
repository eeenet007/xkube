package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type ServiceAccountsController struct {
	beego.Controller
}

func (this *ServiceAccountsController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	List, err := m.ServiceAccountsList(clusterId, nameSpace)
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

func (this *ServiceAccountsController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	saName := this.GetString("saName")
	yamlStr, _ := m.GetServiceAccountsYaml(clusterId, nameSpace, saName)
	this.Ctx.WriteString(yamlStr)
}
