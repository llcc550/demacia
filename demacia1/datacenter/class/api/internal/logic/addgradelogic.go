package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddGradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddGradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddGradeLogic {
	return AddGradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddGradeLogic) AddGrade(req types.Info) (*types.Id, error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	stage, err := l.svcCtx.StageModel.GetStageById(req.Id)
	if err != nil {
		return nil, err
	}
	insertId, err := l.svcCtx.GradeModel.Insert(orgId, req.Id, req.Title, stage.Title)
	if err != nil {
		return nil, err
	}
	yearNum, _ := l.svcCtx.GradeModel.GetStageYearByStageId(req.Id)
	err = l.svcCtx.StageModel.UpdateStageYear(orgId, req.Id, int64(yearNum))
	if err != nil {
		return nil, err
	}
	return &types.Id{Id: insertId}, nil
}
