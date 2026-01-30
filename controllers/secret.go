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

type SecretController struct {
	beego.Controller
}

func (this *SecretController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	secretName := this.GetString("secretName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		secretName = gp.Get("secretName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && secretName == "" && labels == "" {
		resp := common.Get("secList" + clusterId + nameSpace)
		if resp != "" {
			var secList = make([]m.Secret, 0)
			err := json.Unmarshal([]byte(resp), &secList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal secList Error:%s\n", err)
			} else {
				count := len(secList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &secList}
				this.ServeJSON()
				return
			}
		}
	}

	secList, err := m.SecretList(clusterId, nameSpace, secretName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(secList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && secretName == "" && labels == "" {
		//set数量到redis
		if nameSpace == "" {
			_ = common.HSet("count_"+clusterId, "secret", fmt.Sprintf("%d", count))
		}
		bodystr, err := json.Marshal(&secList)
		if err == nil {
			_ = common.SetEx("secretList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &secList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *SecretController) Create() {
	// clusterId := this.GetString("clusterId")
	// namespace := this.GetString("nameSpace")
	// var secret m.Secret

	// err := json.Unmarshal(this.Ctx.Input.RequestBody, &secret)
	// err = m.SecretCreate(clusterId, namespace, &secret)
	// if err != nil {
	// 	log.Println(err)
	// }
	// this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	// this.ServeJSON()

	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	log.Println(string(this.Ctx.Input.RequestBody))
	err := m.SecretCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] configmap Create Fail:%s\n", err)
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *SecretController) CreateByYaml() {
	clusterId := this.GetString("clusterId")
	err := m.SecretYamlCreate(clusterId, this.Ctx.Input.RequestBody)
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

func (this *SecretController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	//nameSpace := gp.Get("nameSpace").String()
	log.Println(string(this.Ctx.Input.RequestBody))
	bodyByte := []byte(strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%"))
	err := m.SecretYamlModify(clusterId, bodyByte)
	msg := "success"
	code := 0
	if err != nil {
		log.Printf("[WARN] SecretYamlModify Fail:%s\n", err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 迁移到其他集群
func (this *SecretController) Clone() {
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
		err := m.SecretClone(clusterId, namespace, vv, target_clusterid, target_namespace, targetObjname)
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
		Restype:         "secret",
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
		log.Printf("[ERROR] Clone secret add Fail:%s\n", err)
		result += fmt.Sprintf("insert log fail")
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "success": isOk, "msg": result}
	this.ServeJSON()
}

func (this *SecretController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	secretName := this.GetString("secretName")
	xdetail, err := m.SecretDetail(clusterId, nameSpace, secretName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *SecretController) Modify() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	var secret m.Secret

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &secret)
	err = m.SecretUpdate(clusterId, namespace, &secret)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *SecretController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	secretName := this.GetString("secretName")

	yamlStr, _ := m.GetSecretYaml(clusterId, namespace, secretName)
	this.Ctx.WriteString(yamlStr)
}

func (this *SecretController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	secretName := this.GetString("secretName")
	err := m.SecretDelete(clusterId, namespace, secretName)
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
