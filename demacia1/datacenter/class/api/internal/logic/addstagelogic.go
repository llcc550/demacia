package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/class/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddStageLogic {
	return AddStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddStageLogic) AddStage(req *types.StageInfo) (*types.Id, error) {

	orgId, _ := baseauth.GetOrgId(l.ctx)
	if req.Year > 20 {
		return nil, errlist.StageYearTooLong
	}
	stageInfo, err := l.svcCtx.StageModel.GetStageByTitleAndOrg(req.Title, orgId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if stageInfo != nil {
		return nil, errlist.StageNameExist
	}
	id, err := l.svcCtx.StageModel.InsertStageOfOrgId(&model.Stage{
		OrgId: orgId,
		Title: req.Title,
		Year:  req.Year,
	})
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.GradeModel.AddGradeNumByStageId(orgId, id, req.Title, req.Year)
	if err != nil {
		return nil, err
	}
	return &types.Id{Id: id}, nil
}
