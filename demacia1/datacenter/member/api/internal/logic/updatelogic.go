package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/datacenter"
	"demacia/common/errlist"
	"demacia/datacenter/card/rpc/card"
	"demacia/datacenter/card/rpc/cardclient"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/member/api/internal/errors"
	"demacia/datacenter/member/api/internal/svc"
	"demacia/datacenter/member/api/internal/types"
	"demacia/datacenter/member/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateLogic {
	return UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req types.UpdateReq) error {
	// 获取机构id
	Oid, GetOrgIdErr := baseauth.GetOrgId(l.ctx)
	if GetOrgIdErr != nil {
		return errlist.InvalidParam
	}
	// 验证卡号
	if len(req.Cards) > 0 {
		// 验证卡号是否被使用
		for _, v := range req.Cards {
			i, CheckErr := l.svcCtx.CardRpc.CheckMember(l.ctx, &cardclient.CheckReq{
				OrgId:    Oid,
				ObjectId: 0,
				CardNum:  v,
			})
			if CheckErr != nil || i.Result == true {
				return errors.CardsExitError
			}
		}
	}

	err := l.svcCtx.MemberModel.UpdateMemberInfo(req.MemberId, model.Member{
		OrgId:    Oid,
		TrueName: req.TrueName,
		UserName: req.UserName,
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		Avatar:   req.Avatar,
		Face:     req.Face,
	})

	if err != nil {
		return err
	}
	// 修改卡号
	if len(req.Cards) > 0 {
		_, _ = l.svcCtx.CardRpc.AddTeacherCard(l.ctx, &card.AddReq{
			OrgId:    Oid,
			ObjectId: req.MemberId,
			CardNum:  req.Cards,
		})
	}

	threading.GoSafe(func() {
		_, _ = l.svcCtx.DataBusRpc.Update(context.Background(), &databus.Req{Topic: datacenter.Member, ObjectId: req.MemberId})
	})
	return nil
}
