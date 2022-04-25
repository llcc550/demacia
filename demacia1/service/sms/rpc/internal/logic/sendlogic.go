package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/service/sms/model"
	"demacia/service/sms/util"
	"strings"

	"demacia/service/sms/rpc/internal/svc"
	"demacia/service/sms/rpc/sms"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type SendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}
type (
	AdminSendSmsRequest struct {
		Mobile     string   `path:"mobile"`
		Params     []string `json:"params"`
		TemplateId string   `json:"template_id"`
	}

	MultiBatchSendSmsRequest struct {
		TemplateId string                 `json:"template_id"`
		MultiMt    []*model.ContentMobile `json:"multimt"`
	}
)

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendLogic) Send(in *sms.SendRequest) (*sms.SmsSendResponse, error) {
	err := l.adminSendSms(&AdminSendSmsRequest{
		Mobile:     in.Mobile,
		Params:     in.Params,
		TemplateId: in.TemplateId,
	})
	if err != nil {
		return nil, err
	}

	return &sms.SmsSendResponse{}, nil
}

func (l *SendLogic) adminSendSms(req *AdminSendSmsRequest) error {
	if !l.svcCtx.Push {
		return nil
	}
	if len(req.Mobile) != 11 {
		return errlist.InvalidParam
	}
	if !util.Allow(l.svcCtx.Limiter, req.Mobile) {
		return errlist.OutOfQuotaErr
	}
	err := l.svcCtx.HuaweiModel.Send(req.Mobile,
		req.TemplateId, req.Params,
	)
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	_ = l.svcCtx.TrafficModel.Insert(&model.SmsTraffic{
		Mobile:      req.Mobile,
		ContentType: req.TemplateId,
		Content:     strings.Join(req.Params, ","),
		Provider:    util.ProviderHuawei,
		Error:       errorString,
	})
	if err != nil {
		return errlist.SendSmsErr
	}
	return nil
}
