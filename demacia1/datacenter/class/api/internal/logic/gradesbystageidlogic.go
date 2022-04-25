package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GradesByStageIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGradesByStageIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) GradesByStageIdLogic {
	return GradesByStageIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GradesByStageIdLogic) GradesByStageId(req types.Id) (*types.ListRespose, error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.GradeModel.ListByOrgIdAndStageId(orgId, req.Id)
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
