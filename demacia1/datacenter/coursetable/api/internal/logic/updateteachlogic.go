package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/member/rpc/member"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"

	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
)

type UpdateTeachLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTeachLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateTeachLogic {
	return UpdateTeachLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTeachLogic) UpdateTeach(req types.UpdateTeachReq) error {
	clazz, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class Api UpdateTeach err :%s", err.Error())
		return errlist.Unknown
	}
	if clazz == nil {
		return errlist.ClassNotFound
	}
	isExistSubject, err := l.svcCtx.TeachModel.FindSubjectBySubjectIdAndClassId(req.SubjectId, clazz.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach select ClassSubjectTeacher exist err:%s", err.Error())
		return errlist.Unknown
	}
	if isExistSubject == nil {
		return errlist.SubjectNotFound
	}
	memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: req.TeacherId})
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("CLass Api Member[Rpc] FindOneById err:%s", err.Error())
		return errlist.Unknown
	}
	if memberInfo == nil {
		return errlist.MemberNotExist
	}
	isExistTeacher, err := l.svcCtx.TeachModel.FindTeacherByMemberIdAndClassId(req.TeacherId, req.ClassId)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("CLass Api Teach select ClassSubjectTeacher exist err:%s", err.Error())
		return errlist.Unknown
	}
	if isExistTeacher != nil {
		return errlist.TeacherExistTeach
	}

	memberCourse, err := l.svcCtx.CourseTableModel.SelectByMemberId(req.TeacherId)
	if err != nil {
		return errlist.Unknown
	}

	classCourse, err := l.svcCtx.CourseTableModel.SelectByClassIdAndSubjectId(req.ClassId, req.SubjectId)
	if err != nil {
		return errlist.Unknown
	}

	for _, memberCourse := range memberCourse {
		for _, classCourse := range classCourse {
			if memberCourse.WeekDay == classCourse.WeekDay && classCourse.CourseSort == memberCourse.CourseSort {
				return errors.CourseConflictErr
			}
		}
	}

	if err := l.svcCtx.CourseTableModel.ChangeTeacherByClassId(req.ClassId, req.TeacherId, req.SubjectId, memberInfo.TrueName); err != nil {
		l.Logger.Errorf("change coursetable member err:%s", err.Error())
		return errlist.Unknown
	}

	err = l.svcCtx.TeachModel.Update(&model.ClassSubjectTeacher{
		ClassId:   req.ClassId,
		SubjectId: req.SubjectId,
		MemberId:  memberInfo.Id,
		TrueName:  memberInfo.TrueName,
	})
	if err != nil {
		l.Logger.Errorf("CLass Api Teach update ClassSubjectTeacher err:%s", err.Error())
		return errlist.Unknown
	}
	return nil
}
