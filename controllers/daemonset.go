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

type DaemonsetController struct {
	beego.Controller
}

func (this *DaemonsetController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	daemonsetName := this.GetString("daemonsetName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		daemonsetName = gp.Get("daemonsetName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && daemonsetName == "" && labels == "" {
		resp := common.Get("dsList" + clusterId + nameSpace)
		if resp != "" {
			var dxList = make([]m.Daemonset, 0)
			err := json.Unmarshal([]byte(resp), &dxList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal dsList Error:%s\n", err)
			} else {
				count := len(dxList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &dxList}
				this.ServeJSON()
				return
			}
		}
	}

	dxList, err := m.DaemonsetList(clusterId, nameSpace, daemonsetName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(dxList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	//set in redis
	if code == 0 && daemonsetName == "" && labels == "" {
		bodystr, err := json.Marshal(&dxList)
		if err == nil {
			_ = common.SetEx("dsList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &dxList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *DaemonsetController) Detail() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	daemonsetName := this.GetString("daemonsetName")
	xdetail, err := m.DaemonsetDetail(clusterId, namespace, daemonsetName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *DaemonsetController) Create() {
	clusterId := this.GetString("clusterId")
	var sts m.Daemonset

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sts)
	err = m.DaemonsetCreate(clusterId, &sts)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *DaemonsetController) Modify() {
	clusterId := this.GetString("clusterId")
	var sts m.Daemonset

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sts)
	err = m.DaemonsetModify(clusterId, &sts)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *DaemonsetController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	daemonsetName := this.GetString("daemonsetName")
	err := m.DaemonsetDel(clusterId, namespace, daemonsetName)
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

func (this *DaemonsetController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	daemonsetName := this.GetString("daemonsetName")

	yamlStr, _ := m.GetDaemonsetYaml(clusterId, namespace, daemonsetName)
	this.Ctx.WriteString(yamlStr)
	//this.Data["yaml"] = &yamlStr
	//this.ServeYAML()
	//this.ServeJSON()
}

func (this *DaemonsetController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	daemonsetName := this.GetString("daemonsetName")

	//backup
	yamlStr, _ := m.GetDaemonsetYaml(clusterId, nameSpace, daemonsetName)
	_ = m.InsertBackup(clusterId, nameSpace, daemonsetName, "daemonset", yamlStr, "Backup before updating")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	code := 0
	msg := "success"
	err := m.DaemonsetYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] DaemonsetYamlModify Fail:%s\n", err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
