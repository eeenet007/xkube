package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type GatewayClassController struct {
	beego.Controller
}

func (this *GatewayClassController) List() {
	clusterId := this.GetString("clusterId")
	
	list, err := m.GatewayClassList(clusterId)
	msg := "success"
	code := 0
	count := len(list)
	
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	
	this.Data["json"] = &map[string]interface{}{
		"code":  code, 
		"msg":   msg, 
		"count": count, 
		"data":  &list,
	}
	this.ServeJSON()
}

func (this *GatewayClassController) Detail() {
	clusterId := this.GetString("clusterId")
	gatewayClassName := this.GetString("gatewayClassName")
	
	detail, err := m.GetGatewayClassDetail(clusterId, gatewayClassName)
	
	if err != nil {
		log.Println(err)
		this.Data["json"] = &map[string]interface{}{
			"code": -1,
			"msg":  err.Error(),
			"data": nil,
		}
		this.ServeJSON()
		return
	}
	
	this.Data["json"] = &map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": detail,
	}
	this.ServeJSON()
}

func (this *GatewayClassController) Yaml() {
	clusterId := this.GetString("clusterId")
	gatewayClassName := this.GetString("gatewayClassName")

	yamlStr, err := m.GetGatewayClassYaml(clusterId, gatewayClassName)
	
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Error: " + err.Error())
		return
	}
	
	this.Ctx.WriteString(yamlStr)
}
