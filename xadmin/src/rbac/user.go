package rbac

import (
	"fmt"
	"strconv"
	"strings"

	//admlib "xkube/xadmin/src/lib"

	m "xkube/xadmin/src/models"
)

type UserController struct {
	CommonController
}

func (this *UserController) List() {

	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin", "count": 0}
		this.ServeJSON()
		return
	}

	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")

	sort := this.GetString("sort")
	order := this.GetString("order")
	if len(order) > 0 {
		if order == "desc" {
			sort = "-" + sort
		}
	} else {
		sort = "Id"
	}
	users, count := m.Getuserlist(page, page_size, sort)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "", "count": count, "data": &users}
	this.ServeJSON()
}

func (this *UserController) Add() {
	u := m.User{}
	Roleids := this.GetString("Roleid")
	userinfo := this.GetSession("userinfo")

	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	roleidarr := strings.Split(Roleids, ",")
	if err := this.ParseForm(&u); err != nil {
		this.Rsp(false, err.Error())
		return
	}

	id, err := m.AddUser(&u)
	if err == nil && id > 0 {
		if len(roleidarr) > 0 {
			for _, v := range roleidarr {
				rid, _ := strconv.Atoi(v)
				_, err := m.AddRoleUser(int64(rid), id)
				if err != nil {
					this.Rsp(false, err.Error())
					return
				}
			}
		}
		content := fmt.Sprintf("adduser:%s,roleId:%s", u.Username, Roleids)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "AddUser", "success", content)
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}

}

func (this *UserController) Update() {
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	if this.Ctx.Request.Method == "POST" {
		u := m.User{}
		if err := this.ParseForm(&u); err != nil {
			//handle error
			this.Rsp(false, err.Error())
			return
		}

		//更新角色
		Id, _ := this.GetInt64("Id")
		Roleids := this.GetString("Roleid")
		roleidarr := strings.Split(Roleids, ",")

		//userinfo := this.GetSession("userinfo")

		//删除user_roles中的用户角色 当用户更改不授权角色时需要删除对应的权限
		err1 := m.DelUserRoleByUserId(0, Id)
		if err1 != nil {
			this.Rsp(false, err1.Error())
		}

		if len(roleidarr) > 0 {
			for _, v := range roleidarr {
				rid, _ := strconv.Atoi(v)
				if rid > 0 {
					_, err := m.AddRoleUser(int64(rid), Id)
					if err != nil {
						this.Rsp(false, err.Error())
						return
					}
				}
			}
		}
		//更新用户信息
		//id, err := m.UpdateUser(&u)
		_, err := m.UpdateUser(&u)
		//if err == nil && id > 0 {
		if err == nil {
			//log record start
			content := fmt.Sprintf("user:%s,roleId:%s", u.Username, Roleids)
			_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "UpdateUser", "success", content)
			//log record end
			this.Rsp(true, "Success")
			return
		} else {
			this.Rsp(false, "UpdateFail")
			return
		}
	} else {
		Id, _ := this.GetInt64("Id")
		user, _ := m.GetUserById(Id)
		this.Data["json"] = &map[string]interface{}{"data": &user}
		this.ServeJSON()
	}
}

func (this *UserController) Delete() {
	Id, _ := this.GetInt64("Id")
	//userinfo := this.GetSession("userinfo")
	if Id == 1 { //admin不允许删除
		this.Rsp(false, "NotAllowDel")
		return
	}
	if Id == 0 { //admin不允许删除
		this.Rsp(false, "errorId")
		return
	}
	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	status, err := m.DelUserById(Id)
	if err == nil && status > 0 {
		//删除user_roles中的用户角色
		err1 := m.DelUserRoleByUserId(0, Id)
		if err1 != nil {
			this.Rsp(false, err1.Error())
			return
		}
		content := fmt.Sprintf("userId:%d", Id)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "DelUser", "success", content)
		this.Rsp(true, "Success")
		return
	} else {
		this.Rsp(false, err.Error())
		return
	}
}
