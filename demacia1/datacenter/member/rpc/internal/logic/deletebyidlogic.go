package logic

import (
	"context"

	"demacia/common/cachemodel"
	"demacia/common/datacenter"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/rpc/internal/svc"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteByIdLogic {
	return &DeleteByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteByIdLogic) DeleteById(in *member.IdReq) (*member.NullResp, error) {
	_, err := l.svcCtx.MemberModel.FindOneById(in.Id)
	if err != nil {
		if err == cachemodel.ErrNotFound {
			return nil, status.Error(codes.NotFound, "教师不存在")
		}
		return nil, err
	}
	_ = l.svcCtx.MemberModel.DeleteById(in.Id)
	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Delete(context.Background(), &databus.Req{Topic: datacenter.Member, ObjectId: in.Id})
	})
	return &member.NullResp{}, nil
}
