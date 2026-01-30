package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type EventController struct {
	beego.Controller
}

func (this *EventController) List() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	kind := this.GetString("kind") //Pod/Namespace/Ingress/Service/Node/Configmap/Deployment/Stateful/Job/Cronjob/Secret可参考阿里云k8s事件中心类型
	objName := this.GetString("objName")
	limitd, _ := this.GetInt64("limit")
	List, err := m.EventListV2(clusterId, namespace, kind, objName, limitd)
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
