// cluster.go
package controllers

import (
	//"fmt"

	"encoding/json"
	"log"

	//"time"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	//"github.com/tidwall/gjson"
)

type CloneController struct {
	beego.Controller
}

// 列表
func (this *CloneController) List() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin", "count": 0}
		this.ServeJSON()
		return
	}
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")

	datas, count := m.GetList_Clone(page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

// 添加
func (this *CloneController) Add() {
	u := m.Xkb_clone_log{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &u)
	id, err := m.Add_Clone(&u)
	if err == nil && id > 0 {
		_ = m.ClearCache(u.TargetClusterid) //创建以后，刷新一下列表缓存
		log.Printf("[INFO] Add success")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Add Fail:%s", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
	}
}

// 删除
func (this *CloneController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	respid, err := m.DelById_Clone(Id)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}
