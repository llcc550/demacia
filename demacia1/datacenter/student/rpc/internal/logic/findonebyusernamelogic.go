package logic

import (
	"context"

	"demacia/datacenter/student/errors"
	"demacia/datacenter/student/rpc/internal/svc"
	"demacia/datacenter/student/rpc/student"

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

func (l *FindOneByUserNameLogic) FindOneByUserName(in *student.UserNameRequest) (*student.StudentResponse, error) {
	info, err := l.svcCtx.StudentModel.FindOneByOrgIdAndUserName(in.OrgId, in.UserName)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.StudentExist.Error())
	}
	return &student.StudentResponse{
		Id:         info.Id,
		Name:       info.TrueName,
		UserName:   info.UserName,
		ClassName:  info.ClassFullName,
		Password:   info.Password,
		OrgId:      info.OrgId,
		Avatar:     info.Avatar,
		Face:       info.Face,
		IdNumber:   info.IdNumber,
		CardNumber: info.CardNumber,
		ClassId:    info.ClassId,
	}, nil
}
