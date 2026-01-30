// Package models provides data models and related functionality for the xkube admin system.
// 包models为xkube管理系统提供数据模型及相关功能。
package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	//beego "github.com/beego/beego/v2/server/web"
)

// Node represents a URL path in the system for role-based access control.
// Node 表示系统中用于基于角色访问控制的URL路径。
type Node struct {
	Id     int64  `json:"id"`                                               // Id is the unique identifier of the node. Id是节点的唯一标识符。
	Title  string `orm:"size(100)" form:"Title"  valid:"Required"`          // Title is the display title of the URL path. Title是URL路径的显示标题。
	Name   string `orm:"size(100)" form:"Name"  valid:"Required"`           // Name is the actual URL path string. Name是实际的URL路径字符串。
	Level  int    `orm:"default(1)" form:"Level"  valid:"Required"`         // Level indicates the hierarchy level of the node (1 for main menu, 2 for submenu, etc.). Level表示节点的层级（1为主菜单，2为子菜单等）。
	Pid    int64  `form:"Pid"  valid:"Required"`                            // Pid is the parent node ID. Pid是父节点ID。
	Remark string `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"` // Remark is additional information about the node. Remark是关于节点的附加信息。
	Icons  string `orm:"null;size(200)" form:"Icons" valid:"MaxSize(200)"`  // Icons is the icon class for the node. Icons是节点的图标类。
	Status int    `orm:"default(2)" form:"Status" valid:"Range(1,2)"`       // Status indicates the status of the node (1 for active, 2 for inactive). Status表示节点的状态（1为激活，2为未激活）。
	Sorts  int64  `form:"Sorts"`                                            // Sorts is the sorting order of the node. Sorts是节点的排序顺序。
	//Groupid int64   `form:"-"`
	Group *Group  `orm:"rel(fk)"`  // Group is the group that this node belongs to. Group是此节点所属的组。
	Role  []*Role `orm:"rel(m2m)"` // Role is the list of roles that have access to this node. Role是有权访问此节点的角色列表。
}

// init registers the Node model with the ORM framework.
// init 将Node模型注册到ORM框架中。
func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Node))
}

// func (n *Node) TableName() string {
// 	vTable, _ := beego.AppConfig.String("rbac_node_table")
// 	return vTable
// }

// checkNode validates the Node entity using the Beego validation framework.
// checkNode 使用Beego验证框架验证Node实体。
// u: pointer to the Node object to be validated, u: 指向待验证的Node对象的指针
// Returns: error if validation fails, nil otherwise. 返回: 如果验证失败则返回错误，否则返回nil。
func checkNode(u *Node) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

// GetNodelist retrieves a paginated list of nodes with specified sorting.
// GetNodelist 检索具有指定排序的节点分页列表。
// page: page number (starting from 1), page: 页码（从1开始）
// page_size: number of items per page, page_size: 每页项目数
// sort: sorting field and order (e.g., "-id" for descending by ID), sort: 排序字段和顺序（例如"-id"表示按ID降序）
// Returns: list of nodes as orm.Params and total count, 返回: 以orm.Params形式的节点列表和总数
func GetNodelist(page int64, page_size int64, sort string) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	qs := o.QueryTable(node)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&nodes, "Id", "Title", "Name", "Status", "Level", "Pid", "Icons", "Sorts", "Remark", "Group__id")
	//qs.Limit(page_size, offset).OrderBy(sort).Values(&nodes)
	count, _ = qs.Count()
	return nodes, count
}

// ReadNode retrieves a node by its ID.
// ReadNode 根据ID检索节点。
// nid: ID of the node to retrieve, nid: 要检索的节点ID
// Returns: Node object and error if any, 返回: Node对象和可能的错误
func ReadNode(nid int64) (Node, error) {
	o := orm.NewOrm()
	node := Node{Id: nid}
	err := o.Read(&node)
	if err != nil {
		return node, err
	}
	return node, nil
}

