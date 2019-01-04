package controllers

import (
	"fmt"
	"go-jenkins/service/database"
	"go-jenkins/service/jenkins"

	"github.com/astaxie/beego"
	"go-jenkins/service"
	"go-jenkins/service/bdlog"
	"go-jenkins/utils"
)

type JobTrriger struct {
	beego.Controller
}

func (c *JobTrriger) Get() {
	fmt.Println("in trriger")
	hId, err := c.GetInt64("id")
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	history, err := database.GetHistoryById(hId)
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	job, err := database.GetJobById(history.JobId)
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	build, err := bdjenkins.GetBuildInfo(job)
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	_, err = build.Stop()
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	history.Status = build.GetResult()

	buildLog := build.GetConsoleOutput()
	fileName := utils.OrgniseUniqName()
	go bdlog.AddLog(&buildLog, fileName) //here need an error charge

	history.Log = fileName
	fmt.Println(build.GetResult())

	_, err = build.Stop()
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	_, err = build.Poll()
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}
	history.Duration = build.GetDuration()
	_, err = database.UpdateHistory(history)
	if err != nil {
		fmt.Println(err.Error())
		c.Ctx.WriteString("error")
		return
	}

	//send Email to users
	//发送邮件放到了下面的TrrigerBuilder中.先检查image的build状态之后再发现邮件
	//go service.SendEmail(history)

	//trigger builder
	go service.TrrigerBuilder(history)

	_, err = bdjenkins.DeleteJenkinsJob(job)
	if err != nil {
		fmt.Println(1, err.Error())
		c.Ctx.WriteString("error")
		return
	}
	fmt.Println("delete job success")
	c.Ctx.WriteString("success")
}
