package logic

import (
	"context"
	"demacia/common/datacenter"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/errlist"
	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/parent/api/internal/svc"
	"demacia/datacenter/parent/api/internal/types"
	"demacia/datacenter/parent/model"
	"demacia/datacenter/student/rpc/student"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/mr"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type InsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) InsertLogic {
	return InsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertLogic) Insert(req types.InsertRequest) error {
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	err = mr.Finish(func() error {
		//查询用户名
		mobile := basefunc.CheckMobile(req.Moblie)
		if !mobile {
			return errlist.MobileErr
		}
		_, err = l.svcCtx.ParentModel.FindOneByMobile(req.Moblie)
		if err == nil {
			return errlist.MobileExist
		}
		return nil
	},
		func() error {
			//判断身份证号是否存在
			if req.IdNumber != "" {
				idNumber := basefunc.CheckIdCard(req.IdNumber)
				if !idNumber {
					return errlist.IdNumberErr
				}
				_, err = l.svcCtx.ParentModel.FindOneByIdNumber(req.IdNumber)
				if err == nil {
					return errlist.StudentIdNumberExist
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	parentId, err := l.svcCtx.ParentModel.InsertOne(&model.Parent{
		TrueName:   req.ParentName,
		UserName:   req.Moblie,
		Password:   basefunc.HashPassword(req.Moblie),
		Face:       req.Face,
		FaceStatus: 0,
		IdNumber:   req.IdNumber,
		Address:    req.Address,
		Mobile:     req.Moblie,
		Pinyin:     basefunc.GetPinyin(req.ParentName),
	})
	if err != nil {
		return err
	}
	threading.GoSafe(func() {
		datacenterStudentInfoSlice := make([]datacenter.StudentInfo, 0)
		for _, v := range req.Student {
			datacenterStudentInfoSlice = append(datacenterStudentInfoSlice, datacenter.StudentInfo{
				StudentId: v.StudentId,
				Relation:  v.Relation,
			})
			classInfo, err := l.svcCtx.ClassRpc.GetClassInfoById(context.Background(), &class.IdReq{Id: v.ClassId})
			if err != nil {
				continue
			}
			studentInfo, err := l.svcCtx.StudentRpc.FindOneById(context.Background(), &student.IdRequest{Id: v.StudentId})
			if err != nil {
				continue
			}
			_, err = l.svcCtx.StudentParentModel.InsertOne(&model.StudentParent{
				OrgId:       orgId,
				ClassId:     v.ClassId,
				ClassName:   classInfo.FullName,
				ParentId:    parentId,
				StudentId:   v.StudentId,
				StudentName: studentInfo.Name,
				Relation:    v.Relation,
			})
			if err != nil {
				continue
			}
		}
		s := datacenter.Marshal(datacenter.Parent, parentId, datacenter.Add, datacenter.ParentData{
			ParentName:  req.ParentName,
			UserName:    req.Moblie,
			Mobile:      req.Moblie,
			StudentInfo: datacenterStudentInfoSlice,
		})
		_ = l.svcCtx.KqPusher.Push(s)
	})
	return nil
}
