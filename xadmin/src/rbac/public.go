package rbac

import (
	//"fmt"
	//"lib"
	//"errors"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
	"time"

	. "xkube/xadmin/src"
	//admlib "xkube/xadmin/src/lib"
	m "xkube/xadmin/src/models"

	common "xkube/common"

	"github.com/beego/beego/v2/client/cache"
	//"github.com/beego/beego/v2/client/orm"
	_ "github.com/beego/beego/v2/client/cache/redis"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/captcha"
	"github.com/tidwall/gjson"
)

type MainController struct {
	CommonController
}

type Tree struct {
	Id         int64      `json:"id"`
	Text       string     `json:"text"`
	IconCls    string     `json:"iconCls"`
	Sorts      int64      `json:"sorts"`
	Checked    string     `json:"checked"`
	State      string     `json:"state"`
	Children   []Tree     `json:"children"`
	Attributes Attributes `json:"attributes"`
}
type Attributes struct {
	Url   string `json:"url"`
	Price int64  `json:"price"`
}

var cpt *captcha.Captcha

var MaxLoginFailUser, _ = beego.AppConfig.Int64("max_login_fail_user")
var MaxLoginFailIp, _ = beego.AppConfig.Int64("max_login_fail_ip")
var UserLockTime, _ = beego.AppConfig.Int64("user_lock_time")
var IpLockTime, _ = beego.AppConfig.Int64("ip_lock_time")
var MobileVerifyCode, _ = beego.AppConfig.Bool("mobile_verify_code")

var redis_host, _ = beego.AppConfig.String("redisDb")
var redis_passwd, _ = beego.AppConfig.String("redisPasswd")

func init() {
	//store := cache.NewMemoryCache()
	store, err := cache.NewCache("redis", fmt.Sprintf(`{"key":"cptcache","conn":"redis://%s@%s","dbNum":"0"}`, redis_passwd, redis_host))
	if err != nil {
		log.Printf("[ERROR] captcha NewCache error:%s\n", err.Error())
	}
	cpt = captcha.NewWithFilter("/public/captcha/", store) //一定要写在构造函数里面，要不然第一次打开页面有可能是X
	// 设置验证码长度
	cpt.ChallengeNums = 5
	// 设置验证码模板高度
	cpt.StdHeight = 40
	// 设置验证码模板宽度
	cpt.StdWidth = 120

	gob.Register(m.User{})
}

// 后台首页
func (this *MainController) Index() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		authGateWay, _ := beego.AppConfig.String("rbac_auth_gateway")
		this.Ctx.Redirect(302, authGateWay)
		return
	}
	this.Data["userinfo"] = userinfo
	this.TplName = "front/xkube_index.html"
}

// 登录
func (this *MainController) Login() {
	isajax := this.GetString("isajax")
	// log.Println(string(this.Ctx.Input.RequestBody))
	if isajax == "1" {

		username := this.GetString("username")
		password := this.GetString("password")
		src := this.GetString("src")
		if src != "xkubeApp" { //app登录不做验证码校验
			if MobileVerifyCode {
				telcode := this.GetString("telcode")
				telCodeOk := m.VerifyTelCode(username, telcode)
				if !telCodeOk {
					this.Rsp(false, "验证码错误")
					return
				}
			} else {
				cptchaOk := cpt.VerifyReq(this.Ctx.Request)
				if !cptchaOk {
					this.Rsp(false, "验证码错误")
					return
				}
			}
		}

		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}
		//判断用户是否被锁
		if common.Get("lock_"+username) == "lock" {
			this.Rsp(false, "max login fail,username locked")
			return
		}
		//判断IP是否被锁
		realIp := strings.Split(userip, ":")[0]
		if common.Get("lock_"+realIp) == "lock" {
			this.Rsp(false, "max login fail,ip locked")
			return
		}

		var ctx context.Context
		sessId := this.Ctx.Input.CruSession.SessionID(ctx)
		user, err := CheckLogin(username, password)
		if err == nil {
			//userinfo := this.GetSession("userinfo")
			//log.Println(userinfo)
			this.SetSession("userinfo", user)
			accesslist, _ := GetAccessList(user.Id)
			this.SetSession("accesslist", accesslist)

			clusterlist, _ := GetClusterList(user.Username)
			this.SetSession("clusterlist", clusterlist)

			//log record start
			//admlib.SecurityLog(userip, username, "AdmLogin", "ok")
			_ = m.InsertLogAudit(userip, username, "login", "success", "")
			//log record end
			logintime := time.Now().Format("2006-01-02T15:04:05")
			m.UpdateLoginIpTime(user.Id, logintime, userip)
			//this.Rsp(true, "登录成功")

			//登录成功清除用户名和IP的计数
			_ = common.Del("lock_" + username)
			_ = common.Del("lock_" + realIp)

			this.Data["json"] = &map[string]interface{}{"status": true, "msg": "登录成功", "sessionId": sessId}
			this.ServeJSON()
			return
		} else {
			//admlib.SecurityLog(userip, username, "AdmLogin", "fail")
			_ = m.InsertLogAudit(userip, username, "login", "fail", "")
			//超出失败次数锁定用户名
			num, _ := common.IncrBy("lock_"+username, 1)
			if num >= MaxLoginFailUser {
				_ = common.SetEx("lock_"+username, "lock", UserLockTime)
			}
			//超出失败次数锁定IP
			num1, _ := common.IncrBy("lock_"+realIp, 1)
			if num1 >= MaxLoginFailIp {
				_ = common.SetEx("lock_"+realIp, "lock", IpLockTime)
			}
			this.Rsp(false, err.Error())
			return
		}
	}

	userinfo := this.GetSession("userinfo")
	if userinfo != nil {
		this.Ctx.Redirect(302, "/index")
	}

	var loginTpl = "/login_code.html"
	if MobileVerifyCode {
		loginTpl = "/login_telcode.html"
	}

	this.TplName = this.GetTemplatetype() + loginTpl
}

