package controllers

import (
	"github.com/astaxie/beego"
)

type PortalController struct {
	beego.Controller
}

func (c *PortalController) Get() {
	c.TplName = "portal.html"
}
