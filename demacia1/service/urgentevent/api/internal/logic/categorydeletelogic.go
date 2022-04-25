package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"

	"demacia/service/urgentevent/api/internal/svc"
	"demacia/service/urgentevent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CategoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) CategoryDeleteLogic {
	return CategoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryDeleteLogic) CategoryDelete(req types.CategoryIdReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	//根据id查询
	category, err := l.svcCtx.CategoryModel.FindOne(req.CategoryId)
	if err != nil {
		return errlist.EventCategoryNotExist
	}
	if category.OrgId != orgId {
		return errlist.NoAuth
	}
	err = l.svcCtx.CategoryModel.Delete(req.CategoryId)
	if err != nil {
		return err
	}
	return nil
}
