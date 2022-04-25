package logic

import (
	"context"

	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type StudentsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) StudentsListLogic {
	return StudentsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StudentsListLogic) StudentsList(req types.ListConditionRequest) (resp *types.ListResponse, err error) {
	resp = &types.ListResponse{
		List:  []*types.List{},
		Count: 0,
	}
	students, count, err := l.svcCtx.StudentModel.FindListByConditions(&model.ListCondition{
		OrgId:       req.OrgId,
		StageId:     req.StageId,
		GradeId:     req.GradeId,
		ClassId:     req.ClassId,
		StudentName: req.StudentName,
		Page:        req.Page,
		Limit:       req.Limit,
		FaceStatus:  req.FaceStatus,
	})
	if err != nil || count == 0 {
		return resp, err
	}
	for _, student := range students {
		resp.List = append(resp.List, &types.List{
			StudentId:   student.Id,
			StudentName: student.TrueName,
			ClassName:   student.ClassFullName,
			UserName:    student.UserName,
			Sex:         student.Sex,
			Face:        student.Face,
		})
	}
	resp.Count = count
	return resp, nil
}
