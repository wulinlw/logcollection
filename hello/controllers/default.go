package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"hello/models"
	"hello/utils"
	"net"
	"strconv"
	"strings"
	"time"
)

type MainController struct {
	beego.Controller
}
type queryParams struct {
	Id    int    `form:"id"`
	Stime string `form:"stime"`
	Etime string `form:"etime"`
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

	//models.Test()

	//models.Queryone()
	maps := models.GetAllApp()
	for k, v := range maps {
		fmt.Println(k, v["from"])
	}

	this.Data["apps"] = maps
	this.Data["a"] = "aaa"
	this.TplNames = "main/test.tpl"
}

//设置layout
func layout(this *MainController) {
	this.Layout = "layout/layout.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["menu"] = "layout/menu.html"
}

func (this *MainController) Query() {
	params := queryParams{}
	if err := this.ParseForm(&params); err != nil {
		//handle error

	}

	tm2, _ := time.Parse("2006-01-02 15:04:05", params.Stime)

	fmt.Println(tm2.Unix())

	fmt.Println(params)
	projectInfo := models.GetAppById(params.Id)
	fmt.Println(projectInfo)
	tableName := models.GetTableName(projectInfo)
	//	fmt.Println(tableName)
	//	paginator := SetPaginator(3, 13)
	//	this.TplNames = "main/query.tpl".

	pno, _ := this.GetInt("pno") //获取当前请求页
	var tlog []models.LogModel
	var conditions string = " and crtime between " + strconv.FormatInt(timeStr2Unix(params.Stime), 10) + " and " + strconv.FormatInt(timeStr2Unix(params.Etime), 10) + " order by id desc" //定义日志查询条件,格式为 " and name='zhifeiya' and age=12 "
	var po pager.PageOptions                                                                                                                                                               //定义一个分页对象
	po.TableName = "`" + tableName + "`"                                                                                                                                                   //指定分页的表名
	po.EnableFirstLastLink = true                                                                                                                                                          //是否显示首页尾页 默认false
	po.EnablePreNexLink = true                                                                                                                                                             //是否显示上一页下一页 默认为false
	po.Conditions = conditions                                                                                                                                                             // 传递分页条件 默认全表
	po.Currentpage = int(pno)                                                                                                                                                              //传递当前页数,默认为1
	po.PageSize = 3                                                                                                                                                                        //页面大小  默认为20

	totalItem, totalpages, rs, htmlStr := pager.GetPagerLinks(&po, this.Ctx)
	rs.QueryRows(&tlog)      //把当前页面的数据序列化进一个切片内
	this.Data["list"] = tlog //把当前页面的数据传递到前台
	this.Data["pagerhtml"] = htmlStr
	this.Data["totalItem"] = totalItem
	this.Data["PageSize"] = po.PageSize
	this.Data["totalPages"] = totalpages
	this.Data["tableName"] = tableName
	this.Data["params"] = params
	//this.Data["htmlStr"] = htmlStr
	//fmt.Println(totalItem, totalpages, tlog)
	layout(this)
	this.TplNames = "main/query.tpl"

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

func timeStr2Unix(timeStr string) int64 {
	tm2, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return tm2.Unix()
}
