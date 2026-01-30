package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type PvController struct {
	beego.Controller
}

func (this *PvController) List() {
	clusterId := this.GetString("clusterId")
	List, err := m.PersistentVolumeList(clusterId)
	msg := "success"
	code := 0
	count := len(List)
	if err != nil {
		//log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &List}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *PvController) Detail() {
	clusterId := this.GetString("clusterId")
	pvName := this.GetString("pvName")
	xdetail, err := m.PersistentVolumeDetail(clusterId, pvName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *PvController) Yaml() {
	clusterId := this.GetString("clusterId")
	pvName := this.GetString("pvName")
	yamlStr, _ := m.GetPersistentVolumeYaml(clusterId, pvName)
	this.Ctx.WriteString(yamlStr)
}
