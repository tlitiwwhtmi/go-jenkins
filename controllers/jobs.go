package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go-jenkins/models/jenkins"

	"github.com/astaxie/beego"
	"github.com/bndr/gojenkins"
)

type JobsController struct {
	beego.Controller
}

//调用jenkins创建job，
func (c *JobsController) Create() {
	defer c.Ctx.WriteString("done")

	jenkins, err := gojenkins.CreateJenkins("http://jenkins.dds.com/", "liubin", "Liubin5%").Init()

	if err != nil {
		panic("Something Went Wrong")
	}

	job, err := jenkins.GetJob("test_jenkins_config")

	if err != nil {
		panic("Job does not exist")
	}

	config, err := job.GetConfig()
	configStruct := jobconfig.JobConfig{}

	err = xml.Unmarshal([]byte(config), &configStruct)

	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	jenkins_test, err := gojenkins.CreateJenkins("http://10.32.173.215:9090/", "admin", "admin").Init()
	output, err := xml.MarshalIndent(configStruct, "  ", "    ")

	configStr := string(output)
	jenkins_test.CreateJob(configStr, "hello")

	fmt.Println(string(output))

	for i := range c.Data {
		fmt.Println(i)
		fmt.Println(c.Data[i])
	}

	//fmt.Println(c.Ctx.Input.RequestBody)
	//c.Ctx.Input.RequestBody

	var s Result
	json.Unmarshal(c.Ctx.Input.RequestBody, &s)
	fmt.Println(string(c.Ctx.Input.RequestBody))
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}
	fmt.Println(string(b))

	fmt.Println(s.Aggs.Project.Terms.Field)
	fmt.Println(s)
}

//删除job

//执行job

type TermsS struct {
	Field string
	Name  string
}

type ProjectS struct {
	Terms TermsS
}

type AggsS struct {
	Project ProjectS
}

type Result struct {
	Aggs AggsS
}

//config
