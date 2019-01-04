package controllers

import (
	"fmt"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
	"go-jenkins/service"
)

type JobRun struct {
	beego.Controller
}

func (c *JobRun) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {

		utils.GetLogger().Info("START RUNNING JOB ")

		jobId, err := c.GetInt64("id")
		if err != nil {
			utils.GetLogger().Info("FAILED RUNNING JOB: missing param jobId")
			fmt.Println(1, err.Error())
			c.Ctx.WriteString("error")
			return
		}
		msg, _ := service.StartBuild(jobId, user)
		if msg != "" {
			c.Ctx.WriteString(msg)
			return
		}
		//bdJob, err := database.GetJobById(jobId)
		//if err != nil {
		//	fmt.Println(2, err.Error())
		//	c.Ctx.WriteString("error")
		//	return
		//}
		////check if is running
		//isRunning, _, _ := bdjenkins.GetJobStatus(bdJob)
		//if isRunning {
		//	c.Ctx.WriteString("is running")
		//	return
		//}
		//
		//branch, err := bdgitlab.GetSingleBranch(int(bdJob.ProjectId), bdJob.BranchName)
		//if err != nil {
		//	c.Ctx.WriteString("branch not exist")
		//	return
		//}
		//
		//history := new(bdmodels.BdBuildHistory)
		//history.JobId = bdJob.Id
		//history.StartTime = time.Now()
		//history.BuildExecutor = user.Account
		//history.Version = branch.Commit.ID
		//history.CommitAuthor = branch.Commit.AuthorName
		//rs := []rune(branch.Commit.Message)
		//if len(rs) > 20{
		//	history.Message = string(rs[0:19])
		//}else{
		//	history.Message = branch.Commit.Message
		//}
		//
		//hId, err := database.AddHistory(history)
		//if err != nil {
		//	fmt.Println(3, err.Error())
		//	c.Ctx.WriteString("error")
		//	return
		//}
		//
		//history.Id = hId
		//_, err = bdjenkins.RunJob(bdJob, history, c.Ctx.Input.Port())
		//if err != nil {
		//	fmt.Println(err.Error())
		//	c.Ctx.WriteString("error")
		//	return
		//}
		c.Ctx.WriteString("success")
		return
	}
	c.Ctx.WriteString("-1")
}
