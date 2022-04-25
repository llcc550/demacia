package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type TeacherCourseTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTeacherCourseTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) TeacherCourseTableLogic {
	return TeacherCourseTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TeacherCourseTableLogic) TeacherCourseTable(req types.MemberIdReq) (types.TeachersCourseTableReply, error) {

	resp := types.TeachersCourseTableReply{TeacherCourseTableInfo: types.TeacherCourseTableInfo{}}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.InvalidParam
	}

	if req.MemberId == 0 {
		return resp, errlist.InvalidParam
	}

	deploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(oid)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, errors.MustConfig
		}
	}

	courseTables, err := l.svcCtx.CourseTableModel.SelectByMemberId(req.MemberId)
	if err != nil {
		return resp, errlist.Unknown
	}

	if len(courseTables) == 0 {
		return resp, errors.CourseNotFoundErr
	}

	for _, deploy := range deploys {
		if deploy.CourseFlag == 1 {
			startTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.StartTime)
			endTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.EndTime)
			resp.TeacherCourseTableInfo.Classes = append(resp.TeacherCourseTableInfo.Classes, &types.OrgCourseInfo{
				Weekday:    deploy.Weekday,
				CourseSort: deploy.CourseSort,
				StartTime:  startTime.Format("15:04:05"),
				EndTime:    endTime.Format("15:04:05"),
				CourseNote: deploy.Note,
			})
		}
	}

	for _, table := range courseTables {
		for _, r := range resp.TeacherCourseTableInfo.Classes {
			if r.Weekday == table.WeekDay && r.CourseSort == table.CourseSort {
				r.ClassName = table.ClassName
				r.PositionName = table.PositionName
				r.SubjectName = table.SubjectName
			}
		}
	}

	return resp, nil
}
