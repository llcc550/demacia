package photos

import (
	"context"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EditlockscreenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditlockscreenLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditlockscreenLogic {
	return EditlockscreenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditlockscreenLogic) Editlockscreen(req types.EditReq) error {
	// todo: add your logic here and delete this line

	return nil
}
