package controllers

import (
	//"encoding/json"
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type HpaV2beta2Controller struct {
	beego.Controller
}

func (this *HpaV2beta2Controller) List() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")
	targetName := this.GetString("targetName")
	List, err := m.HpaListV2beta2(clusterId, namespace, hpaName, targetName)
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

func (this *HpaV2beta2Controller) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")

	yamlStr, _ := m.GetHpaYamlV2beta2(clusterId, namespace, hpaName)
	this.Ctx.WriteString(yamlStr)
}

func (this *HpaV2beta2Controller) Create() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	err := m.HpaCreateV2beta2(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] hpa Create Fail:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *HpaV2beta2Controller) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	hpaName := this.GetString("hpaName")
	err := m.HpaDeleteV2beta2(clusterId, namespace, hpaName)
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

func (this *HpaV2beta2Controller) ModifyByYaml() {
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

	err := m.HpaYamlModifyV2beta2(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] HpaYamlModifyV2beta2 Fail:%s\n", err)
		code = 1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
