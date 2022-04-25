package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/class/model"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
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
	class, err := l.svcCtx.ClassModel.GetClassById(req.ClassId)
	if err != nil || class == nil {
		return errlist.ClassNotFound
	}
	isExistSubject, err := l.svcCtx.TeachModel.FindSubjectBySubjectIdAndClassId(req.SubjectId, class.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach select ClassSubjectTeacher exist err:%s", err.Error())
		return errlist.Unknown
	}
	if isExistSubject == nil {
		return errlist.SubjectNotFound
	}
	memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: req.TeacherId})
	if err != nil || memberInfo == nil {
		return errlist.MemberNotExist
	}
	isExistTeacher, err := l.svcCtx.TeachModel.FindTeacherByMemberIdAndClassId(req.TeacherId, req.ClassId)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach select ClassSubjectTeacher exist err:%s", err.Error())
		return errlist.Unknown
	}
	if isExistTeacher != nil {
		return errlist.TeacherExistTeach
	}
	err = l.svcCtx.TeachModel.Update(&model.ClassSubjectTeacher{
		ClassId:   req.ClassId,
		SubjectId: req.SubjectId,
		MemberId:  memberInfo.Id,
		TrueName:  memberInfo.TrueName,
	})
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach update ClassSubjectTeacher err:%s", err.Error())
		return errlist.Unknown
	}
	return nil
}
