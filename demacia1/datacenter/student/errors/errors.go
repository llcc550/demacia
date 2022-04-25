package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	StudentExist         = httpx.NewCodeError(10701, "用户名存在")
	StudentNotExist      = httpx.NewCodeError(10702, "学生不存在")
	StudentIdNumberExist = httpx.NewCodeError(10703, "学生身份证已存在")
	StudentCardExist     = httpx.NewCodeError(10704, "卡号存在")
	StudentClassNotExist = httpx.NewCodeError(10705, "班级不存在")
)
