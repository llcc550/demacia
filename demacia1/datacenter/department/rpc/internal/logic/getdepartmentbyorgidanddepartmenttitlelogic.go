package logic

import (
	"context"

	"demacia/common/errlist"
	"demacia/datacenter/department/rpc/department"
	"demacia/datacenter/department/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetDepartmentByOrgIdAndDepartmentTitleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentByOrgIdAndDepartmentTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentByOrgIdAndDepartmentTitleLogic {
	return &GetDepartmentByOrgIdAndDepartmentTitleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDepartmentByOrgIdAndDepartmentTitle 根据机构ID和部门名称获取部门信息
func (l *GetDepartmentByOrgIdAndDepartmentTitleLogic) GetDepartmentByOrgIdAndDepartmentTitle(in *department.OrgIdAndDepartmentTitleReq) (*department.DepartmentInfo, error) {
	info, err := l.svcCtx.DepartmentModel.GetDepartmentByTitle(in.OrgId, in.DepartmentTitle)
	if err != nil {
		return nil, status.Error(codes.NotFound, errlist.DepartmentNotExit.Error())
	}
	return &department.DepartmentInfo{
		DepartmentId:    info.Id,
		OrgId:           info.OrgId,
		DepartmentTitle: info.Title,
		Sort:            info.Sort,
		MemberCount:     info.MemberCount,
	}, nil
}
