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

type ConfigMapController struct {
	beego.Controller
}

func (this *ConfigMapController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	configmapName := this.GetString("configmapName")

	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		configmapName = gp.Get("configmapName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && configmapName == "" && labels == "" {
		resp := common.Get("cmList" + clusterId + nameSpace)
		if resp != "" {
			var cmList = make([]m.Configmap, 0)
			err := json.Unmarshal([]byte(resp), &cmList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal cmList Error:%s\n", err)
			} else {
				count := len(cmList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &cmList}
				this.ServeJSON()
				return
			}
		}
	}

	cmList, err := m.CmList(clusterId, nameSpace, configmapName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(cmList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && configmapName == "" && labels == "" {
		//set数量到redis
		if nameSpace == "" {
			_ = common.HSet("count_"+clusterId, "configmap", fmt.Sprintf("%d", count))
		}
		bodystr, err := json.Marshal(&cmList)
		if err == nil {
			_ = common.SetEx("cmList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &cmList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *ConfigMapController) Create() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	//log.Println(string(this.Ctx.Input.RequestBody))
	err := m.CmCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] configmap Create Fail:%s\n", err)
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *ConfigMapController) CreateByYaml() {
	clusterId := this.GetString("clusterId")
	err := m.CmYamlCreate(clusterId, this.Ctx.Input.RequestBody)
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

func (this *ConfigMapController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	configmapName := this.GetString("configmapName")
	code := 0
	msg := "success"

	//backup
	yamlStr, _ := m.GetCmYaml(clusterId, namespace, configmapName)
	_ = m.InsertBackup(clusterId, namespace, configmapName, "configmap", yamlStr, "Backup before updating")

	//log.Println(string(this.Ctx.Input.RequestBody))

	//bodyByte := []byte(strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%"))
	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.CmYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] configmap ModifyByYaml Fail:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 迁移到其他集群
func (this *ConfigMapController) Clone() {
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
		err := m.CmClone(clusterId, namespace, vv, target_clusterid, target_namespace, targetObjname)
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
		Restype:         "configmap",
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
		log.Printf("[ERROR] Clone cm add Fail:%s\n", err)
		result += fmt.Sprintf("fail:%s", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "success": isOk, "msg": result}
	this.ServeJSON()
}

func (this *ConfigMapController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	configmapName := this.GetString("configmapName")
	xdetail, err := m.CmDetail(clusterId, nameSpace, configmapName)
	if err != nil {
		log.Println(err)
	}
	_, err = json.Marshal(&xdetail)
	if err != nil {
		fmt.Printf("err=%v\n", err)
	}
	//fmt.Printf("%s", string(data))

	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *ConfigMapController) Modify() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	var configmap m.Configmap

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &configmap)
	err = m.CmUpdate(clusterId, namespace, &configmap)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *ConfigMapController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	configmapName := this.GetString("configmapName")

	yamlStr, _ := m.GetCmYaml(clusterId, namespace, configmapName)
	this.Ctx.WriteString(yamlStr)
}

func (this *ConfigMapController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	configmapName := this.GetString("configmapName")
	err := m.CmDelete(clusterId, namespace, configmapName)
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
