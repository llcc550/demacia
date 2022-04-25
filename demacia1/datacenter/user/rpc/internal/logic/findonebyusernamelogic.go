package logic

import (
	"context"

	"demacia/datacenter/user/rpc/internal/svc"
	"demacia/datacenter/user/rpc/user"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindOneByUserNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneByUserNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneByUserNameLogic {
	return &FindOneByUserNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneByUserNameLogic) FindOneByUserName(in *user.UserNameReq) (*user.UserInfo, error) {
	memberInfo, err := l.svcCtx.UserModel.FindOneByUserName(in.UserName)
	if err != nil {
		return nil, err
	}
	return &user.UserInfo{
		Id:       memberInfo.Id,
		UserName: memberInfo.UserName,
		TrueName: memberInfo.TrueName,
		Mobile:   memberInfo.Mobile,
		Password: memberInfo.Password,
	}, nil
}
