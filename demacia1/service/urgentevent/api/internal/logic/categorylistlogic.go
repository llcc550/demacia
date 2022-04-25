package logic

import (
	"context"
	"demacia/common/baseauth"

	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) CategoryListLogic {
	return CategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryListLogic) CategoryList(req types.CategoryInsertReq) (resp *types.CategoryListResponse, err error) {
	resp = &types.CategoryListResponse{List: []*types.CategoryList{}}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, err
	}
	list, err := l.svcCtx.CategoryModel.FindList(orgId, req.CategoryName)
	if err != nil {
		return resp, err
	}
	for _, item := range *list {
		resp.List = append(resp.List, &types.CategoryList{
			CategoryId:   item.Id,
			CategoryName: item.Name,
		})
	}
	return resp, nil
}
