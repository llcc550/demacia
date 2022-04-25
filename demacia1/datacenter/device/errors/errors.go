package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	DeviceNotExist   = httpx.NewCodeError(10601, "设备不存在")
	DeviceTitleExist = httpx.NewCodeError(10602, "设备名重复")
	DeviceSnExist    = httpx.NewCodeError(10603, "设备SN重复")
	GroupExist       = httpx.NewCodeError(10604, "设备组名称已存在")
)
