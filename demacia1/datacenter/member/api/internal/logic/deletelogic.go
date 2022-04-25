package logic

import (
	"context"

	"demacia/common/datacenter"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteLogic {
	return DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req types.MemberIds) error {
	err := l.svcCtx.MemberModel.DelMembers(req.MemberIds)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		for _, memberId := range req.MemberIds {
			_, _ = l.svcCtx.DataBusRpc.Delete(context.Background(), &databus.Req{Topic: datacenter.Member, ObjectId: memberId})
		}
	})
	return nil
}
