package logic

import (
	"context"

	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type FindGradeByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindGradeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindGradeByIdLogic {
	return &FindGradeByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindGradeByIdLogic) FindGradeById(in *class.IdReq) (*class.GradeInfo, error) {
	// todo: add your logic here and delete this line
	grade, err := l.svcCtx.GradeModel.FindGradeById(in.Id)
	if err != nil {
		return nil, err
	}
	return &class.GradeInfo{
		Id:    grade.Id,
		Title: grade.Title,
	}, nil
}
