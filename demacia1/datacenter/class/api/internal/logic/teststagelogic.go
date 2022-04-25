package logic

import (
	"context"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type TestStageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestStageLogic(ctx context.Context, svcCtx *svc.ServiceContext) TestStageLogic {
	return TestStageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestStageLogic) TestStage(req types.Id) error {

	return nil
}
