// lockuser.go
package rbac

import (
	"fmt"
	"log"
	"strings"
	"xkube/common"
	m "xkube/xadmin/src/models"
)

type LockInfo struct {
	LockKey  string `json:"lockKey"`
	LockTime int64  `json:"lockTime"`
}

func (this *UserController) LockAct() {

	userinfo := this.GetSession("userinfo")
	userip := this.Ctx.Request.Header.Get("x-forwarded-for")
	if userip == "" {
		userip = this.Ctx.Request.RemoteAddr
	}

	act := this.GetString("act")
	key := this.GetString("key")
	if act == "list" {
		respList := common.Keys("lock_*")
		var lockArry = []LockInfo{}
		for _, vv := range respList {
			ext, _ := common.Ttl(vv)
			//log.Println(ext)
			xitems := &LockInfo{
				LockKey:  strings.Replace(vv, "lock_", "", -1),
				LockTime: ext,
			}
			lockArry = append(lockArry, *xitems)
		}
		code := 0
		msg := "success"
		count := len(lockArry)

		if count == 0 {
			code = -1
			msg = "noLockInfo"
		}

		this.Data["json"] = &map[string]interface{}{"code": code, "msg": msg, "count": count, "data": &lockArry}
		this.ServeJSON()
	} else if act == "unLock" {
		resp := common.Del("lock_" + key)
		log.Printf("[INFO] unLock %s,result:%d\n", key, resp)

		content := fmt.Sprintf("info:%s", key)
		_ = m.InsertLogAudit(userip, userinfo.(m.User).Username, "unLock", "success", content)

		this.Rsp(true, "success")
		return
	} else if act == "longTermLock" {
		resp, _ := common.Persist("lock_" + key)
		log.Printf("[INFO] longTermLock %s,result:%v\n", key, resp)
		this.Rsp(true, "success")
		return
	}
}
