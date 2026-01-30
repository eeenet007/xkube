// cluster.go
package controllers

import (
	"fmt"

	"encoding/json"
	"log"

	//"time"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type CicdController struct {
	beego.Controller
}

// 列表
func (this *CicdController) List() {
	id := this.GetString("id")
	cicdName := this.GetString("cicdName")
	appname := this.GetString("appname")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")
	datas, count := m.GetList_Cicd(id, cicdName, appname, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) ListAppname() {
	datas, count := m.CicdListAppname()
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

// 获取cicd的配置及类似
func (this *CicdController) GetCicdInfo() {
	cicdName := this.GetString("cicdName")
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	code := 0
	msg := "success"
	if cicdName == "" {
		code = -1
		msg = "need cicdName"
	}
	datas, err := m.GetCicdByCicdName(cicdName, clusterId, nameSpace)
	if err != nil {
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) Add() {
	u := m.Xkb_cicd{}
	b := m.Xkb_cicd_pipelines{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	if err != nil {
		log.Printf("[ERROR] Umarshal Fail:%s\n", err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "UmarshalFail"}
		this.ServeJSON()
		return
	}
	err2 := json.Unmarshal(this.Ctx.Input.RequestBody, &b)
	if err2 != nil {
		log.Printf("[ERROR] Umarshal Fail:%s\n", err2)
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "UmarshalFail2"}
		this.ServeJSON()
		return
	}
	id, err := m.Add_Cicd_Pipelines(&u, &b)
	if err == nil && id > 0 {
		log.Printf("[INFO] Add success")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Add Fail:%s", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
	}
}

func (this *CicdController) PostStatus() {
	cicdId := this.GetString("cicdId")
	status := this.GetString("status")
	lastRunTime := this.GetString("lastRunTime")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, err := m.UpdateStatus_Cicd(cicdId, status, lastRunTime)
	if err == nil {
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		this.Ctx.WriteString(fmt.Sprintf(`{"code":-1,"msg":"%s"}`, err.Error()))
	}
}

func (this *CicdController) Update() {
	u := m.Xkb_cicd{}
	b := m.Xkb_cicd_pipelines{}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	if err != nil {
		log.Printf("[ERROR] Umarshal Fail:%s\n", err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "UmarshalFail"}
		this.ServeJSON()
		return
	}
	err2 := json.Unmarshal(this.Ctx.Input.RequestBody, &b)
	if err2 != nil {
		log.Printf("[ERROR] Umarshal Fail:%s\n", err2)
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "UmarshalFail2"}
		this.ServeJSON()
		return
	}
	id, err := m.Update_Cicd(&u, &b)
	if err == nil && id > 0 {
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] updatecicd Fail:%s\n", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
	}
}

// 删除
func (this *CicdController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	respid, err := m.DelById_Cicd(Id)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}

// 更新某个字段
func (this *CicdController) Edit() {
	id := this.GetString("id")
	table := this.GetString("tb")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	if gp.Get("name").Exists() {
		name := gp.Get("name").String()
		value := gp.Get("value").String()
		_, err := m.Edit_data(id, name, value, table)
		if err == nil {
			this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
		} else {
			log.Printf("[ERROR] update Fail:%s\n", err)
			this.Ctx.WriteString(fmt.Sprintf(`{"code":-1,"msg":"%s"}`, err.Error()))
		}
	} else {
		this.Ctx.WriteString(`{"code":-1,"msg":"bodyParseFail"}`)
	}
}

// 根据获取aliyun pipelines配置
func (this *CicdController) GetPipelines() {
	cicdId, _ := this.GetInt64("cicdId")
	code := 0
	msg := "success"
	if cicdId <= 0 {
		code = -1
		msg = "cicdId error"
	}
	datas, err := m.GetPipelinesByCicdId(cicdId)
	if err != nil {
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) GetAliyunIdList() {
	aliyun_id := this.GetString("aliyun_id")
	code := 0
	msg := "success"
	datas, num := m.GetAliyunIdList(aliyun_id)
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": num, "data": &datas}
	this.ServeJSON()
}

// xkb_cicd_ak列表
func (this *CicdController) AkList() {
	id := this.GetString("id")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")
	datas, count := m.GetAkList(id, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) AkAdd() {
	u := m.Xkb_cicd_ak{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	if err != nil {
		log.Printf("[ERROR] Umarshal Fail:%s\n", err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "UmarshalFail"}
		this.ServeJSON()
		return
	}
	id, err := m.Add_Ak(&u)
	if err == nil && id > 0 {
		log.Printf("[INFO] Add success")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Add Fail:%s", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
	}
}

// 删除
func (this *CicdController) AkDel() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	aliyunId := this.GetString("aliyunId")
	respid, err := m.DelById_Ak(Id, aliyunId)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}

// xkb_cicd_jks列表
func (this *CicdController) JksList() {
	id := this.GetString("id")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")
	datas, count := m.GetJksList(id, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) GetJksList() {
	datas, count := m.GetJksListNopasswd()
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) JksAdd() {
	u := m.Xkb_cicd_jks{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	if err != nil {
		log.Printf("[ERROR] Umarshal Fail:%s\n", err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "UmarshalFail"}
		this.ServeJSON()
		return
	}
	id, err := m.Add_Jks(&u)
	if err == nil && id > 0 {
		//log.Printf("[INFO] Add success")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		//log.Printf("[ERROR] Add Fail:%s", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
	}
}

func (this *CicdController) JksDel() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	jksId := this.GetString("jksId")
	respid, err := m.DelById_Jks(Id, jksId)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}
