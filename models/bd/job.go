package bdmodels

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type BdJob struct {
	Id             int64     `orm:"pk;column(id)"`
	ProjectId      int64     `orm:"column(project_id)"`
	JobName        string    `orm:"column(job_name)"`
	BranchName     string    `orm:"column(branch_name)"`
	Shell          string    `orm:"column(shell)"`
	Email          string    `orm:"column(email)"`
	IsGitlab       int       `orm:"column(is_gitlab)"`
	CreatorAccount string    `orm:"column(creator_account)"`
	CreateTime     time.Time `orm:"column(create_time)"`
	Modifier       string    `orm:"column(modifier)"`
	ModifyTime     time.Time `orm:"column(modify_time)"`
	Build          string    `orm:"column(build)"`
	PomPath        string    `orm:"column(pom_path)"`
	Profile        string    `orm:"column(profile)"`
	Deploy         int       `orm:"column(deploy)"`
	AutoBuild      int       `orm:"column(auto_build)"`
	AutoTime       string    `orm:"column(auto_time)"`
	LastBuild      time.Time `orm:"-"`
	Duration       int64     `orm:"-"`
	Status         string    `orm:"-"`
	HistoryId      int64     `orm:"-"`
	IsRunning      bool      `orm:"-"`
	What           string    `orm:"-"`
	Log            string    `orm:"-"`
	Version        string    `orm:"-"`
	CommitMessage  string    `orm:"-"`
	CommitAuthor   string    `orm:"-"`
	StartTime      time.Time `orm:"-"`
	AveDura        int       `orm:"-"`
	Progress       float32   `orm:"-"`
	//BuildTimes      int64       `orm:"-"`
}

func init() {
	orm.RegisterModel(new(BdJob))
}
