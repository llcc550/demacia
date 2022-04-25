package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetSubjectTeacherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSubjectTeacherLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSubjectTeacherLogic {
	return GetSubjectTeacherLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSubjectTeacherLogic) GetSubjectTeacher(req types.Id) (resp *types.ListSubjectTeacherResp, err error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	subjectTeacher, err := l.svcCtx.SubjectTeacherModel.ListSubjectTeacherByOrgIdAndSubject(orgId, req.Id)
	if err != nil {
		l.Logger.Errorf("subject Api SubjectTeacher[Model] ListSubjectTeacherByOrgIdAndSubject err :%s", err.Error())
		return nil, errlist.Unknown
	}
	resp = &types.ListSubjectTeacherResp{List: []types.SubjectTeacherResp{}}

	for _, v := range *subjectTeacher {
		resp.List = append(resp.List, types.SubjectTeacherResp{
			Id:       v.MemberId,
			TrueName: v.TrueName,
		})
	}
	return resp, nil
}
