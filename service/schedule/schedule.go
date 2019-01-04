package schedule

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
	"go-jenkins/models/bd"
	"go-jenkins/service"
	"go-jenkins/service/database"
	"go-jenkins/utils"
	"net/http"
	"strconv"
)

var jobMap map[int64]*cron.Cron

func startJob(job *bdmodels.BdJob) {
	c := cron.New()
	c.AddFunc(job.AutoTime, func() {
		msg, _ := service.StartBuild(job.Id, nil)
		if msg != "" {
			fmt.Println("AutoBuild", job.Id, job.JobName, msg)
		}
	})
	c.Start()
	if jobMap == nil {
		jobMap = make(map[int64]*cron.Cron)
	}
	jobMap[job.Id] = c
	fmt.Println("Schedule:", len(jobMap), "jobs are in schedule")
}

func checkScheduleJob(job *bdmodels.BdJob) {
	if jobMap == nil {
		return
	}
	if v, ok := jobMap[job.Id]; ok {
		v.Stop()
	}
	delete(jobMap, job.Id)
	fmt.Println("Schedule:", len(jobMap), "jobs are in schedule")
}

func ScheduledJobDelete(jobId int64) {
	/*job,err := database.GetJobById(jobId)
	if err != nil{
		return;
	}
	if job.AutoBuild != 1{
		return;
	}
	checkScheduleJob(job)*/
	go http.Get(beego.AppConfig.String("schedule_address") + "/scheduledelete?id=" + fmt.Sprintf("%d", jobId))
}

func ScheduledJob(job *bdmodels.BdJob) {
	/*if job.AutoBuild != 1{
		return;
	}
	checkScheduleJob(job)
	startJob(job)*/
	go http.Get(beego.AppConfig.String("schedule_address") + "/schedulejob?id=" + fmt.Sprintf("%d", job.Id))
}

func ScheduledJobs(jobs *[]bdmodels.BdJob) {
	for _, job := range *jobs {
		ScheduledJob(&job)
	}
	utils.GetLogger().Info("SCHEDULE: " + strconv.Itoa(len(jobMap)) + "jobs are in schedule")
}

func CheckAllJob() {
	jobs, err := database.ListAutoBuildJob()
	if err != nil {
		fmt.Println("init autobuild job failure")
		utils.GetLogger().Error("INIT SCHEDULE FAILED: get jobs from database failed")
		return
	}
	ScheduledJobs(jobs)
}
