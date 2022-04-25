package logic

import (
	"context"
	"fmt"

	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/common/rpc/common"
	"demacia/datacenter/common/rpc/commonclient"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/organization/api/internal/errors"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"
	"demacia/datacenter/organization/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateLogic {
	return UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req types.OrgUpdateReq) error {
	_, err := l.svcCtx.OrganizationModel.FindOneById(req.Id)
	if err != nil {
		return errors.OrganizationNotExist
	}
	var areaInfo *commonclient.AreaResp
	err = mr.Finish(func() error {
		orgId, err := l.svcCtx.OrganizationModel.FindIdByTitle(req.Title)
		if err == nil && orgId != req.Id {
			return errors.OrganizationExist
		}
		return nil
	}, func() error {
		areaResp, err := l.svcCtx.CommonRpc.FindAreaInfo(l.ctx, &common.AreaReq{AreaId: req.AreaId})
		if err != nil {
			return errlist.InvalidParam
		}
		areaInfo = areaResp
		return nil
	})
	if err != nil {
		return err
	}
	err = l.svcCtx.OrganizationModel.Update(&model.Organization{
		Id:           req.Id,
		Title:        req.Title,
		ActivateDate: req.ActivateDate,
		ExpireDate:   req.ExpireDate,
		ProvinceId:   areaInfo.ProvinceId,
		CityId:       areaInfo.CityId,
		AreaId:       areaInfo.AreaId,
		AreaTitle:    fmt.Sprintf("%s-%s-%s", areaInfo.ProvinceTitle, areaInfo.CityTitle, areaInfo.AreaTitle),
		Addr:         req.Addr,
		Msg:          req.Msg,
		TrueName:     req.TrueName,
		Mobile:       req.Mobile,
	})
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Update(context.Background(), &databus.Req{Topic: datacenter.Organization, ObjectId: req.Id})
	})
	return nil
}
