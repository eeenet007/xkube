// clusterModel.go
package models

import (
	//"bytes"
	"errors"
	"fmt"

	//"io/ioutil"
	"log"
	"reflect"
	"time"

	"xkube/common"

	//"github.com/beego/beego/v2/core/validation"

	"github.com/beego/beego/v2/client/orm"
)

type Xkb_cicd_config struct { //cicd的全部配置
	// Id CICD configuration primary key
	// Id CICD配置主键
	Id int64 `json:"id"`
	// CicdName CICD name
	// CicdName CICD名称
	CicdName string `json:"cicd_name"`
	// CicdType CICD type
	// CicdType CICD类型
	CicdType int64 `json:"cicd_type"`
	// Appname application name
	// Appname 应用名称
	Appname string `json:"appname"`
	// Status CICD status
	// Status CICD状态
	Status int64 `json:"status"`
	// Remarks remarks or description
	// Remarks 备注或描述
	Remarks string `json:"remarks"`
	// Createtime creation time
	// Createtime 创建时间
	Createtime string `json:"createtime"`
}

type Xkb_cicd struct { //查询时会将大写以下划线分开
	// Id CICD primary key
	// Id CICD主键
	Id int64 `json:"id"`
	// CicdName CICD name
	// CicdName CICD名称
	CicdName string `json:"cicd_name"`
	// CicdType CICD type
	// CicdType CICD类型
	CicdType int64 `json:"cicd_type"`
	// Appname application name
	// Appname 应用名称
	Appname string `json:"appname"`
	// ClusterId cluster identifier
	// ClusterId 集群标识符
	ClusterId string `json:"cluster_id"`
	// Namespace Kubernetes namespace
	// Namespace Kubernetes命名空间
	Namespace string `json:"namespace"`
	// Status CICD status
	// Status CICD状态
	Status int64 `json:"status"`
	// LastRuntime last run time
	// LastRuntime 最后运行时间
	LastRuntime string `json:"last_runtime"`
	// Remarks remarks or description
	// Remarks 备注或描述
	Remarks string `json:"remarks"`
	// Createtime creation time
	// Createtime 创建时间
	Createtime string `json:"createtime"`
}

type Xkb_cicd_pipelines struct {
	// Id pipeline primary key
	// Id 流水线主键
	Id int64 `json:"id"`
	// CicdId CICD identifier
	// CicdId CICD标识符
	CicdId int64 `json:"cicd_id"`
	// AliyunId Alibaba Cloud identifier
	// AliyunId 阿里云标识符
	AliyunId string `json:"aliyun_id"`
	// OrganizationId organization identifier
	// OrganizationId 组织标识符
	OrganizationId string `json:"organization_id"`
	// jenkins name,帐号
	JksId string `json:"jks_id"`
	// PipelineId pipeline identifier
	// PipelineId 流水线标识符
	PipelineId string `json:"pipeline_id"`
}

// 阿里云 ak登记表
type Xkb_cicd_ak struct {
	// Id access key primary key
	// Id 访问密钥主键
	Id int64 `json:"id"`
	// AliyunId Alibaba Cloud identifier
	// AliyunId 阿里云标识符
	AliyunId string `json:"aliyun_id"`
	// AccesskeyId access key ID
	// AccesskeyId 访问密钥ID
	AccesskeyId string `json:"accesskey_id"`
	// AccesskeySecret access key secret
	// AccesskeySecret 访问密钥密钥
	AccesskeySecret string `json:"accesskey_secret"`
	// OrganizationId organization identifier
	// OrganizationId 组织标识符
	OrganizationId string `json:"organization_id"`
	// Remarks remarks or description
	// Remarks 备注或描述
	Remarks string `json:"remarks"`
	// Createtime creation time
	// Createtime 创建时间
	Createtime string `json:"createtime"`
}

// jenkins url、用户名、密码记录表
type Xkb_cicd_jks struct {
	Id         int64  `json:"id"`
	JksId      string `json:"jks_id"`
	JksUrl     string `json:"jks_url"`
	JksUser    string `json:"jks_user"`
	JksPasswd  string `json:"jks_passwd"`
	Remarks    string `json:"remarks"`
	Createtime string `json:"createtime"`
}

