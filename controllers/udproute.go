package controllers

import (
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type UDPRouteController struct {
	beego.Controller
}

// List 获取UDPRoute列表
func (this *UDPRouteController) List() {
	clusterId := this.GetString("clusterId")

	list, err := m.UDPRouteList(clusterId)
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

// Detail 获取UDPRoute详情
func (this *UDPRouteController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	udprouteName := this.GetString("udprouteName")

	detail, err := m.GetUDPRouteDetail(clusterId, nameSpace, udprouteName)

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

// Yaml 获取UDPRoute的YAML配置
func (this *UDPRouteController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	udprouteName := this.GetString("udprouteName")

	yamlStr, err := m.GetUDPRouteYaml(clusterId, nameSpace, udprouteName)

	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Error: " + err.Error())
		return
	}

	this.Ctx.WriteString(yamlStr)
}

// Delete 删除UDPRoute
func (this *UDPRouteController) Delete() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	udprouteName := this.GetString("udprouteName")

	err := m.DeleteUDPRoute(clusterId, nameSpace, udprouteName)

	if err != nil {
		log.Println(err)
		this.Data["json"] = &map[string]interface{}{
			"code": -1,
			"msg":  err.Error(),
		}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &map[string]interface{}{
		"code": 0,
		"msg":  "success",
	}
	this.ServeJSON()
}

// Create 创建UDPRoute
func (this *UDPRouteController) Create() {
	clusterId := this.GetString("clusterId")

	err := m.CreateUDPRoute(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		log.Println(err)
		this.Data["json"] = &map[string]interface{}{
			"code": -1,
			"msg":  err.Error(),
		}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &map[string]interface{}{
		"code": 0,
		"msg":  "success",
	}
	this.ServeJSON()
}

// UpdateByYaml 通过YAML更新UDPRoute
func (this *UDPRouteController) UpdateByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	udprouteName := this.GetString("udprouteName")

	//backup
	yamlStr, _ := m.GetUDPRouteYaml(clusterId, nameSpace, udprouteName)
	_ = m.InsertBackup(clusterId, nameSpace, udprouteName, "udproute", yamlStr, "Backup before updating")

	// 获取请求体中的YAML内容
	//yamlContent := string(this.Ctx.Input.RequestBody)

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.UpdateUDPRouteByYaml(clusterId, []byte(reqBody))

	if err != nil {
		log.Printf("[WARN] UpdateUDPRouteByYaml Fail:%s\n", err)
		this.Data["json"] = &map[string]interface{}{
			"code": -1,
			"msg":  err.Error(),
		}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &map[string]interface{}{
		"code": 0,
		"msg":  "success",
	}
	this.ServeJSON()
}
