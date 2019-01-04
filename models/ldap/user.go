package ldapmodels

type LdapUser struct {
	Account           string `json:"account"`
	Name              string `json:"name"`
	Cn                string `json:"-"`
	GitToken          string `json:"gitToken"`
	GitId             int    `json:"gitId"`
	DistinguishedName string `json:"-"`
}
