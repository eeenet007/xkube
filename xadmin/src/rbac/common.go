package rbac

import (
	//"fmt"
	"sort"

	. "xkube/xadmin/src"

	beego "github.com/beego/beego/v2/server/web"

	l "xkube/xadmin/src/lib"
	m "xkube/xadmin/src/models"
)

type CommonController struct {
	beego.Controller
	Templatetype string //ui template type
}

func (this *CommonController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJSON()
}

func (this *CommonController) GetTemplatetype() string {
	templatetype, _ := beego.AppConfig.String("template_type")
	if templatetype == "" {
		templatetype = "easyui"
	}
	return templatetype
}

// 排序数据结构
type ByCount []Tree

func (c ByCount) Len() int           { return len(c) }
func (c ByCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool { return c[i].Sorts > c[j].Sorts }

func (this *CommonController) GetTree(Uid int64) []Tree {
	//func (this *CommonController) GetTree() []Tree {
	nodes, _ := m.GetNodeTree(0, 1)
	tree := make([]Tree, len(nodes))
	//tree := make([]Tree, 0, 0)
	var nodeidarry []int64
	if Uid != 111111 { //超级管理员显示所有菜单
		accesslist, _ := m.AccessList(Uid)
		//fmt.Println(accesslist)
		for _, vv := range accesslist {
			nodeidarry = append(nodeidarry, vv["Id"].(int64))
		}
	}
	//fmt.Println(nodes)
	for k, v := range nodes {
		if l.IsExistInArry(nodeidarry, v["Id"].(int64)) || (Uid == 111111) {
			tree[k].Id = v["Id"].(int64)
			tree[k].Text = v["Title"].(string)
			if _, ok := v["Icons"]; ok {
				tree[k].IconCls = v["Icons"].(string) //图标
			}
			if _, ok := v["Sorts"]; ok {
				tree[k].Sorts = v["Sorts"].(int64) //排序
			}
			children, _ := m.GetNodeTree(v["Id"].(int64), 2)
			tree[k].Children = make([]Tree, len(children))
			for k1, v1 := range children {
				if l.IsExistInArry(nodeidarry, v1["Id"].(int64)) || (Uid == 111111) {
					tree[k].Children[k1].Id = v1["Id"].(int64)
					tree[k].Children[k1].Text = v1["Title"].(string)
					if _, ok := v1["Sorts"]; ok {
						tree[k].Children[k1].Sorts = v1["Sorts"].(int64) //排序
					}
					tree[k].Children[k1].Attributes.Url = "/" + v["Name"].(string) + "/" + v1["Name"].(string)
				}
			}
			sort.Sort(sort.Reverse(ByCount(tree[k].Children)))
		}
	}
	//sort.Sort(ByCount(tree)) //降序
	sort.Sort(sort.Reverse(ByCount(tree))) //升序
	return tree
}

func init() {

	//验证权限
	AccessRegister()
}
