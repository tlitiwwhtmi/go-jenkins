package database

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//userName := beego.AppConfig.String("user")
	//passWord := beego.AppConfig.String("password")
	//url := beego.AppConfig.String("url")
	//defaultdb := beego.AppConfig.String("db")
	//orm.RegisterDataBase("default", "mysql", userName+":"+passWord+"@tcp("+url+")/"+defaultdb)

	orm.RegisterDataBase("default", "mysql", "pm:pm@@1234@tcp(10.32.135.62:3306)/pm?loc=Local")
}

func getConnection() *orm.Ormer {
	o := orm.NewOrm()
	//o.Using(beego.AppConfig.String("db"))
	o.Using("pm")
	return &o
}
