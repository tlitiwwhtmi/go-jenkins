package bdjenkins

import (
	"encoding/xml"
	"fmt"
	"go-jenkins/models/bd"
	"go-jenkins/models/jenkins"

	"github.com/astaxie/beego"
	"github.com/bndr/gojenkins"
	"github.com/xanzy/go-gitlab"
	"go-jenkins/service/database"
	"go-jenkins/service/gitlab"
	"go-jenkins/utils"
	"strconv"
	"strings"
)

var labelIndex int

func getConnection() (*gojenkins.Jenkins, error) {
	con, err := gojenkins.CreateJenkins(beego.AppConfig.String("jenkins_url"), "admin", "chexiangadmin").Init()
	return con, err
}

func ReadJob() {
	jenkins, err := gojenkins.CreateJenkins("http://10.32.173.215:9090", "admin", "admin").Init()

	if err != nil {
		fmt.Println("Something Went Wrong")
	}

	job, err := jenkins.GetJob("dh_test1")

	if err != nil {
		fmt.Println("Job does not exist")
	}

	config, err := job.GetConfig()
	fmt.Println(config)
	configStruct := jobconfig.JobConfig{}

	_ = xml.Unmarshal([]byte(config), &configStruct)
}

func createJenkinsJob(bdjob *bdmodels.BdJob, history *bdmodels.BdBuildHistory, port int) (*gojenkins.Job, error) {
	//backIp := utils.GetLocalIp()
	//backIp := "10.47.18.19"
	//backIp := "build.dds.com"
	backIp := "10.47.12.116"

	_, projects, err := database.GetProjectsByProjectId(bdjob.ProjectId)
	if err != nil {
		return nil, err
	}
	bdProject := (*projects)[0]

	if labelIndex == 11 {
		labelIndex = 0
	}
	labelIndex++
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	config := new(jobconfig.JobConfig)
	//config.AssignedNode = "slave-" + strconv.Itoa(labelIndex)

	if bdProject.Language == "Java" {
		config.AssignedNode = "slave-1"
	}
	if bdProject.Language == "Android" {
		config.AssignedNode = "slave-3"
	}

	config.Disabled = "false"
	config.ConcurrentBuild = "false"
	buildShell := ""
	buildShell += "backIp=" + backIp + "\r\n"
	buildShell += "backPort=" + strconv.Itoa(port) + "\r\n"
	buildShell += "backId=" + fmt.Sprintf("%d", history.Id) + "\r\n"
	buildShell += "echo \"10.32.135.82 maven.dds.com\" >> /etc/hosts \r\n"
	buildShell += "echo \"10.47.12.218 git01.dds.com\" >> /etc/hosts \r\n"
	buildShell += "echo \"10.32.135.102 jira.dds.com\" >> /etc/hosts \r\n"
	buildShell += "echo \"10.32.135.82 docs.dds.com\" >> /etc/hosts \r\n"
	buildShell += "npm config set registry http://registry.npm.taobao.org/ \r\n"
	buildShell += "echo 'nameserver 10.47.12.60' > /etc/resolv.conf \r\n"
	buildShell += "ssh-keyscan -H docs.dds.com >> /root/.ssh/known_hosts \r\n"
	buildShell += "export TZ=Asia/Shanghai \r\n"
	gitproject, err := bdgitlab.GetProjectById(bdjob.ProjectId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	buildShell += "git clone " + *gitproject.SSHURLToRepo + "\r\n"
	buildShell += "cd $WORKSPACE/" + *gitproject.Path + "\r\n"
	//buildShell += "git checkout origin/" + bdjob.BranchName + "\r\n"
	buildShell += "git checkout " + history.Version + "\r\n"

	generateShell, err := generateBuildShell(bdProject, bdjob, gitproject)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	buildShell += generateShell

	config.Builders.Shell.Command = buildShell
	config.Publishers.Script.BuildSteps.Shell.Command = "curl http://" + backIp + ":" + strconv.Itoa(port) + "/jobtrriger?id=" + fmt.Sprintf("%d", history.Id)
	config.Publishers.Script.Plugin = "postbuildscript@0.17"
	config.Publishers.Script.MarkBuildUnstable = "false"
	config.Publishers.Script.ScriptOnlyIfFailure = "false"
	config.Publishers.Script.ScriptOnlyIfSuccess = "false"
	output, err := xml.Marshal(config)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	job, err := con.CreateJob(string(output), fmt.Sprintf("%d", bdjob.ProjectId)+"_"+bdjob.JobName)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return job, nil
}

func generateBuildShell(bdProject bdmodels.BdProject, bdjob *bdmodels.BdJob, gitproject *gitlab.Project) (string, error) {
	if bdProject.Language == "Java" {
		return generateJavaBuildShell(bdjob, gitproject)
	}
	if bdProject.Language == "Android" {
		return generateAndroidBuildShell(bdjob, gitproject)
	}
	return "", nil
}

func generateAndroidBuildShell(bdjob *bdmodels.BdJob, gitproject *gitlab.Project) (string, error) {
	buildShell := ""
	buildShell += bdjob.Shell
	return buildShell, nil
}

func generateJavaBuildShell(bdjob *bdmodels.BdJob, gitproject *gitlab.Project) (string, error) {
	buildShell := ""
	fileName := "pom.xml"
	if bdjob.Build == "1" { //build type standard
		if bdjob.PomPath == "" {
			if bdjob.Deploy == 1 {
				buildShell += "mvn -T 4 clean package -Dmaven.test.skip=true"
			} else if bdjob.Deploy == 2 {
				buildShell += "mvn -T 4 clean deploy -Dmaven.test.skip=true"
			} else {
				buildShell += "mvn clean compile org.codehaus.mojo:sonar-maven-plugin:3.0.2:sonar"
			}
		} else {
			pomPath := bdjob.PomPath
			index := strings.LastIndex(pomPath, "/")
			if index < 0 {
				if bdjob.Deploy == 1 {
					buildShell += "mvn -T 4 clean package -f " + pomPath + " -Dmaven.test.skip=true"
				} else if bdjob.Deploy == 2 {
					buildShell += "mvn -T 4 clean deploy -f " + pomPath + " -Dmaven.test.skip=true"
				} else {
					buildShell += "mvn clean compile -f " + pomPath + " org.codehaus.mojo:sonar-maven-plugin:3.0.2:sonar"
				}
			} else {
				rs := []rune(pomPath)
				end := len(rs)
				filePath := string(rs[0:index])
				if filePath == "." {
					filePath = ""
				}
				fileName = string(rs[index+1 : end])
				buildShell += "cd $WORKSPACE/" + *gitproject.Path + "/" + filePath + "\r\n"
				if bdjob.Deploy == 1 {
					buildShell += "mvn -T 4 clean package -f " + fileName + " -Dmaven.test.skip=true"
				} else if bdjob.Deploy == 2 {
					buildShell += "mvn -T 4 clean deploy -f " + fileName + " -Dmaven.test.skip=true"
				} else {
					buildShell += "mvn clean compile  -f " + fileName + " org.codehaus.mojo:sonar-maven-plugin:3.0.2:sonar"
				}
			}

		}
		if bdjob.Profile != "" {
			buildShell += " -P" + bdjob.Profile + " -X -U \r\n"
		} else {
			buildShell += " -X -U \r\n"
		}
	} else {
		buildShell += "cd $WORKSPACE/" + "\r\n"
		buildShell += bdjob.Shell
	}
	moreShell, err := utils.GetBuildShell()
	if err != nil {
		return buildShell, err
	}
	buildShell += "\r\n"
	if bdjob.Deploy == 1 {
		buildShell += moreShell
	}
	if bdjob.Build == "1" && bdjob.Deploy != 3 { //build type standard
		buildShell += "mvn build-helper:remove-project-artifact -f " + fileName + " \r\n"
	}
	return buildShell, nil
}

func RunJob(bdjob *bdmodels.BdJob, history *bdmodels.BdBuildHistory, port int) (*gojenkins.Job, error) {
	job, err := createJenkinsJob(bdjob, history, port)
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: creating JOB failed " + err.Error())
		fmt.Println("duanhao error")
		return nil, err
	}
	utils.GetLogger().Info("CREATING JENKINS JOB SUCCESS")

	con, err := getConnection()
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: get jenkins connection failed " + err.Error())
		fmt.Println(err.Error())
		return nil, err
	}
	isRunning, err := con.BuildJob(fmt.Sprintf("%d", bdjob.ProjectId) + "_" + bdjob.JobName)
	if err != nil {
		utils.GetLogger().Error("FAILED RUNNING JOB: starting jenkins job failed " + err.Error())

		fmt.Println(err.Error())
		return nil, err
	}

	utils.GetLogger().Info("SUCCESS RUNNING JOB")

	fmt.Println(isRunning)
	return job, nil
}

