package routers

import (
	"github.com/astaxie/beego"
	"hello/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.MainController{})

	//beego.SetStaticPath("/images", "images")
	//beego.SetStaticPath("/css", "css")
	//beego.SetStaticPath("/js", "js")
}
