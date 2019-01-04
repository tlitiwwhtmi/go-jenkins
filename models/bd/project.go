package bdmodels

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type BdProject struct {
	Id          int64     `orm:"pk;column(id)"`
	ProjectId   int64     `orm:"column(project_id)"`
	UserAccount string    `orm:"column(user_account)"`
	CreateTime  time.Time `orm:"column(create_time)"`
	Language    string    `orm:"column(language)"`
	ProjectName string    `orm:"-"`
}

func init() {
	orm.RegisterModel(new(BdProject))
}
