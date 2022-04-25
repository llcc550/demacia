package logic

import (
	"context"
	"fmt"

	"demacia/common/baseauth"
	"demacia/datacenter/parent/api/internal/svc"
	"demacia/datacenter/parent/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ParentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewParentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ParentListLogic {
	return ParentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParentListLogic) ParentList(req types.ListConditionRequest) (resp *types.ParentList, err error) {
	res := &types.ParentList{
		List:  []types.List{},
		Count: 0,
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return res, err
	}
	ids := make([]int64, 0)
	studentParentMap := map[int64][]types.StudentInfo{}
	studentParentInfo, _, err := l.svcCtx.StudentParentModel.FindListByConditions(orgId, req.ClassId, req.StudentName, 0, 0)
	if err == nil {
		for _, studentParent := range studentParentInfo {
			ids = append(ids, studentParent.ParentId)
			studentParentMap[studentParent.ParentId] = append(studentParentMap[studentParent.ParentId], types.StudentInfo{
				StudentId:   studentParent.StudentId,
				StudentName: studentParent.StudentName,
				ClassName:   studentParent.ClassName,
				Relation:    studentParent.Relation,
			})
		}
	}

	parentInfo, parentCount, err := l.svcCtx.ParentModel.FindListByConditions(ids, req.ParentName, req.FaceStatus, req.Page, req.Limit)
	if err != nil || parentCount == 0 {
		return res, err
	}
	for _, parent := range parentInfo {
		studentInfo := make([]types.StudentInfo, 0)
		if _, ok := studentParentMap[parent.Id]; ok {
			studentInfo = studentParentMap[parent.Id]
		}
		res.List = append(res.List, types.List{
			ParentId:    parent.Id,
			ParentName:  parent.TrueName,
			Mobile:      parent.Mobile,
			StudentInfo: studentInfo,
			FaceStatus:  parent.FaceStatus,
		})
	}
	fmt.Println(parentCount)
	res.Count = parentCount
	return res, nil
}
