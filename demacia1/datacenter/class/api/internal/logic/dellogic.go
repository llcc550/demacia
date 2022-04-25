package logic

import (
	"context"

	"demacia/common/datacenter"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type DelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) DelLogic {
	return DelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLogic) Del(req types.Id) error {
	// todo: 参数合法性验证
	err := l.svcCtx.ClassModel.DeleteById(req.Id)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		s := datacenter.Marshal(datacenter.Class, req.Id, datacenter.Delete, nil)
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
