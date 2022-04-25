package logic

import (
	"context"
	"fmt"

	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/common/rpc/common"
	"demacia/datacenter/common/rpc/commonclient"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/rpc/member"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/datacenter/organization/api/internal/errors"
	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"
	"demacia/datacenter/organization/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddLogic {
	return AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req types.OrgAddReq) (*types.IdReq, error) {
	var areaInfo *commonclient.AreaResp
	err := mr.Finish(func() error {
		_, err := l.svcCtx.OrganizationModel.FindIdByTitle(req.Title)
		if err == nil {
			return errors.OrganizationExist
		}
		return nil
	}, func() error {
		checkManager, err := l.svcCtx.MemberRpc.FindOneByUserName(l.ctx, &memberclient.UserNameReq{UserName: req.UserName})
		if err == nil && checkManager.Id != 0 {
			return errors.ManagerExist
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
		return nil, err
	}

	orgId, err := l.svcCtx.OrganizationModel.Insert(&model.Organization{
		Title:        req.Title,
		OrgType:      req.IsSchool,
		ActivateDate: req.ActivateDate,
		ExpireDate:   req.ExpireDate,
		ProvinceId:   areaInfo.ProvinceId,
		CityId:       areaInfo.CityId,
		AreaId:       areaInfo.AreaId,
		AgentId:      req.AgentId,
		AreaTitle:    fmt.Sprintf("%s-%s-%s", areaInfo.ProvinceTitle, areaInfo.CityTitle, areaInfo.AreaTitle),
		Addr:         req.Addr,
		Msg:          req.Msg,
		TermId:       req.TermId,
		TrueName:     req.TrueName,
		Mobile:       req.Mobile,
	})
	if err != nil {
		l.Logger.Errorf("insert organization err:%s", err.Error())
		return nil, errlist.Unknown
	}
	memberId, err := l.svcCtx.MemberRpc.Insert(l.ctx, &member.InsertReq{
		UserName: req.UserName,
		OrgId:    orgId,
		TrueName: req.TrueName,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("insert organization manager err:%s", err.Error())
		_ = l.svcCtx.OrganizationModel.DeleteByIdForce(orgId)
		return nil, errlist.Unknown
	}
	_ = l.svcCtx.OrganizationModel.UpdateManagerMember(orgId, memberId.Id, req.UserName)
	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Create(context.Background(), &databus.Req{Topic: datacenter.Organization, ObjectId: orgId})
	})
	return &types.IdReq{Id: orgId}, nil
}
