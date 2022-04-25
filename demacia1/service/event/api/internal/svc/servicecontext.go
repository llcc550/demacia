package svc

import (
	"demacia/datacenter/member/rpc/memberclient"
	"demacia/service/event/api/internal/config"
	"demacia/service/event/model"
	"demacia/service/position/positionclient"

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
