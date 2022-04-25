package photos

import (
	"context"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EditpublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) EditpublishLogic {
	return EditpublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditpublishLogic) Editpublish(req types.EditReq) error {
	// todo: add your logic here and delete this line

	return nil
}
