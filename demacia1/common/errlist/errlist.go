package errlist

import "gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"

var (
	Unknown      = httpx.NewCodeError(10000, "服务暂不可用，请稍后再试")
	ImportSuc    = httpx.NewCodeError(10002, "%s导入成功")
	ImportErr    = httpx.NewCodeError(10003, "%s导入失败。%s")
	Excel        = httpx.NewCodeError(10004, "excel文件错误")
	InvalidParam = httpx.NewCodeError(10005, "传入参数有误")
	NoAuth       = httpx.NewCodeError(10006, "授权信息不合法")
	IdNumberErr  = httpx.NewCodeError(10007, "身份证格式有误")
	MobileErr    = httpx.NewCodeError(10008, "手机号格式错误")

	WebsocketNodeOffline = httpx.NewCodeError(10101, "暂无可用节点")

	CommonAreaIdErr = httpx.NewCodeError(10201, "区域参数错误")

	AuthLoginFail = httpx.NewCodeError(10301, "登录失败，账号不存在或密码错误")

	MemberNotExist = httpx.NewCodeError(10501, "人员不存在")

	DeviceNotExist   = httpx.NewCodeError(10601, "设备不存在")
	DeviceTitleExist = httpx.NewCodeError(10602, "设备名重复")
	DeviceSnExist    = httpx.NewCodeError(10603, "设备SN重复")
	GroupExist       = httpx.NewCodeError(10604, "设备组名称已存在")

	InvalidMobilesErr = httpx.NewCodeError(10601, "没有有效的手机号码")
	OutOfQuotaErr     = httpx.NewCodeError(10602, "超出每日短信数量限制")
	SendSmsErr        = httpx.NewCodeError(10603, "发送短信失败")

	DepartmentNotExit = httpx.NewCodeError(10601, "部门不存在")
	DepartmentExit    = httpx.NewCodeError(10601, "部门已存在")

	ClassNameExist    = httpx.NewCodeError(10901, "班级名称已存在")
	StageGradeErr     = httpx.NewCodeError(10902, "学段年级数据异常")
	StageNameExist    = httpx.NewCodeError(10903, "学段名称已存在")
	ClassNotFound     = httpx.NewCodeError(10904, "班级不存在")
	TeacherExistTeach = httpx.NewCodeError(10905, "教师已存在任课")
	StageYearTooLong  = httpx.NewCodeError(10905, "学段年制最大20")

	SubjectNotFound = httpx.NewCodeError(11001, "学科不存在")
	SubjectExist    = httpx.NewCodeError(11002, "学科已存在")

	StudentIdNumberExist = httpx.NewCodeError(10703, "学生身份证已存在")
	MobileExist          = httpx.NewCodeError(10704, "手机号已存在")

	EventCategoryExist    = httpx.NewCodeError(10701, "栏目名称已存在")
	EventCategoryNotExist = httpx.NewCodeError(10702, "栏目不存在")
	EventNotExist         = httpx.NewCodeError(10703, "通知不存在")
	EventPositionErr      = httpx.NewCodeError(10704, "范围不能为空")
	EventTimeErr          = httpx.NewCodeError(10705, "开始/结束时间错误")
)
