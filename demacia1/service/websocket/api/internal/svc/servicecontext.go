package svc

import (
	"demacia/service/websocket/api/internal/config"
	"demacia/service/websocket/utils"
)

type ServiceContext struct {
	Config config.Config
	Nodes  []utils.NodeInfo
}

func NewServiceContext(c config.Config, Nodes []utils.NodeInfo) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Nodes:  Nodes,
	}
}
