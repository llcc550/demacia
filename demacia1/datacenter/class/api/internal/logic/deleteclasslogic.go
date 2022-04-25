package logic

import (
	"context"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteClassLogic {
	return DeleteClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteClassLogic) DeleteClass(req types.Id) error {

	if req.Id <= 0 {
		return errlist.ClassNotFound
	}
	err := l.svcCtx.ClassModel.DeleteById(req.Id)

	if err != nil {
		return err
	}
	err = l.svcCtx.ClassTeacherModel.DeletedByClassId(req.Id)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		s := datacenter.Marshal(datacenter.Class, req.Id, datacenter.Delete, nil)
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