/*
*

	type Xkb_cicd_build struct {
		// Id build configuration primary key
		// Id 构建配置主键
		Id int64 `json:"id"`
		// CicdId CICD identifier
		// CicdId CICD标识符
		CicdId int64 `json:"cicd_id"`
		// CodeType code repository type
		// CodeType 代码仓库类型
		CodeType string `json:"code_type"`
		// CodeVersion code version
		// CodeVersion 代码版本
		CodeVersion string `json:"code_version"`
		// BuildCmd build command
		// BuildCmd 构建命令
		BuildCmd string `json:"build_cmd"`
		// OtherConfig other configuration
		// OtherConfig 其他配置
		OtherConfig string `json:"other_config"`
	}

	type Xkb_cicd_git struct {
		// Id git configuration primary key
		// Id git配置主键
		Id int64 `json:"id"`
		// CicdId CICD identifier
		// CicdId CICD标识符
		CicdId int64 `json:"cicd_id"`
		// GitUrl git repository URL
		// GitUrl git仓库URL
		GitUrl string `json:"git_url"`
		// GitBranch git branch
		// GitBranch git分支
		GitBranch string `json:"git_branch"`
		// TokenId token identifier
		// TokenId 令牌标识符
		TokenId string `json:"token_id"`
		// Remarks remarks or description
		// Remarks 备注或描述
		Remarks string `json:"remarks"`
	}

	type Xkb_cicd_k8s struct {
		// Id Kubernetes configuration primary key
		// Id Kubernetes配置主键
		Id int64 `json:"id"`
		// CicdId CICD identifier
		// CicdId CICD标识符
		CicdId int64 `json:"cicd_id"`
		// ClusterId cluster identifier
		// ClusterId 集群标识符
		ClusterId string `json:"cluster_id"`
		// Namespace Kubernetes namespace
		// Namespace Kubernetes命名空间
		Namespace string `json:"namespace"`
		// ResType resource type
		// ResType 资源类型
		ResType string `json:"res_type"`
		// ResName resource name
		// ResName 资源名称
		ResName string `json:"res_name"`
		// ContainerName container name
		// ContainerName 容器名称
		ContainerName string `json:"container_name"`
		// Remarks remarks or description
		// Remarks 备注或描述
		Remarks string `json:"remarks"`
	}

	type Xkb_cicd_token struct {
		// Id token primary key
		// Id 令牌主键
		Id int64 `json:"id"`
		// TokenName token name
		// TokenName 令牌名称
		TokenName int64 `json:"token_name"`
		// TokenType token type
		// TokenType 令牌类型
		TokenType string `json:"token_type"`
		// TokenUrl token URL
		// TokenUrl 令牌URL
		TokenUrl string `json:"token_url"`
		// TokenUsername token username
		// TokenUsername 令牌用户名
		TokenUsername string `json:"token_username"`
		// TokenPassword token password
		// TokenPassword 令牌密码
		TokenPassword string `json:"token_password"`
		// Remarks remarks or description
		// Remarks 备注或描述
		Remarks string `json:"remarks"`
		// Createtime creation time
		// Createtime 创建时间
		Createtime string `json:"createtime"`
	}

*
*/
func init() {
	//orm.Debug = true
	orm.RegisterModel(new(Xkb_cicd), new(Xkb_cicd_pipelines), new(Xkb_cicd_ak), new(Xkb_cicd_jks))
}

// GetList_Cicd Get CICD list with pagination and filtering capabilities
// GetList_Cicd 获取带分页和过滤功能的CICD列表
// id: CICD ID for filtering, empty means no filtering
// id: 用于过滤的CICD ID，空表示不过滤
// cicdName: CICD name for filtering, empty means no filtering
// cicdName: 用于过滤的CICD名称，空表示不过滤
// appname: application name for filtering, empty means no filtering
// appname: 用于过滤的应用名称，空表示不过滤
// page: page number (starting from 1)
// page: 页码（从1开始）
// page_size: number of items per page
// page_size: 每页条目数
// Returns:
//   - []Xkb_cicd: CICD list
//   - int64: total count of CICD items matching the criteria
//
// 返回值:
//   - []Xkb_cicd: CICD列表
//   - int64: 符合条件的CICD项目总数
func GetList_Cicd(id, cicdName, appname string, page, page_size int64) (dps []Xkb_cicd, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_cicd))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if id != "" {
		cond = cond.And("id", id)
	}

	if cicdName != "" {
		cond = cond.And("cicd_name__contains", cicdName)
	}

	if appname != "" {
		cond = cond.And("appname", appname)
	}
	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&dps)
	//qs.All(&xfr)
	count, _ = qs.Count()
	return dps, count
}

