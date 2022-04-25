package logic

import (
	"context"

	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/department/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetDepartmentIdsByOrgIdAndMemberIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentIdsByOrgIdAndMemberIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentIdsByOrgIdAndMemberIdLogic {
	return &GetDepartmentIdsByOrgIdAndMemberIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDepartmentIdsByOrgIdAndMemberId 根据机构ID和人员ID获取该人员所在的部门ID列表
func (l *GetDepartmentIdsByOrgIdAndMemberIdLogic) GetDepartmentIdsByOrgIdAndMemberId(in *department.OrgIdAndMemberIdReq) (*department.DepartmentIdsResp, error) {
	list, err := l.svcCtx.DepartmentMemberModel.GetDepartmentIdsByOrgIdAndMemberId(in.OrgId, in.MemberId)
	if err != nil {
		return nil, err
	}
	return &department.DepartmentIdsResp{DepartmentIds: list}, nil
}
