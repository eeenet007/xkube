package controllers

import (
	m "xkube/models"

	//beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

func (this *CicdController) Run() {
	jksId := this.GetString("jksId")
	jobName := this.GetString("jobName")
	var paramsMap = make(map[string]string)
	if this.Ctx.Request.Method == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for _, vv := range gp.Array() {
			paramsMap[vv.Get("key").String()] = vv.Get("value").String()
		}
	}
	buildId, err := m.RunJobBuild(jksId, jobName, paramsMap)
	if err != nil {
		this.Data["json"] = &map[string]interface{}{"code": 1, "msg": err.Error()}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "buildId": buildId}
	this.ServeJSON()
}

// func (this *CicdController) GetJobInfo() {
// 	jobName := this.GetString("jobName")
// 	result, err := m.GetJobInfo(jobName)
// 	if err != nil {
// 		this.Data["json"] = &map[string]interface{}{"code": 1, "msg": err.Error()}
// 	} else {
// 		this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "data": result}
// 	}
// 	this.ServeJSON()
// }

func (this *CicdController) BuildList() {
	jobName := this.GetString("jobName")
	jksId := this.GetString("jksId")
	datas, _ := m.GetBuildList(jksId, jobName)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(datas), "data": &datas}
	this.ServeJSON()
}

func (this *CicdController) BuildLog() {
	jksId := this.GetString("jksId")
	jobName := this.GetString("jobName")
	buildId, _ := this.GetInt64("buildId")
	logtext, err := m.GetBuildLog(jksId, jobName, buildId)
	if err != nil {
		this.Ctx.WriteString(err.Error())
	}
	this.Ctx.WriteString(logtext)
}

func (this *CicdController) BuildState() {
	jksId := this.GetString("jksId")
	jobName := this.GetString("jobName")
	buildId, _ := this.GetInt64("buildId")
	bs, err := m.GetBuildState(jksId, jobName, buildId)
	if err != nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": err.Error()}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "data": &bs}
	this.ServeJSON()
}
