package logic

import (
	"context"

	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChangeStudentNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeStudentNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeStudentNumLogic {
	return &ChangeStudentNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeStudentNumLogic) ChangeStudentNum(in *class.StudentNumReq) (*class.NullResp, error) {
	// todo: add your logic here and delete this line
	classInfo, err := l.svcCtx.ClassModel.GetClassById(in.ClassId)
	if err != nil {
		return nil, status.Error(codes.NotFound, errlist.ClassNotFound.Error())
	}
	err = l.svcCtx.ClassModel.ChangeStudentNumByClassId(classInfo, in.StudentNum)
	if err != nil {
		return nil, err
	}
	total, err := l.svcCtx.ClassModel.GetStudentTotalNumByGradeId(classInfo.GradeId)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.GradeModel.ChangeStudentNumById(classInfo.GradeId, total)
	if err != nil {
		return nil, err
	}
	return &class.NullResp{}, nil
}
