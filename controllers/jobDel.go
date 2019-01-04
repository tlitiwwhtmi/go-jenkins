package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"go-jenkins/service/database"
	"go-jenkins/service/schedule"
	"go-jenkins/utils"
)

type JobDel struct {
	beego.Controller
}

func (c *JobDel) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		jobId, err := c.GetInt64("id")
		if err != nil {
			c.Ctx.WriteString("wrong param")
			return
		}

		//check schedulejobs
		schedule.ScheduledJobDelete(jobId)

		err = database.RemoveJobById(jobId)
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("delete error")
			return
		}

		//don't delete the histories
		/*err = database.RemoveHistoryByJob(jobId)
		if err != nil{
			fmt.Println(err.Error())
			c.Ctx.WriteString("delete error")
			return;
		}*/

		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("error")
}
