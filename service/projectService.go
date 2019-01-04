package service

import (
	"go-jenkins/models/bd"
	"go-jenkins/service/database"
	"go-jenkins/service/gitlab"
)

func GetAllProjects() (*[]bdmodels.BdProject, error) {
	projects, err := database.GetProjects()
	if err != nil {
		return nil, err
	}
	bdgitlab.RawProjects(projects)
	return projects, nil
}
