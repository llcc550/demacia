package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/cachemodel"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/class/model"
	"demacia/datacenter/member/rpc/member"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type UpdateClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateClassLogic {
	return UpdateClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateClassLogic) UpdateClass(req types.AddClassReq) error {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return err
	}
	if req.Id <= 0 {
		// 班级不存在
		return errlist.InvalidParam
	}

	var teachers string
	teacherMap := make([]types.Teacher, 0)
	BatchTeacherSlices := make(model.ClassTeachers, 0)
	if len(req.TeacherId) > 0 {
		for _, rt := range req.TeacherId {
			teacher, err := l.svcCtx.MemberRpc.FindOneById(context.Background(), &member.IdReq{Id: rt})
			if err != nil {
				continue
			}
			teacherMap = append(teacherMap, types.Teacher{
				Id:       teacher.Id,
				TrueName: teacher.TrueName,
			})
		}
		num := 0
		if len(teacherMap) > 0 {
			for _, tm := range teacherMap {
				num++
				if num == len(teacherMap) {
					teachers += tm.TrueName
				} else {
					teachers += tm.TrueName + ","
				}
			}
		}
	}

	classId := req.Id
	err = l.svcCtx.ClassModel.Update(&model.Class{
		Id:        classId,
		OrgId:     orgId,
		AliasName: req.AliasName,
		Desc:      req.Description,
		Teachers:  teachers,
	})
	if err != nil {
		if err == cachemodel.ErrNotFound {
			return errlist.ClassNotFound
		}
		l.Logger.Errorf("Class Update Class err:" + err.Error())
		return errlist.Unknown
	}
	err = l.svcCtx.ClassTeacherModel.DeletedByClassId(classId)
	if err != nil {
		l.Logger.Errorf("Class Delete ClassTeacher err:" + err.Error())
		return errlist.Unknown
	}
	for _, v := range teacherMap {
		BatchTeacherSlices = append(BatchTeacherSlices, &model.ClassTeacher{
			ClassId:     classId,
			TeacherId:   v.Id,
			TeacherName: v.TrueName,
		})
	}

	err = l.svcCtx.ClassTeacherModel.BatchInsert(BatchTeacherSlices)
	if err != nil {
		l.Logger.Errorf("Class batchInsert ClassTeacher err:" + err.Error())
		return errlist.Unknown
	}
	return nil
}
