package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type PodController struct {
	beego.Controller
}

func (this *PodController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")

	resName := this.GetString("resName")
	// deployName := this.GetString("deployName")
	// statefulsetName := this.GetString("statefulsetName")
	// daemonsetName := this.GetString("daemonsetName")
	// if deployName == "" && statefulsetName != "" {
	// 	deployName = statefulsetName
	// }

	resType := this.GetString("resType")
	podName := this.GetString("podName")
	nodeName := this.GetString("nodeName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		podName = gp.Get("podName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":") //yaml里没有对应的appname,deployment不会传递过来
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && podName == "" && resName == "" && labels == "" && nodeName == "" {
		resp := common.Get("podList" + clusterId + nameSpace)
		if resp != "" {
			var podList = make([]m.Podinfo, 0)
			err := json.Unmarshal([]byte(resp), &podList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal podList Error:%s\n", err)
			} else {
				count := len(podList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &podList}
				this.ServeJSON()
				return
			}
		}
	}

	//podList, err := m.PodList(clusterId, nameSpace, deployName, podName, labelsKey, labelsValue, nodeName)
	podList, err := m.PodListV2(clusterId, nameSpace, resName, resType, podName, labelsKey, labelsValue, nodeName)
	msg := "success"
	code := 0
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}

	count := len(podList)
	//set in redis
	//if code == 0 && podName == "" && deployName == "" && labels == "" && nodeName == "" {
	if code == 0 && podName == "" && resName == "" && labels == "" && nodeName == "" {
		//set数量到redis
		if nameSpace == "" {
			_ = common.HSet("count_"+clusterId, "pod", fmt.Sprintf("%d", count))
		}
		bodystr, err := json.Marshal(&podList)
		if err == nil {
			_ = common.SetEx("podList"+clusterId+nameSpace, string(bodystr), 900)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &podList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *PodController) Detail() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	podName := this.GetString("podName")
	xdetail, err := m.PodDetail(clusterId, nameSpace, podName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *PodController) ContainerList() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	podName := this.GetString("podName")
	code := 0
	msg := "success"
	xList, err := m.PodContainerList(clusterId, nameSpace, podName)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] ContainerList error:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": len(xList), "data": &xList}
	this.ServeJSON()
}

func (this *PodController) Log() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	podName := this.GetString("podName")
	download, _ := this.GetBool("download")
	container := this.GetString("container")
	logLine, _ := this.GetInt64("logLine")
	encode := this.GetString("encode")

	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		podName = gp.Get("podName").String()
		logLine = gp.Get("logLine").Int()
		container = gp.Get("container").String()
		encode = gp.Get("encode").String()
	}
	if logLine == 0 { //没设置行时，默认为100
		logLine = 100
	}
	if download { //当日志时下载时，不限制行数
		logLine = 0
	}
	log := m.PodLog(clusterId, nameSpace, podName, container, logLine, encode)
	this.Ctx.WriteString(log)
}

func (this *PodController) Del() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	podName := this.GetString("podName")
	err := m.PodDel(clusterId, nameSpace, podName)
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

func (this *PodController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	podName := this.GetString("podName")

	yamlStr, _ := m.GetPodYaml(clusterId, namespace, podName)
	this.Ctx.WriteString(yamlStr)
}

func (this *PodController) Exec() {
	if this.Ctx.Request.Method == "POST" {
		clusterId := this.GetString("clusterId")
		nameSpace := this.GetString("nameSpace")
		podName := this.GetString("podName")
		containerName := this.GetString("containerName")
		command := string(this.Ctx.Input.RequestBody)
		stdout, stderr, err := m.PodExec(clusterId, nameSpace, podName, containerName, command, 120)
		if err != nil {
			log.Printf("[ERROR] PodExec Fail:%v\n", err)
			this.Ctx.WriteString("podExecFail")
			return
		}
		bodyStr := fmt.Sprintf("stdout:\n%s\nstderr:\n%s\n", stdout, stderr)
		this.Ctx.WriteString(bodyStr)
	} else {
		this.Ctx.WriteString("noSupportMethod")
	}
}

