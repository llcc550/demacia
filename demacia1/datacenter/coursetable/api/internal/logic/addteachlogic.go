package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"demacia/datacenter/coursetable/errors"
	"demacia/datacenter/coursetable/model"
	"demacia/datacenter/member/rpc/member"
	"demacia/datacenter/subject/rpc/subject"

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
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	clazz, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class Api addTeach class[model] GetClassById err:%s", err.Error())
		return nil, errlist.Unknown
	}
	if clazz == nil {
		return nil, errlist.ClassNotFound
	}
	subjectInfo, err := l.svcCtx.SubjectRpc.GetSubjectById(l.ctx, &subject.IdReq{Id: req.SubjectId})
	if err != nil {
		l.Logger.Errorf("Class Api addTeach Subject[Rpc] GetSubjectById err:%s", err.Error())
		return nil, errlist.Unknown
	}
	if subjectInfo == nil {
		return nil, errlist.SubjectNotFound
	}
	subInfo := types.Info{
		Id:    subjectInfo.Id,
		Title: subjectInfo.Title,
	}
	memberInfo, err := l.svcCtx.MemberRpc.FindOneById(l.ctx, &member.IdReq{Id: req.TeacherId})
	if err != nil {
		l.Logger.Errorf("Class Api addTeach Member[Rpc] FindOneById err:%s", err.Error())
		return nil, errlist.Unknown
	}
	if memberInfo == nil {
		return nil, errlist.MemberNotExist
	}
	isExistSubject, err := l.svcCtx.TeachModel.FindSubjectBySubjectIdAndClassId(req.SubjectId, clazz.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class Api Teach select SubjectTeacher exist err:%s", err.Error())
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

	memberCourse, err := l.svcCtx.CourseTableModel.SelectByMemberId(req.TeacherId)
	if err != nil {
		return nil, errlist.Unknown
	}

	classCourse, err := l.svcCtx.CourseTableModel.SelectByClassIdAndSubjectId(req.ClassId, req.SubjectId)
	if err != nil {
		return nil, errlist.Unknown
	}

	for _, memberCourse := range memberCourse {
		for _, classCourse := range classCourse {
			if memberCourse.WeekDay == classCourse.WeekDay && classCourse.CourseSort == memberCourse.CourseSort {
				return nil, errors.CourseConflictErr
			}
		}
	}

	if err := l.svcCtx.CourseTableModel.ChangeTeacherByClassId(req.ClassId, req.TeacherId, req.SubjectId, memberInfo.TrueName); err != nil {
		l.Logger.Errorf("change coursetable member err:%s", err.Error())
		return nil, errlist.Unknown
	}

	insertId, err := l.svcCtx.TeachModel.Insert(&model.ClassSubjectTeacher{
		ClassId:      clazz.Id,
		SubjectId:    subInfo.Id,
		SubjectTitle: subInfo.Title,
		MemberId:     memberInfo.Id,
		TrueName:     memberInfo.TrueName,
		OrgId:        orgId,
	})
	if err != nil || insertId < 0 {
		l.Logger.Errorf("Teach insert ClassSubjectTeacher error:%s", err.Error())
		return nil, errlist.Unknown
	}
	return &types.Id{Id: insertId}, nil
}
