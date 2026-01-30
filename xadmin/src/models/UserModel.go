// Package models provides data models and related functionality for the xkube admin system.
// 包models为xkube管理系统提供数据模型及相关功能。
package models

import (
	"errors"
	"log"
	"time"

	lib "xkube/xadmin/src/lib"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	//beego "github.com/beego/beego/v2/server/web"
)

// User represents a user in the system with authentication and authorization information.
// User 表示系统中具有身份验证和授权信息的用户。
type User struct {
	Id            int64     // Id is the unique identifier of the user. Id是用户的唯一标识符。
	Sessionid     string    `orm:"size(32)"`                                                                 // Sessionid is the session identifier for the user. Sessionid是用户的会话标识符。
	Username      string    `orm:"unique;size(32)" form:"Username"  valid:"Required;MaxSize(20);MinSize(6)"` // Username is the login name of the user. Username是用户的登录名。
	Password      string    `orm:"size(32)" form:"Password" valid:"Required;MaxSize(20);MinSize(6)"`         // Password is the encrypted password of the user. Password是用户的加密密码。
	Repassword    string    `orm:"-" form:"Repassword" valid:"Required"`                                     // Repassword is the password confirmation field (not stored in database). Repassword是密码确认字段（不存储在数据库中）。
	Nickname      string    `orm:"unique;size(32)" form:"Nickname" valid:"Required;MaxSize(20);MinSize(2)"`  // Nickname is the display name of the user. Nickname是用户的显示名称。
	Telphone      string    `orm:"unique;size(11)" form:"Telphone" valid:"Required;MaxSize(11);MinSize(11)"` // Telphone is the phone number of the user. Telphone是用户的电话号码。
	Email         string    `orm:"size(32)" form:"Email" valid:"Email"`                                      // Email is the email address of the user. Email是用户的电子邮件地址。
	Company       string    `orm:"size(128)" form:"Company"`                                                 // Company is the company name of the user. Company是用户的公司名称。
	Department    string    `orm:"size(128)" form:"Department"`                                              // Department is the department of the user. Department是用户的部门。
	Remark        string    `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`                        // Remark is additional information about the user. Remark是关于用户的附加信息。
	Status        int       `orm:"default(1)" form:"Status" valid:"Range(0,1)"`                              // Status indicates the status of the user (0 for inactive, 1 for active). Status表示用户的状态（0为未激活，1为激活）。
	Lastlogintime string    `orm:"null;size(255)" form:"-"`                                                  // Lastlogintime is the last login time of the user. Lastlogintime是用户的最后登录时间。
	Lastloginip   string    `orm:"null;size(255)" form:"-"`                                                  // Lastloginip is the last login IP address of the user. Lastloginip是用户的最后登录IP地址。
	Createtime    time.Time `orm:"type(datetime);auto_now_add"`                                              // Createtime is the creation time of the user account. Createtime是用户账户的创建时间。
	Role          []*Role   `orm:"rel(m2m)"`                                                                 // Role is the list of roles assigned to the user. Role是分配给用户的角色列表。
}

// func (u *User) TableName() string {
// 	vTable, _ := beego.AppConfig.String("rbac_user_table")
// 	return vTable
// }

// Valid performs custom validation on the User entity.
// Valid 对User实体执行自定义验证。
// v: validation object to set errors on, v: 用于设置错误的验证对象
func (u *User) Valid(v *validation.Validation) {
	if u.Password != u.Repassword {
		v.SetError("Repassword", "两次输入的密码不一样")
		// v.SetError("Repassword", "Passwords do not match")
	}
}

