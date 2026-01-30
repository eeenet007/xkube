package main

import (
	//"xkube/informers"
	"fmt"
	_ "xkube/routers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/redis"
)

const VERSION = "4.0"

var GoVersion string
var BuildTime string

func main() {
	//go informers.WatchDeploy("zx-k8s")
	fmt.Println("xkube:" + VERSION)
	fmt.Println("BuildTime: " + BuildTime)
	fmt.Println("GoLangVersion: " + GoVersion)
	fmt.Println("Start ok")
	beego.SetStaticPath("/", "views/front")
	beego.Run()
}
