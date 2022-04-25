package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/subject/api/internal/svc"
	"demacia/datacenter/subject/api/internal/types"
	"demacia/datacenter/subject/model"
	"strings"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type SubjectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubjectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) SubjectListLogic {
	return SubjectListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubjectListLogic) SubjectList(req types.TitleReq) (resp *types.SubjectListsResp, err error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, errlist.NoAuth
	}
	subjectList := &model.Subjects{}
	resp = &types.SubjectListsResp{List: []types.SubjectListResp{}}
	if req.GradeId > 0 {
		subjectGrades, err := l.svcCtx.SubjectGradeModel.ListSubjectByGradeId(req.GradeId)
		if err != nil {
			l.Logger.Errorf("Subject Api SubjectGrade[model] ListSubjectByGradeId err:%s", err.Error())
			return nil, errlist.Unknown
		}
		for _, subjectGrade := range *subjectGrades {
			resp.List = append(resp.List, types.SubjectListResp{
				Id:     subjectGrade.SubjectId,
				Title:  subjectGrade.SubjectTitle,
				Grades: []types.Grade{},
			})
		}
		return resp, nil
	}
	if strings.Trim(req.Title, " ") != "" {
		subjectList, err = l.svcCtx.SubjectModel.ListSubjectByTitleAndOrgId(req.Title, orgId)
		if err != nil {
			l.Logger.Errorf("Subject Api Subject[model] ListSubjectByTitleAndOrgId err:%s", err.Error())
			return nil, errlist.Unknown
		}
	} else {
		subjectList, err = l.svcCtx.SubjectModel.ListSubjectByOrgId(orgId)
		if err != nil {
			l.Logger.Errorf("Subject Api Subject[model] ListSubjectByOrgId err:%s", err.Error())
			return nil, errlist.Unknown
		}
	}

	subjectGradeList, err := l.svcCtx.SubjectGradeModel.ListSubjectGradeByOrgId(orgId)
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("Subject Select SubjectGrade err:%s", err.Error())
		return nil, errlist.Unknown
	}

	subjectMap := map[int64]model.Subject{}
	subjectGradeMap := map[int64][]types.Grade{}
	for _, subject := range *subjectList {
		subjectMap[subject.Id] = model.Subject{
			Id:    subject.Id,
			Title: subject.Title,
		}
	}
	for _, subjectGrade := range *subjectGradeList {
		if _, ok := subjectMap[subjectGrade.SubjectId]; ok {
			subjectGradeMap[subjectGrade.SubjectId] = append(subjectGradeMap[subjectGrade.SubjectId], types.Grade{
				Id:    subjectGrade.GradeId,
				Title: subjectGrade.GradeTitle,
			})
		}

	}
	for _, v := range subjectMap {
		grades := make([]types.Grade, 0)
		if _, ok := subjectGradeMap[v.Id]; ok {
			grades = subjectGradeMap[v.Id]
		}
		resp.List = append(resp.List, types.SubjectListResp{
			Id:     v.Id,
			Title:  v.Title,
			Grades: grades,
		})
	}
	return resp, nil
}
