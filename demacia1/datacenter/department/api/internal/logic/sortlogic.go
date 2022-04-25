package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/errlist"
	"demacia/datacenter/department/api/internal/svc"
	"demacia/datacenter/department/api/internal/types"
	"demacia/datacenter/department/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type SortLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortLogic(ctx context.Context, svcCtx *svc.ServiceContext) SortLogic {
	return SortLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortLogic) Sort(req types.SortReq) error {
	if len(req.List) == 0 {
		return errlist.InvalidParam
	}
	orgId, err := baseauth.GetOrgId(l.ctx)
	if err != nil {
		return errlist.AuthLoginFail
	}
	list, err := l.svcCtx.DepartmentModel.GetDepartmentsByOrgId(orgId)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	departmentIdMap := map[int64]interface{}{}
	for _, item := range list {
		departmentIdMap[item.Id] = nil
	}
	reqDepartmentIdMap := map[int64]interface{}{}
	reqSortMap := map[int64]interface{}{}
	maxSort := 0
	sortData := make(model.Departments, 0, len(departmentIdMap))
	for _, item := range req.List {
		if item.Id <= 0 || item.Sort <= 0 {
			return errlist.InvalidParam
		}

		// 参数的departmentId不是本校的，参数错误
		if _, ok := departmentIdMap[item.Id]; !ok {
			return errlist.InvalidParam
		}

		// 参数中的departmentId重复，参数错误
		if _, ok := reqDepartmentIdMap[item.Id]; ok {
			return errlist.InvalidParam
		}
		reqDepartmentIdMap[item.Id] = nil

		// 参数中的sort重复，参数错误
		if _, ok := reqSortMap[item.Sort]; ok {
			return errlist.InvalidParam
		}
		reqSortMap[item.Sort] = nil

		// 记录最大的sort
		if int(item.Sort) > maxSort {
			maxSort = int(item.Sort)
		}

		//  一切正常时，model更新数据需要的参数
		sortData = append(sortData, &model.Department{
			Id:   item.Id,
			Sort: item.Sort,
		})
	}
	// 最后再校验下参数
	if maxSort != len(departmentIdMap) || len(reqDepartmentIdMap) != len(departmentIdMap) || len(reqSortMap) != len(departmentIdMap) {
		return errlist.InvalidParam
	}

	// 更新数据
	err = l.svcCtx.DepartmentModel.UpdateDepartmentSort(orgId, sortData)
	if err != nil {
		return err
	}
	return nil
}
