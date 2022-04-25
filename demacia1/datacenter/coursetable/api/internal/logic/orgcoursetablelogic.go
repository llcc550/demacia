package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/classclient"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"time"
)

type OrgCourseTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrgCourseTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) OrgCourseTableLogic {
	return OrgCourseTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrgCourseTableLogic) OrgCourseTable() (*types.OrgCourseTableReply, error) {

	resp := &types.OrgCourseTableReply{OrgCourseTable: []*types.OrgCourseTable{}}

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.InvalidParam
	}

	courseTables, err := l.svcCtx.CourseTableModel.SelectCourseTableByOrgId(oid)
	if err != nil {
		return resp, errlist.Unknown
	}

	deploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(oid)
	if err != nil {
		return resp, errors.MustConfig
	}

	classes, err := l.svcCtx.ClassRpc.ListByOrgId(l.ctx, &classclient.OrgIdReq{OrgId: oid})
	if err != nil {
		return resp, errors.ClassNotFoundErr
	}

	for _, deploy := range deploys {
		if deploy.CourseFlag == 1 {
			resp.OrgCourseTable = append(resp.OrgCourseTable, &types.OrgCourseTable{
				Weekday:    deploy.Weekday,
				CourseSort: deploy.CourseSort,
				Note:       deploy.Note,
				CourseInfo: []*types.CourseInfo{},
			})
		}
	}

	for _, courseTable := range resp.OrgCourseTable {
		for _, deploy := range deploys {
			if courseTable.Weekday == deploy.Weekday && courseTable.CourseSort == deploy.CourseSort {
				for _, info := range classes.List {
					startTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.StartTime)
					endTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.EndTime)
					courseTable.CourseInfo = append(courseTable.CourseInfo, &types.CourseInfo{
						ClassId:   info.Id,
						ClassName: info.FullName,
						StartTime: startTime.Format("15:04:05"),
						EndTime:   endTime.Format("15:04:05"),
					})
				}
			}
		}
	}

	for _, table := range resp.OrgCourseTable {
		for _, courseTable := range courseTables {
			if table.CourseSort == courseTable.CourseSort && table.Weekday == courseTable.WeekDay {
				for _, info := range table.CourseInfo {
					if info.ClassId == courseTable.ClassId {
						info.TeacherName = courseTable.TeacherName
						info.PositionName = courseTable.PositionName
						info.SubjectName = courseTable.SubjectName
						break
					}
				}
			}
		}
	}

	return resp, nil
}