// CicdListAppname Get list of unique application names from CICD configurations
// CicdListAppname 从CICD配置中获取唯一的应用名称列表
// Returns:
//   - []Xkb_cicd: list of CICD items with unique application names
//   - int64: total count of unique application names
//
// 返回值:
//   - []Xkb_cicd: 包含唯一应用名称的CICD项目列表
//   - int64: 唯一应用名称的总数
func CicdListAppname() (dps []Xkb_cicd, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_cicd))
	qs.GroupBy("appname").All(&dps, "appname")
	count, _ = qs.Count()
	return dps, count
}

// UpdateStatus_Cicd Update CICD status and last run time
// UpdateStatus_Cicd 更新CICD状态和最后运行时间
// id: CICD ID
// id: CICD ID
// status: new status value
// status: 新状态值
// lastRunTime: last run time, empty means no update
// lastRunTime: 最后运行时间，空表示不更新
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func UpdateStatus_Cicd(id, status, lastRunTime string) (int64, error) {
	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("UPDATE xkb_cicd SET status = '%s' WHERE id = '%s'", status, id)
	if lastRunTime != "" && lastRunTime != "undefined" {
		sqlstr = fmt.Sprintf("UPDATE xkb_cicd SET status = '%s',last_runtime = '%s' WHERE id = '%s'", status, lastRunTime, id)
	}
	if status == "2" { //RUNNING装表示第一次启动然后更新一个启动时间
		dd := time.Now().Format("2006-01-02 15:04:05")
		sqlstr = fmt.Sprintf("UPDATE xkb_cicd SET status = '%s',last_runtime = '%s' WHERE id = '%s'", status, dd, id)
	}
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("[INFO] UpdateStatus_Cicd affected nums:%d\n", num)
		return num, nil
	}
	return 0, err
}

// Add_Cicd_Pipelines Add CICD and pipeline configurations in a transaction
// Add_Cicd_Pipelines 通过事务添加CICD和流水线配置
// u: CICD configuration to be added
// u: 要添加的CICD配置
// b: Pipeline configuration to be added
// b: 要添加的流水线配置
// Returns:
//   - int64: inserted pipeline record ID
//   - error: error message
//
// 返回值:
//   - int64: 插入的流水线记录ID
//   - error: 错误信息
func Add_Cicd_Pipelines(u *Xkb_cicd, b *Xkb_cicd_pipelines) (int64, error) {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return 0, err
	}
	Ds := new(Xkb_cicd)
	Bs := new(Xkb_cicd_pipelines)
	Ds = u
	Bs = b
	Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
	Ds.Status = 1
	id, err := tx.Insert(Ds)
	if err != nil {
		log.Printf("[ERROR] Add_Cicd_Pipelines Fail:%s\n", err)
		tx.Rollback()
		return 0, err
	}
	Bs.CicdId = id
	bid, err2 := tx.Insert(Bs)
	if err2 != nil {
		log.Printf("[ERROR] Add_Cicd_Pipelines Fail2:%s\n", err2)
		tx.Rollback()
		return 0, err2
	}
	tx.Commit()
	return bid, nil
}

