package controllers

import (
	"github.com/astaxie/beego"
	"go-jenkins/service/database"
	"go-jenkins/utils"
)

type GetPortalData struct {
	beego.Controller
}

func (c *GetPortalData) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		jobCount, err := database.GetJobCount()
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		pacCount, err := database.GetPacCount()
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		histories, err := database.GetAllHistories()
		if err != nil {
			c.Ctx.WriteString("error")
			return
		}
		var totalDuration int64
		totalDuration = 0
		totalFailure := 0
		for _, history := range *histories {
			totalDuration += history.Duration
			if history.Status == "FAILURE" {
				totalFailure += 1
			}
		}

		userCount := database.GetUsersCount()

		c.Data["json"] = &portalData{jobCount, pacCount, totalDuration / int64(len(*histories)), float32(totalFailure) / float32(len(*histories)), userCount}
		c.ServeJSON()
		return
	}
	c.Ctx.WriteString("error")
}

type portalData struct {
	JobCount  int64
	PacCount  int64
	AveTime   int64
	Faiure    float32
	UserCount int
}
