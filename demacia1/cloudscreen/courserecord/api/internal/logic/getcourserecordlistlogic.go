package logic

import (
	"context"
	"database/sql"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/cloudscreen/courserecord/errors"
	"demacia/common/errlist"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetCourseRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCourseRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCourseRecordListLogic {
	return GetCourseRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCourseRecordListLogic) GetCourseRecordList(req types.CourseRecordReq) (*types.CourseRecordReply, error) {
	resp := &types.CourseRecordReply{CourseRecords: []*types.RecordInfo{}}
	if req.Type == 0 {
		return resp, errlist.InvalidParam
	}
	var userId int64
	if req.Type == 1 {
		userId = req.StudentId
	} else if req.Type == 2 {
		userId = req.MemberId
	}
	if userId == 0 {
		return resp, errlist.InvalidParam
	}
	records, count, err := l.svcCtx.CourseRecordModel.SelectByStudentId(userId, req.Type, req.Page, req.Limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, errors.SignRecordNotFound
		}
		return resp, errlist.Unknown
	}
	for _, record := range records {
		signDate, _ := time.Parse("2006-01-02T15:04:05Z", record.SignDate)
		signTime, _ := time.Parse("2006-01-02T15:04:05Z", record.SignTime)
		resp.CourseRecords = append(resp.CourseRecords, &types.RecordInfo{
			Name:         record.Truename,
			Date:         signDate.Format("2006-01-02"),
			SignTime:     signTime.Format("15:04:05"),
			Note:         record.CourseNote,
			SubjectName:  record.SubjectName,
			PositionName: record.PositionName,
			Status:       record.Status,
			Photo:        record.Photo,
		})
	}
	resp.Count = count
	return resp, nil
}
