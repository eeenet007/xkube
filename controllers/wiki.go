// cluster.go
package controllers

import (
	"fmt"

	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
	m "xkube/models"

	beego "github.com/beego/beego/v2/server/web"
	//"github.com/tidwall/gjson"

	adm "xkube/xadmin/src/models"
)

type WikiController struct {
	beego.Controller
}

//列表
func (this *WikiController) List() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin", "count": 0}
		this.ServeJSON()
		return
	}

	id := this.GetString("id")
	xcolumn := this.GetString("xcolumn")
	key := this.GetString("key")
	page, _ := this.GetInt64("page")
	page_size, _ := this.GetInt64("limit")

	datas, count := m.GetList_Wiki(key, id, xcolumn, page, page_size)
	this.Data["json"] = &map[string]interface{}{"code": 0, "msg": "success", "count": count, "data": &datas}
	this.ServeJSON()
}

//添加
func (this *WikiController) Add() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Data["json"] = &map[string]interface{}{"code": -1, "msg": "NoLogin", "count": 0}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	title := this.GetString("title")
	//author := this.GetString("author")
	authkey := this.GetString("authkey")
	xcolumn := this.GetString("xcolumn")
	content := this.GetString("content")

	content = strings.Replace(content, "\\", "\\\\", -1)
	content = strings.Replace(content, "'", "\\'", -1)
	content = strings.Replace(content, "%", "%%", -1)

	currUser := userinfo.(adm.User).Username

	u := m.Wiki{
		Title:   title,
		Authkey: authkey,
		Author:  currUser,
		Xcolumn: xcolumn,
		Content: content,
		Sketch:  "",
	}

	id, err := m.Add_Wiki(&u)
	if err == nil && id > 0 {
		log.Printf("[INFO] Add Wiki success")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Add Wiki Fail:%s", err)
		resp := fmt.Sprintf(`{"code":-1,"msg":"%f"}`, err)
		this.Ctx.WriteString(resp)
	}
}

func (this *WikiController) Read() {
	Id, _ := this.GetInt64("id")
	key := this.GetString("key")
	var authkey string
	if this.Ctx.Request.Method == "POST" {
		//gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
		//authkey = gp.Get("authkey").String()
		authkey = this.GetString("authkey")
	}
	if key != "content" {
		this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	}

	resp, _ := m.ReadArticleById(Id, key, authkey)
	this.Ctx.WriteString(resp)
}

func (this *WikiController) ReadEncry() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	resp := m.ReadEncryCheck(Id)
	respBody := fmt.Sprintf(`{"result":"%s"}`, resp)
	this.Ctx.WriteString(respBody)
}

//删除
func (this *WikiController) Del() {
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	Id, _ := this.GetInt64("id")
	respid, err := m.DelById_Wiki(Id)
	if err == nil && respid > 0 {
		log.Printf("[INFO] Delsuccess")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Del Error:%s\n", err.Error())
		this.Ctx.WriteString(`{"code":-1,"msg":"Error"}`)
	}
}

//更新某个字段
// func (this *WikiController) Edit() {
// 	id := this.GetString("id")
// 	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
// 	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
// 	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)

// 	if gp.Get("name").Exists() {
// 		name := gp.Get("name").String()
// 		value := gp.Get("value").String()
// 		_, err := m.Edit_Cluster(id, name, value)
// 		if err == nil {
// 			this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
// 		} else {
// 			log.Printf("[ERROR] update Fail:%s\n", err)
// 			this.Ctx.WriteString(fmt.Sprintf(`{"code":-1,"msg":"%s"}`, err.Error()))
// 		}
// 	} else {
// 		this.Ctx.WriteString(`{"code":-1,"msg":"bodyParseFail"}`)
// 	}
// }

func (this *WikiController) Update() {

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")

	id, _ := this.GetInt64("id")
	title := this.GetString("title")
	//author := this.GetString("author")
	authkey := this.GetString("authkey")
	xcolumn := this.GetString("xcolumn")
	content := this.GetString("content")

	if content == "passwordError" {
		this.Ctx.WriteString(`{"code":-1,"msg":"passwordError"}`)
		return
	}

	content = strings.Replace(content, "\\", "\\\\", -1)
	content = strings.Replace(content, "'", "\\'", -1)
	content = strings.Replace(content, "%", "%%", -1)

	u := m.Wiki{
		Id:      id,
		Title:   title,
		Authkey: authkey,
		//Author:  author,
		Xcolumn: xcolumn,
		Content: content,
	}

	id, err := m.Update_Wiki(&u)
	if err == nil && id > 0 {
		log.Printf("[INFO] Update Wiki success")
		this.Ctx.WriteString(`{"code":0,"msg":"success"}`)
	} else {
		log.Printf("[ERROR] Update Wiki Fail:%s", err)
		this.Ctx.WriteString(`{"code":-1,"msg":"error"}`)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateFileFullPath() string {
	rootPath, _ := beego.AppConfig.String("rootPath")
	uploadPath, _ := beego.AppConfig.String("uploadPath")
	dir1 := strings.TrimRight(rootPath, "/") + "/" + strings.TrimRight(uploadPath, "/") + "/" + time.Now().Format("2006/01/02")
	if exist, _ := PathExists(dir1); !exist {
		_ = os.MkdirAll(dir1, os.ModePerm)
	}
	return strings.TrimRight(uploadPath, "/") + "/" + time.Now().Format("2006/01/02") + "/" + fmt.Sprintf("%d", time.Now().Unix())
}

func (this *WikiController) Upload() {
	file, head, err := this.GetFile("editormd-image-file") //此处名称为编辑器默认的name
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	//创建目录并返回不带文件后缀的路径
	fullPath := CreateFileFullPath() + path.Ext(head.Filename)

	rootPath, _ := beego.AppConfig.String("rootPath")

	//创建文件
	fW, err := os.Create(strings.TrimRight(rootPath, "/") + "/" + fullPath)
	if err != nil {
		fmt.Println("文件创建失败")
		this.Ctx.WriteString(`{"success":0,"msg":"create fail","url":""}`)
		return
	}
	defer fW.Close()
	_, err = io.Copy(fW, file)
	if err != nil {
		fmt.Println("文件保存失败")
		this.Ctx.WriteString(`{"success":0,"msg":"save fail","url":""}`)
		return
	}
	domain, _ := beego.AppConfig.String("domain")
	resp := fmt.Sprintf(`{"success":1,"msg":"upload success","url":"%s"}`, strings.TrimRight(domain, "/")+"/"+fullPath)
	this.Ctx.WriteString(resp)
}
