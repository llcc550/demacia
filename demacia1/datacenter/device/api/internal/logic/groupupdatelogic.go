package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"
	"demacia/datacenter/device/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupUpdateLogic {
	return GroupUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUpdateLogic) GroupUpdate(req types.Group) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	groupInfo, err := l.svcCtx.GroupModel.FindOne(req.GroupId)
	if err != nil || groupInfo.OrgId != orgId {
		return errlist.NoAuth
	}
	groupInfoByName, err := l.svcCtx.GroupModel.FindOneByGroupName(orgId, req.GroupName)
	if err == nil && groupInfoByName.Id != req.GroupId {
		return errlist.GroupExist
	}
	err = l.svcCtx.GroupModel.Update(&model.Group{
		Id:    req.GroupId,
		OrgId: orgId,
		Name:  req.GroupName,
	})
	if err != nil {
		return err
	}
	return nil
}
