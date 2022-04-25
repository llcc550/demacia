package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type StagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) StagesLogic {
	return StagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StagesLogic) Stages() (*types.ListRespose, error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.StageModel.ListByOrgId(orgId)
	if err != nil {
		return nil, err
	}
	res := make([]types.Info, 0, len(list))
	for _, i := range list {
		res = append(res, types.Info{
			Id:    i.Id,
			Title: i.Title,
		})
	}
	return &types.ListRespose{List: res}, nil
}
