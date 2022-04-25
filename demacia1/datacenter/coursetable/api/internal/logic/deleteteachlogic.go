package logic

import (
	"context"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteTeachLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTeachLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteTeachLogic {
	return DeleteTeachLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTeachLogic) DeleteTeach(req types.DeleteTeachReq) error {
	threading.GoSafe(func() {
		for _, id := range req.Ids {
			_ = l.svcCtx.TeachModel.DeleteById(id)
			_ = l.svcCtx.CourseTableModel.DeleteTeacherInfo(id)
		}
	})
	return nil
}
