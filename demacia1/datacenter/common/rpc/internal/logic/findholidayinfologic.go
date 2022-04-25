package logic

import (
	"context"
	"demacia/common/errlist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"demacia/datacenter/common/rpc/common"
	"demacia/datacenter/common/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindHolidayInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindHolidayInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindHolidayInfoLogic {
	return &FindHolidayInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindHolidayInfoLogic) FindHolidayInfo(in *common.HolidayReq) (*common.HolidayResp, error) {
	resp := &common.HolidayResp{SpecialDate: []string{}}
	if in.Year == 0 {
		return nil, status.Error(codes.NotFound, errlist.InvalidParam.Error())
	}

	list, err := l.svcCtx.HolidayModel.SelectList(in.Year)
	if err != nil {
		return nil, status.Error(codes.NotFound, errlist.Unknown.Error())
	}

	for _, holiday := range list {
		resp.SpecialDate = append(resp.SpecialDate, holiday.SpecialDate)
	}

	return resp, nil
}
