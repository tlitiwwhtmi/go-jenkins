package bdmodels

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type BdBuildHistory struct {
	Id            int64     `orm:"pk;column(id)"`
	JobId         int64     `orm:"column(job_id)"`
	StartTime     time.Time `orm:"column(start_time)"`
	Duration      int64     `orm:"column(duration)"`
	Status        string    `orm:"column(status)"`
	DownloadTime  int       `orm:"column(download_time)"`
	BuildTime     int       `orm:"column(build_time)"`
	UploadTime    int       `orm:"column(upload_time)"`
	CommitId      string    `orm:"column(commit_id)"`
	BuildExecutor string    `orm:"column(build_executor)"`
	Log           string    `orm:"column(log)"`
	Version       string    `orm:"column(version)"`
	Message       string    `orm:"column(message)"`
	CommitAuthor  string    `orm:"column(commit_author)"`
	Branch        string    `orm:"-"`
}

func init() {
	orm.RegisterModel(new(BdBuildHistory))
}
