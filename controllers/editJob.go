package controllers

import (
	"go-jenkins/service/database"
	"go-jenkins/utils"
	"time"

	"github.com/astaxie/beego"
	"go-jenkins/service/schedule"
)

type EditJob struct {
	beego.Controller
}

func (c *EditJob) Post() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		jobId, err := c.GetInt64("jobid")
		if err != nil {
			c.Ctx.WriteString("param wrong")
			return
		}
		isGitLab, err := c.GetInt("gitlab")
		if err != nil {
			c.Ctx.WriteString("param wrong")
			return
		}
		autoBuild, err := c.GetInt("autobuild")
		if err != nil {
			c.Ctx.WriteString("param wrong")
			return
		}
		deployType, err := c.GetInt("deploytype")
		if err != nil {
			c.Ctx.WriteString("param wrong")
		}

		job, err := database.GetJobById(jobId)
		if err != nil {
			c.Ctx.WriteString("not exist")
			return
		}

		//check if is existed
		jobs, err := database.GetJobByPIdAndBranch(job.ProjectId, c.GetString("branch"))
		if c.GetString("build") == "1" {
			if err == nil {
				for _, tempJob := range *jobs {
					if tempJob.Build == c.GetString("build") && tempJob.Id != job.Id && tempJob.PomPath == job.PomPath && tempJob.Profile == job.Profile {
						c.Ctx.WriteString("already exist:" + tempJob.JobName)
						return
					}
				}
			}
		}
		if c.GetString("build") == "2" {
			if err == nil {
				for _, tempJob := range *jobs {
					if tempJob.Shell == c.GetString("shell") && tempJob.Build == c.GetString("build") && tempJob.Id != job.Id && tempJob.PomPath == job.PomPath && tempJob.Profile == job.Profile {
						c.Ctx.WriteString("already exist:" + tempJob.JobName)
						return
					}
				}
			}
		}

		//check if is existed
		//		jobs, err := database.GetJobByPIdAndBranch(job.ProjectId, c.GetString("branch"))
		//		if err == nil {
		//			for _, tempJob := range *jobs {
		//				if tempJob.Shell == c.GetString("shell") {
		//					c.Ctx.WriteString("already exist:" + tempJob.JobName)
		//					return
		//				}
		//			}
		//		}

		job.BranchName = c.GetString("branch")
		job.Build = c.GetString("build")
		job.PomPath = c.GetString("pompath")
		job.Profile = c.GetString("profile")
		job.Deploy = deployType
		job.JobName = c.GetString("name")
		job.IsGitlab = isGitLab
		job.Email = c.GetString("emails")
		job.Shell = c.GetString("shell")
		job.AutoBuild = autoBuild
		job.AutoTime = c.GetString("autotime")
		job.ModifyTime = time.Now()
		job.Modifier = user.Account
		_, err = database.UpdateJob(job)
		if err != nil {
			c.Ctx.WriteString("save wrong")
			return
		}
		schedule.ScheduledJob(job)
		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("not login")
}
