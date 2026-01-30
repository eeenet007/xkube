package src

import (
	"errors"
	"fmt"
	"log"

	"encoding/gob"

	"net/url"
	"strconv"
	"strings"

	. "xkube/xadmin/src/lib"
	m "xkube/xadmin/src/models"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

var clusterMap = make(map[string]bool)
var accesslist = make(map[string]bool)

func init() {
	gob.Register(clusterMap)
	gob.Register(accesslist)
}

// check access and register user's nodes
func AccessRegister() {
	var Check = func(ctx *context.Context) {
		authType, _ := beego.AppConfig.String("user_auth_type")
		user_auth_type, _ := strconv.Atoi(authType)
		rbac_auth_gateway, _ := beego.AppConfig.String("rbac_auth_gateway")
		var accesslist map[string]bool
		var clusterlist map[string]bool
		if user_auth_type > 0 {
			params := strings.Split(strings.ToLower(strings.Split(ctx.Request.RequestURI, "?")[0]), "/")
			u, err := url.Parse(ctx.Request.RequestURI)
			if err != nil {
				log.Printf("[ERROR] urlParse error:%s\n", err)
				return
			}
			clusterId := u.Query().Get("clusterId")
			//nameSpace := u.Query().Get("nameSpace")
			//if clusterId == "" {
			//	return
			//}
			//log.Printf("clusterId:%s,nameSpace:%s\n", clusterId, nameSpace)

			if CheckAccess(params) {
				uinfo := ctx.Input.Session("userinfo")
				//fmt.Println(uinfo)
				if uinfo == nil {
					//fmt.Println(uinfo)
					ctx.Redirect(302, rbac_auth_gateway)
					return
				}
				//admin用户不用认证权限
				adminuser, _ := beego.AppConfig.String("rbac_admin_user")

				// var rolelist []string
				// for _, vv := range uinfo.(m.User).Role {
				// 	rolelist = append(rolelist, vv.Name)
				// }

				//if uinfo.(m.User).Username == adminuser || strings.Contains(strings.Join(rolelist, ","), "超级管理员") {
				if uinfo.(m.User).Username == adminuser {
					return
				}

				if user_auth_type == 1 {
					listbysession := ctx.Input.Session("accesslist")
					if listbysession != nil {
						accesslist = listbysession.(map[string]bool)
					}
					//集群权限
					clusterbysession := ctx.Input.Session("clusterlist")
					if clusterbysession != nil {
						clusterlist = clusterbysession.(map[string]bool)
					}

				} else if user_auth_type == 2 {
					accesslist, _ = GetAccessList(uinfo.(m.User).Id)
					clusterlist, _ = GetClusterList(uinfo.(m.User).Username)
				}

				//集群权限
				if clusterId != "" {
					//log.Println(clusterId)
					//log.Println(clusterlist)
					vv := AccessClusterDecision(clusterId, clusterlist)
					//log.Println(vv)
					if !vv {
						log.Printf("[WARN] ClusterNoPermission,url:%s, user:%s,clusterId:%s\n", ctx.Request.RequestURI, uinfo.(m.User).Username, clusterId)
						ctx.Output.JSON(&map[string]interface{}{"code": -1, "msg": "集群权限不足"}, true, false)
					} else {
						log.Printf("[INFO] AllowClusterId,url:%s,user:%s,clusterId:%s,value:%v\n", ctx.Request.RequestURI, uinfo.(m.User).Username, clusterId, vv)
					}
				} else {
					log.Printf("[WARN] NoClusterId,url:%s,clusterId:%s\n", ctx.Request.RequestURI, clusterId)
				}

				ret := AccessDecision(params, accesslist)
				if !ret {
					log.Printf("[WARN] UrlNoPermission user:%s,url:%s\n", uinfo.(m.User).Username, params)
					ctx.Output.JSON(&map[string]interface{}{"code": -1, "msg": "模块权限不足"}, true, false)
				}
			}

		}
	}
	beego.InsertFilter("/*", beego.BeforeRouter, Check)
}

// Determine whether need to verify
func CheckAccess(params []string) bool {
	if len(params) < 3 {
		return false
	}
	authType, _ := beego.AppConfig.String("not_auth_package")
	for _, nap := range strings.Split(authType, ",") {
		if params[1] == nap {
			return false
		}
	}
	return true
}

// To test whether permissions
func AccessDecision(params []string, accesslist map[string]bool) bool {
	if CheckAccess(params) {
		//fmt.Println(params)
		paramsLength := len(params)
		s := fmt.Sprintf("%s", params[1])
		if paramsLength == 3 {
			s = fmt.Sprintf("%s/%s", params[1], params[2])
		} else if paramsLength == 4 {
			s = fmt.Sprintf("%s/%s/%s", params[1], params[2], params[3])
		} else if paramsLength == 5 {
			s = fmt.Sprintf("%s/%s/%s/%s", params[1], params[2], params[3], params[4])
		}

		if len(accesslist) < 1 {
			return false
		}
		_, ok := accesslist[s]
		if ok != false {
			return true
		}
	} else {
		return true
	}
	return false
}

type AccessNode struct {
	Id        int64
	Name      string
	Childrens []*AccessNode
}

// Access permissions list
func GetAccessList(uid int64) (map[string]bool, error) {
	list, err := m.AccessList(uid)
	if err != nil {
		return nil, err
	}

	// 创建一个映射以便快速查找节点
	nodeMap := make(map[int64]*AccessNode)
	alist := make([]*AccessNode, 0)

	// 一次遍历处理所有节点
	for _, l := range list {
		// 安全地进行类型断言
		pid, pidOk := l["Pid"].(int64)
		level, levelOk := l["Level"].(int64)
		id, idOk := l["Id"].(int64)
		name, nameOk := l["Name"].(string)

		// 如果类型断言失败，则跳过该节点
		if !pidOk || !levelOk || !idOk || !nameOk {
			continue
		}

		// 查找或创建当前节点
		node, exists := nodeMap[id]
		if !exists {
			node = &AccessNode{
				Id:   id,
				Name: name,
			}
			nodeMap[id] = node
		} else {
			// 如果节点已存在，更新名称
			node.Name = name
		}

		// 根据层级处理节点关系
		if pid == 0 && level == 1 {
			// 添加根节点到列表
			alist = append(alist, node)
		} else if parent, parentExists := nodeMap[pid]; parentExists {
			// 将当前节点添加到父节点的子节点列表中
			parent.Childrens = append(parent.Childrens, node)
		} else {
			// 如果父节点不存在，先创建一个占位符
			parent = &AccessNode{Id: pid}
			nodeMap[pid] = parent
			parent.Childrens = append(parent.Childrens, node)

			// 如果是第二层节点且其父节点是根节点，则加入根节点列表
			if level == 2 {
				alist = append(alist, parent)
			}
		}
	}

	//log.Println(alist)

	// 构建访问权限列表
	accesslist := make(map[string]bool)
	for _, v := range alist {
		for _, v1 := range v.Childrens {
			if len(v1.Childrens) == 0 {
				vname := strings.Split(v.Name, "/")
				if len(vname) > 0 {
					str := fmt.Sprintf("%s/%s", strings.ToLower(vname[0]), strings.ToLower(v1.Name))
					accesslist[str] = true
				}
			} else {
				for _, v2 := range v1.Childrens {
					vname := strings.Split(v.Name, "/")
					v1name := strings.Split(v1.Name, "/")
					if len(vname) > 0 && len(v1name) > 0 {
						str := fmt.Sprintf("%s/%s/%s", strings.ToLower(vname[0]), strings.ToLower(v1name[0]), strings.ToLower(v2.Name))
						accesslist[str] = true
					}
				}
			}
		}
	}
	//log.Println(accesslist)
	return accesslist, nil
}

// access cluster list
func GetClusterList(username string) (map[string]bool, error) {
	//var clusterMap = make(map[string]bool)
	list, count := m.GetUserClusterList(username, "", 0, 0)
	if count <= 0 {
		return clusterMap, fmt.Errorf("GetUserClusterListFail")
	}

	log.Printf("GetClusterList:username:%s,list:%v", username, list)
	accessCluster := make(map[string]bool)
	for _, vv := range list {
		//clusterMap[fmt.Sprintf("%s_check", vv.Userid)] = true
		accessCluster[vv.ClusterId] = true
	}
	//log.Println(clusterMap)
	return accessCluster, nil
}

func AccessClusterDecision(clusterId string, accessCluster map[string]bool) bool {
	//log.Println(accessCluster)
	//if len(accessCluster) < 1 {
	//	return true //当没有给某个用户赋予某个集群时，默认全部都有权限
	//}
	if _, ok1 := accessCluster["all"]; ok1 {
		log.Printf("AccessClusterDecision:vv:all,kk:%v", ok1)
		return true //当给用户授予all时，全部集群都有权限
	}
	vv, ok := accessCluster[clusterId]
	log.Printf("AccessClusterDecision:vv:%v,kk:%v", vv, ok)
	if ok != false {
		return true //能查询到用户有设置某个集群，是返回有权限
	}
	return false //不能查询到则返回无权限
}

// check login
func CheckLogin(username string, password string) (user m.User, err error) {
	user, err = m.GetUserByUsername(username)
	if err != nil {
		return user, err
	}
	if user.Id == 0 {
		return user, errors.New("用户不存在")
	}
	if user.Password != Pwdhash(password) {
		return user, errors.New("密码错误")
	}
	return user, nil
}
