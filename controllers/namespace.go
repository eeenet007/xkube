// node
package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	m "xkube/models"

	"github.com/tidwall/gjson"

	beego "github.com/beego/beego/v2/server/web"
)

type NsController struct {
	beego.Controller
}

func (this *NsController) List() {
	clusterId := this.GetString("clusterId")
	xList, err := m.NsList(clusterId)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *NsController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	xdetail, err := m.NsDetail(clusterId, nameSpace)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *NsController) Yaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	yamlStr, _ := m.NsYaml(clusterId, nameSpace)
	this.Ctx.WriteString(yamlStr)
}

func (this *NsController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")

	yamlStr, _ := m.NsYaml(clusterId, nameSpace)
	_ = m.InsertBackup(clusterId, nameSpace, nameSpace, "namespace", yamlStr, "Backup before updating")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	err := m.NsYamlUpdate(clusterId, []byte(reqBody))
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *NsController) Create() {
	clusterId := this.GetString("clusterId")
	var labelsMap = make(map[string]string)
	//log.Println(string(this.Ctx.Input.RequestBody))
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	for i := 0; i <= 20; i++ {
		kk := gp.Get("labels_key[" + fmt.Sprintf("%d", i) + "]").Str
		vv := gp.Get("labels_value[" + fmt.Sprintf("%d", i) + "]").Str
		if kk != "" && vv != "" {
			//log.Println(kk, vv)
			labelsMap[kk] = vv
		} else {
			break
		}
	}
	code := 0
	msg := "success"
	if clusterId == "" {
		clusterId = gp.Get("clusterId").Str
	}
	namespace := gp.Get("nameSpace").Str
	err := m.CreateNs(clusterId, namespace, labelsMap)
	if err != nil {
		log.Printf("[ERROR] CreateNs Fail:%s,%s\n", clusterId, namespace)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *NsController) Del() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	force := this.GetString("force")
	err := m.DeleteNs(clusterId, nameSpace, force)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success"}
	this.ServeJSON()
}

// func (this *NsController) LimitRangeCreate() {
// 	clusterId := this.GetString("clusterId")
// 	nameSpace := this.GetString("nameSpace")

// 	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
// 	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

// 	var limitRange m.LimitRes
// 	err := json.Unmarshal(this.Ctx.Input.RequestBody, &limitRange)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	err = m.CreateLimitRange(clusterId, nameSpace, limitRange)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success"}
// 	this.ServeJSON()
// }

// func (this *NsController) LimitRangeUpdate() {
// 	clusterId := this.GetString("clusterId")
// 	nameSpace := this.GetString("nameSpace")

// 	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
// 	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

// 	var limitRange m.LimitRes
// 	err := json.Unmarshal(this.Ctx.Input.RequestBody, &limitRange)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	err = m.UpdateLimitRange(clusterId, nameSpace, limitRange)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success"}
// 	this.ServeJSON()
// }

func (this *NsController) LimitRange() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	code := 0
	msg := "success"

	if this.Ctx.Input.Method() == "POST" {
		var limitRange m.LimitRes
		err := json.Unmarshal(this.Ctx.Input.RequestBody, &limitRange)
		if err != nil {
			log.Println(err)
		}
		_, err = m.GetLimitRange(clusterId, nameSpace) //查询是否存在
		if err != nil {
			log.Printf("[WARN] notFound LimitRange:%s,%s\n", clusterId, nameSpace) //不存在就创建
			err = m.CreateLimitRange(clusterId, nameSpace, limitRange)
			if err != nil {
				log.Printf("[ERROR] LimitRange Create error:%s\n", err)
				code = -1
				msg = err.Error()
			}
		} else { //存在就更新
			err = m.UpdateLimitRange(clusterId, nameSpace, limitRange)
			if err != nil {
				log.Printf("[ERROR] LimitRange Update error:%s\n", err)
				code = -1
				msg = err.Error()
			}
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}

	xlr, err := m.GetLimitRange(clusterId, nameSpace)
	if err != nil {
		log.Printf("[ERROR] LimitRange Get error:%s\n", err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &xlr}
	this.ServeJSON()
}
