package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/xanzy/go-gitlab"
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/service/gitlab"
	"go-jenkins/utils"
)

type ListPackage struct {
	beego.Controller
}

func (c *ListPackage) Get() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		projectId, err := c.GetInt64("pid")
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		cId, err := c.GetInt64("cid")
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		pageIndex, err := c.GetInt("page")
		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}
		branch := c.GetString("branch")

		keyword := c.GetString("q")

		pacs := new([]bdmodels.BdPackage)

		if "" != keyword {
			pacs, err = database.SearchPacsByPacName(projectId, pageIndex, 10, keyword)
		} else {
			pacs, err = database.GetPacsByProject(projectId, cId, pageIndex, 10, branch)
		}

		if err != nil {
			fmt.Println(err.Error())
			c.Ctx.WriteString("error")
			return
		}

		var hIds []int64
		for _, pac := range *pacs {
			if isExist(pac.HistoryId, hIds) {
				continue
			} else {
				hIds = append(hIds, pac.HistoryId)
			}
		}

		branches, err := bdgitlab.GetBranchByProjectId(c.GetString("pid"))

		if len(hIds) > 0 {
			histories, err := database.GetHistoriesByIdsnolog(hIds)
			if err != nil {
				fmt.Println("get his", err.Error())
				c.Ctx.WriteString("error")
				return
			}
			var gPacs []genPac
			for _, history := range *histories {
				gPac := new(genPac)
				job, err := database.GetJobById(history.JobId)
				if err != nil {
					fmt.Println("get job", err.Error())
					c.Ctx.WriteString("error")
					return
				}
				history.Branch = job.BranchName
				gPac.History = history
				gPac.Pacs = getPacs(pacs, history.Id)
				gPacs = append(gPacs, *gPac)
			}

			c.Data["json"] = &genData{len(*pacs), gPacs, branches}
			c.ServeJSON()
			return
		}

		c.Data["json"] = &genData{0, nil, branches}

		c.ServeJSON()
		return
	}
	c.Ctx.WriteString("not login")
}

func getPacs(pacs *[]bdmodels.BdPackage, hId int64) []bdmodels.BdPackage {
	var gPacs []bdmodels.BdPackage
	for _, pac := range *pacs {
		if pac.HistoryId == hId {
			gPacs = append(gPacs, pac)
		}
	}
	return gPacs
}

func isExist(id int64, ids []int64) bool {
	for _, tempId := range ids {
		if tempId == id {
			return true
		}
	}
	return false
}

type genData struct {
	Total    int
	GenPacs  []genPac
	Branches []*gitlab.Branch
}

type genPac struct {
	History bdmodels.BdBuildHistory
	Pacs    []bdmodels.BdPackage
}
