package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	InvalidTime            = httpx.NewCodeError(10701, "错误的上课时间")
	CourseTimeErr          = httpx.NewCodeError(10702, "上/下课时间冲突")
	MustConfig             = httpx.NewCodeError(10703, "还未进行节次配置")
	CourseAddErr           = httpx.NewCodeError(10704, "错误的课程设置")
	InvalidMorningSelfTime = httpx.NewCodeError(10705, "早自习时间与上课时间冲突")
	InvalidNightSelfTime   = httpx.NewCodeError(10706, "晚自习时间与上课时间冲突")
	CourseNotFoundErr      = httpx.NewCodeError(10707, "未查询到指定课程")
	DeployNotFoundErr      = httpx.NewCodeError(10708, "未查询到课程配置")
	ClassNotFoundErr       = httpx.NewCodeError(10709, "未查询到班级信息")
	CourseConflictErr      = httpx.NewCodeError(10710, "教师课程冲突")
	SubjectSetErr          = httpx.NewCodeError(10711, "获取学科信息失败")
)
