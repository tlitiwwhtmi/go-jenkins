package bderrors

import (
	"errors"
)

type Bderror struct {
	Code int
	Err  error
}

var LdapError = &Bderror{1, errors.New("connect to ldap server error")}
var WrongAccount = &Bderror{2, errors.New("wrong account")}
var WrongPass = &Bderror{3, errors.New("Wrong password")}