// GetCicdByCicdName Get CICD configuration by CICD name, cluster ID and namespace
// GetCicdByCicdName 根据CICD名称、集群ID和命名空间获取CICD配置
// cicdName: CICD name for filtering, empty means no filtering
// cicdName: 用于过滤的CICD名称，空表示不过滤
// clusterId: cluster ID for filtering, empty means no filtering
// clusterId: 用于过滤的集群ID，空表示不过滤
// nameSpace: namespace for filtering, empty means no filtering
// nameSpace: 用于过滤的命名空间，空表示不过滤
// Returns:
//   - Xkb_cicd: CICD configuration
//   - error: error message
//
// 返回值:
//   - Xkb_cicd: CICD配置
//   - error: 错误信息
func GetCicdByCicdName(cicdName, clusterId, nameSpace string) (Xkb_cicd, error) {
	o := orm.NewOrm()
	var stu Xkb_cicd
	qs := o.QueryTable(new(Xkb_cicd))
	cond := orm.NewCondition()
	if cicdName != "" {
		cond = cond.And("cicd_name", cicdName)
	}
	if clusterId != "" {
		cond = cond.And("cluster_id", clusterId)
	}
	if nameSpace != "" {
		cond = cond.And("namespace", nameSpace)
	}
	qs = qs.SetCond(cond)
	err := qs.One(&stu)
	if err != nil {
		log.Printf("[ERROR] GetCicdByCicdName query Fail:%s\n", err)
	}
	return stu, err
}

// GetPipelinesByCicdId Get pipeline configuration by CICD ID
// GetPipelinesByCicdId 根据CICD ID获取流水线配置
// cicdId: CICD ID
// cicdId: CICD ID
// Returns:
//   - Xkb_cicd_pipelines: pipeline configuration
//   - error: error message
//
// 返回值:
//   - Xkb_cicd_pipelines: 流水线配置
//   - error: 错误信息
func GetPipelinesByCicdId(cicdId int64) (Xkb_cicd_pipelines, error) {
	o := orm.NewOrm()
	var stu Xkb_cicd_pipelines
	qs := o.QueryTable(new(Xkb_cicd_pipelines))
	err := qs.Filter("cicd_id", cicdId).One(&stu)
	if err != nil {
		log.Printf("[ERROR] GetPipelinesByCicdId query Fail:%s\n", err)
	}
	return stu, err
}

// GetAliyunIdList Get Alibaba Cloud ID list with filtering
// GetAliyunIdList 获取带过滤功能的阿里云ID列表
// aliyun_id: Alibaba Cloud ID for filtering, empty means no filtering
// aliyun_id: 用于过滤的阿里云ID，空表示不过滤
// Returns:
//   - []Xkb_cicd_ak: list of Alibaba Cloud access key configurations
//   - int64: total count of items
//
// 返回值:
//   - []Xkb_cicd_ak: 阿里云访问密钥配置列表
//   - int64: 项目总数
func GetAliyunIdList(aliyun_id string) (dps []Xkb_cicd_ak, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_cicd_ak))
	cond := orm.NewCondition()
	if aliyun_id != "" {
		cond = cond.And("aliyun_id", aliyun_id)
	}
	qs = qs.SetCond(cond)
	qs.OrderBy("-id").All(&dps, "id", "aliyun_id", "organization_id")
	count, _ = qs.Count()
	return dps, count
}

// DelById_Cicd Delete CICD configuration by ID along with associated pipeline configurations
// DelById_Cicd 根据ID删除CICD配置及关联的流水线配置
// id: CICD ID to be deleted
// id: 要删除的CICD ID
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func DelById_Cicd(id int64) (int64, error) {
	o := orm.NewOrm()
	tx, err := o.Begin()
	if err != nil {
		return 0, err
	}
	_, err = tx.Delete(&Xkb_cicd{Id: id})
	if err != nil {
		log.Printf("[ERROR] DelById_Cicd Fail:%s\n", err)
		tx.Rollback()
		return 0, err
	}

	res, err := tx.Raw(`DELETE FROM xkb_cicd_pipelines WHERE cicd_id = ?`, id).Exec()
	if err != nil {
		log.Printf("[ERROR] Del DelById_Cicd err2:%s\n", err)
		tx.Rollback()
		return 0, err
	}
	num, err2 := res.RowsAffected()
	if err2 != nil {
		log.Printf("[ERROR] Del customurlFail3 err:%s\n", err2)
		tx.Rollback()
		return 0, err2
	}
	tx.Commit()
	return num, nil
}

