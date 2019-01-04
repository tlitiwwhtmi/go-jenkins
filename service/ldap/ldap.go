package bdldap

import (
	"github.com/astaxie/beego"
	"github.com/axgle/mahonia"
	"github.com/nmcclain/ldap"
	"go-jenkins/errors"
	"go-jenkins/models/ldap"
)

var dec = mahonia.NewEncoder("ISO8859-1")

func getConnection() *ldap.Conn {
	con, err := ldap.Dial("tcp", beego.AppConfig.String("ldap_server")+":"+beego.AppConfig.String("ldap_port"))
	if err != nil {
		return nil
	}
	return con
}

func LdapLogin(account string, passWord string) (*ldapmodels.LdapUser, *bderrors.Bderror) {
	con := getConnection()
	clientCon := getConnection()
	if con == nil || clientCon == nil {
		return nil, bderrors.LdapError
	}
	defer con.Close()
	defer clientCon.Close()

	err := con.Bind(beego.AppConfig.String("ldap_user"), beego.AppConfig.String("ldap_pass"))
	if err != nil {
		return nil, bderrors.LdapError
	}

	searchRequest := &ldap.SearchRequest{
		beego.AppConfig.String("base_dn"),
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		"(sAMAccountName=" + account + ")",
		nil,
		nil}
	result, err := con.Search(searchRequest)
	if err != nil {
		return nil, bderrors.LdapError
	} else {
		if len(result.Entries) == 0 {
			return nil, bderrors.WrongAccount
		}
		entry := result.Entries[0]
		user := new(ldapmodels.LdapUser)
		for _, attr := range entry.Attributes {
			if attr.Name == "sAMAccountName" {
				user.Account = attr.Values[0]
			}
			if attr.Name == "givenName" {
				user.Name = dec.ConvertString(attr.Values[0])
			}
			if attr.Name == "cn" {
				//user.Cn = attr.Values[0]
				user.Cn = dec.ConvertString(attr.Values[0])
			}
			if attr.Name == "distinguishedName" {
				user.DistinguishedName = dec.ConvertString(attr.Values[0])
			}
		}

		err = clientCon.Bind(user.DistinguishedName, passWord)
		if err != nil {
			return nil, bderrors.WrongPass
		}
		return user, nil
	}

}
