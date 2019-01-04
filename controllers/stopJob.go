package controllers

import (
	"fmt"
	"go-jenkins/service/database"
	"go-jenkins/service/jenkins"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
)

type StopJob struct {
	beego.Controller
}

func (c *StopJob) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {

		utils.GetLogger().Info("STARTING STOP JOB")

		jobId, err := c.GetInt64("id")
		if err != nil {

			utils.GetLogger().Error("STOP JOB: missing params jobId")

			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		bdJob, err := database.GetJobById(jobId)
		if err != nil {

			utils.GetLogger().Error("STOP JOB: get job failed")

			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		err = bdjenkins.StopJenkinsJob(bdJob)
		if err != nil {
			//do what you want
			utils.GetLogger().Error("STOP JOB: stop jenkins job failed " + err.Error())

			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}

		history, err := database.GetLatestHistory(jobId)
		if err != nil {
			utils.GetLogger().Error("STOP JOB: find history failed")

			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		if history.Status == "" {
			history.Status = "ABORTED"
			database.UpdateHistory(history)
		}

		utils.GetLogger().Info("STOP JOB SUCCESS")
		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("error")
}
