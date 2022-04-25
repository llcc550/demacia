package photofolder

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"

	"demacia/datacenter/photowall/api/internal/svc"
	"demacia/datacenter/photowall/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeletedPhotowallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletedPhotowallLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletedPhotowallLogic {
	return DeletedPhotowallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletedPhotowallLogic) DeletedPhotowall(req types.Id) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.NoAuth
	}
	_ = l.svcCtx.PhotoFolderModel.Delete(req.Id, orgId)
	// todo : 删除相册下的资源
	return nil
}
