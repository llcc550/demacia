package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type UpdateGradeByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateGradeByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateGradeByIdLogic {
	return UpdateGradeByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGradeByIdLogic) UpdateGradeById(req types.Info) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = l.svcCtx.GradeModel.Update(orgId, req.Id, req.Title)
	if err != nil {
		return err
	}
	return nil
}
