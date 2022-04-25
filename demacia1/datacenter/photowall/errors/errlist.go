package errors

import (
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
)

var (
	PhotoFolderNameExist = httpx.NewCodeError(10701, "相册名称已存在")
)
