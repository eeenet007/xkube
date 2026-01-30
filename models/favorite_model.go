// Package models provides data models and related functions for favorites management
// 包models提供了收藏管理的数据模型和相关函数
package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Xkb_favorite represents a user's favorite item
// Xkb_favorite表示用户的收藏项
type Xkb_favorite struct { //查询时会将大写以下划线分开
	Id         int64  `json:"id"`         // 主键ID
	FavType    string `json:"fav_type"`   //类型：cluster,deploy,sts,cm,svc,ing
	ClusterId  string `json:"cluster_id"` //集群ID
	NameSpace  string `json:"name_space"` //命名空间
	FavName    string `json:"fav_name"`   //名称
	Username   string `json:"username"`   //用户名
	Createtime string `json:"createtime"` //创建时间
}

// init registers the Xkb_favorite model with the ORM
// init将Xkb_favorite模型注册到ORM中
func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Xkb_favorite))
}

// GetList_Favorite retrieves a list of favorites for a user with pagination
// GetList_Favorite为用户检索带分页的收藏列表
// Parameters:
//   - username: 用户名，为空则查询所有用户
//   - page: 页码
//   - page_size: 每页大小
//
// Returns:
//   - []Xkb_favorite: 收藏列表
//   - int64: 总数
func GetList_Favorite(username string, page, page_size int64) (dps []Xkb_favorite, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_favorite))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if username != "" {
		cond = cond.And("username", username)
	}
	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&dps)
	//qs.All(&xfr)
	count, _ = qs.Count()
	return dps, count
}

// Add_Favorite adds a new favorite item
// Add_Favorite添加一个新的收藏项
// Parameters:
//   - u: 要添加的收藏项
//
// Returns:
//   - int64: 插入的记录ID
//   - error: 错误信息
func Add_Favorite(u *Xkb_favorite) (int64, error) {
	o := orm.NewOrm()
	Ds := new(Xkb_favorite)
	Ds.FavType = u.FavType
	Ds.ClusterId = u.ClusterId
	Ds.NameSpace = u.NameSpace
	Ds.FavName = u.FavName
	Ds.Username = u.Username

	Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
	num, err := o.Insert(Ds)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// DelById_Favorite deletes a favorite item by ID and username
// DelById_Favorite根据ID和用户名删除收藏项
// Parameters:
//   - id: 收藏项ID
//   - username: 用户名
//
// Returns:
//   - int64: 删除的记录数
//   - error: 错误信息
func DelById_Favorite(id int64, username string) (int64, error) {
	o := orm.NewOrm()
	res, err := o.Raw("DELETE FROM xkb_favorite WHERE id=? AND username=?", id, username).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		return num, nil
	}
	return 0, err
	//status, err := o.Delete(&Xkb_favorite{Id: id, Username: username})
	//return status, err
}
