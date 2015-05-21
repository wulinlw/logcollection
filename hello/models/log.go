package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

//`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
//  `aid` int(10) DEFAULT NULL,
//  `from` varchar(32) DEFAULT NULL,
//  `file_name` varchar(64) DEFAULT NULL,
//  `crtime` int(10) DEFAULT NULL,
//  `line` int(10) DEFAULT NULL,
//  `content`

type LogModel struct {
	Id        int
	Aid       int
	From      string
	File_name string
	Crtime    int
	Line      int
	Content   string
}

func init() {
	orm.RegisterModel(new(LogModel))
}
func (this *LogModel) test() {
	fmt.Println("this in model")
}
