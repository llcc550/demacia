package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListLogic {
	return ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req types.ListReq) (resp *types.ListResp, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.AuthLoginFail
	}
	list, err := l.svcCtx.DepartmentModel.GetDepartmentsByOrgIdWithTitle(orgId, req.Title)
	if err != nil {
		return nil, err
	}
	resp = &types.ListResp{List: []*types.ListInfo{}}
	for _, item := range list {
		resp.List = append(resp.List, &types.ListInfo{
			Id:          item.Id,
			Title:       item.Title,
			Sort:        item.Sort,
			MemberCount: item.MemberCount,
		})
	}
	return resp, nil
}
