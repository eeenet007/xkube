// node
package controllers

import (
	"fmt"
	"log"
	"strings"
	"xkube/common"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type NodeController struct {
	beego.Controller
}

func (this *NodeController) List() {
	clusterId := this.GetString("clusterId")
	nodeList, err := m.NodeList(clusterId)
	msg := "success"
	code := 0
	count := len(nodeList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	//set in redis
	if code == 0 {
		_ = common.HSet("count_"+clusterId, "node", fmt.Sprintf("%d", count))
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &nodeList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *NodeController) Detail() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	xdetail, err := m.NodeDetail(clusterId, nodeName)
	if err != nil {
		log.Println(err)
	}
	//this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": len(xList), "data": &xList}
	this.Data["json"] = &xdetail
	this.ServeJSON()
}

func (this *NodeController) PoolList() {
	clusterId := this.GetString("clusterId")
	xList, err := m.NodePoolList(clusterId)
	msg := "success"
	code := 0
	count := len(xList)
	if err != nil {
		log.Println(err)
		msg = err.Error()
		code = -1
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &xList}
	//this.Data["json"] = &datas
	this.ServeJSON()
}

func (this *NodeController) Yaml() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	yamlStr, _ := m.GetNodeYaml(clusterId, nodeName)
	this.Ctx.WriteString(yamlStr)
}

func (this *NodeController) ModifyByYaml() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	//backup
	yamlStr, _ := m.GetNodeYaml(clusterId, nodeName)
	_ = m.InsertBackup(clusterId, "", nodeName, "node", yamlStr, "Backup before updating")

	code := 0
	msg := "success"

	reqBody := strings.ReplaceAll(string(this.Ctx.Input.RequestBody), "%25", "%")
	reqBody = strings.ReplaceAll(reqBody, "%3B", ";")
	err := m.NodeYamlModify(clusterId, []byte(reqBody))
	if err != nil {
		log.Printf("[WARN] NodeYamlModify Fail:%s\n", err)
		code = -1
		msg = err.Error()
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 排水
func (this *NodeController) Drain() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	beta := this.GetString("beta") //default:0，When beta1 is 1, use the beta version
	code := 0
	msg := "success"

	//设置不可调度
	err := m.NodeUnschedulable(clusterId, nodeName, true)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] Unschedulable Fail:%s\n", err)
	}
	if beta == "1" {
		err = m.NodeDrainV1beta1(clusterId, nodeName)
		if err != nil {
			code = -1
			msg = err.Error()
			log.Printf("[ERROR] NodeDrainveta1 Fail:%s\n", err)
		}
	} else {
		//排水
		err = m.NodeDrain(clusterId, nodeName)
		if err != nil {
			code = -1
			msg = err.Error()
			log.Printf("[ERROR] NodeDrain Fail:%s\n", err)
		}
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 调度
func (this *NodeController) Unschedulable() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	code := 0
	msg := "success"
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	unschedulableValue := gp.Get("unschedulable").Bool()
	err := m.NodeUnschedulable(clusterId, nodeName, unschedulableValue)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] Unschedulable Fail:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

// 移除
func (this *NodeController) Del() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	//gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	//isDrain := gp.Get("isDrain").Sting()
	code := 0
	msg := "success"
	// if isDrain == "on" {
	// 	err := m.NodeDrain(clusterId, nodeName)
	// 	if err != nil {
	// 		code = -1
	// 		msg = err.Error()
	// 		log.Printf("[ERROR] NodeDel NodeDrain Fail:%s\n", err)
	// 		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	// 		this.ServeJSON()
	// 		return
	// 	}
	// }

	err := m.NodeDelete(clusterId, nodeName)
	if err != nil {
		code = -1
		msg = err.Error()
		log.Printf("[ERROR] NodeDelete Fail:%s\n", err)
	}
	this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
	this.ServeJSON()
}

func (this *NodeController) GetLabels() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	xList, err := m.GetNodeLabels(clusterId, nodeName)
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

func (this *NodeController) UpdateLabels() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	var labelsMap = make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {
		//backup
		yamlStr, _ := m.GetNodeYaml(clusterId, nodeName)
		_ = m.InsertBackup(clusterId, "", nodeName, "node", yamlStr, "Backup before updating")

		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for i := 0; i <= 40; i++ {
			kk := gp.Get("labels_key[" + fmt.Sprintf("%d", i) + "]").Str
			vv := gp.Get("labels_value[" + fmt.Sprintf("%d", i) + "]").Str
			if kk != "" && vv != "" {
				labelsMap[kk] = vv
			}
		}
		msg := "success"
		code := 0
		err := m.UpdateNodeLabels(clusterId, nodeName, labelsMap)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
}

func (this *NodeController) GetTaint() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	xList, err := m.GetNodeTaint(clusterId, nodeName)
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

func (this *NodeController) UpdateTaint() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")
	var taint = make([]m.NodeTaint, 0)
	if this.Ctx.Input.Method() == "POST" {
		//backup
		yamlStr, _ := m.GetNodeYaml(clusterId, nodeName)
		_ = m.InsertBackup(clusterId, "", nodeName, "node", yamlStr, "Backup before updating")

		gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		for i := 0; i <= 40; i++ {
			kk := gp.Get("key[" + fmt.Sprintf("%d", i) + "]").Str
			vv := gp.Get("value[" + fmt.Sprintf("%d", i) + "]").Str
			ee := gp.Get("effect[" + fmt.Sprintf("%d", i) + "]").Str
			if kk != "" && vv != "" {
				taint = append(taint, m.NodeTaint{
					Key:    kk,
					Effect: ee,
					Value:  vv,
				})
			}
		}
		msg := "success"
		code := 0
		err := m.UpdateNodeTaint(clusterId, nodeName, taint)
		if err != nil {
			msg = err.Error()
			code = -1
		}
		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg}
		this.ServeJSON()
		return
	}
}

func (this *NodeController) GetAllocated() {
	clusterId := this.GetString("clusterId")
	nodeName := this.GetString("nodeName")

	result, err := m.NodeAllocated(clusterId, nodeName)
	if err != nil {
		this.Data["json"] = &map[string]interface{}{
			"code":  -1,
			"msg":   err.Error(),
			"count": 0,
			"data":  []string{},
		}
	} else {
		this.Data["json"] = result
	}

	this.ServeJSON()
}
