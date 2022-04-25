package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/rpc/coursetable"
	"demacia/datacenter/coursetable/rpc/internal/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type RemoveCourseTeacherLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveCourseTeacherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCourseTeacherLogic {
	return &RemoveCourseTeacherLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveCourseTeacherLogic) RemoveCourseTeacher(in *coursetable.ClassIdAndMemberIdReq) (*coursetable.SuccessResp, error) {

	if in.MemberId == 0 || in.ClassId == 0 {
		return nil, status.Error(codes.NotFound, errlist.InvalidParam.Error())
	}

	if err := l.svcCtx.CourseTableModel.UnBindTeacherByClassIdAndMemberId(in.ClassId, in.MemberId); err != nil {
		return nil, status.Error(codes.Unknown, errlist.Unknown.Error())
	}

	return &coursetable.SuccessResp{Success: true}, nil
}
