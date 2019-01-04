package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"strings"
)

type SavePackage struct {
	beego.Controller
}

func (c *SavePackage) Get() {
	pacPath := c.GetString("pac")
	pacName := ""
	if pacPath != "" {
		if index := strings.LastIndex(pacPath, "/"); index > 0 {
			rs := []rune(pacPath)
			end := len(rs)
			pacName = string(rs[index+1 : end])
		} else {
			pacName = pacPath
		}
	}
	hId, err := c.GetInt64("hid")
	if err != nil {
		fmt.Println("historyId is not present")
		c.Ctx.WriteString("")
		return
	}
	bdPackage := new(bdmodels.BdPackage)
	bdPackage.HistoryId = hId
	bdPackage.Name = pacName
	database.AddPackage(bdPackage)
	c.Ctx.WriteString("")
}
