package logic

import (
	"context"
	"demacia/common/errlist"
	"demacia/service/sms/model"
	"demacia/service/sms/rpc/internal/svc"
	"demacia/service/sms/rpc/sms"
	"demacia/service/sms/util"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
	"strings"
)

type MultiSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMultiSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiSendLogic {
	return &MultiSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MultiSendLogic) MultiSend(in *sms.MultiSendRequest) (*sms.SmsSendResponse, error) {
	multiRequest := new(MultiBatchSendSmsRequest)
	multiRequest.TemplateId = in.TemplateId
	for _, v := range in.Multi {
		multiRequest.MultiMt = append(multiRequest.MultiMt, &model.ContentMobile{
			Mobile:     v.Mobile,
			Params:     v.Params,
			TemplateId: v.TemplateId,
		})
	}

	err := l.multiBatchSendSms(multiRequest)
	if err != nil {
		return nil, err
	}
	return &sms.SmsSendResponse{}, nil
}

func (l *MultiSendLogic) multiBatchSendSms(req *MultiBatchSendSmsRequest) error {
	if !l.svcCtx.Push {
		return nil
	}
	if len(req.MultiMt) == 0 {
		return errlist.InvalidParam
	}
	var mobiles []*model.ContentMobile
	for _, v := range req.MultiMt {
		for _, s := range v.Mobile {
			if len(s) != 11 {
				logx.Errorf("illegal mobile, mobile=%v", s)
				continue
			}

			if !util.Allow(l.svcCtx.Limiter, s) {
				logx.Errorf("out of quota, mobile=%v", s)
				continue
			}
			mobiles = append(mobiles, &model.ContentMobile{
				Mobile:     v.Mobile,
				TemplateId: v.TemplateId,
				Params:     v.Params,
			})
		}
	}
	if len(mobiles) == 0 {
		return errlist.InvalidMobilesErr
	}
	threading.GoSafe(func() {
		for {
			if len(mobiles) > util.SendLimit {
				to := mobiles[:util.SendLimit]
				err := l.svcCtx.HuaweiModel.MultiSendBatchSms(req.TemplateId, to)
				errorString := ""
				if err != nil {
					errorString = err.Error()
				}
				trafficData := make(model.SmsTrafficList, 0, len(to))
				for _, mobile := range to {
					trafficData = append(trafficData, &model.SmsTraffic{
						Mobile:      strings.Join(mobile.Mobile, ";"),
						ContentType: req.TemplateId,
						Content:     strings.Join(mobile.Params, ";"),
						Provider:    util.ProviderHuawei,
						Error:       errorString,
					})
				}
				_ = l.svcCtx.TrafficModel.InsertBatch(&trafficData)
				mobiles = mobiles[util.SendLimit:]
			} else {
				err := l.svcCtx.HuaweiModel.MultiSendBatchSms(req.TemplateId, mobiles)
				errorString := ""
				if err != nil {
					errorString = err.Error()
				}
				trafficData := make(model.SmsTrafficList, 0, len(mobiles))
				for _, mobile := range mobiles {
					trafficData = append(trafficData, &model.SmsTraffic{
						Mobile:      strings.Join(mobile.Mobile, ";"),
						ContentType: req.TemplateId,
						Content:     strings.Join(mobile.Params, ";"),
						Provider:    util.ProviderHuawei,
						Error:       errorString,
					})
				}
				_ = l.svcCtx.TrafficModel.InsertBatch(&trafficData)
				return
			}
		}
	})
	return nil
}
