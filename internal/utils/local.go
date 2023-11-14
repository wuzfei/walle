package utils

import (
	"os"
	"os/user"
)

var CurrentUser *user.User
var CurrentHostname string
var CurrentUserHost string

func init() {
	var err error
	CurrentUser, err = user.Current()
	if err != nil {
		panic("获取当前用户错误:" + err.Error())
	}
	CurrentHostname, err = os.Hostname()
	if err != nil {
		panic("获取当前主机名错误:" + err.Error())
	}
	CurrentUserHost = CurrentUser.Username + "@" + CurrentHostname
}
