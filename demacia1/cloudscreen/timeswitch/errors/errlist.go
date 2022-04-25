package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	SwitchTimeErr      = httpx.NewCodeError(11001, "错误的开关机时间")
	SwitchTimeRangeErr = httpx.NewCodeError(11002, "开关机时间段不可重叠")
	SwitchDateErr      = httpx.NewCodeError(11003, "错误的日期参数")
)
