package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/model"
	positionErr "demacia/datacenter/position/errors"
	"demacia/datacenter/position/rpc/position"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/rest/httpx"
	"time"
)

type CourseTableAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTableAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseTableAddLogic {
	return CourseTableAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTableAddLogic) CourseTableAdd(req types.CourseTableAddReq) (*types.BoolReply, error) {

	if req.CourseSort == 0 || req.ClassId == 0 || req.Weekday == 0 || req.SubjectId == 0 || req.PositionId == 0 {
		return &types.BoolReply{Success: false}, errlist.InvalidParam
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &types.BoolReply{Success: false}, errlist.NoAuth
	}

	var course model.CourseTable

	deploy, err := l.svcCtx.CourseTableDeployModel.SelectByOrgIdAndWeekdayAndSort(oid, req.Weekday, req.CourseSort)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.BoolReply{Success: false}, errors.MustConfig
		}
		return &types.BoolReply{Success: false}, errlist.Unknown
	}
	if deploy.CourseFlag != 1 {
		return &types.BoolReply{Success: false}, errors.CourseAddErr
	}
	startTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.StartTime)
	endTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.EndTime)
	if startTime.IsZero() || endTime.IsZero() {
		return &types.BoolReply{Success: false}, errors.InvalidTime
	}
	course.StartTime = startTime.Format("15:04:05")
	course.EndTime = endTime.Format("15:04:05")
	course.CourseSort = deploy.CourseSort
	course.WeekDay = deploy.Weekday
	course.OrganizationId = oid
	errState := &httpx.CodeError{}
	err = mr.Finish(func() error {
		classInfo, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
		if err != nil {
			l.Logger.Errorf("call classRpc error:%s", err.Error())
			errState = errlist.ClassNotFound
			return err
		}
		course.ClassId = classInfo.Id
		course.ClassName = classInfo.FullName
		return nil
	}, func() error {
		subjectInfo, err := l.svcCtx.TeachModel.FindTeacherByClassIdAndSubjectId(req.ClassId, req.SubjectId)
		if err != nil {
			errState = errors.SubjectSetErr
			return err
		}
		course.SubjectId = subjectInfo.SubjectId
		course.SubjectName = subjectInfo.SubjectTitle
		course.TeacherId = subjectInfo.MemberId
		course.TeacherName = subjectInfo.TrueName
		return nil
	}, func() error {
		positionInfo, err := l.svcCtx.PositionRpc.FindById(l.ctx, &position.PositionIdReq{PositionId: req.PositionId})
		if err != nil {
			l.Logger.Errorf("call positionRpc error:%s", err.Error())
			errState = positionErr.PositionNotExist
			return err
		}
		course.PositionId = positionInfo.Id
		course.PositionName = positionInfo.PositionName
		return nil
	})
	if err != nil {
		return &types.BoolReply{Success: false}, errState
	}

	courseId, err := l.svcCtx.CourseTableModel.SelectExistByClassIdAndWeekdayAndSort(oid, req.ClassId, req.Weekday, req.CourseSort)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("select course exist err:%s", err.Error())
		return &types.BoolReply{Success: false}, errlist.Unknown
	}
	if courseId == 0 {
		if err := l.svcCtx.CourseTableModel.InsertCourseTable(&course); err != nil {
			l.Logger.Errorf("insert courseTable error:%s", err.Error())
			return &types.BoolReply{Success: false}, errlist.Unknown
		}
	} else {
		course.Id = courseId
		if err := l.svcCtx.CourseTableModel.UpdateCourseTable(&course); err != nil {
			l.Logger.Errorf("update courseTable error:%s", err.Error())
			return &types.BoolReply{Success: false}, errlist.Unknown
		}
	}

	return &types.BoolReply{Success: true}, nil
}
