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

// Role represents a role in the role-based access control system.
// Role 表示基于角色访问控制系统中的角色。
type Role struct {
	Id     int64   `json:"id"`                           // Id is the unique identifier of the role. Id是角色的唯一标识符。
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`  // Title is the display title of the role. Title是角色的显示标题。
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`   // Name is the name identifier of the role. Name是角色的名称标识符。
	Remark string  `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"` // Remark is additional information about the role. Remark是关于角色的附加信息。
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"`       // Status indicates the status of the role (1 for active, 2 for inactive). Status表示角色的状态（1为激活，2为未激活）。
	Node   []*Node `orm:"reverse(many)"`                 // Node is the list of nodes (URL paths) associated with this role. Node是与此角色关联的节点（URL路径）列表。
	User   []*User `orm:"reverse(many)"`                 // User is the list of users assigned to this role. User是分配给此角色的用户列表。
}

// func (r *Role) TableName() string {
// 	vTable, _ := beego.AppConfig.String("rbac_role_table")
// 	return vTable
// }

// init registers the Role model with the ORM framework.
// init 将Role模型注册到ORM框架中。
func init() {
	orm.RegisterModel(new(Role))
}

// checkRole validates the Role entity using the Beego validation framework.
// checkRole 使用Beego验证框架验证Role实体。
// g: pointer to the Role object to be validated, g: 指向待验证的Role对象的指针
// Returns: error if validation fails, nil otherwise. 返回: 如果验证失败则返回错误，否则返回nil。
func checkRole(g *Role) (err error) {
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

// GetRolelist retrieves a paginated list of roles with specified sorting.
// GetRolelist 检索具有指定排序的角色分页列表。
// page: page number (starting from 1), page: 页码（从1开始）
// page_size: number of items per page, page_size: 每页项目数
// sort: sorting field and order (e.g., "-id" for descending by ID), sort: 排序字段和顺序（例如"-id"表示按ID降序）
// Returns: list of roles as orm.Params and total count, 返回: 以orm.Params形式的角色列表和总数
func GetRolelist(page int64, page_size int64, sort string) (roles []orm.Params, count int64) {
	o := orm.NewOrm()
	role := new(Role)
	qs := o.QueryTable(role)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&roles)
	count, _ = qs.Count()
	return roles, count
}

// AddRole adds a new role to the database.
// AddRole 向数据库添加一个新角色。
// r: pointer to the Role object to be added, r: 指向要添加的Role对象的指针
// Returns: ID of the newly inserted role and error if any, 返回: 新插入角色的ID和可能的错误
func AddRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	role := new(Role)
	role.Title = r.Title
	role.Name = r.Name
	role.Remark = r.Remark
	role.Status = r.Status

	id, err := o.Insert(role)
	return id, err
}

// UpdateRole updates an existing role in the database.
// UpdateRole 更新数据库中的现有角色。
// r: pointer to the Role object containing updated information, r: 指向包含更新信息的Role对象的指针
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func UpdateRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	role := make(orm.Params)
	if len(r.Title) > 0 {
		role["Title"] = r.Title
	}
	if len(r.Name) > 0 {
		role["Name"] = r.Name
	}
	if len(r.Remark) > 0 {
		role["Remark"] = r.Remark
	}
	if r.Status != 0 {
		role["Status"] = r.Status
	}
	if len(role) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Role
	num, err := o.QueryTable(table).Filter("Id", r.Id).Update(role)
	return num, err
}

// DelRoleById deletes a role by its ID.
// DelRoleById 根据ID删除角色。
// Id: ID of the role to be deleted, Id: 要删除的角色ID
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func DelRoleById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Role{Id: Id})
	return status, err
}

// GetNodelistByRoleId retrieves a list of nodes associated with a specific role.
// GetNodelistByRoleId 检索与特定角色关联的节点列表。
// Id: ID of the role to retrieve nodes for, Id: 要检索节点的角色ID
// Returns: list of nodes as orm.Params and total count, 返回: 以orm.Params形式的节点列表和总数
func GetNodelistByRoleId(Id int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Role__Role__Id", Id).Values(&nodes)
	return nodes, count
}

