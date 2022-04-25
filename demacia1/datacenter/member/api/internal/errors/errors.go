package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	NotExist           = httpx.NewCodeError(10501, "人员不存在")
	MobileFormatError  = httpx.NewCodeError(10502, "手机号格式错误")
	MobileExitError    = httpx.NewCodeError(10503, "手机号已存在")
	CardsExitError     = httpx.NewCodeError(10504, "卡号已存在")
	NotNormal          = httpx.NewCodeError(10505, "人员状态不可用")
	GetUserDetailError = httpx.NewCodeError(10506, "获取人员信息错误")
	Inoperable         = httpx.NewCodeError(10507, "不可操作")
	ParamsFail         = httpx.NewCodeError(10508, "参数错误")
	OrgNotExist        = httpx.NewCodeError(10509, "单位不存在")
)
