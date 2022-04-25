package logic

import (
	"context"
	"demacia/cloudscreen/courserecord/api/internal/svc"
	"demacia/cloudscreen/courserecord/api/internal/types"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/classclient"
	"fmt"
	"math"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type CourseRecordCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseRecordCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) CourseRecordCountLogic {
	return CourseRecordCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseRecordCountLogic) CourseRecordCount(req types.CourseRecordCountReq) (*types.StudentRecordReply, error) {

	resp := &types.StudentRecordReply{StudentRecords: []*types.StudentRecord{}}

	var classIds []int64

	oid, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return resp, errlist.NoAuth
	}

	classList, err := l.svcCtx.ClassRpc.ListByOrgId(l.ctx, &classclient.OrgIdReq{OrgId: oid})
	if err != nil {
		return resp, errlist.Unknown
	}

	if req.StartDate == "" && req.EndDate != "" || req.StartDate != "" && req.EndDate == "" {
		return resp, errlist.Unknown
	}

	if req.Type == 0 {
		for _, info := range classList.List {
			classIds = append(classIds, info.Id)
		}
	} else if req.Type == 1 {
		if req.StageId == 0 {
			return resp, errlist.InvalidParam
		}
		for _, info := range classList.List {
			if info.StageId == req.StageId {
				classIds = append(classIds, info.Id)
			}
		}
	} else if req.Type == 2 {
		if req.GradeId == 0 {
			return resp, errlist.InvalidParam
		}
		for _, info := range classList.List {
			if info.GradeId == req.GradeId {
				classIds = append(classIds, info.Id)
			}
		}
	} else if req.Type == 3 {
		if req.ClassId == 0 {
			return resp, errlist.InvalidParam
		}
		for _, info := range classList.List {
			if info.Id == req.ClassId {
				classIds = append(classIds, info.Id)
			}
		}
	}
	if len(classIds) >= 0 {
		studentCounts, count, err := l.svcCtx.CourseRecordCountModel.SelectByClassIds(classIds, req.StartDate, req.EndDate, req.StudentName, req.Page, req.Limit)
		if err != nil {
			return resp, errlist.Unknown
		}
		var allShouldCount int
		var allNormalCount int
		var allLateCount int
		for _, studentCount := range studentCounts {
			for _, class := range classList.List {
				if studentCount.ClassId == class.Id {
					allShouldCount += studentCount.ShouldCount
					allNormalCount += studentCount.NormalCount
					allLateCount += studentCount.LateCount
					resp.StudentRecords = append(resp.StudentRecords, &types.StudentRecord{
						ClassName:   class.FullName,
						StudentName: studentCount.Truename,
						Attendance:  fmt.Sprintf("%.0f%s", math.Floor(float64(studentCount.NormalCount)/float64(studentCount.ShouldCount)*100), "%"),
						ShouldCount: studentCount.ShouldCount,
						NormalCount: studentCount.NormalCount,
						LateCount:   studentCount.LateCount,
						LackCount:   studentCount.ShouldCount - studentCount.NormalCount - studentCount.LateCount,
					})
				}
			}
		}
		resp.Count = count
		resp.AllShouldCount = allShouldCount
		resp.AllNormalCount = allNormalCount
		resp.AllLateCount = allLateCount
		resp.AllLackCount = allShouldCount - allNormalCount - allLateCount
		if resp.AllShouldCount > 0 {
			resp.AllAttendance = fmt.Sprintf("%.0f%s", math.Floor(float64(allNormalCount)/float64(allShouldCount)*100), "%")
		}
	}

	return resp, nil
}
