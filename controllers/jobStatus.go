package controllers

import (
	"fmt"
	"go-jenkins/models/bd"
	"go-jenkins/service"
	"go-jenkins/service/database"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
	"time"
)

type JobStatus struct {
	beego.Controller
}

func (c *JobStatus) Post() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		//fmt.Println(c.GetStrings("ids"), c.GetString("data"))
		//pageJobs := new([]bdmodels.BdJob)
		//json.Unmarshal([]byte(c.GetString("data")), pageJobs)

		jobs, err := getJobs(c.GetString("ids"))
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		//for !checkJobData(jobs, pageJobs) {
		//	jobs, err = getJobs(c.GetString("ids"))
		//	if err != nil {
		//		c.Ctx.WriteString("error")
		//		return
		//	}
		//	fmt.Println(123)
		//}

		c.Data["json"] = jobs
		c.ServeJSON()
		return
		//		jobId, err := c.GetInt64("id")
		//		if err != nil {
		//			fmt.Println(err.Error())
		//			c.Ctx.WriteString("error")
		//			return
		//		}
		//		job, err := database.GetJobById(jobId)
		//		if err != nil {
		//			fmt.Println(err.Error())
		//			c.Ctx.WriteString("error")
		//			return
		//		}
		//		history, err := database.GetLatestHistory(job.Id)
		//		if err != nil {
		//			fmt.Println(err.Error())
		//			c.Ctx.WriteString("error")
		//			return
		//		}
		//		if history != nil {
		//			job.HistoryId = history.Id
		//			job.Status = history.Status
		//			job.LastBuild = history.StartTime
		//			job.Duration = history.Duration
		//			job.Log = history.Log
		//		}
		//		isRunning, msg := bdjenkins.GetJobStatus(job)
		//		job.IsRunning = isRunning
		//		job.What = msg
		//		c.Data["json"] = job
		//		c.ServeJson()
		//		return
	}
	c.Ctx.WriteString("error")
}

func getJobs(ids string) (*[]bdmodels.BdJob, error) {
	jobs, err := database.GetJobByIds(ids)
	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}

	service.RawJob(jobs)
	service.GetJobStatus(jobs)

	for i, _ := range *jobs {
		job := &(*jobs)[i]
		if job.IsRunning {
			spentTime := int(time.Now().Sub(job.StartTime).Seconds() * 1000)
			if job.AveDura == 0 {
				job.AveDura = 3 * 60 * 1000
			}
			if spentTime >= job.AveDura {
				job.Progress = 0.99
			} else {
				job.Progress = float32(spentTime) / float32(job.AveDura)
			}
		}
	}

	return jobs, nil
}

func checkJobData(jobs *[]bdmodels.BdJob, pageJobs *[]bdmodels.BdJob) bool {

	for _, job := range *jobs {
		for _, tempjob := range *pageJobs {
			if job.Id == tempjob.Id {

				if job.Duration != tempjob.Duration || job.HistoryId != tempjob.HistoryId || job.IsRunning != tempjob.IsRunning || job.JobName != tempjob.JobName || !job.LastBuild.Equal(tempjob.LastBuild) || job.Log != tempjob.Log || job.Modifier != tempjob.Modifier || !job.ModifyTime.Equal(tempjob.ModifyTime) || job.Shell != tempjob.Shell || job.Status != tempjob.Status {
					return true
				}
			}
		}
	}
	return false
}
