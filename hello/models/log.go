package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"strconv"
	"strings"
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
func getAppLogPage(appid int, page int) {

}

func GetAppById(id int) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	num, _ := o.Raw("SELECT * FROM `app` where id=?", id).Values(&maps)
	if num > 0 {
		//		fmt.Println("==============>", maps[0]["from"])
		//		for k, v := range maps {
		//			fmt.Println(k, v["from"])
		//		}
		return maps
	}
	return maps
}

func GetTableName(projectInfo []orm.Params) (tableName2 string) {
	var tableName string = ""
	if projectInfo[0] != nil {
		ip, _ := strconv.Atoi(projectInfo[0]["ip"].(string)) //interface 转换成其他类型
		tableName = fmt.Sprintf("log_%s_%s", inet_ntoa(int64(ip)), projectInfo[0]["from"])
	}
	return tableName
}

//NewQueryBuilder稳定版不存在，属于开发分支
func Queryone() {
	//	var logs []LogModel
	//	qb, _ := orm.NewQueryBuilder("mysql")

}

// Convert uint to net.IP http://www.sharejs.com
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Convert net.IP to int64 ,  http://www.sharejs.com
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}
