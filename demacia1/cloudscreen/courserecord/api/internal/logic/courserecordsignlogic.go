package logic

import (
	"context"
	"database/sql"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/cloudscreen/courserecord/errors"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CourseRecordSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseRecordSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseRecordSignLogic {
	return CourseRecordSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseRecordSignLogic) CourseRecordSign(req types.RecordsReq) (*types.RecordsReply, error) {
	resp := &types.RecordsReply{
		Count:      0,
		RecordList: []*types.RecordInfo{},
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}
	if req.UserType == 0 {
		return resp, errlist.InvalidParam
	}
	if req.StartDate == "" && req.EndDate != "" || req.StartDate != "" && req.EndDate == "" {
		return resp, errlist.InvalidParam
	}

	records, count, err := l.svcCtx.CourseRecordModel.SelectByParam(req.StartDate, req.EndDate, req.Truename, req.UserType, req.Status, req.Page, req.Limit, oid)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, errors.SignRecordNotFound
		}
		return resp, errlist.Unknown
	}
	resp.Count = count

	for _, record := range records {
		signDate, _ := time.Parse("2006-01-02T15:04:05Z", record.SignDate)
		signTime, _ := time.Parse("2006-01-02T15:04:05Z", record.SignTime)
		resp.RecordList = append(resp.RecordList, &types.RecordInfo{
			Name:         record.Truename,
			ClassName:    record.ClassName,
			Date:         signDate.Format("2006-01-02"),
			SignTime:     signTime.Format("15:04:05"),
			Note:         record.CourseNote,
			SubjectName:  record.SubjectName,
			PositionName: record.PositionName,
			Status:       record.Status,
			Photo:        record.Photo,
		})
	}

	return resp, nil
}
