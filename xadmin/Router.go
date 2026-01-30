package admin

import (
	"xkube/xadmin/src/rbac"

	beego "github.com/beego/beego/v2/server/web"
)

// router configures the routes for the admin backend
// router 配置后台管理的路由
func router() {
	//beego.Router("/", &rbac.MainController{}, "*:Index")
	beego.Router("/index", &rbac.MainController{}, "*:Index")                      // 后台首页 | Admin home page
	beego.Router("/public/login", &rbac.MainController{}, "*:Login")               // 登录接口 | Login endpoint
	beego.Router("/public/appLogin", &rbac.MainController{}, "*:AppLogin")         // App登录 | App login
	beego.Router("/public/logout", &rbac.MainController{}, "*:Logout")             // 登出接口 | Logout endpoint
	beego.Router("/public/changepwd", &rbac.MainController{}, "*:Changepwd")       // 修改密码 | Change password
	beego.Router("/public/myinfo", &rbac.MainController{}, "*:MyInfo")             // 用户信息 | User information
	beego.Router("/public/UpdateMyInfo", &rbac.MainController{}, "*:UpdateMyInfo") // 更新个人信息 | Update my information

	beego.Router("/public/CheckAppLogin", &rbac.MainController{}, "*:CheckAppLogin") // 检查App登录 | Check app login
	beego.Router("/public/check", &rbac.MainController{}, "*:Check")                 // 服务存活检测 | Service health check
	beego.Router("/public/sendCode", &rbac.MainController{}, "*:SendCode")           // 发送手机验证码 | Send SMS verification code

	// User management routes
	// 用户管理路由
	beego.Router("/rbac/user/Add", &rbac.UserController{}, "*:Add")         // 添加用户 | Add user
	beego.Router("/rbac/user/Update", &rbac.UserController{}, "*:Update")   // 更新用户 | Update user
	beego.Router("/rbac/user/Delete", &rbac.UserController{}, "*:Delete")   // 删除用户 | Delete user
	beego.Router("/rbac/user/List", &rbac.UserController{}, "*:List")       // 用户列表 | User list
	beego.Router("/rbac/user/lockAct", &rbac.UserController{}, "*:LockAct") // 处理被锁住的用户或IP | Handle locked users or IPs
	//beego.Router("/public/user/GetuList", &rbac.UserController{}, "*:GetuList")

	// Audit log routes
	// 审计日志路由
	beego.Router("/rbac/audit/List", &rbac.AuditController{}, "*:List") // 审计日志列表 | Audit log list

	//beego.Router("/rbac/node/AddEdit", &rbac.NodeController{}, "*:AddEdit")
	// Node management routes
	// 节点管理路由
	beego.Router("/rbac/node/Add", &rbac.NodeController{}, "*:Add")         // 添加节点 | Add node
	beego.Router("/rbac/node/Edit", &rbac.NodeController{}, "*:Edit")       // 编辑节点 | Edit node
	beego.Router("/rbac/node/Delete", &rbac.NodeController{}, "*:Delete")   // 删除节点 | Delete node
	beego.Router("/rbac/node/List", &rbac.NodeController{}, "*:List")       // 节点列表 | Node list
	beego.Router("/rbac/node/Getlist", &rbac.NodeController{}, "*:Getlist") // 获取节点列表 | Get node list
	beego.Router("/rbac/node/GetPid", &rbac.NodeController{}, "*:GetPid")   // 获取节点PID | Get node PID

	// Cluster authorization routes
	// 集群授权路由
	beego.Router("/rbac/cluster/Delete", &rbac.ClusterController{}, "*:Delete")               // 删除集群授权 | Delete cluster authorization
	beego.Router("/rbac/cluster/List", &rbac.ClusterController{}, "*:List")                   // 集群授权列表 | Cluster authorization list
	beego.Router("/rbac/cluster/MyClusterList", &rbac.ClusterController{}, "*:MyClusterList") // 我的集群列表 | My cluster list
	beego.Router("/rbac/cluster/Add", &rbac.ClusterController{}, "*:Add")                     // 添加集群授权 | Add cluster authorization

	// Group management routes
	// 组管理路由
	beego.Router("/rbac/group/Add", &rbac.GroupController{}, "*:Add")       // 添加组 | Add group
	beego.Router("/rbac/group/Update", &rbac.GroupController{}, "*:Update") // 更新组 | Update group
	beego.Router("/rbac/group/Delete", &rbac.GroupController{}, "*:Delete") // 删除组 | Delete group
	beego.Router("/rbac/group/List", &rbac.GroupController{}, "*:List")     // 组列表 | Group list

	// Role management routes
	// 角色管理路由
	beego.Router("/rbac/role/AddAndEdit", &rbac.RoleController{}, "*:AddAndEdit")         // 添加和编辑角色 | Add and edit role
	beego.Router("/rbac/role/Delete", &rbac.RoleController{}, "*:Delete")                 // 删除角色 | Delete role
	beego.Router("/rbac/role/AddAccess", &rbac.RoleController{}, "*:AddAccess")           // 添加访问权限 | Add access
	beego.Router("/rbac/role/RoleToUserList", &rbac.RoleController{}, "*:RoleToUserList") // 角色到用户列表 | Role to user list
	beego.Router("/rbac/role/RoleToNodeList", &rbac.RoleController{}, "*:RoleToNodeList") // 角色到节点列表 | Role to node list
	beego.Router("/rbac/role/AddRoleToUser", &rbac.RoleController{}, "*:AddRoleToUser")   // 添加角色到用户 | Add role to user
	beego.Router("/rbac/role/DelRoleToUser", &rbac.RoleController{}, "*:DelRoleToUser")   // 删除用户的角色 | Delete user's role
	beego.Router("/rbac/role/DelRoleToNode", &rbac.RoleController{}, "*:DelRoleToNode")   // 删除节点的角色 | Delete node's role
	beego.Router("/rbac/role/Getlist", &rbac.RoleController{}, "*:Getlist")               // 获取角色列表 | Get role list
	beego.Router("/rbac/role/GetridByuid", &rbac.RoleController{}, "*:GetRoleidByUserId") // 根据用户ID获取角色ID | Get role ID by user ID
	beego.Router("/rbac/role/List", &rbac.RoleController{}, "*:List")                     // 角色列表 | Role list
	//beego.Router("/rbac/role/AccessToNode", &rbac.RoleController{}, "*:AccessToNode") // 未用到 | Not used
}
