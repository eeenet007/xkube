package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"xkube/common"
	m "xkube/models"

	"github.com/tidwall/gjson"

	beego "github.com/beego/beego/v2/server/web"
)

type IngressController struct {
	beego.Controller
}

func (this *IngressController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	serviceName := this.GetString("serviceName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		ingressName = gp.Get("ingressName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}
	//get from redis
	if clusterId != "" && ingressName == "" && labels == "" {
		resp := common.Get("ingressList" + clusterId + nameSpace)
		if resp != "" {
			var ingressList = make([]m.Ingress, 0)
			err := json.Unmarshal([]byte(resp), &ingressList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal ingressList Error:%s\n", err)
			} else {
				count := len(ingressList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &ingressList}
				this.ServeJSON()
				return
			}
		}
	}

	ingressList, err := m.IngList(clusterId, nameSpace, ingressName, serviceName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(ingressList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && ingressName == "" && labels == "" {
		bodystr, err := json.Marshal(&ingressList)
		if err == nil {
			_ = common.SetEx("ingressList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": msg, "count": count, "data": &ingressList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *IngressController) Create() {
	clusterId := this.GetString("clusterId")
	msg := "success"
	code := 0
	err := m.IngCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *IngressController) CreateByYaml() {
	clusterId := this.GetString("clusterId")
	err := m.IngYamlCreate(clusterId, this.Ctx.Input.RequestBody)
	msg := "success"
	code := 0
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *IngressController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	//nameSpace := gp.Get("nameSpace").String()
	//log.Println(string(this.Ctx.Input.RequestBody))

	//backup
	yamlStr, _ := m.GetIngYaml(clusterId, nameSpace, ingressName)
	_ = m.InsertBackup(clusterId, nameSpace, ingressName, "ingress", yamlStr, "Backup before updating")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	err := m.IngYamlModify(clusterId, []byte(reqBody))
	code := 0
	msg := "success"
	if err != nil {
		log.Printf("[WARN] IngYamlModify Fail:%s\n", err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *IngressController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	xdetail, err := m.IngDetail(clusterId, nameSpace, ingressName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *IngressController) Del() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	//backup
	yamlStr, _ := m.GetIngYaml(clusterId, nameSpace, ingressName)
	_ = m.InsertBackup(clusterId, nameSpace, ingressName, "ingress", yamlStr, "Backup before updating")

	err := m.IngDelete(clusterId, nameSpace, ingressName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success"}
	this.ServeJSON()
}

func (this *IngressController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")

	yamlStr, _ := m.GetIngYaml(clusterId, namespace, ingressName)
	this.Ctx.WriteString(yamlStr)
}

func (this *IngressController) GetRule() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	xList, err := m.GetIngRule(clusterId, nameSpace, ingressName)
	code := 0
	msg := "success"
	if err != nil {
		log.Println(err)
		code = -1
		msg = "fail"
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": len(xList), "data": &xList}
	this.ServeJSON()
}

func (this *IngressController) UpdateRule() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	//backup
	yamlStr, _ := m.GetIngYaml(clusterId, nameSpace, ingressName)
	_ = m.InsertBackup(clusterId, nameSpace, ingressName, "ingress", yamlStr, "Backup before updating")

	//log.Println(string(this.Ctx.Input.RequestBody))
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	// type IngHostPath struct {
	// 	Path        string `json:"path"`
	// 	PathType    string `json:"pathType"`
	// 	ServiceName string `json:"serviceName"`
	// 	ServicePort int32  `json:"servicePort"`
	// }

	var ingRuleArry []m.IngressRule
	for i := 1; i <= 40; i++ {
		hostStr := gp.Get("host_" + fmt.Sprintf("%d", i)).Str
		if hostStr == "" {
			break
		}
		//log.Println(hostStr)
		var rPath = make([]m.IngHostPath, 0)
		for j := 1; j <= 40; j++ {
			pathStr := gp.Get("path_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			pathTypeStr := gp.Get("pathType_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			serviceNameStr := gp.Get("serviceName_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			servicePortStr := gp.Get("servicePort_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			//var svcPort int32 = 0
			//var err1 error
			//log.Println(pathStr, pathTypeStr, serviceNameStr, servicePortStr)
			if pathStr != "" && pathTypeStr != "" && serviceNameStr != "" {
				var svcPort int64 = 80
				var err error
				if servicePortStr != "" {
					svcPort, err = strconv.ParseInt(servicePortStr, 10, 32)
					if err != nil {
						log.Printf("[WARN] UpdateRule ParseInt servicePort Fail:%s\n", err.Error())
						this.Data["json"] = &map[string]interface{}{"code": -1, "msg": err.Error()}
						this.ServeJSON()
						return
					}
				}
				rPath = append(rPath, m.IngHostPath{
					Path:        pathStr,
					PathType:    pathTypeStr,
					ServiceName: serviceNameStr,
					ServicePort: int32(svcPort),
				})
			}
		}
		if len(rPath) > 0 {
			ingRuleArry = append(ingRuleArry, m.IngressRule{
				Host:  hostStr,
				Paths: rPath,
			})
		}
	}
	//log.Println(ingRuleArry)
	err := m.UpdateIngRule(clusterId, nameSpace, ingressName, ingRuleArry)
	code := 0
	msg := "success"
	if err != nil {
		log.Printf("[WARN] UpdateIngRule Fail:%s\n", err.Error())
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *IngressController) GetTlsHost() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	xList, err := m.GetIngTlsHost(clusterId, nameSpace, ingressName)
	code := 0
	msg := "success"
	if err != nil {
		log.Println(err)
		code = -1
		msg = "fail"
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": len(xList), "data": &xList}
	this.ServeJSON()
}

func (this *IngressController) UpdateTlsHost() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	ingressName := this.GetString("ingressName")
	//backup
	yamlStr, _ := m.GetIngYaml(clusterId, nameSpace, ingressName)
	_ = m.InsertBackup(clusterId, nameSpace, ingressName, "ingress", yamlStr, "Backup before updating")

	//log.Println(string(this.Ctx.Input.RequestBody))
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	var tlsHostArry []m.IngressTlsHost
	for i := 0; i <= 30; i++ {
		hostStr := gp.Get("hosts[" + fmt.Sprintf("%d", i) + "]").Str
		var hosts = []string{}
		for _, vv := range strings.Split(hostStr, "\n") {
			if len(strings.TrimSpace(vv)) > 0 {
				hosts = append(hosts, strings.TrimSpace(vv))
			}
		}
		tlsName := gp.Get("tlsName[" + fmt.Sprintf("%d", i) + "]").Str
		if tlsName != "" && len(hosts) > 0 {
			tlsHostArry = append(tlsHostArry, m.IngressTlsHost{
				Hosts: hosts,
				Tls:   tlsName,
			})
		} else {
			break
		}
	}

	err := m.UpdateIngTlsHost(clusterId, nameSpace, ingressName, tlsHostArry)
	code := 0
	msg := "success"
	if err != nil {
		log.Printf("[WARN] UpdateIngTlsHost Fail:%s\n", err.Error())
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
