package photofolder

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/cachemodel"
	"demacia/common/errlist"
	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"
	"demacia/datacenter/photowall/errors"
	"demacia/datacenter/photowall/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type InsertPhotowallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertPhotowallLogic(ctx context.Context, svcCtx *svc.ServiceContext) InsertPhotowallLogic {
	return InsertPhotowallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertPhotowallLogic) InsertPhotowall(req types.TitleReq) (*types.Id, error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	folder, err := l.svcCtx.PhotoFolderModel.FindOneByName(req.Title, orgId)
	if err != nil {
		if err != cachemodel.ErrNotFound {
			return nil, errlist.Unknown
		}
	}
	if folder != nil {
		return nil, errors.PhotoFolderNameExist
	}
	insertId, err := l.svcCtx.PhotoFolderModel.Insert(&model.PhotoFolder{
		OrgId: orgId,
		Title: req.Title,
	})
	if err != nil {
		return nil, errlist.Unknown
	}
	return &types.Id{Id: insertId}, nil
}
