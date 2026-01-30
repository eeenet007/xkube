package controllers

import (
	//"encoding/json"
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type GatewayController struct {
	beego.Controller
}

// List 获取Gateway列表
func (this *GatewayController) List() {
	clusterId := this.GetString("clusterId")

	list, err := m.GatewayList(clusterId)
	msg := "success"
	code := 0
	count := len(list)

	if err != nil {
		//log.Println(err)
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

// Detail 获取Gateway详情
func (this *GatewayController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	gatewayName := this.GetString("gatewayName")

	detail, err := m.GetGatewayDetail(clusterId, nameSpace, gatewayName)

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

// Yaml 获取Gateway的YAML配置
func (this *GatewayController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	gatewayName := this.GetString("gatewayName")

	yamlStr, err := m.GetGatewayYaml(clusterId, nameSpace, gatewayName)

	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Error: " + err.Error())
		return
	}

	this.Ctx.WriteString(yamlStr)
}

// Delete 删除Gateway
func (this *GatewayController) Delete() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	gatewayName := this.GetString("gatewayName")

	err := m.DeleteGateway(clusterId, nameSpace, gatewayName)

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

// Create 创建Gateway
func (this *GatewayController) Create() {
	clusterId := this.GetString("clusterId")
	err := m.CreateGateway(clusterId, this.Ctx.Input.RequestBody)
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

// UpdateByYaml 通过YAML更新Gateway
func (this *GatewayController) UpdateByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	gatewayName := this.GetString("gatewayName")

	//backup
	yamlStr, _ := m.GetGatewayYaml(clusterId, nameSpace, gatewayName)
	_ = m.InsertBackup(clusterId, nameSpace, gatewayName, "gateway", yamlStr, "Backup before updating")

	// 获取请求体中的YAML内容
	//yamlContent := string(this.Ctx.Input.RequestBody)

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.UpdateGatewayByYaml(clusterId, []byte(reqBody))

	if err != nil {
		log.Printf("[WARN] UpdateGatewayByYaml Fail:%s\n", err)
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
