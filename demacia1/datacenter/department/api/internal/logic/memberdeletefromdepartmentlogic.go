package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type MemberDeleteFromDepartmentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMemberDeleteFromDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) MemberDeleteFromDepartmentLogic {
	return MemberDeleteFromDepartmentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberDeleteFromDepartmentLogic) MemberDeleteFromDepartment(req types.MemberIdsReq) error {
	if req.DepartmentId <= 0 || len(req.MemberIds) == 0 {
		return errlist.InvalidParam
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.AuthLoginFail
	}
	departmentInfo, err := l.svcCtx.DepartmentModel.GetDepartmentById(req.DepartmentId)
	if err != nil || departmentInfo.OrgId != orgId {
		return errlist.DepartmentNotExit
	}
	err = l.svcCtx.DepartmentMemberModel.DeleteByOrgIdAndDepartmentIdAndMemberIds(orgId, req.DepartmentId, req.MemberIds)
	if err != nil {
		return err
	}
	list, err := l.svcCtx.DepartmentMemberModel.GetMembersByOrgIdAndDepartmentId(orgId, req.DepartmentId)
	if err != nil {
		return err
	}
	err = l.svcCtx.DepartmentModel.UpdateDepartmentMemberCount(req.DepartmentId, int64(len(list)))
	if err != nil {
		return err
	}
	return nil
}
