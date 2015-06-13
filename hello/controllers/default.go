package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"hello/models"
)

type MainController struct {
	beego.Controller
}

func init() {
	//启动的时候生效,方法请求时不生效
	fmt.Println("init !!! ")
}
func Prepare(this *MainController) {
	//自动路由的时候不生效
	fmt.Println("Prepare !!!")

}
func (this *MainController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplNames = "index.tpl"
}

func (this *MainController) Test() {
	//this.Layout = "layout/layout.html"
	//this.Ctx.WriteString("hello")
	this.TplNames = "main/test.tpl"
}

func (this *MainController) Index() {
	//this.Layout = "layout/layout.html"
	//this.LayoutSections = make(map[string]string)
	//this.LayoutSections["menu"] = "layout/menu.html"
	//this.Ctx.WriteString("hello")
	layout(this)

	models.Test()
	models.Queryone()

	this.TplNames = "main/test.tpl"
}

//设置layout
func layout(this *MainController) {
	this.Layout = "layout/layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["menu"] = "layout/menu.html"
}
