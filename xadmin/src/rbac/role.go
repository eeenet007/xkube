package rbac

import (
	//"encoding/json"
	"fmt"
	"strconv"
	"strings"

	m "xkube/xadmin/src/models"

	"github.com/beego/beego/v2/client/orm"
)

type RoleController struct {
	CommonController
}

func (this *RoleController) List() {
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
	roles, count := m.GetRolelist(page, page_size, sort)
	if len(roles) < 1 {
		roles = []orm.Params{}
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &roles}
	this.ServeJSON()
}

func (this *RoleController) AddAndEdit() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	r := m.Role{}
	if err := this.ParseForm(&r); err != nil {
		//handle error
		this.Rsp(false, err.Error())
		return
	}
	var id int64
	var err error
	Rid, _ := this.GetInt64("Id")
	if Rid > 0 {
		id, err = m.UpdateRole(&r)
	} else {
		id, err = m.AddRole(&r)
	}
	if err == nil && id > 0 {
		content := fmt.Sprintf("roleId:%d", Rid)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "UpdateRole", "success", content)
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *RoleController) Delete() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	Id, _ := this.GetInt64("Id")
	err := m.DelUserRole(Id) //增加删除user_role中的角色
	if err != nil {
		this.Rsp(false, err.Error())
	}
	err2 := m.DelNodeRole(Id) //增加删除node_role中的角色
	if err2 != nil {
		this.Rsp(false, err.Error())
	}
	status, err := m.DelRoleById(Id)
	if err == nil && status > 0 {
		content := fmt.Sprintf("roleId:%d", Id)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "DelRole", "success", content)
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}

func (this *RoleController) Getlist() {
	roles, _ := m.GetRolelist(1, 1000, "Id")
	if len(roles) < 1 {
		roles = []orm.Params{}
	}
	this.Data["json"] = &roles
	this.ServeJSON()
}

func (this *RoleController) AccessToNode() {
	roleid, _ := this.GetInt64("Id")
	groupid, _ := this.GetInt64("group_id")
	nodes, count := m.GetNodelistByGroupid(groupid)
	list, _ := m.GetNodelistByRoleId(roleid)
	for i := 0; i < len(nodes); i++ {
		if nodes[i]["Pid"] != 0 {
			nodes[i]["_parentId"] = nodes[i]["Pid"]
		} else {
			nodes[i]["state"] = "closed"
		}
		for x := 0; x < len(list); x++ {
			if nodes[i]["Id"] == list[x]["Id"] {
				nodes[i]["checked"] = 1
			}
		}
	}
	if len(nodes) < 1 {
		nodes = []orm.Params{}
	}
	//this.Data["json"] = &map[string]interface{}{"total": count, "rows": &nodes}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &nodes}
	this.ServeJSON()

}

func (this *RoleController) AddAccess() {
	roleid, _ := this.GetInt64("roleid")
	//group_id, _ := this.GetInt64("group_id")
	//err := m.DelGroupNode(roleid, group_id)
	//if err != nil {
	//	this.Rsp(false, err.Error())
	//}

	//err := m.DelNodeRole(roleid) //用户在修改授权时需要删除旧的授权
	//if err != nil {
	//	this.Rsp(false, err.Error())
	//}

	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	ids := this.GetString("ids")
	nodeids := strings.Split(ids, ",")
	for _, v := range nodeids {
		id, _ := strconv.Atoi(v)
		_, err := m.AddRoleNode(roleid, int64(id))
		if err != nil {
			this.Rsp(false, err.Error())
		}
	}
	content := fmt.Sprintf("role:%d,path:%s\n", roleid, ids)
	_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "RoleGrant", "success", content)
	this.Rsp(true, "success")

}

func (this *RoleController) RoleToUserList() {
	roleid, _ := this.GetInt64("Id")
	//users, count := m.Getuserlist(1, 1000, "Id")
	users, count := m.GetUserByRoleId(roleid)
	/*
		for i := 0; i < len(users); i++ {
			for x := 0; x < len(list); x++ {
				if users[i]["Id"] == list[x]["Id"] {
					users[i]["checked"] = 1
				}
			}
		}
		if len(users) < 1 {
			users = []orm.Params{}
		}
	*/
	if len(users) < 1 {
		users = []orm.Params{}
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &users}
	this.ServeJSON()
}

// add
func (this *RoleController) RoleToNodeList() {
	roleid, _ := this.GetInt64("Id")
	nodes, count := m.GetNodeByRoleId(roleid)
	if len(nodes) < 1 {
		nodes = []orm.Params{}
	}
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &nodes}
	this.ServeJSON()
}

func (this *RoleController) AddRoleToUser() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	roleid, _ := this.GetInt64("Id")
	ids := this.GetString("ids")
	userids := strings.Split(ids, ",")
	//err := m.DelUserRole(roleid)
	//if err != nil {
	//	this.Rsp(false, err.Error())
	//}
	if len(ids) > 0 {
		for _, v := range userids {
			id, _ := strconv.Atoi(v)
			_, err := m.AddRoleUser(roleid, int64(id))
			if err != nil {
				this.Rsp(false, err.Error())
			}
		}
	}
	content := fmt.Sprintf("roleId:%d,userId:%s", roleid, userids)
	_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "AddUserInRole", "success", content)
	this.Rsp(true, "success")
}

func (this *RoleController) DelRoleToUser() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	roleid, _ := this.GetInt64("Id")
	ids := this.GetString("ids")
	userids := strings.Split(ids, ",")
	if len(ids) > 0 {
		for _, v := range userids {
			id, _ := strconv.Atoi(v)
			err := m.DelUserRoleByUserId(roleid, int64(id))
			if err != nil {
				this.Rsp(false, err.Error())
			}
		}
	}
	content := fmt.Sprintf("roleId:%d,userId:%s", roleid, userids)
	_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "DelUserFromRole", "success", content)
	this.Rsp(true, "success")
}

func (this *RoleController) DelRoleToNode() {
	roleid, _ := this.GetInt64("Id")
	ids := this.GetString("ids")

	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	nodeids := strings.Split(ids, ",")
	if len(ids) > 0 {
		for _, v := range nodeids {
			id, _ := strconv.Atoi(v)
			err := m.DelNodeRoleByNodeId(roleid, int64(id))
			if err != nil {
				this.Rsp(false, err.Error())
			}
		}
	}
	content := fmt.Sprintf("role:%d,path:%s\n", roleid, ids)
	_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "DelRoleGrant", "success", content)
	this.Rsp(true, "success")
}

// 编辑用户信息时根据userid获取角色
func (this *RoleController) GetRoleidByUserId() {
	userid, _ := this.GetInt64("Id")
	roleids := m.GetRoleByUserId(userid)
	if len(roleids) < 1 {
		roleids = []orm.Params{}
	}
	this.Data["json"] = &roleids
	this.ServeJSON()
	//return
}
