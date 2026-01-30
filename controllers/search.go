// cluster.go
package controllers

import (
	"fmt"

	"log"
	"time"
	"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type SearchController struct {
	beego.Controller
}

//更新搜索表
func (this *SearchController) UpdateIndex() {
	if common.GetLock("updateIndex") {
		this.Ctx.WriteString(`{"code":-1,"msg":"alreadyRunning"}`)
		return
	}
	go func() {
		log.Printf("[INFO] updateIndex start\n")
		_ = common.AddLock("updateIndex", "1", 7200)
		sTime := time.Now().Unix()
		clusterList, _ := m.GetList_Cluster("", "", 1, 100)
		for _, vv := range clusterList {
			num, err := m.DeployInsert(vv.ClusterId)
			log.Printf("cluster:%s,DeployInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.StsInsert(vv.ClusterId)
			log.Printf("cluster:%s,stsInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.DsInsert(vv.ClusterId)
			log.Printf("cluster:%s,dsInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.CronjobInsert(vv.ClusterId)
			log.Printf("cluster:%s,cronjobInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.PodInsert(vv.ClusterId)
			log.Printf("cluster:%s,podInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.SvcInsert(vv.ClusterId)
			log.Printf("cluster:%s,svcInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.IngressInsert(vv.ClusterId)
			log.Printf("cluster:%s,ingressInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.CmInsert(vv.ClusterId)
			log.Printf("cluster:%s,cmInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.SecretInsert(vv.ClusterId)
			log.Printf("cluster:%s,secretInsert num:%d, err:%v", vv.ClusterId, num, err)
			num, err = m.PvcInsert(vv.ClusterId)
			log.Printf("cluster:%s,pvcInsert num:%d, err:%v", vv.ClusterId, num, err)
		}
		eTime := time.Now().Unix()
		_ = common.Del("updateIndex")
		log.Printf("[INFO] updateIndex end time:%d\n", eTime-sTime)
	}()
	resp := fmt.Sprintf(`{"code":0,"msg":"success"}`)
	this.Ctx.WriteString(resp)
}

//列表
func (this *SearchController) List() {
	keyword := this.GetString("keyword")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	if gp.Get("keyword").Exists() {
		keyword = gp.Get("keyword").String()
	}

	msg := "success"
	code := 0
	datas, count := m.GetList_Index(keyword)
	if count == 0 {
		code = -1
		msg = "NoSearchFound"
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &datas}
	this.ServeJSON()
}

//删除
func (this *SearchController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	respid, err := m.DelById_Index(Id)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}
