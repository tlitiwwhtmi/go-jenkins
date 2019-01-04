package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"go-jenkins/models/ldap"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

func GetUserFromSession(ctx *context.Context) *ldapmodels.LdapUser {
	body, err := Base64Decode(ctx.GetCookie("user"))
	/*fmt.Println(string(body))
	if user := ctx.Input.Session("loginUser"); user != nil {
		if value, ok := user.(*ldapmodels.LdapUser); ok {
			return value
		}
	}
	return nil*/
	if err != nil {
		return nil
	}
	user := new(ldapmodels.LdapUser)
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil
	}
	return user
}

func GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
				//fmt.Println(ipnet.IP.String())
			}
		}
	}
	return ""
}

func GetBuildShell() (string, error) {
	fileBytes, err := ioutil.ReadFile("./conf/build.sh")
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

func Covert2Email(account string) string {
	return account + "@chexiang.com"
}

func CoverDurationTime(mSeconds int64) string {
	var timeStr string
	hour := 0
	min := 0
	sec := 0
	sec = int(mSeconds) / 1000
	if sec/3600 > 0 {
		hour = sec / 3600
		sec = sec % 3600
	}
	if sec/60 > 0 {
		min = sec / 60
		sec = sec % 60
	}
	if hour != 0 {
		timeStr += strconv.Itoa(hour) + "h"
	}
	if min != 0 {
		timeStr += strconv.Itoa(min) + "m"
	}
	timeStr += strconv.Itoa(sec) + "s"
	return timeStr
}

func OrgniseUniqName() string {
	t := time.Now().UnixNano()
	return strconv.FormatInt(t, 10)
}

func CompressStr(log string) (bytes.Buffer, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(log)); err != nil {
		return b, err
	}
	if err := gz.Flush(); err != nil {
		return b, err
	}
	if err := gz.Close(); err != nil {
		return b, err
	}
	return b, nil
}

func CoverArray2str(addrs *[]string) string {
	addrStr := ""
	for _, add := range *addrs {
		addrStr += add + ";"
	}
	return addrStr
}
