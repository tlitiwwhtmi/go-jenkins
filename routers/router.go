package routers

import (
	"go-jenkins/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	/*
	登录页面
	*/
	beego.Router("/login", &controllers.LoginController{})
	/*
		登录动作,获取GitId和token
	*/
	beego.Router("/loginaction", &controllers.LoginAction{})
	/*
		登出,并删除cookie
	*/
	beego.Router("/logoutaction", &controllers.LogoutAction{})
	/*
		通过cookie获取user,通过user获取项目
	*/
	beego.Router("/getprojects", &controllers.GetProjects{})
	/*
		增加项目
	*/
	beego.Router("/addproject", &controllers.AddProject{})
	beego.Router("/ownedproject", &controllers.OwnedProject{})
	beego.Router("/deleteproject", &controllers.DeleteProject{})
	beego.Router("/listjob", &controllers.ListJob{})
	beego.Router("/addjob", &controllers.AddJob{})
	beego.Router("/editjob", &controllers.EditJob{})
	beego.Router("/jobrun", &controllers.JobRun{})
	beego.Router("/jobstatus", &controllers.JobStatus{})
	beego.Router("/jobtrriger", &controllers.JobTrriger{})
	beego.Router("/savepac", &controllers.SavePackage{})
	beego.Router("/stopjob", &controllers.StopJob{})
	beego.Router("/jobs", &controllers.JobsController{}, "post:Create")
	beego.Router("/logout", &controllers.LogOutController{})
	beego.Router("/getlog", &controllers.GetLogController{})
	beego.Router("/jobdel", &controllers.JobDel{})
	beego.Router("/lisyhistories", &controllers.ListHistory{})
	beego.Router("/getspelog", &controllers.GetspeLogController{})
	beego.Router("/spelog", &controllers.SpeLogController{})
	beego.Router("/listpac", &controllers.ListPackage{})

	beego.Router("/portal", &controllers.PortalController{})
	beego.Router("/portaldata", &controllers.GetPortalData{})
	beego.Router("/buildsday", &controllers.BuildsDay{})

	beego.Router("/startbuild", &controllers.StartBuild{})

}
