package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/class/model"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"strconv"
)

type BatchAddClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchAddClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) BatchAddClassLogic {
	return BatchAddClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchAddClassLogic) BatchAddClass(req types.BatchAddClassReq) error {

	orgId, err := baseauth.GetOrgId(l.ctx)

	// 检验学段年级是否匹配
	StageGrade, err := l.svcCtx.GradeModel.GetGradeByIdAndStageId(req.GradeId, req.StageId)
	if err != nil || StageGrade == nil {
		// 学段年级不匹配
		return errlist.StageGradeErr
	}

	classList, err := l.svcCtx.ClassModel.ListByOrgIdAndGradeId(orgId, req.GradeId)
	noSort := make([]int8, 0)
	sortClassList := map[int8]interface{}{}
	maxClassSort := int8(0)

	for _, i := range classList {
		sortClassList[i.Sort] = nil
		if maxClassSort < i.Sort {
			maxClassSort = i.Sort
		}
	}
	sortList := make([]int8, 0)
	// 是否补全，不补全
	if req.Completion == 0 {
		for insNum := maxClassSort + 1; insNum <= maxClassSort+req.Num; insNum++ {
			sortList = append(sortList, insNum)
		}
		//直接插入数据库 sortList
	} else {
		// 补全
		for index := int8(1); index <= maxClassSort; index++ {
			if _, ok := sortClassList[index]; !ok {
				noSort = append(noSort, index)
			}
		}
		for i := int8(0); i < req.Num; i++ {
			if int(i) <= len(noSort)-1 {
				sortList = append(sortList, noSort[i])
			} else {
				maxClassSort++
				sortList = append(sortList, maxClassSort)
			}
		}
	}
	//data := make([]model.Class, 0)
	data := make(model.Classs, 0)
	for _, v := range sortList {
		ClassTitle := "(" + strconv.Itoa(int(v)) + ")班"
		data = append(data, &model.Class{
			OrgId:          orgId,
			Title:          ClassTitle,
			StageId:        StageGrade.StageId,
			StageTitle:     StageGrade.StageTitle,
			GradeId:        StageGrade.Id,
			GradeTitle:     StageGrade.Title,
			FullName:       StageGrade.StageTitle + StageGrade.Title + ClassTitle,
			ClassMemberNum: 0,
			Desc:           "",
			Sort:           v,
		})
	}
	err = l.svcCtx.ClassModel.BatchInsert(&data)
	if err != nil {
		return err
	}

	return nil
}
