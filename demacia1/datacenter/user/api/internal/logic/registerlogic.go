package logic

import (
	"context"

	"demacia/common/basefunc"
	"demacia/datacenter/user/api/internal/svc"
	"demacia/datacenter/user/api/internal/types"
	"demacia/datacenter/user/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterLogic {
	return RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req types.RegisterReq) error {
	_, _ = l.svcCtx.UserModel.Insert(&model.User{
		UserName: req.UserName,
		Password: basefunc.HashPassword(req.Password),
		Mobile:   req.Mobile,
		TrueName: req.TrueName,
	})
	return nil
}
