package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"demacia/datacenter/subject/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type AddSubjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSubjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddSubjectLogic {
	return AddSubjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSubjectLogic) AddSubject(req types.AddSubjectReq) (*types.Id, error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	subject, err := l.svcCtx.SubjectModel.GetSubjectByTitleAndOrgId(req.Title, orgId)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Subject Api subject[model] GetSubjectByTitleAndOrgId err:%s", err.Error())
		return nil, errlist.Unknown
	}

	if subject != nil {
		return nil, errlist.SubjectExist
	}
	insSubjectId, err := l.svcCtx.SubjectModel.Insert(&model.Subject{
		Title: req.Title,
		OrgId: orgId,
	})
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Subject Insert subject err:%s", err.Error())
		return nil, errlist.Unknown
	}
	threading.GoSafe(func() {
		subjectGrades := make(model.SubjectGrades, 0)
		for _, v := range req.Grades {
			grade, err := l.svcCtx.ClassRpc.FindGradeById(context.Background(), &class.IdReq{
				Id: v,
			})
			if err != nil {
				l.Logger.Errorf("Subject Api class[rpc] FindGradeById err:%s", err.Error())
				continue
			}
			subjectGrades = append(subjectGrades, &model.SubjectGrade{
				OrgId: orgId,
				// 年级
				GradeId: grade.Id,
				// 学科
				SubjectId:    insSubjectId,
				SubjectTitle: req.Title,
				GradeTitle:   grade.Title,
			})
		}
		err = l.svcCtx.SubjectGradeModel.BatchInsert(subjectGrades)
		if err != nil {
			l.Logger.Errorf("Subject BatchInsert SubjectGrade err:%s", err.Error())
		}
	})
	return &types.Id{Id: insSubjectId}, nil
}
