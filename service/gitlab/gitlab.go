package bdgitlab

import (
	"fmt"
	"go-jenkins/models/bd"
	"go-jenkins/models/ldap"

	"github.com/astaxie/beego"
	"github.com/xanzy/go-gitlab"
)

var pClient *gitlab.Client

func getClient() *gitlab.Client {
	if pClient != nil {
		return pClient
	}

	pClient := gitlab.NewClient(nil, beego.AppConfig.String("gitlabToken"))
	pClient.SetBaseURL(beego.AppConfig.String("gitlabApiUrl"))

	return pClient
}

func GetProjectById(pId int64) (*gitlab.Project, error) {
	client := getClient()
	gitProject, _, err := client.Projects.GetProject(int(pId))
	return gitProject, err
}

func RawProjects(projects *[]bdmodels.BdProject) error {
	client := getClient()
	for i, _ := range *projects {
		project := &(*projects)[i]
		gitProject, _, err := client.Projects.GetProject(int(project.ProjectId))
		if err != nil {
			return err
		}
		project.ProjectName = *gitProject.NameWithNamespace
	}
	return nil
}

func GetBranchByProjectId(pId string) ([]*gitlab.Branch, error) {
	client := getClient()
	branches, _, err := client.Branches.ListBranches(pId)
	return branches, err
}

func GetProjectsByUser(user *ldapmodels.LdapUser) []*gitlab.Project {
	client := gitlab.NewClient(nil, user.GitToken)
	client.SetBaseURL(beego.AppConfig.String("gitlabApiUrl"))

	var projects []*gitlab.Project
	for projectPage := 1; true; projectPage++ {
		pageListOptions := gitlab.ListOptions{
			Page:    projectPage,
			PerPage: 100,
		}
		projectOptions := gitlab.ListProjectsOptions{
			ListOptions: pageListOptions,
		}

		gitProjects, _, err := client.Projects.ListProjects(&projectOptions)

		if len(gitProjects) == 0 {
			break
		}

		if err != nil {
			fmt.Println("get git project error:", err.Error())
			break
		}

		projects = appendProjectArray(gitProjects, projects)

	}
	return projects
}

func appendProjectArray(array1, array2 []*gitlab.Project) []*gitlab.Project {
	//var array3 []*gitlab.Project
	var array3 = make([]*gitlab.Project, len(array1)+len(array2))
	if len(array1) != 0 {
		for i, project := range array1 {
			array3[i] = project
		}
	}
	if len(array2) != 0 {
		for i, project := range array2 {
			array3[len(array1)+i] = project
		}
	}
	return array3
}

func GetSingleBranch(pId int, branchName string) (*gitlab.Branch, error) {
	client := getClient()
	branch, _, err := client.Branches.GetBranch(pId, branchName)
	return branch, err
}

func GetAuthedUserForProject(pId int64) []*gitlab.ProjectMember {
	client := getClient()
	var members []*gitlab.ProjectMember
	for userPage := 1; true; userPage++ {
		pageListOptions := gitlab.ListOptions{
			Page:    userPage,
			PerPage: 100,
		}
		userOptions := gitlab.ListProjectMembersOptions{
			ListOptions: pageListOptions,
			Query:       "",
		}

		gitUsers, _, err := client.Projects.ListProjectMembers(int(pId), &userOptions)

		if len(gitUsers) == 0 {
			break
		}

		if err != nil {
			fmt.Println("get git auto user error:", err.Error())
			break
		}

		members = appendMemberArray(gitUsers, members)

	}

	group := GetGroupByProject(pId)
	if group != nil {
		fmt.Println("success get the group:", group.Name)
		gMembers, err := GetGroupMembers(group.ID)
		if err != nil {
			fmt.Println(err.Error())
			return members
		}
		members = appendGroupMandProjectm(gMembers, members)

	}

	return members
}

func GetGroupMembers(gId int) ([]*gitlab.GroupMember, error) {
	client := getClient()
	members, _, err := client.Groups.ListGroupMembers(gId)
	return members, err
}

func GetGroupByProject(pId int64) *gitlab.Group {
	project, err := GetProjectById(pId)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	groups := GetAllGroup()
	for _, group := range groups {
		if group.Name == *project.Namespace.Name {
			return group
		}
	}
	return nil
}

func GetAllGroup() []*gitlab.Group {
	client := getClient()
	var groups []*gitlab.Group
	for groupPage := 1; true; groupPage++ {
		pageListOptions := gitlab.ListOptions{
			Page:    groupPage,
			PerPage: 100,
		}
		groupOptions := gitlab.ListGroupsOptions{
			ListOptions: pageListOptions,
			Search:      "",
		}

		gitGroups, _, err := client.Groups.ListGroups(&groupOptions)

		if len(gitGroups) == 0 {
			break
		}

		if err != nil {
			fmt.Println("get git auto user error:", err.Error())
			break
		}

		groups = appendGroupArray(gitGroups, groups)

	}
	return groups
}

func appendGroupMandProjectm(array1 []*gitlab.GroupMember, array2 []*gitlab.ProjectMember) []*gitlab.ProjectMember {
	var array3 = make([]*gitlab.ProjectMember, len(array1)+len(array2))
	if len(array1) != 0 {
		for i, member := range array1 {
			tempMember := new(gitlab.ProjectMember)
			tempMember.Name = member.Name
			tempMember.AccessLevel = member.AccessLevel
			tempMember.CreatedAt = member.CreatedAt
			tempMember.Email = member.Email
			tempMember.ID = member.ID
			tempMember.State = member.State
			tempMember.Username = member.Username
			array3[i] = tempMember
		}
	}
	if len(array2) != 0 {
		for i, member := range array2 {
			array3[len(array1)+i] = member
		}
	}
	return array3
}

func appendGroupArray(array1, array2 []*gitlab.Group) []*gitlab.Group {
	var array3 = make([]*gitlab.Group, len(array1)+len(array2))
	if len(array1) != 0 {
		for i, member := range array1 {
			array3[i] = member
		}
	}
	if len(array2) != 0 {
		for i, member := range array2 {
			array3[len(array1)+i] = member
		}
	}
	return array3
}

func appendMemberArray(array1, array2 []*gitlab.ProjectMember) []*gitlab.ProjectMember {
	var array3 = make([]*gitlab.ProjectMember, len(array1)+len(array2))
	if len(array1) != 0 {
		for i, member := range array1 {
			array3[i] = member
		}
	}
	if len(array2) != 0 {
		for i, member := range array2 {
			array3[len(array1)+i] = member
		}
	}
	return array3
}
