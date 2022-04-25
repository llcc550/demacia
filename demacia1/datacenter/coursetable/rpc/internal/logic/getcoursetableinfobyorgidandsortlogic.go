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

type GetCourseTableInfoByOrgIdAndSortLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseTableInfoByOrgIdAndSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseTableInfoByOrgIdAndSortLogic {
	return &GetCourseTableInfoByOrgIdAndSortLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseTableInfoByOrgIdAndSortLogic) GetCourseTableInfoByOrgIdAndSort(in *coursetable.CourseInfoReq) (*coursetable.CourseTableInfoResp, error) {
	if in.OrgId == 0 || in.ClassId == 0 || in.Weekday == 0 || in.CourseSort == 0 {
		return nil, status.Error(codes.InvalidArgument, errlist.InvalidParam.Error())
	}
	course, err := l.svcCtx.CourseTableModel.SelectByClassIdAndWeekdayAndSort(in.OrgId, in.ClassId, int8(in.Weekday), int8(in.CourseSort))
	if err != nil {
		l.Logger.Errorf("select courseInfo err:%s", err.Error())
		return nil, status.Error(codes.NotFound, errors.CourseNotFoundErr.Error())
	}

	return &coursetable.CourseTableInfoResp{
		PositionName: course.PositionName,
		SubjectName:  course.SubjectName,
		ClassName:    course.ClassName,
		TeacherName:  course.TeacherName,
	}, nil
}
