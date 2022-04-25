package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"demacia/datacenter/organization/rpc/organization"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type TreeClassLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeClassLogic(ctx context.Context, svcCtx *svc.ServiceContext) TreeClassLogic {
	return TreeClassLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TreeClassLogic) TreeClass() (resp *types.TreeClassResp, err error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	orgInfo, err := l.svcCtx.OrgRpc.FindOne(l.ctx, &organization.IdReply{Id: orgId})
	if err != nil {
		return nil, err
	}
	classList, _ := l.svcCtx.ClassModel.ListByOrgId(orgId)
	classTreeResp := map[int64][]types.TreeClassResp{}
	for _, class := range classList {
		classTreeResp[class.GradeId] = append(classTreeResp[class.GradeId], types.TreeClassResp{
			Id:       class.Id,
			Name:     class.FullName,
			Num:      class.ClassMemberNum,
			Children: []types.TreeClassResp{},
		})
	}

	gradeList, _ := l.svcCtx.GradeModel.GetGradeListByOrgId(orgId)
	gradeTreeResp := map[int64][]types.TreeClassResp{}
	classChildren := make([]types.TreeClassResp, 0)
	for _, grade := range *gradeList {
		if _, ok := classTreeResp[grade.Id]; ok {
			classChildren = classTreeResp[grade.Id]
		}
		gradeTreeResp[grade.StageId] = append(gradeTreeResp[grade.StageId], types.TreeClassResp{
			Id:       grade.Id,
			Name:     grade.Title,
			Num:      grade.GradeMemberNum,
			Children: classChildren,
		})
		classChildren = []types.TreeClassResp{}
	}
	stageList, _ := l.svcCtx.StageModel.ListByOrgId(orgId)
	stageTreeResp := map[int64][]types.TreeClassResp{}
	gradeChildren := make([]types.TreeClassResp, 0)
	var stageMemberNum int64
	for _, stage := range stageList {

		if _, ok := gradeTreeResp[stage.Id]; ok {
			gradeChildren = gradeTreeResp[stage.Id]
			for _, gradeV := range gradeTreeResp[stage.Id] {
				stageMemberNum = gradeV.Num + stageMemberNum
			}
		}

		stageTreeResp[stage.OrgId] = append(stageTreeResp[stage.OrgId], types.TreeClassResp{
			Id:       stage.Id,
			Name:     stage.Title,
			Num:      stageMemberNum,
			Children: gradeChildren,
		})
		gradeChildren = []types.TreeClassResp{}
		stageMemberNum = 0
	}
	var orgMemberNum int64
	stageChildren := make([]types.TreeClassResp, 0)
	if _, ok := stageTreeResp[orgId]; ok {
		stageChildren = stageTreeResp[orgId]
		for _, stageV := range stageTreeResp[orgId] {
			orgMemberNum = stageV.Num + orgMemberNum
		}
	}
	orgresp := &types.TreeClassResp{
		Id:       orgId,
		Name:     orgInfo.Title,
		Num:      orgMemberNum,
		Children: stageChildren,
	}
	stageChildren = []types.TreeClassResp{}
	orgMemberNum = 0
	return orgresp, nil
}
