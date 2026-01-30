package controllers

import (
	//"encoding/json"
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	//gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type HTTPRouteController struct {
	beego.Controller
}

// List 获取HTTPRoute列表
func (this *HTTPRouteController) List() {
	clusterId := this.GetString("clusterId")

	list, err := m.HTTPRouteList(clusterId)
	msg := "success"
	code := 0
	count := len(list)
	if err != nil {
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

// Detail 获取HTTPRoute详情
func (this *HTTPRouteController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	httprouteName := this.GetString("httprouteName")

	detail, err := m.GetHTTPRouteDetail(clusterId, nameSpace, httprouteName)

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

// Yaml 获取HTTPRoute的YAML配置
func (this *HTTPRouteController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	httprouteName := this.GetString("httprouteName")

	yamlStr, err := m.GetHTTPRouteYaml(clusterId, nameSpace, httprouteName)

	if err != nil {
		log.Println(err)
		this.Ctx.WriteString("Error: " + err.Error())
		return
	}

	this.Ctx.WriteString(yamlStr)
}

// Delete 删除HTTPRoute
func (this *HTTPRouteController) Delete() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	httprouteName := this.GetString("httprouteName")

	err := m.DeleteHTTPRoute(clusterId, nameSpace, httprouteName)

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

// Create 创建HTTPRoute
func (this *HTTPRouteController) Create() {
	clusterId := this.GetString("clusterId")
	err := m.CreateHTTPRoute(clusterId, this.Ctx.Input.RequestBody)

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

// UpdateByYaml 通过YAML更新HTTPRoute
func (this *HTTPRouteController) UpdateByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	httprouteName := this.GetString("httprouteName")
	//backup
	yamlStr, _ := m.GetHTTPRouteYaml(clusterId, nameSpace, httprouteName)
	_ = m.InsertBackup(clusterId, nameSpace, httprouteName, "httproute", yamlStr, "Backup before updating")

	// 获取请求体中的YAML内容
	//yamlContent := string(this.Ctx.Input.RequestBody)

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.UpdateHTTPRouteByYaml(clusterId, []byte(reqBody))

	if err != nil {
		log.Printf("[WARN] UpdateHTTPRouteByYaml Fail:%s\n", err)
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
