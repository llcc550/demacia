package svc

import (
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/service/position/positionclient"
	"demacia/service/urgentevent/api/internal/config"
	"demacia/service/urgentevent/model"

	"gitlab.u-jy.cn/xiaoyang/go-zero/core/stores/postgres"
	"gitlab.u-jy.cn/xiaoyang/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	EventModel         *model.EventModel
	CategoryModel      *model.CategoryModel
	EventPositionModel *model.EventPositionModel
	MemberRpc          memberclient.Member
	PositionRpc        positionclient.Position
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Postgres.DataSource)
	cacheRedis := c.CacheRedis.NewRedis()
	return &ServiceContext{
		Config:             c,
		EventModel:         model.NewEventModel(conn, cacheRedis),
		CategoryModel:      model.NewCategoryModel(conn, cacheRedis),
		EventPositionModel: model.NewEventPositionModel(conn, cacheRedis),
		MemberRpc:          memberclient.NewMember(zrpc.MustNewClient(c.MemberRpc)),
		PositionRpc:        positionclient.NewPosition(zrpc.MustNewClient(c.PositionRpc)),
	}
}
