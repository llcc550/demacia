package logic

import (
	"encoding/json"
	"strings"

	"demacia/common/errlist"
	"demacia/service/sms/model"
	"demacia/service/sms/util"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/threading"
)

type (
	RmqData struct {
		Type int8   `json:"type"`
		Data string `json:"data"`
	}

	MultiBatchSendSmsRequest struct {
		TemplateId string                 `json:"template_id"`
		MultiMt    []*model.ContentMobile `json:"multimt"`
	}
)

func (l *Consumer) Consumer(body string) error {
	if !l.svcCtx.Push {
		return nil
	}
	logx.Infof("rmq body: %s", body)
	var req RmqData
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return err
	}
	var data MultiBatchSendSmsRequest
	if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
		logx.Errorf("data unmarshal %s, failed: %s", body, err.Error())
		return err
	}
	err := l.MultiBatchSendSms(&data)
	return err
}

func (l *Consumer) MultiBatchSendSms(req *MultiBatchSendSmsRequest) error {
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