// app登录
func (this *MainController) AppLogin() {
	isajax := this.GetString("isajax")
	// log.Println(string(this.Ctx.Input.RequestBody))
	if isajax == "1" {
		username := this.GetString("username")
		password := this.GetString("password")
		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}
		//判断用户是否被锁
		if common.Get("lock_"+username) == "lock" {
			this.Rsp(false, "max login fail,username locked")
			return
		}
		//判断IP是否被锁
		realIp := strings.Split(userip, ":")[0]
		if common.Get("lock_"+realIp) == "lock" {
			this.Rsp(false, "max login fail,ip locked")
			return
		}

		var ctx context.Context
		sessId := this.Ctx.Input.CruSession.SessionID(ctx)
		user, err := CheckLogin(username, password)
		if err == nil {
			//userinfo := this.GetSession("userinfo")
			//log.Println(userinfo)
			this.SetSession("userinfo", user)
			accesslist, _ := GetAccessList(user.Id)
			this.SetSession("accesslist", accesslist)

			clusterlist, _ := GetClusterList(user.Username)
			this.SetSession("clusterlist", clusterlist)

			//log record start
			//admlib.SecurityLog(userip, username, "AdmLogin", "ok")
			_ = m.InsertLogAudit(userip, username, "login", "success", "app login")
			//log record end
			logintime := time.Now().Format("2006-01-02T15:04:05")
			m.UpdateLoginIpTime(user.Id, logintime, userip)
			//this.Rsp(true, "登录成功")

			//登录成功清除用户名和IP的计数
			_ = common.Del("lock_" + username)
			_ = common.Del("lock_" + realIp)

			this.Data["json"] = &map[string]interface{}{"status": true, "msg": "登录成功", "sessionId": sessId}
			this.ServeJSON()
			return
		} else {
			//admlib.SecurityLog(userip, username, "AdmLogin", "fail")
			_ = m.InsertLogAudit(userip, username, "login", "fail", "app login")
			//超出失败次数锁定用户名
			num, _ := common.IncrBy("lock_"+username, 1)
			if num >= MaxLoginFailUser {
				_ = common.SetEx("lock_"+username, "lock", UserLockTime)
			}
			//超出失败次数锁定IP
			num1, _ := common.IncrBy("lock_"+realIp, 1)
			if num1 >= MaxLoginFailIp {
				_ = common.SetEx("lock_"+realIp, "lock", IpLockTime)
			}
			this.Rsp(false, err.Error())
			return
		}
	}

	this.Rsp(false, "ajaxFail")
}

// 退出
func (this *MainController) Logout() {
	isajax := this.GetString("isajax")

	userinfo := this.GetSession("userinfo")
	//log record start
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}
	_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "Logout", "success", "")
	this.DelSession("userinfo")
	if isajax == "1" {
		this.Rsp(true, "success")
	} else {
		//this.Ctx.Redirect(302, "/login.html")
		rbac_auth_gateway, _ := beego.AppConfig.String("rbac_auth_gateway")
		this.Ctx.Redirect(302, rbac_auth_gateway)
	}
}

// 修改密码
func (this *MainController) Changepwd() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		authGateway, _ := beego.AppConfig.String("rbac_auth_gateway")
		this.Ctx.Redirect(302, authGateway)
	}
	if this.IsAjax() {
		oldpassword := this.GetString("oldpassword")
		newpassword := this.GetString("newpassword")
		repeatpassword := this.GetString("repeatpassword")
		if newpassword != repeatpassword {
			this.Rsp(false, "两次输入密码不一致")
		}

		//log record start
		userip := this.Ctx.Request.Header.Get("x-forwarded-for")
		if userip == "" {
			userip = this.Ctx.Request.RemoteAddr
		}

		user, err := CheckLogin(userinfo.(m.User).Username, oldpassword)
		if err == nil {
			var u m.User
			u.Id = user.Id
			u.Password = newpassword
			id, err := m.UpdateUser(&u)
			if err == nil && id > 0 {
				_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "updatePassword", "success", "")
				this.Rsp(true, "密码修改成功")
				return
			} else {
				_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "updatePassword", "fail", "")
				this.Rsp(false, err.Error())
				return
			}
		}
		this.Rsp(false, "密码有误")
	}
}

