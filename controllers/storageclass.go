package controllers

import (
	//"encoding/json"
	"log"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type StorageClassController struct {
	beego.Controller
}

func (this *StorageClassController) List() {
	clusterId := this.GetString("clusterId")
	List, err := m.StorageclassList(clusterId)
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

func (this *StorageClassController) Detail() {
	clusterId := this.GetString("clusterId")
	storageclassName := this.GetString("storageclassName")
	xdetail, err := m.StorageclassDetail(clusterId, storageclassName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *StorageClassController) Yaml() {
	clusterId := this.GetString("clusterId")
	storageclassName := this.GetString("storageclassName")

	yamlStr, _ := m.GetStorageclassYaml(clusterId, storageclassName)
	this.Ctx.WriteString(yamlStr)
}
