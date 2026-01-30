// cluster.go
package controllers

import (
	"fmt"

	"log"
	//"time"
	"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type ClusterController struct {
	beego.Controller
}

// 列表
func (this *ClusterController) List() {

	// userinfo := this.GetSession("userinfo")
	// if userinfo == nil {
	// 	this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin", "count": 0}
	// 	this.ServeJSON()
	// 	return
	// }
	// for kk, vvs := range this.Ctx.Request.Header {
	// 	for _, v1 := range vvs {
	// 		log.Printf("%s:%s\n", kk, v1)
	// 	}
	// }

	id := this.GetString("id")
	clusterId := this.GetString("clusterId")
	//clusterName := this.GetString("clusterName")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")

	datas, count := m.GetList_Cluster(id, clusterId, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &datas}
	this.ServeJSON()
}

func (this *ClusterController) Detail() {
	// userinfo := this.GetSession("userinfo")
	// if userinfo == nil {
	// 	this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin"}
	// 	this.ServeJSON()
	// 	return
	// }
	id, _ := this.GetInt64("id")
	datas, err := m.GetDetail_Cluster(id)
	if err != nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "DeailFail"}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "data": &datas}
	this.ServeJSON()
}

// 添加
func (this *ClusterController) Add() {
	u := m.Xkb_cluster{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	if gp.Get("cluster_id").Exists() {
		u.ClusterId = gp.Get("cluster_id").String()
		u.ClusterName = gp.Get("cluster_name").String()
		u.KubeVersion = gp.Get("kube_version").String()
		u.IdcName = gp.Get("idc_name").String()
		u.KubeConfig = gp.Get("kube_config").String()
		u.BearerToken = gp.Get("bearer_token").String()
		u.LanSlbip = gp.Get("lan_slbip").String()
		u.WanSlbip = gp.Get("wan_slbip").String()
		u.Status = gp.Get("status").Int()
		u.Remarks = gp.Get("remarks").String()

		id, err := m.Add_Cluster(&u)
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

// 删除
func (this *ClusterController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	cluster_id := this.GetString("cluster_id")
	respid, err := m.DelById_Cluster(Id)
	if err == nil && respid > 0 {
		_ = common.Del("count_" + cluster_id)
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}

// 更新某个字段
func (this *ClusterController) Edit() {
	id := this.GetString("id")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	if gp.Get("name").Exists() {
		name := gp.Get("name").String()
		value := gp.Get("value").String()
		_, err := m.Edit_Cluster(id, name, value)
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

func (this *ClusterController) Update() {
	u := m.Xkb_cluster{}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	// log.Println(string(this.Ctx.Input.RequestBody))
	// log.Println(gp.Get("cluster_id").String())
	// log.Println(gp.Get("id").Int())
	if gp.Get("id").Int() > 0 {
		u.Id = gp.Get("id").Int()
		u.ClusterId = gp.Get("cluster_id").String()
		u.ClusterName = gp.Get("cluster_name").String()
		u.IdcName = gp.Get("idc_name").String()
		u.KubeVersion = gp.Get("kube_version").String()
		u.KubeConfig = gp.Get("kube_config").String()
		u.BearerToken = gp.Get("bearer_token").String()
		u.LanSlbip = gp.Get("lan_slbip").String()
		u.WanSlbip = gp.Get("wan_slbip").String()
		u.Status = gp.Get("status").Int()
		u.Remarks = gp.Get("remarks").String()

		id, err := m.Update_Cluster(&u)
		if err == nil && id > 0 {
			this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
		} else {
			log.Printf("[ERROR] updatekubeCluster Fail:%s\n", err)
			this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
		}
	} else {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "bodyParseFail"}
		this.ServeJSON()
	}
}

func (this *ClusterController) Count() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	resp, _ := m.CountWorkLoad()
	count := len(resp)
	code := 0
	msg := "success"
	if count < 1 {
		code = -1
		msg = "nodata"
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &resp}
	this.ServeJSON()

}
