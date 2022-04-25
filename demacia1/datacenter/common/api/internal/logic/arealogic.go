package logic

import (
	"context"

	"demacia/datacenter/common/api/internal/svc"
	"demacia/datacenter/common/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AreaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaLogic(ctx context.Context, svcCtx *svc.ServiceContext) AreaLogic {
	return AreaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaLogic) Area(req types.AreaReq) (*types.AreaResp, error) {
	list, err := l.svcCtx.AreaModel.FindListByPid(req.Pid)
	if err != nil {
		return nil, err
	}
	res := make([]*types.Area, 0, len(list))
	for _, item := range list {
		res = append(res, &types.Area{
			Id:   item.Id,
			Pid:  item.Pid,
			Name: item.Name,
		})
	}
	return &types.AreaResp{List: res}, nil
}
