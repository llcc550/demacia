package logic

import (
	"context"
	"demacia/common/errlist"

	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetClassInfoByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetClassInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClassInfoByIdLogic {
	return &GetClassInfoByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetClassInfoByIdLogic) GetClassInfoById(in *class.IdReq) (*class.ClassInfo, error) {

	classInfo, err := l.svcCtx.ClassModel.GetClassById(in.Id)
	if err != nil {
		return nil, errlist.ClassNotFound
	}
	return &class.ClassInfo{
		Id:       classInfo.Id,
		StageId:  classInfo.StageId,
		GradeId:  classInfo.GradeId,
		FullName: classInfo.FullName,
		OrgId:    classInfo.OrgId,
	}, nil
}
