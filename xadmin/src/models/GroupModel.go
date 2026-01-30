// Package models provides data models and related functionality for the xkube admin system.
// 包models为xkube管理系统提供数据模型及相关功能。
package models

import (
	"errors"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	//beego "github.com/beego/beego/v2/server/web"
)

// Group represents a group entity for role-based access control.
// Group 表示基于角色访问控制的组实体。
type Group struct {
	Id     int64   `json:"id"`                                         // Id is the unique identifier of the group. Id是组的唯一标识符。
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`     // Name is the name of the group. Name是组的名称。
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`    // Title is the display title of the group. Title是组的显示标题。
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"` // Status indicates the status of the group (1 for active, 2 for inactive). Status表示组的状态（1为激活，2为未激活）。
	Sort   int     `orm:"default(1)" form:"Sort" valid:"Numeric"`      // Sort is the sorting order of the group. Sort是组的排序顺序。
	Nodes  []*Node `orm:"reverse(many)"`                               // Nodes is a collection of nodes associated with this group. Nodes是与此组关联的节点集合。
}

// func (g *Group) TableName() string {
// 	vTable, _ := beego.AppConfig.String("rbac_group_table")
// 	return vTable
// }

// init registers the Group model with the ORM framework.
// init 将Group模型注册到ORM框架中。
func init() {
	orm.RegisterModel(new(Group))
}

// checkGroup validates the Group entity using the Beego validation framework.
// checkGroup 使用Beego验证框架验证Group实体。
// g: pointer to the Group object to be validated, g: 指向待验证的Group对象的指针
// Returns: error if validation fails, nil otherwise. 返回: 如果验证失败则返回错误，否则返回nil。
func checkGroup(g *Group) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&g)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

// GetGrouplist retrieves a paginated list of groups with specified sorting.
// GetGrouplist 检索具有指定排序的组的分页列表。
// page: page number (starting from 1), page: 页码（从1开始）
// page_size: number of items per page, page_size: 每页项目数
// sort: sorting field and order (e.g., "-id" for descending by ID), sort: 排序字段和顺序（例如"-id"表示按ID降序）
// Returns: list of groups as orm.Params and total count, 返回: 以orm.Params形式的组列表和总数
func GetGrouplist(page int64, page_size int64, sort string) (groups []orm.Params, count int64) {
	o := orm.NewOrm()
	group := new(Group)
	qs := o.QueryTable(group)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&groups)
	count, _ = qs.Count()
	return groups, count
}

// AddGroup adds a new group to the database.
// AddGroup 向数据库添加一个新组。
// g: pointer to the Group object to be added, g: 指向要添加的Group对象的指针
// Returns: ID of the newly inserted group and error if any, 返回: 新插入组的ID和可能的错误
func AddGroup(g *Group) (int64, error) {
	if err := checkGroup(g); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	group := new(Group)
	group.Name = g.Name
	group.Title = g.Title
	group.Sort = g.Sort
	group.Status = g.Status
	id, err := o.Insert(group)
	return id, err
}

// UpdateGroup updates an existing group in the database.
// UpdateGroup 更新数据库中的现有组。
// g: pointer to the Group object containing updated information, g: 指向包含更新信息的Group对象的指针
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func UpdateGroup(g *Group) (int64, error) {
	if err := checkGroup(g); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	group := make(orm.Params)
	if len(g.Name) > 0 {
		group["Name"] = g.Name
	}
	if len(g.Title) > 0 {
		group["Title"] = g.Title
	}
	if g.Status != 0 {
		group["Status"] = g.Status
	}
	if g.Sort != 0 {
		group["Sort"] = g.Sort
	}
	if len(group) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Group
	num, err := o.QueryTable(table).Filter("Id", g.Id).Update(group)
	return num, err
}

// DelGroupById deletes a group by its ID.
// DelGroupById 根据ID删除组。
// Id: ID of the group to be deleted, Id: 要删除的组的ID
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func DelGroupById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Group{Id: Id})
	return status, err
}

// GroupList retrieves a list of all groups with only ID and Title fields.
// GroupList 检索所有组的列表，仅包含ID和Title字段。
// Returns: list of groups as orm.Params, 返回: 以orm.Params形式的组列表
func GroupList() (groups []orm.Params) {
	o := orm.NewOrm()
	group := new(Group)
	qs := o.QueryTable(group)
	qs.Values(&groups, "id", "title")
	return groups
}
