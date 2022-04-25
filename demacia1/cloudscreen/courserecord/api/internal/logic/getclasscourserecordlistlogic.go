package logic

import (
	"context"
	"database/sql"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/common/errlist"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetClassCourseRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClassCourseRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetClassCourseRecordListLogic {
	return GetClassCourseRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClassCourseRecordListLogic) GetClassCourseRecordList(req types.ClassCourseRecordReq) (*types.ClassCourseRecordReply, error) {
	resp := &types.ClassCourseRecordReply{ClassCourseRecords: []*types.RecordInfo{}}
	if req.ClassId == 0 {
		return resp, errlist.InvalidParam
	}
	records, count, err := l.svcCtx.CourseRecordModel.SelectByClassIdAndParam(req.ClassId, req.Truename, req.SubjectName, req.QueryDate, req.Status, req.Page, req.Limit)
	if err != nil {
		if err != sql.ErrNoRows {
			return resp, errlist.Unknown
		}
		return resp, nil
	}
	for _, record := range records {
		date, _ := time.Parse("2006-01-02T15:04:05Z", record.SignDate)
		signTime, _ := time.Parse("2006-01-02T15:04:05Z", record.SignTime)
		resp.ClassCourseRecords = append(resp.ClassCourseRecords, &types.RecordInfo{
			Name:         record.Truename,
			ClassName:    record.ClassName,
			Date:         date.Format("2006-01-02"),
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
