package admin

import (
	"fmt"
	"mime"
	"os"

	. "xkube/xadmin/src/lib"
	"xkube/xadmin/src/models"

	beego "github.com/beego/beego/v2/server/web"
)

func Run() {
	initialize()
}
func initialize() {
	fmt.Println("Starting....")
	mime.AddExtensionType(".css", "text/css")
	//判断初始化参数
	initArgs()
	//初始化计数
	InitCnt()
	models.Connect()
	router()
	beego.AddFuncMap("stringsToJson", StringsToJson)
}

func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.SyncdbInfo()
			os.Exit(0)
		}
	}
}
