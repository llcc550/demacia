package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
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

type InsertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInsertLogic(ctx context.Context, svcCtx *svc.ServiceContext) InsertLogic {
	return InsertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InsertLogic) Insert(req types.InsertReq) (resp *types.MemberId, err error) {
	// 获取机构id
	Oid, GetOrgIdErr := baseauth.GetOrgId(l.ctx)
	if GetOrgIdErr != nil {
		return nil, errlist.InvalidParam
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

			if CheckErr == nil && i.Result {
				return nil, errors.CardsExitError
			}
		}
	}
	// 且手机号格式是否正确
	if !basefunc.CheckMobile(req.Mobile) {
		return nil, errors.MobileFormatError
	}
	// 验证手机号是否唯一
	_, err = l.svcCtx.MemberModel.FindOneByMobile(req.Mobile, Oid)
	if err == nil {
		return nil, errors.MobileExitError
	}

	// 个别字段不能为空
	if req.TrueName == "" || req.UserName == "" {
		return nil, errlist.InvalidParam
	}

	var Member = model.Member{
		OrgId:    Oid,
		TrueName: req.TrueName,
		UserName: req.UserName,
		Password: basefunc.HashPassword(req.Mobile),
		Mobile:   req.Mobile,
		Sex:      req.Sex,
		Avatar:   req.Avatar,
		Face:     req.Face,
	}

	// 添加成员
	memberId, err := l.svcCtx.MemberModel.Insert(&Member)
	if err != nil {
		return nil, err
	}
	res := types.MemberId{MemberId: memberId}

	threading.GoSafe(func() {
		if len(req.Cards) > 0 {
			_, _ = l.svcCtx.CardRpc.AddTeacherCard(context.Background(), &card.AddReq{
				OrgId:    Oid,
				ObjectId: memberId,
				CardNum:  req.Cards,
			})
		}
		_, _ = l.svcCtx.DataBusRpc.Create(context.Background(), &databus.Req{Topic: datacenter.Member, ObjectId: memberId})
	})
	return &res, nil
}
