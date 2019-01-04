package bdemail

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/xanzy/go-gitlab"
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/utils"
	"net/smtp"
)

func sendEmail(addrs *[]string, content string) {
	c, err := smtp.Dial("10.32.135.32:25")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer c.Close()
	c.Mail("building@chexiang.com")
	for _, addr := range *addrs {
		c.Rcpt(addr)
	}
	writer, err := c.Data()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer writer.Close()
	content = "To: " + utils.CoverArray2str(addrs) + "\r\n" + content
	buf := bytes.NewBufferString(content)
	if _, err = buf.WriteTo(writer); err != nil {
		fmt.Println(err.Error())
	}

}

func SendEmails(addrs *[]string, history *bdmodels.BdBuildHistory, job *bdmodels.BdJob, gitProject *gitlab.Project) {
	if addrs == nil {
		return
	}
	emailContent := generateEmailContent(history, job, gitProject)
	go sendEmail(addrs, emailContent)
	/*for _, addr := range *addrs {
		if addr != "" {
			go sendEmail(addr, emailContent)
		}
	}*/

}

func generateEmailContent(history *bdmodels.BdBuildHistory, job *bdmodels.BdJob, gitProject *gitlab.Project) string {
	status := ""
	if history.Status == "ABORTED" {
		status = "取消"
	}
	if history.Status == "SUCCESS" {
		status = "成功"
	}
	if history.Status == "FAILURE" {
		status = "失败"
	}
	if history.Status == "" {
		status = "异常"
	}

	emailContent := "Subject: " + *gitProject.Name + "分支" + job.BranchName + "编译" + status + "! \r\n \r\n"
	_, err := beego.AppConfig.Int("httpport")
	if err != nil {
		fmt.Println(err.Error())
	}
	emailContent += "您好! \r\n您在building上的任务:" + job.JobName + "的编译" + status + " \r\n"
	emailContent += "状态:" + history.Status + "\r\n"
	emailContent += "时间:" + history.StartTime.String() + "\r\n"
	emailContent += "耗时:" + utils.CoverDurationTime(history.Duration) + " \r\n"
	emailContent += "操作者:" + history.BuildExecutor + "\r\n"

	//add pacakage list
	pacs, err := database.GetPacListByHistory(history.Id)
	if err == nil && len(*pacs) > 0 {
		emailContent += "生成包： \r\n"
		for _, pac := range *pacs {
			emailContent += "\t " + pac.Name + "\r\n"
		}
	} else {
		emailContent += "暂时无法获取生成包信息 \r\n"
	}

	if err != nil {
		emailContent += "请登录系统查看详细日志!"
	} else {
		//emailContent += "日志访问地址:http://" + utils.GetLocalIp() + ":" + fmt.Sprintf("%d", port) + "/spelog?id=" + fmt.Sprintf("%d", history.Id) + "\r\n"
		emailContent += "日志访问地址:http://build.dds.com/spelog?id=" + fmt.Sprintf("%d", history.Id) + "\r\n"
	}
	return emailContent
}
