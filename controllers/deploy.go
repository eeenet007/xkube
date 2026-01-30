package controllers

import (
	//"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	//"net/http"
	"strings"
	"time"
	"xkube/common"
	m "xkube/models"
	adm "xkube/xadmin/src/models"

	"github.com/tidwall/gjson"

	beego "github.com/beego/beego/v2/server/web"
)

type DeployController struct {
	beego.Controller
}

func (this *DeployController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		deployName = gp.Get("deployName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	//get from redis
	if clusterId != "" && deployName == "" && labels == "" {
		resp := common.Get("deployList" + clusterId + nameSpace)
		if resp != "" {
			var depList = make([]m.Deploy, 0)
			err := json.Unmarshal([]byte(resp), &depList)
			if err != nil {
				log.Printf("[ERROR] Unmarshal depList Error:%s\n", err)
			} else {
				count := len(depList)
				this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &depList}
				this.ServeJSON()
				return
			}
		}
	}

	depList, err := m.DeployList(clusterId, nameSpace, deployName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(depList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}

	//set in redis
	if code == 0 && deployName == "" && labels == "" {
		//set数量到redis
		if nameSpace == "" {
			_ = common.HSet("count_"+clusterId, "deploy", fmt.Sprintf("%d", count))
		}

		bodystr, err := json.Marshal(&depList)
		if err == nil {
			_ = common.SetEx("deployList"+clusterId+nameSpace, string(bodystr), 600)
		}
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &depList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *DeployController) Detail() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	//deploy := this.GetString("deploy")
	deployName := this.GetString("deployName")
	xdetail, err := m.DeployDetail(clusterId, namespace, deployName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *DeployController) Create() {
	clusterId := this.GetString("clusterId")
	code := 0
	msg := "success"

	deployName, err := m.DeployCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] Deploy Create Fail:%s\n", err)
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "create", fmt.Sprintf("%d", code), "createDeployment:"+deployName+",result:"+msg)

	//clear cache
	_ = m.ClearCache(clusterId) //创建新的deploy以后，刷新一下列表缓存，要不然不会马上再列表中显示

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, namespace, deployName)
	_ = m.InsertBackup(clusterId, namespace, deployName, "deployment", yamlStr, "ModifyByYaml")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")

	code := 0
	msg := "success"
	err := m.DeployYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] DeployYamlModify Fail:%s\n", err)
		code = -1
		msg = err.Error()
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateYaml:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Modify() {
	clusterId := this.GetString("clusterId")
	var deploy m.Deploy
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &deploy)
	err = m.DeployModify(clusterId, &deploy)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

// 迁移到其他集群
func (this *DeployController) Clone() {
	//clusterId := this.GetString("clusterId")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	clusterId := gp.Get("clusterid").String()
	namespace := gp.Get("namespace").String()
	target_clusterid := gp.Get("target_clusterid").String()
	target_namespace := gp.Get("target_namespace").String()
	target_objname := gp.Get("target_objname").String()
	objName := gp.Get("objname").String()
	deployArry := strings.Split(objName, ",")
	targetDeployArry := strings.Split(target_objname, ",")
	var result string
	var isOk = true

	for i, vv := range deployArry {
		deployName := vv
		var targetObjname string
		if len(targetDeployArry)-1 >= i {
			targetObjname = targetDeployArry[i]
		}
		err := m.DeployClone(clusterId, namespace, deployName, target_clusterid, target_namespace, targetObjname)
		if err != nil {
			isOk = false
		} else {
			err = fmt.Errorf("ok")
		}
		result += fmt.Sprintf("%s result:%s,", deployName, err)
	}
	uinfo := this.GetSession("userinfo")
	u := m.Xkb_clone_log{
		Clusterid:       clusterId,
		Namespace:       namespace,
		Restype:         "deployment",
		Objname:         objName,
		TargetClusterid: target_clusterid,
		TargetNamespace: target_namespace,
		TargetObjname:   target_objname,
		Status:          fmt.Sprintf("%v", isOk),
		User:            uinfo.(adm.User).Username,
		Result:          result,
	}
	_, err := m.Add_Clone(&u) //插入clone 结果
	if err != nil {
		log.Printf("[ERROR] Clone add Fail:%s\n", err)
		result += fmt.Sprintf("insert log fail")
	}

	this.Data["json"] = &map[string]interface{}{"code": 0, "success": isOk, "msg": result}
	this.ServeJSON()
}

