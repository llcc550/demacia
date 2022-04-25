package logic

import (
	"context"
	"time"

	"demacia/common/basefunc"
	"demacia/datacenter/databus/model"
	"demacia/datacenter/databus/rpc/databus"
	"demacia/datacenter/databus/rpc/internal/svc"

	"github.com/globalsign/mgo/bson"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

type LogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogLogic {
	return &LogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogLogic) Log(in *databus.LogReq) (*databus.Res, error) {
	var jwt, req interface{}
	_ = bson.Unmarshal([]byte(basefunc.Json2Bson(in.Jwt)), &jwt)
	_ = bson.Unmarshal([]byte(basefunc.Json2Bson(in.Req)), &req)
	_ = l.svcCtx.LogModel.Insert(l.ctx, &model.Log{
		Ip:         in.Ip,
		Route:      in.Route,
		Jwt:        jwt,
		Req:        req,
		CreateTime: time.Now().Local(),
	})
	return &databus.Res{Result: true}, nil
}
