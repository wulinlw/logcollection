package main

import (
	"github.com/astaxie/beego"
	_ "hello/routers"
	"time"
)

func main() {
	beego.AddFuncMap("showTime", showTime) //模板函数
	beego.Run()
	//orm.RegisterDataBase("default", "mysql", "root:@/log_gather?charset=utf8", 30)
}
func showTime(timeInt int) string {
	return time.Unix(int64(timeInt), 0).Format("2006-01-02 15:04:05")
}
