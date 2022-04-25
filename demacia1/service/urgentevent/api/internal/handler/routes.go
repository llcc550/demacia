// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"demacia/service/urgentevent/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/category/insert",
				Handler: CategoryInsertHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/category/update",
				Handler: CategoryUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/category/delete",
				Handler: CategoryDeleteHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/category/list",
				Handler: CategoryListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/insert",
				Handler: EventInsertHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/update",
				Handler: EventUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/delete",
				Handler: EventDeleteHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/detail",
				Handler: EventDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/urgentevent/list",
				Handler: EventListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
