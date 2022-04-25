package logic

import (
	"context"

	"demacia/common/basefunc"
	"demacia/datacenter/card/model"
	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/card/rpc/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type AddStudentCardLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddStudentCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddStudentCardLogic {
	return &AddStudentCardLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddStudentCardLogic) AddStudentCard(in *card.AddReq) (*card.BoolResp, error) {
	list, err := l.svcCtx.CardModel.List(in.OrgId, in.ObjectId, roleStudent)
	if err != nil {
		return nil, err
	}
	nowCardNums := make([]string, 0, len(list))
	for _, i := range list {
		nowCardNums = append(nowCardNums, i.CardNum)
	}
	if basefunc.StringSliceEq(nowCardNums, in.CardNum) {
		return &card.BoolResp{Result: true}, nil
	}
	data := make(model.Cards, 0, len(in.CardNum))
	for _, i := range in.CardNum {
		data = append(data, &model.Card{
			CardNum:    i,
			ObjectRole: roleStudent,
			ObjectId:   in.ObjectId,
			OrgId:      in.OrgId,
		})
	}
	_ = l.svcCtx.CardModel.Delete(in.OrgId, in.ObjectId, roleStudent)
	_ = l.svcCtx.CardModel.BatchInsert(data)
	_ = l.svcCtx.CardModel.RemoveCache(in.OrgId, in.ObjectId, roleStudent)
	return &card.BoolResp{Result: true}, nil
}
