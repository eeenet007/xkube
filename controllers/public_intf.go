// cluster.go
package controllers

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	//"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type PublicIntfController struct {
	beego.Controller
}

// 列表
// "cmList" + clusterId + nameSpace)
// "deployList" + clusterId + nameSpace
// "svcList" + clusterId + nameSpace
// "podList" + clusterId + nameSpace
// "ingressList"+clusterId+nameSpace
// "dsList" + clusterId + nameSpace
// "jobList" + clusterId + nameSpace
// "cronjobList" + clusterId + nameSpace
func (this *PublicIntfController) Clear() {

	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin"}
		this.ServeJSON()
		return
	}
	clusterId := this.Ctx.GetCookie("clusterId")
	num := m.ClearCache(clusterId)
	msg := fmt.Sprintf("%s:cache clear success:%d", clusterId, num)
	this.Data["json"] = &map[string]interface{}{"code": 1, "msg": msg}
	this.ServeJSON()
}

func (this *PublicIntfController) UpdateImage() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployType := this.GetString("deployType")
	deployName := this.GetString("deployName")

	authKey := this.Ctx.Request.Header.Get("xkey")
	apiKey, _ := beego.AppConfig.String("apiKey")

	if authKey != apiKey {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "errorKey"}
		this.ServeJSON()
		return
	}
	//userAgent := this.Ctx.Input.UserAgent()

	if this.Ctx.Input.Method() == "POST" {
		log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		containerName := gp.Get("containerName").Str
		image := gp.Get("image").Str
		deployName = gp.Get("deployName").Str
		msg := "success"
		code := 0
		if deployType == "deployment" {
			err := m.DeploymentUpdateImage(clusterId, nameSpace, deployName, containerName, image)
			if err != nil {
				msg = err.Error()
				code = -1
			}
		} else {
			// else if deployType == "statefulset" {
			// _, err := m.StatefulsetUpdateImage(clusterId, nameSpace, deployName, containerName)
			// if err != nil {
			// msg = err.Error()
			// code = -1
			// }
			// } else if deployType == "cronjob" {
			// _, err := m.CronjobUpdateImage(clusterId, nameSpace, deployName, containerName)
			// if err != nil {
			// msg = err.Error()
			// code = -1
			// }
			// }
			msg = fmt.Sprintf("UnSupport:%s", deployType)
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "needPOST"}
	this.ServeJSON()
}

// for app
func (this *PublicIntfController) Check() {

	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin"}
		this.ServeJSON()
		return
	}

	clusterId := this.GetString("clusterId")
	checkType := this.GetString("checkType")

	//this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	var CmdStr string
	if checkType == "ping" {
		pingIp := gp.Get("pingIp").String()
		if pingIp == "" {
			msg := "ip null"
			this.Ctx.WriteString(msg)
			return
		}
		CmdStr = fmt.Sprintf("ping -c 4 %s", pingIp)
		log.Println(CmdStr)
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.WriteString(respBody)
	} else if checkType == "tcp" {
		ipPort := gp.Get("ipPort").String()
		if ipPort == "" {
			msg := "ip port null"
			this.Ctx.WriteString(msg)
			return
		}
		ipstr := strings.Split(ipPort, ":")
		if len(ipstr) != 2 {
			msg := "ip:port error"
			this.Ctx.WriteString(msg)
			return
		}
		CmdStr = fmt.Sprintf("nmap -sT %s -p %s", ipstr[0], ipstr[1])
		log.Println(CmdStr)
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.ResponseWriter.Write([]byte(respBody))
		this.Ctx.ResponseWriter.Flush()
	} else if checkType == "dns" {
		domain := gp.Get("domain").String()
		dnsIp := gp.Get("dnsIp").String()
		if domain == "" {
			msg := "domain null"
			this.Ctx.WriteString(msg)
			return
		}

		CmdStr = fmt.Sprintf("nslookup %s %s", domain, dnsIp)
		log.Println(CmdStr)
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.WriteString(respBody)
	} else if checkType == "cmd" {
		command := gp.Get("command").String()
		if command == "" {
			msg := "command null"
			this.Ctx.WriteString(msg)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(command)
		if err != nil {
			msg := "shell base64 Decode err"
			this.Ctx.WriteString(msg)
			return
		}
		CmdStr = string(decoded)
		log.Println(CmdStr)
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.WriteString(respBody)
	} else if checkType == "http" {
		httpUrl := gp.Get("httpUrl").String()
		httpMethod := gp.Get("httpMethod").String()

		var bodyStr, headerStr string
		if httpMethod == "POST" {
			decoded, err := base64.StdEncoding.DecodeString(gp.Get("httpReqBody").String())
			if err != nil {
				log.Println("[WARN] base64.StdEncoding.Decode err:", err)
			} else {
				bodyStr = string(decoded)
			}
			if bodyStr != "" {
				bodyStr = fmt.Sprintf("-d '%s' ", bodyStr)
			}
			headerStr = "-H 'Content-Type:application/json'"
		}

		var methodStr string
		if httpMethod == "POST" && bodyStr == "" {
			methodStr = "-X POST"
		}
		CmdStr := fmt.Sprintf("curl %s -v %s %s --connect-timeout 3 -m 5 '%s'", methodStr, headerStr, bodyStr, httpUrl)
		log.Println(CmdStr)
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody := fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.WriteString(respBody)
	} else {
		this.Ctx.WriteString("noSupport")
	}
}

// 通过接口调用来进行部署定时重启服务
func (this *PublicIntfController) ScheduleRestart() {

	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	resName := this.GetString("resName")
	resType := this.GetString("resType")

	authKey := this.Ctx.Request.Header.Get("xkey")
	apiKey, _ := beego.AppConfig.String("apiKey")
	if authKey != apiKey {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "errorKey"}
		this.ServeJSON()
		return
	}

	code := 0
	msg := "success"

	switch resType {
	case "deploy", "deployment":
		err := m.DeployRestart(clusterId, namespace, resName)
		if err != nil {
			log.Println(err)
			code = -1
			msg = err.Error()
		}
	case "sts", "statefulset":
		err := m.StatefulsetRestart(clusterId, namespace, resName)
		if err != nil {
			log.Println(err)
			code = -1
			msg = err.Error()
		}
	default:
		code = -1
		msg = "noSupport"
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