func (this *PodController) Check() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	podName := this.GetString("podName")
	podIp := this.GetString("podIp")
	containerName := this.GetString("containerName")
	code := 0
	msg := "success"
	//log.Println(string(this.Ctx.Input.RequestBody))
	if podIp == "" {
		podDetail, err := m.PodDetail(clusterId, nameSpace, podName)
		if err != nil {
			code = -1
			msg = "podIp is Null"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}
		podIp = podDetail.PodIp
	}

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/event-stream")

	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	checkType := gp.Get("checkType").String()
	if containerName == "" {
		containerName = gp.Get("containerName").String()
	}

	var CmdStr string
	if checkType == "http" {
		httpUrl := gp.Get("httpUrl").String()
		httpPort := gp.Get("httpPort").String()
		httpMethod := gp.Get("httpMethod").String()
		printHeader := gp.Get("printHeader").Bool()

		if httpPort == "" {
			httpPort = "80"
		}
		if httpUrl == "" {
			httpUrl = "/"
			// code = -1
			// msg = "url is Null"
			// this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			// this.ServeJSON()
			// return
		}

		var headerStr string
		httpHeaderKey := gp.Get("httpHeaderKey").String()
		httpHeaderValue := gp.Get("httpHeaderValue").String()
		if httpHeaderKey != "" && httpHeaderValue != "" {
			headerStr += fmt.Sprintf("-H '%s:%v' ", httpHeaderKey, httpHeaderValue)
		}
		for i := 0; i <= 10; i++ {
			//log.Println(i)
			kk := gp.Get("httpHeaderKey[" + fmt.Sprintf("%d", i) + "]").Str
			vv := gp.Get("httpHeaderValue[" + fmt.Sprintf("%d", i) + "]").Str
			//log.Println(kk, vv)
			if kk != "" && vv != "" {
				//log.Println(kk, vv)
				headerStr += fmt.Sprintf("-H '%s:%v' ", kk, vv)
			} else {
				break
			}
		}
		var bodyStr string
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
		}
		var printHeaderStr string
		if printHeader {
			printHeaderStr = "-v"
		}
		var methodStr string
		if httpMethod == "POST" && bodyStr == "" {
			methodStr = "-X POST"
		}

		cmdPrefix := fmt.Sprintf("curl %s %s %s %s --connect-timeout 3 -m 3", methodStr, printHeaderStr, headerStr, bodyStr)
		CmdStr = fmt.Sprintf("%s 'http://%s:%s%s'", cmdPrefix, podIp, httpPort, httpUrl)
		log.Println(CmdStr)
		var respBody string
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", podIp, podName)
		if printHeader {
			respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		} else {
			respBody += fmt.Sprintf("%s", stdout)
		}
		//respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s\n", stdout, stderr)
		//log.Println(respBody)
		//this.Ctx.WriteString(respBody)
		this.Ctx.ResponseWriter.Write([]byte(respBody))
		this.Ctx.ResponseWriter.Flush()

	} else if checkType == "tcp" {
		tcpPort := gp.Get("tcpPort").String()
		tcpCmd := gp.Get("tcpCmd").String()
		if tcpPort == "" && tcpCmd == "" {
			code = -1
			msg = "ip port null"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}

		if tcpCmd == "" {
			CmdStr = fmt.Sprintf("nmap -sT %s -p %s", podIp, tcpPort)
		} else {
			CmdStr = strings.Replace(tcpCmd, "$IP", podIp, -1)
			//CmdStr = strings.Replace(tcpCmd, "$PORT", tcpPort, -1)
		}
		log.Println(CmdStr)
		var respBody string
		stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
		respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", podIp, podName)
		respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.ResponseWriter.Write([]byte(respBody))
		this.Ctx.ResponseWriter.Flush()

	} else if checkType == "ping" {
		pingIp := gp.Get("pingIp").String()
		pingPacketNum := gp.Get("pingPacketNum").Int()
		if pingIp == "" {
			code = -1
			msg = "ip null"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}
		if pingPacketNum == 0 && pingPacketNum > 100 {
			pingPacketNum = 2
		}

		CmdStr = fmt.Sprintf("ping -c %d %s", pingPacketNum, pingIp)
		//log.Println(CmdStr)
		var respBody string
		stdout, stderr, _ := m.PodExec(clusterId, nameSpace, podName, containerName, CmdStr, 300)
		respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", podIp, podName)
		respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.ResponseWriter.Write([]byte(respBody))
		this.Ctx.ResponseWriter.Flush()
	} else if checkType == "dns" {
		dnsDomain := gp.Get("dnsDomain").String()
		dnsIp := gp.Get("dnsIp").String()
		if dnsDomain == "" {
			code = -1
			msg = "domain null"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}

		CmdStr = fmt.Sprintf("nslookup %s %s", dnsDomain, dnsIp)
		log.Println(CmdStr)
		var respBody string
		stdout, stderr, _ := m.PodExec(clusterId, nameSpace, podName, containerName, CmdStr, 300)
		respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", podIp, podName)
		respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.ResponseWriter.Write([]byte(respBody))
		this.Ctx.ResponseWriter.Flush()

	} else if checkType == "shell" {
		shellTxt := gp.Get("shellTxt").String()
		if shellTxt == "" {
			code = -1
			msg = "shell null"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(shellTxt)
		if err != nil {
			log.Println("[WARN] base64.StdEncoding.Decode err:", err)
			code = -1
			msg = "shell base64 Decode err"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}
		CmdStr = string(decoded)
		log.Println(CmdStr)
		var respBody string
		stdout, stderr, _ := m.PodExec(clusterId, nameSpace, podName, containerName, CmdStr, 300)
		respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", podIp, podName)
		respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
		this.Ctx.ResponseWriter.Write([]byte(respBody))
		this.Ctx.ResponseWriter.Flush()

	} else {
		log.Println("noSupport")
		code = -1
		msg = "noSupport"
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
}

func (this *PodController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	podName := this.GetString("podName")

	//backup
	yamlStr, _ := m.GetPodYaml(clusterId, namespace, podName)
	_ = m.InsertBackup(clusterId, namespace, podName, "pod", yamlStr, "ModifyByYaml")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	msg := "success"
	code := 0
	err := m.PodYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] PodYamlModify Fail:%s\n", err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}
