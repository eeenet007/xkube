package controllers

import (
	//"encoding/json"
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	//gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type GRPCRouteController struct {
	beego.Controller
}

// List 获取GRPCRoute列表
func (this *GRPCRouteController) List() {
	clusterId := this.GetString("clusterId")

	list, err := m.GRPCRouteList(clusterId)
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

// Detail 获取GRPCRoute详情
func (this *GRPCRouteController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	grpcrouteName := this.GetString("grpcrouteName")

	detail, err := m.GetGRPCRouteDetail(clusterId, nameSpace, grpcrouteName)

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

// Yaml 获取GRPCRoute的YAML配置
func (this *GRPCRouteController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	grpcrouteName := this.GetString("grpcrouteName")

	yamlStr, err := m.GetGRPCRouteYaml(clusterId, nameSpace, grpcrouteName)

	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Error: " + err.Error())
		return
	}

	this.Ctx.WriteString(yamlStr)
}

// Delete 删除GRPCRoute
func (this *GRPCRouteController) Delete() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	grpcrouteName := this.GetString("grpcrouteName")

	err := m.DeleteGRPCRoute(clusterId, nameSpace, grpcrouteName)

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

// Create 创建GRPCRoute
func (this *GRPCRouteController) Create() {
	clusterId := this.GetString("clusterId")

	err := m.CreateGRPCRoute(clusterId, this.Ctx.Input.RequestBody)

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

// UpdateByYaml 通过YAML更新GRPCRoute
func (this *GRPCRouteController) UpdateByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	grpcrouteName := this.GetString("grpcrouteName")

	//backup
	yamlStr, _ := m.GetGRPCRouteYaml(clusterId, nameSpace, grpcrouteName)
	_ = m.InsertBackup(clusterId, nameSpace, grpcrouteName, "grpcroute", yamlStr, "Backup before updating")

	// 获取请求体中的YAML内容
	//yamlContent := string(this.Ctx.Input.RequestBody)

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.UpdateGRPCRouteByYaml(clusterId, []byte(reqBody))

	if err != nil {
		log.Printf("[WARN] UpdateGRPCRouteByYaml Fail:%s\n", err)
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
