package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"
	"demacia/datacenter/department/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type InsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) InsertLogic {
	return InsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertLogic) Insert(req types.InsertReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.AuthLoginFail
	}
	_, err = l.svcCtx.DepartmentModel.GetDepartmentByTitle(orgId, req.Title)
	if err == nil {
		return errlist.DepartmentExit
	}
	departmentId, err := l.svcCtx.DepartmentModel.InsertDepartment(&model.Department{
		OrgId:       orgId,
		Title:       req.Title,
		Sort:        0,
		MemberCount: 0,
	})
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		s := datacenter.Marshal(datacenter.Department, departmentId, datacenter.Update, datacenter.DepartmentData{Title: req.Title})
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
