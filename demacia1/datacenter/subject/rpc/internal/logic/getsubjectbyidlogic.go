package logic

import (
	"context"

	"demacia/datacenter/subject/rpc/internal/svc"
	"demacia/datacenter/subject/rpc/subject"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetSubjectByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSubjectByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSubjectByIdLogic {
	return &GetSubjectByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSubjectByIdLogic) GetSubjectById(in *subject.IdReq) (*subject.SubjectInfo, error) {
	// todo: add your logic here and delete this line
	subjectInfo, err := l.svcCtx.SubjectModel.FindSubjectById(in.Id)
	if err != nil {
		return nil, err
	}
	res := &subject.SubjectInfo{
		Id:    subjectInfo.Id,
		Title: subjectInfo.Title,
	}
	return res, nil
}
