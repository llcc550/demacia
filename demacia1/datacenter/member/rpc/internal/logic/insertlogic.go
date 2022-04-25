package logic

import (
	"context"

	"demacia/common/basefunc"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/model"
	"demacia/datacenter/member/rpc/internal/svc"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type InsertLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InsertLogic {
	return &InsertLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InsertLogic) Insert(in *member.InsertReq) (*member.IdReq, error) {
	data := &model.Member{
		OrgId:    in.OrgId,
		UserName: in.UserName,
		Password: basefunc.HashPassword(in.Password),
		Mobile:   in.Mobile,
		TrueName: in.TrueName,
		Status:   0,
	}

	memberId, err := l.svcCtx.MemberModel.Insert(data)
	if err != nil {
		return nil, errlist.Unknown
	}
	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Create(context.Background(), &databus.Req{Topic: datacenter.Member, ObjectId: memberId})
	})
	return &member.IdReq{Id: memberId}, nil
}
