package logic

import (
	"context"
	"database/sql"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/class/model"
	"demacia/datacenter/member/rpc/member"
	"strconv"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddClassLogic {
	return AddClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddClassLogic) AddClass(req types.AddClassReq) (*types.Id, error) {

	orgId, err := baseauth.GetOrgId(l.ctx)

	title := "(" + strconv.Itoa(int(req.Sort)) + ")班"
	if err != nil {
		return nil, err
	}
	// 检验学段年级是否匹配
	StageGrade, err := l.svcCtx.GradeModel.GetGradeByIdAndStageId(req.GradeId, req.StageId)
	if err != nil || StageGrade == nil {
		// 学段年级不匹配
		return nil, errlist.StageGradeErr
	}
	FullName := StageGrade.StageTitle + StageGrade.Title + "(" + strconv.Itoa(int(req.Sort)) + ")班"
	//// 过滤班级全称(别名)
	//if strings.Trim(req.AliasName, " ") == "" {
	//	req.AliasName = StageGrade.StageTitle + StageGrade.Title + "(" + strconv.Itoa(int(req.Sort)) + ")班"
	//}

	// 检验班级名称是否存在
	checkClassName, err := l.svcCtx.ClassModel.GetClassByFullNameAndOrgId(orgId, FullName)
	if err != nil && err != sql.ErrNoRows {
		l.Logger.Errorf("Class select Class err:" + err.Error())
		return nil, errlist.Unknown
	}
	if checkClassName != nil {
		// 班级名称已存在
		return nil, errlist.ClassNameExist
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

	insertClassId, err := l.svcCtx.ClassModel.Insert(&model.Class{
		OrgId:          orgId,
		Title:          title,
		StageId:        StageGrade.StageId,
		StageTitle:     StageGrade.StageTitle,
		GradeId:        StageGrade.Id,
		GradeTitle:     StageGrade.Title,
		FullName:       FullName,
		AliasName:      req.AliasName,
		ClassMemberNum: 0,
		Desc:           req.Description,
		Sort:           req.Sort,
		Teachers:       teachers,
	})
	if err != nil {
		l.Logger.Errorf("Class Insert Class err:" + err.Error())
		return nil, errlist.Unknown
	}

	if len(req.TeacherId) > 0 {
		for _, v := range teacherMap {
			BatchTeacherSlices = append(BatchTeacherSlices, &model.ClassTeacher{
				ClassId:     insertClassId,
				TeacherId:   v.Id,
				TeacherName: v.TrueName,
			})
		}
		err = l.svcCtx.ClassTeacherModel.BatchInsert(BatchTeacherSlices)
		if err != nil {
			l.Logger.Errorf("Class batchInsert ClassTeach err:" + err.Error())
			goto Res
		}
	}

Res:
	return &types.Id{
		Id: insertClassId,
	}, nil
}
