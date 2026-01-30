package controllers

import (
	"encoding/json"
	"log"
	m "xkube/models"

	//"fmt"
	"strings"
	"xkube/common"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type JobController struct {
	beego.Controller
}

func (this *JobController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	jobName := this.GetString("jobName")
	cronjobName := this.GetString("cronjobName")

	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		jobName = gp.Get("jobName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && cronjobName == "" && jobName == "" && labels == "" {
		resp := common.Get("jobList" + clusterId + nameSpace)
		if resp != "" {
			var jobList = make([]m.Job, 0)
			err := json.Unmarshal([]byte(resp), &jobList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal jobList Error:%s\n", err)
			} else {
				count := len(jobList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &jobList}
				this.ServeJSON()
				return
			}
		}
	}

	//xList, err := m.JobList(clusterId, nameSpace, jobName, labelsKey, labelsValue)
	xList, err := m.JobListv2(clusterId, nameSpace, cronjobName, jobName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && cronjobName == "" && jobName == "" && labels == "" {
		bodystr, err := json.Marshal(&xList)
		if err == nil {
			_ = common.SetEx("jobList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *JobController) Detail() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	jobName := this.GetString("jobName")
	xdetail, err := m.JobDetail(clusterId, namespace, jobName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *JobController) Create() {
	clusterId := this.GetString("clusterId")
	var job m.NewJob

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &job)
	err = m.JobCreate(clusterId, &job)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *JobController) Modify() {
	clusterId := this.GetString("clusterId")
	var sts m.Job

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sts)
	err = m.JobModify(clusterId, &sts)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *JobController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	jobName := this.GetString("jobName")
	err := m.JobDel(clusterId, namespace, jobName)
	code := 0
	msg := "success"
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *JobController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	jobName := this.GetString("jobName")

	yamlStr, _ := m.GetJobYaml(clusterId, namespace, jobName)
	this.Ctx.WriteString(yamlStr)
	//this.Data["yaml"] = &yamlStr
	//this.ServeYAML()
	//this.ServeJSON()
}

func (this *JobController) Log() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	resName := this.GetString("resName")
	encode := this.GetString("encode")
	logstr := m.JobLog(clusterId, nameSpace, resName, encode)
	this.Ctx.WriteString(logstr)
}

func (this *JobController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	jobName := this.GetString("jobName")

	//backup
	yamlStr, _ := m.GetJobYaml(clusterId, namespace, jobName)
	_ = m.InsertBackup(clusterId, namespace, jobName, "job", yamlStr, "ModifyByYaml")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	code := 0
	msg := "success"
	err := m.JobYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] JobYamlModify Fail:%s\n", err)
		code = 1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
