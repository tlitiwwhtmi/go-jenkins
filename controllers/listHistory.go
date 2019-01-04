package controllers

import (
	"fmt"
	"go-jenkins/service/database"
	"go-jenkins/utils"

	"github.com/astaxie/beego"
)

type ListHistory struct {
	beego.Controller
}

func (c *ListHistory) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		pageIndex, err := c.GetInt("page")
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		jobId, err := c.GetInt64("jobid")
		if err != nil {
			fmt.Println(err.Error())
		}
		currentId, err := c.GetInt64("cid")
		if err != nil {
			fmt.Println(err.Error())
		}
		perPage := 10

		histories, err := database.ListHistories(jobId, currentId, pageIndex, perPage)
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}

		c.Data["json"] = &histories
		c.ServeJSON()
		return

	}
	c.Ctx.WriteString("error")
}
