package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	SignDateErr        = httpx.NewCodeError(10901, "错误的打卡日期")
	SignTimeErr        = httpx.NewCodeError(10902, "错误的打卡时间")
	SignRecordNotFound = httpx.NewCodeError(10903, "未获取到打卡信息")
	SignRecordExist    = httpx.NewCodeError(10904, "已打卡，无需重复打卡")
	NoEnable           = httpx.NewCodeError(10905, "未开启打卡功能")
)
