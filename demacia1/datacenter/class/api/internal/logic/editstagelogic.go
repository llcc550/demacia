package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EditStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditStageLogic {
	return EditStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditStageLogic) EditStage(req types.Info) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = l.svcCtx.StageModel.UpdateStageById(orgId, req.Id, req.Title)
	if err != nil {
		return err
	}
	err = l.svcCtx.GradeModel.UpdateStageTitleByStageId(orgId, req.Id, req.Title)
	if err != nil {
		return err
	}
	return nil
}
