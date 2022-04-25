package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	OrganizationNotExist = httpx.NewCodeError(10401, "机构不存在")
	OrganizationExist    = httpx.NewCodeError(10402, "机构已存在")
	ManagerExist         = httpx.NewCodeError(10403, "管理员账号已存在")
)
