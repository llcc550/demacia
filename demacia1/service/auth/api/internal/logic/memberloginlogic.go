package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/errlist"
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/service/auth/api/internal/svc"
	"demacia/service/auth/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type MemberLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMemberLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) MemberLoginLogic {
	return MemberLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberLoginLogic) MemberLogin(req types.LoginRequest) (resp *types.TokenResponse, err error) {
	if req.UserName == "" || req.Password == "" {
		return nil, errlist.InvalidParam
	}
	tokenMap, err := l.login(req)
	if err != nil {
		// todo: 限制登录失败
		return nil, err
	}

	token, err := buildTokens(l.svcCtx.Config.Auth, tokenMap)
	if err != nil {
		return nil, errlist.Unknown
	}
	return &types.TokenResponse{Token: token}, nil
}

func (l *MemberLoginLogic) login(req types.LoginRequest) (res map[string]interface{}, resErr error) {
	resErr = errlist.AuthLoginFail
	memberInfo, err := l.svcCtx.MemberRpc.FindOneByUserName(l.ctx, &memberclient.UserNameReq{UserName: req.UserName})
	if err != nil {
		return
	}
	if !basefunc.CheckPassword(memberInfo.Password, req.Password) {
		return
	}
	return map[string]interface{}{
		baseauth.OrgIdField:    memberInfo.OrgId,
		baseauth.MemberIdField: memberInfo.Id,
		baseauth.RoleField:     1,
	}, nil
}
