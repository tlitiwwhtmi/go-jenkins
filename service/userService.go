package service

import (
	"fmt"
	"go-jenkins/errors"
	"go-jenkins/models/ldap"
	"go-jenkins/service/database"
	"go-jenkins/service/ldap"

	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"go-jenkins/models/bd"
	"go-jenkins/service/bdemail"
	"go-jenkins/service/gitlab"
	"go-jenkins/service/jenkins"
	"go-jenkins/utils"
	"net/http"
	"strings"
	"time"
)

func Login(username, password string) (*ldapmodels.LdapUser, *bderrors.Bderror) {
	return bdldap.LdapLogin(username, password)
}

func RawUserWithGit(user *ldapmodels.LdapUser) {
	pqcon := database.NewPostGres(beego.AppConfig.String("postgressUrl"))
	gitUser, err := pqcon.GetUserByAccount(user.Account)
	if err != nil {
		fmt.Println("connect to gitlab db error")
		return

	}
	user.GitId = gitUser.GitId
	user.GitToken = gitUser.GitToken
}

func StartBuild(jobId int64, user *ldapmodels.LdapUser) (string, error) {
	bdJob, err := database.GetJobById(jobId)
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: query job failed")
		fmt.Println(2, err.Error())
		return "not exist", err
	}
	//check if is running
	isRunning, _, _ := bdjenkins.GetJobStatus(bdJob)
	if isRunning {
		utils.GetLogger().Error("FAILED RUNNING JOB: job is running")
		return "is running", nil
	}

	branch, err := bdgitlab.GetSingleBranch(int(bdJob.ProjectId), bdJob.BranchName)
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: get branch failed")
		return "branch not exist", err
	}

	history := new(bdmodels.BdBuildHistory)
	history.JobId = bdJob.Id
	history.StartTime = time.Now()
	if user != nil {
		history.BuildExecutor = user.Account
	} else {
		history.BuildExecutor = "autobuild"
	}
	history.Version = branch.Commit.ID
	history.CommitAuthor = branch.Commit.AuthorName
	rs := []rune(branch.Commit.Message)
	if len(rs) > 20 {
		history.Message = string(rs[0:19])
	} else {
		history.Message = branch.Commit.Message
	}

	hId, err := database.AddHistory(history)
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: save history failed")
		return "save history error", err
	}
	utils.GetLogger().Info("STARTING RUNNING JOB: success saving the history")
	history.Id = hId
	port, err := beego.AppConfig.Int("httpport")
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: get http port failed")
		return "get port error", err
	}
	_, err = bdjenkins.RunJob(bdJob, history, port)
	if err != nil {
		return "run error", err
	}
	return "", nil
}

func SendEmail(history *bdmodels.BdBuildHistory) {
	var addrs []string

	if history.BuildExecutor != "autobuild" {
		if addr := utils.Covert2Email(history.BuildExecutor); !checkAddrExist(&addrs, addr) {
			addrs = append(addrs, addr)
		}
	}
	job, err := database.GetJobById(history.JobId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	gitProject, err := bdgitlab.GetProjectById(job.ProjectId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer bdemail.SendEmails(&addrs, history, job, gitProject)

	if job.IsGitlab == 1 {
		members := bdgitlab.GetAuthedUserForProject(job.ProjectId)
		if members != nil {
			for _, member := range members {
				if !checkAddrExist(&addrs, member.Email) {
					addrs = append(addrs, member.Email)
				}
			}
		}
	}
	if job.Email != "" {
		if strings.IndexAny(job.Email, ";") < 0 {
			if strings.IndexAny(job.Email, "@") > 0 {
				if !checkAddrExist(&addrs, job.Email) {
					addrs = append(addrs, job.Email)
				}
			}
		} else {
			emails := strings.Split(job.Email, ";")
			for _, email := range emails {
				if strings.IndexAny(email, "@") > 0 {
					if !checkAddrExist(&addrs, email) {
						addrs = append(addrs, email)
					}
				}
			}
		}
	}

}

func TrrigerBuilder(history *bdmodels.BdBuildHistory) bool {
	outcome := false
	pacs, err := database.GetPacListByHistory(history.Id)
	if err == nil && len(*pacs) > 0 {
		pacakges := make([]string, len(*pacs))
		for i, pac := range *pacs {
			pacakges[i] = pac.Name
		}
		requestData := RequestData{
			Packages: pacakges,
		}
		jsonStr, err := json.Marshal(requestData)
		if err != nil {
			outcome = false
		}
		req, err := http.NewRequest("POST", beego.AppConfig.String("builder_url"), bytes.NewBuffer(jsonStr))
		if err != nil {
			outcome = false
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		r, err := client.Do(req)
		if err != nil {
			outcome = false
		}
		defer r.Body.Close()
		var result map[string]interface{}
		err = json.NewDecoder(r.Body).Decode(&result)
		if err != nil {
			outcome = false
		}
		fmt.Println(result["status"])

		outcome = true

	}

	//改为在这里检查image的build状态然后发送邮件
	if err == nil && len(*pacs) > 0 {
		pacName := ""
		for _, pac := range *pacs {
			if strings.Contains(pac.Name, ".war") {
				pacName = pac.Name
			}
		}
		if pacName != "" {
			index := strings.LastIndex(pacName, ".war")
			dirName := ""
			if index > 12 {
				rs := []rune(pacName)
				dirName = string(rs[index-12 : index-4])
			}

			if dirName != "" {
				startTime := time.Now()
				for {
					executeTime := time.Now()
					duration := executeTime.Sub(startTime)
					//如果运行时间超过5分钟则退出循环
					if duration.Minutes() > 5 {
						break
					}
					requestData := PageurlData{
						PackageUrl: "http://10.32.135.110/package/" + dirName + "/" + pacName,
					}
					jsonStr, err := json.Marshal(requestData)
					if err != nil {
						fmt.Println(err.Error())
						break
					}

					res, err := http.Post("http://darth.chexiang.com/builder/_/get-build-status/", "application/json;charset=utf-8", bytes.NewBuffer([]byte(jsonStr)))
					if err != nil {
						fmt.Println(err.Error())
						break
					}
					defer res.Body.Close()
					var result map[string]interface{}
					err = json.NewDecoder(res.Body).Decode(&result)
					if err != nil {
						fmt.Println(err.Error())
						break
					}
					if result["data"].(map[string]interface{})["status"] != "success" && result["data"].(map[string]interface{})["status"] != "fail" {
						//隔5S请求一次
						time.Sleep(5 * time.Second)
					} else {
						break
					}
				}
			}
			//req,err := http.NewRequest("POST","darth.chexiang.com/builder/_/get-build-status/",bytes.NewBuffer("{package_url:http://10.32.135.110/package/20160510/mvp-webapp-1.0-SNAPSHOT-201605101144.war}"))
		}
	}

	go SendEmail(history)

	return outcome
}

func checkAddrExist(addrs *[]string, addr string) bool {
	for _, ad := range *addrs {
		if ad == addr {
			return true
		}
	}
	return false
}

type RequestData struct {
	Packages []string `json:"packages"`
}

type PageurlData struct {
	PackageUrl string `json:"package_url"`
}
