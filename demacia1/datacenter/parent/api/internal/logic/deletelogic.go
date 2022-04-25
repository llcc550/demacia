package logic

import (
	"context"
	"demacia/common/datacenter"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/parent/api/internal/svc"
	"demacia/datacenter/parent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteLogic {
	return DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req types.IdsRequest) error {
	if len(req.ParentIds) == 0 {
		return nil
	}
	err := l.svcCtx.ParentModel.DeleteByParentIds(req.ParentIds)
	if err != nil {
		return err
	}
	err = l.svcCtx.StudentParentModel.DeleteByParentIds(req.ParentIds)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		for _, id := range req.ParentIds {
			s := datacenter.Marshal(datacenter.Parent, id, datacenter.Delete, nil)
			_ = l.svcCtx.KqPusher.Push(s)
		}
	})
	return nil
}
