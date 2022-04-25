package logic

import (
	"context"
	"database/sql"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/class/model"
	"demacia/datacenter/member/rpc/member"
	subject2 "demacia/datacenter/subject/rpc/subject"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddTeachLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddTeachLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddTeachLogic {
	return AddTeachLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTeachLogic) AddTeach(req types.AddTeachReq) (*types.Id, error) {

	class, err := l.svcCtx.ClassModel.GetClassById(req.ClassId)
	if err != nil || class == nil {
		return nil, errlist.ClassNotFound
	}
	subjectInfo, err := l.svcCtx.SubjectRpc.GetSubjectById(l.ctx, &subject2.IdReq{Id: req.SubjectId})
	if err != nil || subjectInfo == nil {
		return nil, errlist.SubjectNotFound
	}
	subject := types.Info{
		Id:    subjectInfo.Id,
		Title: subjectInfo.Title,
	}
	memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: req.TeacherId})
	if err != nil || memberInfo == nil {
		return nil, errlist.MemberNotExist
	}
	isExistSubject, err := l.svcCtx.TeachModel.FindSubjectBySubjectIdAndClassId(req.SubjectId, class.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach select SubjectTeacher exist err:%s", err.Error())
		return nil, errlist.Unknown
	}
	if isExistSubject != nil {
		return nil, errlist.SubjectExist
	}

	isExistTeacher, err := l.svcCtx.TeachModel.FindTeacherByMemberIdAndClassId(req.TeacherId, req.ClassId)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Teach select SubjectTeacher exist err:%s", err.Error())
		return nil, errlist.Unknown
	}
	if isExistTeacher != nil {
		return nil, errlist.TeacherExistTeach
	}

	insertId, err := l.svcCtx.TeachModel.Insert(&model.ClassSubjectTeacher{
		ClassId:      class.Id,
		SubjectId:    subject.Id,
		SubjectTitle: subject.Title,
		MemberId:     memberInfo.Id,
		TrueName:     memberInfo.TrueName,
	})
	if err != nil || insertId < 0 {
		l.Logger.Errorf("Teach insert ClassSubjectTeacher error:%s", err.Error())
		return nil, errlist.Unknown
	}
	return &types.Id{Id: insertId}, nil
}
