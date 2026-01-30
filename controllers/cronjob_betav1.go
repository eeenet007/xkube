package controllers

import (
	//"encoding/json"
	"fmt"
	"log"
	"strings"

	//"xkube/common"
	m "xkube/models"
	xm "xkube/xadmin/src/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type CronjobBeta1Controller struct {
	beego.Controller
}

func (this *CronjobBeta1Controller) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		cronjobName = gp.Get("cronjobName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	// if clusterId != "" && cronjobName == "" {
	// 	resp := common.Get("cronjobList" + clusterId + nameSpace)
	// 	if resp != "" {
	// 		var cronjobList = make([]m.Cronjob, 0)
	// 		err := json.Unmarshal([]byte(resp), &cronjobList)
	// 		if err != nil {
	// 			log.Printf("[ERROR] Unmarshal cronjobList Error:%s\n", err)
	// 		} else {
	// 			count := len(cronjobList)
	// 			this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &cronjobList}
	// 			this.ServeJSON()
	// 			return
	// 		}
	// 	}
	// }

	dxList, err := m.CronjobListBeta1(clusterId, nameSpace, cronjobName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(dxList)
	if err != nil {
		log.Printf("[ERROR] CronjobList error:%s \n", err)
		msg = err.Error()
		code = -1
	}
	//set in redis
	// if code == 0 && cronjobName == "" && labels == "" {
	// 	//set数量到redis
	// 	if nameSpace == "" {
	// 		_ = common.HSet("count_"+clusterId, "cronjob", fmt.Sprintf("%d", count))
	// 	}
	// 	bodystr, err := json.Marshal(&dxList)
	// 	if err == nil {
	// 		_ = common.SetEx("cronjobList"+clusterId+nameSpace, string(bodystr), 600)
	// 	}
	// }

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &dxList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *CronjobBeta1Controller) Detail() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")
	xdetail, err := m.CronjobDetailBeta1(clusterId, namespace, cronjobName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *CronjobBeta1Controller) Create() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	log.Println(string(this.Ctx.Input.RequestBody))
	err := m.CronjobCreateBeta1(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] configmap Create Fail:%s\n", err)
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *CronjobBeta1Controller) Modify() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	var key, value string
	if gp.Get("key").Exists() {
		key = gp.Get("key").String()
		value = gp.Get("value").String()
	}
	var code = 0
	var msg = "ok"
	err := m.CronjobModifyBeta1(clusterId, nameSpace, cronjobName, key, value)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] cronjob Modify error:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *CronjobBeta1Controller) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")

	//backup
	yamlStr, _ := m.GetCronjobYamlBeta1(clusterId, nameSpace, cronjobName)
	_ = m.InsertBackup(clusterId, nameSpace, cronjobName, "cronjob_betav1", yamlStr, "Backup before updating")

	code := 0
	msg := "success"
	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.CronjobYamlModifyBeta1(clusterId, []byte(reqBody))
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] cronjob ModifyByYaml Fail:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 迁移到其他集群
func (this *CronjobBeta1Controller) Clone() {
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
		err := m.CronjobCloneBeta1(clusterId, namespace, vv, target_clusterid, target_namespace, targetObjname)
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
		Restype:         "cronjob",
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
		log.Printf("[ERROR] Clone cronjob add Fail:%s\n", err)
		result += fmt.Sprintf("fail:%s", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "success": isOk, "msg": result}
	this.ServeJSON()
}

func (this *CronjobBeta1Controller) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")
	err := m.CronjobDelBeta1(clusterId, namespace, cronjobName)
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

func (this *CronjobBeta1Controller) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")

	yamlStr, _ := m.GetCronjobYamlBeta1(clusterId, namespace, cronjobName)
	this.Ctx.WriteString(yamlStr)
	//this.Data["yaml"] = &yamlStr
	//this.ServeYAML()
	//this.ServeJSON()
}

func (this *CronjobBeta1Controller) Labels() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")

	var labelsMap = make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {
		//log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for i := 0; i <= 20; i++ {
			//log.Println(i)
			kk := gp.Get("labels_key[" + fmt.Sprintf("%d", i) + "]").Str
			vv := gp.Get("labels_value[" + fmt.Sprintf("%d", i) + "]").Str
			if kk != "" && vv != "" {
				log.Println(kk, vv)
				labelsMap[kk] = vv
			} else {
				break
			}
		}
		msg := "success"
		code := 0
		_, err := m.CronjobLabelsBeta1(clusterId, nameSpace, cronjobName, "POST", labelsMap)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	xList, err := m.CronjobLabelsBeta1(clusterId, nameSpace, cronjobName, "GET", labelsMap)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *CronjobBeta1Controller) Run() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	cronjobName := this.GetString("cronjobName")
	err := m.CronjobRunBeta1(clusterId, namespace, cronjobName)
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
