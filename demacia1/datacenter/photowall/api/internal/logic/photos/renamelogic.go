package photos

import (
	"context"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type RenameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenameLogic(ctx context.Context, svcCtx *svc.ServiceContext) RenameLogic {
	return RenameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenameLogic) Rename(req types.RenameReq) error {
	// todo: add your logic here and delete this line

	return nil
}
