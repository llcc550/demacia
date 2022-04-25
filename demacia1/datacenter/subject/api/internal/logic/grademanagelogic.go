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
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GradeManageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGradeManageLogic(ctx context.Context, svcCtx *svc.ServiceContext) GradeManageLogic {
	return GradeManageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GradeManageLogic) GradeManage(req types.GradeManageReq) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	subject, err := l.svcCtx.SubjectModel.FindSubjectById(req.Id)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("GradeManage Select subject err:%s", err.Error())
		return errlist.Unknown
	}
	threading.GoSafe(func() {
		subjectGrades := make(model.SubjectGrades, 0)
		for _, v := range req.Grades {
			grade, err := l.svcCtx.ClassRpc.FindGradeById(context.Background(), &class.IdReq{
				Id: v,
			})
			if err != nil {
				l.Logger.Errorf("GradeManage FindGradRpc err:%s", err.Error())
				continue
			}
			subjectGrades = append(subjectGrades, &model.SubjectGrade{
				OrgId: orgId,
				// 年级
				GradeId: grade.Id,
				// 学科
				SubjectId:    subject.Id,
				SubjectTitle: subject.Title,
				GradeTitle:   grade.Title,
			})
		}
		err = l.svcCtx.SubjectGradeModel.Deleted(&model.SubjectGrade{
			OrgId:     orgId,
			SubjectId: req.Id,
		})
		if err != nil && err != sql.ErrNoRows {
			l.Logger.Errorf("GradeManage deleted err:%s", err.Error())
			return
		}
		err = l.svcCtx.SubjectGradeModel.BatchInsert(subjectGrades)
		if err != nil && err != sql.ErrNoRows {
			l.Logger.Errorf("GradeManage BatchInsert SubjectGrade err:%s", err.Error())
			return
		}
	})
	return nil
}
