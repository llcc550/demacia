package logic

import (
	"context"

	"demacia/datacenter/organization/api/internal/svc"
	"demacia/datacenter/organization/api/internal/types"
	"demacia/datacenter/organization/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) ListLogic {
	return ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req types.OrgListReq) (resp *types.OrgListResp, err error) {
	count, list, err := l.svcCtx.OrganizationModel.List(req.Page, req.Limit, &model.ListCond{
		Title:      req.Title,
		ProvinceId: req.ProvinceId,
		CityId:     req.CityId,
		AreaId:     req.AreaId,
	})
	if err != nil {
		logx.Errorf("get organization list error. req is %v, error is %s", req, err.Error())
		return nil, err
	}
	if count == 0 || len(list) == 0 {
		return &types.OrgListResp{Count: count, List: []*types.OrgInfo{}}, nil
	}
	respList := make([]*types.OrgInfo, 0, req.Limit)
	for _, i := range list {
		respList = append(respList, &types.OrgInfo{
			Id:           i.Id,
			AreaTitle:    i.AreaTitle,
			Title:        i.Title,
			ActivateDate: i.ActivateDate,
			ExpireDate:   i.ExpireDate,
			TrueName:     i.TrueName,
			Mobile:       i.Mobile,
			UserName:     i.ManagerMemberUserName,
			Addr:         i.Addr,
			Msg:          i.Msg,
			CreatedTime:  i.CreateTime,
			OrgStatus:    i.OrgStatus,
		})
	}
	return &types.OrgListResp{
		Count: count,
		List:  respList,
	}, nil
}
