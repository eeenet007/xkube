package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	m "xkube/models"
	xm "xkube/xadmin/src/models"

	"github.com/tidwall/gjson"

	beego "github.com/beego/beego/v2/server/web"
)

type StatefulsetController struct {
	beego.Controller
}

func (this *StatefulsetController) List() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	labels := this.GetString("labels")
	if this.Ctx.Input.Method() == "POST" {
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		statefulsetName = gp.Get("statefulsetName").String()
		nameSpace = gp.Get("nameSpace").String()
	}
	labelsKV := strings.Split(labels, ":")
	var labelsKey, labelsValue string
	if len(labelsKV) == 2 {
		labelsKey = labelsKV[0]
		labelsValue = labelsKV[1]
	}

	dxList, err := m.StatefulsetList(clusterId, nameSpace, statefulsetName, labelsKey, labelsValue)
	msg := "success"
	code := 0
	count := len(dxList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &dxList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *StatefulsetController) Detail() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	xdetail, err := m.StatefulsetDetail(clusterId, namespace, statefulsetName)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *StatefulsetController) Create() {
	clusterId := this.GetString("clusterId")

	code := 0
	msg := "success"
	err := m.StatefulsetCreate(clusterId, this.Ctx.Input.RequestBody)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] Deploy Create Fail:%s\n", err)
	}
	_ = m.ClearCache(clusterId) //创建以后，刷新一下列表缓存
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *StatefulsetController) Modify() {
	clusterId := this.GetString("clusterId")
	var sts m.Statefulset

	err := json.Unmarshal(this.Ctx.Input.RequestBody, &sts)
	err = m.StatefulsetModify(clusterId, &sts)
	if err != nil {
		log.Println(err)
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "ok"}
	this.ServeJSON()
}

func (this *StatefulsetController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	//backup
	yamlStr, _ := m.GetStatefulsetYaml(clusterId, namespace, statefulsetName)
	_ = m.InsertBackup(clusterId, namespace, statefulsetName, "statefulset", yamlStr, "Backup before updating")

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	code := 0
	msg := "success"
	err := m.StatefulsetYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] StatefulsetYamlModify Fail:%s\n", err)
		code = 1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 迁移到其他集群
func (this *StatefulsetController) Clone() {
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	clusterId := gp.Get("clusterid").String()
	namespace := gp.Get("namespace").String()
	target_clusterid := gp.Get("target_clusterid").String()
	target_namespace := gp.Get("target_namespace").String()
	target_objname := gp.Get("target_objname").String()
	objName := gp.Get("objname").String()
	objNameArry := strings.Split(objName, ",")
	targetObjnameArry := strings.Split(target_objname, ",")
	var result string
	var isOk = true

	for i, vv := range objNameArry {
		var targetObjname string
		if len(targetObjnameArry)-1 >= i {
			targetObjname = targetObjnameArry[i]
		}
		err := m.StatefulsetClone(clusterId, namespace, vv, target_clusterid, target_namespace, targetObjname)
		if err != nil {
			isOk = false
		} else {
			err = fmt.Errorf("ok")
		}
		result += fmt.Sprintf("%s result:%s,", vv, err)
	}
	uinfo := this.GetSession("userinfo")
	u := m.Xkb_clone_log{
		Clusterid:       clusterId,
		Namespace:       namespace,
		Restype:         "statefulset",
		Objname:         objName,
		TargetClusterid: target_clusterid,
		TargetNamespace: target_namespace,
		TargetObjname:   target_objname,
		Status:          fmt.Sprintf("%v", isOk),
		User:            uinfo.(xm.User).Username,
		Result:          result,
	}
	_, err := m.Add_Clone(&u) //插入clone 结果
	if err != nil {
		log.Printf("[ERROR] Clone sts add Fail:%s\n", err)
		result += fmt.Sprintf("insert log fail")
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "success": isOk, "msg": result}
	this.ServeJSON()
}

func (this *StatefulsetController) Del() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	err := m.StatefulsetDel(clusterId, namespace, statefulsetName)
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

