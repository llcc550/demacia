package logic

import (
	"context"
	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/class/api/internal/svc"
	"demacia/datacenter/class/api/internal/types"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"strconv"
)

type ClassCompletionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClassCompletionLogic(ctx context.Context, svcCtx *svc.ServiceContext) ClassCompletionLogic {
	return ClassCompletionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClassCompletionLogic) ClassCompletion(req types.ClassCompletionReq) (*types.ListRespose, error) {

	// 检验学段年级是否匹配
	orgId, err := baseauth.GetOrgId(l.ctx)

	if err != nil {
		return nil, err
	}
	StageGrade, err := l.svcCtx.GradeModel.GetGradeByIdAndStageId(req.GradeId, req.StageId)
	if err != nil || StageGrade == nil {
		// 学段年级不匹配
		return nil, errlist.StageGradeErr
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
	for index := int8(1); index <= maxClassSort; index++ {
		if _, ok := sortClassList[index]; !ok {
			noSort = append(noSort, index)
		}
	}
	if req.Num > 0 {
		noSortLen := int8(len(noSort))
		if req.Num == 1 && noSortLen >= 1 {
			noSort = append(noSort, maxClassSort+1)
		}
		if req.Num > noSortLen {
			for insNum := maxClassSort + 1; insNum <= maxClassSort+(req.Num-noSortLen); insNum++ {
				noSort = append(noSort, insNum)
			}
		}
	}
	resp := &types.ListRespose{}
	for _, n := range noSort {
		resp.List = append(resp.List, types.Info{
			Id:    int64(n),
			Title: StageGrade.StageTitle + StageGrade.Title + "(" + strconv.Itoa(int(n)) + ")班",
		})
	}
	return resp, nil
}
