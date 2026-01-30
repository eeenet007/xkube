package controllers

import (
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type TCPRouteController struct {
	beego.Controller
}

// List 获取TCPRoute列表
func (this *TCPRouteController) List() {
	clusterId := this.GetString("clusterId")

	list, err := m.TCPRouteList(clusterId)
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

// Detail 获取TCPRoute详情
func (this *TCPRouteController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	tcprouteName := this.GetString("tcprouteName")

	detail, err := m.GetTCPRouteDetail(clusterId, nameSpace, tcprouteName)

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

// Yaml 获取TCPRoute的YAML配置
func (this *TCPRouteController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	tcprouteName := this.GetString("tcprouteName")

	yamlStr, err := m.GetTCPRouteYaml(clusterId, nameSpace, tcprouteName)

	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Error: " + err.Error())
		return
	}

	this.Ctx.WriteString(yamlStr)
}

// Delete 删除TCPRoute
func (this *TCPRouteController) Delete() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	tcprouteName := this.GetString("tcprouteName")

	err := m.DeleteTCPRoute(clusterId, nameSpace, tcprouteName)

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

// Create 创建TCPRoute
func (this *TCPRouteController) Create() {
	clusterId := this.GetString("clusterId")

	err := m.CreateTCPRoute(clusterId, this.Ctx.Input.RequestBody)
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

// UpdateByYaml 通过YAML更新TCPRoute
func (this *TCPRouteController) UpdateByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	tcprouteName := this.GetString("tcprouteName")

	//backup
	yamlStr, _ := m.GetTCPRouteYaml(clusterId, nameSpace, tcprouteName)
	_ = m.InsertBackup(clusterId, nameSpace, tcprouteName, "tcproute", yamlStr, "Backup before updating")

	// 获取请求体中的YAML内容
	//yamlContent := string(this.Ctx.Input.RequestBody)

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.UpdateTCPRouteByYaml(clusterId, []byte(reqBody))

	if err != nil {
		log.Printf("[WARN] UpdateTCPRouteByYaml Fail:%s\n", err)
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
