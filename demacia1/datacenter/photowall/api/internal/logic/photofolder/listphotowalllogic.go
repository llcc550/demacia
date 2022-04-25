package photofolder

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListPhotowallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPhotowallLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListPhotowallLogic {
	return ListPhotowallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPhotowallLogic) ListPhotowall(req types.ListReq) (resp *types.List, err error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	list, total, err := l.svcCtx.PhotoFolderModel.List(req.Title, orgId, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	resp = &types.List{List: []types.Info{}, Total: 0}
	resp.Total = total
	for _, v := range list {
		resp.List = append(resp.List, types.Info{
			Id:    v.Id,
			Title: v.Title,
			Url:   "",
		})
	}
	return resp, nil
}
