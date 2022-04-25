package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"

	"github.com/gogo/protobuf/sortkeys"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"time"
)

type ClassCourseTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClassCourseTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) ClassCourseTableLogic {
	return ClassCourseTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClassCourseTableLogic) ClassCourseTable(req types.ClassIdReq) (*types.ClassCourseTableReply, error) {
	resp := &types.ClassCourseTableReply{ClassCourses: []*types.ClassCourse{}}
	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}

	if req.ClassId == 0 {
		return resp, errlist.InvalidParam
	}

	courseTables, err := l.svcCtx.CourseTableModel.SelectByCid(req.ClassId)
	if err != nil {
		return resp, errlist.Unknown
	}
	if len(courseTables) == 0 {
		return resp, errors.CourseNotFoundErr
	}

	deploys, err := l.svcCtx.CourseTableDeployModel.SelectByOrgId(oid)
	if err != nil {
		if err == sql.ErrNoRows {
			return resp, errors.MustConfig
		}
	}

	var classIds []int64

	classNames := map[int64]string{}

	for _, table := range courseTables {
		classIds = append(classIds, table.ClassId)
		if _, ok := classNames[table.ClassId]; !ok {
			classNames[table.ClassId] = table.ClassName
		}
	}

	sortkeys.Int64s(classIds)

	classIds = removeDup(classIds)

	res := make([]*types.ClassCourse, 0, len(classIds))

	for _, r := range classIds {
		res = append(res, &types.ClassCourse{ClassId: r, ClassName: classNames[r], CourseList: []*types.OrgCourseInfo{}})
	}

	for _, r := range res {
		for _, deploy := range deploys {
			if deploy.CourseFlag == 1 {
				startTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.StartTime)
				endTime, _ := time.Parse("2006-01-02T15:04:05Z", deploy.EndTime)
				r.CourseList = append(r.CourseList, &types.OrgCourseInfo{
					Weekday:    deploy.Weekday,
					CourseSort: deploy.CourseSort,
					StartTime:  startTime.Format("15:04:05"),
					EndTime:    endTime.Format("15:04:05"),
					CourseNote: deploy.Note,
				})
			}
		}
	}

	for _, r := range res {
		for _, table := range courseTables {
			if r.ClassId == table.ClassId {
				for _, info := range r.CourseList {
					if info.CourseSort == table.CourseSort && info.Weekday == table.WeekDay {
						info.TeacherName = table.TeacherName
						info.SubjectName = table.SubjectName
						info.PositionName = table.PositionName
						break
					}
				}
			}
		}
	}
	resp.ClassCourses = res
	return resp, nil
}

func removeDup(a []int64) []int64 {
	i := 0
	for j := 1; j < len(a); j++ {
		if a[i] != a[j] {
			i++
			a[i] = a[j]
		}
	}
	return a[:i+1]
}
