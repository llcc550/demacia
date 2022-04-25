package logic

import (
	"context"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/coursetable/rpc/coursetable"
	"demacia/datacenter/coursetable/rpc/internal/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetCourseTableRecordInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseTableRecordInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseTableRecordInfoLogic {
	return &GetCourseTableRecordInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseTableRecordInfoLogic) GetCourseTableRecordInfo(in *coursetable.OrgIdAndClassIdReq) (*coursetable.CourseTableRecordResp, error) {
	var resp coursetable.CourseTableRecordResp

	courses, err := l.svcCtx.CourseTableModel.SelectCourseTableByOrgId(in.OrgId)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.CourseNotFoundErr.Error())
	}

	deploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(in.OrgId)
	if err != nil {
		return nil, status.Error(codes.NotFound, errors.DeployNotFoundErr.Error())
	}
	weekday := int8(time.Now().Weekday())
	if weekday == 0 {
		weekday = 7
	}
	var deploysVal model.CourseTableDeploys
	for _, deploy := range deploys {
		if deploy.Weekday == weekday {
			deploysVal = append(deploysVal, &model.CourseTableDeploy{
				Id:         deploy.Id,
				OrgId:      deploy.OrgId,
				CourseSort: deploy.CourseSort,
				Weekday:    deploy.Weekday,
				Note:       deploy.Note,
				Grouping:   deploy.Grouping,
				StartTime:  deploy.StartTime,
				EndTime:    deploy.EndTime,
				CourseFlag: deploy.CourseFlag,
			})
		}
	}

	for _, course := range courses {
		if course.ClassId == in.ClassId && course.WeekDay == weekday {
			for _, deploy := range deploysVal {
				if course.StartTime == deploy.StartTime {
					startTime, _ := time.Parse("2006-01-02T15:04:05Z", course.StartTime)
					endTime, _ := time.Parse("2006-01-02T15:04:05Z", course.EndTime)
					note := deploy.Note
					if deploy.Grouping == 2 || deploy.Grouping == 5 {
						note = "第" + note
					}
					resp.List = append(resp.List, &coursetable.CourseTableRecord{
						StartTime:    startTime.Format("15:04:05"),
						EndTime:      endTime.Format("15:04:05"),
						SubjectName:  course.SubjectName,
						PositionName: course.PositionName,
						CourseNote:   note,
						ClassName:    course.ClassName,
						OrgId:        course.OrganizationId,
						MemberId:     course.TeacherId,
						MemberName:   course.TeacherName,
					})
				}
			}
		}
	}

	return &resp, nil
}
