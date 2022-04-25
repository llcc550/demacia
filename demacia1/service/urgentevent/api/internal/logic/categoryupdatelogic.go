package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/service/urgentevent/model"

	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CategoryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) CategoryUpdateLogic {
	return CategoryUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryUpdateLogic) CategoryUpdate(req types.CategoryUpdateReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	//验证id
	category, err := l.svcCtx.CategoryModel.FindOne(req.CategoryId)
	if (err == nil && category.OrgId != orgId) || err != nil {
		return errlist.EventCategoryNotExist
	}
	//验证名称唯一性
	categoryInfo, err := l.svcCtx.CategoryModel.FindOneByName(orgId, req.CategoryName)
	if err == nil && categoryInfo.Id != req.CategoryId {
		return errlist.EventCategoryExist
	}
	err = l.svcCtx.CategoryModel.Update(&model.Category{
		Id:    req.CategoryId,
		OrgId: 0,
		Name:  req.CategoryName,
	})
	if err != nil {
		return err
	}
	return nil
}
