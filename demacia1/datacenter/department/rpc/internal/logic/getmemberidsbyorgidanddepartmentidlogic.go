package logic

import (
	"context"

	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/department/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetMemberIdsByOrgIdAndDepartmentIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMemberIdsByOrgIdAndDepartmentIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemberIdsByOrgIdAndDepartmentIdLogic {
	return &GetMemberIdsByOrgIdAndDepartmentIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMemberIdsByOrgIdAndDepartmentId 根据机构ID和部门ID获取部门下的人员ID列表
func (l *GetMemberIdsByOrgIdAndDepartmentIdLogic) GetMemberIdsByOrgIdAndDepartmentId(in *department.OrgIdAndDepartmentIdReq) (*department.MemberIdsResp, error) {
	list, err := l.svcCtx.DepartmentMemberModel.GetMembersByOrgIdAndDepartmentId(in.OrgId, in.DepartmentId)
	if err != nil {
		return nil, err
	}
	resp := &department.MemberIdsResp{}
	for _, item := range list {
		resp.MemberIds = append(resp.MemberIds, item.MemberId)
	}
	return resp, nil
}
