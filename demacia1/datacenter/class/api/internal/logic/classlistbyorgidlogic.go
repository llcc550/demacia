package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"strings"
)

type ClassListByOrgIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClassListByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) ClassListByOrgIdLogic {
	return ClassListByOrgIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClassListByOrgIdLogic) ClassListByOrgId(req types.PageListReq) (resp *types.ClassListRespose, err error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	var stageId, gradeId, classId int64
	var teacherName string
	if req.StageId > 0 {
		stageId = req.StageId
		if req.GradeId > 0 {
			gradeId = req.GradeId
			if req.ClassId > 0 {
				classId = req.ClassId
			}
		}
	}
	if strings.Trim(req.TeacherName, " ") != "" {
		teacherName = strings.Trim(req.TeacherName, "")
	}
	list, total, err := l.svcCtx.ClassModel.PageListByOrgId(orgId, stageId, gradeId, classId, teacherName, req.Page, req.Limit)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class Select List err:%s", err.Error())
		return nil, errlist.Unknown
	}
	resp = &types.ClassListRespose{
		List:  []*types.ClassInfo{},
		Total: total,
		Limit: req.Limit,
	}

	classIds := make([]int64, 0)
	for _, i := range list {
		classIds = append(classIds, i.Id)
	}
	//classTeaches, err := l.svcCtx.ClassTeacherModel.GetTeacherWithClassIds(classIds)
	//if err != nil {
	//	return nil, err
	//}
	for _, i := range list {
		teachers := make([]*types.Teacher, 0)
		if i.Teachers != "" {
			teacherArr := strings.Split(i.Teachers, ",")
			for _, t := range teacherArr {
				teachers = append(teachers, &types.Teacher{
					Id:       0,
					TrueName: t,
				})
			}
		}

		//for _, t := range *classTeaches {
		//	if t.ClassId == i.Id {
		//		teachers = append(teachers, &types.Teacher{
		//			Id:       t.TeacherId,
		//			TrueName: t.TeacherName,
		//		})
		//	}
		//}
		resp.List = append(resp.List, &types.ClassInfo{
			Id:            i.Id,
			StageId:       i.StageId,
			GradeId:       i.GradeId,
			StageTitle:    i.StageTitle,
			FullName:      i.FullName,
			AliasName:     i.AliasName,
			Class_teacher: teachers,
			Desc:          i.Desc,
			MemberNum:     i.ClassMemberNum,
		})
	}
	return
}
