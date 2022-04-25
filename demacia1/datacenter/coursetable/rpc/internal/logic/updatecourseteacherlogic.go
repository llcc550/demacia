package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/rpc/coursetable"
	"demacia/datacenter/coursetable/rpc/internal/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type UpdateCourseTeacherLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCourseTeacherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCourseTeacherLogic {
	return &UpdateCourseTeacherLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCourseTeacherLogic) UpdateCourseTeacher(in *coursetable.ClassIdAndMemberInfoReq) (*coursetable.SuccessResp, error) {

	if in.MemberId == 0 || in.ClassId == 0 {
		return nil, status.Error(codes.NotFound, errlist.InvalidParam.Error())
	}

	memberCourse, err := l.svcCtx.CourseTableModel.SelectByMemberId(in.MemberId)
	if err != nil {
		return nil, status.Error(codes.Unknown, errlist.Unknown.Error())
	}

	classCourse, err := l.svcCtx.CourseTableModel.SelectByCid(in.ClassId)
	if err != nil {
		return nil, status.Error(codes.Unknown, errlist.Unknown.Error())
	}

	for _, memberCourse := range memberCourse {
		for _, classCourse := range classCourse {
			if memberCourse.WeekDay == classCourse.WeekDay && classCourse.CourseSort == memberCourse.CourseSort {
				return nil, status.Error(codes.Unknown, errors.CourseConflictErr.Error())
			}
		}
	}

	if err := l.svcCtx.CourseTableModel.ChangeTeacherByClassId(in.ClassId, in.MemberId, in.SubjectId, in.Truename); err != nil {
		l.Logger.Errorf("change coursetable member err:%s", err.Error())
		return nil, status.Error(codes.Unknown, errlist.Unknown.Error())
	}

	return &coursetable.SuccessResp{Success: true}, nil
}
