package bdmodels

import "github.com/astaxie/beego/orm"

type BdPackage struct {
	Id        int64  `orm:"pk;column(id)"`
	HistoryId int64  `orm:"column(history_id)"`
	Name      string `orm:"column(name)"`
}

func init() {
	orm.RegisterModel(new(BdPackage))
}
