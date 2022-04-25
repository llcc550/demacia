package logic

import (
	"context"
	"strings"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/student/api/internal/svc"
	"demacia/datacenter/student/api/internal/types"
	"demacia/datacenter/student/errors"
	"demacia/datacenter/student/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type InsertStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) InsertStudentsLogic {
	return InsertStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertStudentsLogic) InsertStudents(req types.InsertRequest) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	classInfo, err := l.svcCtx.ClassRpc.GetClassInfoById(l.ctx, &class.IdReq{Id: req.ClassId})
	if err != nil || classInfo.OrgId != req.ClassId {
		return errors.StudentClassNotExist
	}

	err = mr.Finish(func() error {
		//查询用户名
		_, err := l.svcCtx.StudentModel.FindOneByOrgIdAndUserName(orgId, req.UserName)
		if err == nil {
			return errors.StudentExist
		}
		return nil
	}, func() error {
		//判断卡号是否存在
		if len(req.CardNumber) == 0 {
			return nil
		}
		for _, cardNumber := range req.CardNumber {
			result, err := l.svcCtx.CardRpc.CheckStudent(context.Background(), &card.CheckReq{
				OrgId:    orgId,
				CardNum:  cardNumber,
				ObjectId: 0,
			})
			if err == nil && result.Result {
				return errors.StudentCardExist
			}
		}
		return nil
	}, func() error {
		//判断身份证号是否存在
		if req.IdNumber != "" {
			return nil
		}
		idNumber := basefunc.CheckIdCard(req.IdNumber)
		if !idNumber {
			return errlist.IdNumberErr
		}
		_, err = l.svcCtx.StudentModel.FindOneByIdNumber(req.IdNumber)
		if err == nil {
			return errors.StudentIdNumberExist
		}
		return nil
	})

	if err != nil {
		return err
	}
	studentId, err := l.svcCtx.StudentModel.InsertOne(&model.Student{
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
		sex := int8(1)
		if req.Sex != 1 {
			sex = 0
		}
		s := datacenter.Marshal(datacenter.Student, studentId, datacenter.Add, datacenter.StudentData{
			StudentName:    req.StudentName,
			UserName:       req.UserName,
			OrganizationId: orgId,
			ClassId:        req.ClassId,
			Sex:            sex,
		})
		_ = l.svcCtx.KqPusher.Push(s)
	})
	count, err := l.svcCtx.StudentModel.GetClassStudentCount(req.ClassId)
	if err == nil {
		_, err := l.svcCtx.ClassRpc.ChangeStudentNum(l.ctx, &class.StudentNumReq{
			ClassId:    req.ClassId,
			StudentNum: count,
		})
		if err != nil {
			logx.Errorf("update class student err:%s", err.Error())
		}
	}
	//插入卡号
	if len(req.CardNumber) > 0 {
		_, err = l.svcCtx.CardRpc.AddStudentCard(l.ctx, &card.AddReq{
			ObjectId: studentId,
			OrgId:    orgId,
			CardNum:  req.CardNumber,
		})
	}
	return nil
}
