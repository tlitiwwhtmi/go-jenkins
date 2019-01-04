package service

import (
	"fmt"
	"go-jenkins/models/bd"
	"go-jenkins/service/bdlog"
	"go-jenkins/service/database"
	"go-jenkins/service/jenkins"
)

func RawJob(jobs *[]bdmodels.BdJob) {
	for i, _ := range *jobs {
		job := &(*jobs)[i]
		count, historys, err := database.GetHistoryByJob(job.Id)
		if err != nil || count == 0 {
			continue
		}
		history := (*historys)[0]
		job.HistoryId = history.Id
		job.LastBuild = history.StartTime
		job.Duration = history.Duration
		job.Status = history.Status
		job.Version = history.Version
		job.CommitAuthor = history.CommitAuthor
		job.CommitMessage = history.Message
		job.StartTime = history.StartTime
		if count >= 2 {
			job.AveDura = int((*historys)[1].Duration)
		}
		//job.BuildTimes = count
		//job.Log = history.Log
	}
}

func GetJobLog(jobId int64) (string, string, bool) {
	job, err := database.GetJobById(jobId)
	if err != nil {
		fmt.Println(err.Error())
		return "", "", false
	}
	isRunning, _, log := bdjenkins.GetJobStatus(job)
	if isRunning {
		return *log, "", isRunning
	}

	count, historys, err := database.GetHistoryByJob(job.Id)
	if err != nil || count == 0 {
		return "", "", false
	}
	history := (*historys)[0]
	log, err = bdlog.ReadLog(history.Log)
	if err != nil {
		return "", "", false
	}
	return *log, history.Status, false
}