// 个人信息
func (this *MainController) MyInfo() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		authGateway, _ := beego.AppConfig.String("rbac_auth_gateway")
		this.Ctx.Redirect(302, authGateway)
	}
	if this.IsAjax() {
		user, _ := m.GetUserById(userinfo.(m.User).Id)
		this.Data["json"] = &map[string]interface{}{"data": &user}
		this.ServeJSON()
	} else {
		this.TplName = this.GetTemplatetype() + "/myinfo.html"
	}
}

func (this *MainController) UpdateMyInfo() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		authGateway, _ := beego.AppConfig.String("rbac_auth_gateway")
		this.Ctx.Redirect(302, authGateway)
	}
	if this.Ctx.Request.Method == "POST" && this.IsAjax() {
		u := m.User{}
		if err := this.ParseForm(&u); err != nil {
			//handle error
			this.Rsp(false, err.Error())
			return
		}

		if userinfo.(m.User).Username == u.Username && userinfo.(m.User).Id == u.Id {
			//更新用户信息
			id, err := m.UpdateUser(&u)
			if err == nil && id > 0 {
				this.Rsp(true, "Success")
				return
			} else {
				this.Rsp(false, "UpdateFail")
				return
			}
		} else {
			this.Rsp(false, "ErrorId")
			return
		}
	}
	this.Rsp(false, "Error")
}

// 检测app登录状态
func (this *MainController) CheckAppLogin() {
	//cookeToken := this.Ctx.GetCookie("BsessionId")
	//var ctx context.Context
	//log.Println(this.Ctx.Input.CruSession.SessionID(ctx))
	userinfo := this.GetSession("userinfo")
	if userinfo != nil {
		this.Rsp(true, "ok")
	} else {
		this.Rsp(false, "fail")
	}
}

func (this *MainController) SendCode() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "*")
	// log.Println(string(this.Ctx.Input.RequestBody))

	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	username := gjson.Get(string(this.Ctx.Input.RequestBody), "username").String()
	password := gjson.Get(string(this.Ctx.Input.RequestBody), "password").String()
	//判断用户是否被锁
	if common.Get("lock_"+username) == "lock" {
		this.Rsp(false, "max login fail,username locked")
		return
	}
	//判断IP是否被锁
	realIp := strings.Split(userip, ":")[0]
	if common.Get("lock_"+realIp) == "lock" {
		this.Rsp(false, "max login fail,ip locked")
		return
	}
	_, err := CheckLogin(username, password)
	if err == nil {
		err := m.SendSmsV2(username)
		if err != nil {
			//admlib.SecurityLog(userip, username, "sendCode", "fail")
			_ = m.InsertLogAudit(userip, username, "sendCode", "fail", "code login")
			//超出失败次数锁定用户名
			num, _ := common.IncrBy(username, 1)
			if num >= MaxLoginFailUser {
				_ = common.SetEx("lock_"+username, "lock", UserLockTime)
			}
			//超出失败次数锁定IP
			num1, _ := common.IncrBy(realIp, 1)
			if num1 >= MaxLoginFailIp {
				_ = common.SetEx("lock_"+realIp, "lock", IpLockTime)
			}

			this.Rsp(false, fmt.Sprintf("%s", err))
			return
		}
		//admlib.SecurityLog(userip, username, "sendCode", "success")
		_ = m.InsertLogAudit(userip, username, "sendCode", "success", "code login")
	} else {
		//admlib.SecurityLog(userip, username, "sendCode:userPassword", "fail")
		_ = m.InsertLogAudit(userip, username, "sendCode", "fail", "code login and userPassword")
		//超出失败次数锁定用户名
		num, _ := common.IncrBy("lock_"+username, 1)
		if num >= MaxLoginFailUser {
			_ = common.SetEx("lock_"+username, "lock", UserLockTime)
		}
		//超出失败次数锁定IP
		num1, _ := common.IncrBy("lock_"+realIp, 1)
		if num1 >= MaxLoginFailIp {
			_ = common.SetEx("lock_"+realIp, "lock", IpLockTime)
		}
		this.Rsp(false, "用户名或密码错误")
		return
	}
	//登录成功清除用户名和IP的计数
	_ = common.Del("lock_" + username)
	_ = common.Del("lock_" + realIp)
	this.Rsp(true, "发送成功,3分钟内有效")
}
