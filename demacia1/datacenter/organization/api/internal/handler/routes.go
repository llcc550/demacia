// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"demacia/datacenter/organization/api/internal/svc"

	"gitlab.u-jy.cn/xiaoyang/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Log},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/api/organization/add",
					Handler: addHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/organization/update",
					Handler: updateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/organization/del",
					Handler: delHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/organization/list",
					Handler: listHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/api/organization/detail",
					Handler: detailHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}