func (this *DeployController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, namespace, deployName)
	_ = m.InsertBackup(clusterId, namespace, deployName, "deployment", yamlStr, "Del")

	err := m.DeployDel(clusterId, namespace, deployName)
	code := 0
	msg := "success"
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "delete", fmt.Sprintf("%d", code), "deleteDeployment:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	yamlStr, _ := m.GetDeployYaml(clusterId, namespace, deployName)
	this.Ctx.WriteString(yamlStr)
	//this.Data["yaml"] = &yamlStr
	//this.ServeYAML()
	//this.ServeJSON()
}

func (this *DeployController) Labels() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	var labelsMap = make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {

		//backup
		yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
		_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Labels")

		//log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for i := 0; i <= 40; i++ {
			//log.Println(i)
			kk := gp.Get("labels_key[" + fmt.Sprintf("%d", i) + "]").Str
			vv := gp.Get("labels_value[" + fmt.Sprintf("%d", i) + "]").Str
			if kk != "" && vv != "" {
				//log.Println(kk, vv)
				labelsMap[kk] = vv
			} else {
				break
			}
		}
		msg := "success"
		code := 0
		_, err := m.DeployLabels(clusterId, nameSpace, deployName, "POST", labelsMap)
		if err != nil {
			msg = err.Error()
			code = -1
		}

		//auditlog
		userinfo := this.GetSession("userinfo")
		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}
		_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateLables deploy:"+deployName+",result:"+msg)

		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	xList, err := m.DeployLabels(clusterId, nameSpace, deployName, "GET", labelsMap)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) Image() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	var bbb = make([]m.ImageKv, 0)
	if this.Ctx.Input.Method() == "POST" {
		//backup
		yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
		_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Image")

		//log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for i := 0; i <= 10; i++ {
			cn := gp.Get("containerName[" + fmt.Sprintf("%d", i) + "]").Str
			ci := gp.Get("containerId[" + fmt.Sprintf("%d", i) + "]").Int()
			img := gp.Get("image[" + fmt.Sprintf("%d", i) + "]").Str
			//log.Printf("[DEBUG]%s,%d,%s\n", cn, ci, img)
			if cn != "" && ci >= 0 && img != "" {
				bbb = append(bbb, *&m.ImageKv{
					ContainerName: cn,
					ContainerId:   int(ci),
					Image:         img,
				})
			} else {
				break
			}
		}
		msg := "success"
		code := 0
		log.Println(bbb)
		_, err := m.DeployImage(clusterId, nameSpace, deployName, "POST", bbb)
		if err != nil {
			msg = err.Error()
			code = -1
		}

		//auditlog
		userinfo := this.GetSession("userinfo")
		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}
		_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateImage deploy:"+deployName+",result:"+msg)

		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	xList, err := m.DeployImage(clusterId, nameSpace, deployName, "GET", bbb)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) Replicas() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	if this.Ctx.Input.Method() == "POST" {
		//log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		podNum := gp.Get("podNumber").Int()
		msg := "success"
		code := 0
		_, err := m.DeployReplicas(clusterId, nameSpace, deployName, "POST", int32(podNum))
		if err != nil {
			msg = err.Error()
			code = -1
		}

		//auditlog
		userinfo := this.GetSession("userinfo")
		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}
		_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateReplicas deploy:"+deployName+",result:"+msg)

		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	resp, err := m.DeployReplicas(clusterId, nameSpace, deployName, "GET", 0)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "number": resp}
	this.ServeJSON()
}

func (this *DeployController) Strategy() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	var sst m.StrategyST
	if this.Ctx.Input.Method() == "POST" {

		//backup
		yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
		_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Strategy")

		msg := "success"
		code := 0
		//log.Println(string(this.Ctx.Input.RequestBody))
		//gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		err := json.Unmarshal(this.Ctx.Input.RequestBody, &sst)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		_, err = m.DeployStrategy(clusterId, nameSpace, deployName, "POST", sst)
		if err != nil {
			msg = err.Error()
			code = -1
		}

		//auditlog
		userinfo := this.GetSession("userinfo")
		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}
		_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateStrategy deploy:"+deployName+",result:"+msg)

		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	resp, err := m.DeployStrategy(clusterId, nameSpace, deployName, "GET", sst)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": resp}
	this.ServeJSON()
}

func (this *DeployController) ReplicasetYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	replicasetName := this.GetString("replicasetName")

	yamlStr, _ := m.GetReplicasetYaml(clusterId, namespace, replicasetName)
	this.Ctx.WriteString(yamlStr)
}

