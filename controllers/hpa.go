package controllers

import (
	//"encoding/json"
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type HpaController struct {
	beego.Controller
}

func (this *HpaController) List() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")
	targetName := this.GetString("targetName")
	List, err := m.HpaList(clusterId, namespace, hpaName, targetName)
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

func (this *HpaController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")

	yamlStr, _ := m.GetHpaYaml(clusterId, namespace, hpaName)
	this.Ctx.WriteString(yamlStr)
}

func (this *HpaController) Create() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	err := m.HpaCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] hpa Create Fail:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *HpaController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")
	err := m.HpaDelete(clusterId, namespace, hpaName)
	code := 0
	msg := "success"
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *HpaController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")
	//backup
	yamlStr, _ := m.GetHpaYaml(clusterId, nameSpace, hpaName)
	_ = m.InsertBackup(clusterId, nameSpace, hpaName, "hpa", yamlStr, "Backup before updating")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	code := 0
	msg := "success"
	err := m.HpaYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] HpaYamlModify Fail:%s\n", err)
		code = 1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
