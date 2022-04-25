package logic

import (
	"context"
	"demacia/common/errlist"

	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetClassInfoByFullNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetClassInfoByFullNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClassInfoByFullNameLogic {
	return &GetClassInfoByFullNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetClassInfoByFullNameLogic) GetClassInfoByFullName(in *class.FullNameReq) (*class.ClassInfo, error) {
	orgId := in.OrgId
	fullName := in.FullName
	classInfo, err := l.svcCtx.ClassModel.GetClassByFullNameAndOrgId(orgId, fullName)
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
