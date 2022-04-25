package logic

import (
	"context"

	"demacia/datacenter/common/api/internal/svc"
	"demacia/datacenter/common/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EthnicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEthnicLogic(ctx context.Context, svcCtx *svc.ServiceContext) EthnicLogic {
	return EthnicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EthnicLogic) Ethnic() (*types.EthnicResp, error) {
	list, err := l.svcCtx.EthnicModel.FindList()
	if err != nil {
		return nil, err
	}
	res := make([]*types.Ethnic, 0, len(list))
	for _, item := range list {
		res = append(res, &types.Ethnic{
			Id:   item.Id,
			Name: item.Name,
		})
	}
	return &types.EthnicResp{List: res}, nil
}