// checkUser validates the User entity using the Beego validation framework.
// checkUser 使用Beego验证框架验证User实体。
// u: pointer to the User object to be validated, u: 指向待验证的User对象的指针
// Returns: error if validation fails, nil otherwise. 返回: 如果验证失败则返回错误，否则返回nil。
func checkUser(u *User) (err error) {
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

// init registers the User model with the ORM framework.
// init 将User模型注册到ORM框架中。
func init() {
	//orm.Debug = true
	orm.RegisterModel(new(User))
}

/************************************************************/

// Getuserlist retrieves a paginated list of users with specified sorting.
// Getuserlist 检索具有指定排序的用户分页列表。
// page: page number (starting from 1), page: 页码（从1开始）
// page_size: number of items per page, page_size: 每页项目数
// sort: sorting field and order (e.g., "-id" for descending by ID), sort: 排序字段和顺序（例如"-id"表示按ID降序）
// Returns: list of users as orm.Params and total count, 返回: 以orm.Params形式的用户列表和总数
func Getuserlist(page int64, page_size int64, sort string) (users []orm.Params, count int64) {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable(user)
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	qs.Limit(page_size, offset).OrderBy(sort).Values(&users, "id", "sessionid", "username", "nickname", "telphone", "email", "department", "company", "lastloginip", "lastlogintime", "remark", "status", "createtime")
	count, _ = qs.Count()
	return users, count
}

// AddUser adds a new user to the database.
// AddUser 向数据库添加一个新用户。
// u: pointer to the User object to be added, u: 指向要添加的User对象的指针
// Returns: ID of the newly inserted user and error if any, 返回: 新插入用户的ID和可能的错误
func AddUser(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	user := new(User)
	user.Sessionid = lib.GetRandomString(32)
	user.Username = u.Username
	user.Password = lib.Strtomd5(u.Password)
	user.Nickname = u.Nickname
	user.Telphone = u.Telphone
	user.Email = u.Email
	user.Company = u.Company
	user.Department = u.Department
	user.Remark = u.Remark
	user.Status = u.Status

	id, err := o.Insert(user)
	return id, err
}

// UpdateUser updates an existing user in the database.
// UpdateUser 更新数据库中的现有用户。
// u: pointer to the User object containing updated information, u: 指向包含更新信息的User对象的指针
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func UpdateUser(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	user := make(orm.Params)
	if len(u.Username) > 0 {
		user["Username"] = u.Username
	}
	if len(u.Nickname) > 0 {
		user["Nickname"] = u.Nickname
	}
	if len(u.Email) > 0 {
		user["Email"] = u.Email
	}
	if len(u.Telphone) > 0 {
		user["Telphone"] = u.Telphone
	}
	if len(u.Company) > 0 {
		user["Company"] = u.Company
	}
	if len(u.Department) > 0 {
		user["Department"] = u.Department
	}
	if len(u.Remark) > 0 {
		user["Remark"] = u.Remark
	}
	if len(u.Password) > 0 {
		user["Password"] = lib.Strtomd5(u.Password)
	}
	//if u.Status != 0 {
	//	user["Status"] = u.Status
	//}
	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table User
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(user)
	return num, err
}

// UpdateLoginIpTime updates the last login time and IP address of a user.
// UpdateLoginIpTime 更新用户的最后登录时间和IP地址。
// id: ID of the user to update, id: 要更新的用户ID
// logintime: last login time, logintime: 最后登录时间
// loginip: last login IP address, loginip: 最后登录IP地址
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func UpdateLoginIpTime(id int64, logintime, loginip string) (int64, error) {
	o := orm.NewOrm()
	user := make(orm.Params)
	user["Lastlogintime"] = logintime
	user["Lastloginip"] = loginip
	var table User
	num, err := o.QueryTable(table).Filter("Id", id).Update(user)
	return num, err
}

// DelUserById deletes a user by its ID.
// DelUserById 根据ID删除用户。
// Id: ID of the user to be deleted, Id: 要删除的用户ID
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func DelUserById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&User{Id: Id})
	return status, err
}

// GetUserByUsername retrieves a user by username.
// GetUserByUsername 根据用户名检索用户。
// username: username of the user to retrieve, username: 要检索的用户名
// Returns: User object and error if any, 返回: User对象和可能的错误
func GetUserByUsername(username string) (user User, err error) {
	//user = User{Username: username}
	o := orm.NewOrm()
	usert := new(User)
	err = o.QueryTable(usert).Filter("Username", username).Filter("Status", 1).One(&user)
	if err != nil {
		return user, err
	}
	//err = o.Raw("select * from user where status = '2' AND username = ?", username).QueryRow(&user)
	//if err != nil {
	//	return user, err
	//}
	//o.Read(&user, "Username")

	return user, nil
}

