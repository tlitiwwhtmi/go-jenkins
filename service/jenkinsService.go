package service

import (
	"go-jenkins/models/bd"
	"go-jenkins/service/jenkins"
)

func GetJobStatus(jobs *[]bdmodels.BdJob) {
	for i, _ := range *jobs {
		job := &(*jobs)[i]
		isRunning, msg, _ := bdjenkins.GetJobStatus(job)
		job.IsRunning = isRunning
		job.What = msg
		if isRunning {
			//job.Log = *log
		}
	}
}
