package logic

import (
	"context"
	"strings"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/card/rpc/cardclient"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/errors"
	"demacia/datacenter/student/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type UpdateStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateStudentsLogic {
	return UpdateStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateStudentsLogic) UpdateStudents(req types.InsertRequest) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	studentInfo, err := l.svcCtx.StudentModel.FindOneById(req.StudentId)
	if err != nil || studentInfo.OrgId != orgId {
		return errors.StudentNotExist
	}
	classInfo, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
	if err != nil || classInfo.OrgId != orgId {
		return errors.StudentClassNotExist
	}
	err = mr.Finish(func() error {
		//查询用户名
		student, err := l.svcCtx.StudentModel.FindOneByOrgIdAndUserName(orgId, req.UserName)
		if err == nil && student.Id != req.StudentId {
			return errors.StudentExist
		}
		return nil
	}, func() error {
		//判断卡号是否存在
		if len(req.CardNumber) > 0 {
			for _, card := range req.CardNumber {
				result, err := l.svcCtx.CardRpc.CheckStudent(
					l.ctx, &cardclient.CheckReq{
						OrgId:    orgId,
						CardNum:  card,
						ObjectId: req.StudentId,
					})
				if err == nil && result.Result {
					return errors.StudentCardExist
				}
			}
		}
		return nil
	}, func() error {
		//判断身份证号是否存在
		if req.IdNumber != "" {
			idNumber := basefunc.CheckIdCard(req.IdNumber)
			if !idNumber {
				return errlist.IdNumberErr
			}
			_, err = l.svcCtx.StudentModel.FindOneByIdNumber(req.IdNumber)
			if err == nil && studentInfo.Id != req.StudentId {
				return errors.StudentIdNumberExist
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = l.svcCtx.StudentModel.UpdateOne(&model.Student{
		Id:            req.StudentId,
		OrgId:         orgId,
		ClassId:       req.ClassId,
		TrueName:      req.StudentName,
		UserName:      req.UserName,
		Deleted:       false,
		StageId:       classInfo.StageId,
		GradeId:       classInfo.GradeId,
		ClassFullName: classInfo.FullName,
		Password:      basefunc.HashPassword(req.UserName),
		Face:          req.Face,
		Avatar:        req.Avatar,
		FaceStatus:    0,
		IdNumber:      req.IdNumber,
		CardNumber:    strings.Join(req.CardNumber, ","),
		Sex:           req.Sex,
	})
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		sex := int8(0)
		if req.Sex != 0 {
			sex = 1
		}
		s := datacenter.Marshal(datacenter.Student, req.StudentId, datacenter.Update, datacenter.StudentData{
			StudentName:    req.StudentName,
			UserName:       req.UserName,
			OrganizationId: orgId,
			ClassId:        req.ClassId,
			Sex:            sex,
		})
		_ = l.svcCtx.KqPusher.Push(s)
	})
	//插入卡号
	if len(req.CardNumber) > 0 {
		_, err = l.svcCtx.CardRpc.AddStudentCard(l.ctx, &cardclient.AddReq{
			ObjectId: req.StudentId,
			OrgId:    orgId,
			CardNum:  req.CardNumber,
		})
	}
	return nil
}
