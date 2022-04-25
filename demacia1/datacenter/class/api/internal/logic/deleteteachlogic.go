package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type DeleteTeachLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTeachLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteTeachLogic {
	return DeleteTeachLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTeachLogic) DeleteTeach(req types.DeleteTeachReq) error {
	isExistSubject, err := l.svcCtx.TeachModel.FindSubjectBySubjectIdAndClassId(req.SubjectId, req.ClassId)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach select Subject Exist err:%s", err.Error())
		return errlist.Unknown
	}
	if isExistSubject == nil {
		return errlist.SubjectNotFound
	}
	err = l.svcCtx.TeachModel.DeleteBySubjectIdAndClassId(req.SubjectId, req.ClassId)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach Deleted Teach err:%s", err.Error())
		return errlist.Unknown
	}
	return nil
}
