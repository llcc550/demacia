package logic

import (
	"context"

	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListGradeByOrgIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGradeByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGradeByOrgIdLogic {
	return &ListGradeByOrgIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListGradeByOrgIdLogic) ListGradeByOrgId(in *class.OrgIdReq) (*class.ListGradeResp, error) {
	// todo: add your logic here and delete this line
	list, err := l.svcCtx.GradeModel.GetGradeListByOrgId(in.OrgId)
	if err != nil {
		return nil, err
	}
	listResp := &class.ListGradeResp{}
	for _, v := range *list {
		listResp.List = append(listResp.List, &class.GradeInfo{
			Id:    v.Id,
			Title: v.Title,
		})
	}

	return listResp, nil
}
