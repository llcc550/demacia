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

type GetDepartmentByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentByIdLogic {
	return &GetDepartmentByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDepartmentById 根据id获取部门详情
func (l *GetDepartmentByIdLogic) GetDepartmentById(in *department.DepartmentIdReq) (*department.DepartmentInfo, error) {
	info, err := l.svcCtx.DepartmentModel.GetDepartmentById(in.DepartmentId)
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
