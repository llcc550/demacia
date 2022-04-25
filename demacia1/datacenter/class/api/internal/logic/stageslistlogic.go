package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type StagesListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStagesListLogic(ctx context.Context, svcCtx *svc.ServiceContext) StagesListLogic {
	return StagesListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StagesListLogic) StagesList() (resp *types.StageListRespose, err error) {
	orgId, _ := baseauth.GetOrgId(l.ctx)
	list, err := l.svcCtx.StageModel.ListByOrgId(orgId)
	if err != nil {
		return nil, err
	}

	resp = &types.StageListRespose{List: []*types.StageInfo{}}

	for _, item := range list {
		resp.List = append(resp.List, &types.StageInfo{
			Id:    item.Id,
			Title: item.Title,
			Year:  item.Year,
		})
	}
	return resp, nil
}
