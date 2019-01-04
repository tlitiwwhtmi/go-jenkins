package controllers

import (
	"go-jenkins/service"

	"encoding/json"
	"github.com/astaxie/beego"
	"go-jenkins/utils"
	"strings"
)

type LoginAction struct {
	beego.Controller
}

func (c *LoginAction) Post() {
	if user := utils.GetUserFromSession(c.Ctx); user != nil {
		c.DelSession("loginUser")
	}
	//获取用户名和错误
	user, bderr := service.Login(c.GetString("user"), c.GetString("pass"))
	if bderr != nil {
		c.Data["json"] = &loginMsg{"failure", bderr.Err.Error()}
		c.ServeJSON()
		return
	}
	service.RawUserWithGit(user)
	//c.SetSession("loginUser", user)
	body, _ := json.Marshal(user)
	c.Ctx.SetCookie("user", utils.Base64Encode([]byte(body)))
	refer := "/"
	if c.GetSession("refer") != nil {
		refer = c.GetSession("refer").(string)
		if refer == "" || strings.Index(refer, "spelog") < 0 {
			refer = "/"
		}
	}
	c.Data["json"] = &loginMsg{"success", refer}
	c.ServeJSON()
	return
}

type loginMsg struct {
	Status string
	Msg    string
}
