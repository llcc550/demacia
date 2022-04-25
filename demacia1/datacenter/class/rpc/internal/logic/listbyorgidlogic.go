package logic

import (
	"context"

	"demacia/datacenter/class/rpc/class"
	"demacia/datacenter/class/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type ListByOrgIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListByOrgIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListByOrgIdLogic {
	return &ListByOrgIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListByOrgIdLogic) ListByOrgId(in *class.OrgIdReq) (*class.ListByOrgIdResp, error) {
	// todo: add your logic here and delete this line
	list, err := l.svcCtx.ClassModel.ListByOrgId(in.OrgId)
	if err != nil {
		return nil, err
	}
	resp := class.ListByOrgIdResp{}
	for _, v := range list {
		resp.List = append(resp.List, &class.ClassInfo{
			Id:       v.Id,
			StageId:  v.StageId,
			GradeId:  v.GradeId,
			FullName: v.FullName,
		})
	}
	return &resp, nil
}
