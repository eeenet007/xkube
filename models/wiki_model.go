// wikiModel.go
package models

import (
	//"bytes"
	"crypto/md5"
	//"errors"
	"fmt"

	//"io/ioutil"
	"log"
	"time"

	//"xkube/common"

	//"github.com/beego/beego/v2/core/validation"

	"github.com/beego/beego/v2/client/orm"
)

// Wiki - 文档中心模型结构体(Wiki - Document center model struct)
type Wiki struct {
	Id         int64  `json:"id"`         // 文档ID(Document ID)
	Title      string `json:"title"`      // 文档标题(Document title)
	Xcolumn    string `json:"xcolumn"`    // 文档分类(Document category)
	Sketch     string `json:"sketch"`     // 文档摘要(Document abstract)
	Content    string `json:"content"`    // 文档内容(Document content)
	Author     string `json:"author"`     // 文档作者(Document author)
	Status     int    `json:"status"`     // 文档状态(Document status)
	Authkey    string `json:"authkey"`    // 文档访问密钥(Document access key)
	Createtime string `json:"createtime"` // 创建时间(Creation time)
	Updatetime string `json:"updatetime"` // 更新时间(Update time)
}

// init - 初始化函数，注册Wiki模型(Init function, register Wiki model)
func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Wiki))
}

// GetList_Wiki - 获取Wiki文档列表(Get Wiki document list)
// key - 查询关键字(Query keyword)
// id - 文档ID(Document ID)
// xcolumn - 文档分类(Document category)
// page - 页码(Page number)
// page_size - 每页大小(Page size)
// 返回值：文档列表和总数(Return: document list and total count)
func GetList_Wiki(key, id, xcolumn string, page, page_size int64) (dps []Wiki, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Wiki))
	cond := orm.NewCondition()
	var offset int64
	// 计算分页偏移量(Calculate page offset)
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	// 如果查询关键字为column，则查询所有分类(If query keyword is column, query all categories)
	if key == "column" {
		qs.Distinct().All(&dps, "Xcolumn")
		count, _ = qs.Count()
		return dps, count
	}

	// 根据ID条件过滤(Filter by ID condition)
	if id != "" {
		cond = cond.And("id", id)
	}

	// 根据分类条件过滤(Filter by category condition)
	if xcolumn != "" {
		cond = cond.And("xcolumn", xcolumn)
	}

	qs = qs.SetCond(cond)
	// 按更新时间倒序排列并分页查询(Order by update time descending and query with pagination)
	qs.Limit(page_size, offset).OrderBy("-updatetime").All(&dps, "Id", "Title", "Xcolumn", "Sketch", "Author", "Authkey", "Status", "Createtime", "Updatetime")
	//qs.All(&xfr)
	//fmt.Println(dps)
	count, _ = qs.Count()
	var ttt = make([]Wiki, 0)
	// 处理访问密钥显示(Processing access key display)
	for _, vv := range dps {
		authkey := "false"
		if vv.Authkey != "" {
			authkey = "true"
		}
		ttt = append(ttt, Wiki{
			Id:         vv.Id,
			Title:      vv.Title,
			Xcolumn:    vv.Xcolumn,
			Sketch:     vv.Sketch,
			Author:     vv.Author,
			Status:     vv.Status,
			Authkey:    authkey,
			Createtime: vv.Createtime,
			Updatetime: vv.Updatetime,
		})
	}
	return ttt, count
}

// Add_Wiki - 添加Wiki文档(Add Wiki document)
// u - Wiki文档结构体指针(Wiki document struct pointer)
// 返回值：插入的ID和错误信息(Return: inserted ID and error information)
func Add_Wiki(u *Wiki) (int64, error) {
	o := orm.NewOrm()
	Ds := new(Wiki)
	Ds.Title = u.Title     // 设置文档标题(Set document title)
	Ds.Xcolumn = u.Xcolumn // 设置文档分类(Set document category)
	Ds.Sketch = u.Sketch   // 设置文档摘要(Set document abstract)
	Ds.Content = u.Content // 设置文档内容(Set document content)
	Ds.Author = u.Author   // 设置文档作者(Set document author)
	// 如果访问密钥不为空且不是默认值，则进行MD5加密(If access key is not empty and not default value, perform MD5 encryption)
	if u.Authkey != "" && u.Authkey != "******" {
		Ds.Authkey = fmt.Sprintf("%x", md5.Sum([]byte(u.Authkey)))
	}
	// 设置创建和更新时间(Set creation and update time)
	Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
	Ds.Updatetime = time.Now().Format("2006-01-02 15:04:05")
	num, err := o.Insert(Ds)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// DelById_Wiki - 根据ID删除Wiki文档(Delete Wiki document by ID)
// id - 文档ID(Document ID)
// 返回值：删除状态和错误信息(Return: delete status and error information)
func DelById_Wiki(id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Wiki{Id: id})
	return status, err
}

