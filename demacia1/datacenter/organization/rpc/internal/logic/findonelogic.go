package logic

import (
	"context"

	"demacia/common/cachemodel"
	"demacia/datacenter/organization/rpc/internal/svc"
	"demacia/datacenter/organization/rpc/organization"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FindOneLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneLogic {
	return &FindOneLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneLogic) FindOne(in *organization.IdReply) (*organization.OrgInfo, error) {
	orgInfo, err := l.svcCtx.OrganizationModel.FindOneById(in.Id)
	if err != nil {
		if err == cachemodel.ErrNotFound {
			return nil, status.Error(codes.NotFound, "机构不存在")
		}
		return nil, err
	}
	return &organization.OrgInfo{
		Id:                    orgInfo.Id,
		Title:                 orgInfo.Title,
		OrgStatus:             int64(orgInfo.OrgStatus),
		ManagerMemberId:       orgInfo.ManagerMemberId,
		ManagerMemberUserName: orgInfo.ManagerMemberUserName,
		AreaTitle:             orgInfo.AreaTitle,
		TrueName:              orgInfo.TrueName,
		Mobile:                orgInfo.Mobile,
		Addr:                  orgInfo.Addr,
		Msg:                   orgInfo.Msg,
		TermId:                orgInfo.TermId,
		ProvinceId:            orgInfo.ProvinceId,
		CityId:                orgInfo.CityId,
		AreaId:                orgInfo.AreaId,
		AgentId:               orgInfo.AgentId,
		CreateTime:            orgInfo.CreateTime,
		ActivateDate:          orgInfo.ActivateDate,
		ExpireDate:            orgInfo.ExpireDate,
	}, nil
}
