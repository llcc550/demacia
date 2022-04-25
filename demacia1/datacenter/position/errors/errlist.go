package errors

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	PositionNotExist      = httpx.NewCodeError(10801, "地点不存在")
	PositionBindClassErr  = httpx.NewCodeError(10802, "该班级已被绑定")
	PositionBindDeviceErr = httpx.NewCodeError(10803, "云屏绑定冲突")
	PositionExistErr      = httpx.NewCodeError(10804, "地点已存在")
	PositionGenerateErr   = httpx.NewCodeError(10805, "%s生成失败")
	PositionGenerateSuc   = httpx.NewCodeError(10806, "%s生成成功")
	PositionEditDeviceErr = httpx.NewCodeError(10807, "云屏绑定失败")
)
