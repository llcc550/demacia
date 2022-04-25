package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type MemberByDepartmentIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMemberByDepartmentIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) MemberByDepartmentIdLogic {
	return MemberByDepartmentIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberByDepartmentIdLogic) MemberByDepartmentId(req types.MemberListReq) (resp *types.MemberListResp, err error) {
	if req.DepartmentId <= 0 || req.Page <= 0 || req.Limit <= 0 {
		return nil, errlist.InvalidParam
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.AuthLoginFail
	}
	departmentInfo, err := l.svcCtx.DepartmentModel.GetDepartmentById(req.DepartmentId)
	if err != nil || departmentInfo.OrgId != orgId {
		return nil, errlist.DepartmentNotExit
	}
	list, err := l.svcCtx.DepartmentMemberModel.GetMembersByOrgIdAndDepartmentId(orgId, req.DepartmentId)
	if err != nil {
		return nil, err
	}
	resp = &types.MemberListResp{
		Count: int64(len(list)),
		List:  []*types.MemberInfo{},
	}
	begin, end := basefunc.PageLimit(len(list), req.Page, req.Limit)
	if end == 0 {
		return
	}
	for _, item := range list[begin:end] {
		resp.List = append(resp.List, &types.MemberInfo{
			MemberId: item.MemberId,
			TrueName: item.TrueName,
			Mobile:   item.Mobile,
		})
	}
	return
}