// DelGroupNode removes the relationship between a role and all nodes in a group.
// DelGroupNode 删除角色与组中所有节点之间的关系。
// roleid: ID of the role, roleid: 角色ID
// groupid: ID of the group, groupid: 组ID
// Returns: error if any, 返回: 可能的错误
func DelGroupNode(roleid int64, groupid int64) error {
	var nodes []*Node
	var node Node
	role := Role{Id: roleid}
	o := orm.NewOrm()
	num, err := o.QueryTable(node).Filter("Group", groupid).RelatedSel().All(&nodes)
	if err != nil {
		return err
	}
	if num < 1 {
		return nil
	}
	for _, n := range nodes {
		m2m := o.QueryM2M(n, "Role")
		_, err1 := m2m.Remove(&role)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

// AddRoleNode adds a relationship between a role and a node.
// AddRoleNode 添加角色与节点之间的关系。
// roleid: ID of the role, roleid: 角色ID
// nodeid: ID of the node, nodeid: 节点ID
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func AddRoleNode(roleid int64, nodeid int64) (int64, error) {
	o := orm.NewOrm()
	if !o.QueryTable("node_roles").Filter("node_id", nodeid).Filter("role_id", roleid).Exist() {
		role := Role{Id: roleid}
		node := Node{Id: nodeid}
		m2m := o.QueryM2M(&node, "Role")
		num, err := m2m.Add(&role)
		return num, err
	} else {
		return 0, nil
	}
}

// DelUserRole deletes the relationship between a role and all users.
// DelUserRole 删除角色与所有用户之间的关系。
// roleid: ID of the role, roleid: 角色ID
// Returns: error if any, 返回: 可能的错误
func DelUserRole(roleid int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user_roles").Filter("role_id", roleid).Delete()
	return err
}

// DelUserRoleByUserId deletes the relationship between a user and a role.
// DelUserRoleByUserId 删除用户与角色之间的关系。
// roleid: ID of the role (0 means all roles), roleid: 角色ID（0表示所有角色）
// userid: ID of the user, userid: 用户ID
// Returns: error if any, 返回: 可能的错误
func DelUserRoleByUserId(roleid, userid int64) error {
	o := orm.NewOrm()
	var err error
	if roleid == 0 {
		_, err = o.QueryTable("user_roles").Filter("user_id", userid).Delete()
	} else {
		_, err = o.QueryTable("user_roles").Filter("user_id", userid).Filter("role_id", roleid).Delete()
	}
	return err
}

// DelNodeRoleByNodeId deletes the relationship between a node and a role.
// DelNodeRoleByNodeId 删除节点与角色之间的关系。
// roleid: ID of the role, roleid: 角色ID
// nodeid: ID of the node, nodeid: 节点ID
// Returns: error if any, 返回: 可能的错误
func DelNodeRoleByNodeId(roleid, nodeid int64) error {
	o := orm.NewOrm()
	//_, err := o.QueryTable("user_roles").Filter("user_id", userid).Delete()
	_, err := o.QueryTable("node_roles").Filter("node_id", nodeid).Filter("role_id", roleid).Delete()
	return err
}

// AddRoleUser adds a relationship between a role and a user.
// AddRoleUser 添加角色与用户之间的关系。
// roleid: ID of the role, roleid: 角色ID
// userid: ID of the user, userid: 用户ID
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func AddRoleUser(roleid int64, userid int64) (int64, error) {
	o := orm.NewOrm()
	//if !o.QueryTable("user_roles").Filter("user_id", userid).Filter("role_id", roleid).Exist() {
	role := Role{Id: roleid}
	user := User{Id: userid}
	m2m := o.QueryM2M(&user, "Role")
	num, err := m2m.Add(&role)
	return num, err
	//} else {
	//	return 0, nil
	//}
}

// GetUserByRoleId retrieves a list of users associated with a specific role.
// GetUserByRoleId 检索与特定角色关联的用户列表。
// roleid: ID of the role to retrieve users for, roleid: 要检索用户的角色ID
// Returns: list of users as orm.Params and total count, 返回: 以orm.Params形式的用户列表和总数
func GetUserByRoleId(roleid int64) (users []orm.Params, count int64) {
	o := orm.NewOrm()
	user := new(User)
	count, _ = o.QueryTable(user).Filter("Role__Role__Id", roleid).Values(&users)
	return users, count
}

// GetRoleByUserId retrieves a list of role IDs associated with a specific user.
// GetRoleByUserId 检索与特定用户关联的角色ID列表。
// userid: ID of the user to retrieve roles for, userid: 要检索角色的用户ID
// Returns: list of role IDs as orm.Params, 返回: 以orm.Params形式的角色ID列表
func GetRoleByUserId(userid int64) (roleids []orm.Params) {
	o := orm.NewOrm()
	_, err2 := o.QueryTable("user_roles").Filter("user_id", userid).Values(&roleids, "role_id")
	if err2 != nil {
		roleids = []orm.Params{}
	}
	return roleids
}

// GetNodeByRoleId retrieves a list of nodes associated with a specific role.
// GetNodeByRoleId 检索与特定角色关联的节点列表。
// roleid: ID of the role to retrieve nodes for, roleid: 要检索节点的角色ID
// Returns: list of nodes as orm.Params and total count, 返回: 以orm.Params形式的节点列表和总数
func GetNodeByRoleId(roleid int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Role__Role__Id", roleid).Values(&nodes)
	return nodes, count
}

// AccessList retrieves a list of nodes accessible to a user based on their roles.
// AccessList 根据用户的角色检索用户可访问的节点列表。
// uid: ID of the user to retrieve accessible nodes for, uid: 要检索可访问节点的用户ID
// Returns: list of accessible nodes as orm.Params and error if any, 返回: 以orm.Params形式的可访问节点列表和可能的错误
func AccessList(uid int64) (list []orm.Params, err error) {
	var roles []orm.Params
	o := orm.NewOrm()
	role := new(Role)
	_, err = o.QueryTable(role).Filter("User__User__Id", uid).Values(&roles)
	if err != nil {
		return nil, err
	}
	var nodes []orm.Params
	node := new(Node)
	for _, r := range roles {
		_, err := o.QueryTable(node).Filter("Role__Role__Id", r["Id"]).Values(&nodes)
		if err != nil {
			return nil, err
		}
		list = append(list, nodes...)
	}
	return list, nil
}