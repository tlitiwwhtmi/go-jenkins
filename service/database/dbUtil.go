package database

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"go-jenkins/models/bd"
	"strconv"
)

func GetProjects() (*[]bdmodels.BdProject, error) {
	o := getConnection()
	var bdProjects []bdmodels.BdProject
	count, err := (*o).Raw("select * from bd_project").QueryRows(&bdProjects)
	fmt.Println(count)
	return &bdProjects, err
}

func AddProject(project *bdmodels.BdProject) (int64, error) {
	o := getConnection()
	id, err := (*o).Insert(project)
	if err != nil {
		fmt.Println(err.Error())
	}
	return id, err
}

func AddPackage(packege *bdmodels.BdPackage) (int64, error) {
	o := getConnection()
	id, err := (*o).Insert(packege)
	if err != nil {
		fmt.Println(err.Error())
	}
	return id, err
}

func AddJob(job *bdmodels.BdJob) (int64, error) {
	o := getConnection()
	id, err := (*o).Insert(job)
	if err != nil {
		fmt.Println(err.Error())
	}
	return id, err
}

func UpdateJob(job *bdmodels.BdJob) (int64, error) {
	o := getConnection()
	id, err := (*o).Update(job)
	if err != nil {
		fmt.Println(err.Error())
	}
	return id, err
}

func AddHistory(history *bdmodels.BdBuildHistory) (int64, error) {
	o := getConnection()
	id, err := (*o).Insert(history)
	if err != nil {
		fmt.Println(err.Error())
	}
	return id, err
}

func GetProjectByAccountAndPid(project *bdmodels.BdProject) (*bdmodels.BdProject, error) {
	o := getConnection()
	pro := new(bdmodels.BdProject)
	err := (*o).Raw("select * from bd_project where project_id = ? and user_account = ?", project.ProjectId, project.UserAccount).QueryRow(pro)
	return pro, err
}

func DeleteProject(project *bdmodels.BdProject) error {
	o := getConnection()
	_, err := (*o).Delete(project)
	return err
}

func GetProjectsByAccount(account string) (int64, *[]bdmodels.BdProject, error) {
	o := getConnection()
	var bdProjects []bdmodels.BdProject
	count, err := (*o).Raw("select * from bd_project where user_account = ?", account).QueryRows(&bdProjects)
	return count, &bdProjects, err
}

func GetProjectsByProjectId(projectId int64) (int64, *[]bdmodels.BdProject, error) {
	o := getConnection()
	var bdProjects []bdmodels.BdProject
	count, err := (*o).Raw("select * from bd_project where project_id = ?", projectId).QueryRows(&bdProjects)
	return count, &bdProjects, err
}

func SearchPacsByPacName(projectId int64, pageIndex, perPage int, keyword string) (*[]bdmodels.BdPackage, error) {
	o := getConnection()
	var bdPacs []bdmodels.BdPackage
	_, err := (*o).Raw("select * from bd_package where history_id in (select id from bd_build_history where job_id in (select id from bd_job where project_id = ?)) and name like %?% order by id desc limit ?,?", projectId, keyword, (pageIndex-1)*perPage, perPage).QueryRows(&bdPacs)
	if err != nil {
		return nil, err
	}
	return &bdPacs, nil
}

func GetPacsByProject(projectId, cId int64, pageIndex, perPage int, branch string) (*[]bdmodels.BdPackage, error) {
	o := getConnection()
	var bdPacs []bdmodels.BdPackage
	var err error
	if branch == "" {
		if cId != 0 {
			_, err = (*o).Raw("select * from bd_package where history_id in (select id from bd_build_history where job_id in (select id from bd_job where project_id = ?)) and id < ? order by id desc limit ?,?", projectId, cId, (pageIndex-1)*perPage, perPage).QueryRows(&bdPacs)
		} else {
			_, err = (*o).Raw("select * from bd_package where history_id in (select id from bd_build_history where job_id in (select id from bd_job where project_id = ?)) order by id desc limit ?,?", projectId, (pageIndex-1)*perPage, perPage).QueryRows(&bdPacs)
		}

	} else {
		if cId != 0 {
			_, err = (*o).Raw("select * from bd_package where history_id in (select id from bd_build_history where job_id in (select id from bd_job where project_id = ? and branch_name = ?)) and id < ? order by id desc limit ?,?", projectId, branch, cId, (pageIndex-1)*perPage, perPage).QueryRows(&bdPacs)
		} else {
			_, err = (*o).Raw("select * from bd_package where history_id in (select id from bd_build_history where job_id in (select id from bd_job where project_id = ? and branch_name = ?)) order by id desc limit ?,?", projectId, branch, (pageIndex-1)*perPage, perPage).QueryRows(&bdPacs)
		}
	}
	if err != nil {
		return nil, err
	}
	return &bdPacs, nil
}

