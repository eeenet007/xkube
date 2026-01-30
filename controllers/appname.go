// cluster.go
package controllers

import (
	"fmt"

	"log"
	//"time"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type AppnameController struct {
	beego.Controller
}

//列表
func (this *AppnameController) List() {
	appname := this.GetString("appname")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")

	datas, count := m.GetList_Appname(appname, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

//添加
func (this *AppnameController) Add() {
	u := m.XkbAppname{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	if gp.Get("appname").Exists() {
		u.Appname = gp.Get("appname").String()
		u.Remarks = gp.Get("remarks").String()
		id, err := m.Add_Appname(&u)
		if err == nil && id > 0 {
			log.Printf("[INFO] Add success")
			this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
		} else {
			log.Printf("[ERROR] Add Fail:%s", err)
			this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
		}
	} else {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "bodyParseFail"}
		this.ServeJSON()
	}
}

//删除
func (this *AppnameController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	respid, err := m.DelById_Appname(Id)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}

//更新某个字段
func (this *AppnameController) Edit() {
	id := this.GetString("id")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	if gp.Get("name").Exists() {
		name := gp.Get("name").String()
		value := gp.Get("value").String()
		_, err := m.Edit_Appname(id, name, value)
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

func (this *AppnameController) Update() {
	u := m.XkbAppname{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	if gp.Get("id").Int() > 0 {
		u.Id = gp.Get("id").Int()
		u.Appname = gp.Get("appname").String()
		u.Remarks = gp.Get("remarks").String()
		id, err := m.Update_Appname(&u)
		if err == nil && id > 0 {
			this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
		} else {
			log.Printf("[ERROR] updateAppname Fail:%s\n", err)
			this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
		}
	} else {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "bodyParseFail"}
		this.ServeJSON()
	}
}