func (this *DeployController) ReplicasetList() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	labels := this.GetString("labels")
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}
	msg := "success"
	code := 0
	xList, err := m.GetReplicasetList(clusterId, namespace, labelsKey, labelsValue)
	count := len(xList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) RollBack() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	replicasetName := this.GetString("replicasetName")

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, namespace, deployName)
	_ = m.InsertBackup(clusterId, namespace, deployName, "deployment", yamlStr, "RollBack")

	code := 0
	msg := "success"
	err := m.DeployRollBack(clusterId, namespace, deployName, replicasetName)
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "rollback", fmt.Sprintf("%d", code), "RollBack deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Restart() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	code := 0
	msg := "success"
	err := m.DeployRestart(clusterId, namespace, deployName)
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "restart", fmt.Sprintf("%d", code), "Restart deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Check() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/event-stream")
	log.Println(string(this.Ctx.Input.RequestBody))
	code := 0
	msg := "success"

	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	checkType := gp.Get("checkType").String()
	podList, _ := m.PodListV2(clusterId, nameSpace, deployName, "deploy", "", "", "", "")

	if len(podList) == 0 {
		code = -1
		msg = "get Pod ip fail"
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}

	maxPodNum := gp.Get("maxPodNum").Int()
	if maxPodNum == 0 {
		maxPodNum = 10
	}

	var CmdStr string
	if checkType == "http" {
		httpUrl := gp.Get("httpUrl").String()
		httpPort := gp.Get("httpPort").String()
		httpMethod := gp.Get("httpMethod").String()
		printHeader := gp.Get("printHeader").Bool()

		if httpUrl == "" {
			//httpUrl = m.DeployAnnotations(clusterId, nameSpace, deployName, "url")
			httpUrl = "/"
		}
		if httpPort == "" {
			httpPort = "80"
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
			log.Println(kk, vv)
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
		//log.Println(cmdPrefix)
		//this.Ctx.WriteString("ok")
		//return

		for i, vv := range podList {
			ip := vv.PodIp
			CmdStr = fmt.Sprintf("%s 'http://%s:%s%s'", cmdPrefix, ip, httpPort, httpUrl)
			log.Println(CmdStr)
			var respBody string
			stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
			respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", ip, vv.PodName)
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
			if i == int(maxPodNum)-1 {
				break
			}
		}
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
		for i, vv := range podList {
			ip := vv.PodIp
			if tcpCmd == "" {
				CmdStr = fmt.Sprintf("nmap -sT %s -p %s", ip, tcpPort)
			} else {
				CmdStr = strings.Replace(tcpCmd, "$IP", ip, -1)
				//CmdStr = strings.Replace(tcpCmd, "$PORT", tcpPort, -1)
			}
			log.Println(CmdStr)
			var respBody string
			stdout, stderr, _ := m.PodExec(clusterId, "xkube", "xkube-check", "xkube-check", CmdStr, 300)
			respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", ip, vv.PodName)
			respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
			this.Ctx.ResponseWriter.Write([]byte(respBody))
			this.Ctx.ResponseWriter.Flush()
			if i == int(maxPodNum)-1 {
				break
			}
		}
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
		for i, vv := range podList {
			CmdStr = fmt.Sprintf("ping -c %d %s", pingPacketNum, pingIp)
			//log.Println(CmdStr)
			var respBody string
			xList, err := m.PodContainerList(clusterId, nameSpace, vv.PodName)
			if err != nil || len(xList) == 0 {
				log.Printf("[ERROR] Get ContainerList error:%s\n", err)
				respBody = fmt.Sprintf("%s Get ContainerList error:%s\n", vv.PodName, err)
				this.Ctx.ResponseWriter.Write([]byte("Get ContainerList error"))
				this.Ctx.ResponseWriter.Flush()
				continue
			}
			containerName := xList[0].ContainerName
			stdout, stderr, _ := m.PodExec(clusterId, vv.NameSpace, vv.PodName, containerName, CmdStr, 300)
			respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", vv.PodIp, vv.PodName)
			respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
			this.Ctx.ResponseWriter.Write([]byte(respBody))
			this.Ctx.ResponseWriter.Flush()
			if i == int(maxPodNum)-1 {
				break
			}
		}
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

		for i, vv := range podList {
			CmdStr = fmt.Sprintf("nslookup %s %s", dnsDomain, dnsIp)
			//log.Println(CmdStr)
			var respBody string
			xList, err := m.PodContainerList(clusterId, nameSpace, vv.PodName)
			if err != nil || len(xList) == 0 {
				log.Printf("[ERROR] Get ContainerList error:%s\n", err)
				respBody = fmt.Sprintf("%s Get ContainerList error:%s\n", vv.PodName, err)
				this.Ctx.ResponseWriter.Write([]byte("Get ContainerList error"))
				this.Ctx.ResponseWriter.Flush()
				continue
			}
			containerName := xList[0].ContainerName
			stdout, stderr, _ := m.PodExec(clusterId, vv.NameSpace, vv.PodName, containerName, CmdStr, 300)
			respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", vv.PodIp, vv.PodName)
			respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
			this.Ctx.ResponseWriter.Write([]byte(respBody))
			this.Ctx.ResponseWriter.Flush()
			if i == int(maxPodNum)-1 {
				break
			}
		}
	} else if checkType == "shell" {
		shellTxt := gp.Get("shellTxt").String()
		if shellTxt == "" {
			code = -1
			msg = "shell null"
			this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
			this.ServeJSON()
			return
		}

		for i, vv := range podList {
			decoded, err := base64.StdEncoding.DecodeString(shellTxt)
			if err != nil {
				log.Println("[WARN] base64.StdEncoding.Decode err:", err)
				code = -1
				msg = "shell base64 Decode err"
				this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
				this.ServeJSON()
				break
			}
			CmdStr = string(decoded)
			log.Println(CmdStr)
			var respBody string
			xList, err := m.PodContainerList(clusterId, nameSpace, vv.PodName)
			if err != nil || len(xList) == 0 {
				log.Printf("[ERROR] Get ContainerList error:%s\n", err)
				respBody = fmt.Sprintf("%s Get ContainerList error:%s\n", vv.PodName, err)
				this.Ctx.ResponseWriter.Write([]byte("Get ContainerList error"))
				this.Ctx.ResponseWriter.Flush()
				continue
			}
			containerName := xList[0].ContainerName
			stdout, stderr, _ := m.PodExec(clusterId, vv.NameSpace, vv.PodName, containerName, CmdStr, 300)
			respBody += fmt.Sprintf("\n=======check %s,%s result ==========\n", vv.PodIp, vv.PodName)
			respBody += fmt.Sprintf("stdout:\n%s\nstderr:\n%s", stdout, stderr)
			this.Ctx.ResponseWriter.Write([]byte(respBody))
			this.Ctx.ResponseWriter.Flush()
			if i == int(maxPodNum)-1 {
				break
			}
		}
	} else {
		log.Println("noSupport")
		code = -1
		msg = "noSupport"
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
}

// test
func (this *DeployController) Check2() {

	log.Println(string(this.Ctx.Input.RequestBody))
	// 设置流式响应头
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/event-stream")

	// 模拟循环任务
	for i := 1; i <= 30; i++ {
		// 执行任务并获取结果
		result := fmt.Sprintf("run num %d result: %s\n",
			i, time.Now().Format("2006-01-02 15:04:05"))

		// 流式写入
		this.Ctx.WriteString(result)
		this.Ctx.ResponseWriter.Flush()

		// 模拟任务间隔
		time.Sleep(1 * time.Second)
	}
}

func (this *DeployController) Host() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var hostArry []m.HostKv
	xList, err := m.DeployHost(clusterId, nameSpace, deployName, "GET", hostArry)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateHost() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var hostArry []m.HostKv
	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Host")

	//log.Println(string(this.Ctx.Input.RequestBody))
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	for i := 0; i <= 20; i++ {
		kk := gp.Get("domain[" + fmt.Sprintf("%d", i) + "]").Str
		vv := gp.Get("ip[" + fmt.Sprintf("%d", i) + "]").Str
		if kk != "" && vv != "" {
			//log.Println(kk, vv)
			hostArry = append(hostArry, m.HostKv{
				Domain: kk,
				Ip:     vv,
			})
		} else {
			break
		}
	}
	msg := "success"
	code := 0
	_, err := m.DeployHost(clusterId, nameSpace, deployName, "POST", hostArry)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateHost deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Resource() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	//containerName := this.GetString("containerName")

	var resource m.ResourceSt
	xList, err := m.DeployResource(clusterId, nameSpace, deployName, "GET", resource)
	msg := "success"
	code := 0
	count := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}
	count = len(xList)
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateResource() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")
	//containerName := this.GetString("containerName")

	var resource m.ResourceSt
	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Resource")

	//log.Println(string(this.Ctx.Input.RequestBody))
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	resource = m.ResourceSt{
		ContainerId: int(gp.Get("containerId").Int()),
		Request_cpu: gp.Get("request_cpu").Str,
		Request_mem: gp.Get("request_mem").Str,
		Limit_cpu:   gp.Get("limit_cpu").Str,
		Limit_mem:   gp.Get("limit_mem").Str,
	}
	msg := "success"
	code := 0
	_, err := m.DeployResource(clusterId, nameSpace, deployName, "POST", resource)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateResource deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Probe() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var probe m.ProbeST
	xList, err := m.DeployProbe(clusterId, nameSpace, deployName, "GET", probe)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateProbe() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var probe m.ProbeST
	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Probe")

	//log.Println(string(this.Ctx.Input.RequestBody))
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	probe = m.ProbeST{
		ContainerId:                   int(gp.Get("containerId").Int()),
		Readiness_checkType:           gp.Get("readiness_checkType").Str,
		Readiness_path:                gp.Get("readiness_path").Str,
		Readiness_httpPort:            int32(gp.Get("readiness_httpPort").Int()),
		Readiness_tcpPort:             int32(gp.Get("readiness_tcpPort").Int()),
		Readiness_cmd:                 gp.Get("readiness_cmd").Str,
		Readiness_initialDelaySeconds: int32(gp.Get("readiness_initialDelaySeconds").Int()),
		Readiness_periodSeconds:       int32(gp.Get("readiness_periodSeconds").Int()),
		Readiness_successThreshold:    int32(gp.Get("readiness_successThreshold").Int()),
		Readiness_failureThreshold:    int32(gp.Get("readiness_failureThreshold").Int()),
		Readiness_timeoutSeconds:      int32(gp.Get("readiness_timeoutSeconds").Int()),
		Liveness_checkType:            gp.Get("liveness_checkType").Str,
		Liveness_path:                 gp.Get("liveness_path").Str,
		Liveness_httpPort:             int32(gp.Get("liveness_httpPort").Int()),
		Liveness_tcpPort:              int32(gp.Get("liveness_tcpPort").Int()),
		Liveness_cmd:                  gp.Get("liveness_cmd").Str,
		Liveness_initialDelaySeconds:  int32(gp.Get("liveness_initialDelaySeconds").Int()),
		Liveness_periodSeconds:        int32(gp.Get("liveness_periodSeconds").Int()),
		Liveness_successThreshold:     int32(gp.Get("liveness_successThreshold").Int()),
		Liveness_failureThreshold:     int32(gp.Get("liveness_failureThreshold").Int()),
		Liveness_timeoutSeconds:       int32(gp.Get("liveness_timeoutSeconds").Int()),
	}
	msg := "success"
	code := 0
	_, err := m.DeployProbe(clusterId, nameSpace, deployName, "POST", probe)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateProbe deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Env() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var envst m.EnvSt
	xList, err := m.DeployEnv(clusterId, nameSpace, deployName, "GET", envst)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateEnv() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var envArry []m.EnvKv
	var envst m.EnvSt

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Env")

	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	for i := 0; i <= 40; i++ {
		kk := gp.Get("env_key[" + fmt.Sprintf("%d", i) + "]").Str
		vv := gp.Get("env_value[" + fmt.Sprintf("%d", i) + "]").Str
		if kk != "" && vv != "" {
			envArry = append(envArry, m.EnvKv{
				Key:   kk,
				Value: vv,
			})
		} else {
			break
		}
	}
	envst = m.EnvSt{
		ContainerId:   int(gp.Get("containerId").Int()),
		ContainerName: gp.Get("containerName").String(),
		Envs:          envArry,
	}
	msg := "success"
	code := 0
	_, err := m.DeployEnv(clusterId, nameSpace, deployName, "POST", envst)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateEnv deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) Lifecycle() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var lifecycle m.LifecycleSt
	xList, err := m.DeployLifecycle(clusterId, nameSpace, deployName, "GET", lifecycle)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateLifecycle() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var lifecycle m.LifecycleSt

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "update Lifecycle")

	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	lifecycle = m.LifecycleSt{
		ContainerId:             int(gp.Get("containerId").Int()),
		ContainerName:           gp.Get("containerName").String(),
		PreStop_type:            gp.Get("preStop_type").String(),
		PreStop_execCommand:     gp.Get("preStop_execCommand").String(),
		PreStop_httpGetPath:     gp.Get("preStop_httpGetPath").String(),
		PreStop_httpGetPort:     int32(gp.Get("preStop_httpGetPort").Int()),
		PreStop_httpGetScheme:   gp.Get("preStop_httpGetScheme").String(),
		PreStop_tcpSocketPort:   int32(gp.Get("preStop_tcpSocketPort").Int()),
		PostStart_type:          gp.Get("postStart_type").String(),
		PostStart_execCommand:   gp.Get("postStart_execCommand").String(),
		PostStart_httpGetPath:   gp.Get("postStart_httpGetPath").String(),
		PostStart_httpGetPort:   int32(gp.Get("postStart_httpGetPort").Int()),
		PostStart_httpGetScheme: gp.Get("postStart_httpGetScheme").String(),
		PostStart_tcpSocketPort: int32(gp.Get("postStart_tcpSocketPort").Int()),
	}
	msg := "success"
	code := 0
	_, err := m.DeployLifecycle(clusterId, nameSpace, deployName, "POST", lifecycle)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateLifecycle deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *DeployController) GetTolerations() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var TolerationArry []m.TolerationSt
	xList, err := m.DeployTolerations(clusterId, nameSpace, deployName, "GET", TolerationArry)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateTolerations() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "UpdateTolerations")

	var TolerationArry []m.TolerationSt
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	for i := 0; i <= 20; i++ {
		key := gp.Get("key[" + fmt.Sprintf("%d", i) + "]").Str
		operator := gp.Get("operator[" + fmt.Sprintf("%d", i) + "]").Str
		value := gp.Get("value[" + fmt.Sprintf("%d", i) + "]").Str
		effect := gp.Get("effect[" + fmt.Sprintf("%d", i) + "]").Str
		tolerationSeconds := gp.Get("tolerationSeconds[" + fmt.Sprintf("%d", i) + "]").Int()
		if key != "" && value != "" {
			//log.Println(kk, vv)
			TolerationArry = append(TolerationArry, m.TolerationSt{
				Key:               key,
				Operator:          operator,
				Value:             value,
				Effect:            effect,
				TolerationSeconds: &tolerationSeconds,
			})
		} else {
			break
		}
	}

	xList, err := m.DeployTolerations(clusterId, nameSpace, deployName, "POST", TolerationArry)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateTolerations deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) GetNodeAffinity() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var nodeAffinitySt m.DeployNodeAffinitySt
	xList, err := m.DeployNodeAffinity(clusterId, nameSpace, deployName, "GET", nodeAffinitySt)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdateNodeAffinity() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "UpdateNodeAffinity")

	var nodeAffinitySt m.DeployNodeAffinitySt
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	nodeName := gp.Get("nodeNames").Str

	//解析nodeSelector
	var nodeSelectors []m.LabelsKv
	for i := 0; i <= 10; i++ {
		kk := gp.Get("labels_key[" + fmt.Sprintf("%d", i) + "]").Str
		vv := gp.Get("labels_value[" + fmt.Sprintf("%d", i) + "]").Str
		if kk != "" && vv != "" {
			nodeSelectors = append(nodeSelectors, m.LabelsKv{
				Key:   kk,
				Value: vv,
			})
		}
	}

	//解析硬亲和
	var matchExpression []m.MatchExpression
	for i := 1; i <= 10; i++ {
		var exressions = make([]m.Expression, 0)
		for j := 1; j <= 10; j++ {
			rKey := gp.Get("rKey_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			rOperator := gp.Get("rOperator_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			rValues := gp.Get("rValues_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			if rKey != "" && rOperator != "" {
				exressions = append(exressions, m.Expression{
					Key:      rKey,
					Operator: rOperator,
					Value:    strings.Split(rValues, ","),
				})
			}
		}
		if len(exressions) > 0 {
			matchExpression = append(matchExpression, m.MatchExpression{
				Expressions: exressions,
			})
		}
	}
	//解析软亲和
	var preferences []m.Preference
	for i := 1; i <= 10; i++ {
		weightStr := gp.Get("weight_" + fmt.Sprintf("%d", i)).Str
		var weight int64
		var err error
		if weightStr != "" {
			weight, err = strconv.ParseInt(weightStr, 10, 64)
			if err != nil {
				log.Printf("[WARN] UpdateNodeAffinity ParseInt weightStr Fail:%s\n", err.Error())
				this.Data["json"] = &map[string]interface{}{"code": -1, "msg": err.Error()}
				this.ServeJSON()
				return
			}
		}
		var exressions = make([]m.Expression, 0)
		for j := 1; j <= 10; j++ {
			pKey := gp.Get("pKey_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			pOperator := gp.Get("pOperator_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			pValues := gp.Get("pValues_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			if pKey != "" && pOperator != "" {
				exressions = append(exressions, m.Expression{
					Key:      pKey,
					Operator: pOperator,
					Value:    strings.Split(pValues, ","),
				})
			}
		}
		if len(exressions) > 0 {
			preferences = append(preferences, m.Preference{
				Weight:      int32(weight),
				Expressions: exressions,
			})
		}
	}

	nodeAffinitySt = m.DeployNodeAffinitySt{
		NodeNames:         nodeName,
		NodeSelector:      nodeSelectors,
		RequiredAffinity:  matchExpression,
		PreferredAffinity: preferences,
	}

	xList, err := m.DeployNodeAffinity(clusterId, nameSpace, deployName, "POST", nodeAffinitySt)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdateNodeAffinity deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) GetPodAffinity() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	var podAffinitySt m.DeployPodAffinitySt
	xList, err := m.DeployPodAffinity(clusterId, nameSpace, deployName, "GET", podAffinitySt)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &xList}
	this.ServeJSON()
}

