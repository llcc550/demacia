package logic

import (
	"context"
	"database/sql"
	"demacia/cloudscreen/timeswitch/api/internal/svc"
	"demacia/cloudscreen/timeswitch/api/internal/types"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetTimeSwitchInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTimeSwitchInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTimeSwitchInfoLogic {
	return GetTimeSwitchInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTimeSwitchInfoLogic) GetTimeSwitchInfo(req types.IdReq) (*types.TimeSwitchInfo, error) {

	resp := &types.TimeSwitchInfo{
		SpecialDate:   []*types.DateRange{},
		TimeRangeList: []*types.TimeRangeList{},
	}

	if req.DeviceId == 0 {
		return resp, errlist.InvalidParam
	}

	_, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}

	config, err := l.svcCtx.TimeSwitchConfigModel.SelectByDeviceId(req.DeviceId)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			return resp, errlist.Unknown
		}
	}

	dates, err := l.svcCtx.TimeSwitchDateModel.SelectByDeviceId(req.DeviceId)
	if err != nil {
		return resp, errlist.Unknown
	}
	for _, dateRange := range dates {
		resp.SpecialDate = append(resp.SpecialDate, &types.DateRange{
			StartDate: dateRange.StartDate.Format("2006-01-02"),
			EndDate:   dateRange.EndDate.Format("2006-01-02"),
		})
	}

	timeSwitches, err := l.svcCtx.TimeSwitchModel.SelectByDeviceId(req.DeviceId)
	if err != nil {
		return resp, errlist.Unknown
	}

	for i := 1; i <= 7; i++ {
		resp.TimeRangeList = append(resp.TimeRangeList, &types.TimeRangeList{
			Weekday:   int8(i),
			TimeRange: []*types.TimeRange{},
		})
	}

	for _, timeSwitch := range timeSwitches {
		for _, timeRangeInfo := range resp.TimeRangeList {
			if timeSwitch.Weekday == timeRangeInfo.Weekday {

				startTime, _ := time.Parse("2006-01-02T15:04:05Z", timeSwitch.StartTime)
				endTime, _ := time.Parse("2006-01-02T15:04:05Z", timeSwitch.EndTime)
				timeRangeInfo.TimeRange = append(timeRangeInfo.TimeRange, &types.TimeRange{
					StartTime: startTime.Format("15:04:05"),
					EndTime:   endTime.Format("15:04:05"),
				})
			}
		}
	}

	resp.HolidayFlag = config.HolidayFlag

	return resp, nil
}