func (this *StatefulsetController) Yaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	yamlStr, _ := m.GetStatefulsetYaml(clusterId, namespace, statefulsetName)
	this.Ctx.WriteString(yamlStr)
	//this.Data["yaml"] = &yamlStr
	//this.ServeYAML()
	//this.ServeJSON()
}

func (this *StatefulsetController) ReplicaSetYaml() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	replicaSet := this.GetString("replicaSet")

	yamlStr, _ := m.GetReplicasetYaml(clusterId, namespace, replicaSet)
	this.Ctx.WriteString(yamlStr)
}

func (this *StatefulsetController) RollBack() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	replicaSet := this.GetString("replicaSet")

	code := 0
	msg := "success"
	err := m.StatefulsetRollBack(clusterId, namespace, statefulsetName, replicaSet)
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *StatefulsetController) Restart() {
	clusterId := this.GetString("clusterId")
	namespace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	code := 0
	msg := "success"
	err := m.StatefulsetRestart(clusterId, namespace, statefulsetName)
	if err != nil {
		log.Println(err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *StatefulsetController) Labels() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	var labelsMap = make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {
		//log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for i := 0; i <= 20; i++ {
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
		_, err := m.StatefulsetLabels(clusterId, nameSpace, statefulsetName, "POST", labelsMap)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	xList, err := m.StatefulsetLabels(clusterId, nameSpace, statefulsetName, "GET", labelsMap)
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

func (this *StatefulsetController) Image() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	var bbb = make([]m.ImageKv, 0)
	if this.Ctx.Input.Method() == "POST" {
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
		_, err := m.StatefulsetImage(clusterId, nameSpace, statefulsetName, "POST", bbb)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	xList, err := m.StatefulsetImage(clusterId, nameSpace, statefulsetName, "GET", bbb)
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

func (this *StatefulsetController) Replicas() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")
	if this.Ctx.Input.Method() == "POST" {
		//log.Println(string(this.Ctx.Input.RequestBody))
		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		podNum := gp.Get("podNumber").Int()
		msg := "success"
		code := 0
		_, err := m.StatefulsetReplicas(clusterId, nameSpace, statefulsetName, "POST", int32(podNum))
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
	resp, err := m.StatefulsetReplicas(clusterId, nameSpace, statefulsetName, "GET", 0)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "number": resp}
	this.ServeJSON()
}

func (this *StatefulsetController) GetHost() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var hostArry []m.HostKv
	xList, err := m.StatefulsetHost(clusterId, nameSpace, statefulsetName, "GET", hostArry)
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

func (this *StatefulsetController) UpdateHost() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var hostArry []m.HostKv
	//backup
	yamlStr, _ := m.GetStatefulsetYaml(clusterId, nameSpace, statefulsetName)
	_ = m.InsertBackup(clusterId, nameSpace, statefulsetName, "statefulset", yamlStr, "update Host")

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
	_, err := m.StatefulsetHost(clusterId, nameSpace, statefulsetName, "POST", hostArry)
	if err != nil {
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *StatefulsetController) GetEnv() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var envst m.EnvSt
	xList, err := m.StatefulsetEnv(clusterId, nameSpace, statefulsetName, "GET", envst)
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

func (this *StatefulsetController) UpdateEnv() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var envArry []m.EnvKv
	var envst m.EnvSt
	//backup
	yamlStr, _ := m.GetStatefulsetYaml(clusterId, nameSpace, statefulsetName)
	_ = m.InsertBackup(clusterId, nameSpace, statefulsetName, "statefulset", yamlStr, "update Env")

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
	_, err := m.StatefulsetEnv(clusterId, nameSpace, statefulsetName, "POST", envst)
	if err != nil {
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()

}

func (this *StatefulsetController) GetResource() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var resource m.ResourceSt
	xList, err := m.StatefulsetResource(clusterId, nameSpace, statefulsetName, "GET", resource)
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

func (this *StatefulsetController) UpdateResource() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var resource m.ResourceSt
	//backup
	yamlStr, _ := m.GetStatefulsetYaml(clusterId, nameSpace, statefulsetName)
	_ = m.InsertBackup(clusterId, nameSpace, statefulsetName, "statefulset", yamlStr, "update Resource")

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
	_, err := m.StatefulsetResource(clusterId, nameSpace, statefulsetName, "POST", resource)
	if err != nil {
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *StatefulsetController) GetLifecycle() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var lifecycle m.LifecycleSt
	xList, err := m.StatefulsetLifecycle(clusterId, nameSpace, statefulsetName, "GET", lifecycle)
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

func (this *StatefulsetController) UpdateLifecycle() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var lifecycle m.LifecycleSt
	//backup
	yamlStr, _ := m.GetStatefulsetYaml(clusterId, nameSpace, statefulsetName)
	_ = m.InsertBackup(clusterId, nameSpace, statefulsetName, "statefulset", yamlStr, "update Lifecycle")

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
	_, err := m.StatefulsetLifecycle(clusterId, nameSpace, statefulsetName, "POST", lifecycle)
	if err != nil {
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *StatefulsetController) GetProbe() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var probe m.ProbeST

	xList, err := m.StatefulsetProbe(clusterId, nameSpace, statefulsetName, "GET", probe)
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

func (this *StatefulsetController) UpdateProbe() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var probe m.ProbeST
	//backup
	yamlStr, _ := m.GetStatefulsetYaml(clusterId, nameSpace, statefulsetName)
	_ = m.InsertBackup(clusterId, nameSpace, statefulsetName, "statefulset", yamlStr, "update Probe")

	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	probe = m.ProbeST{
		ContainerId:                   int(gp.Get("containerId").Int()),
		ContainerName:                 gp.Get("containerName").String(),
		Readiness_checkType:           gp.Get("readiness_checkType").String(),
		Readiness_path:                gp.Get("readiness_path").String(),
		Readiness_httpPort:            int32(gp.Get("readiness_httpPort").Int()),
		Readiness_tcpPort:             int32(gp.Get("readiness_tcpPort").Int()),
		Readiness_cmd:                 gp.Get("readiness_cmd").String(),
		Readiness_initialDelaySeconds: int32(gp.Get("readiness_initialDelaySeconds").Int()),
		Readiness_periodSeconds:       int32(gp.Get("readiness_periodSeconds").Int()),
		Readiness_successThreshold:    int32(gp.Get("readiness_successThreshold").Int()),
		Readiness_failureThreshold:    int32(gp.Get("readiness_failureThreshold").Int()),
		Readiness_timeoutSeconds:      int32(gp.Get("readiness_timeoutSeconds").Int()),
		Liveness_checkType:            gp.Get("liveness_checkType").String(),
		Liveness_path:                 gp.Get("liveness_path").String(),
		Liveness_httpPort:             int32(gp.Get("liveness_httpPort").Int()),
		Liveness_tcpPort:              int32(gp.Get("liveness_tcpPort").Int()),
		Liveness_cmd:                  gp.Get("liveness_cmd").String(),
		Liveness_initialDelaySeconds:  int32(gp.Get("liveness_initialDelaySeconds").Int()),
		Liveness_periodSeconds:        int32(gp.Get("liveness_periodSeconds").Int()),
		Liveness_successThreshold:     int32(gp.Get("liveness_successThreshold").Int()),
		Liveness_failureThreshold:     int32(gp.Get("liveness_failureThreshold").Int()),
		Liveness_timeoutSeconds:       int32(gp.Get("liveness_timeoutSeconds").Int()),
	}
	msg := "success"
	code := 0
	_, err := m.StatefulsetProbe(clusterId, nameSpace, statefulsetName, "POST", probe)
	if err != nil {
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
	return
}

func (this *StatefulsetController) GetNodeAffinity() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var affinity m.AffinitySt
	xList, err := m.StatefulsetNodeAffinity(clusterId, nameSpace, statefulsetName, "GET", affinity)
	msg := "success"
	code := 0
	if err != nil {
		msg = err.Error()
		code = -1
	}

	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "data": xList}
	this.ServeJSON()
}

func (this *StatefulsetController) UpdateNodeAffinity() {
	clusterId := this.GetString("clusterId")
	nameSpace := this.GetString("nameSpace")
	statefulsetName := this.GetString("statefulsetName")

	var affinity m.AffinitySt
	if this.Ctx.Input.Method() == "POST" {
		//backup
		yamlStr, _ := m.GetStatefulsetYaml(clusterId, nameSpace, statefulsetName)
		_ = m.InsertBackup(clusterId, nameSpace, statefulsetName, "statefulset", yamlStr, "update NodeAffinity")

		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

		// 处理nodeSelector
		for i := 0; i <= 10; i++ {
			kk := gp.Get("labels_key[" + fmt.Sprintf("%d", i) + "]").Str
			vv := gp.Get("labels_value[" + fmt.Sprintf("%d", i) + "]").Str
			if kk != "" && vv != "" {
				affinity.NodeSelector = append(affinity.NodeSelector, m.LabelsKv{
					Key:   kk,
					Value: vv,
				})
			}
		}

		// 处理requiredDuringSchedulingIgnoredDuringExecution
		for i := 0; ; i++ {
			keyPrefix := fmt.Sprintf("rKey_%d_", i+1)
			opPrefix := fmt.Sprintf("rOperator_%d_", i+1)
			valPrefix := fmt.Sprintf("rValues_%d_", i+1)

			// 检查是否存在这个规则组
			if !gp.Get(keyPrefix + "1").Exists() {
				break
			}

			var expressions []m.Expression
			for j := 1; ; j++ {
				key := gp.Get(fmt.Sprintf("%s%d", keyPrefix, j)).Str
				op := gp.Get(fmt.Sprintf("%s%d", opPrefix, j)).Str
				val := gp.Get(fmt.Sprintf("%s%d", valPrefix, j)).Str

				if key == "" || op == "" {
					break
				}

				expressions = append(expressions, m.Expression{
					Key:      key,
					Operator: op,
					Value:    strings.Split(val, ","),
				})
			}

			if len(expressions) > 0 {
				affinity.RequiredAffinity = append(affinity.RequiredAffinity, m.MatchExpression{
					Expressions: expressions,
				})
			}
		}

		// 处理preferredDuringSchedulingIgnoredDuringExecution
		for i := 0; ; i++ {
			weightKey := fmt.Sprintf("weight_%d", i+1)
			keyPrefix := fmt.Sprintf("pKey_%d_", i+1)
			opPrefix := fmt.Sprintf("pOperator_%d_", i+1)
			valPrefix := fmt.Sprintf("pValues_%d_", i+1)

			// 检查是否存在这个规则组
			weight := int32(gp.Get(weightKey).Int())
			if weight == 0 {
				break
			}

			var expressions []m.Expression
			for j := 1; ; j++ {
				key := gp.Get(fmt.Sprintf("%s%d", keyPrefix, j)).Str
				op := gp.Get(fmt.Sprintf("%s%d", opPrefix, j)).Str
				val := gp.Get(fmt.Sprintf("%s%d", valPrefix, j)).Str

				if key == "" || op == "" {
					break
				}

				expressions = append(expressions, m.Expression{
					Key:      key,
					Operator: op,
					Value:    strings.Split(val, ","),
				})
			}

			if len(expressions) > 0 {
				affinity.PreferredAffinity = append(affinity.PreferredAffinity, m.Preference{
					Weight:      weight,
					Expressions: expressions,
				})
			}
		}

		// 处理nodeName
		affinity.NodeNames = gp.Get("nodeNames").String()

		//log.Println(affinity)

		msg := "success"
		code := 0
		_, err := m.StatefulsetNodeAffinity(clusterId, nameSpace, statefulsetName, "POST", affinity)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
}
