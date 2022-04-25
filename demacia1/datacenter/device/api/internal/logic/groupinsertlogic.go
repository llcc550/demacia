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

type GroupInsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupInsertLogic {
	return GroupInsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInsertLogic) GroupInsert(req types.GroupInsertReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.GroupModel.FindOneByGroupName(orgId, req.GroupName)
	if err == nil {
		return errlist.GroupExist
	}
	_, err = l.svcCtx.GroupModel.Insert(&model.Group{
		OrgId: orgId,
		Name:  req.GroupName,
	})
	return err
}
