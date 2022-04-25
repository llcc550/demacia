package logic

import (
	"context"

	"demacia/common/errlist"
	"demacia/datacenter/common/rpc/common"
	"demacia/datacenter/common/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FindAreaInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAreaInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAreaInfoLogic {
	return &FindAreaInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindAreaInfoLogic) FindAreaInfo(in *common.AreaReq) (resp *common.AreaResp, errRes error) {
	errRes = status.Error(codes.NotFound, errlist.CommonAreaIdErr.Error())
	areaInfo, err := l.svcCtx.AreaModel.FindOneById(in.AreaId)
	if err != nil {
		return
	}
	cityInfo, err := l.svcCtx.AreaModel.FindOneById(areaInfo.Pid)
	if err != nil {
		return
	}
	provinceInfo, err := l.svcCtx.AreaModel.FindOneById(cityInfo.Pid)
	if err != nil {
		return
	}
	return &common.AreaResp{
		ProvinceId:    provinceInfo.Id,
		CityId:        cityInfo.Id,
		AreaId:        areaInfo.Id,
		ProvinceTitle: provinceInfo.Name,
		CityTitle:     cityInfo.Name,
		AreaTitle:     areaInfo.Name,
	}, nil
}
