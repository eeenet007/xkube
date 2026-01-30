package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"xkube/common"
	m "xkube/models"
	xm "xkube/xadmin/src/models"

	"github.com/tidwall/gjson"

	beego "github.com/beego/beego/v2/server/web"
)

type SvcController struct {
	beego.Controller
}

func (this *SvcController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	serviceName := this.GetString("serviceName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		serviceName = gp.Get("serviceName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}
	//get from redis
	if clusterId != "" && serviceName == "" && labels == "" {
		resp := common.Get("svcList" + clusterId + nameSpace)
		if resp != "" {
			var svcList = make([]m.Service, 0)
			err := json.Unmarshal([]byte(resp), &svcList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal svcList Error:%s\n", err)
			} else {
				count := len(svcList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &svcList}
				this.ServeJSON()
				return
			}
		}
	}

	svcList, err := m.SvcList(clusterId, nameSpace, serviceName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(svcList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && serviceName == "" && labels == "" {
		//set数量到redis
		if nameSpace == "" {
			_ = common.HSet("count_"+clusterId, "service", fmt.Sprintf("%d", count))
		}
		bodystr, err := json.Marshal(&svcList)
		if err == nil {
			_ = common.SetEx("svcList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": msg, "count": count, "data": &svcList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *SvcController) Create() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	log.Println(string(this.Ctx.Input.RequestBody))
	err := m.SvcCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] configmap Create Fail:%s\n", err)
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *SvcController) CreateByYaml() {
	clusterId := this.GetString("clusterId")
	err := m.SvcYamlCreate(clusterId, this.Ctx.Input.RequestBody)
	msg := "success"
	code := 0
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *SvcController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	serviceName := this.GetString("serviceName")

	//backup
	yamlStr, _ := m.GetSvcYaml(clusterId, nameSpace, serviceName)
	_ = m.InsertBackup(clusterId, nameSpace, serviceName, "service", yamlStr, "Backup before updating")

	//nameSpace := gp.Get("nameSpace").String()
	//log.Println(string(this.Ctx.Input.RequestBody))
	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	code := 0
	msg := "success"
	err := m.SvcYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] SvcYamlModify Fail:%s\n", err)
		code = 1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 迁移到其他集群
func (this *SvcController) Clone() {
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	clusterId := gp.Get("clusterid").String()
	namespace := gp.Get("namespace").String()
	target_clusterid := gp.Get("target_clusterid").String()
	target_namespace := gp.Get("target_namespace").String()
	target_objname := gp.Get("target_objname").String()
	objName := gp.Get("objname").String()
	objNameArry := strings.Split(objName, ",")
	targetObjnameArry := strings.Split(target_objname, ",")
	var result string
	var isOk = true

	for i, vv := range objNameArry {
		var targetObjname string
		if len(targetObjnameArry)-1 >= i {
			targetObjname = targetObjnameArry[i]
		}
		err := m.SvcClone(clusterId, namespace, vv, target_clusterid, target_namespace, targetObjname)
		if err != nil {
			isOk = false
		} else {
			err = fmt.Errorf("ok")
		}
		result += fmt.Sprintf("%s result:%s,", vv, err)
	}
	uinfo := this.GetSession("userinfo")
	u := m.Xkb_clone_log{
		Clusterid:       clusterId,
		Namespace:       namespace,
		Restype:         "service",
		Objname:         objName,
		TargetClusterid: target_clusterid,
		TargetNamespace: target_namespace,
		TargetObjname:   target_objname,
		Status:          fmt.Sprintf("%v", isOk),
		User:            uinfo.(xm.User).Username,
		Result:          result,
	}
	_, err := m.Add_Clone(&u) //插入clone 结果
	if err != nil {
		log.Printf("[ERROR] Clone service add Fail:%s\n", err)
		result += fmt.Sprintf("insert log fail")
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "success": isOk, "msg": result}
	this.ServeJSON()
}

func (this *SvcController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	serviceName := this.GetString("serviceName")
	xdetail, err := m.SvcDetail(clusterId, nameSpace, serviceName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *SvcController) Del() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	serviceName := this.GetString("serviceName")
	err := m.SvcDelete(clusterId, nameSpace, serviceName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success"}
	this.ServeJSON()
}

func (this *SvcController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	serviceName := this.GetString("serviceName")

	yamlStr, _ := m.GetSvcYaml(clusterId, nameSpace, serviceName)
	this.Ctx.WriteString(yamlStr)
}
