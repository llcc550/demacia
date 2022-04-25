package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/service/event/model"

	"demacia/service/event/api/internal/svc"
	"demacia/service/event/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CategoryInsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) CategoryInsertLogic {
	return CategoryInsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryInsertLogic) CategoryInsert(req types.CategoryInsertReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	//验证名称唯一性
	_, err = l.svcCtx.CategoryModel.FindOneByName(orgId, req.CategoryName)
	if err == nil {
		return errlist.EventCategoryExist
	}

	_, err = l.svcCtx.CategoryModel.Insert(&model.Category{
		OrgId: orgId,
		Name:  req.CategoryName,
	})
	if err != nil {
		return err
	}
	return nil
}
