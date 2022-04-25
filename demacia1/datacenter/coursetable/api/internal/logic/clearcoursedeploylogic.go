package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ClearCourseDeployLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearCourseDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) ClearCourseDeployLogic {
	return ClearCourseDeployLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearCourseDeployLogic) ClearCourseDeploy() (*types.BoolReply, error) {
	resp := &types.BoolReply{}
	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}
	if err := l.svcCtx.CourseTableDeployModel.DeleteByOrgId(oid); err != nil {
		return nil, errlist.Unknown
	}
	if err := l.svcCtx.CourseTableModel.DeleteByOrgId(oid); err != nil {
		return nil, errlist.Unknown
	}
	resp.Success = true
	return resp, nil
}
