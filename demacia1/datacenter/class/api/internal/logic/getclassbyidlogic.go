package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"

	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetClassByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClassByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetClassByIdLogic {
	return GetClassByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClassByIdLogic) GetClassById(req types.Id) (*types.ClassInfo, error) {
	class, err := l.svcCtx.ClassModel.GetClassById(req.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class select class err:%s", err.Error())
		return nil, errlist.Unknown
	}
	if class == nil {
		return nil, errlist.ClassNotFound
	}
	teachers, err := l.svcCtx.ClassTeacherModel.GetTeacherByClassId(req.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class select classTeacher err:%s", err.Error())
		return nil, errlist.Unknown
	}
	classTeacher := make([]*types.Teacher, 0)
	for _, v := range *teachers {
		classTeacher = append(classTeacher, &types.Teacher{
			Id:       v.TeacherId,
			TrueName: v.TeacherName,
		})
	}

	resp := types.ClassInfo{
		Id:            class.Id,
		StageId:       class.StageId,
		GradeId:       class.GradeId,
		StageTitle:    class.StageTitle,
		FullName:      class.FullName,
		AliasName:     class.AliasName,
		Class_teacher: classTeacher,
		Desc:          class.Desc,
		MemberNum:     class.ClassMemberNum,
	}

	return &resp, nil
}
