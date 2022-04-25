package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/student/rpc/student"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type MyCourseTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMyCourseTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) MyCourseTableLogic {
	return MyCourseTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MyCourseTableLogic) MyCourseTable(req types.MyCourseTableReq) (*types.CourseTableInfoReply, error) {
	resp := &types.CourseTableInfoReply{CourseTable: []*types.CourseTableInfo{}}
	if req.Type == 0 || req.Type == 1 && req.StudentId == 0 || req.Type == 2 && req.MemberId == 0 {
		return resp, errlist.InvalidParam
	}
	for i := 1; i < 8; i++ {
		resp.CourseTable = append(resp.CourseTable, &types.CourseTableInfo{
			Weekday:    int8(i),
			CourseInfo: []*types.CourseInfo{},
		})
	}
	var courseTables model.CourseTables
	var err error
	if req.Type == 1 {
		studentInfo, err := l.svcCtx.StudentRpc.FindOneById(l.ctx, &student.IdRequest{Id: req.StudentId})
		if err != nil {
			return resp, errlist.Unknown
		}
		courseTables, err = l.svcCtx.CourseTableModel.SelectByCid(studentInfo.ClassId)
		if err != nil {
			if err != sql.ErrNoRows {
				return resp, errors.CourseNotFoundErr
			}
			return resp, errlist.Unknown
		}
	} else if req.Type == 2 {
		courseTables, err = l.svcCtx.CourseTableModel.SelectByMemberId(req.MemberId)
		if err != nil {
			if err != sql.ErrNoRows {
				return resp, errors.CourseNotFoundErr
			}
			return resp, errlist.Unknown
		}
	}
	for _, info := range resp.CourseTable {
		for _, table := range courseTables {
			if info.Weekday == table.WeekDay {
				endTime, _ := time.Parse("2006-01-02T15:04:05Z", table.EndTime)
				startTime, _ := time.Parse("2006-01-02T15:04:05Z", table.StartTime)
				if req.Type == 1 {
					info.CourseInfo = append(info.CourseInfo, &types.CourseInfo{
						CourseName:  table.SubjectName,
						TeacherName: table.TeacherName,
						StartTime:   startTime.Format("15:04:05"),
						EndTime:     endTime.Format("15:04:05"),
					})
				} else if req.Type == 2 {
					info.CourseInfo = append(info.CourseInfo, &types.CourseInfo{
						ClassName:    table.ClassName,
						StartTime:    startTime.Format("15:04:05"),
						EndTime:      endTime.Format("15:04:05"),
						CourseName:   table.SubjectName,
						PositionName: table.PositionName,
					})
				}
			}
		}
	}
	return resp, nil
}
