package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GradesByOrgIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGradesByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) GradesByOrgIdLogic {
	return GradesByOrgIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GradesByOrgIdLogic) GradesByOrgId() (resp *types.ListGradeByOrgResp, err error) {

	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.GradeModel.GetGradeListByOrgId(orgId)
	if err != nil {
		return nil, err
	}
	if len(*list) < 1 {
		return &types.ListGradeByOrgResp{}, nil
	}
	resp = &types.ListGradeByOrgResp{}
	gradeMap := map[int64][]*types.Info{}
	stageMap := map[int64]types.Info{}
	for _, v := range *list {
		gradeMap[v.StageId] = append(gradeMap[v.StageId], &types.Info{
			Id:    v.Id,
			Title: v.Title,
		})
		stageMap[v.StageId] = types.Info{
			Id:    v.StageId,
			Title: v.StageTitle,
		}
	}
	for gradeK, gradev := range gradeMap {
		if _, ok := stageMap[gradeK]; ok {
			resp.List = append(resp.List, &types.ListStageGradeInfo{
				Id:     stageMap[gradeK].Id,
				Title:  stageMap[gradeK].Title,
				Grades: gradev,
			})
		}
	}
	return resp, nil
}
