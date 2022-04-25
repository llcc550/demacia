package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteGradeByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteGradeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteGradeByIdLogic {
	return DeleteGradeByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGradeByIdLogic) DeleteGradeById(req types.Id) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}

	err = l.svcCtx.GradeModel.DeleteGradeById(orgId, req.Id)

	if err != nil {
		return err
	}

	return nil
}
