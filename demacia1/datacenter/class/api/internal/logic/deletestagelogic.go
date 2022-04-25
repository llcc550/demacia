package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteStageLogic {
	return DeleteStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteStageLogic) DeleteStage(req types.Id) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = l.svcCtx.StageModel.DeleteStageById(orgId, req.Id)
	if err != nil {
		return err
	}
	err = l.svcCtx.GradeModel.DeleteGradeByStageId(orgId, req.Id)
	if err != nil {
		return err
	}
	return nil
}
