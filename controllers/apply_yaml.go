package controllers

import (
	//"encoding/json"
	"fmt"
	"log"
	"strings"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
)

type ApplyYamlController struct {
	beego.Controller
}

func (this *ApplyYamlController) ApplyYaml() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"
	if clusterId == "" {
		code = -1
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": "Need clusterId"}
		this.ServeJSON()
		return
	}
	//nameSpace := gp.Get("nameSpace").String()
	//log.Println(string(this.Ctx.Input.RequestBody))
	//bodyByte := []byte(strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")) //将url编码后的%替换成%
	err := m.ApplyYaml(clusterId, strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%"))
	if err != nil {
		code = -1
		msg = fmt.Sprintf("%s", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *ApplyYamlController) CreateByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	msg := "success"
	code := 0
	if clusterId == "" {
		msg = "need ClusterId"
		code = 0
	}
	//log.Println(string(this.Ctx.Input.RequestBody))
	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	err := m.CreateByYaml(clusterId, nameSpace, []byte(reqBody))
	if err != nil {
		log.Printf("[INFO] CreateByYaml err:%s\n", err)
		msg = err.Error()
		code = -1
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
