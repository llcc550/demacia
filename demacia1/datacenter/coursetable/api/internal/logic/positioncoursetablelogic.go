package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/errors"
	"time"

	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type PositionCourseTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPositionCourseTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) PositionCourseTableLogic {
	return PositionCourseTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PositionCourseTableLogic) PositionCourseTable(req types.PositionIdReq) (*types.PositionCourseTableReply, error) {

	resp := &types.PositionCourseTableReply{PositionCourseTableInfo: []*types.OrgCourseInfo{}}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.InvalidParam
	}

	if req.PositionId == 0 {
		return resp, errlist.InvalidParam
	}
	deploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(oid)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, errors.MustConfig
		}
	}

	courseTables, err := l.svcCtx.CourseTableModel.SelectByPositionId(req.PositionId)
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
			resp.PositionCourseTableInfo = append(resp.PositionCourseTableInfo, &types.OrgCourseInfo{
				Weekday:    deploy.Weekday,
				CourseSort: deploy.CourseSort,
				StartTime:  startTime.Format("15:04:05"),
				EndTime:    endTime.Format("15:04:05"),
				CourseNote: deploy.Note,
			})
		}
	}

	for _, table := range courseTables {
		for _, info := range resp.PositionCourseTableInfo {
			if info.Weekday == table.WeekDay && info.CourseSort == table.CourseSort {
				info.SubjectName = table.SubjectName
				info.TeacherName = table.TeacherName
				info.ClassName = table.ClassName
			}
		}
	}

	return resp, nil
}
