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

type FindOneByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOneByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOneByIdLogic {
	return &FindOneByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOneByIdLogic) FindOneById(in *member.IdReq) (*member.MemberInfo, error) {
	memberInfo, err := l.svcCtx.MemberModel.FindOneById(in.Id)
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
		Avatar:   memberInfo.Avatar,
		Role:     memberInfo.Role,
	}, nil
}