func GetJobsByProjectId(projectId int64, pageIndex, perPage int, branch string, keyword string) (int64, *[]bdmodels.BdJob, error) {
	o := getConnection()
	var bdJobs []bdmodels.BdJob
	var totalJobs []bdmodels.BdJob
	var total int64
	var err error
	if branch == "" {
		total, err = (*o).Raw("select * from bd_job where project_id = ? and job_name like '%"+keyword+"%'", projectId).QueryRows(&totalJobs)
		_, err = (*o).Raw("select * from bd_job where project_id = ? and job_name like '%"+keyword+"%' limit ?,?", projectId, (pageIndex-1)*perPage, perPage).QueryRows(&bdJobs)
	} else {
		total, err = (*o).Raw("select * from bd_job where project_id = ? and branch_name = ? and job_name like '%"+keyword+"%'", projectId, branch).QueryRows(&totalJobs)
		_, err = (*o).Raw("select * from bd_job where project_id = ? and branch_name = ? and job_name like '%"+keyword+"%' limit ?,?", projectId, branch, (pageIndex-1)*perPage, perPage).QueryRows(&bdJobs)
	}

	return total, &bdJobs, err
}

func GetHistoryByJob(jobId int64) (int64, *[]bdmodels.BdBuildHistory, error) {
	o := getConnection()
	var bdHistorys []bdmodels.BdBuildHistory
	count, err := (*o).Raw("select * from bd_build_history where job_id = ? order by id desc", jobId).QueryRows(&bdHistorys)
	return count, &bdHistorys, err
}

func GetLatestHistory(jobId int64) (*bdmodels.BdBuildHistory, error) {
	count, bdHistorys, err := GetHistoryByJob(jobId)
	if err != nil || count == 0 {
		return nil, err
	}
	return &(*bdHistorys)[0], nil
	return nil, nil
}

