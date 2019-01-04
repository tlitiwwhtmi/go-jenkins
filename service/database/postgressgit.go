package database

import (
	"database/sql"
	"fmt"
	"go-jenkins/models/ldap"

	_ "github.com/lib/pq"
	"go-jenkins/models/bd"
	"strconv"
)

type PostGresConn struct {
	url string
}

func NewPostGres(urlstr string) *PostGresConn {
	pqConn := PostGresConn{
		url: urlstr,
	}
	return &pqConn
}

func getPQConnection(urlstr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", urlstr)
	return db, err
}

func (pqConn *PostGresConn) GetUserByAccount(account string) (*ldapmodels.LdapUser, error) {
	db, err := getPQConnection(pqConn.url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id,authentication_token from users where username = '" + account + "'")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := new(ldapmodels.LdapUser)
		if err := rows.Scan(&user.GitId, &user.GitToken); err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

func (pqConn *PostGresConn) GetGitProjectDetail(projects *[]bdmodels.BdProject) (*[]*bdmodels.GitProject, error) {
	var gitProjects []*bdmodels.GitProject
	db, err := getPQConnection(pqConn.url)
	if err != nil {
		return nil, err
	}

	idstr := ""
	sql := ""
	if len(*projects) == 0 {
		return nil, err
	}
	if len(*projects) > 1 {
		idstr += "("
		for _, project := range *projects {
			idstr += strconv.Itoa(int(project.ProjectId))
			idstr += ","
		}
		idstr = idstr[0 : len(idstr)-1]
		idstr += ")"
		sql = "SELECT a.id, a.name, b.path as namespace FROM projects as a join namespaces as b on a.namespace_id = b.id where a.id in " + idstr
	} else {
		idstr = strconv.Itoa(int((*projects)[0].ProjectId))
		sql = "SELECT a.id, a.name, b.path as namespace FROM projects as a join namespaces as b on a.namespace_id = b.id where a.id = " + idstr
	}
	defer db.Close()
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		gitPro := new(bdmodels.GitProject)
		if err := rows.Scan(&gitPro.Id, &gitPro.Name, &gitPro.Namespace); err != nil {
			return nil, err
		}
		gitProjects = append(gitProjects, gitPro)
	}
	return &gitProjects, nil
}
