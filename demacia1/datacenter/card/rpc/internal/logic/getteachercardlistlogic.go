package logic

import (
	"context"

	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/card/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetTeacherCardListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTeacherCardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTeacherCardListLogic {
	return &GetTeacherCardListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTeacherCardListLogic) GetTeacherCardList(in *card.ListReq) (*card.ListResp, error) {
	list, err := l.svcCtx.CardModel.List(in.OrgId, in.ObjectId, roleTeacher)
	if err != nil {
		return nil, err
	}
	resp := &card.ListResp{CardNum: []string{}}
	for _, i := range list {
		resp.CardNum = append(resp.CardNum, i.CardNum)
	}
	return resp, nil
}
