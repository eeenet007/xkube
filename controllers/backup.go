// cluster.go
package controllers

import (
	//"fmt"

	"log"
	//"time"
	//"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type BackupController struct {
	beego.Controller
}

// 备份
func (this *BackupController) Backup() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	resName := this.GetString("resName")
	resType := this.GetString("resType")
	remarks := this.GetString("remarks")

	if this.Ctx.Request.Method == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		clusterId = gp.Get("clusterId").String()
		nameSpace = gp.Get("nameSpace").String()
		resName = gp.Get("resName").String()
		resType = gp.Get("resType").String()
		remarks = gp.Get("remarks").String()
	}

	log.Println(resType)

	msg := "success"
	code := 0
	var yaml string
	var err error
	switch resType {
	case "deploy", "deployment":
		yaml, err = m.GetDeployYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "cm", "configmap":
		yaml, err = m.GetCmYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "secret":
		yaml, err = m.GetSecretYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "ing", "ingress":
		yaml, err = m.GetIngYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "ing_beta1", "ingress_beta1":
		yaml, err = m.GetIngYamlV1beta1(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "cronjob":
		yaml, err = m.GetCronjobYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "cronjob_beta1":
		yaml, err = m.GetCronjobYamlBeta1(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "sts", "statefulset":
		yaml, err = m.GetStatefulsetYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "hpa":
		yaml, err = m.GetHpaYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "hpa_beta2":
		yaml, err = m.GetHpaYamlV2beta2(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "svc", "service":
		yaml, err = m.GetSvcYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "node":
		yaml, err = m.GetNodeYaml(clusterId, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "gtw", "gateway":
		yaml, err = m.GetGatewayYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "httproute":
		yaml, err = m.GetHTTPRouteYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "grpcroute":
		yaml, err = m.GetGRPCRouteYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "tcproute":
		yaml, err = m.GetTCPRouteYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	case "udproute":
		yaml, err = m.GetUDPRouteYaml(clusterId, nameSpace, resName)
		if err != nil {
			code = -1
			msg = err.Error()
		}
	default:
		code = -1
		msg = "noSupport"
	}
	err1 := m.InsertBackup(clusterId, nameSpace, resName, resType, yaml, remarks)
	if err1 != nil {
		log.Printf("[ERROR] InsertBackup Fail:%s", err1.Error())
		code = -1
		msg = "InsertBackup Fail"
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 恢复备份
func (this *BackupController) Recover() {
	clusterId := this.GetString("clusterId")
	//namespace := this.GetString("nameSpace")
	//resName := this.GetString("resName")
	resType := this.GetString("resType")
	id, _ := this.GetInt64("id")

	bkst, err := m.GetBackupById(id)
	if err != nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "getBackupYaml Fail"}
		this.ServeJSON()
		return
	}

	code := 0
	msg := "success"
	yaml := bkst.Content
	switch resType {
	case "deploy", "deployment":
		err := m.DeployYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover DeployYaml Fail"
		}
	case "cm", "configmap":
		err := m.CmYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover cm Fail"
		}
	case "secret":
		err := m.SecretYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover secret Fail"
		}
	case "ing", "ingress":
		err := m.IngYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover ing Fail"
		}
	case "ing_beta1", "ingress_beta1":
		err := m.IngYamlModifyV1beta1(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover ing_beta1 Fail"
		}
	case "cronjob":
		err := m.CronjobYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover cronjob Fail"
		}
	case "cronjob_beta1":
		err := m.CronjobYamlModifyBeta1(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover cronjob_beta1 Fail"
		}
	case "sts", "statefulset":
		err := m.StatefulsetYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover Statefulset Fail"
		}
	case "hpa":
		err := m.HpaYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover hpa Fail"
		}
	case "hpa_beta2":
		err := m.HpaYamlModifyV2beta2(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover hpa_beta2 Fail"
		}
	case "svc", "service":
		err := m.SvcYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover svc Fail"
		}
	case "node":
		err := m.NodeYamlModify(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover node Fail"
		}
	case "gateway":
		err := m.UpdateGatewayByYaml(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover gateway Fail"
		}
	case "httproute":
		err := m.UpdateHTTPRouteByYaml(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover httproute Fail"
		}
	case "grpcroute":
		err := m.UpdateGRPCRouteByYaml(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover grpcroute Fail"
		}
	case "tcproute":
		err := m.UpdateTCPRouteByYaml(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover tcproute Fail"
		}
	case "udproute":
		err := m.UpdateUDPRouteByYaml(clusterId, []byte(yaml))
		if err != nil {
			code = -1
			msg = "Recover tcproute Fail"
		}
	default:
		code = -1
		msg = "noSupport"
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 列表
func (this *BackupController) List() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	resName := this.GetString("resName")
	resType := this.GetString("resType")
	starttime := this.GetString("starttime")
	endtime := this.GetString("endtime")
	id, _ := this.GetInt64("id")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")

	if this.Ctx.Request.Method == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		clusterId = gp.Get("clusterId").String()
		nameSpace = gp.Get("nameSpace").String()
		resName = gp.Get("resName").String()
		resType = gp.Get("resType").String()
		starttime = gp.Get("starttime").String()
		endtime = gp.Get("endtime").String()
		page = gp.Get("page").Int()
		page_size = gp.Get("limit").Int()
	}

	msg := "success"
	code := 0
	datas, count := m.GetList_Backup(clusterId, nameSpace, resName, resType, starttime, endtime, id, page, page_size)
	if count == 0 {
		code = -1
		msg = "NoFound"
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &datas}
	this.ServeJSON()
}

func (this *BackupController) View() {
	id, _ := this.GetInt64("id")
	msg := "success"
	code := 0
	data, err := m.GetBackupById(id)
	if err != nil {
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &data}
	this.ServeJSON()
}

// 删除
// func (this *BackupController) Del() {
// 	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
// 	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

// 	Id, _ := this.GetInt64("id")
// 	respid, err := m.DelById_Backup(Id)
// 	if err == nil && respid > 0 {
// 		log.Printf("[INFO] Delsuccess")
// 		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
// 	} else {
// 		log.Printf("[ERROR] Del Error:%s\n", err.Error())
// 		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
// 	}
// }
