package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
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

func (l *DelLogic) Del(req types.DelReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.AuthLoginFail
	}
	mr.ForEach(func(source chan<- interface{}) {
		for _, id := range req.Ids {
			source <- id
		}
	}, func(item interface{}) {
		id := item.(int64)
		departmentInfo, err := l.svcCtx.DepartmentModel.GetDepartmentById(id)
		if err != nil || departmentInfo.OrgId != orgId {
			return
		}
		err = l.svcCtx.DepartmentModel.DeleteById(id)
		if err != nil {
			return
		}
		s := datacenter.Marshal(datacenter.Department, id, datacenter.Delete, nil)
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
