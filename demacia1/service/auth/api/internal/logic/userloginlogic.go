package logic

import (
	"context"

	"demacia/common/baseauth"
	"demacia/common/basefunc"
	"demacia/common/errlist"
	"demacia/datacenter/user/rpc/userclient"
	"demacia/service/auth/api/internal/svc"
	"demacia/service/auth/api/internal/types"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type UserLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserLoginLogic {
	return UserLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLoginLogic) UserLogin(req types.LoginRequest) (resp *types.TokenResponse, err error) {
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

func (l *UserLoginLogic) login(req types.LoginRequest) (res map[string]interface{}, resErr error) {
	resErr = errlist.AuthLoginFail
	userInfo, err := l.svcCtx.UserRpc.FindOneByUserName(l.ctx, &userclient.UserNameReq{UserName: req.UserName})
	if err != nil {
		return
	}
	if !basefunc.CheckPassword(userInfo.Password, req.Password) {
		return
	}
	return map[string]interface{}{
		baseauth.UserIdField: userInfo.Id,
		baseauth.RoleField:   1,
	}, nil
}
