package util

import "go-zentao-task/pkg/config"

func GetRobotUrl() string {
	//获取机器人webhook地址
	robotUrl := config.Get("wxRobotUrl")
	return robotUrl
}
