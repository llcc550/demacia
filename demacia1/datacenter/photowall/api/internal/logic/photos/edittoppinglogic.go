package photos

import (
	"context"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type EdittoppingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEdittoppingLogic(ctx context.Context, svcCtx *svc.ServiceContext) EdittoppingLogic {
	return EdittoppingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EdittoppingLogic) Edittopping(req types.EditReq) error {
	// todo: add your logic here and delete this line

	return nil
}
