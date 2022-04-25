package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ClassesByGradeIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type ()

func NewClassesByGradeIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) ClassesByGradeIdLogic {
	return ClassesByGradeIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClassesByGradeIdLogic) ClassesByGradeId(req types.Id) (resp *types.ClassListRespose, err error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.ClassModel.ListByOrgIdAndGradeId(orgId, req.Id)
	if err != nil {
		return nil, err
	}
	resp = &types.ClassListRespose{List: []*types.ClassInfo{}}

	classIds := make([]int64, 0)
	for _, i := range list {
		classIds = append(classIds, i.Id)
	}
	classTeaches, err := l.svcCtx.ClassTeacherModel.GetTeacherWithClassIds(classIds)
	if err != nil {
		return nil, err
	}
	for _, i := range list {
		teachers := make([]*types.Teacher, 0)
		for _, t := range *classTeaches {
			if t.ClassId == i.Id {
				teachers = append(teachers, &types.Teacher{
					Id:       t.Id,
					TrueName: t.TeacherName,
				})
			}
		}
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
	return resp, nil
}