// AddNode adds a new URL path node to the database.
// AddNode 向数据库添加新的URL路径节点。
// n: pointer to the Node object to be added, n: 指向要添加的Node对象的指针
// Returns: ID of the newly inserted node and error if any, 返回: 新插入节点的ID和可能的错误
func AddNode(n *Node) (int64, error) {
	if err := checkNode(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := new(Node)
	node.Title = n.Title
	node.Name = n.Name
	node.Level = n.Level
	node.Pid = n.Pid
	node.Remark = n.Remark
	node.Status = n.Status
	//node.Groupid = n.Groupid
	node.Group = n.Group

	id, err := o.Insert(node)
	return id, err
}

// UpdateNode updates an existing URL path node in the database.
// UpdateNode 更新数据库中的现有URL路径节点。
// n: pointer to the Node object containing updated information, n: 指向包含更新信息的Node对象的指针
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func UpdateNode(n *Node) (int64, error) {
	if err := checkNode(n); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	node := make(orm.Params)
	if len(n.Title) > 0 {
		node["Title"] = n.Title
	}
	if len(n.Name) > 0 {
		node["Name"] = n.Name
	}
	if len(n.Remark) > 0 {
		node["Remark"] = n.Remark
	}
	if len(n.Icons) > 0 {
		node["Icons"] = n.Icons
	}
	if n.Level != 0 {
		node["Level"] = n.Level
	}
	if n.Pid >= 0 {
		node["Pid"] = n.Pid
	}
	if n.Sorts != 0 {
		node["Sorts"] = n.Sorts
	}
	if n.Status != 0 {
		node["Status"] = n.Status
	}
	if len(node) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Node
	num, err := o.QueryTable(table).Filter("Id", n.Id).Update(node)
	return num, err
}

// EditNode updates a specific field of a node.
// EditNode 更新节点的特定字段。
// id: ID of the node to update, id: 要更新的节点ID
// key: field name to update, key: 要更新的字段名
// value: new value for the field, value: 字段的新值
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func EditNode(id, key, value string) (int64, error) {
	o := orm.NewOrm()
	if key == "group__id" {
		key = "group_id"
	}
	sqlstr := fmt.Sprintf("UPDATE node SET %s = '%s' WHERE  id = %s", key, value, id)
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("[INFO] EditApp affected nums:%d\n", num)
		return num, nil
	}
	return 0, err
}

// DelNodeById deletes a node by its ID.
// DelNodeById 根据ID删除节点。
// Id: ID of the node to be deleted, Id: 要删除的节点ID
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func DelNodeById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Node{Id: Id})
	return status, err
}

// GetNodelistByGroupid retrieves a list of nodes belonging to a specific group.
// GetNodelistByGroupid 检索属于特定组的节点列表。
// Groupid: ID of the group to retrieve nodes for, Groupid: 要检索节点的组ID
// Returns: list of nodes as orm.Params and total count, 返回: 以orm.Params形式的节点列表和总数
func GetNodelistByGroupid(Groupid int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Group", Groupid).RelatedSel().Values(&nodes)
	//count, _ = o.QueryTable(node).Filter("Group_id", Groupid).Values(&nodes)
	return nodes, count
}

// GetNodeTree retrieves nodes of a specific parent and level for building navigation tree.
// GetNodeTree 检索特定父节点和层级的节点以构建导航树。
// pid: parent node ID, pid: 父节点ID
// level: node level, level: 节点层级
// Returns: list of nodes as orm.Params and error if any, 返回: 以orm.Params形式的节点列表和可能的错误
func GetNodeTree(pid int64, level int64) ([]orm.Params, error) {
	o := orm.NewOrm()
	node := new(Node)
	var nodes []orm.Params
	_, err := o.QueryTable(node).Filter("Pid", pid).Filter("Level", level).Filter("Status", 2).Values(&nodes)
	if err != nil {
		return nodes, err
	}
	return nodes, nil
}

// DelNodeRole deletes the relationship between a node and a role.
// DelNodeRole 删除节点和角色之间的关系。
// roleid: ID of the role to remove from nodes, roleid: 要从节点中移除的角色ID
// Returns: error if any, 返回: 可能的错误
func DelNodeRole(roleid int64) error { //add by kang @20181121
	o := orm.NewOrm()
	_, err := o.QueryTable("node_roles").Filter("role_id", roleid).Delete()
	return err
}

// GetNodePid retrieves parent nodes with level less than 3.
// GetNodePid 检索层级小于3的父节点。
// Returns: list of nodes as orm.Params and total count, 返回: 以orm.Params形式的节点列表和总数
func GetNodePid() (nodes []orm.Params, count int64) { //add by kang@20181122
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Level__lt", 3).Distinct().RelatedSel().Values(&nodes, "Id", "Title", "Level")
	//count, _ = o.QueryTable(node).Filter("Level__lt", 3).Distinct().Values(&nodes, "Id", "Title", "Level")
	return nodes, count
}
