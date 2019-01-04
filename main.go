//主函数入口
package main

import (
	_ "go-jenkins/routers"
	"go-jenkins/utils"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func main() {
	var loginFilter = func(ctx *context.Context) {
		ctx.Input.CruSession.Set("refer", ctx.Request.RequestURI)
		//不过滤的页面列表
		if utils.GetUserFromSession(ctx) == nil && ctx.Request.RequestURI != "/login" && ctx.Request.RequestURI != "/loginaction" && !strings.Contains(ctx.Request.RequestURI, "jobtrriger") && !strings.Contains(ctx.Request.RequestURI, "savepac") && !strings.Contains(ctx.Request.RequestURI, "png") && !strings.Contains(ctx.Request.RequestURI, "startbuild") {
			ctx.Redirect(302, "/login")
		}
	}
	var onlineFilter = func(ctx *context.Context) {
		if user := utils.GetUserFromSession(ctx); user != nil {
			ctx.Input.SetData("user", user)
			//ctx.Input.Data["user"] = user
		}
	}
	beego.InsertFilter("", beego.BeforeRouter, loginFilter)
	beego.InsertFilter("/*", beego.BeforeRouter, loginFilter)
	beego.InsertFilter("", beego.BeforeExec, onlineFilter)
	beego.InsertFilter("/*", beego.BeforeExec, onlineFilter)
	//schedule.CheckAllJob()

	utils.InitLogger()

	utils.GetLogger().Info("SYSTEM STARTED")

	beego.Run()
}