func GetJobById(jobId int64) (*bdmodels.BdJob, error) {
	o := getConnection()
	var job bdmodels.BdJob
	err := (*o).Raw("select * from bd_job where id = ?", jobId).QueryRow(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func RemoveJobById(jobId int64) error {
	o := getConnection()
	job, err := GetJobById(jobId)
	if err != nil {
		return err
	}
	_, err = (*o).Delete(job)
	return err
}

func RemoveHistoryByJob(jobId int64) error {
	o := getConnection()
	_, historys, err := GetHistoryByJob(jobId)
	if err != nil {
		return err
	}
	for _, history := range *historys {
		_, err = (*o).Delete(&history)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetHistoryById(hId int64) (*bdmodels.BdBuildHistory, error) {
	o := getConnection()
	var bdHistory bdmodels.BdBuildHistory
	err := (*o).Raw("select * from bd_build_history where id = ?", hId).QueryRow(&bdHistory)
	if err != nil {
		return nil, err
	}
	return &bdHistory, nil
}

func UpdateHistory(history *bdmodels.BdBuildHistory) (int64, error) {
	o := getConnection()
	hId, err := (*o).Update(history)
	return hId, err
}

func GetJobByIds(ids string) (*[]bdmodels.BdJob, error) {
	o := getConnection()
	var jobs []bdmodels.BdJob
	_, err := (*o).Raw("select * from bd_job where id in (" + ids + ")").QueryRows(&jobs)
	return &jobs, err
}

func GetJobByPIdAndBranch(id int64, bName string) (*[]bdmodels.BdJob, error) {
	o := getConnection()
	var jobs []bdmodels.BdJob
	_, err := (*o).Raw("select * from bd_job where project_id = ? and branch_name = ?", id, bName).QueryRows(&jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func ListAutoBuildJob() (*[]bdmodels.BdJob, error) {
	o := getConnection()
	var jobs []bdmodels.BdJob
	_, err := (*o).Raw("select * from bd_job where auto_build = ?", 1).QueryRows(&jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func ListHistories(jobId, cId int64, pageIndex, perPage int) (*[]bdmodels.BdBuildHistory, error) {
	o := getConnection()
	var bdHistorys []bdmodels.BdBuildHistory
	var err error
	if cId != 0 {
		_, err = (*o).Raw("select id,job_id,start_time,duration,status,commit_id,build_executor,version,message,commit_author from bd_build_history where job_id = ? and id < ? order by id desc limit ?,? ", jobId, cId, (pageIndex-1)*perPage, perPage).QueryRows(&bdHistorys)
	} else {
		_, err = (*o).Raw("select id,job_id,start_time,duration,status,commit_id,build_executor,version,message,commit_author from bd_build_history where job_id = ? order by id desc limit ?,? ", jobId, (pageIndex-1)*perPage, perPage).QueryRows(&bdHistorys)
	}
	if err != nil {
		return nil, err
	}
	return &bdHistorys, err
}

func GetHistoriesByIdsnolog(ids []int64) (*[]bdmodels.BdBuildHistory, error) {
	idStr := ""
	for _, id := range ids {
		idStr += strconv.Itoa(int(id)) + ","
	}
	rs := []rune(idStr)
	end := len(rs)
	tranferIds := string(rs[0 : end-1])
	o := getConnection()
	var bdHistorys []bdmodels.BdBuildHistory
	_, err := (*o).Raw("select id,job_id,start_time,duration,status,commit_id,build_executor,version,message,commit_author from bd_build_history where id in (" + tranferIds + ") order by id desc").QueryRows(&bdHistorys)
	return &bdHistorys, err
}

func GetJobCount() (int64, error) {
	o := getConnection()
	var jobs []bdmodels.BdJob
	count, err := (*o).Raw("select * from bd_job").QueryRows(&jobs)
	return count, err
}

func GetPacCount() (int64, error) {
	o := getConnection()
	var pacs []bdmodels.BdPackage
	count, err := (*o).Raw("select * from bd_package").QueryRows(&pacs)
	return count, err
}

func GetPacListByHistory(id int64) (*[]bdmodels.BdPackage, error) {
	o := getConnection()
	var pacs []bdmodels.BdPackage
	_, err := (*o).Raw("select * from bd_package where history_id = ?", id).QueryRows(&pacs)
	return &pacs, err
}

func GetAllHistories() (*[]bdmodels.BdBuildHistory, error) {
	o := getConnection()
	var bdHistorys []bdmodels.BdBuildHistory
	_, err := (*o).Raw("select * from bd_build_history").QueryRows(&bdHistorys)
	return &bdHistorys, err
}

func GetBuildsPerday() (*[]bdmodels.BuildsDay, error) {
	o := getConnection()
	var buildsDays []bdmodels.BuildsDay
	_, err := (*o).Raw("select count(*) as total,DATE_FORMAT(start_time,'%Y-%m-%d') as start_date from bd_build_history group by start_date order by start_date asc").QueryRows(&buildsDays)
	return &buildsDays, err
}

func getHistoryUsersCount() int {
	count := 0
	o := getConnection()
	var maps []orm.Params
	num, err := (*o).Raw("select count(distinct(build_executor)) as users from bd_build_history").Values(&maps)
	if err == nil && num > 0 {
		count, _ = strconv.Atoi(maps[0]["users"].(string))

	}
	return count
}

func getJobUsersCount() int {
	count := 0
	o := getConnection()
	var maps []orm.Params
	num, err := (*o).Raw("select count(*) as users from bd_job where creator_account not in (select distinct(build_executor) from bd_build_history)").Values(&maps)
	if err == nil && num > 0 {
		count, _ = strconv.Atoi(maps[0]["users"].(string))

	}
	return count
}

func GetUsersCount() int {
	return getHistoryUsersCount() + getJobUsersCount()
}
