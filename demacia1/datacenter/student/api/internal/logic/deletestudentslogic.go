package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type DeleteStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteStudentsLogic {
	return DeleteStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteStudentsLogic) DeleteStudents(req types.IdsRequest) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = l.svcCtx.StudentModel.DeleteByStudentIds(req.StudentIds)
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		l.deleteStudents(orgId, req.StudentIds)
	})
	return nil
}

func (l *DeleteStudentsLogic) deleteStudents(orgId int64, studentIds []int64) {
	for _, studentId := range studentIds {
		err := l.updateClassStudentNum(studentId)
		if err != nil {
			logx.Errorf("update class student number err:%s", err.Error())
		}
		err = l.deleteStudentCard(orgId, studentId)
		if err != nil {
			logx.Errorf("delete student card err:%s", err.Error())
		}
		l.pushDeleteStudentActionToDataCenter(studentId)
	}
}

func (l *DeleteStudentsLogic) updateClassStudentNum(studentId int64) error {
	studentInfo, err := l.svcCtx.StudentModel.FindOneById(studentId)
	if err != nil {
		return err
	}
	count, err := l.svcCtx.StudentModel.GetClassStudentCount(studentInfo.ClassId)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.ClassRpc.ChangeStudentNum(l.ctx, &class.StudentNumReq{
		ClassId:    studentInfo.ClassId,
		StudentNum: count,
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *DeleteStudentsLogic) deleteStudentCard(orgId, studentId int64) error {
	_, err := l.svcCtx.CardRpc.AddStudentCard(context.Background(), &card.AddReq{
		OrgId:    orgId,
		ObjectId: studentId,
		CardNum:  []string{},
	})
	return err
}

func (l *DeleteStudentsLogic) pushDeleteStudentActionToDataCenter(studentId int64) {
	s := datacenter.Marshal(datacenter.Student, studentId, datacenter.Delete, nil)
	_ = l.svcCtx.KqPusher.Push(s)
}
