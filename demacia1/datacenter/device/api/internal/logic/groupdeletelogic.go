package logic

import (
	"context"
	"demacia/common/baseauth"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/device/api/internal/svc"
	"demacia/datacenter/device/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GroupDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) GroupDeleteLogic {
	return GroupDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupDeleteLogic) GroupDelete(req types.GroupIdReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	groupInfo, err := l.svcCtx.GroupModel.FindOne(req.GroupId)
	if err != nil || groupInfo.OrgId != orgId {
		return err
	}
	err = l.svcCtx.GroupModel.Delete(req.GroupId)
	threading.GoSafe(func() {
		_ = l.svcCtx.DeviceGroupModel.DeleteByGroupId(req.GroupId)
	})
	return err
}
