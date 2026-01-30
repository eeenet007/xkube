// metrics.go
package controllers

import (
	//"encoding/json"
	//"log"
	//"strings"
	//"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type MetricsController struct {
	beego.Controller
}

func (this *MetricsController) PodList() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	sort := this.GetString("sort") //mem or cpu

	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		sort = gp.Get("sort").String()
		nameSpace = gp.Get("nameSpace").String()
	}

	xList, err := m.GetPodMetricList(clusterId, nameSpace, sort)
	msg := "success"
	code := 0
	if err != nil {
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": len(xList), "data": &xList}
	this.ServeJSON()
}

func (this *MetricsController) NodeList() {
	clusterId := this.GetString("clusterId")
	//nameSpace := this.GetString("nameSpace")
	xList, err := m.GetNodeMetricList(clusterId)
	msg := "success"
	code := 0
	if err != nil {
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": len(xList), "data": &xList}
	this.ServeJSON()
}

func (this *MetricsController) PodUsage() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	podName := this.GetString("podName")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	metricStr := m.GetPodMetric(clusterId, nameSpace, podName)
	this.Ctx.WriteString(metricStr)
}

func (this *MetricsController) NodeUsage() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	metricStr := m.GetNodeMetric(clusterId, nodeName)
	this.Ctx.WriteString(metricStr)
}

// func (this *MetricsController) CpuUsageList() {
// 	clusterId := this.GetString("clusterId")
// 	nameSpace := this.GetString("nameSpace")
// 	podName := this.GetString("podName")
// 	xList := m.GetPodCpuMetric(clusterId, nameSpace, podName)
// 	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": len(xList), "data": &xList}
// 	this.ServeJSON()
// }
