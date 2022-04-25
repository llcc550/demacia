package logic

import (
	"context"
	"demacia/cloudscreen/timeswitch/errors"
	"demacia/cloudscreen/timeswitch/model"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"time"

	"demacia/cloudscreen/timeswitch/api/internal/svc"
	"demacia/cloudscreen/timeswitch/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EditTimeSwitchInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditTimeSwitchInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditTimeSwitchInfoLogic {
	return EditTimeSwitchInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditTimeSwitchInfoLogic) EditTimeSwitchInfo(req types.TimeSwitchInfo) (*types.SuccessResp, error) {
	resp := &types.SuccessResp{}
	if len(req.DeviceIds) == 0 {
		return resp, errlist.InvalidParam
	}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}

	for _, rangeList := range req.TimeRangeList {
		for _, timeRange := range rangeList.TimeRange {
			startTime, err := time.Parse("15:04:05", timeRange.StartTime)
			if err != nil {
				return resp, errors.SwitchTimeErr
			}
			endTime, err := time.Parse("15:04:05", timeRange.EndTime)
			if err != nil {
				return resp, errors.SwitchTimeErr
			}
			if startTime.After(endTime) {
				return resp, errors.SwitchTimeErr
			}
		}
		for i := 0; i < len(rangeList.TimeRange)-1; i++ {
			startTime, err := time.Parse("15:04:05", rangeList.TimeRange[i+1].StartTime)
			if err != nil {
				return resp, errors.SwitchTimeErr
			}
			endTime, err := time.Parse("15:04:05", rangeList.TimeRange[i].EndTime)
			if err != nil {
				return resp, errors.SwitchTimeErr
			}
			if startTime.Before(endTime) {
				return resp, errors.SwitchTimeRangeErr
			}
		}
	}

	dateRange := model.TimeSwitchDates{}
	for _, d := range req.SpecialDate {
		for _, id := range req.DeviceIds {
			startDate, err := time.Parse("2006-01-02", d.StartDate)
			if err != nil {
				return resp, errors.SwitchDateErr
			}
			endDate, err := time.Parse("2006-01-02", d.EndDate)
			if err != nil {
				return resp, errors.SwitchDateErr
			}
			dateRange = append(dateRange, &model.TimeSwitchDate{
				DeviceId:  id,
				StartDate: startDate,
				EndDate:   endDate,
			})
		}
	}
	if err := l.svcCtx.TimeSwitchDateModel.InsertTimeSwitchDates(req.DeviceIds, dateRange); err != nil {
		return resp, errlist.Unknown
	}

	configs := model.TimeSwitchConfigs{}
	for _, id := range req.DeviceIds {
		configs = append(configs, &model.TimeSwitchConfig{
			OrgId:       oid,
			DeviceId:    id,
			HolidayFlag: req.HolidayFlag,
		})
	}

	if err := l.svcCtx.TimeSwitchConfigModel.InsertTimeSwitchConfigs(req.DeviceIds, configs); err != nil {
		return resp, errlist.Unknown
	}

	timeSwitches := model.TimeSwitches{}
	for _, id := range req.DeviceIds {
		for _, rangeList := range req.TimeRangeList {
			for _, timeRange := range rangeList.TimeRange {
				startTime, err := time.Parse("15:04:05", timeRange.StartTime)
				if err != nil {
					return resp, errors.SwitchTimeErr
				}
				endTime, err := time.Parse("15:04:05", timeRange.EndTime)
				if err != nil {
					return resp, errors.SwitchTimeErr
				}
				timeSwitches = append(timeSwitches, &model.TimeSwitch{
					DeviceId:  id,
					StartTime: startTime.Format("15:04:05"),
					EndTime:   endTime.Format("15:04:05"),
					Weekday:   rangeList.Weekday,
				})
			}
		}
	}
	if err := l.svcCtx.TimeSwitchModel.InsertTimeSwitches(req.DeviceIds, timeSwitches); err != nil {
		return resp, errlist.Unknown
	}

	resp.Success = true
	return resp, nil
}