func GetBuildInfo(bdJob *bdmodels.BdJob) (*gojenkins.Build, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}
	job, err := con.GetJob(fmt.Sprintf("%d", bdJob.ProjectId) + "_" + bdJob.JobName)
	if err != nil {
		return nil, err
	}
	build, err := job.GetLastBuild()
	if err != nil {
		return nil, err
	}
	return build, nil
}

func DeleteJenkinsJob(bdJob *bdmodels.BdJob) (bool, error) {
	con, err := getConnection()
	if err != nil {
		return false, err
	}
	job, err := con.GetJob(fmt.Sprintf("%d", bdJob.ProjectId) + "_" + bdJob.JobName)
	if err != nil {
		return false, err
	}
	isDelete, err := job.Delete()
	return isDelete, err
}

func GetJobStatus(bdJob *bdmodels.BdJob) (bool, string, *string) {
	con, err := getConnection()
	var log string
	if err != nil {
		return false, "not running", nil
	}
	job, err := con.GetJob(fmt.Sprintf("%d", bdJob.ProjectId) + "_" + bdJob.JobName)
	if err != nil {
		return false, "not running", nil
	}
	build, err := job.GetLastBuild()
	if err != nil || build == nil {
		log = "connecting..."
		return true, "waiting", &log
	}
	log = build.GetConsoleOutput()
	return true, "building", &log
}

func StopJenkinsJob(bdJob *bdmodels.BdJob) error {
	con, err := getConnection()
	if err != nil {
		return err
	}
	job, err := con.GetJob(fmt.Sprintf("%d", bdJob.ProjectId) + "_" + bdJob.JobName)
	if err != nil {
		return err
	}
	build, _ := job.GetLastBuild()

	if build != nil {
		build.Stop()
		//DeleteJenkinsJob(bdJob)
		return nil
	}
	queue, err := con.GetQueue()
	if err != nil {
		return err
	}
	tasks := queue.GetTasksForJob(fmt.Sprintf("%d", bdJob.ProjectId) + "_" + bdJob.JobName)
	for _, task := range tasks {
		task.Cancel()
	}
	DeleteJenkinsJob(bdJob)
	return nil
}
