package main

import (
	"github.com/astaxie/beego"
	_ "hello/routers"

)

func main() {
	beego.Run()
	//orm.RegisterDataBase("default", "mysql", "root:@/log_gather?charset=utf8", 30)
}
