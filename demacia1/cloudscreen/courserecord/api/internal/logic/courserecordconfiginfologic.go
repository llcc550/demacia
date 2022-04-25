package logic

import (
	"context"
	"database/sql"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CourseRecordConfigInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseRecordConfigInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseRecordConfigInfoLogic {
	return CourseRecordConfigInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseRecordConfigInfoLogic) CourseRecordConfigInfo() (*types.CourseRecordConfigReply, error) {
	resp := types.CourseRecordConfigReply{SpecialDates: []*types.SpecialDate{}}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return &resp, errlist.NoAuth
	}

	config, err := l.svcCtx.CourseRecordConfigModel.SelectByOrgId(oid)
	if err != nil {
		if err == sql.ErrNoRows {
			config, _ = l.svcCtx.CourseRecordConfigModel.SelectDefaultConfig()
		} else {
			return &resp, errlist.Unknown
		}
	}

	if config.SignHoliday {
		resp.SignHoliday = 1
	}
	if config.Enable {
		resp.Enable = 1
	}
	resp.SignPerson = int(config.SignPerson)
	resp.SignBeforeTime = int(config.SignBeforeTime)
	dates, err := l.svcCtx.CourseRecordDateModel.SelectByOrgId(oid)
	for _, date := range dates {
		d, _ := time.Parse("2006-01-02T15:04:05Z", date.SpecialDate)
		resp.SpecialDates = append(resp.SpecialDates, &types.SpecialDate{
			Year: int(date.Year),
			Date: d.Format("2006-01-02"),
			Type: int(date.Type),
		})
	}
	return &resp, nil
}
