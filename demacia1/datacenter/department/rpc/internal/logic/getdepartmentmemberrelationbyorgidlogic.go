package logic

import (
	"context"

	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/department/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetDepartmentMemberRelationByOrgIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentMemberRelationByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentMemberRelationByOrgIdLogic {
	return &GetDepartmentMemberRelationByOrgIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDepartmentMemberRelationByOrgId 根据机构ID获取所有的部门人员关系
func (l *GetDepartmentMemberRelationByOrgIdLogic) GetDepartmentMemberRelationByOrgId(in *department.OrgIdReq) (*department.DepartmentMembersResp, error) {
	list, err := l.svcCtx.DepartmentMemberModel.GetDepartmentMembersByOrgId(in.OrgId)
	if err != nil {
		return nil, err
	}
	resp := &department.DepartmentMembersResp{}
	for _, item := range list {
		resp.DepartmentMembers = append(resp.DepartmentMembers, &department.DepartmentMember{
			DepartmentId: item.DepartmentId,
			MemberId:     item.MemberId,
		})
	}
	return resp, nil
}
