package logic

import (
	"context"

	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/department/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetDepartmentsByOrgIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentsByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentsByOrgIdLogic {
	return &GetDepartmentsByOrgIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDepartmentsByOrgId 根据机构ID获取所有的部门
func (l *GetDepartmentsByOrgIdLogic) GetDepartmentsByOrgId(in *department.OrgIdReq) (*department.DepartmentListResp, error) {
	list, err := l.svcCtx.DepartmentModel.GetDepartmentsByOrgId(in.OrgId)
	if err != nil {
		return nil, err
	}
	resp := &department.DepartmentListResp{}
	for _, item := range list {
		resp.Departments = append(resp.Departments, &department.DepartmentInfo{
			DepartmentId:    item.Id,
			OrgId:           item.OrgId,
			DepartmentTitle: item.Title,
			Sort:            item.Sort,
			MemberCount:     item.MemberCount,
		})
	}
	return resp, nil
}