// Edit_data Edit field value by ID in specified table
// Edit_data 在指定表中根据ID编辑字段值
// id: record ID
// id: 记录ID
// key: field name to be updated
// key: 要更新的字段名
// value: new value for the field
// value: 字段的新值
// table: table name
// table: 表名
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func Edit_data(id, key, value, table string) (int64, error) {
	o := orm.NewOrm()
	sqlstr := fmt.Sprintf("UPDATE %s SET %s = '%s' WHERE id = %s", table, key, value, id)
	res, err := o.Raw(sqlstr).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		log.Printf("[INFO] Edit affected nums:%d\n", num)
		return num, nil
	}
	return 0, err
}

// Update_Cicd Update CICD and pipeline configurations in a transaction
// Update_Cicd 通过事务更新CICD和流水线配置
// u: CICD configuration to be updated
// u: 要更新的CICD配置
// b: Pipeline configuration to be updated
// b: 要更新的流水线配置
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func Update_Cicd(u *Xkb_cicd, b *Xkb_cicd_pipelines) (int64, error) {
	o := orm.NewOrm()
	tx, _ := o.Begin()
	//更新表一
	op := make(orm.Params)
	t := reflect.TypeOf(*u)
	v := reflect.ValueOf(*u)
	for i := 0; i < t.NumField(); i++ {
		//op[string(t.Field(i).Name)] = v.Field(i).Interface()
		vt := v.Field(i).Interface()
		var isChg bool
		switch vt := vt.(type) { //将interface进行断言判断，判断值有更改再加入到map。
		case int:
			if vt > 0 {
				isChg = true
			}
		case string:
			if vt != "" {
				isChg = true
			}
		default:
			isChg = false
		}
		if isChg {
			op[string(t.Field(i).Name)] = v.Field(i).Interface()
		}
	}
	if len(op) == 0 {
		return 0, errors.New("Xkb_cicd update field is empty")
	}

	//更新表二
	op2 := make(orm.Params)
	tt := reflect.TypeOf(*b)
	vv := reflect.ValueOf(*b)
	for i := 0; i < tt.NumField(); i++ {
		vvt := vv.Field(i).Interface()
		var isChg bool
		switch vvt := vvt.(type) { //将interface进行断言判断，判断值有更改再加入到map。
		case int:
			if vvt > 0 {
				isChg = true
			}
		case string:
			if vvt != "" {
				isChg = true
			}
		default:
			isChg = false
		}
		if isChg {
			op2[string(tt.Field(i).Name)] = vv.Field(i).Interface()
		}
	}
	if len(op2) == 0 {
		return 0, errors.New("Xkb_cicd_pipelines update field is empty")
	}

	var table Xkb_cicd
	var table2 Xkb_cicd_pipelines

	num1, err2 := tx.QueryTable(table).Filter("Id", u.Id).Update(op)
	if err2 != nil {
		log.Printf("[ERROR] Update_Cicd Fail:%s\n", err2)
		tx.Rollback()
		return 0, err2
	}
	//log.Println(u.Id)

	num2, err3 := tx.QueryTable(table2).Filter("CicdId", u.Id).Update(op2)
	if err3 != nil {
		log.Printf("[ERROR] Update_Cicd_Pipelines Fail:%s\n", err3)
		tx.Rollback()
		return 0, err3
	}
	tx.Commit()
	return num1 + num2, nil
}

// GetAkList Get access key list with pagination and filtering
// GetAkList 获取带分页和过滤功能的访问密钥列表
// id: access key ID for filtering, empty means no filtering
// id: 用于过滤的访问密钥ID，空表示不过滤
// page: page number (starting from 1)
// page: 页码（从1开始）
// page_size: number of items per page
// page_size: 每页条目数
// Returns:
//   - []Xkb_cicd_ak: list of access key configurations
//   - int64: total count of items
//
// 返回值:
//   - []Xkb_cicd_ak: 访问密钥配置列表
//   - int64: 项目总数
func GetAkList(id string, page, page_size int64) (dps []Xkb_cicd_ak, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_cicd_ak))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if id != "" {
		cond = cond.And("id", id)
	}
	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&dps)
	count, _ = qs.Count()
	return dps, count
}

