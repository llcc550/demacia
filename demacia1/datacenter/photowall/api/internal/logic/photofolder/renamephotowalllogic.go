package photofolder

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/cachemodel"
	"demacia/common/errlist"
	"demacia/datacenter/photowall/errors"
	"demacia/datacenter/photowall/model"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type RenamePhotowallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenamePhotowallLogic(ctx context.Context, svcCtx *svc.ServiceContext) RenamePhotowallLogic {
	return RenamePhotowallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenamePhotowallLogic) RenamePhotowall(req types.IdTitleReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.NoAuth
	}
	folder, err := l.svcCtx.PhotoFolderModel.FindOneByName(req.Title, orgId)
	if err != nil {
		if err != cachemodel.ErrNotFound {
			return errlist.Unknown
		}
	}
	if folder != nil {
		if folder.Id != req.Id {
			return errors.PhotoFolderNameExist
		}
	}
	_ = l.svcCtx.PhotoFolderModel.Rename(&model.PhotoFolder{
		Id:    req.Id,
		OrgId: orgId,
		Title: req.Title,
	})
	return nil
}
