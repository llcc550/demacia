package logic

import (
	"context"

	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/card/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type GetStudentCardListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStudentCardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStudentCardListLogic {
	return &GetStudentCardListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetStudentCardListLogic) GetStudentCardList(in *card.ListReq) (*card.ListResp, error) {
	list, err := l.svcCtx.CardModel.List(in.OrgId, in.ObjectId, roleStudent)
	if err != nil {
		return nil, err
	}
	resp := &card.ListResp{CardNum: []string{}}
	for _, i := range list {
		resp.CardNum = append(resp.CardNum, i.CardNum)
	}
	return resp, nil
}
