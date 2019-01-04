package controllers

import (
	"github.com/astaxie/beego"
)

type LogoutAction struct {
	beego.Controller
}

func (c *LogoutAction) Post() {
	c.Ctx.SetCookie("user", "")
	//c.DelSession("loginUser")
	c.Ctx.WriteString("success")
}
