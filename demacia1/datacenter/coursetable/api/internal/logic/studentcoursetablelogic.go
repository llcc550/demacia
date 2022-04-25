package logic

import (
	"context"
	"demacia/datacenter/coursetable/api/internal/svc"
	"demacia/datacenter/coursetable/api/internal/types"
	"fmt"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type StudentCourseTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentCourseTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) StudentCourseTableLogic {
	return StudentCourseTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StudentCourseTableLogic) StudentCourseTable(req types.PageReq) (*types.StudentTableReply, error) {

	var resp types.StudentTableReply

	//todo 获取学生列表
	for i := 1; i <= 10; i++ {
		resp.StudentCourseTables = append(resp.StudentCourseTables, &types.StudentCourseTable{StudentName: fmt.Sprintf("学生%d", i), ClassName: fmt.Sprintf("班级%d", i), Username: fmt.Sprintf("用户名%d", i)})
	}

	return nil, nil
}
