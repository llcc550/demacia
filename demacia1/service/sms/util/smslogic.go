package util

import (
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/limit"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

const (
	ProviderHuawei = "huawei"
	SendLimit      = 1000
)

func Allow(limiter *limit.PeriodLimit, mobile string) bool {
	code, err := limiter.Take(mobile)
	if err != nil {
		// we can't discard the message when the limit redis is out of service,
		// just let the message go
		logx.Error(err)
		return true
	}
	switch code {
	case limit.OverQuota:
		return false
	case limit.Allowed:
		return true
	case limit.HitQuota:
		// todo: maybe we need to let users know they hit the quota
		return false
	default:
		// unknown response, we just let the sms go
		return true
	}
}
