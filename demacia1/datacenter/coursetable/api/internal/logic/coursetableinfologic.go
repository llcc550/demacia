package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/position/rpc/position"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"time"
)

type CourseTableInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTableInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseTableInfoLogic {
	return CourseTableInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTableInfoLogic) CourseTableInfo(req types.CourseTableInfoReq) (*types.CourseTableInfoReply, error) {

	resp := &types.CourseTableInfoReply{CourseTable: []*types.CourseTableInfo{}}

	if req.PositionId == 0 {
		return resp, errlist.InvalidParam
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.InvalidParam
	}

	deploy, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(oid)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, errors.MustConfig
		}
		return resp, errlist.Unknown
	}

	positionInfo, err := l.svcCtx.PositionRpc.FindById(l.ctx, &position.PositionIdReq{PositionId: req.PositionId})
	if err != nil {
		return resp, errlist.Unknown
	}

	courseTables, err := l.svcCtx.CourseTableModel.SelectByPositionId(positionInfo.Id)
	if err != nil && err != sql.ErrNoRows {
		return resp, errlist.Unknown
	}
	resp.ClassId = positionInfo.ClassId
	resp.PositionName = positionInfo.PositionName
	resp.ClassName = positionInfo.ClassName
	for i := 1; i < 8; i++ {
		resp.CourseTable = append(resp.CourseTable, &types.CourseTableInfo{
			Weekday:    int8(i),
			CourseInfo: []*types.CourseInfo{},
		})
	}

	for _, info := range resp.CourseTable {
		for _, tableDeploy := range deploy {
			if info.Weekday == tableDeploy.Weekday {
				endTime, _ := time.Parse("2006-01-02T15:04:05Z", tableDeploy.EndTime)
				startTime, _ := time.Parse("2006-01-02T15:04:05Z", tableDeploy.StartTime)
				info.CourseInfo = append(info.CourseInfo, &types.CourseInfo{
					CourseSort: tableDeploy.CourseSort,
					StartTime:  startTime.Format("15:04:05"),
					EndTime:    endTime.Format("15:04:05"),
					CourseName: tableDeploy.Note,
					CourseFlag: tableDeploy.CourseFlag,
				})
			}
		}
	}

	for _, info := range resp.CourseTable {
		for _, table := range courseTables {
			if table.WeekDay == info.Weekday {
				for _, courseInfo := range info.CourseInfo {
					if courseInfo.CourseSort == table.CourseSort {
						courseInfo.Id = table.Id
						courseInfo.SubjectName = table.SubjectName
						courseInfo.TeacherName = table.TeacherName
						courseInfo.ClassId = table.ClassId
						courseInfo.ClassName = table.ClassName
						courseInfo.SubjectId = table.SubjectId
					}
				}
			}
		}
	}

	return resp, nil
}
