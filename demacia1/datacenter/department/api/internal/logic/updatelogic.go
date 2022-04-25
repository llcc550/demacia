package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateLogic {
	return UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req types.UpdateReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.AuthLoginFail
	}
	departmentInfo, err := l.svcCtx.DepartmentModel.GetDepartmentById(req.Id)
	if err != nil || departmentInfo.OrgId != orgId {
		return errlist.DepartmentNotExit
	}
	if departmentInfo.Title == req.Title {
		return nil
	}
	_, err = l.svcCtx.DepartmentModel.GetDepartmentByTitle(orgId, req.Title)
	if err == nil {
		return errlist.DepartmentExit
	}
	err = l.svcCtx.DepartmentModel.UpdateDepartmentTitle(req.Id, req.Title)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		s := datacenter.Marshal(datacenter.Department, req.Id, datacenter.Update, datacenter.DepartmentData{Title: req.Title})
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
