package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/model"

	"github.com/tealeg/xlsx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ExportStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) ExportStudentsLogic {
	return ExportStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportStudentsLogic) ExportStudents(req types.ListConditionRequest) (*types.TemplateFile, error) {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	file := xlsx.NewFile()
	dataSheet, _ := file.AddSheet("sheet(1)")
	firstRow := dataSheet.AddRow()
	firstRow.AddCell().Value = "班级名称"
	firstRow.AddCell().Value = "用户名"
	firstRow.AddCell().Value = "姓名"
	firstRow.AddCell().Value = "性别"
	students, count, err := l.svcCtx.StudentModel.FindListByConditions(&model.ListCondition{
		OrgId:       orgId,
		StageId:     req.StageId,
		GradeId:     req.GradeId,
		ClassId:     req.ClassId,
		StudentName: req.StudentName,
		Page:        req.Page,
		Limit:       req.Limit,
		FaceStatus:  req.FaceStatus,
	})
	if err != nil || count == 0 {
		return &types.TemplateFile{
			File: file,
			Name: "学生",
		}, err
	}
	for _, student := range students {
		sex := "男"
		if student.Sex == 0 {
			sex = "女"
		}
		row := dataSheet.AddRow()
		row.AddCell().Value = student.ClassFullName
		row.AddCell().Value = student.UserName
		row.AddCell().Value = student.TrueName
		row.AddCell().Value = sex
	}
	return &types.TemplateFile{
		File: file,
		Name: "学生",
	}, nil
}