// Edit_Wiki - 编辑Wiki文档(Edit Wiki document)
// id - 文档ID(Document ID)
// key - 要更新的字段名(Field name to update)
// value - 要更新的字段值(Field value to update)
// 返回值：影响的行数和错误信息(Return: number of affected rows and error information)
func Edit_Wiki(id, key, value string) (int64, error) {
	o := orm.NewOrm()
	// 构造SQL更新语句(Construct SQL update statement)
	sqlstr := fmt.Sprintf("UPDATE Wiki SET %s = '%s' WHERE id = %s", key, value, id)
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("[INFO] EditCluster affected nums:%d\n", num)
		return num, nil
	}
	return 0, err
}

// Update_Wiki - 更新Wiki文档(Update Wiki document)
// u - Wiki文档结构体指针(Wiki document struct pointer)
// 返回值：影响的行数和错误信息(Return: number of affected rows and error information)
func Update_Wiki(u *Wiki) (int64, error) {
	o := orm.NewOrm()
	op := make(orm.Params)
	// 设置要更新的字段(Set fields to update)
	op["Id"] = u.Id
	op["Title"] = u.Title
	op["Xcolumn"] = u.Xcolumn
	//op["Sketch"] = u.Sketch
	op["Content"] = u.Content
	//op["Author"] = u.Author
	// 如果访问密钥不为空且不是默认值，则进行MD5加密(If access key is not empty and not default value, perform MD5 encryption)
	if u.Authkey != "" && u.Authkey != "******" {
		op["Authkey"] = fmt.Sprintf("%x", md5.Sum([]byte(u.Authkey)))
	} else {
		op["Authkey"] = ""
	}
	// 更新时间(Update time)
	op["Updatetime"] = time.Now().Format("2006-01-02 15:04:05")

	var table Wiki
	// 执行更新操作(Execute update operation)
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(op)

	return num, err
}

// ReadArticleById - 根据ID读取文章(Read article by ID)
// id - 文档ID(Document ID)
// key - 查询字段(Query field)
// authkey - 访问密钥(Access key)
// 返回值：查询结果和错误信息(Return: query result and error information)
func ReadArticleById(id int64, key, authkey string) (string, error) {
	o := orm.NewOrm()
	// 如果访问密钥不为空，则进行MD5加密(If access key is not empty, perform MD5 encryption)
	if authkey != "" {
		authkey = fmt.Sprintf("%x", md5.Sum([]byte(authkey)))
	}
	var wk Wiki
	// 根据内容和密钥查询(QUERY by content and key)
	if key == "content" {
		err := o.Raw("SELECT content from wiki where id = ? AND authkey = ?", id, authkey).QueryRow(&wk)
		if err != nil {
			log.Printf("[ERROR] ReadArticleById Exec:%s\n", err)
			return "passwordError", err
		}
		return wk.Content, nil
	} else {
		// 查询文档基本信息(Query basic document information)
		err := o.Raw("SELECT title,author,xcolumn,authkey from wiki where id = ?", id).QueryRow(&wk)
		if err != nil {
			log.Printf("[ERROR] ReadArticleById2 Exec:%s\n", err)
			return "noFound", err
		}
		authkey := "false"
		if wk.Authkey != "" {
			authkey = "true"
		}
		return fmt.Sprintf(`{"title":"%s","author":"%s","xcolumn":"%s","authkey":"%s"}`, wk.Title, wk.Author, wk.Xcolumn, authkey), nil
	}
}

// ReadEncryCheck - 检查文档是否加密(Check if document is encrypted)
// id - 文档ID(Document ID)
// 返回值：加密状态(Return: encryption status)
func ReadEncryCheck(id int64) string {
	o := orm.NewOrm()
	var wk Wiki
	// 查询文档的访问密钥(Query document access key)
	err := o.Raw("SELECT authkey from wiki where id = ?", id).QueryRow(&wk)
	if err != nil {
		log.Printf("[ERROR] ReadArticleById Exec:%s\n", err)
		return "fail"
	}
	//log.Println(wk.Authkey)
	//log.Println(len(wk.Authkey))
	// 判断密钥长度确定是否加密(Determine if encrypted by key length)
	if len(wk.Authkey) > 30 {
		return "true"
	}
	return "false"
}
