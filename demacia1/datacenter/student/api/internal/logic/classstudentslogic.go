package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/errors"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ClassStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClassStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ClassStudentsLogic {
	return ClassStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClassStudentsLogic) ClassStudents(req types.ClassIdRequest) (*types.ListResponse, error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	res := types.ListResponse{
		List:  []*types.List{},
		Count: 0,
	}
	classInfo, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
	if err != nil || classInfo.OrgId != orgId {
		return nil, errors.StudentClassNotExist
	}
	students, err := l.svcCtx.StudentModel.FindListByClassId(req.ClassId)
	if err != nil {
		return nil, err
	}
	for _, student := range students {
		res.List = append(res.List, &types.List{
			StudentId:   student.Id,
			StudentName: student.TrueName,
			ClassName:   student.ClassFullName,
			UserName:    student.UserName,
			Sex:         student.Sex,
			Face:        student.Face,
		})
	}
	res.Count = len(students)
	return &res, nil
}
