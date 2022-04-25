package logic

import (
	"context"

	"demacia/common/cachemodel"
	"demacia/datacenter/member/rpc/internal/svc"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (l *FindOneByUserNameLogic) FindOneByUserName(in *member.UserNameReq) (*member.MemberInfo, error) {
	memberInfo, err := l.svcCtx.MemberModel.FindOneByUserName(in.UserName)
	if err != nil {
		if err == cachemodel.ErrNotFound {
			return nil, status.Error(codes.NotFound, "教师不存在")
		}
		return nil, err
	}
	if memberInfo.Status != 1 {
		return nil, status.Error(codes.NotFound, "教师不存在")
	}
	return &member.MemberInfo{
		Id:       memberInfo.Id,
		OrgId:    memberInfo.OrgId,
		UserName: memberInfo.UserName,
		TrueName: memberInfo.TrueName,
		Mobile:   memberInfo.Mobile,
		Password: memberInfo.Password,
	}, nil
}
