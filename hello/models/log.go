package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Tb string = "log_tpl"
)

type LogModel struct {
	Id        int
	Aid       int
	From      string
	File_name string
	Crtime    int
	Line      int
	Content   string
}

func (this *LogModel) TableName() string {
	return this.Tbn("default")
}

//func (this *LogModel) TableName() string {
//	return "log_127.0.0.1_apache"
//}

func (this *LogModel) Tbn(tbn string) string {
	//	if tbn == "default" {
	//		Tb = "log_tpl"
	//	}
	fmt.Println(Tb)
	return Tb
}

func init() {
	orm.RegisterModel(new(LogModel))
	//注册驱动
	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", "root:@/log_gather?charset=utf8", 30)
}

func Test() {
	fmt.Println("this in model")
	orm.Debug = true
	o := orm.NewOrm()
	//o.Using("default")
	var maps []orm.Params
	num, _ := o.Raw("SELECT id FROM `log_127.0.0.1_apache` WHERE id = 12").Values(&maps)
	if num > 0 {
		fmt.Println("==============>", maps[0]["id"])
		for k, v := range maps {
			fmt.Println(k, v["from"])
		}
	}
}

//获取所有app
func GetAllApp() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	num, _ := o.Raw("SELECT * FROM `app`").Values(&maps)
	if num > 0 {
		//		fmt.Println("==============>", maps[0]["from"])
		//		for k, v := range maps {
		//			fmt.Println(k, v["from"])
		//		}
		return maps
	}
	return maps
}

//获取某一个app的某一页log
func getAppLogPage(appid, page) {

}

//NewQueryBuilder稳定版不存在，属于开发分支
func Queryone() {
	//	var logs []LogModel
	//	qb, _ := orm.NewQueryBuilder("mysql")

}
