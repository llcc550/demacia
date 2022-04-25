package logic

import (
	"context"
	"database/sql"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/cloudscreen/courserecord/model"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"fmt"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CourseRecordConfigSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseRecordConfigSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseRecordConfigSaveLogic {
	return CourseRecordConfigSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseRecordConfigSaveLogic) CourseRecordConfigSave(req types.CourseRecordConfigSaveReq) (*types.SuccessReply, error) {
	resp := &types.SuccessReply{}
	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}
	if req.SignPerson == 0 || req.SignBeforeTime == 0 || req.Enable == 0 || req.SignHoliday == 0 {
		return resp, errlist.InvalidParam
	}

	config, err := l.svcCtx.CourseRecordConfigModel.SelectByOrgId(oid)
	if err != nil {
		if err != sql.ErrNoRows {
			return resp, errlist.Unknown
		}
	}
	config.OrgId = oid
	config.SignPerson = int8(req.SignPerson)
	config.SignBeforeTime = int8(req.SignBeforeTime)
	if req.Enable == 1 {
		config.Enable = true
	}
	if req.SignHoliday == 1 {
		config.SignHoliday = true
	}
	if config.Id == 0 {
		if err := l.svcCtx.CourseRecordConfigModel.InsertCourseRecordConfig(&config); err != nil {
			return resp, errlist.Unknown
		}
	} else {
		if err := l.svcCtx.CourseRecordConfigModel.UpdateCourseRecordConfig(&config); err != nil {
			return resp, errlist.Unknown
		}
	}
	fmt.Println(len(req.SpecialDates))
	if len(req.SpecialDates) != 0 {
		specialDates := make(model.CourseRecordDates, 0, len(req.SpecialDates))
		for _, date := range req.SpecialDates {
			specialDates = append(specialDates, &model.CourseRecordDate{
				OrgId:       oid,
				SpecialDate: date.Date,
				Type:        int8(date.Type),
				Year:        int64(date.Year),
				IsHoliday:   2,
			})
		}
		if err := l.svcCtx.CourseRecordDateModel.InsertDates(oid, specialDates); err != nil {
			l.Logger.Errorf("insert course_record_dates err:%s", err.Error())
		}
	}
	resp.Success = true
	return resp, nil
}
