// cluster.go
package rbac

import (
	"fmt"

	"log"
	//"time"
	m "xkube/xadmin/src/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type ClusterController struct {
	beego.Controller
}

// 列表
func (this *ClusterController) List() {
	//id := this.GetString("id")

	if this.Ctx.Request.Method == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		datas, count := m.GetUserClusterListV2(gp.Map())
		this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
		this.ServeJSON()
		return
	} else {
		username := this.GetString("username")
		clusterId := this.GetString("clusterId")
		page, _ := this.GetInt64("page")
		page_size, _ := this.GetInt64("limit")
		msg := "success"
		datas, count := m.GetUserClusterList(username, clusterId, page, page_size)
		this.Data["json"] = &map[string]interface{}{"code": 0, "msg": msg, "count": count, "data": &datas}
		this.ServeJSON()
	}
}

// 添加
func (this *ClusterController) Add() {
	u := m.UserCluster{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	xff := this.Ctx.Request.Header.Get("x-forwarded-for")
	remoteAddr := this.Ctx.Request.RemoteAddr
	user_ip := fmt.Sprintf("%s,%s", xff, remoteAddr)
	userinfo := this.GetSession("userinfo")

	if gp.Get("clusterId").Exists() {
		u.ClusterId = gp.Get("clusterId").String()
		u.Username = gp.Get("username").String()

		id, err := m.AddUserCluster(&u)
		if err == nil && id > 0 {
			log.Printf("[INFO] Add success")
			_ = m.InsertLogAudit(user_ip, userinfo.(m.User).Username, "grantPermission", "success", fmt.Sprintf("clusetrId:%s,user:%s,id:%d", u.ClusterId, u.Username, id))
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

// 删除
func (this *ClusterController) Delete() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	xff := this.Ctx.Request.Header.Get("x-forwarded-for")
	remoteAddr := this.Ctx.Request.RemoteAddr
	user_ip := fmt.Sprintf("%s,%s", xff, remoteAddr)
	userinfo := this.GetSession("userinfo")

	Id, _ := this.GetInt64("id")
	respid, err := m.DelUserClusterById(Id)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		_ = m.InsertLogAudit(user_ip, userinfo.(m.User).Username, "delPermission", "success", fmt.Sprintf("id:%d", Id))
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}

func (this *ClusterController) MyClusterList() {
	userinfo := this.GetSession("userinfo")
	username := userinfo.(m.User).Username
	//page, _ := this.GetInt64("page")
	//page_size, _ := this.GetInt64("limit")
	msg := "success"
	code := 0
	//datas, count := m.GetUserClusterList(username, "", page, page_size)
	datas, _ := m.GetMyClusterList(username)
	count := len(datas)
	if count == 0 {
		if username == "admin" {
			msg = "AllClustersPermissions"
		} else {
			msg = "NoAnyClusterPermissions"
		}
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &datas}
	this.ServeJSON()

}
