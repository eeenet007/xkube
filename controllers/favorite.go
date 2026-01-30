// cluster.go
package controllers

import (
	//"fmt"

	"log"
	//"time"
	m "xkube/models"
	adm "xkube/xadmin/src/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type FavController struct {
	beego.Controller
}

//列表
func (this *FavController) List() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin", "count": 0}
		this.ServeJSON()
		return
	}
	username := userinfo.(adm.User).Username
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")
	datas, count := m.GetList_Favorite(username, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

//添加
func (this *FavController) Add() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin"}
		this.ServeJSON()
		return
	}

	u := m.Xkb_favorite{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	//fav_type := this.GetString("fav_type")
	username := userinfo.(adm.User).Username
	if gp.Get("fav_type").Exists() {
		u.FavType = gp.Get("fav_type").String()
		u.ClusterId = gp.Get("cluster_id").String()
		u.NameSpace = gp.Get("name_space").String()
		u.FavName = gp.Get("fav_name").String()
		u.Username = username

		id, err := m.Add_Favorite(&u)
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
func (this *FavController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin"}
		this.ServeJSON()
		return
	}
	username := userinfo.(adm.User).Username
	Id, _ := this.GetInt64("id")
	respid, err := m.DelById_Favorite(Id, username)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}