func (this *DeployController) UpdatePodAffinity() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	deployName := this.GetString("deployName")

	//backup
	yamlStr, _ := m.GetDeployYaml(clusterId, nameSpace, deployName)
	_ = m.InsertBackup(clusterId, nameSpace, deployName, "deployment", yamlStr, "UpdatePodAffinity")

	var podAffinitySt m.DeployPodAffinitySt
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	// 解析 PodAffinity Required
	var podAffinityRequired []m.PodAffinityReq
	for i := 1; i <= 10; i++ {
		namespacesStr := gp.Get("podAffinityRequiredNamespaces_" + fmt.Sprintf("%d", i)).Str
		topologyKey := gp.Get("podAffinityRequiredTopologyKey_" + fmt.Sprintf("%d", i)).Str
		matchLabelKeysStr := gp.Get("podAffinityRequiredMatchLabelKeys_" + fmt.Sprintf("%d", i)).Str
		if topologyKey != "" {
			var labelSelectors []m.Expression
			for j := 1; j <= 10; j++ {
				paKey := gp.Get("podAffinityRequiredLabelKey_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
				paOperator := gp.Get("podAffinityRequiredLabelOperator_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
				paValues := gp.Get("podAffinityRequiredLabelValue_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
				if paKey != "" && paOperator != "" {
					var values []string
					if paValues != "" {
						values = strings.Split(paValues, ",")
					}
					labelSelectors = append(labelSelectors, m.Expression{
						Key:      paKey,
						Operator: paOperator,
						Value:    values,
					})
				}
			}

			var namespaces []string
			if namespacesStr != "" {
				namespaces = strings.Split(namespacesStr, ",")
			}

			var matchLabelKeys []string
			if matchLabelKeysStr != "" {
				matchLabelKeys = strings.Split(matchLabelKeysStr, ",")
			}

			podAffinityRequired = append(podAffinityRequired, m.PodAffinityReq{
				LabelSelector:  labelSelectors,
				TopologyKey:    topologyKey,
				NameSpaces:     namespaces,
				MatchLableKeys: matchLabelKeys,
			})
		}
	}

	// 解析 PodAffinity Preferred
	var podAffinityPreferred []m.PodAffinityTerm
	for i := 1; i <= 10; i++ {
		weightStr := gp.Get("podAffinityPreferredWeight_" + fmt.Sprintf("%d", i)).Str
		var weight int64
		var err error
		if weightStr != "" {
			weight, err = strconv.ParseInt(weightStr, 10, 64)
			if err != nil {
				log.Printf("[WARN] UpdatePodAffinity ParseInt podAffinityPreferredWeight Fail:%s\n", err.Error())
				this.Data["json"] = &map[string]interface{}{"code": -1, "msg": err.Error()}
				this.ServeJSON()
				return
			}
		}

		var labelSelectors = make([]m.Expression, 0)
		for j := 1; j <= 10; j++ {
			paKey := gp.Get("podAffinityPreferredLabelKey_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			paOperator := gp.Get("podAffinityPreferredLabelOperator_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			paValues := gp.Get("podAffinityPreferredLabelValue_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			if paKey != "" && paOperator != "" {
				var values []string
				if paValues != "" {
					values = strings.Split(paValues, ",")
				}
				labelSelectors = append(labelSelectors, m.Expression{
					Key:      paKey,
					Operator: paOperator,
					Value:    values,
				})
			}
		}

		var namespaces []string
		namespacesStr := gp.Get("podAffinityPreferredNamespaces_" + fmt.Sprintf("%d", i)).Str
		if namespacesStr != "" {
			namespaces = strings.Split(namespacesStr, ",")
		}

		var matchLabelKeys []string
		matchLabelKeysStr := gp.Get("podAffinityPreferredMatchLabelKeys_" + fmt.Sprintf("%d", i)).Str
		if matchLabelKeysStr != "" {
			matchLabelKeys = strings.Split(matchLabelKeysStr, ",")
		}

		if len(labelSelectors) > 0 || weight > 0 {
			podAffinityPreferred = append(podAffinityPreferred, m.PodAffinityTerm{
				Weight:         int32(weight),
				LabelSelector:  labelSelectors,
				TopologyKey:    gp.Get("podAffinityPreferredTopologyKey_" + fmt.Sprintf("%d", i)).Str,
				NameSpaces:     namespaces,
				MatchLableKeys: matchLabelKeys,
			})
		}
	}

	// 解析 PodAntiAffinity Required
	var podAntiAffinityRequired []m.PodAffinityReq
	for i := 1; i <= 10; i++ {
		namespacesStr := gp.Get("podAntiAffinityRequiredNamespaces_" + fmt.Sprintf("%d", i)).Str
		topologyKey := gp.Get("podAntiAffinityRequiredTopologyKey_" + fmt.Sprintf("%d", i)).Str
		matchLabelKeysStr := gp.Get("podAntiAffinityRequiredMatchLabelKeys_" + fmt.Sprintf("%d", i)).Str
		if topologyKey != "" {
			var labelSelectors []m.Expression
			for j := 1; j <= 10; j++ {
				paaKey := gp.Get("podAntiAffinityRequiredLabelKey_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
				paaOperator := gp.Get("podAntiAffinityRequiredLabelOperator_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
				paaValues := gp.Get("podAntiAffinityRequiredLabelValue_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
				if paaKey != "" && paaOperator != "" {
					var values []string
					if paaValues != "" {
						values = strings.Split(paaValues, ",")
					}
					labelSelectors = append(labelSelectors, m.Expression{
						Key:      paaKey,
						Operator: paaOperator,
						Value:    values,
					})
				}
			}

			var namespaces []string
			if namespacesStr != "" {
				namespaces = strings.Split(namespacesStr, ",")
			}

			var matchLabelKeys []string
			if matchLabelKeysStr != "" {
				matchLabelKeys = strings.Split(matchLabelKeysStr, ",")
			}

			podAntiAffinityRequired = append(podAntiAffinityRequired, m.PodAffinityReq{
				LabelSelector:  labelSelectors,
				TopologyKey:    topologyKey,
				NameSpaces:     namespaces,
				MatchLableKeys: matchLabelKeys,
			})
		}
	}

	// 解析 PodAntiAffinity Preferred
	var podAntiAffinityPreferred []m.PodAffinityTerm
	for i := 1; i <= 10; i++ {
		weightStr := gp.Get("podAntiAffinityPreferredWeight_" + fmt.Sprintf("%d", i)).Str
		var weight int64
		var err error
		if weightStr != "" {
			weight, err = strconv.ParseInt(weightStr, 10, 64)
			if err != nil {
				log.Printf("[WARN] UpdatePodAffinity ParseInt podAntiAffinityPreferredWeight Fail:%s\n", err.Error())
				this.Data["json"] = &map[string]interface{}{"code": -1, "msg": err.Error()}
				this.ServeJSON()
				return
			}
		}

		var labelSelectors = make([]m.Expression, 0)
		for j := 1; j <= 10; j++ {
			paaKey := gp.Get("podAntiAffinityPreferredLabelKey_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			paaOperator := gp.Get("podAntiAffinityPreferredLabelOperator_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			paaValues := gp.Get("podAntiAffinityPreferredLabelValue_" + fmt.Sprintf("%d", i) + "_" + fmt.Sprintf("%d", j)).Str
			if paaKey != "" && paaValues != "" {
				var values []string
				if paaValues != "" {
					values = strings.Split(paaValues, ",")
				}
				labelSelectors = append(labelSelectors, m.Expression{
					Key:      paaKey,
					Operator: paaOperator,
					Value:    values,
				})
			}
		}

		var namespaces []string
		namespacesStr := gp.Get("podAntiAffinityPreferredNamespaces_" + fmt.Sprintf("%d", i)).Str
		if namespacesStr != "" {
			namespaces = strings.Split(namespacesStr, ",")
		}

		var matchLabelKeys []string
		matchLabelKeysStr := gp.Get("podAntiAffinityPreferredMatchLabelKeys_" + fmt.Sprintf("%d", i)).Str
		if matchLabelKeysStr != "" {
			matchLabelKeys = strings.Split(matchLabelKeysStr, ",")
		}

		if len(labelSelectors) > 0 || weight > 0 {
			podAntiAffinityPreferred = append(podAntiAffinityPreferred, m.PodAffinityTerm{
				Weight:         int32(weight),
				LabelSelector:  labelSelectors,
				NameSpaces:     namespaces,
				MatchLableKeys: matchLabelKeys,
				TopologyKey:    gp.Get("podAntiAffinityPreferredTopologyKey_" + fmt.Sprintf("%d", i)).Str,
			})
		}
	}

	podAffinitySt = m.DeployPodAffinitySt{
		PodAffinity: m.PodAffinitySt{
			Required:  podAffinityRequired,
			Preferred: podAffinityPreferred,
		},
		PodAntiAffinity: m.PodAntiAffinitySt{
			Required:  podAntiAffinityRequired,
			Preferred: podAntiAffinityPreferred,
		},
	}

	xList, err := m.DeployPodAffinity(clusterId, nameSpace, deployName, "POST", podAffinitySt)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	//auditlog
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = adm.InsertLogAudit(userip, userinfo.(adm.User).Username, "update", fmt.Sprintf("%d", code), "UpdatePodAffinity deploy:"+deployName+",result:"+msg)

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": &xList}
	this.ServeJSON()
}