// Add_Ak Add or update access key configuration
// Add_Ak 添加或更新访问密钥配置
// u: access key configuration to be added or updated
// u: 要添加或更新的访问密钥配置
// Returns:
//   - int64: inserted or updated record ID
//   - error: error message
//
// 返回值:
//   - int64: 插入或更新的记录ID
//   - error: 错误信息
func Add_Ak(u *Xkb_cicd_ak) (int64, error) {
	o := orm.NewOrm()
	Ds := new(Xkb_cicd_ak)
	Ds = u
	if u.Id > 0 { //有ID就更新，无ID就添加
		op := make(orm.Params)
		t := reflect.TypeOf(*u)
		v := reflect.ValueOf(*u)
		for i := 0; i < t.NumField(); i++ {
			op[string(t.Field(i).Name)] = v.Field(i).Interface()
		}
		if len(op) == 0 {
			return 0, errors.New("Xkb_cicd update field is empty")
		}
		var table Xkb_cicd_ak
		num, err := o.QueryTable(table).Filter("Id", u.Id).Update(op)
		if err != nil {
			log.Printf("[ERROR] UpdateAk Fail:%s\n", err)
			return 0, err
		}
		return num, nil
	} else {
		Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
		id, err := o.Insert(Ds)
		if err != nil {
			log.Printf("[ERROR] AddAk Fail:%s\n", err)
			return 0, err
		}
		return id, nil
	}
}

// DelById_Ak Delete access key configuration by ID
// DelById_Ak 根据ID删除访问密钥配置
// id: access key ID to be deleted
// id: 要删除的访问密钥ID
// aliyunId: Alibaba Cloud ID to be removed from cache
// aliyunId: 要从缓存中移除的阿里云ID
// Returns:
//   - int64: number of affected rows
//   - error: error message
//
// 返回值:
//   - int64: 受影响的行数
//   - error: 错误信息
func DelById_Ak(id int64, aliyunId string) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&Xkb_cicd_ak{Id: id})
	if err != nil {
		log.Printf("[ERROR] DelById_Ak Fail:%s\n", err)
		return 0, err
	}
	delete(common.AliyunIdMap, aliyunId)
	return num, nil
}

func GetJksList(id string, page, page_size int64) (dps []Xkb_cicd_jks, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_cicd_jks))
	cond := orm.NewCondition()
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if id != "" {
		cond = cond.And("id", id)
	}
	qs = qs.SetCond(cond)
	qs.Limit(page_size, offset).OrderBy("-id").All(&dps)
	count, _ = qs.Count()
	return dps, count
}

func GetJksListNopasswd() (dps []Xkb_cicd_jks, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Xkb_cicd_jks))
	cond := orm.NewCondition()
	qs = qs.SetCond(cond)
	qs.OrderBy("-id").All(&dps, "id", "jks_id", "jks_url")
	count, _ = qs.Count()
	return dps, count
}

func Add_Jks(u *Xkb_cicd_jks) (int64, error) {
	o := orm.NewOrm()
	Ds := new(Xkb_cicd_jks)
	Ds = u
	if u.Id > 0 { //有ID就更新，无ID就添加
		op := make(orm.Params)
		t := reflect.TypeOf(*u)
		v := reflect.ValueOf(*u)
		for i := 0; i < t.NumField(); i++ {
			op[string(t.Field(i).Name)] = v.Field(i).Interface()
		}
		if len(op) == 0 {
			return 0, errors.New("Xkb_cicd_jks update field is empty")
		}
		var table Xkb_cicd_jks
		num, err := o.QueryTable(table).Filter("Id", u.Id).Update(op)
		if err != nil {
			log.Printf("[ERROR] UpdateJks Fail:%s\n", err)
			return 0, err
		}
		return num, nil
	} else {
		Ds.Createtime = time.Now().Format("2006-01-02 15:04:05")
		id, err := o.Insert(Ds)
		if err != nil {
			log.Printf("[ERROR] AddJks Fail:%s\n", err)
			return 0, err
		}
		return id, nil
	}
}

func DelById_Jks(id int64, jksId string) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Delete(&Xkb_cicd_jks{Id: id})
	if err != nil {
		log.Printf("[ERROR] DelById_Jks Fail:%s\n", err)
		return 0, err
	}
	delete(JksIdMap, jksId)
	return num, nil
}
