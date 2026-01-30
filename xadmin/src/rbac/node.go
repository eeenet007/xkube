package rbac

import (
	//"encoding/json"
	"fmt"
	"log"
	"strings"
	m "xkube/xadmin/src/models"

	"github.com/tidwall/gjson"

	"github.com/beego/beego/v2/client/orm"
)

type NodeController struct {
	CommonController
}

func (this *NodeController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJSON()
}

func (this *NodeController) List() {
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("rows")
	sort := this.GetString("sort")
	order := this.GetString("order")
	if len(order) > 0 {
		if order == "desc" {
			sort = "-" + sort
		}
	} else {
		sort = "Id"
	}
	nodes, count := m.GetNodelist(page, page_size, sort)
	for i := 0; i < len(nodes); i++ {
		if nodes[i]["Pid"] != 0 {
			nodes[i]["_parentId"] = nodes[i]["Pid"]
		} else {
			nodes[i]["state"] = "closed"
		}
	}
	if len(nodes) < 1 {
		nodes = []orm.Params{}
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &nodes}
	this.ServeJSON()

}

// func (this *NodeController) AddEdit() {
// 	n := m.Node{}
// 	if err := this.ParseForm(&n); err != nil {
// 		//handle error
// 		this.Rsp(false, err.Error())
// 		return
// 	}
// 	var id int64
// 	var err error
// 	Nid, _ := this.GetInt64("Id")
// 	if Nid > 0 {
// 		id, err = m.UpdateNode(&n)
// 	} else {
// 		group_id, _ := this.GetInt64("Group_id")
// 		group := new(m.Group)
// 		group.Id = group_id
// 		n.Group = group
// 		if n.Pid != 0 {
// 			n1, _ := m.ReadNode(n.Pid)
// 			n.Level = n1.Level + 1
// 		} else {
// 			n.Level = 1
// 		}
// 		id, err = m.AddNode(&n)
// 	}
// 	if err == nil && id > 0 {
// 		this.Rsp(true, "Success")
// 		return
// 	} else {
// 		this.Rsp(false, err.Error())
// 		return
// 	}
// }

func (this *NodeController) Add() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	n := m.Node{}
	if err := this.ParseForm(&n); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	group_id, _ := this.GetInt64("Group_id")
	group := new(m.Group)
	group.Id = group_id
	n.Group = group
	if n.Pid != 0 {
		n1, _ := m.ReadNode(n.Pid)
		n.Level = n1.Level + 1
	} else {
		n.Level = 1
	}
	id, err := m.AddNode(&n)
	if err == nil && id > 0 {
		content := fmt.Sprintf("nodeId:%d,groupId:%d", id, group_id)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "AddNodeUrl", "success", content)
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}

func (this *NodeController) Edit() {

	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	id := this.GetString("id")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

	if gp.Get("name").Exists() {
		id = gp.Get("id").String()
		name := gp.Get("name").String()
		value := gp.Get("value").String()
		_, err := m.EditNode(id, strings.ToLower(name), value)
		if err == nil {
			content := fmt.Sprintf("nodeId:%s,info:%s:%s", id, name, value)
			_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "EditNodeUrl", "success", content)
			this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
		} else {
			log.Printf("[ERROR] Edit rbac Node Fail:%s\n", err)
			this.Ctx.WriteString(fmt.Sprintf(`{"code":-1,"msg":"%s"}`, err.Error()))
		}
	} else {
		this.Ctx.WriteString(`{"code":-1,"msg":"bodyParseFail"}`)
	}
}

func (this *NodeController) Delete() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	Id, _ := this.GetInt64("Id")
	status, err := m.DelNodeById(Id)
	if err == nil && status > 0 {
		content := fmt.Sprintf("nodeId:%d", Id)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "DelNodeUrl", "success", content)
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}

func (this *NodeController) Getlist() {
	groupId, _ := this.GetInt64("Id")
	nodes, count := m.GetNodelistByGroupid(groupId)
	if len(nodes) < 1 {
		nodes = []orm.Params{}
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &nodes}
	this.ServeJSON()
}

func (this *NodeController) GetPid() {
	pids, count := m.GetNodePid()
	if len(pids) < 1 {
		pids = []orm.Params{}
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &pids}
	this.ServeJSON()

}
