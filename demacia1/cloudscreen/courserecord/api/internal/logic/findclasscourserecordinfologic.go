package logic

import (
	"context"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/cloudscreen/courserecord/errors"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindClassCourseRecordInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFindClassCourseRecordInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) FindClassCourseRecordInfoLogic {
	return FindClassCourseRecordInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindClassCourseRecordInfoLogic) FindClassCourseRecordInfo(req types.InitReq) (*types.ClassCourseRecordInfoReply, error) {
	resp := &types.ClassCourseRecordInfoReply{
		StartTime:   "",
		EndTime:     "",
		NormalCount: 0,
		LateCount:   0,
		LackCount:   0,
		Course:      &types.CourseInfo{},
		Teacher:     &types.TeacherInfo{},
		Students:    []*types.StudentInfo{},
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}

	if req.ClassId == 0 {
		return resp, errlist.InvalidParam
	}

	config, err := l.svcCtx.CourseRecordConfigModel.SelectByOrgId(oid)
	if err != nil {
		return resp, errlist.Unknown
	}
	if !config.Enable {
		return resp, errors.NoEnable
	}
	if config.SignPerson > 0 {
		startTime := time.Now().Add(time.Duration(config.SignBeforeTime) * time.Minute).Format("15:04:05")
		records, err := l.svcCtx.CourseRecordModel.SelectBetweenStartTime(req.ClassId, startTime)
		if err != nil {
			return resp, errlist.Unknown
		}
		if len(records) == 0 {
			return resp, errors.SignTimeErr
		}
		var normalCount int
		var lackCount int
		var lateCount int
		for _, record := range records {
			if record.Status == 1 {
				normalCount++
			} else if record.Status == 2 {
				lackCount++
			} else {
				lateCount++
			}
			if record.UserType == 1 {
				resp.Students = append(resp.Students, &types.StudentInfo{
					StudentName: record.Truename,
					Status:      0,
					Photo:       "",
				})
			} else if record.UserType == 2 {
				resp.Teacher.TeacherName = record.Truename
				resp.Teacher.Photo = ""
				resp.Teacher.Status = 0
			}
		}
		resp.NormalCount = normalCount
		resp.LackCount = lackCount
		resp.LateCount = lateCount
		signTime, _ := time.Parse("2006-01-02T15:04:05Z", records[0].StartTime)
		resp.EndTime = signTime.Format("15:04:05")
		resp.StartTime = signTime.Add(-time.Duration(config.SignBeforeTime) * time.Minute).Format("15:04:05")
		resp.Course.Note = records[0].CourseNote
		resp.Course.SubjectName = records[0].SubjectName
	}
	return resp, nil
}
