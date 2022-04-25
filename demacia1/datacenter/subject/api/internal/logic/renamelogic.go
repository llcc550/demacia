package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"demacia/datacenter/subject/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type RenameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRenameLogic(ctx context.Context, svcCtx *svc.ServiceContext) RenameLogic {
	return RenameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RenameLogic) Rename(req types.Info) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.NoAuth
	}
	err = l.svcCtx.SubjectModel.Rename(&model.Subject{
		Id:    req.Id,
		Title: req.Title,
		OrgId: orgId,
	})
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Subject Rename err:%s", err.Error())
		return errlist.Unknown
	}
	_, _ = l.svcCtx.DataBusRpc.Update(l.ctx, &databus.Req{
		Topic:    datacenter.Subject,
		ObjectId: req.Id,
	})

	err = l.svcCtx.SubjectGradeModel.RenameSubjectBySubjectIdAndOrgId(orgId, req.Id, req.Title)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Subject Rename err:%s", err.Error())
		return errlist.Unknown
	}
	return nil
}
