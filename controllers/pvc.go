package controllers

import (
	"encoding/json"
	"log"
	"strings"
	"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type PvcController struct {
	beego.Controller
}

func (this *PvcController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	pvcName := this.GetString("pvcName")

	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		pvcName = gp.Get("pvcName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && pvcName == "" {
		resp := common.Get("pvcList" + clusterId + nameSpace)
		if resp != "" {
			var pvcList = make([]m.PersistentVolumeClaim, 0)
			err := json.Unmarshal([]byte(resp), &pvcList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal pvcList Error:%s\n", err)
			} else {
				count := len(pvcList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &pvcList}
				this.ServeJSON()
				return
			}
		}
	}

	pvcList, err := m.PersistentVolumeClaimList(clusterId, nameSpace, pvcName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(pvcList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && pvcName == "" && labels == "" {
		bodystr, err := json.Marshal(&pvcList)
		if err == nil {
			_ = common.SetEx("pvcList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &pvcList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *PvcController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	pvcName := this.GetString("pvcName")
	xdetail, err := m.PersistentVolumeClaimDetail(clusterId, nameSpace, pvcName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *PvcController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	pvcName := this.GetString("pvcName")
	yamlStr, _ := m.GetPersistentVolumeClaimYaml(clusterId, nameSpace, pvcName)
	this.Ctx.WriteString(yamlStr)
}