// GetUserByPhoneNo retrieves a user by phone number.
// GetUserByPhoneNo 根据电话号码检索用户。
// phoneNo: phone number of the user to retrieve, phoneNo: 要检索的用户电话号码
// Returns: User object and error if any, 返回: User对象和可能的错误
func GetUserByPhoneNo(phoneNo string) (user User, err error) {
	o := orm.NewOrm()
	usert := new(User)
	err = o.QueryTable(usert).Filter("Telphone", phoneNo).Filter("Status", 1).One(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetUserById retrieves a user by ID.
// GetUserById 根据ID检索用户。
// id: ID of the user to retrieve, id: 要检索的用户ID
// Returns: User object and error if any, 返回: User对象和可能的错误
func GetUserById(id int64) (user User, err error) {
	// user = User{Id: id}
	// o := orm.NewOrm()
	// o.Read(&user, "Id") //默认填充所有字段
	// return user
	o := orm.NewOrm()
	usert := new(User)
	err = o.QueryTable(usert).Filter("Id", id).One(&user, "id", "sessionid", "username", "nickname", "telphone", "email", "department", "company", "lastloginip", "lastlogintime", "remark", "status", "createtime") // 仅返回 Id 和 Name
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetuList retrieves a list of users with basic information.
// GetuList 检索包含基本信息的用户列表。
// Returns: list of users as orm.Params and total count, 返回: 以orm.Params形式的用户列表和总数
func GetuList() (obj []orm.Params, count int64) {
	o := orm.NewOrm()
	sqlstr := "SELECT username,nickname,department FROM user"
	count, err := o.Raw(sqlstr).Values(&obj)
	if err == nil && count > 0 {
		return obj, count
	}
	return nil, 0
}

// UeditMyInfo updates personal information of a user.
// UeditMyInfo 更新用户的个人信息。
// u: pointer to the User object containing updated information, u: 指向包含更新信息的User对象的指针
// Returns: number of affected rows and error if any, 返回: 受影响的行数和可能的错误
func UeditMyInfo(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	user := make(orm.Params)

	if len(u.Nickname) > 0 {
		user["Nickname"] = u.Nickname
	}
	if len(u.Email) > 0 {
		user["Email"] = u.Email
	}
	if len(u.Telphone) > 0 {
		user["Telphone"] = u.Telphone
	}
	if len(u.Company) > 0 {
		user["Company"] = u.Company
	}
	if len(u.Department) > 0 {
		user["Department"] = u.Department
	}

	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table User
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(user)
	return num, err
}

// CheckExitUserOrUpdate checks if a user exists and updates or creates the user accordingly.
// CheckExitUserOrUpdate 检查用户是否存在并相应地更新或创建用户。
// u: pointer to the User object to check or create, u: 指向要检查或创建的User对象的指针
// Returns: User object and error if any, 返回: User对象和可能的错误
func CheckExitUserOrUpdate(u *User) (User, error) {
	user, err := GetUserByUsername(u.Username)
	if err != nil {
		//return 0,err
		log.Printf("[ERROR] CheckExitUserOrUpdate getuserByUsername Fail:%s\n", err)
	}
	//没有用户就添加
	if user.Id == 0 {
		_, err := AddUser(u) //添加用户
		if err != nil {
			log.Printf("[ERROR] CheckExitUserOrUpdate AddUser Fail:%s\n", err)
			return user, err
		}
		xuser, err2 := GetUserByUsername(u.Username) //插入后查询是否存在用户
		if err2 != nil {
			log.Printf("[ERROR] CheckExitUserOrUpdate GetUser Fail:%s\n", err2)
			return user, err2
		}
		if xuser.Id == 0 {
			return user, errors.New("CheckExitUserOrUpdate GetUserId is 0")
		}
		_, err3 := AddRoleUser(6, xuser.Id) //6的角色ID为只读角色
		if err3 != nil {
			log.Printf("[ERROR] CheckExitUserOrUpdate AddRoleUser Fail:%s\n", err3)
			//return 0, err
		}
		return xuser, nil
	}

	//有用户更新密码
	u.Id = user.Id
	_, err4 := UpdateUser(u)
	if err4 != nil {
		log.Printf("[ERROR] CheckExitUserOrUpdate UpdateUser Fail:%s\n", err4)
		//return 0, err
	}
	return user, nil
}
