package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type CdrController struct {
	beego.Controller
}

func (this *CdrController) List() {
	clusterId := this.GetString("clusterId")
	List, err := m.CdrList(clusterId)
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

func (this *CdrController) Yaml() {
	clusterId := this.GetString("clusterId")
	cdrName := this.GetString("cdrName")

	yamlStr, _ := m.GetCdrYaml(clusterId, cdrName)
	this.Ctx.WriteString(yamlStr)
}

func (this *CdrController) Del() {
	clusterId := this.GetString("clusterId")
	cdrName := this.GetString("cdrName")
	err := m.CdrDelete(clusterId, cdrName)
	code := 0
	msg := "success"
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
