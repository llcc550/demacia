package logic

import (
	"context"

	"demacia/datacenter/organization/api/internal/errors"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) DetailLogic {
	return DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req types.IdReq) (*types.OrgInfo, error) {
	detail, err := l.svcCtx.OrganizationModel.FindOneById(req.Id)
	if err != nil {
		return nil, errors.OrganizationNotExist
	}

	return &types.OrgInfo{
		Id:           detail.Id,
		AreaTitle:    detail.AreaTitle,
		Title:        detail.Title,
		ActivateDate: detail.ActivateDate,
		ExpireDate:   detail.ExpireDate,
		TrueName:     detail.TrueName,
		Mobile:       detail.Mobile,
		UserName:     detail.ManagerMemberUserName,
		Addr:         detail.Addr,
		Msg:          detail.Msg,
		CreatedTime:  detail.CreateTime,
		OrgStatus:    detail.OrgStatus,
	}, nil
}
