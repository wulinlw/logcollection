package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
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

func (this *LogModel) TableName(tbn string) string {
	return tbn
}

//func (this *LogModel) TableName() string {
//	return "log_127.0.0.1_apache"
//}

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
		fmt.Println(maps[0]["id"])
	}
}

func Queryone() {
	db, err := orm.GetDB()
	if err != nil {
		fmt.Println("get default DataBase")
	}
	fmt.Println(db)

	o := orm.NewOrm()
	//	var r orm.RawSeter
	//	r = o.Raw("SELECT id FROM `log_127.0.0.1_apache` WHERE id = 12")
	//	var maps []orm.Params
	//	r.Values(&maps)
	//	fmt.Println(maps)

	logs := new(LogModel)
	logs.Aid = 1
	logs.From = "apache"
	logs.File_name = "test"
	logs.Crtime = 123456789
	logs.Line = 111
	logs.Content = "test str"
	logs.TableName("log_apache") //设置表名有问题
	o.Insert(logs)

}
